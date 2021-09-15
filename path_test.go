package gitconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpendHome(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "gitconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	UnsetHome()
	name, err = HomeDir()
	assert.NotNil(err)
	assert.Equal("", name)

	name, err = ExpendHome("")
	assert.NotNil(err)
	assert.Equal("", name)

	SetHome(tmpdir)

	name, err = HomeDir()
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = ExpendHome("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "~a"), name)

	name, err = ExpendHome("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = ExpendHome("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = ExpendHome(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = ExpendHome(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

}

func TestAbs(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "gitconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	UnsetHome()
	name, err = Abs("~/")
	assert.NotNil(err)
	assert.Equal("", name)

	SetHome(tmpdir)
	cwd, err := os.Getwd()
	assert.Nil(err)

	name, err = Abs("")
	assert.Nil(err, fmt.Sprintf("err should be nil, but got: %s", err))
	assert.Equal(cwd, name)

	name, err = Abs("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = Abs("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = Abs("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = Abs("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = Abs("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = Abs("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = Abs(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = Abs(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)
}

func TestAbsJoin(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "gitconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	SetHome(tmpdir)

	cwd := "/some/dir"
	if runtime.GOOS == "windows" {
		cwd = "c:\\some\\dir"
	}

	name, err = AbsJoin(cwd, "")
	assert.Nil(err)
	assert.Equal(cwd, name)

	name, err = AbsJoin(cwd, "a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = AbsJoin(cwd, "~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = AbsJoin(cwd, "~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = AbsJoin(cwd, "~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = AbsJoin(cwd, "~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = AbsJoin(cwd, "ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = AbsJoin(cwd, inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = AbsJoin(cwd, inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)
}

func TestFindGitDir(t *testing.T) {
	var (
		err     error
		dir     string
		gitdir  string
		workdir string
		cfg     string
		home    string
		assert  = assert.New(t)
	)

	tmpdir, err := ioutil.TempDir("", "gitconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	SetHome(tmpdir)

	// find in: bare.git
	gitdir = filepath.Join(tmpdir, "bare.git")
	cmd := exec.Command("git", "init", "--bare", gitdir, "--")
	assert.Nil(cmd.Run())
	dir, err = FindGitDir(gitdir)
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(gitdir)
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// find in: bare.git/objects/pack
	dir, err = FindGitDir(filepath.Join(gitdir, "objects", "pack"))
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(gitdir, "objects", "pack"))
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// create repo2 with gitdir file repo2/.git
	repo2 := filepath.Join(tmpdir, "repo2")
	err = os.MkdirAll(filepath.Join(repo2, "a", "b"), 0755)
	assert.Equal(nil, err)
	err = ioutil.WriteFile(filepath.Join(repo2, ".git"),
		[]byte("gitdir: ../bare.git"),
		0644)
	assert.Nil(err)

	// find in: repo2/a/b/c
	dir, err = FindGitDir(filepath.Join(repo2, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(repo2, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// create bad gitdir file: repo2.git
	err = ioutil.WriteFile(filepath.Join(repo2, ".git"),
		[]byte("../bare.git"),
		0644)
	assert.Nil(err)

	// fail to find in repo2/a/b/c (bad gitdir file)
	dir, err = FindGitDir(filepath.Join(repo2, "a", "b", "c"))
	assert.NotNil(err)
	assert.Equal("", dir)

	cfg, err = FindGitConfig(filepath.Join(repo2, "a", "b", "c"))
	assert.NotNil(err)
	assert.Equal("", cfg)

	// create worktree
	workdir = filepath.Join(tmpdir, "workdir")
	cmd = exec.Command("git", "init", workdir, "--")
	assert.Nil(cmd.Run())

	gitdir = filepath.Join(workdir, ".git")
	err = os.MkdirAll(filepath.Join(workdir, "a", "b"), 0755)
	assert.Nil(err)

	// find in workdir
	dir, err = FindGitDir(workdir)
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(workdir)
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// find in workdir/.git
	dir, err = FindGitDir(gitdir)
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(gitdir)
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// find in workdir/.git
	dir, err = FindGitDir(filepath.Join(workdir, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(gitdir, dir)

	cfg, err = FindGitConfig(filepath.Join(workdir, "a", "b", "c"))
	assert.Nil(err)
	assert.Equal(filepath.Join(gitdir, "config"), cfg)

	// fail to find in tmpdir
	dir, err = FindGitDir(tmpdir)
	assert.Equal("", dir)
	assert.Equal(ErrNotInGitDir, err)

	cfg, err = FindGitConfig(tmpdir)
	assert.Equal(ErrNotInGitDir, err)
	assert.Equal("", cfg)
}
