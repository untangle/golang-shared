package mocks

import "github.com/untangle/golang-shared/services/credentialsmanager"

type mockCredentialsManager struct{}

func NewMockCredentialsManager() credentialsmanager.CredentialsManager {
	return &mockCredentialsManager{}
}

func (m *mockCredentialsManager) Startup() error                     { return nil }
func (m *mockCredentialsManager) Shutdown() error                    { return nil }
func (m *mockCredentialsManager) GetAlertsAuthToken() string         { return "" }
func (m *mockCredentialsManager) GetCloudReportingAuthToken() string { return "" }

func (cm *mockCredentialsManager) Name() string {
	return "Mocked Credentials Manager"
}
