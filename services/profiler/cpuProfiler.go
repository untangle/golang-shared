package profiler

import (
	"errors"
	"os"
	"runtime/pprof"

	"github.com/untangle/golang-shared/services/logger"
)

// CPUProfiler struct wraps functionality for profiling the CPU
type CPUProfiler struct {
	CPUProfileFileName string
	IsRunning          bool
}

// StartCPUProfile sets up and starts cpu profiling
func (cpuProfiler *CPUProfiler) StartCPUProfile() error {
	if cpuProfiler.CPUProfileFileName == "" {
		return errors.New("Cannot start cpu profiling. CPUProfileFileName must be specified!\n")
	}
	cpu, err := os.Create(cpuProfiler.CPUProfileFileName)
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
	cpuProfiler.IsRunning = false
	logger.Alert("+++++ CPU profiling is finished. Output file:% s  +++++\n", cpuProfiler.CPUProfileFileName)
}
