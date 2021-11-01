package licensemanager

import (
	"encoding/json"
	"io/ioutil"
	"os"

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

// ReadCommandFileAndGetStatus TODO
func ReadCommandFileAndGetStatus(name string) (bool, error) {
	var state ServiceState
	filename := ServicesAllowedStatesLocation + name
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warn("Not able to find service state file %s: %s\n", name, err.Error())
		return false, err
	}

	err = json.Unmarshal(fileContent, &state)
	if err != nil {
		logger.Warn("Not able to read content of service command file.%v \n ", err)
		return false, err
	}

	// remove file
	err = os.Remove(filename)
	if err != nil {
		logger.Warn("Failure removing file %s, continueing anyways: %s\n", filename, err.Error())
	}

	return state.getAllowedState() == StateEnable, nil
}

// writeOutServiceToEnableOrDisable() TODO
func (state *ServiceState) writeOutServiceToEnableOrDisable() error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	fileLocation := ServicesAllowedStatesLocation + state.Name
	logger.Debug("Location of service command: %s\n", fileLocation)
	err = ioutil.WriteFile(fileLocation, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (state *ServiceState) setAllowedState(newState State) {
	state.AllowedState = newState
}

func (state *ServiceState) getAllowedState() State {
	return state.AllowedState
}
