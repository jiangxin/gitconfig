package gitconfig

import (
	"io/ioutil"
	"os"
)

// LoadFile loads specific git config file.
func LoadFile(name string) (GitConfig, error) {
	if cfg, ok := CacheGet(name); ok {
		return cfg, nil
	}

	// cache will be updated using this time
	fi, err := os.Stat(name)
	if err != nil {
		return nil, ErrNotExist
	}

	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	cfg, _, err := Parse(buf, name)
	if err != nil {
		return cfg, err
	}

	// update cache
	CacheSet(name, cfg, fi.Size(), fi.ModTime())
	return cfg, nil
}

// LoadFileWithDefault loads specific git config file and fallback
// to default config (user level config or system level).
func LoadFileWithDefault(name string) (GitConfig, error) {
	cfg := DefaultConfig()

	repoConfig, err := LoadFile(name)
	if err == nil {
		cfg.Merge(repoConfig, ScopeSelf)
	}
	return cfg, nil
}

// LoadDir only loads git config file found in gitdir.
func LoadDir(dir string) (GitConfig, error) {
	var (
		err error
	)

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	configFile, err := FindGitConfig(dir)
	if err != nil {
		return nil, ErrNotExist
	}

	return LoadFile(configFile)
}

// LoadDirWithDefault loads git config file found in gitdir, and
// fallback to default (global and system level git config).
func LoadDirWithDefault(dir string) (GitConfig, error) {
	cfg := DefaultConfig()

	repoConfig, err := LoadDir(dir)
	if err == nil {
		cfg.Merge(repoConfig, ScopeSelf)
	}
	return cfg, nil
}

// SystemConfig returns system git config, reload if necessary
func SystemConfig() (GitConfig, error) {
	file := SystemConfigFile()
	if file == "" {
		return nil, nil
	}

	if _, err := os.Stat(file); err != nil {
		return nil, nil
	}
	return LoadFile(file)
}

// GlobalConfig returns global user config, reload if necessary
func GlobalConfig() (GitConfig, error) {
	file, err := GlobalConfigFile()
	if err != nil {
		return nil, nil
	}

	if _, err := os.Stat(file); err != nil {
		return nil, nil
	}
	return LoadFile(file)
}

// DefaultConfig returns global and system wide config
func DefaultConfig() GitConfig {
	cfg := NewGitConfig()
	if sysCfg, err := SystemConfig(); err == nil && sysCfg != nil {
		cfg.Merge(sysCfg, ScopeSystem)
	}
	if globalCfg, err := GlobalConfig(); err == nil && globalCfg != nil {
		cfg.Merge(globalCfg, ScopeGlobal)
	}
	return cfg
}
