package mocks

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

// NewMockLogger returns a new logger mock.
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

type LoggerHelper struct {
	MockLogger
	test *testing.T
}

func InitLoggerHelper(loggerHelper *LoggerHelper, test *testing.T) {
	loggerHelper.test = test
}

// MockPassthrough sets up the mock to allow any calls to any
// interface functions.
func (m *MockLogger) MockPassthrough() {
	m.On("Emerg", mock.Anything, mock.Anything, mock.Anything).Return(nil)
}

// Emerg mock method will return nil
func (m *MockLogger) Emerg(format string, args ...interface{}) {
	m.On("Emerg", mock.Anything, mock.Anything).Return(nil)
}

// Alert mock method will return nil
func (m *MockLogger) Alert(format string, args ...interface{}) {
	m.On("Alert", mock.Anything, mock.Anything).Return(nil)
}

// Crit mock method will return nil
func (m *MockLogger) Crit(format string, args ...interface{}) {
	m.On("Crit", mock.Anything, mock.Anything).Return(nil)
}

// Err mock method will return nil
func (m *MockLogger) Err(format string, args ...interface{}) {
	m.On("Err", mock.Anything, mock.Anything).Return(nil)
}

// Warn mock method will return nil
func (m *MockLogger) Warn(format string, args ...interface{}) {
	m.On("Warn", mock.Anything, mock.Anything).Return(nil)
}

// Notice mock method will return nil
func (m *MockLogger) Notice(format string, args ...interface{}) {
	m.On("Notice", mock.Anything, mock.Anything).Return(nil)
}

// Info mock method will return nil
func (m *MockLogger) Info(format string, args ...interface{}) {
	m.On("Info", mock.Anything, mock.Anything).Return(nil)
}

// Debug mock method will return nil
func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.On("Debug", mock.Anything, mock.Anything).Return(nil)
}

// Trace mock method will return nil
func (m *MockLogger) Trace(format string, args ...interface{}) {
	m.On("Trace", mock.Anything, mock.Anything).Return(nil)
}

// Trace mock method will return nil
func (m *MockLogger) OCWarn(format string, name string, limit int64, args ...interface{}) {
	m.Called()
	m.On("OCWarn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
}
