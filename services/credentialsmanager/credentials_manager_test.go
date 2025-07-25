package credentialsmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/logger"
)

func TestCredentialsManager(t *testing.T) {
	m := NewCredentialsManager(logger.GetLoggerInstance(), os.DirFS(".")).(*credentialsManager)

	testBadFileStartup(t, m)
	testGoodFileStartup(t, m)
	testAlertsToken(t, m)
	testCloudReportToken(t, m)
	testBadKeyToken(t, m)
	testShutdown(t, m)
	testNoValuesAfterShutdown(t, m)
}

// testBadFileStartup assert that a bad file path prevents startup
func testBadFileStartup(t *testing.T, m *credentialsManager) {
	m.fileLocation = "some/path/that/should/not/exist.json"

	err := m.Startup()
	assert.Nil(t, err, "Startup bad file")
}

// testGoodFileStartup assert that it starts when the file exists and is in the right format
func testGoodFileStartup(t *testing.T, m *credentialsManager) {
	m.fileLocation = "test_files/test_credentials.json"

	err := m.Startup()
	assert.Nil(t, err, "Startup good file")
}

// testAlertsToken assert it returns the alert token
func testAlertsToken(t *testing.T, m *credentialsManager) {
	token := m.GetToken("alertsAuthToken")
	assert.Equal(t, "a13R-T5A-uTh-T0k-3N", token, "alertsAuthToken")
}

// testCloudReportToken asert it returns the cloud reporting token
func testCloudReportToken(t *testing.T, m *credentialsManager) {
	token := m.GetToken("cloudReportingAuthToken")
	assert.Equal(t, "CL0UDR-3P0R-T1NG-AUTH-T0K3N", token, "cloudReportingAuthToken")
}

func testBadKeyToken(t *testing.T, m *credentialsManager) {
	token := m.GetToken("vErYrEaLtOkEn")
	assert.Equal(t, "", token, "vErYrEaLtOkEn")
}

// testShutdown assert it shuts down properly
func testShutdown(t *testing.T, m *credentialsManager) {
	err := m.Shutdown()
	assert.Nil(t, err, "Shutdown")
}

// testNoValuesAfterShutdown it should return no values after shutdown
func testNoValuesAfterShutdown(t *testing.T, m *credentialsManager) {
	at := m.GetToken("alertsAuthToken")
	assert.Equal(t, "", at, `GetToken("alertsAuthToken")`)
	crt := m.GetToken("cloudReportingAuthToken")
	assert.Equal(t, "", crt, `GetToken("cloudReportingAuthToken")`)
	vrt := m.GetToken("vErYrEaLtOkEn")
	assert.Equal(t, "", vrt, `GetToken("vErYrEaLtOkEn")`)
}
