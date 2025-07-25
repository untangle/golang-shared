package logger

import (
	"bytes"
	"testing"
)

func TestSimpleLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewSimpleLogger(LevelDebug, buf)

	logger.Info("info message\n")
	logger.Debug("debug message\n")
	logger.Trace("trace message\n") // Should not be logged

	expected := "info message\ndebug message\n"
	if buf.String() != expected {
		t.Errorf("Log output mismatch. Expected: %q, Actual: %q", expected, buf.String())
	}

	buf.Reset()
	logger.Level = LevelInfo
	logger.Info("info message\n")
	logger.Debug("debug message\n") // Should now not be logged
	if buf.String() != "info message\n" {
		t.Errorf("Expected only info message, but also got debug message.  Actual: %q", buf.String())
	}

	buf.Reset()
	logger = NewSimpleLogger(LevelDebug, nil) // test default output
	logger.Debug("debug message, default output\n")
	if buf.String() != "" {
		t.Error("Expected empty buffer because default logger to Stderr")
	}

}
