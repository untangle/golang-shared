package signalhandler

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSIGTERMSignalHandler(t *testing.T) {
	// Create a new SignalHandler
	sh := NewSignalHandler()

	// Check the shutdownFlag value before starting HandleSignals() function
	assert.Equal(t, sh.shutdownFlag, uint32(0))

	// Simulate a SIGTERM signal after a delay
	go func() {
		time.Sleep(1 * time.Second)
		sh.SetShutdownFlag()
	}()

	// Wait for the shutdown to complete or timeout
	select {
	case <-sh.ShutdownChannel:
		logger.Warn("Shutdown completed.\n")
	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for shutdown.")
	}

	// Ensure that the shutdownFlag is set
	assert.Equal(t, sh.GetShutdownFlag(), true)
}

func TestSIGQUITSignalHandler(t *testing.T) {
	// Create a new SignalHandler
	sh := NewSignalHandler()

	// Check the shutdownFlag value before starting HandleSignals() function
	assert.Equal(t, sh.shutdownFlag, uint32(0))

	// Simulate a SIGQUIT signal after a delay
	go func() {
		time.Sleep(1 * time.Second)
		sh.dumpStack()
	}()

	fileInfo, err := os.Stat("/tmp/reportd.stack")
	if os.IsNotExist(err) {
		assert.Error(t, err, "File not found at /tmp/reportd.stack")
	} else {
		expectedFileName := "reportd.stack"

		// check File name
		fileName := fileInfo.Name()
		assert.Equal(t, expectedFileName, fileName)

		// check File size (in bytes) if it is greater than zero
		fileSize := fileInfo.Size()
		assert.True(t, fileSize > 0, "File should not be empty")
	}
}
