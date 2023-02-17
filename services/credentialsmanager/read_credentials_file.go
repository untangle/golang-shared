package credentialsmanager

import (
	"encoding/json"
	"io/ioutil"
)

// the type of the credential file's content
type credentialsFile struct {
	AlertsAuthToken         string `json:"alertsAuthToken"`
	CloudReportingAuthToken string `json:"cloudReportingAuthToken"`
}

// readFile reads the credentials file and saves the values
func (cm *credentialsManager) readFile() error {
	raw, err := ioutil.ReadFile(cm.fileLocation)
	if err != nil {
		cm.logger.Err("Error reading file at path %s: %s\n", cm.fileLocation, err)
		return err
	}

	var credentials credentialsFile
	if err := json.Unmarshal(raw, &credentials); err != nil {
		cm.logger.Err("Error unmarshalling file at path %s: %s", cm.fileLocation, err)
		return err
	}

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.credentials = &credentials

	return nil
}
