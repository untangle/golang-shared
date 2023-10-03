package cacher

import (
	"container/list"
	"sync"

	logService "github.com/untangle/golang-shared/services/logger"
)

var logger = logService.GetLoggerInstance()

// Attach the key, used to look up a value in the cache, to its value.
type KeyPair struct {
	Key   string
	Value interface{}
}

// A simple LRU Cache implementation. The least recently used element
// in the cache are removed if the cache's max capacity is hit. The cache's
// mutex cannot be a RWMutex since Gets alter the cache's underlying data structures.
// Slower concurrent performance than a Random Removal cache, but more cache hits.
// O(1) reads and O(1) insertions.
type LruCache struct {
	capacity   uint
	list       *list.List
	cacheMutex sync.Mutex
	elements   map[string]*list.Element

	// Name of your cache. Only used to provide accurate logging
	cacheName string
}

// Returns a pointer to a newly initialized LruCache with it's capacity and name set to
// those provided by capactiy and cacheName, respectively
func NewLruCache(capacity uint, cacheName string) *LruCache {
	return &LruCache{
		capacity:  capacity,
		list:      list.New(),
		elements:  make(map[string]*list.Element),
		cacheName: cacheName,
	}
}

// Iterates over each key, value pair in the cache and runs them through
// the provided cleanup function. If the cleanup function provided returns true,
// The element will be removed from the cache. The cleanupFunction provided
// will be given the key and the value of the current element.
func (cache *LruCache) ForEach(cleanupFunc func(string, interface{}) bool) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	for key, val := range cache.elements {
		// Remove element if the cleanUp func returns true
		if cleanupFunc(key, val.Value.(*list.Element).Value.(KeyPair).Value) {
			cache.removeElementNoLock(key)
		}
	}

}

// Gets an item from the cache using a provided key. Once an item has been
// retrieved, move it to the front of the cache's queue. Return a bool to
// signify a value was found since a key could be mapped to nil.
func (cache *LruCache) Get(key string) (interface{}, bool) {
	var value interface{}
	var found bool

	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	if node, ok := cache.elements[key]; ok {
		value = node.Value.(*list.Element).Value.(KeyPair).Value
		found = true
		cache.list.MoveToFront(node)
	}

	return value, found
}

// Add an item to the cache and move it to the front of the queue.
// If the item's key is already in the cache, update the key's value
// and move the the item to the front of the queue.
func (cache *LruCache) Put(key string, value interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()

	// Update key's value if already present in the cache
	if node, ok := cache.elements[key]; ok {
		cache.list.MoveToFront(node)
		node.Value.(*list.Element).Value = KeyPair{Key: key, Value: value}

		logger.Debug("Updated the element with key %s in the cache named %s\n", key, cache.cacheName)
	} else {
		// Remove least recently used item in cache if the cache's capacity has reached its limit
		if uint(cache.list.Len()) >= cache.capacity {
			// Remove node from the cache's internal map
			elementToRemove := cache.list.Back().Value.(*list.Element).Value.(KeyPair).Key
			delete(cache.elements, elementToRemove)

			cache.list.Remove(cache.list.Back())
			logger.Debug("Removed element with key %s from the cache named %s\n", key, cache.cacheName)
		}

		newNode := &list.Element{
			Value: KeyPair{
				Key:   key,
				Value: value,
			},
		}

		mostRecentlyUsed := cache.list.PushFront(newNode)
		cache.elements[key] = mostRecentlyUsed
		logger.Debug("Added element with key %s to the cache named %s\n", key, cache.cacheName)
	}
}

// Removes an element from the cache.
// Does NOT take the cache's lock. Functions calling removeElementNoLock()
// need to do it themselves
func (cache *LruCache) removeElementNoLock(key string) {
	if node, ok := cache.elements[key]; ok {
		delete(cache.elements, key)
		cache.list.Remove(node)
		logger.Debug("Removed element with key %s from the cache name %s\n", key, cache.cacheName)
	}
}

// Delete an item from the cache based off the key
func (cache *LruCache) Remove(key string) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.removeElementNoLock(key)
}

// Clear all all internal data structures
func (cache *LruCache) Clear() {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	cache.elements = make(map[string]*list.Element)
	cache.list.Init()
	logger.Debug("Cleared cache of name %s\n", cache.cacheName)
}

// Gets the most recently looked up value in the cache
func (cache *LruCache) GetMostRecentlyUsed() (string, interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	keyPair := cache.list.Front().Value.(*list.Element).Value.(KeyPair)
	return keyPair.Key, keyPair.Value
}

// Get the least recently looked up value on the cache
func (cache *LruCache) GetLeastRecentlyUsed() (string, interface{}) {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	keyPair := cache.list.Back().Value.(*list.Element).Value.(KeyPair)
	return keyPair.Key, keyPair.Value
}

// Gets the total number of elements currently in the cache
func (cache *LruCache) GetTotalElements() int {
	cache.cacheMutex.Lock()
	defer cache.cacheMutex.Unlock()
	return cache.list.Len()
}
