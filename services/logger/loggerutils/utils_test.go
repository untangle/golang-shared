package loggerutils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/logger/loggerutils"
	"github.com/untangle/golang-shared/services/settings"
)

func TestConfigureLoggerFromSettings(t *testing.T) {
	// Arrange
	content := []byte(`{
		"system": {
			"logging": {
				"foo": { "logname": "DEBUG" },
				"bar": { "logname": "WARN" },
				"*":   { "logname": "INFO" }
			}
		}
	}`)

	tmpfile, err := os.CreateTemp("", "settings-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	sf := settings.NewSettingsFileForTesting(tmpfile.Name())
	log := logger.NewLogger()
	path := []string{"system", "logging"}

	// Act
	err = loggerutils.ConfigureLoggerFromSettings(log, sf, path...)

	// Assert
	require.NoError(t, err)

	config := log.GetConfig()
	assert.Equal(t, logger.LogLevelDebug, config.LogLevelHighest)

	// for "foo", DEBUG is enabled (7), TRACE is not (8)
	assert.True(t, log.IsLogEnabledSource(logger.LogLevelDebug, "foo"))
	assert.False(t, log.IsLogEnabledSource(logger.LogLevelTrace, "foo"))

	// for "bar", WARN is enabled (4), INFO is not (6)
	assert.True(t, log.IsLogEnabledSource(logger.LogLevelWarn, "bar"))
	assert.False(t, log.IsLogEnabledSource(logger.LogLevelInfo, "bar"))

	// for others ("baz"), INFO is enabled (6), DEBUG is not (7)
	assert.True(t, log.IsLogEnabledSource(logger.LogLevelInfo, "baz"))
	assert.False(t, log.IsLogEnabledSource(logger.LogLevelDebug, "baz"))
}

func TestConfigureLoggerFromSettings_Error(t *testing.T) {
	// Arrange: Non-existent file
	sf := settings.NewSettingsFileForTesting("non-existent-file.json")
	log := logger.NewLogger()
	path := []string{"system", "logging"}

	// Act
	err := loggerutils.ConfigureLoggerFromSettings(log, sf, path...)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to find logger configs in path")
	assert.Contains(t, err.Error(), "settings file: unable to open file")
}
