package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilenameLocator finds files on the local filesytem, allowing the
// system to be in hybrid or non-hybid mode and concealing the
// diffrences.
type FilenameLocator struct {
	fileExists func(filename string) bool
}

const (
	// prefix for generic filepaths in hybrid mode
	hybridModeGenericPrefix = "/mfw"

	// kernel forwarding mode/BST container mode path prefix.
	kernelModeSettingsPrefix = "/etc/config"

	// prefix specifically for config files in hybrid mode
	hybridModeSettingsPrefix = "/mnt/flash/mfw-settings"
)

// FileExists returns true if we can Stat the filename. We don't
// distinguish between various kinds of errors, but do log them, on
// the theory that if you can't Stat the filename, for most purposes,
// that is the same as it not existing, and isn't a common case.
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

// LocateFile locates the input filename on the filesystem,
// automatically translating it to hybrid mode filenames when needed.
func (f *FilenameLocator) LocateFile(filename string) (string, error) {
	if f.fileExists(filename) {
		return filename, nil
	}
	return f.findEOSFileName(filename)

}

var defaultLocator = &FilenameLocator{
	fileExists: FileExists,
}

// LocateFile calls FilenameLocator.LocateFile on the default filename
// locator.
func LocateFile(filename string) (string, error) {
	return defaultLocator.LocateFile(filename)
}