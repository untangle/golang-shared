package logger

/** Remove these once logger is being used properly **/

// For backward compatibility
func Trace(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Trace(format, args)
}

// For backward compatibility
func IsTraceEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelTrace)
}

// For backward compatibility
func Debug(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Debug(format, args)
}

// For backward compatibility
func IsDebugEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelDebug)
}

// For backward compatibility
func Info(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Info(format, args)
}

// For backward compatibility
func IsInfoEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelInfo)
}

// For backward compatibility
func Notice(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Notice(format, args)
}

// For backward compatibility
func IsNoticeEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelNotice)
}

// For backward compatibility
func Warn(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Warn(format, args)
}

// For backward compatibility
func IsWarnEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelWarn)
}

// For backward compatibility
func Err(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Err(format, args)
}

// For backward compatibility
func IsErrEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelErr)
}

// For backward compatibility
func Crit(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Crit(format, args)
}

// For backward compatibility
func IsCritEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelCrit)
}

// For backward compatibility
func Alert(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Alert(format, args)
}

// For backward compatibility
func IsAlertEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelAlert)
}

// For backward compatibility
func Emerg(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.Emerg(format, args)
}

// For backward compatibility
func IsEmergEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelEmerg)
}
