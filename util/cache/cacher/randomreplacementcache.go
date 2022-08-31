package cacher

import (
	"math/rand"
	"sync"

	"github.com/untangle/golang-shared/services/logger"
)

type Value struct {
	keyIndex uint
	value    interface{}
}

// Simple cache that removes elements randomly when the cache capacity is met.
// O(1) lookups and insertions. For O(1) insertions, space complexity had to be increased.
// The size of the cache will grow for every cache deletion since the keys slice can't
// have elements removed from it.
// The cache can be read from by multiple threads, but written to by one.
type RandomReplacementCache struct {
	maxCapacity uint
	elements    map[string]*Value
	cacheName   string
	cacheMutex  sync.RWMutex

	keys          []string
	totalElements uint
}

func NewRandomReplacementCache(capacity uint, cacheName string) *RandomReplacementCache {
	return &RandomReplacementCache{
		maxCapacity:   capacity,
		elements:      make(map[string]*Value, capacity),
		cacheName:     cacheName,
		keys:          make([]string, capacity),
		totalElements: 0,
	}
}

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

// func (cache *RandomReplacementCache) GetIterator() func() (string, interface{}, bool) {
// 	// Once an iterator has been retrieved, it captures the state of
// 	// of the cache. If the cache is updated the iterator won't contain
// 	// the update
// 	cache.cacheMutex.RLock()
// 	keys := getMapKeys(&cache.elements)
// 	cache.cacheMutex.RUnlock()

// 	i := 0
// 	// Return key, val, and if there is anything left to iterate over
// 	return func() (string, interface{}, bool) {
// 		if i == len(keys) {
// 			return "", nil, false
// 		}

// 		currentKey := keys[i]

// 		// The value could be nil if the map was altered
// 		value, _ := cache.Get(currentKey)
// 		i += 1
// 		return currentKey, &value, true
// 	}
// }

func (cache *RandomReplacementCache) Get(key string) (interface{}, bool) {
	cache.cacheMutex.RLock()
	defer cache.cacheMutex.RUnlock()

	if value, ok := cache.elements[key]; ok {
		return value.value, ok
	} else {
		return nil, ok
	}
}

func (cache *RandomReplacementCache) Put(key string, value interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	// Update element if already present in cache
	if _, ok := cache.elements[key]; ok {
		cache.elements[key].value = value
		logger.Debug("Updated the element with key %s in the cache named %s", key, cache.cacheName)
	} else {
		// Remove element randomly if the capacity has been met
		if cache.totalElements >= cache.maxCapacity {
			indexToSwap := rand.Intn(len(cache.keys))
			keyToRemove := cache.keys[indexToSwap]
			cache.removeWithoutLock(keyToRemove)

			logger.Debug("Removed element with key %s from the cache named %s", key, cache.cacheName)
		}

		// Add new element
		cache.totalElements += 1
		cache.elements[key] = &Value{keyIndex: cache.totalElements - 1, value: value}
		cache.keys[cache.totalElements-1] = key
		logger.Debug("Added element with key %s to the cache named %s", key, cache.cacheName)
	}
}

func (cache *RandomReplacementCache) removeWithoutLock(key string) {
	if _, ok := cache.elements[key]; ok {
		indexToRemove := cache.elements[key].keyIndex

		// Order doesn't matter for the keys slice, so delete the fast way.
		// Which is just swapping the element to delete with the last element
		// then ignoring the last element of the slice
		cache.keys[indexToRemove] = cache.keys[len(cache.keys)-1]
		cache.keys[len(cache.keys)-1] = ""

		// Update index of moved element
		cache.elements[key].keyIndex = indexToRemove

		delete(cache.elements, key)
		cache.totalElements -= 1
	}
	// else the key didn't exists in the cache and nothing should be done
}

// Remove is an O(n) operation since the key to be removed must be found first
func (cache *RandomReplacementCache) Remove(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.removeWithoutLock(key)
}

func (cache *RandomReplacementCache) Clear() {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.elements = make(map[string]*Value, cache.maxCapacity)
	cache.keys = make([]string, cache.maxCapacity)
	cache.totalElements = 0
	logger.Debug("Cleared cache of name %s", cache.cacheName)
}

func (cache *RandomReplacementCache) GetCurrentCapacity() int {
	cache.cacheMutex.RLock()
	defer cache.cacheMutex.RUnlock()
	return int(cache.totalElements)
}
