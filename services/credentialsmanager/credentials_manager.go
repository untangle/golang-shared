package credentialsmanager

import (
	"sync"

	"github.com/untangle/golang-shared/services/logger"
)

const fileLocation = "/etc/config/credentials.json"

// interface for the credentials manager service
type CredentialsManager interface {
	Startup() error
	Shutdown() error
	GetAlertsAuthToken() string
	GetCloudReportingAuthToken() string
	Name() string
}

type credentialsManager struct {
	fileLocation string
	logger       logger.LoggerLevels
	credentials  *credentialsFile
	mutex        sync.Mutex
}

// GetCredentialsManager creates a new manager and returns it
func NewCredentialsManager(logger logger.LoggerLevels) CredentialsManager {
	return &credentialsManager{
		fileLocation: fileLocation,
		logger:       logger,
		mutex:        sync.Mutex{},
	}
}

// Startup starts the credentials manager service by reading the credentials file
func (cm *credentialsManager) Startup() error {
	cm.logger.Info("Starting the credentials service\n")

	return cm.readFile()
}

// Shutdown shuts down the credentials manager service
func (cm *credentialsManager) Shutdown() error {
	cm.logger.Info("Shutting down the credentials service\n")

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.credentials = nil

	return nil
}

// Name returns the service name
func (cm *credentialsManager) Name() string {
	return "Credentials Manager"
}

// GetAlertsAuthToken returns the alerts authentication token, if present
func (cm *credentialsManager) GetAlertsAuthToken() string {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.credentials == nil {
		cm.logger.Err("GetAlertsAuthToken error: Credential configs are missing!\n")
		return ""
	}

	return cm.credentials.AlertsAuthToken
}

// GetCloudReportingAuthToken returns the cloud reporting authentication token, if present
func (cm *credentialsManager) GetCloudReportingAuthToken() string {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.credentials == nil {
		cm.logger.Err("GetCloudReportingAuthToken error: Credential configs are missing!\n")
		return ""
	}

	return cm.credentials.CloudReportingAuthToken
}
