package sweeper

import (
	"time"

	"github.com/untangle/golang-shared/services/logger"
)

type SweepOnTime struct {
	shutdownChannel chan bool
	waitTime        time.Duration
}

func NewSweepOnTime(waitTime time.Duration) *SweepOnTime {
	return &SweepOnTime{
		shutdownChannel: make(chan bool),
		waitTime:        waitTime,
	}
}

func (sweeper *SweepOnTime) StartSweeping(cleanupFunc func()) {
	go sweeper.runCleanup(cleanupFunc)
}

func (sweeper *SweepOnTime) StopSweeping() {
	sweeper.shutdownChannel <- true

	select {
	case <-sweeper.shutdownChannel:
		logger.Info("Successful shutdown of clean up \n")
	case <-time.After(10 * time.Second):
		logger.Warn("Failed to properly shutdown cleanupTask\n")
	}
}

func (sweeper *SweepOnTime) runCleanup(cleanupFunc func()) {
	for {
		select {
		case <-sweeper.shutdownChannel:
			sweeper.shutdownChannel <- true
			return
		case <-time.After(sweeper.waitTime):

			cleanupFunc()
		}
	}
}
