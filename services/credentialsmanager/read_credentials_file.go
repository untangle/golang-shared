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
		return err
	}

	var credentials credentialsFile
	if err := json.Unmarshal(raw, &credentials); err != nil {
		return err
	}

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.credentials = &credentials

	return nil
}
