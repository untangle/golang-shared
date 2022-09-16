package settings

import (
	"fmt"
	"os"
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
