package cache

import (
	"time"

	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/util/cache/cacher"
)

var logger = logService.GetLoggerInstance()

// Adds the ability to sweep elements on a timer to a cache.
// A goroutine is started which sweeps through
// the underlying cache on a set interval.
type TimeSweptCache struct {
	cacher.Cacher

	shutdownChannel chan bool
	waitTime        time.Duration
}

// Returns a pointer to an initialized TimeSweptCache. The underlying cache type is set to
// what is provided by cache. The interval in which the sweeper runs is set by waitTime.
func NewTimeSweptCache(cache cacher.Cacher, waitTime time.Duration) *TimeSweptCache {
	return &TimeSweptCache{
		Cacher:          cache,
		shutdownChannel: make(chan bool),
		waitTime:        waitTime,
	}
}

// Starts sweeping the cache at the provided interval with the
// cleanupFunc provided. The cleanUp func will be provided with the key/val of
// each element in the cache. If the element should be deleted from the cache,
// return true. If the cache value is a pointer to an object, this method
// can be used to mutate a cache value.
// Does not do a sweep on call, only after waitTime.
func (sweeper *TimeSweptCache) StartSweeping(cleanupFunc func(string, interface{}) bool) {
	go sweeper.runCleanup(cleanupFunc)
}

// Runs the cleanup function
func (sweeper *TimeSweptCache) runCleanup(cleanupFunc func(string, interface{}) bool) {
	for {
		select {
		case <-sweeper.shutdownChannel:
			sweeper.shutdownChannel <- true
			return
		case <-time.After(sweeper.waitTime):
			sweeper.Cacher.ForEach(cleanupFunc)
		}
	}
}

// Stops sweepign the cache
func (sweeper *TimeSweptCache) StopSweeping() {
	sweeper.shutdownChannel <- true

	select {
	case <-sweeper.shutdownChannel:
		logger.Info("Successful shutdown of clean up \n")
	case <-time.After(10 * time.Second):
		logger.Warn("Failed to properly shutdown cleanupTask\n")
	}
}
