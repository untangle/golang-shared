package cache

import (
	"github.com/untangle/golang-shared/structs/cache/cacher"
	"github.com/untangle/golang-shared/structs/cache/sweeper"
)

// Cache with a struct responsible for kicking off sweeps of the cache.
type SweptCache struct {
	cacher.Cacher

	sweeper sweeper.Sweeper
}

// The sweeper runs a function when certain criteria is met. To give the function being run by the sweeper
// access to the cache, use a closure.
func (sweptCache *SweptCache) generateCleanupTask(cleanupFunc func(string, interface{}) bool) func() {

	cache := sweptCache.Cacher

	return func() {

		getNext := cache.GetIterator()

		for key, value, ok := getNext(); ok; key, value, ok = getNext() {
			if cleanupFunc(key, value) {
				sweptCache.Cacher.Remove(key)
			}
		}

	}
}

func NewSweptCache(cache cacher.Cacher, sweeper sweeper.Sweeper) *SweptCache {

	return &SweptCache{cache, sweeper}
}

// Starts sweeping the cache. The function provided will be run on every element in the cache
// once a sweep is triggered.
// The key of the cache element, and a pointer to it, must be handled by the provided function.
// If false is returned from the provided function, the cache element will be removed.
func (sweptCache *SweptCache) StartSweeping(cleanupFunc func(string, interface{}) bool) {
	sweptCache.sweeper.StartSweeping(sweptCache.generateCleanupTask(cleanupFunc))
}

func (sweptCache *SweptCache) StopSweeping() {
	sweptCache.sweeper.StopSweeping()
}
