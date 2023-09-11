package licensemanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
