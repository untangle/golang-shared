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
func ConfigureLoggerFromSettings(settingsPath ...string) error {
	settingsFile, err := settings.GetSettingsFileSingleton()
	if err != nil {
		// The logger singleton will exist with default settings, so we can use it.
		logger.GetLoggerInstance().Warn("loggerutils: Could not get settings file singleton: %v", err)
		// This is not a fatal error; the logger will just use defaults.
	}

	logLevelMap := make(map[string]logger.LogLevel)

	if err := settingsFile.UnmarshalSettingsAtPath(&logLevelMap, settingsPath...); err != nil {
		return fmt.Errorf("unable to find logger configs in path %s: %w", strings.Join(settingsPath, ","), err)
	}

	conf := logger.DefaultLoggerConfig()
	conf.SetLogLevelMap(logLevelMap)

	logger.GetLoggerInstance().LoadConfig(conf)
	return nil
}

// StartConfigReloadingOnSIGHUP sets up a listener for the SIGHUP signal to reload the logger configuration.
// This should be called once during application startup.
func StartConfigReloadingOnSIGHUP(settingsPath ...string) {
	go func() {
		hupch := make(chan os.Signal, 1)
		signal.Notify(hupch, syscall.SIGHUP)

		for {
			sig := <-hupch
			log := logger.GetLoggerInstance()
			log.Info("Received signal [%v]. Refreshing logger config\n", sig)
			if err := ConfigureLoggerFromSettings(settingsPath...); err != nil {
				log.Warn("Failed to refresh logger config on SIGHUP: %v\n", err)
			}
		}
	}()
}
