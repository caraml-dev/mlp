package enforcer

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	store *cache.Cache
}

func newInMemoryCache(keyExpirySeconds int, cacheCleanUpIntervalSeconds int) *InMemoryCache {
	return &InMemoryCache{
		store: cache.New(
			time.Duration(keyExpirySeconds)*time.Second,
			time.Duration(cacheCleanUpIntervalSeconds)*time.Second,
		),
	}
}

// LookUpUserPermissions returns the cached permission check result for a user / permission pair.
// The returned value indicates whether the result is cached.
func (c *InMemoryCache) LookUpUserPermissions(user string, permission string) (*bool, bool) {
	if cachedValue, ok := c.store.Get(c.buildCacheKey(user, permission)); ok {
		if allowed, ok := cachedValue.(*bool); ok {
			return allowed, true
		}
	}
	return nil, false
}

// StoreUserPermissions stores the permission check result for a user / permission pair.
func (c *InMemoryCache) StoreUserPermissions(user string, permission string, result bool) {
	c.store.Set(c.buildCacheKey(user, permission), &result, cache.DefaultExpiration)
}

func (c *InMemoryCache) buildCacheKey(user string, permission string) string {
	return fmt.Sprintf("%s-%s", user, permission)
}
