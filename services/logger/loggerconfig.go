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
	// I think these can be removed? they are only accessible from a single function
	file *os.File
	info os.FileInfo
}

type LogLevel struct {
	Name string `json:logname`
	Id   uint8
}

func (conf *LoggerConfig) GetLogID(name string) uint8 {
	switch name {
	case "EMERG":
		return 0
	case "ALERT":
		return 1
	case "CRIT":
		return 2
	case "ERROR":
		return 3
	case "WARN":
		return 4
	case "NOTICE":
		return 5
	case "INFO":
		return 6
	case "DEBUG":
		return 7
	case "TRACE":
		return 8
	default:
		return 9
	}
}

// loadLoggerConfig loads the logger configuration file
func (conf *LoggerConfig) LoadConfigFromFile() []byte {
	if conf.FileLocation == "" {
		GetLoggerInstance().Err("FileLocation must be set\n")
		return nil
	}

	var err error

	// open the logger configuration file
	conf.file, err = os.Open(conf.FileLocation)
	fmt.Print(err)

	// if there was an error create the default config and try the open again
	if err != nil {
		conf.writeLoggerConfigToJSON()
		conf.file, err = os.Open(conf.FileLocation)

		// if there is still an error we are out of options
		if err != nil {
			GetLoggerInstance().Err("Unable to load Log configuration file: %s\n", conf.FileLocation)
			return nil
		}
	}

	// make sure the file gets closed
	defer conf.file.Close()

	// get the file status
	conf.info, err = conf.file.Stat()
	if err != nil {
		GetLoggerInstance().Err("Unable to query file information\n")
		return nil
	}
	data := make([]byte, conf.info.Size())
	len, err := conf.file.Read(data)

	if (err != nil) || (len < 1) {
		GetLoggerInstance().Err("Unable to read Log configuration\n")
		return nil
	}

	return data
}

// split -> Mock Json pass to the function below
// read the raw configuration json from the file
func (conf *LoggerConfig) LoadConfigFromJSON(data []byte) {
	serviceMap := make(map[string]LogLevel)
	conf.LogLevelMap = make(map[string]LogLevel)

	// unmarshal the configuration into a map
	err := json.Unmarshal(data, &serviceMap)
	if err != nil {
		GetLoggerInstance().Err("Unable to parse Log configuration\n")
		return
	}
}

func (conf *LoggerConfig) writeLoggerConfigToJSON() {
	GetLoggerInstance().Alert("Log configuration not found. Creating default File: %s\n", conf.FileLocation)

	// convert the config map to a json object
	jstr, err := json.MarshalIndent(conf.LogLevelMap, "", "")
	if err != nil {
		GetLoggerInstance().Alert("Log failure creating default configuration: %s\n", err.Error())
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
