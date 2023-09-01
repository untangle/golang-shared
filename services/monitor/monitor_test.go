package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoutineStarted(t *testing.T) {
	// Create context for testing
	ctx := context.Background()

	// create tests for storing the routines information
	tests := []struct {
		name                   string
		routineNames           []string
		numberOfActiveRoutines int
	}{
		{
			name:                   "test",
			routineNames:           []string{"routineInfoWatcher1", "routineInfoWatcher2", "routineInfoWatcher3", "routineInfoWatcher4"},
			numberOfActiveRoutines: 4,
		},
	}

	monitorRelation := CreateRoutineContextRelation(ctx, tests[0].name, tests[0].routineNames)

	routines := tests[0].routineNames
	for _, currentRoutineName := range routines {
		go monitorRoutineEvents(monitorRelation.Contexts[currentRoutineName], handleRoutineWatcherEvents)

		RoutineStarted(currentRoutineName)
		<-time.After(150 * time.Millisecond)

		activeRoutinesMutex.RLock()
		_, exist := activeRoutines[currentRoutineName]
		activeRoutinesMutex.RUnlock()

		// check current status of routine
		assert.True(t, exist)
	}

	// check number of the activeRoutines
	assert.Equal(t, len(activeRoutines), tests[0].numberOfActiveRoutines)

	// stop all started routines
	for _, routineName := range routines {
		RoutineEnd(routineName)
		<-time.After(100 * time.Millisecond)
	}
}

func TestRoutineEnd(t *testing.T) {
	// Create context for testing
	ctx := context.Background()

	// create tests for storing the routines information
	tests := []struct {
		name                   string
		routineNames           []string
		numberOfActiveRoutines int
	}{
		{
			name:                   "test",
			routineNames:           []string{"routineInfoWatcher1", "routineInfoWatcher2", "routineInfoWatcher3"},
			numberOfActiveRoutines: 3,
		},
	}

	monitorRelation := CreateRoutineContextRelation(ctx, tests[0].name, tests[0].routineNames)

	routines := tests[0].routineNames
	for _, routineName := range routines {
		go monitorRoutineEvents(monitorRelation.Contexts[routineName], handleRoutineWatcherEvents)

		RoutineStarted(routineName)
		<-time.After(150 * time.Millisecond)
	}

	RoutineEnd("routineInfoWatcher2")
	<-time.After(150 * time.Millisecond)

	activeRoutinesMutex.RLock()
	_, exist := activeRoutines["routineInfoWatcher2"]
	activeRoutinesMutex.RUnlock()

	// check current status of routine
	assert.False(t, exist)

	// check number of the activeRoutines
	assert.Equal(t, len(activeRoutines), tests[0].numberOfActiveRoutines-1)

	// stop all started routines
	for routineName, _ := range activeRoutines {
		RoutineEnd(routineName)
		<-time.After(100 * time.Millisecond)
	}
}

func TestRoutineError(t *testing.T) {
	// Create  context for testing
	ctx := context.Background()

	// create tests for storing the routines information
	tests := []struct {
		name                   string
		routineNames           []string
		numberOfActiveRoutines int
	}{
		{
			name:                   "test",
			routineNames:           []string{"routineInfoWatcher1", "routineInfoWatcher2", "routineInfoWatcher3"},
			numberOfActiveRoutines: 3,
		},
	}

	monitorRelation := CreateRoutineContextRelation(ctx, tests[0].name, tests[0].routineNames)

	routines := tests[0].routineNames
	for _, routineName := range routines {
		go monitorRoutineEvents(monitorRelation.Contexts[routineName], handleRoutineWatcherEvents)

		RoutineStarted(routineName)
		<-time.After(150 * time.Millisecond)
	}

	RoutineError("routineInfoWatcher3")
	<-time.After(150 * time.Millisecond)

	activeRoutinesMutex.RLock()
	_, exist := activeRoutines["routineInfoWatcher3"]
	activeRoutinesMutex.RUnlock()

	// check current status of routine
	assert.False(t, exist)

	// check number of the activeRoutines
	assert.Equal(t, len(activeRoutines), tests[0].numberOfActiveRoutines-1)

	for routineName, _ := range activeRoutines {
		RoutineEnd(routineName)
		<-time.After(100 * time.Millisecond)
	}
}

func handleRoutineWatcherEvents(rtEvt *RoutineInfo) {
	logger.Info("Taking action on %v event", rtEvt.Name)
}

func TestRoutineContextGroup(t *testing.T) {
	// Create context for testing
	ctx := context.Background()

	// create tests for storing the routines information
	tests := []struct {
		name                   string
		routineNames           []string
		numberOfActiveRoutines int
	}{
		{
			name:                   "test",
			routineNames:           []string{"routineInfoWatcher1", "routineInfoWatcher2", "routineInfoWatcher3"},
			numberOfActiveRoutines: 3,
		},
	}

	monitorRelation := CreateRoutineContextRelation(ctx, tests[0].name, tests[0].routineNames)

	routines := tests[0].routineNames
	for _, routineName := range routines {
		go monitorRoutineEvents(monitorRelation.Contexts[routineName], handleRoutineWatcherEvents)

		RoutineStarted(routineName)
		<-time.After(150 * time.Millisecond)
	}

	beforeShutdownRCG := monitorRelation
	Shutdown()
	afterShutdownRCG := monitorRelation

	// check RoutineContextGroup's status before and after shutdown
	assert.Equal(t, beforeShutdownRCG, afterShutdownRCG)

	// stop all started routines
	for _, routineName := range routines {
		RoutineEnd(routineName)
		<-time.After(100 * time.Millisecond)
	}
}
