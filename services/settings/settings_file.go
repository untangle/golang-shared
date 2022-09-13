package settings

import (
	"fmt"
	"os"
	"sync"
)

type SettingsFile struct {
	filename string
	mutex    *sync.RWMutex
}

type SettingsOption func(*SettingsFile)

func WithLock(mutex *sync.RWMutex) SettingsOption {
	return func(file *SettingsFile) {
		file.mutex = mutex
	}
}

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
