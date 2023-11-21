package logger

/** Remove these once logger is being used properly **/

// For backward compatibility
func Trace(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelTrace, format, Ocname{}, args...)
}

// For backward compatibility
func IsTraceEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelTrace)
}

// For backward compatibility
func Debug(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelDebug, format, Ocname{}, args...)
}

// For backward compatibility
func IsDebugEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelDebug)
}

// For backward compatibility
func Info(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelInfo, format, Ocname{}, args...)
}

// For backward compatibility
func IsInfoEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelInfo)
}

// For backward compatibility
func Notice(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelNotice, format, Ocname{}, args...)
}

// For backward compatibility
func IsNoticeEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelNotice)
}

// For backward compatibility
func Warn(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelWarn, format, Ocname{}, args...)
}

// For backward compatibility
func IsWarnEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelWarn)
}

// For backward compatibility
func Err(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelErr, format, Ocname{}, args...)
}

// For backward compatibility
func IsErrEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelErr)
}

// For backward compatibility
func Crit(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelCrit, format, Ocname{}, args...)
}

// For backward compatibility
func IsCritEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelCrit)
}

// For backward compatibility
func Alert(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelAlert, format, Ocname{}, args...)
}

// For backward compatibility
func IsAlertEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelAlert)
}

// For backward compatibility
func Emerg(format string, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	currentLogger.logMessage(LogLevelEmerg, format, Ocname{}, args...)
}

// For backward compatibility
func IsEmergEnabled() bool {
	currentLogger := GetLoggerInstance()
	return currentLogger.isLogEnabled(LogLevelEmerg)
}

// For backward compatibility
func OCCrit(format string, name string, limit int64, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	newOcname := Ocname{name, limit}
	currentLogger.logMessage(LogLevelCrit, format, newOcname, args...)
}

// For backward compatibility
func OCErr(format string, name string, limit int64, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	newOcname := Ocname{name, limit}
	currentLogger.logMessage(LogLevelErr, format, newOcname, args...)
}

// For backward compatibility
func OCDebug(format string, name string, limit int64, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	newOcname := Ocname{name, limit}
	currentLogger.logMessage(LogLevelDebug, format, newOcname, args...)
}

// For backward compatibility
func OCWarn(format string, name string, limit int64, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	newOcname := Ocname{name, limit}
	currentLogger.logMessage(LogLevelWarn, format, newOcname, args...)
}

// For backward compatibility
func OCTrace(format string, name string, limit int64, args ...interface{}) {
	currentLogger := GetLoggerInstance()
	newOcname := Ocname{name, limit}
	currentLogger.logMessage(LogLevelTrace, format, newOcname, args...)
}
