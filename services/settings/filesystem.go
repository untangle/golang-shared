package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilenameLocator struct {
	fileExists func(filename string) bool
}

const (
	hybridModeGenericPrefix  = "/mfw"
	kernelModeSettingsPrefix = "/etc/config"
	hybridModeSettingsPrefix = "/mnt/flash/mfw-settings"
)

func FileExists(fname string) bool {
	if _, err := os.Stat(fname); err != nil {
		if !os.IsNotExist(err) {
			logger.Warn("Unexpected error code from os.Stat: %s",
				err)
		}
		return false
	}
	return true
}

func (f *FilenameLocator) findEOSFileName(filename string) (string, error) {
	if strings.HasPrefix(filename, kernelModeSettingsPrefix) {
		newFileName := strings.Replace(
			filename,
			kernelModeSettingsPrefix,
			hybridModeSettingsPrefix,
			1)
		if !f.fileExists(newFileName) {
			return "", fmt.Errorf("unable to find config file: %s", filename)
		}
		return newFileName, nil
	} else {
		newFileName := filepath.Join(
			hybridModeGenericPrefix,
			filename)
		if !f.fileExists(newFileName) {
			return "", fmt.Errorf(
				"unable to locate file: %s", filename)
		}
		return newFileName, nil
	}
}

func (f *FilenameLocator) LocateFile(filename string) (string, error) {
	if f.fileExists(filename) {
		return filename, nil
	}
	return f.findEOSFileName(filename)

}

var defaultLocator = &FilenameLocator{
	fileExists: FileExists,
}

func LocateFile(filename string) (string, error) {
	return defaultLocator.LocateFile(filename)
}
