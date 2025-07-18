package credentialsmanager

import (
	"io/fs"
	"sync"

	"github.com/untangle/golang-shared/logger"
	"github.com/untangle/golang-shared/plugins"
)

const fileLocation = "/etc/config/credentials.json"

// interface for the credentials manager service
type CredentialsManager interface {
	plugins.Plugin
	GetToken(key string) string
}

type credentialsManager struct {
	fileLocation string
	logger       logger.LoggerLevels
	credentials  map[string]string
	mutex        sync.Mutex
	fileSystem   fs.FS
}

// GetCredentialsManager creates a new manager and returns it
func NewCredentialsManager(logger logger.LoggerLevels, fs fs.FS) CredentialsManager {
	return &credentialsManager{
		fileLocation: fileLocation,
		logger:       logger,
		mutex:        sync.Mutex{},
	}
}

// Startup starts the credentials manager service by reading the credentials file
func (m *credentialsManager) Startup() error {
	m.logger.Info("Starting the credentials service\n")

	if err := m.readFile(); err != nil {
		m.logger.Err("Unable to start credentials service; assuming no credentials - %v\n", err)
		m.credentials = nil
	}

	return nil
}

// Shutdown shuts down the credentials manager service
func (m *credentialsManager) Shutdown() error {
	m.logger.Info("Shutting down the credentials service\n")

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.credentials = nil

	return nil
}

// Name returns the service name
func (m *credentialsManager) Name() string {
	return "Credentials Manager"
}

// GetToken returns the auth token found in the credentials file under the `key` field
func (m *credentialsManager) GetToken(key string) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	token, ok := m.credentials[key]
	if !ok {
		m.logger.OCWarn("Could not get token for key %s\n", "getTokenFailure", 100, key)
	}
	return token
}
