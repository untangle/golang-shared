package logger

// Logger Model.

// LoggerLevels logger API interface.
//
// It is implemented by the logger service (github.com/untangle/golang-shared/services/logger).
//
// Defines how the logger should behave.
// This model is kept separate from the logger service package,
// because it allows using other components in the logger
// service that also need logging (avoid circular dependencies).
type LoggerLevels interface {
	Emerg(format string, args ...interface{})
	Alert(format string, args ...interface{})
	Crit(format string, args ...interface{})
	Err(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	OCWarn(format string, name string, limit int64, args ...interface{})
	IsDebugEnabled() bool
}
