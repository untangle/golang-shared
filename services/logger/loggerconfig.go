package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"

	"github.com/untangle/golang-shared/util"
)

// LoggerConfig struct retains information about the where the log level map is stored, default log levels and writer that should be used
type LoggerConfig struct {
	FileLocation string
	LogLevelMap  map[string]string
	OutputWriter io.Writer
	// I think these can be removed? they are only accessible from a single function
	file *os.File
	info os.FileInfo
}

// validateConfig ensures all the log levels set in config.LogLevelMap are valid
func (conf *LoggerConfig) Validate() error {
	if conf.FileLocation == "" {
		return errors.New("FileLocation must be set")
	}
	for key, value := range conf.LogLevelMap {
		if !util.ContainsString(logLevelName[:], value) {
			return fmt.Errorf("%s is using an incorrect log level: %s", key, value)
		}
	}
	return nil
}

// loadLoggerConfig loads the logger configuration file
func (conf *LoggerConfig) LoadConfigFromFile() []byte {

	var err error

	// open the logger configuration file
	conf.file, err = os.Open(conf.FileLocation)

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
	serviceMap := make(map[string]string)
	conf.LogLevelMap = make(map[string]*int32)

	// unmarshal the configuration into a map
	err := json.Unmarshal(data, &serviceMap)
	if err != nil {
		GetLoggerInstance().Err("Unable to parse Log configuration\n")
		return
	}

	// put the name/string pairs from the file into the name/int lookup map we us in the Log function
	for cfgname, cfglevel := range serviceMap {
		// ignore any comment strings that start and end with underscore
		if strings.HasPrefix(cfgname, "_") && strings.HasSuffix(cfgname, "_") {
			continue
		}

		// find the index of the logLevelName that matches the configured level
		found := false
		for levelvalue, levelname := range logLevelName {
			// if the string matches the level will be the index of the matched name
			if strings.Compare(levelname, strings.ToUpper(cfglevel)) == 0 {
				conf.LogLevelMap[cfgname] = new(int32)
				atomic.StoreInt32(conf.LogLevelMap[cfgname], int32(levelvalue))
				found = true
			}
		}
		if !found {
			GetLoggerInstance().Warn("Invalid Log configuration entry: %s=%s\n", cfgname, cfglevel)
		}
	}
}

func (conf *LoggerConfig) writeLoggerConfigToJSON() {
	GetLoggerInstance().Alert("Log configuration not found. Creating default File: %s\n", logger.config.FileLocation)

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
