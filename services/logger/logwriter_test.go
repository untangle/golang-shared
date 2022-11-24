package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LogWriterTestSuite struct {
	suite.Suite
}

func TestLogWriterSuitet(t *testing.T) {
	suite.Run(t, new(LogWriterTestSuite))
}

func (suite *LogWriterTestSuite) SetupSuite() {
	loggerSingleton = NewLogger()
}

// Test log levels and log level processor.
func (suite *LogWriterTestSuite) TestLogLevels() {
	// Init log writer
	logWriter := DefaultLogWriter("test")
	logWriter.SetLogLevelProcessor(func(currentLevel LogLevel, message string) LogLevel {
		if strings.Contains(message, "err") {
			return NewLogLevel("ERROR")
		}

		return currentLevel
	})

	log.SetOutput(logWriter)

	for _, testLevel := range []string{"INFO", "ALERT", "CRIT"} {
		// Set Log level
		logWriter.SetLogLevel(NewLogLevel(testLevel))

		result := suite.logAndGetOutput("Test message.")

		// And assert that the message was logged with the right log level.
		assert.Contains(suite.T(), string(result), testLevel)
	}

	result := suite.logAndGetOutput("This is an error.")

	// And assert that the message was logged with the right log level.
	assert.Contains(suite.T(), string(result), "ERROR")
}

func (*LogWriterTestSuite) logAndGetOutput(message string) string {
	// Update stdout so we can catch the output.
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	log.Println(message)

	// Read results.
	_ = w.Close()
	result, _ := io.ReadAll(r)

	// Set stdout to the original value.
	os.Stdout = stdOut

	return string(result)
}
