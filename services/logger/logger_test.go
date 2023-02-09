package logger

import (
	"bytes"
	"encoding/json"
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
}

type TestLogger struct {
	suite.Suite
	configFile MockConfigFile
}

func (m *MockConfigFile) MockLoadConfigFromFile(logger *Logger) {
	logger.config = &LoggerConfig{}
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

// createTestConfig creates the logger config
func createTestConfig() LoggerConfig {
	return LoggerConfig{
		FileLocation: "/tmp/logconfig_test.json",
		LogLevelMap:  createTestMap(),
	}
}

func createTestMap() map[string]LogLevel {
	return map[string]LogLevel{
		"test1": {Name: "INFO"},
		"test2": {Name: "WARN"},
		"test3": {Name: "ERROR"},
		"test4": {Name: "DEBUG"},
		"test5": {Name: "INFO"},
	}
}

func (suite *TestLogger) SetupTest() {
	suite.configFile = MockConfigFile{}
}

func (suite *TestLogger) TestStartup() {
	logger := Logger{}
	logger.config = &LoggerConfig{}

	//Startup on a new logger will use the default config options
	logger.Startup()
	assert.Equal(suite.T(), DefaultLogWriter("system"), logger.config.OutputWriter)

	logger.config.FileLocation = "LoggerConfig.json"
	logger.Startup()

	var MockWriter bytes.Buffer
	logger.config.OutputWriter = &MockWriter
	logger.Startup()
}

//Test default service name
func (suite *TestLogger) TestName() {
	logger := Logger{}
	assert.Equal(suite.T(), "logger", logger.Name())
}

func (suite *TestLogger) TestIsLogEnabledSourceTrue() {
	logger := Logger{}
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&logger)
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
		assert.Equal(suite.T(), tt.output, logger.IsLogEnabledSource(tt.level, tt.source))
	}
}

func (suite *TestLogger) TestEnabled() {
	logger := Logger{}
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&logger)
	assert.Equal(suite.T(), true, logger.IsEmergEnabled())
	assert.Equal(suite.T(), false, logger.IsDebugEnabled())
	assert.Equal(suite.T(), true, logger.IsInfoEnabled())
	assert.Equal(suite.T(), true, logger.IsNoticeEnabled())
	assert.Equal(suite.T(), true, logger.IsWarnEnabled())
	assert.Equal(suite.T(), true, logger.IsErrEnabled())
	assert.Equal(suite.T(), true, logger.IsCritEnabled())
	assert.Equal(suite.T(), true, logger.IsAlertEnabled())
	assert.Equal(suite.T(), false, logger.IsTraceEnabled())

	logger.Shutdown()
}

func (suite *TestLogger) TestFindLogLevelName() {
	logger := Logger{}
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&logger)
	assert.Equal(suite.T(), "UNDEFINED", FindLogLevelName(-1))
	assert.Equal(suite.T(), "11", FindLogLevelName(11))
	assert.Equal(suite.T(), "TRACE", FindLogLevelName(8))
}

func (suite *TestLogger) TestFindLogLevelID() {
	logger := Logger{}
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&logger)
	traceLevel := LogLevel{Name: "TRACE"}
	badLevel := LogLevel{Name: "test"}
	infoLevel := LogLevel{Name: "INFO"}

	assert.Equal(suite.T(), int32(8), traceLevel.GetId())
	assert.Equal(suite.T(), int32(-1), badLevel.GetId())
	assert.Equal(suite.T(), int32(6), infoLevel.GetId())
}

func (suite *TestLogger) TestWrite() {
	logger := NewLogger()
	int_result, error_result := logger.config.OutputWriter.Write([]byte("test\n"))
	assert.Equal(suite.T(), 5, int_result)
	assert.Equal(suite.T(), nil, error_result)
}

func (suite *TestLogger) TestDefaultLogWriter() {
	assert.Equal(suite.T(), &LogWriter{buffer: []uint8{}, source: "System", logLevel: NewLogLevel("INFO")}, DefaultLogWriter("System"))
}

func (suite *TestLogger) TestLoadConfigFromFile() {
	logger := NewLogger()
	//Test load from default file that may or may not exist
	assert.Equal(suite.T(), []uint8([]byte(nil)), logger.config.LoadConfigFromFile())
	// Test that load config from file works
	logger.config.FileLocation = "LoggerConfig.json"
	assert.Equal(suite.T(), 791, len(logger.config.LoadConfigFromFile()))

	logger.getLogLevel("discovery", "discovery")
	//Test that the LoggerConfig.json matches some properties
	assert.Equal(suite.T(), int32(6), logger.getLogLevel("discovery", "discovery"))

}

func (suite *TestLogger) TestLoadConfigFromJSON() {
	loggerConf := LoggerConfig{}

	testMap := createTestMap()

	jsonData, err := json.Marshal(testMap)

	if err != nil {
		suite.T().Log("unable to convert test map data into json")
	}

	suite.T().Logf("testing data: %v", jsonData)

	loggerConf.LoadConfigFromJSON(jsonData)

	assert.NotNil(suite.T(), loggerConf.LogLevelMap)
	assert.Equal(suite.T(), testMap, loggerConf.LogLevelMap)
}

func (suite *TestLogger) TestDefaultInstance() {
	logInstance := GetLoggerInstance()

	//Verify the pointers match
	assert.Equal(suite.T(), DefaultLoggerConfig(), logInstance.config)

	// Verify other properties on default instance
	assert.Equal(suite.T(), false, logInstance.timestampEnabled)

	// Default log level names exist
	assert.Equal(suite.T(), logLevelName, logInstance.logLevelName)

}

func (suite *TestLogger) TestInstanceModifications() {
	logInstance := GetLoggerInstance()
	testConfig := createTestConfig()

	//overwrite config
	logInstance.LoadConfig(&testConfig)

	assert.Equal(suite.T(), logInstance.config, &testConfig)

	//new instance - should use singleton
	logInstance2 := GetLoggerInstance()

	//config matches
	assert.Equal(suite.T(), logInstance2.config, &testConfig)

}
