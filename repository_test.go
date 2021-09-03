package gitconfig

import (
	"path/filepath"
	"testing"

	testspace "github.com/Jiu2015/gotestspace"
	"github.com/stretchr/testify/assert"
)

var repoTestSpace testspace.Space

func initRepoTestSpace(t *testing.T) {
	var err error

	repoTestSpace, err = testspace.Create(
		testspace.WithPathOption("testspace-*"),
		testspace.WithShellOption(`
			git config --global core.abbrev 10 &&
			git config --global init.defaultBranch main &&
			git init --bare repo.git &&
			git clone repo.git workdir &&
			(
				cd workdir &&
				mkdir a &&
				printf "A\n" >a/a.txt &&
				git add a &&
				test_tick &&
				git commit -m "A" &&
				mkdir a/b &&
				printf "B\n" >a/b/b.txt &&
				git add a/b &&
				test_tick &&
				git commit -m "B" &&
				mkdir a/b/c &&
				printf "C\n" >a/b/c/c.txt &&
				git add a/b/c &&
				test_tick &&
				git commit -m "C" &&
				git push -u origin HEAD
			) &&
			(
				cd workdir &&
				git worktree add -b topic1 ../topic1 &&
				cd ../topic1 &&
				printf "topic1\n" >topic1.txt &&
				git add topic1.txt &&
				test_tick &&
				git commit -m "topic1"
			) &&
			(
				cd workdir &&
				git worktree add -b topic2 ../topic2
			)
		`),
	)
	assert.Nil(t, err)
}
func TestWithSubtest(t *testing.T) {
	// Setup
	initRepoTestSpace(t)

	t.Run("T=1", subTestRepositoryIsBare)
	t.Run("T=2", subTestRepositoryGitPath)

	// Tear-down
	repoTestSpace.Cleanup()
}

func subTestRepositoryIsBare(t *testing.T) {
	var (
		repo *Repository
		err  error
	)
	for _, tc := range []struct {
		Path   string
		IsBare bool
	}{
		{"workdir", false},
		{"workdir/a/b/c", false},
		{"workdir/a/b/c/non-exist/dir", false},
		{"repo.git", true},
		{"repo.git/refs", true},
		{"repo.git/refs/non-exist/dir", true},
		{"topic1/a/b/c", false},
		{"topic2/a/b/c", false},
	} {
		repo, err = FindRepository(repoTestSpace.GetPath(tc.Path))
		if assert.Nil(t, err) {
			assert.Equal(t,
				tc.IsBare,
				repo.IsBare(),
				"Repository at '%s' is a bare? expect: %v, actual: %v",
				tc.Path,
				tc.IsBare,
				repo.IsBare(),
			)

			assert.Equal(t,
				tc.IsBare,
				repo.Config().GetBool("core.bare", true),
				"Repository at '%s' is a bare? expect: %v, config: %v",
				tc.Path,
				tc.IsBare,
				repo.Config().GetBool("core.bare", true),
			)

		}
	}
}

func subTestRepositoryGitPath(t *testing.T) {
	var (
		repo *Repository
		err  error
	)

	baseDir := repoTestSpace.GetPath("")
	baseDir, _ = filepath.EvalSymlinks(baseDir)
	getRelDir := func(dir string) string {
		dir, _ = filepath.EvalSymlinks(dir)
		dir, _ = filepath.Rel(baseDir, dir)
		return dir
	}
	for _, tc := range []struct {
		Path      string
		GitDir    string
		CommonDir string
		WorkDir   string
	}{
		{"workdir", "workdir/.git", "workdir/.git", "workdir"},
		{"workdir/a/b/c", "workdir/.git", "workdir/.git", "workdir"},
		{"workdir/a/b/c/non-exist/dir", "workdir/.git", "workdir/.git", "workdir"},
		{"repo.git", "repo.git", "repo.git", ""},
		{"repo.git/refs", "repo.git", "repo.git", ""},
		{"repo.git/refs/non-exist/dir", "repo.git", "repo.git", ""},
		{"topic1/a/b/c", "workdir/.git/worktrees/topic1", "workdir/.git", "topic1"},
		{"topic2/a", "workdir/.git/worktrees/topic2", "workdir/.git", "topic2"},
	} {
		repo, err = FindRepository(repoTestSpace.GetPath(tc.Path))
		if assert.Nil(t, err) {
			assert.Equal(t,
				tc.GitDir,
				getRelDir(repo.GitDir()),
			)

			assert.Equal(t,
				tc.CommonDir,
				getRelDir(repo.GitCommonDir()),
			)

			assert.Equal(t,
				tc.WorkDir,
				getRelDir(repo.WorkDir()),
			)
		}

	}
}
