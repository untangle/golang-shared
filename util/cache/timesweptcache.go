package cache

import (
	"time"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/cache/cacher"
)

type TimeSweptCache struct {
	cacher.Cacher

	shutdownChannel chan bool
	waitTime        time.Duration
}

func NewTimeSweptCache(cache cacher.Cacher, waitTime time.Duration) *TimeSweptCache {
	return &TimeSweptCache{
		Cacher:          cache,
		shutdownChannel: make(chan bool),
		waitTime:        waitTime,
	}
}

func (sweeper *TimeSweptCache) StartSweeping(cleanupFunc func(string, interface{}) bool) {
	go sweeper.runCleanup(cleanupFunc)
}

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

func (sweeper *TimeSweptCache) StopSweeping() {
	sweeper.shutdownChannel <- true

	select {
	case <-sweeper.shutdownChannel:
		logger.Info("Successful shutdown of clean up \n")
	case <-time.After(10 * time.Second):
		logger.Warn("Failed to properly shutdown cleanupTask\n")
	}
}
