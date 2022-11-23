package logger

// LogWriter is used to send an output stream to the Log facility
type LogWriter struct {
	buffer []byte
	source string
	// defines the message logging level
	logLevel int32
}

// DefaultLogWriter creates an io Writer to steam output to the Log facility
func DefaultLogWriter(name string) *LogWriter {
	writer := new(LogWriter)
	writer.buffer = make([]byte, 0)
	writer.source = name
	writer.logLevel = LogLevelInfo

	return writer
}

// SetLogLevel allows to modify the log level of messages
func (writer *LogWriter) SetLogLevel(logLevel int32) {
	writer.logLevel = logLevel
}

// Write takes written data and stores it in a buffer and writes to the log when a line feed is detected
func (writer *LogWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		writer.buffer = append(writer.buffer, b)
		if b == '\n' {
			LogMessageSource(writer.logLevel, writer.source, string(writer.buffer))
			writer.buffer = make([]byte, 0)
		}
	}

	return len(p), nil
}
