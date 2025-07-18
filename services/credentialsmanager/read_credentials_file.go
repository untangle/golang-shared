package credentialsmanager

import (
	"encoding/json"
	"io"
)

// readFile reads the credentials file and saves the values
func (cm *credentialsManager) readFile() error {
	file, err := cm.fileSystem.Open(cm.fileLocation)
	if err != nil {
		cm.logger.Err("Error reading file at path %s: %w\n", cm.fileLocation, err)
		return err
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		cm.logger.Err("Error reading file at path %s: %w\n", cm.fileLocation, err)
		return err
	}

	credentials := map[string]string{}
	if err := json.Unmarshal(raw, &credentials); err != nil {
		cm.logger.Err("Error unmarshalling file at path %s: %w\n", cm.fileLocation, err)
		return err
	}

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.credentials = credentials

	return nil
}
