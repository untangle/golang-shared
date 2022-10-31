package logger

// LogWriter is used to send an output stream to the Log facility
type LogWriter struct {
	buffer []byte
	source string
}

// DefaultLogWriter creates an io Writer to steam output to the Log facility
func DefaultLogWriter(name string) *LogWriter {
	writer := new(LogWriter)
	writer.buffer = make([]byte, 0)
	writer.source = name
	return writer
}

// Write takes written data and stores it in a buffer and writes to the log when a line feed is detected
func (writer *LogWriter) Write(p []byte) (int, error) {
	logger := Logger{}
	for _, b := range p {
		writer.buffer = append(writer.buffer, b)
		if b == '\n' {
			logger.LogMessageSource(LogLevelInfo, writer.source, string(writer.buffer))
			writer.buffer = make([]byte, 0)
		}
	}

	return len(p), nil
}
