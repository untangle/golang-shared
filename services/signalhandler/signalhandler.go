package signalhandler

import (
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"

	logService "github.com/untangle/golang-shared/services/logger"
)

var logger = logService.GetLoggerInstance()

// SignalHandler is the type that holds the channel and flag for a shutdown
type SignalHandler struct {
	shutdownFlag    uint32
	ShutdownChannel chan bool // ShutdownChannel is used to signal to other routines that the system is shutting down
}

// NewSignalHandler creates a new SignalHandler with channel and flag set
func NewSignalHandler() *SignalHandler {
	hs := new(SignalHandler)
	hs.shutdownFlag = 0
	hs.ShutdownChannel = make(chan bool)

	return hs
}

// HandleSignals adds functionality to handle system signals
func (hs *SignalHandler) HandleSignals() {
	// Add SIGINT & SIGTERM handler (exit)
	termch := make(chan os.Signal, 1)
	signal.Notify(termch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-termch
		go func() {
			logger.Warn("Received signal [%v]. Shutting down routines...\n", sig)
			hs.SetShutdownFlag()
		}()
	}()

	// Add SIGQUIT handler (dump thread stack trace)
	quitch := make(chan os.Signal, 1)
	signal.Notify(quitch, syscall.SIGQUIT)
	go func() {
		for {
			sig := <-quitch
			logger.Info("Received signal [%v]. Calling dumpStack()\n", sig)
			go hs.dumpStack()
		}
	}()
}

// dumpStack dumps the stack trace to /tmp/reportd.stack and log
func (hs *SignalHandler) dumpStack() {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)
	err := ioutil.WriteFile("/tmp/reportd.stack", buf[:stacklen], 0644)
	if err != nil {
		logger.Warn("Failed to write data to file /tmp/reportd/stack with error : %v\n", err)
	}
	logger.Warn("Printing Thread Dump...\n")
	logger.Warn("Printing Thread Dump...\n")
	logger.Warn("\n\n%s\n\n", buf[:stacklen])
	logger.Warn("Thread dump complete.\n")
}

// PrintStats prints some basic stats about the running package
func (hs *SignalHandler) PrintStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	logger.Debug("Memory Stats:\n")
	logger.Debug("Memory Alloc: %d kB\n", (mem.Alloc / 1024))
	logger.Debug("Memory TotalAlloc: %d kB\n", (mem.TotalAlloc / 1024))
	logger.Debug("Memory HeapAlloc: %d kB\n", (mem.HeapAlloc / 1024))
	logger.Debug("Memory HeapSys: %d kB\n", (mem.HeapSys / 1024))
}

// GetShutdownFlag returns the shutdown flag for kernel
func (hs *SignalHandler) GetShutdownFlag() bool {
	return atomic.LoadUint32(&hs.shutdownFlag) != 0
}

// SetShutdownFlag sets the shutdown flag for kernel
func (hs *SignalHandler) SetShutdownFlag() {
	hs.ShutdownChannel <- true
	atomic.StoreUint32(&hs.shutdownFlag, 1)
}
