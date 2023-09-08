package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
)

type CmdAlertDetail struct {
	severity Alerts.AlertSeverity
	logType  Alerts.AlertType
}

var CmdAlertDefaultSetup = map[int32]CmdAlertDetail{
	LogLevelCrit: {
		severity: Alerts.AlertSeverity_CRITICAL,
		logType:  Alerts.AlertType_CRITICALERROR,
	},
	LogLevelErr: {
		severity: Alerts.AlertSeverity_ERROR,
		logType:  Alerts.AlertType_CRITICALERROR,
	},
}

// LoggerConfig struct retains information about the where the log level map is stored, default log levels and writer that should be used
type LoggerConfig struct {
	FileLocation    string
	LogLevelHighest int32
	LogLevelMap     map[string]LogLevel
	OutputWriter    io.Writer
	CmdAlertSetup   map[int32]CmdAlertDetail
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

	// return error if one exists
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

	return conf.LoadConfigFromJSON(data)
}

// split -> Mock Json pass to the function below
// read the raw configuration json from the file
func (conf *LoggerConfig) LoadConfigFromJSON(data []byte) error {
	conf.LogLevelMap = make(map[string]LogLevel)

	// unmarshal the configuration into a map
	err := json.Unmarshal(data, &conf.LogLevelMap)
	if err != nil {
		return err
	}

	// set the highest log level
	conf.SetLogLevelHighest()
	return nil
}

// SetLogLevelHighest will set the highest log level in the log config
func (conf *LoggerConfig) SetLogLevelHighest() {
	for _, v := range conf.LogLevelMap {
		if v.GetId() > conf.LogLevelHighest {
			conf.LogLevelHighest = v.GetId()
		}
	}
}

// SaveConfig will write the current loglevelmap to disk
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
	conf.SetLogLevelHighest()
}

// SetLogLevel can set the log level in the log config
func (conf *LoggerConfig) SetLogLevel(key string, newLevel LogLevel) {
	conf.LogLevelMap[key] = newLevel
	if newLevel.GetId() > conf.LogLevelHighest {
		conf.LogLevelHighest = newLevel.GetId()
	}
}

// removeConfigFile will remove the config file from disk
func (conf *LoggerConfig) removeConfigFile() {
	os.Remove(conf.FileLocation)
}
