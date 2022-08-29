package cacher

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/util"
)

// Simple cache that removes elements randomly when the cache capacity is met.
// O(1) lookups and insertions. For O(1) insertions, space complexity had to be increased.
// The size of the cache will grow for every cache deletion since the keys slice can't
// have elements removed from it.
// The cache can be read from by multiple threads, but written to by one.
type RandomReplacementCache struct {
	maxCapacity uint
	elements    map[string]interface{}
	cacheName   string
	cacheMutex  sync.RWMutex

	keys []string
}

func NewRandomReplacementCache(capacity uint, cacheName string) *RandomReplacementCache {
	return &RandomReplacementCache{
		maxCapacity: capacity,
		elements:    make(map[string]interface{}, capacity),
		cacheName:   cacheName,

		// Removing elements from the keys slice would cause an entire update of the
		// keyToIndex map. For a performance bump, just set removed element's keys to
		// nil when removed. Since keys capacity will exceed that of the maps, give it
		// a much larger size to avoid too many copies
		keys: make([]string, 2*capacity),
	}
}

func (cache *RandomReplacementCache) ForEach(cleanupFunc func(string, interface{}) bool) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	for key, val := range cache.elements {
		// Remove element if the cleanUp func returns true
		if cleanupFunc(key, val) {
			//cache.removeElement(key)
			fmt.Print("whateverI ")
		}
	}

}

func (cache *RandomReplacementCache) GetIterator() func() (string, interface{}, bool) {
	// Once an iterator has been retrieved, it captures the state of
	// of the cache. If the cache is updated the iterator won't contain
	// the update
	cache.cacheMutex.RLock()
	keys := util.GetMapKeys(cache.elements)
	cache.cacheMutex.RUnlock()

	i := 0
	// Return key, val, and if there is anything left to iterate over
	return func() (string, interface{}, bool) {
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

func (cache *RandomReplacementCache) Put(key string, value interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	// Update element if already present in cache
	if _, ok := cache.elements[key]; ok {
		cache.elements[key] = value
		logger.Debug("Updated the element with key %s in the cache named %s", key, cache.cacheName)
	} else {
		// Remove element randomly if the capacity has been met
		if uint(len(cache.elements)) >= cache.maxCapacity {
			indexToSwap := rand.Intn(len(cache.keys))
			keyToRemove := cache.keys[indexToSwap]

			// Swap the last key with a random key. Delete the new last key.
			lastElementCopy := cache.keys[len(cache.keys)-1]
			cache.keys[len(cache.keys)-1] = cache.keys[indexToSwap]
			cache.keys[indexToSwap] = lastElementCopy
			cache.keys = cache.keys[:len(cache.keys)-1]

			delete(cache.elements, keyToRemove)

			logger.Debug("Removed element with key %s from the cache named %s", key, cache.cacheName)
		}

		// Add new element
		cache.elements[key] = value
		cache.keys = append(cache.keys, key)

		logger.Debug("Added element with key %s to the cache named %s", key, cache.cacheName)

	}
}

// Remove is an O(n) operation since the key to be removed must be found first
func (cache *RandomReplacementCache) Remove(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	if _, ok := cache.elements[key]; ok {
		delete(cache.elements, key)

		for i := 0; i < len(cache.keys); i++ {
			if cache.keys[i] == key {
				// Order doesn't matter for the keys slice, so delete the fast way.
				// Which is just swapping the element to delete with the last element
				// then ignoring the last element of the slice
				cache.keys[i] = cache.keys[len(cache.keys)-1]
				cache.keys = cache.keys[:len(cache.keys)-1]
			}
		}

	}
	// else the key didn't exists in the cache and nothing should be done
}

func (cache *RandomReplacementCache) Clear() {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.elements = make(map[string]interface{}, cache.maxCapacity)
	cache.keys = make([]string, cache.maxCapacity)
	logger.Debug("Cleared cache of name %s", cache.cacheName)
}

func (cache *RandomReplacementCache) GetCurrentCapacity() int {
	cache.cacheMutex.RLock()
	defer cache.cacheMutex.RUnlock()
	return len(cache.elements)
}
