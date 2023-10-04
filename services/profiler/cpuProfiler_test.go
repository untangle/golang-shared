package profiler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartAndStopCPUProfile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-cpu-profile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	cpuProfiler := &CPUProfiler{
		CPUProfileFileName: tempFile.Name(),
	}

	err = cpuProfiler.StartCPUProfile()
	if err != nil {
		t.Errorf("StartCPUProfile returned an error: %v", err)
	}

	assert.True(t, cpuProfiler.IsRunning, "CPUProfiler should be running")

	cpuProfiler.StopCPUProfile()
	assert.False(t, cpuProfiler.IsRunning, "CPUProfiler should not be running")
}

func TestStartAndStopCPUProfileWithoutFileName(t *testing.T) {
	cpuProfiler := &CPUProfiler{}

	err := cpuProfiler.StartCPUProfile()
	if err == nil {
		t.Error("StartCPUProfile should return an error when CPUProfileFileName is not specified")
	}

	assert.False(t, cpuProfiler.IsRunning, "CPUProfiler should not be running")

	cpuProfiler.StopCPUProfile() // Stopping when not running should not result in an error
	assert.False(t, cpuProfiler.IsRunning, "CPUProfiler should not be running")
}

func TestStopCPUProfileWithInvalidFileName(t *testing.T) {
	cpuProfiler := &CPUProfiler{
		CPUProfileFileName: "/nonexistent-directory/nonexistent-file.prof",
	}

	cpuProfiler.StopCPUProfile() // Stopping with an invalid file name should not result in an error
	assert.False(t, cpuProfiler.IsRunning, "CPUProfiler should not be running")
}

func TestMain(m *testing.M) {
	// Run tests with the default M.Run() function
	os.Exit(m.Run())
}
