package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/untangle/discoverd/plugins/arp"
	"github.com/untangle/discoverd/plugins/lldp"
	"github.com/untangle/discoverd/plugins/nmap"
	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/discoverd/services/example"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/profiler"
)

var shutdownFlag uint32
var shutdownChannel = make(chan bool)
var cpuProfileFilename = ""
var cpuProfiler *profiler.CPUProfiler

/* main function for discoverd */
func main() {
	// Check we are root user
	userinfo, err := user.Current()
	if err != nil {
		panic(err)
	}

	userid, err := strconv.Atoi(userinfo.Uid)
	if err != nil {
		panic(err)
	}

	if userid != 0 {
		panic("This application must be run as root\n")
	}

	// Start up logger
	loggerConfig := createLoggerConfig()
	logger.Startup(loggerConfig)
	logger.Info("Starting up discoverd...\n")

	parseArguments()

	// setup CPU profiling
	cpuProfiler = profiler.CreateCPUProfiler(cpuProfileFilename)
	err = cpuProfiler.StartCPUProfile()
	if err != nil {
		logger.Warn("CPU Profiler could not start: %s\n", err.Error())
	}

	// Start services
	startServices()

	startPlugins()

	// Handle the stop signals
	handleSignals()

	// Keep discoverd running while the shutdown flag is false
	// shutdown once flag is true or the shutdownChannel indicates a shutdown
	for !GetShutdownFlag() {
		select {
		case <-shutdownChannel:
			logger.Info("Shutdown channel initiated... %v\n", GetShutdownFlag())
		case <-time.After(2 * time.Minute):
			logger.Debug("discoverd is running...\n")
			logger.Info("\n")
			printStats()
		}
	}

	logger.Info("Shutdown discoverd...\n")

	stopServices()
	stopPlugins()

	cpuProfiler.StopCPUProfile()
}

func getPluginSettings() map[string]interface{} {
	pluginSettings := make(map[string]interface{})

}

/* startServices starts the gin server and cert manager */
func startServices() {
	example.Startup()
	discovery.Startup()
}

/* stopServices stops the gin server, cert manager, and logger*/
func stopServices() {
	example.Shutdown()
}

func startPlugins() {
	arp.Start()
	lldp.Start()
	nmap.Start()
}

func stopPlugins() {
	arp.Stop()
	lldp.Stop()
	nmap.Stop()
}

/* handleSignals handles SIGINT, SIGTERM, and SIGQUIT signals */
func handleSignals() {
	// Add SIGINT & SIGTERM handler (exit)
	termch := make(chan os.Signal, 1)
	signal.Notify(termch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-termch
		go func() {
			logger.Info("Received signal [%v]. Setting shutdown flag\n", sig)
			SetShutdownFlag()
		}()
	}()

	// Add SIGQUIT handler (dump thread stack trace)
	quitch := make(chan os.Signal, 1)
	signal.Notify(quitch, syscall.SIGQUIT)
	go func() {
		for {
			sig := <-quitch
			logger.Info("Received signal [%v]. Calling dumpStack()\n", sig)
			go dumpStack()
		}
	}()

	// Add SIGHUP handler (call handlers)
	hupch := make(chan os.Signal, 1)
	signal.Notify(hupch, syscall.SIGHUP)
	go func() {
		for {
			sig := <-hupch
			logger.Info("Received signal [%v]. Calling handlers\n", sig)
			targets := []func(syscall.Signal){
				stats.PluginSignal, threatprevention.PluginSignal,
				webfilter.PluginSignal}
			sig.Signal()
			plugins.GlobalPluginControl().Signal(syscall.SIGHUP)
			notifyTargets(syscall.SIGHUP, targets)

		}
	}()
}

// notifyTargets signals all plugins with a handler (in parallel)
func notifyTargets(message syscall.Signal, targets []func(syscall.Signal)) {
	var wg sync.WaitGroup

	for _, f := range targets {
		wg.Add(1)
		go func(f func(syscall.Signal)) {
			f(message)
			wg.Done()
		}(f)
	}

	wg.Wait()
}

// dumpStack to /tmp/discoverd.stack and log
func dumpStack() {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)
	_ = ioutil.WriteFile("/tmp/discoverd.stack", buf[:stacklen], 0644)
	logger.Warn("Printing Thread Dump...\n")
	logger.Warn("\n\n%s\n\n", buf[:stacklen])
	logger.Warn("Thread dump complete.\n")
}

// prints some basic stats about discoverd
func printStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	logger.Debug("Memory Stats:\n")
	logger.Debug("Memory Alloc: %d kB\n", (mem.Alloc / 1024))
	logger.Debug("Memory TotalAlloc: %d kB\n", (mem.TotalAlloc / 1024))
	logger.Debug("Memory HeapAlloc: %d kB\n", (mem.HeapAlloc / 1024))
	logger.Debug("Memory HeapSys: %d kB\n", (mem.HeapSys / 1024))
}

// GetShutdownFlag returns the shutdown flag for kernel
func GetShutdownFlag() bool {
	return atomic.LoadUint32(&shutdownFlag) != 0
}

// SetShutdownFlag sets the shutdown flag for kernel
func SetShutdownFlag() {
	shutdownChannel <- true
	atomic.StoreUint32(&shutdownFlag, 1)
}

// parseArguments parses the command line arguments
func parseArguments() {
	logger.Debug("Parsing cmd arguments\n")

	cpuProfilePtr := flag.String("cpuprofile", "", "filename for CPU pprof output")
	nmapNetworkPtr := flag.String("nmap-network", "", "network to scan for hosts")

	flag.Parse()

	if len(*cpuProfilePtr) > 0 {
		cpuProfileFilename = *cpuProfilePtr
	}

	if len(*nmapNetworkPtr) > 0 {
		nmap.SetNetwork(*nmapNetworkPtr)
	}
}

// createLoggerConfig creates the logger config
func createLoggerConfig() logger.Config {
	config := logger.Config{
		FileLocation: "/tmp/logconfig_discoverd.json",
		LogLevelMap:  getLogLevels(),
	}

	return config
}

// getLogLevels returns the default log levels for each service and plugin
func getLogLevels() map[string]string {
	return map[string]string{
		// services
		"example":   "INFO",
		"discovery": "INFO",
		"arp":       "INFO",
		"lldp":      "INFO",
		"nmap":      "INFO",
	}
}
