package logger

import (
	"io"

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
	LogLevelHighest int32
	OutputWriter    io.Writer
	CmdAlertSetup   map[int32]CmdAlertDetail

	logLevelMap map[string]LogLevel
}


// SetLogLevelMap sets the log level map and updates the highest level
func (conf *LoggerConfig) SetLogLevelMap(logLevelMap map[string]LogLevel) {
	conf.logLevelMap = logLevelMap
	conf.SetLogLevelHighest()
}

// SetLogLevelHighest will set the highest log level in the log config
func (conf *LoggerConfig) SetLogLevelHighest() {
	for _, v := range conf.logLevelMap {
		if v.GetId() > conf.LogLevelHighest {
			conf.LogLevelHighest = v.GetId()
		}
	}
}

// SetLogLevel can set the log level in the log config
func (conf *LoggerConfig) SetLogLevel(key string, newLevel LogLevel) {
	conf.logLevelMap[key] = newLevel
	if newLevel.GetId() > conf.LogLevelHighest {
		conf.LogLevelHighest = newLevel.GetId()
	}
}
