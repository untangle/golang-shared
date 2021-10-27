package profiler

import (
	"os"
	"runtime/pprof"

	"github.com/untangle/golang-shared/services/logger"
)

// CPUProfiler struct wraps functionality for profiling the CPU
type CPUProfiler struct {
	CPUProfileFileName string
}

// StartCPUProfile sets up and starts cpu profiling
func (cpuProfiler *CPUProfiler) StartCPUProfile() error {
	cpu, err := os.Create(cpuProfiler.CPUProfileFileName)
	if err != nil {
		logger.Alert("+++++ Error creating file for CPU profile:%v ++++++\n", err)
		return err
	}
	logger.Alert("+++++ CPU profiling is active. Output file:%s +++++\n", cpuProfiler.CPUProfileFileName)
	pprof.StartCPUProfile(cpu)
	return nil
}

// StopCPUProfile stops CPU profiling
func (cpuProfiler *CPUProfiler) StopCPUProfile() {
	pprof.StopCPUProfile()
	logger.Alert("+++++ CPU profiling is finished. Output file:%s  +++++\n", cpuProfiler.CPUProfileFileName)
}
