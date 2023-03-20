package licensemanager

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test writing a service file and reading it back in
// Could be done independently, but  ends up creating
// a lot of duplicate code
func TestWriteAndReadServiceFile(t *testing.T) {
	ServicesAllowedStatesLocation = "./testdata/allowedstates/"
	testFile := "test-file"

	serviceState := &ServiceState{
		Name:         testFile,
		AllowedState: StateEnable,
	}

	err := serviceState.writeOutServiceToEnableOrDisable()
	assert.NoError(t, err)

	enabled, err := ReadCommandFileAndGetStatus(testFile)
	assert.NoError(t, err)
	assert.True(t, enabled)

	_ = os.Remove(fmt.Sprintf("%s%s", ServicesAllowedStatesLocation, testFile))
}

func TestSetAllowedState(t *testing.T) {
	// Enable service, then disable it.
	// Check if it was successfully disabled
	serviceState := &ServiceState{
		AllowedState: StateEnable,
	}

	expected := StateDisable
	serviceState.setAllowedState(expected)

	assert.Equal(t, expected, serviceState.AllowedState)
}

func TestGetAllowedState(t *testing.T) {
	expected := StateEnable
	serviceState := &ServiceState{
		AllowedState: expected,
	}
	assert.Equal(t, expected, serviceState.getAllowedState())

	expected = StateDisable
	serviceState.AllowedState = expected

	assert.Equal(t, expected, serviceState.getAllowedState())
}
