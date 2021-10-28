package profiler

// CreateCPUProfiler creates a cpu profiler struct and returns a pointer to it
func CreateCPUProfiler(cpuProfileFilename string) *CPUProfiler {
	return &CPUProfiler{
		CPUProfileFileName: cpuProfileFilename,
	}
}
