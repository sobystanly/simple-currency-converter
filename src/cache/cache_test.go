package cache

import (
	"github.com/golang/groupcache/lru"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewCache(t *testing.T) {
	t.Run("successfully initialize the cache", func(t *testing.T) {
		cache := NewCache()
		assert.Equal(t, ratesCache{cache: lru.New(100)}, cache)
	})
}

func TestRatesCache_Get_And_Add(t *testing.T) {
	t.Run("Given conversion of euro to usd is in cache fetch that successfully", func(t *testing.T) {
		cache := NewCache()
		cache.Add("eur-usd", 0.011991018)

		actual, found := cache.Get("eur-usd")
		assert.Equal(t, true, found)
		assert.Equal(t, 0.011991018, actual)
	})
}
