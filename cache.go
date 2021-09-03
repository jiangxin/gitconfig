package gitconfig

import (
	"os"
	"time"

	"github.com/golang/groupcache/lru"
)

var cache *lru.Cache

// cacheItem holds cache for git config.
type cacheItem struct {
	config   GitConfig
	filename string
	time     time.Time
	size     int64
}

func (v *cacheItem) uptodate() bool {
	fi, err := os.Stat(v.filename)
	if err == nil && fi.ModTime() == v.time && fi.Size() == v.size {
		return true
	}
	return false
}

// CacheSet will set cache entry
func CacheSet(key string, cfg GitConfig, size int64, modTime time.Time) {
	if cache == nil {
		return
	}
	cache.Add(key, &cacheItem{
		config:   cfg,
		filename: key,
		time:     modTime,
		size:     size,
	})
}

// CacheGet returns git config if config file is up-to-date
func CacheGet(key string) (GitConfig, bool) {
	value, ok := cache.Get(key)
	if !ok {
		return nil, false
	}
	item, ok := value.(*cacheItem)
	if !ok {
		return nil, false
	}
	if !item.uptodate() {
		cache.Remove(key)
		return nil, false
	}
	return item.config, true
}

func init() {
	cache = lru.New(128)
}
