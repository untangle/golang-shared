package settings

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/untangle/golang-shared/util"
)

// SettingsFile is an object representing the system-wide settings and
// operations on them, which can be configured to point to a local
// file. It is best used as a singleton so that the locking works and
// prevents concurrent writes/reads.
type SettingsFile struct {
	// filename for the settings file.
	filename string

	// Mutex to lock the file.
	mutex *sync.RWMutex
}

// SettingsOption is an option for the constructor of SettingsFile.
type SettingsOption func(*SettingsFile)

// WithLock will use the given lock in the SettingsFile to lock for
// reads/writes.
func WithLock(mutex *sync.RWMutex) SettingsOption {
	return func(file *SettingsFile) {
		file.mutex = mutex
	}
}

// NewSettingsFile is a constructor for SettingsFile, give it filename
// as a path to the settings file, and any supplemental options you
// want.
func NewSettingsFile(filename string, opts ...SettingsOption) *SettingsFile {
	file := &SettingsFile{
		filename: filename,
	}
	for _, opt := range opts {
		opt(file)
	}
	if file.mutex == nil {
		file.mutex = &sync.RWMutex{}
	}
	return file
}

// UnmarshalSettingsAtPath wraps the PathUnmarshaller object, taking
// out a read lock on the file object's lock first.
func (file *SettingsFile) UnmarshalSettingsAtPath(value interface{}, settings ...string) error {
	file.mutex.RLock()
	defer file.mutex.RUnlock()
	reader, err := os.Open(file.filename)
	if err != nil {
		return fmt.Errorf("settings file: unable to open file %s: %w",
			file.filename,
			err)
	}
	unmarshaller := NewPathUnmarshaller(reader)
	return unmarshaller.UnmarshalAtPath(value, settings...)
}

// Generates a backup of a settings file using a provided script. Locks the settings file before generation.
// Returns the file name as the full path to it, the settings file data as []byte, and an error. If any error occurs
// "", nil, err will be returned.
// The script provided must output a line specifying the location of the generated backup file in the format:
// 	Backup location: <full path of file> \n
func (file *SettingsFile) GenerateBackupFile(backupGenerationScript string, scriptArgs ...string) (string, []byte, error) {
	file.mutex.RLock()
	defer file.mutex.RUnlock()

	cmd, err := exec.Command(backupGenerationScript, scriptArgs...).Output()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create the settings file with command %s %v", backupGenerationScript, scriptArgs)
	}

	scanner := bufio.NewScanner(bytes.NewReader(cmd))
	var settingsFile string
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Backup location:") {
			settingsFile = strings.Trim(strings.Split(scanner.Text(), ": ")[1], " \n")
		}
	}
	if settingsFile == "" {
		return "", nil, fmt.Errorf("failed to create the default settings file")
	}

	fileData, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read the default settings file")
	}

	return settingsFile, fileData, nil
}

// Returns a JSON structure(map[string]interface{}) of the current settings
func (file *SettingsFile) GetAllSettings() (map[string]interface{}, error) {
	file.mutex.RLock()
	defer file.mutex.RUnlock()

	raw, err := ioutil.ReadFile(file.filename)
	if err != nil {
		return nil, err
	}
	var jsonObject interface{}
	err = json.Unmarshal(raw, &jsonObject)
	if err != nil {
		return nil, err
	}
	j, ok := jsonObject.(map[string]interface{})
	if ok {
		return j, nil
	}

	return nil, errors.New("invalid settings file format")
}

// SetSettings updates the settings. Calls lock/unlock on the SettingsFile's mutex
func (file *SettingsFile) SetSettings(segments []string, value interface{}, force bool) (interface{}, error) {
	var ok bool
	var err error
	var jsonSettings map[string]interface{}
	var newSettings interface{}

	jsonSettings, err = file.GetAllSettings()
	if err != nil {
		return createJSONErrorObject(err), err
	}

	newSettings, err = setSettingsInJSON(jsonSettings, segments, value)
	if err != nil {
		return createJSONErrorObject(err), err
	}
	jsonSettings, ok = newSettings.(map[string]interface{})
	if !ok {
		err = errors.New("invalid global settings object")
		return createJSONErrorObject(err), err
	}

	file.mutex.Lock()
	output, err := syncAndSave(jsonSettings, file.filename, force)
	file.mutex.Unlock()
	if err != nil {
		var errJSON map[string]interface{}
		marshalErr := json.Unmarshal([]byte(err.Error()), &errJSON)
		if marshalErr != nil {
			logger.Warn("Failed to marshal into json: %s\n", marshalErr.Error())
			if strings.Contains(err.Error(), "CONFIRM") {
				return determineSetSettingsError(err, output, file.filename, jsonSettings)
			}
		} else {
			if _, ok := errJSON["CONFIRM"]; ok {
				return determineSetSettingsError(err, output, file.filename, jsonSettings)
			}
		}
		logger.Warn("Failed to save settings: %s\n", err.Error())
		responseErr := err.Error()
		if len(responseErr) == 0 {
			responseErr = "failed_sync_settings"
		}
		return map[string]interface{}{"error": responseErr, "output": output}, err
	}

	return map[string]interface{}{"output": output}, err
}

// Restores settings from a backups file. The backup file should be in the form of a tar.gz with structure
// /<directory named after date/time created>/settings.json. Initially, backups were restored with just the
// settings.json file, so for the time being the old settings.json backups are still supported.
func (file *SettingsFile) RestoreSettingsFromFile(fileData []byte, exceptions ...string) (interface{}, error) {
	// Set settings data with what could potentially be a JSON string.
	// If it's not, the settings data will get swapped out for what was in
	// a tar file
	settingsData := fileData

	// A user uploaded settings file was originally a single JSON file. For
	// backwards compatibility, if the file uploaded wasn't a JSON
	// try treating it as a JSON
	fileName := "settings.json"
	foundFiles, err := util.ExtractFilesFromTar(fileData, true, fileName)
	if err != nil {
		logger.Warn("Failed to extract the settings restore file as a tar. Attempting to use the settings restore file as a JSON\n")
	} else {
		if data, ok := foundFiles[fileName]; ok {
			settingsData = data
			logger.Debug("Retrieved settings restore file from tar.")
		}
	}

	var settingsJson map[string]interface{}
	if err := json.Unmarshal(settingsData, &settingsJson); err != nil {
		return createJSONErrorObject(err), err
	}

	return file.SetAllSettingsWithExceptions(settingsJson, exceptions...)
}

// Updates settings with the new settings passed in. newSettings needs to be a valid
// 	Json structure(map[string]interface{}) of all the settings. For each exception, the current settings will be
// 	used instead of the what was in newSettings. Returns an error if something went wrong, along
// 	with an error JSON. If the settings were set, no error will be returned and a JSON response
//  object will be. !!!Only works for settings at the highest level in the settings json
func (file *SettingsFile) SetAllSettingsWithExceptions(newSettings map[string]interface{}, exceptions ...string) (interface{}, error) {
	currentSettings, err := file.GetAllSettings()
	if err != nil {
		return createJSONErrorObject(err), err
	}

	for _, exception := range exceptions {
		newSettings[exception] = currentSettings[exception]
	}

	return file.SetSettings(nil, newSettings, true)
}
