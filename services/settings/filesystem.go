package settings

import (
	"fmt"
	"os"
	"strings"
)

// FilenameLocator finds files on the local filesytem, allowing the
// system to be in hybrid or non-hybid mode and concealing the
// diffrences.
type FilenameLocator struct {
	fileExists func(filename string) bool
}

const (
	// Present of file indicates we are in native mode
	nativeEOSIndicatorFile = "/etc/EfwNativeEos"

	// Standard prefix for natic EOS
	nativeEOSPrefix = "/mnt/flash/mfw-settings/"

	// Standard prefix for OpenWRT
	openWRTPrefix = "/etc/config/"
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

// getPlatformFileName translates the filename to the appropriate path based on platform.
// it assumes that called uses default to OpenWRT platform. Only paths in mappings or paths which
// starts with /etc/config are translated.
func (f *FilenameLocator) getPlatformFileName(filename string) (string, error) { // Check if we are in native mode, most likely since there is only native and OpenWRT mode
	if f.fileExists(nativeEOSIndicatorFile) { // In EOS mode, try maping
		if nativePath, exists := openWRTFileToNativeEOS[filename]; exists {
			if !f.fileExists(nativePath) {
				return nativePath, &NoFileAtPath{name: nativePath}
			} else {
				return nativePath, nil
			}
			// Still on EOS, if file contains /etc/config then translate, otherwise return
		} else if strings.Contains(filename, openWRTPrefix) {
			nativePath := nativeEOSPrefix + filename[strings.LastIndex(filename, "/")+1:]
			if !f.fileExists(nativePath) {
				return nativePath, &NoFileAtPath{name: nativePath}
			} else {
				return nativePath, nil
			}
		}
	}
	// On OpenWRT, no translation needed
	if !f.fileExists(filename) {
		return filename, &NoFileAtPath{name: filename}
	}
	return filename, nil
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
