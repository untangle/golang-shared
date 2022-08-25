package cacher

import (
	"math/rand"
	"sync"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/util"
)

// Simple cache that removes elements randomly when the cache capacity is met.
// O(1) lookups, but insertions are O(n) when the capacity is met.
// The cache can be read from by multiple threads, but written to by one.
type RandomReplacementCache struct {
	maxCapacity uint
	elements    map[string]interface{}
	cacheName   string
	cacheMutex  sync.RWMutex
}

func NewRandomReplacementCache(capacity uint, cacheName string) *RandomReplacementCache {
	return &RandomReplacementCache{
		maxCapacity: capacity,
		elements:    make(map[string]interface{}),
		cacheName:   cacheName,
	}
}

func (cache *RandomReplacementCache) GetIterator() func() (string, *interface{}, bool) {
	// Once an iterator has been retrieved, it captures the state of
	// of the cache. If the cache is updated the iterator won't contain
	// the update
	cache.cacheMutex.RLock()
	keys := util.GetMapKeys(cache.elements)
	cache.cacheMutex.RUnlock()

	i := 0
	// Return key, val, and if there is anything left to iterate over
	return func() (string, *interface{}, bool) {
		if i == len(keys) {
			return "", nil, false
		}

		currentKey := keys[i]

		// The value could be nil if the map was altered
		value, _ := cache.Get(currentKey)
		i += 1
		return currentKey, &value, true
	}
}

func (cache *RandomReplacementCache) Get(key string) (interface{}, bool) {
	cache.cacheMutex.RLock()
	defer cache.cacheMutex.RUnlock()

	value, ok := cache.elements[key]

	return value, ok
}

func (cache *RandomReplacementCache) getRandomElement() string {
	// rand's range is exclusive, and so is range. In order to
	// randomly select all elements in the cache, add one to
	// the range given to rand
	indexToRemove := rand.Intn((len(cache.elements) + 1))
	var keyToRemove string

	count := 0
	for key := range cache.elements {
		if count == indexToRemove-1 {
			keyToRemove = key
		}

		count += 1
	}

	return keyToRemove
}

func (cache *RandomReplacementCache) Put(key string, value interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	// Update element if already present in cache5
	if _, ok := cache.elements[key]; ok {
		cache.elements[key] = value
		logger.Debug("Updated the element with key %s in the cache named %s", key, cache.cacheName)
	} else {
		// Remove element if the capacity has been met
		if uint(len(cache.elements)) >= cache.maxCapacity {
			delete(cache.elements, cache.getRandomElement())
			logger.Debug("Removed element with key %s from the cache named %s", key, cache.cacheName)
		}

		// Add new element
		cache.elements[key] = value
		logger.Debug("Added element with key %s to the cache named %s", key, cache.cacheName)

	}

}

func (cache *RandomReplacementCache) Remove(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	delete(cache.elements, key)
	logger.Debug("Removed element with key %s from the cache named %s", key, cache.cacheName)
}

func (cache *RandomReplacementCache) Clear() {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.elements = make(map[string]interface{})
	logger.Debug("Cleared cache of name %s", cache.cacheName)
}

func (cache *RandomReplacementCache) GetCurrentCapacity() int {
	cache.cacheMutex.RLock()
	defer cache.cacheMutex.RUnlock()
	return len(cache.elements)
}
