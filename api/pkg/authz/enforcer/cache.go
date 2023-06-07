package enforcer

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ory/keto-client-go/models"
	cache "github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	store *cache.Cache

	mapLock sync.Mutex

	// Instead of caching the full names of the subject, resource and action, ids will be
	// generated (using an incrementing counter) for each unique value. This will aid in
	// generating smaller cache keys.
	// The following maps store the mapping of the name -> internal id.
	subjectMap  map[string]string
	resourceMap map[string]string
	actionMap   map[string]string
}

func newInMemoryCache(keyExpirySeconds int, cacheCleanUpIntervalSeconds int) *InMemoryCache {
	return &InMemoryCache{
		store: cache.New(
			time.Duration(keyExpirySeconds)*time.Second,
			time.Duration(cacheCleanUpIntervalSeconds)*time.Second,
		),
		subjectMap:  map[string]string{},
		resourceMap: map[string]string{},
		actionMap:   map[string]string{},
	}
}

func (c *InMemoryCache) LookUpPermission(input models.OryAccessControlPolicyAllowedInput) (*bool, bool) {
	if cachedValue, ok := c.store.Get(c.buildCacheKey(input)); ok {
		if allowed, ok := cachedValue.(*bool); ok {
			return allowed, true
		}
	}
	return nil, false
}

func (c *InMemoryCache) StorePermission(input models.OryAccessControlPolicyAllowedInput, isAllowed *bool) {
	c.store.Set(c.buildCacheKey(input), isAllowed, cache.DefaultExpiration)
}

func (c *InMemoryCache) buildCacheKey(input models.OryAccessControlPolicyAllowedInput) string {
	return fmt.Sprintf("%s:%s:%s",
		c.getActionID(input.Action),
		c.getSubjectID(input.Subject),
		c.getResourceID(input.Resource),
	)
}

func (c *InMemoryCache) getSubjectID(name string) string {
	c.mapLock.Lock()
	defer c.mapLock.Unlock()

	if val, ok := c.subjectMap[name]; ok {
		return val
	}
	newID := strconv.Itoa(countMapKeys(c.subjectMap) + 1)
	c.subjectMap[name] = newID
	return newID
}

func (c *InMemoryCache) getResourceID(name string) string {
	c.mapLock.Lock()
	defer c.mapLock.Unlock()

	if val, ok := c.resourceMap[name]; ok {
		return val
	}
	newID := strconv.Itoa(countMapKeys(c.resourceMap) + 1)
	c.resourceMap[name] = newID
	return newID
}

func (c *InMemoryCache) getActionID(name string) string {
	c.mapLock.Lock()
	defer c.mapLock.Unlock()

	if val, ok := c.actionMap[name]; ok {
		return val
	}
	newID := strconv.Itoa(countMapKeys(c.actionMap) + 1)
	c.actionMap[name] = newID
	return newID
}

func countMapKeys(m map[string]string) int {
	count := 0
	for range m {
		count++
	}
	return count
}
