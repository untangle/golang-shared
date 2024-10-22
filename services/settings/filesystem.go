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

	nativeEOSIndicatorFile = "/etc/EfwNativeEos"
)

var openWRTFileToNativeEOS = map[string]string{
	"/etc/config/categories.json": "/usr/share/bctid/categories.json",
}

// NoFileAtPath is an error for if a file doesn't exist. In this case
// platform detection may have gone okay but we didn't see the file.
type NoFileAtPath struct {
	name string
}

// Error returns the error string
func (n *NoFileAtPath) Error() string {
	return fmt.Sprintf("no file at path: %s", n.name)
}

var _ error = &NoFileAtPath{}

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

func (f *FilenameLocator) getPlatformFileName(filename string) (string, error) {
	// Determine platform
	var newFileName string
	if f.fileExists(kernelModeSettingsPrefix) {
		// Kernel/OpenWRT mode
		newFileName = kernelModeSettingsPrefix + "/" + filename[strings.LastIndex(filename, "/")+1:]
	} else {
		// Hybrid mode
		if f.fileExists(nativeEOSIndicatorFile) {
			// Packetd running natively EOS mode
			if nativePath, exists := openWRTFileToNativeEOS[filename]; exists {
				return nativePath, nil
			}

			// Fall through to Hybrid mode handling
		}

		if !strings.HasPrefix(filename, kernelModeSettingsPrefix) {
			// Not a config file, use generic prefix
			newFileName = filepath.Join(hybridModeGenericPrefix, filename)
		} else {
			newFileName = hybridModeSettingsPrefix + "/" + filename[strings.LastIndex(filename, "/")+1:]
		}
	}
	if !f.fileExists(newFileName) {
		// File doesn't exist, but the caller may not care
		return newFileName, &NoFileAtPath{newFileName}
	}
	return newFileName, nil
}

// LocateFile locates the input filename on the filesystem,
// automatically translating it to hybrid mode filenames when needed.
func (f *FilenameLocator) LocateFile(filename string) (string, error) {
	if f.fileExists(filename) {
		return filename, nil
	}
	return f.getPlatformFileName(filename)

}

var defaultLocator = &FilenameLocator{
	fileExists: FileExists,
}

// LocateFile calls FilenameLocator.LocateFile on the default filename
// locator.
func LocateFile(filename string) (string, error) {
	return defaultLocator.LocateFile(filename)
}
