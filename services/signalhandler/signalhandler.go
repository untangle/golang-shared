package signalhandler

import (
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/untangle/golang-shared/services/logger"
)

// SignalHandler is the type that holds the channel and flag for a shutdown
type SignalHandler struct {
	shutdownFlag    uint32
	ShutdownChannel chan bool // ShutdownChannel is used to signal to other routines that the system is shutting down
	Targets         []func(syscall.Signal)
}

// NewSignalHandler creates a new SignalHandler with channel and flag set
func NewSignalHandler() *SignalHandler {
	hs := new(SignalHandler)
	hs.shutdownFlag = 0
	hs.ShutdownChannel = make(chan bool)

	return hs
}

func (hs *SignalHandler) SetTargets(targets []func(syscall.Signal)) {
	hs.Targets = targets
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

	// Add SIGHUP handler (call handlers)
	hupch := make(chan os.Signal, 1)
	signal.Notify(hupch, syscall.SIGHUP)
	go func() {
		for {
			sig := <-hupch
			if len(hs.Targets) != 0 {
				logger.Info("Received signal [%v]. Calling handlers\n", sig)
				hs.signalPlugins(syscall.SIGHUP)
			} else {
				logger.Info("No targets for signal, doing nothing...\n")
			}
		}
	}()
}

// signalPlugins signals all plugins with a handler (in parallel)
func (hs *SignalHandler) signalPlugins(message syscall.Signal) {
	var wg sync.WaitGroup

	for _, f := range hs.Targets {
		wg.Add(1)
		go func(f func(syscall.Signal)) {
			f(message)
			wg.Done()
		}(f)
	}

	wg.Wait()
}

// dumpStack dumps the stack trace to /tmp/reportd.stack and log
func (hs *SignalHandler) dumpStack() {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)
	ioutil.WriteFile("/tmp/reportd.stack", buf[:stacklen], 0644)
	logger.Warn("Printing Thread Dump...\n")
	logger.Warn("\n\n%s\n\n", buf[:stacklen])
	logger.Warn("Thread dump complete.\n")
}

// PrintStats prints some basic stats about the running package
func (hs *SignalHandler) PrintStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	logger.Info("Memory Stats:\n")
	logger.Info("Memory Alloc: %d kB\n", (mem.Alloc / 1024))
	logger.Info("Memory TotalAlloc: %d kB\n", (mem.TotalAlloc / 1024))
	logger.Info("Memory HeapAlloc: %d kB\n", (mem.HeapAlloc / 1024))
	logger.Info("Memory HeapSys: %d kB\n", (mem.HeapSys / 1024))
}

// GetShutdownFlag returns the shutdown flag for kernel
func (hs *SignalHandler) GetShutdownFlag() bool {
	if atomic.LoadUint32(&hs.shutdownFlag) != 0 {
		return true
	}
	return false
}

// SetShutdownFlag sets the shutdown flag for kernel
func (hs *SignalHandler) SetShutdownFlag() {
	hs.ShutdownChannel <- true
	atomic.StoreUint32(&hs.shutdownFlag, 1)
}
