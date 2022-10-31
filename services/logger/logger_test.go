package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(TestLogger))
}

type MockConfigFile struct {
	mock.Mock
	logger *Logger
}

type TestLogger struct {
	suite.Suite
	logger Logger
	write  LogWriter
}

func (m *MockConfigFile) MockLoadConfigFromFile(logger *Logger) {
	logger.config.LogLevelMap = map[string]LogLevel{
		"Emergtest":  {"EMERG", 0},
		"Alerttest":  {"ALERT", 1},
		"Crittest":   {"CRIT", 2},
		"Errtest":    {"ERROR", 3},
		"Warntest":   {"WARN", 4},
		"Noticetest": {"NOTICE", 5},
		"Infotest":   {"INFO", 6},
		"Debugtest":  {"DEBUG", 7},
		"Tracetest":  {"TRACE", 8},
	}
}

func (suite *TestLogger) TestStartup() {
	suite.logger.Startup()
	assert.Equal(suite.T(), nil, suite.logger.config.OutputWriter)
	suite.logger.config.FileLocation = "LoggerConfig.json"
	suite.logger.Startup()
	var MockWriter bytes.Buffer
	suite.logger.config.OutputWriter = &MockWriter
	suite.logger.Startup()
}

func (suite *TestLogger) TestName() {
	assert.Equal(suite.T(), "logger", suite.logger.Name())
}

func (suite *TestLogger) TestIsLogEnabledSourceTrue() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	tests := []struct {
		level  int32
		source string
		output bool
	}{
		{
			level:  5,
			source: "Tracetest",
			output: true,
		},
		{
			level:  6,
			source: "Infotest",
			output: true,
		},
		{
			level:  7,
			source: "Emergtest",
			output: false,
		},
		{
			level:  7,
			source: "test", //test isn't present in loglevelmap
			output: false,
		},
	}
	for _, tt := range tests {
		assert.Equal(suite.T(), tt.output, suite.logger.IsLogEnabledSource(tt.level, tt.source))
	}
}

func (suite *TestLogger) TestEnabled() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	assert.Equal(suite.T(), true, suite.logger.IsEmergEnabled())
	assert.Equal(suite.T(), false, suite.logger.IsDebugEnabled())
	assert.Equal(suite.T(), true, suite.logger.IsInfoEnabled())
	assert.Equal(suite.T(), true, suite.logger.IsNoticeEnabled())
	assert.Equal(suite.T(), true, suite.logger.IsWarnEnabled())
	assert.Equal(suite.T(), true, suite.logger.IsErrEnabled())
	assert.Equal(suite.T(), true, suite.logger.IsCritEnabled())
	assert.Equal(suite.T(), true, suite.logger.IsAlertEnabled())
	assert.Equal(suite.T(), false, suite.logger.IsTraceEnabled())

	suite.logger.Shutdown()
}

func (suite *TestLogger) TestFindLogLevelName() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	assert.Equal(suite.T(), "UNDEFINED", FindLogLevelName(-1))
	assert.Equal(suite.T(), "11", FindLogLevelName(11))
	assert.Equal(suite.T(), "TRACE", FindLogLevelName(8))
}

func (suite *TestLogger) TestFindLogLevelValue() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	assert.Equal(suite.T(), int32(8), FindLogLevelValue("TRACE"))
	assert.Equal(suite.T(), int32(-1), FindLogLevelValue("test"))
	assert.Equal(suite.T(), int32(6), FindLogLevelValue("INFO"))
}

func (suite *TestLogger) TestWrite() {
	int_result, error_result := suite.write.Write([]byte("test\n"))
	assert.Equal(suite.T(), 5, int_result)
	assert.Equal(suite.T(), nil, error_result)
}

func (suite *TestLogger) TestDefaultLogWriter() {
	assert.Equal(suite.T(), &LogWriter{buffer: []uint8{}, source: "System"}, DefaultLogWriter("System"))
}

func (suite *TestLogger) TestLoadConfigFromFile() {
	assert.Equal(suite.T(), []uint8([]byte(nil)), suite.logger.config.LoadConfigFromFile())
	suite.logger.config.FileLocation = "LoggerConfig.json"
	assert.Equal(suite.T(), 791, len(suite.logger.config.LoadConfigFromFile()))
}

func (suite *TestLogger) TestLoadConfigFromJSON() {
	mockstring := `"certcache": "INFO",
	"certfetch": "INFO",
	"certmanager": "INFO",
	"certsniff": "INFO",
	"classify": "INFO",
	"common": "INFO",
	"conntrack": "INFO",
	"device": "INFO",
	"dict": "INFO",
	"discovery": "INFO",
	"dispatch": "INFO",
	"dns": "INFO",
	"example": "INFO",
	"geoip": "INFO",
	"gin": "INFO",
	"kernel": "INFO",
	"logger": "INFO"`
	Mockdata := []byte(mockstring)
	suite.logger.config.LoadConfigFromJSON(Mockdata)
}
