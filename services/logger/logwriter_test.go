package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LogWriterTestSuite struct {
	suite.Suite
	logInstance *log.Logger
}

func TestLogWriterSuite(t *testing.T) {
	suite.Run(t, new(LogWriterTestSuite))
}

func (suite *LogWriterTestSuite) SetupSuite() {
	suite.logInstance = log.New(DefaultLogWriter("test"), "test", 1)
}

// Test log levels and log level processor.
func (suite *LogWriterTestSuite) TestLogLevels() {
	// Init log writer
	logWriter := DefaultLogWriter("test")

	suite.logInstance.SetOutput(logWriter)

	for _, testLevel := range []string{"INFO", "ALERT", "CRIT"} {
		// Set Log level
		assert.NoError(suite.T(), logWriter.SetLogLevel(NewLogLevel(testLevel)))

		result := suite.logAndGetOutput(logWriter, fmt.Sprintf("Test message level: %s\n", testLevel))

		// And assert that the message was logged with the right log level.
		assert.Contains(suite.T(), string(result), testLevel)
	}
}

func (suite *LogWriterTestSuite) logAndGetOutput(writer *LogWriter, message string) string {

	// Use multiwriter to copy writes to both local log buffer instance and output writer
	var buf bytes.Buffer

	// Set the suite loginstance output to this buffer
	suite.logInstance.SetOutput(&buf)

	// Assign both of these IO.Writer interfaces to a multiwriter
	w := io.MultiWriter(writer, suite.logInstance.Writer())

	// Call write on the multiwriter interface
	count, err := w.Write([]byte(message))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, len(message), "Count of written characters(%d) didn't match message size(%d)", count, len(message))

	// Return the output from the bytesBuffer
	return buf.String()
}
