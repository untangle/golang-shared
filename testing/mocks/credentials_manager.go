package mocks

import (
	"github.com/stretchr/testify/mock"
)

type mockCredentialsManager struct {
	mock.Mock
}

func NewMockCredentialsManager() *mockCredentialsManager {
	return &mockCredentialsManager{}
}

func (m *mockCredentialsManager) Startup() error  { return nil }
func (m *mockCredentialsManager) Shutdown() error { return nil }
func (m *mockCredentialsManager) GetToken(key string) string {
	// this is required for mocking this function's call result
	args := m.Called(key)
	return args.String(0)
}
func (m *mockCredentialsManager) Name() string {
	return "Mocked Credentials Manager"
}
