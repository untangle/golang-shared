package credentialsmanager

import (
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/logger"
)

func TestCredntialsManager(t *testing.T) {
	m := NewCredentialsManager(logger.GetLoggerInstance()).(*credentialsManager)

	testBadFileStartup(t, m)
	testGoodFileStartup(t, m)
	testAlertsToken(t, m)
	testCloudReportToken(t, m)
	testShutdown(t, m)
	testNoValuesAfterShutdown(t, m)
}

// testBadFileStartup assert that a bad file path prevents startup
func testBadFileStartup(t *testing.T, m *credentialsManager) {
	m.fileLocation = "/some/path/that/should/not/exist.json"

	err := m.Startup()
	assert.NotNil(t, err, "Startup bad file")
}

// testGoodFileStartup assert that it starts when the file exists and is in the right format
func testGoodFileStartup(t *testing.T, m *credentialsManager) {
	abs, err := filepath.Abs("./test_files/test_credentials.json")
	assert.Nil(t, err)

	m.fileLocation = abs

	err = m.Startup()
	assert.Nil(t, err, "Startup good file")
}

// testAlertsToken assert it returns the alert token
func testAlertsToken(t *testing.T, m *credentialsManager) {
	token := m.GetAlertsAuthToken()
	assert.Equal(t, "a13R-T5A-uTh-T0k-3N", token, "GetAlertsAuthToken")
}

// testCloudReportToken asert it returns the cloud reporting token
func testCloudReportToken(t *testing.T, m *credentialsManager) {
	token := m.GetCloudReportingAuthToken()
	assert.Equal(t, "CL0UDR-3P0R-T1NG-AUTH-T0K3N", token, "GetAlertsAuthToken")
}

// testShutdown assert it shuts down properly
func testShutdown(t *testing.T, m *credentialsManager) {
	err := m.Shutdown()
	assert.Nil(t, err, "Shutdown")
}

// testNoValuesAfterShutdown it should return no values after shutdown
func testNoValuesAfterShutdown(t *testing.T, m *credentialsManager) {
	at := m.GetAlertsAuthToken()
	assert.Equal(t, "", at, "GetAlertsAuthToken after shutdown")
	crt := m.GetCloudReportingAuthToken()
	assert.Equal(t, "", crt, "GetCloudReportingAuthToken after shutdown")
}
