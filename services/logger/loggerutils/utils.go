package loggerutils

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

// ConfigureLoggerFromSettings loads the logger configuration from the settings file
// and applies it to the logger singleton. It should be called during application startup.
func ConfigureLoggerFromSettings(
	log *logger.Logger,
	settingsFile *settings.SettingsFile,
	settingsPath ...string) error {

	logLevelMap := make(map[string]logger.LogLevel)

	if err := settingsFile.UnmarshalSettingsAtPath(&logLevelMap, settingsPath...); err != nil {
		return fmt.Errorf("unable to find logger configs in path %s: %w\n", strings.Join(settingsPath, ","), err)
	}

	conf := logger.DefaultLoggerConfig()
	conf.SetLogLevelMap(logLevelMap)

	log.LoadConfig(conf)
	return nil
}
