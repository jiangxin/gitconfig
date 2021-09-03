package gitconfig

import (
	"context"
	"os"
	"testing"

	testspace "github.com/Jiu2015/gotestspace"
	"github.com/stretchr/testify/assert"
)

var cacheTestSpace testspace.Space

func TestCacheCreateConfig(t *testing.T) {
	var err error

	cacheTestSpace, err = testspace.Create(
		testspace.WithPathOption("testspace-*"),
		testspace.WithShellOption(`
			git config -f config.test a.b "value of a.b"
			git config -f config.test a.b.c "value of a/b/c"
		`),
	)
	assert.Nil(t, err)
}

func TestCacheSet(t *testing.T) {
	fn := cacheTestSpace.GetPath("config.test")
	fi, err := os.Stat(fn)
	if assert.Nil(t, err) {
		cfg, err := loadConfigFile(fn)
		if assert.Nil(t, err) {
			CacheSet(fn, cfg, fi.Size(), fi.ModTime())
		}
	}
}

func TestCacheGet(t *testing.T) {
	fn := cacheTestSpace.GetPath("config.test")
	cfg, ok := CacheGet(fn)
	if assert.True(t, ok) {
		for _, tc := range []struct {
			Key   string
			Value string
		}{
			{"a.b", "value of a.b"},
			{"a.b.c", "value of a/b/c"},
		} {
			assert.Equal(t,
				tc.Value,
				cfg.Get(tc.Key),
			)
		}
	}
}

func TestCacheUpdateConfig(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stdout, stderr, err := cacheTestSpace.Execute(
		cancelCtx,
		`git config -f config.test foo.bar "foo bar"`,
	)
	assert.Nil(t, err, "stdout: %s\nstderr: %s", stdout, stderr)
}

func TestCacheOutOfDate(t *testing.T) {
	fn := cacheTestSpace.GetPath("config.test")
	cfg, ok := CacheGet(fn)
	assert.False(t, ok)
	assert.Nil(t, cfg)
}

func TestCacheSetAgain(t *testing.T) {
	fn := cacheTestSpace.GetPath("config.test")
	fi, err := os.Stat(fn)
	if assert.Nil(t, err) {
		cfg, err := loadConfigFile(fn)
		if assert.Nil(t, err) {
			CacheSet(fn, cfg, fi.Size(), fi.ModTime())
		}
	}
	cfg, ok := CacheGet(fn)
	if assert.True(t, ok) {
		for _, tc := range []struct {
			Key   string
			Value string
		}{
			{"a.b", "value of a.b"},
			{"a.b.c", "value of a/b/c"},
			{"foo.bar", "foo bar"},
		} {
			assert.Equal(t,
				tc.Value,
				cfg.Get(tc.Key),
			)
		}

	}
	cfg, ok = CacheGet(fn)
	assert.True(t, ok)
	assert.NotNil(t, cfg)
}
