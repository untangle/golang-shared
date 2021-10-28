package profiler

import (
	"errors"
	"os"
	"runtime/pprof"

	"github.com/untangle/golang-shared/services/logger"
)

/**
Usage:
// On Service Start
cpuProfiler = profiler.CreateCPUProfiler(cpuProfileFilename)
err = cpuProfiler.StartCPUProfile()
if err != nil {
	logger.Warn("CPU Profiler could not start: %s\n", err.Error())
}

// On Service Shutdown
cpuProfiler.StopCPUProfile()
*/

// CPUProfiler struct wraps functionality for profiling the CPU
type CPUProfiler struct {
	CPUProfileFileName string
	file               *os.File
	IsRunning          bool
}

// StartCPUProfile sets up and starts cpu profiling
func (cpuProfiler *CPUProfiler) StartCPUProfile() error {
	if cpuProfiler.CPUProfileFileName == "" {
		return errors.New("Cannot start cpu profiling. CPUProfileFileName must be specified!")
	}
	cpu, err := os.Create(cpuProfiler.CPUProfileFileName)
	cpuProfiler.file = cpu
	if err != nil {
		logger.Alert("+++++ Error creating file for CPU profile: %v ++++++\n", err)
		return err
	}
	logger.Alert("+++++ CPU profiling is active. Output file: %s +++++\n", cpuProfiler.CPUProfileFileName)
	pprof.StartCPUProfile(cpu)
	cpuProfiler.IsRunning = true
	return nil
}

// StopCPUProfile stops CPU profiling
func (cpuProfiler *CPUProfiler) StopCPUProfile() {
	if !cpuProfiler.IsRunning {
		logger.Warn("CPU profiler is not running. Nothing to stop\n")
		return
	}
	pprof.StopCPUProfile()
	cpuProfiler.file.Close()
	cpuProfiler.IsRunning = false
	logger.Alert("+++++ CPU profiling is finished. Output file:% s  +++++\n", cpuProfiler.CPUProfileFileName)
}
