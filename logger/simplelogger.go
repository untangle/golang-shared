package logger

import (
	"fmt"
	"io"
	"os"
)

// SimpleLogger implements LoggerLevels and writes to an io.Writer.
type SimpleLogger struct {
	Level LogLevel
	out   io.Writer
}

type LogLevel int

const (
	LevelEmerg LogLevel = iota
	LevelAlert
	LevelCrit
	LevelErr
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
	LevelTrace
)

// NewSimpleLogger creates a SimpleLogger with configurable level and output.
func NewSimpleLogger(level LogLevel, out io.Writer) *SimpleLogger {
	if out == nil {
		out = os.Stderr
	}
	return &SimpleLogger{Level: level, out: out}
}

func (l *SimpleLogger) log(level LogLevel, format string, args ...interface{}) {
	if l.Level >= level {
		fmt.Fprintf(l.out, format, args...)
	}
}

func (l *SimpleLogger) Emerg(format string, args ...interface{}) {
	l.log(LevelEmerg, format, args...)
}

func (l *SimpleLogger) Alert(format string, args ...interface{}) {
	l.log(LevelAlert, format, args...)
}

func (l *SimpleLogger) Crit(format string, args ...interface{}) {
	l.log(LevelCrit, format, args...)
}

func (l *SimpleLogger) Err(format string, args ...interface{}) {
	l.log(LevelErr, format, args...)
}

func (l *SimpleLogger) Warn(format string, args ...interface{}) {
	l.log(LevelWarning, format, args...)
}

func (l *SimpleLogger) Notice(format string, args ...interface{}) {
	l.log(LevelNotice, format, args...)
}

func (l *SimpleLogger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)

}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

func (l *SimpleLogger) Trace(format string, args ...interface{}) {
	l.log(LevelTrace, format, args...)

}

func (l *SimpleLogger) OCWarn(format string, name string, limit int64, args ...interface{}) {
	l.log(LevelWarning, format, args...)

}

func (l *SimpleLogger) IsDebugEnabled() bool {
	return l.Level >= LevelDebug
}
