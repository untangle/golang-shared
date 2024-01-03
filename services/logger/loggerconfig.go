package logger

import (
	"fmt"
	"io"

	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
)

type CmdAlertDetail struct {
	severity Alerts.AlertSeverity
	logType  Alerts.AlertType
}

var CmdAlertDefaultSetup = map[int32]CmdAlertDetail{
	LogLevelCrit: {
		severity: Alerts.AlertSeverity_CRITICAL,
		logType:  Alerts.AlertType_CRITICALERROR,
	},
	LogLevelErr: {
		severity: Alerts.AlertSeverity_ERROR,
		logType:  Alerts.AlertType_CRITICALERROR,
	},
}

// LoggerConfig struct retains information about the where the log level map is stored, default log levels and writer that should be used
type LoggerConfig struct {
	SettingsFile    *settings.SettingsFile
	SettingsPath    []string
	LogLevelHighest int32
	LogLevelMap     map[string]LogLevel
	OutputWriter    io.Writer
	CmdAlertSetup   map[int32]CmdAlertDetail
}

// loadLoggerConfig loads the logger configuration file
func (conf *LoggerConfig) LoadConfigFromFile() error {
	if len(conf.SettingsPath) == 0 {
		return fmt.Errorf("Logger config settings path is missing")
	}

	logLevelMap := make(map[string]LogLevel)

	if err := conf.SettingsFile.UnmarshalSettingsAtPath(&logLevelMap, conf.SettingsPath...); err != nil {
		return fmt.Errorf("Unable to find logger configs: %s\n", err)
	}

	conf.LogLevelMap = logLevelMap

	// set the highest log level
	conf.SetLogLevelHighest()
	return nil
}

// SetLogLevelHighest will set the highest log level in the log config
func (conf *LoggerConfig) SetLogLevelHighest() {
	for _, v := range conf.LogLevelMap {
		if v.GetId() > conf.LogLevelHighest {
			conf.LogLevelHighest = v.GetId()
		}
	}
}

// SetLogLevel can set the log level in the log config
func (conf *LoggerConfig) SetLogLevel(key string, newLevel LogLevel) {
	conf.LogLevelMap[key] = newLevel
	if newLevel.GetId() > conf.LogLevelHighest {
		conf.LogLevelHighest = newLevel.GetId()
	}
}
