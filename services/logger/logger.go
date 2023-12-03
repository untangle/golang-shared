package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/alerts"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"

	"github.com/untangle/golang-shared/services/overseer"
)

const serviceName = "logger"

// Ocname struct retains information about the counter name and limit
type Ocname struct {
	name  string
	limit int64
}

// Cache for mapping program counters to program/function names
type functionInfoType struct {
	packageName  string
	functionName string
}

var PcFunctionCache = make(map[uintptr]functionInfoType)

// Read/write lock for cache
var mapMutex sync.RWMutex

// Logger struct retains information about the logger related information
type Logger struct {
	config           *LoggerConfig
	defaultConfig    *LoggerConfig
	configLocker     sync.Mutex
	logLevelLocker   sync.RWMutex
	launchTime       time.Time
	timestampEnabled bool
	alerts           alerts.AlertPublisher
	// logCount is added for testing purposes
	logCount uint64
}

// This is only used inernally.
// This was originally stored in every Logger structure which seems wasteful
var logLevelName = [...]string{"EMERG", "ALERT", "CRIT", "ERROR", "WARN", "NOTICE", "INFO", "DEBUG", "TRACE"}

// LogLevelEmerg = syslog.h/LOG_EMERG
const LogLevelEmerg int32 = 0

// LogLevelAlert = syslog.h/LOG_ALERT
const LogLevelAlert int32 = 1

// LogLevelCrit = syslog.h/LOG_CRIT
const LogLevelCrit int32 = 2

// LogLevelErr = syslog.h/LOG_ERR
const LogLevelErr int32 = 3

// LogLevelWarn = syslog.h/LOG_WARNING
const LogLevelWarn int32 = 4

// LogLevelNotice = syslog.h/LOG_NOTICE
const LogLevelNotice int32 = 5

// LogLevelInfo = syslog.h/LOG_INFO
const LogLevelInfo int32 = 6

// LogLevelDebug = syslog.h/LOG_DEBUG
const LogLevelDebug int32 = 7

// LogLevelTrace = custom value
const LogLevelTrace int32 = 8

var loggerSingleton *Logger
var once sync.Once

// GetLoggerInstance returns a logger object that is singleton
// with a wildcard loglevelmap as default.
func GetLoggerInstance() *Logger {
	once.Do(func() {
		loggerSingleton = NewLogger()
	})

	return loggerSingleton
}

// GetLoggerInstanceWithConfig returns a logger object, that's loaded with the config as well
func GetLoggerInstanceWithConfig(conf *LoggerConfig) *Logger {
	instance := GetLoggerInstance()
	instance.LoadConfig(conf)
	return instance
}

// SetLoggerInstance will override the singleton instance with a new instance reference
// This is mostly used for testing concurrency
func SetLoggerInstance(newSingleton *Logger) {
	loggerSingleton = newSingleton
}

// NewLogger creates an new instance of the logger struct with wildcard config
func NewLogger() *Logger {
	return &Logger{
		defaultConfig:    DefaultLoggerConfig(),
		config:           DefaultLoggerConfig(),
		logLevelLocker:   sync.RWMutex{},
		launchTime:       time.Time{},
		timestampEnabled: false,
	}
}

// DefaultLoggerConfig generates a default config with no file location, and INFO log for all log lines
func DefaultLoggerConfig() *LoggerConfig {

	return &LoggerConfig{
		FileLocation: "",
		LogLevelMap:  map[string]LogLevel{"*": {Name: "INFO"}},
		// Default logLevelMask is set to LogLevelInfo
		LogLevelHighest: LogLevelInfo,
		OutputWriter:    DefaultLogWriter("system"),
		CmdAlertSetup:   CmdAlertDefaultSetup,
	}
}

// LoadConfig loads the config to the current logger
// the new config will be set to the defaultConfig
// if we are able to load the config from file, we will
// if the file does not exist, we will store the default config in the conf.FileLocation
func (logger *Logger) LoadConfig(conf *LoggerConfig) {

	logger.defaultConfig = conf
	// load from file - if this is missing or errors - then save the new default config to OS
	// Load config from file if it exists
	err := conf.LoadConfigFromFile()
	if err != nil {
		logger.Warn("No existing config found - using default as current, err: %s\n", err)
		conf.SaveConfig()
	}

	logger.configLocker.Lock()
	defer logger.configLocker.Unlock()
	//Set the instance config to this config
	logger.config = conf
}

// GetConfig returns the logger config
func (logger *Logger) GetConfig() LoggerConfig {
	defer logger.configLocker.Unlock()
	logger.configLocker.Lock()
	return *logger.config
}

// GetConfig returns the logger config
func (logger *Logger) GetDefaultConfig() LoggerConfig {
	return *logger.defaultConfig
}

// Return a count of the number of logs that were actually printed
// This is used for testing purposes.

func (logger *Logger) getLogCount() uint64 {
	return logger.logCount
}

// Startup starts the logging service
func (logger *Logger) Startup() {

	// capture startup time
	logger.launchTime = time.Now()
	logger.alerts = alerts.Publisher(logger)

	if logger.config != nil {

		// Set system logger to use our logger
		if logger.config.OutputWriter != nil {
			log.SetOutput(logger.config.OutputWriter)
		}
	}
}

// Name returns the service name
func (logger *Logger) Name() string {
	return serviceName
}

// Shutdown stops the logging service
func (logger *Logger) Shutdown() {
	alerts.Shutdown()
	fmt.Println("Shutting down the logger service")
}

// Emerg is called for log level EMERG messages
func (logger *Logger) Emerg(format string, args ...interface{}) {
	logger.logMessage(LogLevelEmerg, format, Ocname{}, args...)
}

// IsEmergEnabled returns true if EMERG logging is enable for the caller
func (logger *Logger) IsEmergEnabled() bool {
	if LogLevelEmerg > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelEmerg)
}

// Alert is called for log level ALERT messages
func (logger *Logger) Alert(format string, args ...interface{}) {
	logger.logMessage(LogLevelAlert, format, Ocname{}, args...)
}

// IsAlertEnabled returns true if ALERT logging is enable for the caller
func (logger *Logger) IsAlertEnabled() bool {
	if LogLevelAlert > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelAlert)
}

// Crit is called for log level CRIT messages
func (logger *Logger) Crit(format string, args ...interface{}) {
	logger.logMessage(LogLevelCrit, format, Ocname{}, args...)
}

// IsCritEnabled returns true if CRIT logging is enable for the caller
func (logger *Logger) IsCritEnabled() bool {
	if LogLevelCrit > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelCrit)
}

// Err is called for log level ERR messages
func (logger *Logger) Err(format string, args ...interface{}) {
	logger.logMessage(LogLevelErr, format, Ocname{}, args...)
}

// IsErrEnabled returns true if ERR logging is enable for the caller
func (logger *Logger) IsErrEnabled() bool {
	if LogLevelErr > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelErr)
}

// Warn is called for log level WARNING messages
func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.logMessage(LogLevelWarn, format, Ocname{}, args...)
}

// IsWarnEnabled returns true if WARNING logging is enable for the caller
func (logger *Logger) IsWarnEnabled() bool {
	if LogLevelWarn > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelWarn)
}

// Notice is called for log level NOTICE messages
func (logger *Logger) Notice(format string, args ...interface{}) {
	logger.logMessage(LogLevelNotice, format, Ocname{}, args...)
}

// IsNoticeEnabled returns true if NOTICE logging is enable for the caller
func (logger *Logger) IsNoticeEnabled() bool {
	if LogLevelNotice > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelNotice)
}

// Info is called for log level INFO messages
func (logger *Logger) Info(format string, args ...interface{}) {
	logger.logMessage(LogLevelInfo, format, Ocname{}, args...)
}

// IsInfoEnabled returns true if INFO logging is enable for the caller
func (logger *Logger) IsInfoEnabled() bool {
	if LogLevelInfo > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelInfo)
}

// Debug is called for log level DEBUG messages
func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.logMessage(LogLevelDebug, format, Ocname{}, args...)
}

// IsDebugEnabled returns true if DEBUG logging is enable for the caller
func (logger *Logger) IsDebugEnabled() bool {
	if LogLevelDebug > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelDebug)
}

// Trace is called for log level TRACE messages
func (logger *Logger) Trace(format string, args ...interface{}) {
	logger.logMessage(LogLevelTrace, format, Ocname{}, args...)
}

// OCTrace is called for overseer messages
func (logger *Logger) OCTrace(format string, name string, limit int64, args ...interface{}) {
	newOcname := Ocname{name, limit}
	logger.logMessage(LogLevelTrace, format, newOcname, args...)
}

// OCWarn is called for overseer warn messages
func (logger *Logger) OCWarn(format string, name string, limit int64, args ...interface{}) {
	newOcname := Ocname{name, limit}
	logger.logMessage(LogLevelTrace, format, newOcname, args...)
}

// OCDebug is called for overseer warn messages
func (logger *Logger) OCDebug(format string, name string, limit int64, args ...interface{}) {
	newOcname := Ocname{name, limit}
	logger.logMessage(LogLevelTrace, format, newOcname, args...)
}

// OCErr is called for overseer err messages
func (logger *Logger) OCErr(format string, name string, limit int64, args ...interface{}) {
	newOcname := Ocname{name, limit}
	logger.logMessage(LogLevelTrace, format, newOcname, args...)
}

// OCCrit is called for overseer crit messages
func (logger *Logger) OCCrit(format string, name string, limit int64, args ...interface{}) {
	newOcname := Ocname{name, limit}
	logger.logMessage(LogLevelTrace, format, newOcname, args...)
}

// IsTraceEnabled returns true if TRACE logging is enable for the caller
func (logger *Logger) IsTraceEnabled() bool {
	if LogLevelTrace > logger.config.LogLevelHighest {
		return false
	}
	return logger.isLogEnabled(LogLevelTrace)
}

// LogMessageSource is for the netfilter interface functions written in C
// and our LogWriter type that can be created and passed to anything that
// expects an object with output stream support. The logging source is passed
// directly rather than determined from the call stack.
func LogMessageSource(level int32, source string, format string, args ...interface{}) {
	logger := GetLoggerInstance()
	logger.logMessage(level, format, Ocname{}, args...)
}

// IsLogEnabledSource returns true if logging is enabled at the argumented level for the argumented source
func (logger *Logger) IsLogEnabledSource(level int32, source string) bool {
	lvl := logger.getLogLevel(source, "")
	return (lvl >= level)
}

// DisableTimestamp disable the elapsed time in output
func (logger *Logger) DisableTimestamp() {
	logger.timestampEnabled = false
}

// getLogLevel returns the log level for the specified package or function
// It checks function first allowing individual functions to be configured
// for a higher level of logging than the package that owns them.
func (logger *Logger) getLogLevel(packageName string, functionName string) int32 {

	if len(functionName) != 0 {
		logger.logLevelLocker.RLock()
		level, ok := logger.config.LogLevelMap[functionName]
		logger.logLevelLocker.RUnlock()
		if ok {
			return int32(level.GetId())
		}
	}

	if len(packageName) != 0 {
		logger.logLevelLocker.RLock()
		level, ok := logger.config.LogLevelMap[packageName]
		logger.logLevelLocker.RUnlock()
		if ok {
			return int32(level.GetId())
		} else {
			if val, ok := logger.config.LogLevelMap["*"]; ok {
				return int32(val.GetId())
			}
		}
	}
	// nothing found so return default level
	return LogLevelInfo
}

// logFormatter creats a log message using the format and arguments provided
// We look for and handle special format verbs that trigger additional processing
func logFormatter(format string, newOcname Ocname, args ...interface{}) string {

	total := overseer.AddCounter(newOcname.name, 1)

	// only format the message on the first and every nnn messages thereafter
	// or if limit is zero which means no limit on logging
	if total == 1 || newOcname.limit == 0 || (total%newOcname.limit) == 0 {
		// if there are only two arguments everything after the verb is the message

		// more than two arguments so use the remaining format and arguments
		buffer := fmt.Sprint(format)
		return buffer
	}
	// return empty string when a repeat is limited
	return ""
}

// isLogEnabled returns true if logging is enabled for the caller at the specified level, false otherwise
func (logger *Logger) isLogEnabled(level int32) bool {
	packageName, functionName := findCallingFunction()
	if logger.IsLogEnabledSource(level, packageName) {
		return true
	}
	if logger.IsLogEnabledSource(level, functionName) {
		return true
	}
	return false
}

// logMessage is called to write messages to the system log
func (logger *Logger) logMessage(level int32, format string, newOcname Ocname, args ...interface{}) {
	// logger.config.LogLevelMask keeps track of the logger levels that have been
	// requested across the entire logger confguration so that we can drop out of this
	// function quickly if the log is for something unlikely like a trace or debug.
	if level > logger.config.LogLevelHighest {
		return
	}
	packageName, functionName := findCallingFunction()

	testLevel := logger.getLogLevel(packageName, functionName)

	if level > testLevel {
		return
	}

	logger.logLevelLocker.RLock()

	var logMessage string

	// If the Ocname is an empty struct, then we are not running %OC logic
	if (newOcname == Ocname{}) {
		logMessage = fmt.Sprintf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], packageName, fmt.Sprintf(format, args...))
	} else { //Handle %OC - buffer the logs on this logger instance until we hit the limit
		buffer := logFormatter(format, newOcname, args...)
		if len(buffer) == 0 {
			logger.logLevelLocker.RUnlock()
			return
		}
		logMessage = fmt.Sprintf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], packageName, buffer)
	}
	fmt.Print(logMessage)

	logger.configLocker.Lock()

	// This is protected by the configLogger.Lock() to avoid concurrency problems
	logger.logCount++

	if alert, ok := logger.config.CmdAlertSetup[level]; ok && logger.alerts != nil {
		logger.configLocker.Unlock()
		logger.logLevelLocker.RUnlock()
		logger.alerts.Send(&Alerts.Alert{
			Type:          alert.logType,
			Severity:      alert.severity,
			Message:       logMessage,
			IsLoggerAlert: true,
		})
		return
	}
	logger.configLocker.Unlock()
	logger.logLevelLocker.RUnlock()
}

// func findCallingFunction() uses runtime.Callers to get the call stack to determine the calling function
// Our public function heirarchy is implemented so the caller is always at the 4th frame
// Frame 0 = runtime.Callers
// Frame 1 = findCallingFunction
// Frame 2 = logMessage / isLogEnabled
// Frame 3 = Warn, Info / IsWarnEnabled, IsInfoEnabled (etc...)
// Frame 4 = the function that actually called logger.Warn, logger.Info, logger.IsWarnEnabled, logger.IsInfoEnabled, etc...
// Here is an example of what we expect to see in the calling function frame:
// FILE: /home/username/golang/src/github.com/untangle/packetd/services/dict/dict.go
// FUNC: github.com/untangle/packetd/services/dict.cleanDictionary
// LINE: 827
//
// packageName		Name like "dict"
// functionName		Package path from package name to calling function.
//
//	This is meant to be an explict path so you can match very granular on a specific function.
//	This can be:
//	dict.cleanDictionary
//	plugins.(*PluginControl).Startup
//	plugincommon.(*BctidConsumerCommon[...]).registerOrDeregister
//	dispatch.NetloggerCallback.func1
func findCallingFunction() (packageName string, functionName string) {
	// create a single entry array to hold the 5th stack frame and pass 4 as the
	// number of frames to skip over so we get the single program_counters frame we need
	pc := make([]uintptr, 1)
	count := runtime.Callers(4, pc)
	if count != 1 {
		return "unknown", "unknown"
	}

	// See if program counter is in our cache
	mapMutex.RLock()
	functionInfo, found := PcFunctionCache[pc[0]]
	mapMutex.RUnlock()
	if found {
		return functionInfo.packageName, functionInfo.functionName
	}

	// get the frame object for the caller
	frames := runtime.CallersFrames(pc)
	caller, _ := frames.Next()

	// Find the index of the last slash to isolate the package.FunctionName
	end := strings.LastIndex(caller.Function, "/")
	if end < 0 {
		functionName = caller.Function
	} else {
		functionName = caller.Function[end+1:]
	}

	// Find the index of the dot after the package name
	dot := strings.Index(functionName, ".")
	if dot < 0 {
		packageName = "unknown"
	} else {
		packageName = functionName[0:dot]
	}

	// Add to cache
	mapMutex.Lock()
	PcFunctionCache[pc[0]] = functionInfoType{packageName: packageName, functionName: functionName}
	mapMutex.Unlock()

	return packageName, functionName
}

// getPrefix returns a log message prefix
func (logger *Logger) getPrefix() string {
	if !logger.timestampEnabled {
		return ""
	}

	nowtime := time.Now()
	var elapsed = nowtime.Sub(logger.launchTime)
	return fmt.Sprintf("[%11.5f] ", elapsed.Seconds())
}

// FindLogLevelName returns the log level name for the argumented value
func FindLogLevelName(level int32) string {
	if level < 0 {
		return "UNDEFINED"
	}
	if int(level) > len(logLevelName) {
		return fmt.Sprintf("%d", level)
	}
	return logLevelName[level]
}
