package settings

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
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
