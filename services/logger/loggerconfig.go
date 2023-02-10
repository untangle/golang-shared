package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// LoggerConfig struct retains information about the where the log level map is stored, default log levels and writer that should be used
type LoggerConfig struct {
	FileLocation string
	LogLevelMap  map[string]LogLevel
	OutputWriter io.Writer
}

// loadLoggerConfig loads the logger configuration file
func (conf *LoggerConfig) LoadConfigFromFile() error {
	if conf.FileLocation == "" {
		return fmt.Errorf("Logger config FileLocation is missing")
	}

	var err error
	var file *os.File
	var info os.FileInfo

	// open the logger configuration file
	file, err = os.Open(conf.FileLocation)
	fmt.Print(err)

	// if there was an error - return nil
	if err != nil {
		return err
	}

	// make sure the file gets closed
	defer file.Close()

	// get the file status
	info, err = file.Stat()
	if err != nil {
		return err
	}
	data := make([]byte, info.Size())
	len, err := file.Read(data)

	if (err != nil) || (len < 1) {
		return err
	}

	conf.LoadConfigFromJSON(data)
	return nil
}

// split -> Mock Json pass to the function below
// read the raw configuration json from the file
func (conf *LoggerConfig) LoadConfigFromJSON(data []byte) {
	conf.LogLevelMap = make(map[string]LogLevel)

	// unmarshal the configuration into a map
	err := json.Unmarshal(data, &conf.LogLevelMap)
	if err != nil {
		return
	}
}

// writeLoggerConfigToJSON will load the default config
func (conf *LoggerConfig) SaveConfig() {

	// convert the config map to a json object
	jstr, err := json.MarshalIndent(conf.LogLevelMap, "", "")
	if err != nil {
		fmt.Printf("Unable to unmarshal LogLevelMap: %s", err)
		return
	}

	// create the logger configuration file
	file, err := os.Create(conf.FileLocation)
	if err != nil {
		fmt.Printf("Unable to save file: %s, error: %s", conf.FileLocation, err)
		return
	}
	defer file.Close()

	// write the default configuration and close the file
	_, err = file.Write(jstr)
	if err != nil {
		fmt.Printf("Unable to write to file: %s, error: %s", conf.FileLocation, err)
		return
	}
}

// SetLogLevel can set the log level in the log config
func (conf *LoggerConfig) SetLogLevel(key string, newLevel LogLevel) {
	conf.LogLevelMap[key] = newLevel
}
