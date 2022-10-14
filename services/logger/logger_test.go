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

var configSave *MockConfigFile = &MockConfigFile{}

type TestLogger struct {
	suite.Suite
	logger         Logger
	write          LogWriter
	mockConfigFile MockConfigFile
}

// type MockConfig struct {
// 	FileLocation string
// 	LogLevelMap  map[string]string
// 	// OutputWriter MockOutputWriter
// 	config Config
// }

var MockLogLevelEmerg int32 = 0
var MockLogLevelAlert int32 = 1
var MockLogLevelCrit int32 = 2
var MockLogLevelErr int32 = 3
var MockLogLevelWarn int32 = 4
var MockLogLevelNotice int32 = 5
var MockLogLevelInfo int32 = 6
var MockLogLevelDebug int32 = 7
var MockLogLevelTrace int32 = 8

func (m *MockConfigFile) MockLoadConfigFromFile(logger *Logger) {
	logger.logLevelMap = map[string]*int32{
		"Emergtest":  &MockLogLevelEmerg,
		"Alerttest":  &MockLogLevelAlert,
		"Crittest":   &MockLogLevelCrit,
		"Errtest":    &MockLogLevelErr,
		"Warntest":   &MockLogLevelWarn,
		"Noticetest": &MockLogLevelNotice,
		"Infotest":   &MockLogLevelInfo,
		"Debugtest":  &MockLogLevelDebug,
		"Tracetest":  &MockLogLevelTrace,
	}
}

type MockConfig interface {
	LoadConfigFromFile()
	ValidateConfig()
}

func (m *MockConfigFile) LoadConfigFromFile() []byte {
	data := []byte{102, 97, 108, 99, 111, 110}
	return data
}
func (m *MockConfigFile) ValidateConfig() error {
	returns := m.Mock.Called()
	return returns.Error(0)
}
func (suite *TestLogger) TestStartup() {
	testObj := new(MockConfigFile)
	// configSave.On("ValidateConfig").Return(0)
	testObj.On("logger.LoadConfigFromFile", mock.Anything)
	suite.logger.Startup()
}

//Include the case where functionname != 0
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

func (suite *TestLogger) TestGenerateReport() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	Mockbuffer := bytes.Buffer{}
	assert.Equal(suite.T(), 0, Mockbuffer.Len())
	suite.logger.GenerateReport(&Mockbuffer)
	assert.Equal(suite.T(), 715, Mockbuffer.Len())
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

func (suite *TestLogger) TestAdjustSource() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	assert.Equal(suite.T(), int32(-1), suite.logger.AdjustSourceLogLevel("INFO", 7))
	assert.Equal(suite.T(), int32(0), suite.logger.AdjustSourceLogLevel("Emergtest", 2))
	assert.Equal(suite.T(), int32(2), suite.logger.AdjustSourceLogLevel("Crittest", 4))

}

func (suite *TestLogger) TestSearchSourceLogLevel() {
	testObj := new(MockConfigFile)
	testObj.MockLoadConfigFromFile(&suite.logger)
	assert.Equal(suite.T(), int32(-1), suite.logger.SearchSourceLogLevel("INFO"))
	assert.Equal(suite.T(), int32(2), suite.logger.SearchSourceLogLevel("Emergtest"))
	assert.Equal(suite.T(), int32(4), suite.logger.SearchSourceLogLevel("Crittest"))
}

func (suite *TestLogger) TestWrite() {
	int_result, error_result := suite.write.Write([]byte("test\n"))
	assert.Equal(suite.T(), 5, int_result)
	assert.Equal(suite.T(), nil, error_result)
}

func (suite *TestLogger) TestDefaultLogWriter() {
	assert.Equal(suite.T(), &LogWriter{buffer: []uint8{}, source: "System"}, DefaultLogWriter("System"))
}

func (suite *TestLogger) TestEnableTimestamp() {
	suite.logger.timestampEnabled = false
	suite.logger.EnableTimestamp()
	assert.Equal(suite.T(), true, suite.logger.timestampEnabled)
	suite.logger.timestampEnabled = true
	suite.logger.EnableTimestamp()
	assert.Equal(suite.T(), true, suite.logger.timestampEnabled)
}

func (suite *TestLogger) TestDisableTimestamp() {
	suite.logger.timestampEnabled = false
	suite.logger.DisableTimestamp()
	assert.Equal(suite.T(), false, suite.logger.timestampEnabled)
	suite.logger.timestampEnabled = true
	suite.logger.DisableTimestamp()
	assert.Equal(suite.T(), false, suite.logger.timestampEnabled)
}
