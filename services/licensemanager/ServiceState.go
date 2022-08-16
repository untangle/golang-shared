package licensemanager

import (
	"encoding/json"
	"io/ioutil"

	"github.com/untangle/golang-shared/services/logger"
)

const (
	// ServicesAllowedStatesLocation is the location where we put where services are enabled/disabled
	ServicesAllowedStatesLocation = "/etc/config/"
)

// ServiceState is used for setting the service state
type ServiceState struct {
	Name         string `json:"name"`
	AllowedState State  `json:"allowedState"`
}

// ReadCommandFileAndGetStatus reads a given command file and gets the status
// @param name - service to look for command file for
// @return bool - if service is enabled or not
// @return error - any associated error, nil if none
func ReadCommandFileAndGetStatus(name string) (bool, error) {
	var state ServiceState
	filename := ServicesAllowedStatesLocation + name

	// read file
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warn("Not able to find service state file %s: %s\n", name, err.Error())
		return false, err
	}

	// put into ServiceState struct
	err = json.Unmarshal(fileContent, &state)
	if err != nil {
		logger.Warn("Not able to read content of service command file: %s \n ", err.Error())
		return false, err
	}

	// return if StateEnable
	return state.getAllowedState() == StateEnable, nil
}

// writeOutServiceToEnableOrDisable writes out the state of a ServiceState
// @return error - any associated error, nil if none
func (state *ServiceState) writeOutServiceToEnableOrDisable() error {
	// get data of ServiceState
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	// write out file
	fileLocation := ServicesAllowedStatesLocation + state.Name
	logger.Debug("Location of service command: %s\n", fileLocation)
	err = ioutil.WriteFile(fileLocation, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// setAllowedState sets the AllowedState of a ServiceState
// @param State newState - new allowed state
func (state *ServiceState) setAllowedState(newState State) {
	state.AllowedState = newState
}

// getAllowedState returns the AllowedState of a ServiceState
// @return State - serviceState being returned
func (state *ServiceState) getAllowedState() State {
	return state.AllowedState
}