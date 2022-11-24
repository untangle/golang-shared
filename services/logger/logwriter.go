package logger

// LogWriter is used to send an output stream to the Log facility
type LogWriter struct {
	buffer []byte
	source string

	// defines the log logging level
	logLevel LogLevel

	// Allows to update the log level based on the log message content.
	// If this is not set (== nil), this functionality is just ignored and the logLevel is used  for all messages.
	// The function will reveive 2 prameters:
	//    currentLevel - the current log level set for logging
	//    message - the actual message we want to log
	logLevelProcessor func(currentLevel LogLevel, message string) LogLevel
}

// DefaultLogWriter creates an io Writer to steam output to the Log facility
func DefaultLogWriter(name string) *LogWriter {
	writer := new(LogWriter)
	writer.buffer = make([]byte, 0)
	writer.source = name
	writer.logLevel = NewLogLevel("INFO")

	return writer
}

// SetLogLevel allows to modify the log level of messages
func (writer *LogWriter) SetLogLevel(logLevel LogLevel) error {
	if logLevel.GetId() < 0 {
		return ErrInvalidLogLevel
	}
	writer.logLevel = logLevel

	return nil
}

// SetLogLevelProcessor allows to set a function that can update the log level based on the message content.
func (writer *LogWriter) SetLogLevelProcessor(processor func(currentLevel LogLevel, message string) LogLevel) {
	writer.logLevelProcessor = processor
}

// Write takes written data and stores it in a buffer and writes to the log when a line feed is detected
func (writer *LogWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		writer.buffer = append(writer.buffer, b)
		if b == '\n' {
			writer.logMessage()
		}
	}

	return len(p), nil
}

func (writer *LogWriter) logMessage() {
	message := string(writer.buffer)

	logLevel := writer.logLevel
	if writer.logLevelProcessor != nil {
		// Update log level based on message content.
		logLevel = writer.logLevelProcessor(logLevel, message)
	}

	LogMessageSource(logLevel.GetId(), writer.source, message)
	writer.buffer = make([]byte, 0)
}
