package lrucache

import (
	"container/list"
)

type KeyPair struct {
	Key   string
	Value interface{}
}

type LruCache struct {
	Capacity uint
	List     *list.List
	Elements map[string]*list.Element
}

func NewLruCache(capacity uint) *LruCache {
	return &LruCache{capacity, list.New(), make(map[string]*list.Element)}
}

// Gets an item from the cache using a provided key. Once an item has been
// retrieved, move it to the front of the cache's queue. Return a bool to
// signify a value was found since a key could be mapped to nil.
func (cache *LruCache) Get(key string) (interface{}, bool) {
	var value interface{}
	var found bool

	if node, ok := cache.Elements[key]; ok {
		value = node.Value.(*list.Element).Value.(KeyPair).Value
		found = true
		cache.List.MoveToFront(node)
	}

	return value, found
}

// Add an item to the cache and move it to the front of the queue.
// If the item's key is already in the cache, update the key's value
// and move the the item to the front of the queue.
func (cache *LruCache) Put(key string, value interface{}) {
	// Update key's value if already present in the cache
	if node, ok := cache.Elements[key]; ok {
		cache.List.MoveToFront(node)
		node.Value.(*list.Element).Value = KeyPair{Key: key, Value: value}

	} else {
		// Remove least recently used item in cache if the cache's capacity has reached its limit
		if uint(cache.List.Len()) >= cache.Capacity {
			// Remove node from the cache's internal map
			elementToRemove := cache.List.Back().Value.(*list.Element).Value.(KeyPair).Key
			delete(cache.Elements, elementToRemove)

			cache.List.Remove(cache.List.Back())
		}
	}

	newNode := &list.Element{
		Value: KeyPair{
			Key:   key,
			Value: value,
		},
	}

	mostRecentlyUsed := cache.List.PushFront(newNode)
	cache.Elements[key] = mostRecentlyUsed
}

// Delete an item from the cache based off the key
func (cache *LruCache) Remove(key string) {
	if node, ok := cache.Elements[key]; ok {
		delete(cache.Elements, key)
		cache.List.Remove(node)
	}
}

// Clear all all internal data structures
func (cache *LruCache) Clear() {
	cache.Elements = make(map[string]*list.Element)
	cache.List.Init()
}

func (cache *LruCache) GetMostRecentlyUsed() (interface{}, interface{}) {
	keyPair := cache.List.Front().Value.(*list.Element).Value.(KeyPair)
	return keyPair.Key, keyPair.Value
}

func (cache *LruCache) GetLeastRecentlyUsed() (interface{}, interface{}) {
	keyPair := cache.List.Back().Value.(*list.Element).Value.(KeyPair)
	return keyPair.Key, keyPair.Value
}
