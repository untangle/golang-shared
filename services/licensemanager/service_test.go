package licensemanager

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceStart(t *testing.T) {
	// Monkey patching
	ServicesAllowedStatesLocation = "./testdata/allowedstates/"

	name := "test-name"
	serviceHook := ServiceHook{
		Start: func() {},
		Stop:  func() {},
	}

	state := ServiceState{
		Name:         name,
		AllowedState: StateEnable,
	}

	service := &Service{
		Name:  name,
		State: state,
		Hook:  serviceHook,
	}

	// A Start function is provided, and it's run
	// Expecting to not need to run sighup
	assert.False(t, service.ServiceStart())

	// Start function isn't provided. Sighup should be run
	service.Hook.Start = nil
	assert.True(t, service.ServiceStart())

	// Cleanup file test generated
	_ = os.Remove(fmt.Sprintf("%s%s", ServicesAllowedStatesLocation, name))
}

func TestServiceStop(t *testing.T) {
	// Monkey patching
	ServicesAllowedStatesLocation = "./testdata/allowedstates/"

	name := "test-name"
	serviceHook := ServiceHook{
		Start: func() {},
		Stop:  func() {},
	}

	state := ServiceState{
		Name:         name,
		AllowedState: StateEnable,
	}

	service := &Service{
		Name:  name,
		State: state,
		Hook:  serviceHook,
	}

	// A Start function is provided, and it's run
	// Expecting to not need to run sighup
	assert.False(t, service.ServiceStop())

	// Start function isn't provided. Sighup should be run
	service.Hook.Stop = nil
	assert.True(t, service.ServiceStop())

	// Cleanup file test generated
	_ = os.Remove(fmt.Sprintf("%s%s", ServicesAllowedStatesLocation, name))
}

func TestSetServiceState(t *testing.T) {
	// Monkey patching
	ServicesAllowedStatesLocation = "./testdata/allowedstates/"

	// Use the start and stop functions to test if ServiceStart and ServiceStop
	// were called as expected. They aren't an interface, but if the funcs
	// aren't nil they wrap ServiceHook's Start and Stop
	serviceStartCalled := false
	serviceStopCalled := false

	name := "test-name"
	serviceHook := ServiceHook{
		Start: func() { serviceStartCalled = true },
		Stop:  func() { serviceStopCalled = true },
	}

	state := ServiceState{
		Name:         name,
		AllowedState: StateEnable,
	}

	service := &Service{
		Name:  name,
		State: state,
		Hook:  serviceHook,
	}

	// Set the AllowedState from enabled to disabled
	err := service.setServiceState(StateDisable, "")
	assert.NoError(t, err)
	assert.True(t, serviceStopCalled)
	assert.False(t, serviceStartCalled)

	// Set AllowedState from disabled to enabled
	service.State.AllowedState = StateDisable
	serviceStartCalled = false
	serviceStopCalled = false
	err = service.setServiceState(StateEnable, "")
	assert.NoError(t, err)
	assert.False(t, serviceStopCalled)
	assert.True(t, serviceStartCalled)
}
