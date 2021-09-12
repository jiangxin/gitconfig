package gitconfig

import "path/filepath"

// Repository defines struct for a Git repository.
type Repository struct {
	gitDir       string
	gitCommonDir string
	workDir      string
	gitConfig    GitConfig
}

// GitDir returns GitDir
func (v Repository) GitDir() string {
	return v.gitDir
}

// GitCommonDir returns commondir where contains "config" file
func (v Repository) GitCommonDir() string {
	return v.gitCommonDir
}

// WorkDir returns workdir
func (v Repository) WorkDir() string {
	return v.workDir
}

// IsBare indicates a repository is a bare repository.
func (v Repository) IsBare() bool {
	return v.workDir == ""
}

// Config returns git config object
func (v Repository) Config() GitConfig {
	return v.gitConfig
}

// FindRepository locates repository object search from the given dir.
func FindRepository(dir string) (*Repository, error) {
	var (
		gitDir    string
		commonDir string
		workDir   string
		gitConfig GitConfig
		err       error
	)

	gitDir, err = findGitDir(dir)
	if err != nil {
		return nil, err
	}
	commonDir, err = getGitCommonDir(gitDir)
	if err != nil {
		return nil, err
	}
	gitConfig, err = LoadFileWithDefault(filepath.Join(commonDir, "config"))
	if err != nil {
		return nil, err
	}
	if !gitConfig.GetBool("core.bare", false) {
		workDir, _ = getWorkTree(gitDir)
	}
	return &Repository{
		gitDir:       gitDir,
		gitCommonDir: commonDir,
		workDir:      workDir,
		gitConfig:    gitConfig,
	}, nil
}
