package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/logger"
)

var routineInfoWatcher = make(chan *RoutineInfo)
var activeRoutines = make(map[string]bool)
var activeRoutinesMutex = &sync.RWMutex{}
var monitorRelation = RoutineContextGroup{}

// RoutineInfo is a struct used to signal routine events
type RoutineInfo struct {
	Name   string
	Action routineAction
}

// routineAction is a constant that represents specific routine actions that occur
type routineAction int

const (
	start routineAction = iota
	err
	end
)

// RoutineContextGroup is a collection of contexts and cancels that are associated
// with one specific routine group (ie: all local reporting contexts)
// This can be used with CancelContexts to send a ctx.Done() to all routines
// within a specific group
type RoutineContextGroup struct {
	Name     string
	Contexts map[string]context.Context
	Cancels  map[string]context.CancelFunc
}

// Startup is called to startup the monitor service
func Startup(callbackErrorHandler func(rtInfo *RoutineInfo)) {
	logger.Info("Starting routine monitor service...\n")
	routineInfoWatcher = make(chan *RoutineInfo)
	activeRoutines = make(map[string]bool)
	monitorRelation = CreateRoutineContextRelation(context.Background(), "monitor", []string{"routineInfoWatcher"})

	go monitorRoutineEvents(monitorRelation.Contexts["routineInfoWatcher"], callbackErrorHandler)

}

// RoutineStarted is used when initializing a new goroutine and adding monitoring to that routine
func RoutineStarted(routineName string) {
	defer activeRoutinesMutex.Unlock()
	activeRoutinesMutex.Lock()
	logger.Debug("Start Routine called: %s \n", routineName)
	routineInfoWatcher <- &RoutineInfo{Name: routineName, Action: start}
	activeRoutines[routineName] = true
}

// RoutineEnd is a function to simplify how we can defer calling finishRoutine() at the top of a function,
// instead of having to always call it at the end of a routine
func RoutineEnd(routineName string) {
	defer activeRoutinesMutex.Unlock()
	activeRoutinesMutex.Lock()
	logger.Debug("Finish Routine called: %s \n", routineName)
	routineInfoWatcher <- &RoutineInfo{Name: routineName, Action: end}
	_, ok := activeRoutines[routineName]
	if ok {
		delete(activeRoutines, routineName)
	}

}

// RoutineError signals a routine error to the routineInfoWatcher channel
func RoutineError(routineName string) {
	defer activeRoutinesMutex.Unlock()
	activeRoutinesMutex.Lock()
	logger.Info("Error Routine called: %s \n", routineName)
	routineInfoWatcher <- &RoutineInfo{Name: routineName, Action: err}
	_, ok := activeRoutines[routineName]
	if ok {
		delete(activeRoutines, routineName)
	}
}

// CreateRoutineContextRelation will create a collection of routine contexts and cancels that are related to routineNames passed as input
// This object is used to store context relations among many contexts, allowing for closing all contexts at once if needed
func CreateRoutineContextRelation(ctx context.Context, name string, routineNames []string) RoutineContextGroup {
	returnCtxs := make(map[string]context.Context)
	returnCancels := make(map[string]context.CancelFunc)

	for _, rtName := range routineNames {
		var thisCtx, thisCleaner = context.WithCancel(ctx)

		returnCtxs[rtName] = thisCtx
		returnCancels[rtName] = thisCleaner

	}

	return RoutineContextGroup{Name: name, Contexts: returnCtxs, Cancels: returnCancels}
}

// CancelContexts is a routine that will wait for a cancel call at the top level context
// and then deliver a cancel call to all downstream contexts passed in the cancels param
func CancelContexts(relation RoutineContextGroup) {
	var wg sync.WaitGroup

	for _, f := range relation.Cancels {
		wg.Add(1)
		go func(f func(), wgs *sync.WaitGroup) {
			defer wgs.Done()
			f()
		}(f, &wg)
	}

	wg.Wait()
	logger.Info("All routines associated with %s should be cancelled...\n", relation.Name)

}

// Shutdown is called to close the routine monitor
func Shutdown() {
	// Set shutdown channel for this service
	logger.Info("Shutting down monitor service...\n")
	CancelContexts(monitorRelation)
}

// montitorRoutineEvents is a routine that monitors the routineInfoWatcher queue for any routine events to act on
// THIS IS A ROUTINE FUNCTION
func monitorRoutineEvents(ctx context.Context, callbackErrorHandler func(rtInfo *RoutineInfo)) {

	// Read the routineInfoWatcher channel for any Error types
	for {
		select {
		case rtEvt := <-routineInfoWatcher:
			if rtEvt.Action == err {
				logger.Info("Acting on this event: %v\n", rtEvt)
				callbackErrorHandler(rtEvt)
			}
		case <-time.Tick(60 * time.Second):
			activeRoutinesMutex.Lock()
			logger.Info("There are %v monitored routines.\n", len(activeRoutines))
			activeRoutinesMutex.Unlock()
		case <-ctx.Done():
			logger.Info("Stopping routine monitor\n")
			return
		}
	}
}
