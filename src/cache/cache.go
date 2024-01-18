package cache

import (
	"github.com/golang/groupcache/lru"
	"platform-sre-interview-excercise-master/config"
	"time"
)

type (
	ratesCache struct {
		cache *lru.Cache
	}
	Entry struct {
		Value      float64
		Expiration time.Time
	}
)

func NewCache() ratesCache {
	//creates an LRU cache with a maximum capacity
	return ratesCache{cache: lru.New(100)}
}

func (rc *ratesCache) Get(key string) (float64, bool) {
	value, found := rc.cache.Get(key)
	if found {
		ce := value.(Entry)
		if time.Now().After(ce.Expiration) {
			// entry has expired, remove it from the cache
			rc.cache.Remove(key)
			return 0, false
		}
		return ce.Value, found
	}
	return 0, false
}

func (rc *ratesCache) Add(key string, value float64) {
	expiration := time.Now().Add(config.AppConfig.CacheExpiry)
	ce := Entry{Value: value, Expiration: expiration}
	rc.cache.Add(key, ce)
}
