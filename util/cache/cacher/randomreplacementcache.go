package cacher

import (
	"math/rand"
	"sync"
)

// The value with it's corresponding index in the slice
// used to keep track of all the keys in the cache.
type Value struct {
	keyIndex uint
	value    interface{}
}

// Simple cache that removes elements randomly when the cache capacity is met.
// O(1) lookups and insertions. For O(1) insertions, space complexity had to be increased
// by adding a few data structures for bookkeeping.
// The cache can be read from by multiple threads, but written to by one.
type RandomReplacementCache struct {
	maxCapacity uint
	elements    map[string]*Value
	cacheName   string
	cacheMutex  sync.RWMutex

	// Keys is a slice of all the keys in the cache
	// Used to randomly select which item should be
	// removed from the cache when the capacity is
	// exceeded
	keys          []string
	totalElements uint
}

// Returns a pointer to a newly initialized RandomReplacement and sets its cache capacity
// and name to those provided by capacity and cacheName, respectively
func NewRandomReplacementCache(capacity uint, cacheName string) *RandomReplacementCache {
	return &RandomReplacementCache{
		maxCapacity:   capacity,
		elements:      make(map[string]*Value, capacity),
		cacheName:     cacheName,
		keys:          make([]string, capacity),
		totalElements: 0,
	}
}

// Iterates over each key, value pair in the cache and runs them through
// the provided cleanup function. If the cleanup function provided returns true,
// The element will be removed from the cache
func (cache *RandomReplacementCache) ForEach(cleanupFunc func(string, interface{}) bool) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	for key, val := range cache.elements {
		// Remove element if the cleanUp func returns true
		if cleanupFunc(key, val.value) {
			cache.removeWithoutLock(key)
		}
	}
}

// It's useful to get the keys directly from the map instead of the array of keys
// Since the array of keys will have empty strings in it
func getMapKeys(m *map[string]*Value) []string {
	keys := make([]string, len(*m))

	i := 0
	for key := range *m {
		keys[i] = key
		i++
	}

	return keys
}

// Retrieves a value from the cache corresponding to the provided key
func (cache *RandomReplacementCache) Get(key string) (interface{}, bool) {
	cache.cacheMutex.RLock()
	defer cache.cacheMutex.RUnlock()

	if value, ok := cache.elements[key]; ok {
		return value.value, ok
	} else {
		return nil, ok
	}
}

// Places the provided value in the cache. If it already exists, the new value
// provided replaces the previous one. If the capacity of the cache has been met,
// an element from the cache is randomly deleted and the provided value is added.
func (cache *RandomReplacementCache) Put(key string, value interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	// Update element if already present in cache
	if _, ok := cache.elements[key]; ok {
		cache.elements[key].value = value
		logger.Debug("Updated the element with key %s in the cache named %s\n", key, cache.cacheName)
	} else {
		// Remove element randomly if the capacity has been met
		if cache.totalElements >= cache.maxCapacity {
			indexToSwap := rand.Intn(len(cache.keys))
			keyToRemove := cache.keys[indexToSwap]
			cache.removeWithoutLock(keyToRemove)

			logger.Debug("Removed element with key %s from the cache named %s\n", key, cache.cacheName)
		}

		// Add new element
		cache.totalElements += 1
		cache.elements[key] = &Value{keyIndex: cache.totalElements - 1, value: value}
		cache.keys[cache.totalElements-1] = key
		logger.Debug("Added element with key %s to the cache named %s\n", key, cache.cacheName)
	}
}

// Deletes an element from the cache. Does not acquire the mutex lock
// Any function calling this should acquire the mutex lock there
func (cache *RandomReplacementCache) removeWithoutLock(key string) {
	if _, ok := cache.elements[key]; ok {
		indexToRemove := cache.elements[key].keyIndex

		// Order doesn't matter for the keys slice, so delete the fast way.
		// Which is swapping the element to delete with the last element
		// then ignoring the last element of the slice
		cache.keys[indexToRemove] = cache.keys[cache.totalElements-1]

		// Update index of moved element
		movedElementKey := cache.keys[indexToRemove]
		cache.elements[movedElementKey].keyIndex = indexToRemove

		delete(cache.elements, key)
		cache.keys[cache.totalElements-1] = ""
		cache.totalElements -= 1
	}
}

// Removes an element from the cache.
func (cache *RandomReplacementCache) Remove(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.removeWithoutLock(key)
}

// Removes all elements from the cache
func (cache *RandomReplacementCache) Clear() {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.elements = make(map[string]*Value, cache.maxCapacity)
	cache.keys = make([]string, cache.maxCapacity)
	cache.totalElements = 0
	logger.Debug("Cleared cache of name %s\n", cache.cacheName)
}
