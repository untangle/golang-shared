package logger

import "strings"

// LogLevel struct retains the loglevel information in a string and int.
type LogLevel struct {
	Name string `json:"logname"`
	id   uint8
}

// GetId returns the numeric log level for the arugmented name
// or a negative value if the level is not valid
func (lvl *LogLevel) GetId() int32 {
	for levelvalue, levelname := range logLevelName {
		if strings.Compare(strings.ToUpper(levelname), strings.ToUpper(lvl.Name)) == 0 {
			return (int32(levelvalue))
		}
	}

	return -1
}
