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
func (conf *LoggerConfig) LoadConfigFromFile() []byte {
	if conf.FileLocation == "" {
		return nil
	}

	var err error
	var file *os.File
	var info os.FileInfo

	// open the logger configuration file
	file, err = os.Open(conf.FileLocation)
	fmt.Print(err)

	// if there was an error create the default config and try the open again
	if err != nil {
		conf.writeLoggerConfigToJSON()
		file, err = os.Open(conf.FileLocation)

		// if there is still an error we are out of options
		if err != nil {
			return nil
		}
	}

	// make sure the file gets closed
	defer file.Close()

	// get the file status
	info, err = file.Stat()
	if err != nil {
		return nil
	}
	data := make([]byte, info.Size())
	len, err := file.Read(data)

	if (err != nil) || (len < 1) {
		return nil
	}

	return data
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
func (conf *LoggerConfig) writeLoggerConfigToJSON() {

	// convert the config map to a json object
	jstr, err := json.MarshalIndent(conf.LogLevelMap, "", "")
	if err != nil {
		return
	}

	// create the logger configuration file
	file, err := os.Create(conf.FileLocation)
	if err != nil {
		return
	}

	// write the default configuration and close the file
	file.Write(jstr)
	file.Close()
}
