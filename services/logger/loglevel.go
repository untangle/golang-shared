package logger

import (
	"errors"
	"strings"
)

var ErrInvalidLogLevel = errors.New("invalid log level")

// LogLevel struct retains the loglevel information in a string and int.
type LogLevel struct {
	Name string `json:"logname"`
	id   uint8
}

var logLevelMap = map[string]int32{
	strings.ToUpper(logLevelName[LogLevelEmerg]):  LogLevelEmerg,
	strings.ToUpper(logLevelName[LogLevelAlert]):  LogLevelAlert,
	strings.ToUpper(logLevelName[LogLevelCrit]):   LogLevelCrit,
	strings.ToUpper(logLevelName[LogLevelErr]):    LogLevelErr,
	strings.ToUpper(logLevelName[LogLevelWarn]):   LogLevelWarn,
	strings.ToUpper(logLevelName[LogLevelNotice]): LogLevelNotice,
	strings.ToUpper(logLevelName[LogLevelInfo]):   LogLevelInfo,
	strings.ToUpper(logLevelName[LogLevelDebug]):  LogLevelDebug,
	strings.ToUpper(logLevelName[LogLevelTrace]):  LogLevelTrace,
}

// GetId returns the numeric log level for the arugmented name
// or a negative value if the level is not valid
func (lvl *LogLevel) GetId() int32 {
	// Originally, this was iterating through the logLevelName array
	// repeatedly doing strings.ToUpper() so this should be an improvement
	if retval, ok := logLevelMap[strings.ToUpper(lvl.Name)]; ok {
		return retval
	}
	return -1
}

func NewLogLevel(name string) LogLevel {
	loglevel := LogLevel{Name: name}
	loglevel.id = uint8(loglevel.GetId())
	return loglevel
}
