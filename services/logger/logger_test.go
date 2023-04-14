package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"github.com/untangle/golang-shared/testing/mocks"

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
	assert.Equal(suite.T(), nil, logger.config.OutputWriter)

	logger.config.FileLocation = "LoggerConfig.json"
	logger.Startup()

	var MockWriter bytes.Buffer
	logger.config.OutputWriter = &MockWriter
	logger.Startup()
}

// Test default service name
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
	assert.Error(suite.T(), fmt.Errorf("Logger config FileLocation is missing"), logger.config.LoadConfigFromFile())

	// Test that load config from file works
	logger.config.FileLocation = "LoggerConfig.json"
	err := logger.config.LoadConfigFromFile()
	assert.NoError(suite.T(), err)

	//Test that the LoggerConfig.json matches some properties
	assert.Equal(suite.T(), LogLevelInfo, logger.getLogLevel("discovery", "discovery"))

	//conntrack is debug
	assert.Equal(suite.T(), LogLevelTrace, logger.getLogLevel("conntrack", "conntrack"))

	//classify is trace
	assert.Equal(suite.T(), LogLevelDebug, logger.getLogLevel("classify", "classify"))

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
	logInstance := NewLogger()

	//Verify the pointers match
	assert.Equal(suite.T(), DefaultLoggerConfig(), logInstance.config)

	// Verify other properties on default instance
	assert.Equal(suite.T(), false, logInstance.timestampEnabled)

	// Default log level names exist
	assert.Equal(suite.T(), logLevelName, logInstance.logLevelName)

}

func (suite *TestLogger) TestInstanceModifications() {

	SetLoggerInstance(NewLogger())

	logInstance := GetLoggerInstance()
	testConfig := createTestConfig()
	testConfig.removeConfigFile()

	//overwrite config
	logInstance.LoadConfig(&testConfig)

	assert.Equal(suite.T(), testConfig, logInstance.GetConfig())

	//new instance - should use singleton
	logInstance2 := GetLoggerInstance()

	// Assert singletons have the same pointer value
	assert.Equal(suite.T(), logInstance, logInstance2)

	//config matches
	assert.Equal(suite.T(), testConfig, logInstance2.GetConfig())

}

// func (suite *TestLogger) TestMultiThreadAccess() {
// 	currentCtx := context.Background()
// 	SetLoggerInstance(NewLogger())
// 	logInstance := GetLoggerInstance()
// 	testingOutput := "Testing output for %s\n"
// 	expectedConfig := createTestConfig()
// 	expectedConfig.removeConfigFile()

// 	go func(testingOutput string, expectedConfig LoggerConfig, ctx context.Context) {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:

// 				logInstance := GetLoggerInstance()
// 				logInstance.Debug(testingOutput, logLevelName[LogLevelDebug])
// 				logInstance.Info(testingOutput, logLevelName[LogLevelInfo])

// 				time.Sleep(time.Millisecond * 2)

// 				// config pointer matches after waiting
// 				assert.Equal(suite.T(), expectedConfig, logInstance.GetConfig())
// 			}
// 		}
// 	}(testingOutput, expectedConfig, currentCtx)

// 	time.Sleep(time.Millisecond * 1)
// 	//Change config after routine starts to enable DEBUG
// 	expectedConfig.SetLogLevel("runtime", NewLogLevel("DEBUG"))
// 	expectedConfig.SetLogLevel("reflect", NewLogLevel("DEBUG"))

// 	// Load new config to the instance
// 	logInstance.LoadConfig(&expectedConfig)
// 	logInstance.Debug(testingOutput, logLevelName[LogLevelDebug])
// 	logInstance.Info(testingOutput, logLevelName[LogLevelInfo])

// 	time.Sleep(time.Millisecond * 5)
// 	currentCtx.Done()
// }

func (suite *TestLogger) TestInstanceLoadFromDisk() {
	logInstance := NewLogger()
	testConfig := createTestConfig()
	testConfig.removeConfigFile()

	//overwrite default config
	logInstance.LoadConfig(&testConfig)

	// now load from file
	logInstance.config.FileLocation = "LoggerConfig.json"
	logInstance.config.LoadConfigFromFile()

	// verify these are different
	assert.NotEqual(suite.T(), testConfig, logInstance.config)

	// Verify we loaded the actual config
	//assert.Equal(suite.T(), NewLogLevel("WARN"), logInstance.config.GetLogLevel("test"))
}

func (suite *TestLogger) TestSaveToDisk() {
	logInstance := NewLogger()

	//Verify we loaded the default config
	assert.Equal(suite.T(), DefaultLoggerConfig(), logInstance.config)

	// Create the test config - save it, load it to the new instance and verify it loaded
	testConfig := createTestConfig()
	testConfig.removeConfigFile()
	logInstance.LoadConfig(&testConfig)

	assert.Equal(suite.T(), &testConfig, logInstance.config)

}

func (suite *TestLogger) TestBasicWriters() {
	logInstance := NewLogger()

	testingOutput := "Testing output for %s\n"

	assert.Equal(suite.T(), DefaultLogWriter("system"), logInstance.config.OutputWriter)

	//Change log writer to print to a buffer for us to analyze
	logInstance.Info(testingOutput, logLevelName[LogLevelInfo])
	logInstance.Err(testingOutput, logLevelName[LogLevelErr])
	logInstance.Debug(testingOutput, logLevelName[LogLevelDebug])
	logInstance.Notice(testingOutput, logLevelName[LogLevelNotice])
	logInstance.Warn(testingOutput, logLevelName[LogLevelWarn])

	//Bump reflect config up
	logInstance.config.SetLogLevel("reflect", NewLogLevel("DEBUG"))
	//Change log writer to print to a buffer for us to analyze
	logInstance.Info(testingOutput, logLevelName[LogLevelInfo])
	logInstance.Err(testingOutput, logLevelName[LogLevelErr])
	logInstance.Debug(testingOutput, logLevelName[LogLevelDebug])
	logInstance.Notice(testingOutput, logLevelName[LogLevelNotice])
	logInstance.Warn(testingOutput, logLevelName[LogLevelWarn])
}

func (suite *TestLogger) TestSendAlertToCMD() {
	logInstance := NewLogger()

	mockAlertsPublisher := &mocks.MockAlertPublisher{}
	logInstance.alerts = mockAlertsPublisher

	testingOutput := "Testing output for %s\n"

	assert.Equal(suite.T(), DefaultLogWriter("system"), logInstance.config.OutputWriter)

	// Test alert is sent for critical level logs
	logInstance.Crit(testingOutput, logLevelName[LogLevelCrit])

	assert.Equal(suite.T(), Alerts.AlertType_CRITICALERROR, mockAlertsPublisher.LastAlert.Type)
	assert.Equal(suite.T(), Alerts.AlertSeverity_CRITICAL, mockAlertsPublisher.LastAlert.Severity)
	assert.Equal(suite.T(), "CRIT              reflect: Testing output for CRIT\n", mockAlertsPublisher.LastAlert.Message)

	logInstance.configLocker.Lock()
	// Set alerts for error level logs
	logInstance.config.CmdAlertSetup[LogLevelErr] = CmdAlertDetail{
		severity: Alerts.AlertSeverity_ERROR,
		logType:  Alerts.AlertType_CRITICALERROR,
	}
	logInstance.configLocker.Unlock()

	logInstance.Err(testingOutput, logLevelName[LogLevelErr])

	assert.Equal(suite.T(), Alerts.AlertType_CRITICALERROR, mockAlertsPublisher.LastAlert.Type)
	assert.Equal(suite.T(), Alerts.AlertSeverity_ERROR, mockAlertsPublisher.LastAlert.Severity)
	assert.Equal(suite.T(), "ERROR             reflect: Testing output for ERROR\n", mockAlertsPublisher.LastAlert.Message)
}

func (suite *TestLogger) TestFindCallingFunction() {

	fileName, lineNumber, packageName, functionName := findCallingFunction()

	assert.Contains(suite.T(), fileName, "suite.go")
	assert.Greater(suite.T(), lineNumber, 0)
	assert.Equal(suite.T(), "suite", packageName)
	assert.Contains(suite.T(), functionName, "suite.Run")
}

func (suite *TestLogger) TestGetInstanceWithConfig() {
	SetLoggerInstance(NewLogger())
	logInstance := GetLoggerInstance()

	// Verify default config was loaded
	assert.Equal(suite.T(), DefaultLoggerConfig(), logInstance.config)

	expectedConfig := createTestConfig()
	expectedConfig.removeConfigFile()

	newInstance := GetLoggerInstanceWithConfig(&expectedConfig)

	//Verify new instance has proper config
	assert.Equal(suite.T(), &expectedConfig, newInstance.config)

	// Verify old instance has this config too
	assert.Equal(suite.T(), &expectedConfig, logInstance.config)

}

func (suite *TestLogger) TestGetDefaultConfig() {
	logger := Logger{defaultConfig: &LoggerConfig{FileLocation: "LoggerConfig.json"}}

	config := logger.GetDefaultConfig()

	if config.FileLocation != "test.json" {
		assert.Error(suite.T(), fmt.Errorf("Expected FileLocation to be 'test.json', but got %s", config.FileLocation))
	}

	if config.OutputWriter == nil {
		assert.Error(suite.T(), fmt.Errorf("Expected OutputWriter to be non-nil, but got nil"))
	}

	if config.LogLevelMap == nil {
		assert.Error(suite.T(), fmt.Errorf("Expected LogLevelMap to be non-nil, but got nil"))
	}
}

func (suite *TestLogger) TestLoggerEmerg() {
	var buf bytes.Buffer
	logger := NewLogger()

	logger.Emerg("Test message")

	output := buf.String()
	expectedOutput := fmt.Sprintf("[EMERG] Test message\n")
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerAlert() {
	var buf bytes.Buffer
	logger := NewLogger()

	logger.Alert("Test message")

	output := buf.String()
	expectedOutput := "[ALERT] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerCrit() {
	var buf bytes.Buffer
	logger := NewLogger()

	logger.Crit("Test message")

	output := buf.String()
	expectedOutput := "[CRIT] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerTrace() {
	var buf bytes.Buffer
	logger := NewLogger()

	logger.Trace("Test message")

	output := buf.String()
	expectedOutput := "[TRACE] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestOCTrace() {
	var buf bytes.Buffer
	logger := NewLogger()
	name := "test"
	limit := int64(10)

	logger.OCTrace("Test message", name, limit)

	output := buf.String()
	expectedOutput := fmt.Sprintf("[TRACE] %s { name: %s, limit: %d } Test message\n", time.Now().Format("2006-01-02 15:04:05.000000"), name, limit)
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLogFormatter() {
	ocname := Ocname{"test_ocname", 3}

	// First log message should not include the prefix
	output := logFormatter("Test message", ocname)
	expectedOutput := "Test message"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}

	// Second log message should include the prefix
	output = logFormatter("Test message", ocname)
	expectedOutput = "[test_ocname-2] Test message"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}

	// Third log message should not include the prefix because of the limit
	output = logFormatter("Test message", ocname)
	expectedOutput = ""
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}

	// Fourth log message should include the prefix because the limit is zero
	ocname = Ocname{"test_ocname", 0}
	output = logFormatter("Test message", ocname)
	expectedOutput = "[test_ocname-1] Test message"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerOCWarn() {
	var buf bytes.Buffer
	logger := NewLogger()

	// Set up a new LoggerConfig with an OutputWriter that writes to the buffer
	conf := DefaultLoggerConfig()
	conf.OutputWriter = &buf
	logger.config = conf

	logger.OCWarn("Test message", "test_ocname", 100)

	output := buf.String()
	expectedOutput := "[WARN][test_ocname-100] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerOCDebug() {
	var buf bytes.Buffer
	logger := NewLogger()

	// Set up a new LoggerConfig with an OutputWriter that writes to the buffer
	conf := DefaultLoggerConfig()
	conf.OutputWriter = &buf
	logger.config = conf

	logger.OCDebug("Test message", "test_ocname", 100)

	output := buf.String()
	expectedOutput := "[Debug][test_ocname-100] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerOCErr() {
	var buf bytes.Buffer
	logger := NewLogger()

	// Set up a new LoggerConfig with an OutputWriter that writes to the buffer
	conf := DefaultLoggerConfig()
	conf.OutputWriter = &buf
	logger.config = conf

	logger.OCErr("Test message", "test_ocname", 100)

	output := buf.String()
	expectedOutput := "[Err][test_ocname-100] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerOCCrit() {
	var buf bytes.Buffer
	logger := NewLogger()

	// Set up a new LoggerConfig with an OutputWriter that writes to the buffer
	conf := DefaultLoggerConfig()
	conf.OutputWriter = &buf
	logger.config = conf

	logger.OCCrit("Test message", "test_ocname", 100)

	output := buf.String()
	expectedOutput := "[Crit][test_ocname-100] Test message\n"
	if output != expectedOutput {
		assert.Error(suite.T(), fmt.Errorf("Unexpected output, expected: %s, actual: %s", expectedOutput, output))
	}
}

func (suite *TestLogger) TestLoggerDisableTimestamp() {
	var buf bytes.Buffer
	logger := NewLogger()
	conf := DefaultLoggerConfig()
	conf.OutputWriter = &buf
	logger.config = conf
	logger.EnableTimestamp()
	logger.Info("Test message")
	outputWithTimestamp := buf.String()

	buf.Reset()
	logger.DisableTimestamp()
	logger.Info("Test message")
	outputWithoutTimestamp := buf.String()

	if outputWithTimestamp == outputWithoutTimestamp {
		assert.Error(suite.T(), fmt.Errorf("Timestamp was not disabled"))
	}
}
