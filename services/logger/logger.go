package logger

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/untangle/golang-shared/services/overseer"
)

const serviceName = "logger"

// Ocname struct retains information about the counter name and limit
type Ocname struct {
	name  string
	limit int64
}

// Add description [Nikki]
type Logger struct {
	config           LoggerConfig
	logLevelMap      map[string]*int32
	logLevelLocker   sync.RWMutex
	launchTime       time.Time
	timestampEnabled bool //[Nikki] default value = true
	logLevelName     [9]string
}

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

// GetLoggerInstance returns a logger object that is a
// singleton. Prefer using this if you can.
func GetLoggerInstance() *Logger {
	if loggerSingleton == nil {
		loggerSingleton = NewLogger()
	}
	return loggerSingleton
}

// NewLogger creates an new instance of the logger struct
func NewLogger() *Logger {

	return &Logger{}
}

// Startup starts the logging service
func (logger *Logger) Startup() {
	err := logger.config.Validate()
	fmt.Print("Step 1\n")
	fmt.Print(err)
	if err != nil {
		fmt.Print("Step 3\n")
		logger.Err("Logger Configuration is invalid: %s\n", err.Error())
		return
	}
	fmt.Print("Step 2\n")

	// capture startup time
	logger.launchTime = time.Now()

	// create the map and load the Log configuration
	data := logger.config.LoadConfigFromFile()
	if data != nil {
		logger.config.LoadConfigFromJSON(data)
	}

	// Set system logger to use our logger
	if logger.config.OutputWriter != nil {
		log.SetOutput(logger.config.OutputWriter)
	} else {
		log.SetOutput(DefaultLogWriter("system"))
	}
}

// Name returns the service name
func (logger *Logger) Name() string {
	return serviceName
}

// Shutdown stops the logging service
func (logger *Logger) Shutdown() {
	fmt.Println("Shutting down the logger service")
}

// Emerg is called for log level EMERG messages
func (logger *Logger) Emerg(format string, args ...interface{}) {
	logger.logMessage(LogLevelEmerg, format, Ocname{"", 0}, args...)
}

// IsEmergEnabled returns true if EMERG logging is enable for the caller
func (logger *Logger) IsEmergEnabled() bool {
	return logger.isLogEnabled(LogLevelEmerg)
}

// Alert is called for log level ALERT messages
func (logger *Logger) Alert(format string, args ...interface{}) {
	logger.logMessage(LogLevelAlert, format, Ocname{"", 0}, args...)
}

// IsAlertEnabled returns true if ALERT logging is enable for the caller
func (logger *Logger) IsAlertEnabled() bool {
	return logger.isLogEnabled(LogLevelAlert)
}

// Crit is called for log level CRIT messages
func (logger *Logger) Crit(format string, args ...interface{}) {
	logger.logMessage(LogLevelCrit, format, Ocname{"", 0}, args...)
}

// IsCritEnabled returns true if CRIT logging is enable for the caller
func (logger *Logger) IsCritEnabled() bool {
	return logger.isLogEnabled(LogLevelCrit)
}

// Err is called for log level ERR messages
func (logger *Logger) Err(format string, args ...interface{}) {
	logger.logMessage(LogLevelErr, format, Ocname{"", 0}, args...)
}

// IsErrEnabled returns true if ERR logging is enable for the caller
func (logger *Logger) IsErrEnabled() bool {
	return logger.isLogEnabled(LogLevelErr)
}

// Warn is called for log level WARNING messages
func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.logMessage(LogLevelWarn, format, Ocname{"", 0}, args...)
}

// IsWarnEnabled returns true if WARNING logging is enable for the caller
func (logger *Logger) IsWarnEnabled() bool {
	return logger.isLogEnabled(LogLevelWarn)
}

// Notice is called for log level NOTICE messages
func (logger *Logger) Notice(format string, args ...interface{}) {
	logger.logMessage(LogLevelNotice, format, Ocname{"", 0}, args...)
}

// IsNoticeEnabled returns true if NOTICE logging is enable for the caller
func (logger *Logger) IsNoticeEnabled() bool {
	return logger.isLogEnabled(LogLevelNotice)
}

// Info is called for log level INFO messages
func (logger *Logger) Info(format string, args ...interface{}) {
	logger.logMessage(LogLevelInfo, format, Ocname{"", 0}, args...)
}

// IsInfoEnabled returns true if INFO logging is enable for the caller
func (logger *Logger) IsInfoEnabled() bool {
	return logger.isLogEnabled(LogLevelInfo)
}

// Debug is called for log level DEBUG messages
func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.logMessage(LogLevelDebug, format, Ocname{"", 0}, args...)
}

// IsDebugEnabled returns true if DEBUG logging is enable for the caller
func (logger *Logger) IsDebugEnabled() bool {
	return logger.isLogEnabled(LogLevelDebug)
}

// Trace is called for log level TRACE messages
func (logger *Logger) Trace(format string, args ...interface{}) {
	logger.logMessage(LogLevelTrace, format, Ocname{"", 0}, args...)
}

// OCTrace is called for overseer messages
func (logger *Logger) OCTrace(format string, newOcname Ocname, args ...interface{}) {
	logger.logMessage(LogLevelTrace, format, newOcname, args...)
}

// IsTraceEnabled returns true if TRACE logging is enable for the caller
func (logger *Logger) IsTraceEnabled() bool {
	return logger.isLogEnabled(LogLevelTrace)
}

// LogMessageSource is for the netfilter interface functions written in C
// and our LogWriter type that can be created and passed to anything that
// expects an object with output stream support. The logging source is passed
// directly rather than determined from the call stack.
func (logger *Logger) LogMessageSource(level int32, source string, format string, args ...interface{}) {
	if level > logger.getLogLevel(source, "") {
		return
	}

	if len(args) == 0 {
		fmt.Printf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], source, format)
	} else {
		buffer := logFormatter(format, Ocname{"", 0}, args...)
		if len(buffer) == 0 {
			return
		}
		fmt.Printf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], source, buffer)
	}
}

// IsLogEnabledSource returns true if logging is enabled at the argumented level for the argumented source
func (logger *Logger) IsLogEnabledSource(level int32, source string) bool {
	lvl := logger.getLogLevel(source, "")
	return (lvl >= level)
}

// Couldn't find the function being called anywhere
// EnableTimestamp enables the elapsed time in output
func (logger *Logger) EnableTimestamp() {
	logger.timestampEnabled = true
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
		ptr, stat := logger.logLevelMap[functionName]
		logger.logLevelLocker.RUnlock()
		if stat {
			return atomic.LoadInt32(ptr)
		}
	}

	if len(packageName) != 0 {
		logger.logLevelLocker.RLock()
		ptr, stat := logger.logLevelMap[packageName]
		logger.logLevelLocker.RUnlock()
		if stat {
			return atomic.LoadInt32(ptr)
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
		buffer := fmt.Sprintf(format, args...)
		return buffer
	}
	// return empty string when a repeat is limited
	return ""
}

// isLogEnabled returns true if logging is enabled for the caller at the specified level, false otherwise
func (logger *Logger) isLogEnabled(level int32) bool {
	_, _, packageName, functionName := findCallingFunction()
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
	_, _, packageName, functionName := findCallingFunction()

	if level > logger.getLogLevel(packageName, functionName) {
		return
	}

	// Make sure we have struct variables populated
	if (newOcname == Ocname{}) {
		// fmt.Printf("ERROR: logFormatter OC verb missing arguments:%s", format)
		return
	} else { //Handle %OC
		buffer := logFormatter(format, newOcname, args...)
		if len(buffer) == 0 {
			return
		}
		fmt.Printf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], packageName, buffer)
	}

	if len(args) == 0 {
		fmt.Printf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], packageName, format)
	} else {
		buffer := logFormatter(format, newOcname, args...)
		if len(buffer) == 0 {
			return
		}
		fmt.Printf("%s%-6s %18s: %s", logger.getPrefix(), logLevelName[level], packageName, buffer)
	}
}

// This function uses runtime.Callers to get the call stack to determine the calling function
// Our public function heirarchy is implemented so the caller is always at the 5th frame
// Frame 0 = runtime.Callers
// Frame 1 = findCallingFunction
// Frame 2 = logMessage / isLogEnabled
// Frame 3 = Warn, Info / IsWarnEnabled, IsInfoEnabled (etc...)
// Frame 4 = the function that actually called logger.Warn, logger.Info, logger.IsWarnEnabled, logger.IsInfoEnabled, etc...

// Here is an example of what we expect to see in the calling function frame:
// FILE: /home/username/golang/src/github.com/untangle/packetd/services/dict/dict.go
// FUNC: github.com/untangle/packetd/services/dict.cleanDictionary
// LINE: 827
// We find the last / in caller.Function and use the entire string as the function name (dict.cleanDictionary)
// We find the dot in the function name and use the left side as the package name (dict)
func findCallingFunction() (string, int, string, string) {
	// create a single entry array to hold the 5th stack frame and pass 4 as the
	// number of frames to skip over so we get the single stack frame we need
	stack := make([]uintptr, 1)
	count := runtime.Callers(4, stack)
	if count != 1 {
		return "unknown", 0, "unknown", "unknown"
	}

	// get the frame object for the caller
	frames := runtime.CallersFrames(stack)
	caller, _ := frames.Next()

	var functionName string
	var packageName string

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

	return caller.File, caller.Line, packageName, functionName
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

// SearchSourceLogLevel returns the log level for the specified source
// or a negative value if the source does not exist
func (logger *Logger) SearchSourceLogLevel(source string) int32 {
	logger.logLevelLocker.RLock()
	ptr, stat := logger.logLevelMap[source]
	logger.logLevelLocker.RUnlock()
	if !stat {
		return -1
	}

	return atomic.LoadInt32(ptr)
}

// AdjustSourceLogLevel sets the log level for the specified source and returns
// the previous level or a negative value if the source does not exist
func (logger *Logger) AdjustSourceLogLevel(source string, level int32) int32 {
	logger.logLevelLocker.RLock()
	ptr, stat := logger.logLevelMap[source]
	logger.logLevelLocker.RUnlock()
	if !stat {
		logger.Notice("Adding log level source NAME:%s LEVEL:%d\n", source, level)
		ptr = new(int32)
		atomic.StoreInt32(ptr, -1)
		logger.logLevelLocker.Lock()
		logger.logLevelMap[source] = ptr
		logger.logLevelLocker.Unlock()
	}

	prelvl := atomic.LoadInt32(ptr)
	atomic.StoreInt32(ptr, level)
	return prelvl
}

// FindLogLevelValue returns the numeric log level for the arugmented name
// or a negative value if the level is not valid
func FindLogLevelValue(source string) int32 {
	for levelvalue, levelname := range logLevelName {
		if strings.Compare(strings.ToUpper(levelname), strings.ToUpper(source)) == 0 {
			return (int32(levelvalue))
		}
	}

	return -1
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

// GenerateReport is called to create a dynamic HTTP page that shows all debug sources
func (logger *Logger) GenerateReport(buffer *bytes.Buffer) {
	logger.logLevelLocker.RLock()
	defer logger.logLevelLocker.RUnlock()

	// create a sorted list of the log level names
	namelist := make([]string, 0, len(logger.logLevelMap))
	for name := range logger.logLevelMap {
		namelist = append(namelist, name)
	}
	sort.Strings(namelist)

	buffer.WriteString("<TABLE BORDER=2 CELLPADDING=4 BGCOLOR=#EEEEEE>\r\n")
	buffer.WriteString("<TR><TH COLSPAN=2>Logger Source Levels</TH></TR>\r\n")
	buffer.WriteString("<TR><TD><B>Logger Source</B></TD><TD><B>Log Level</B></TD></TR>\r\n")

	for _, name := range namelist {
		ptr := logger.logLevelMap[name]
		buffer.WriteString("<TR><TD><TT>")
		buffer.WriteString(name)
		buffer.WriteString("</TT></TD><TD><TT>")
		buffer.WriteString(FindLogLevelName(atomic.LoadInt32(ptr)))
		buffer.WriteString("</TT></TD></TR>\n\n")
	}

	buffer.WriteString("</TABLE>\r\n")
}
