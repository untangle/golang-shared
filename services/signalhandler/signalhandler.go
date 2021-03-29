package signalhandler

import (
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"

	"github.com/untangle/golang-shared/services/logger"
)

// HandleShutdown is the type that holds the channel and flag for a shutdown
type HandleShutdown struct {
	shutdownFlag    uint32
	ShutdownChannel chan bool // ShutdownChannel is used to signal to other routines that the system is shutting down
}

// NewHandleShutdown creates a new HandleShutdown with channel and flag set
func NewHandleShutdown() *HandleShutdown {
	hs := new(HandleShutdown)
	hs.shutdownFlag = 0
	hs.ShutdownChannel = make(chan bool)

	return hs
}

// HandleSignals adds functionality to handle system signals
func (hs *HandleShutdown) HandleSignals() {
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
func (hs *HandleShutdown) dumpStack() {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)
	ioutil.WriteFile("/tmp/reportd.stack", buf[:stacklen], 0644)
	logger.Warn("Printing Thread Dump...\n")
	logger.Warn("\n\n%s\n\n", buf[:stacklen])
	logger.Warn("Thread dump complete.\n")
}

// PrintStats prints some basic stats about the running package
func (hs *HandleShutdown) PrintStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	logger.Info("Memory Stats:\n")
	logger.Info("Memory Alloc: %d kB\n", (mem.Alloc / 1024))
	logger.Info("Memory TotalAlloc: %d kB\n", (mem.TotalAlloc / 1024))
	logger.Info("Memory HeapAlloc: %d kB\n", (mem.HeapAlloc / 1024))
	logger.Info("Memory HeapSys: %d kB\n", (mem.HeapSys / 1024))
}

// GetShutdownFlag returns the shutdown flag for kernel
func (hs *HandleShutdown) GetShutdownFlag() bool {
	if atomic.LoadUint32(&hs.shutdownFlag) != 0 {
		return true
	}
	return false
}

// SetShutdownFlag sets the shutdown flag for kernel
func (hs *HandleShutdown) SetShutdownFlag() {
	hs.ShutdownChannel <- true
	atomic.StoreUint32(&hs.shutdownFlag, 1)
}
