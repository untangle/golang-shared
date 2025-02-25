package settings

import (
	"fmt"
	"os"
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
)

var openWRTFileToNativeEOS = map[string]string{
	"/etc/config/categories.json": "/usr/share/bctid/categories.json",
	"/etc/config/appstate.json":   "/mnt/flash/mfw-settings/appstate.json",
	"/etc/config/settings.json":   "/mnt/flash/mfw-settings/settings.json",
	"/etc/config/current.json":    "/mnt/flash/mfw-settings/current.json",
	"/etc/config/default.json":    "/mnt/flash/mfw-settings/default.json",
	"/etc/config/uid":             "/mnt/flash/mfw-settings/uid",
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

func (f *FilenameLocator) getPlatformFileName(filename string) (string, error) { // Check if we are in native mode, most likely since there is only native and OpenWRT mode
	if f.fileExists(nativeEOSIndicatorFile) { // In Native mode, do translation
		if nativePath, exists := openWRTFileToNativeEOS[filename]; exists {
			if !f.fileExists(nativePath) {
				return nativePath, &NoFileAtPath{name: nativePath}
			}
			return nativePath, nil
		}
		return filename, fmt.Errorf("In Native mode, not file translation found for %v", filename)
	}
	// In OpenWRT mode, no translation since the defaults are openwrt paths.
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
