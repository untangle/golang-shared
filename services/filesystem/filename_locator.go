package filesystem

import "strings"

// FilenameLocator finds files on the local filesystem, allowing the
// system to be EOS or OpenWrt mode and concealing the
// differences.
type FilenameLocator struct {
	fileExists func(filename string) bool
}

func NewFilenameLocator() {

}

// getPlatformFileName translates the filename to the appropriate path based on platform.
// it assumes that called uses default to OpenWRT platform. Only paths in mappings or paths which
// starts with /etc/config are translated.
func (f *FilenameLocator) getPlatformFileName(filename string) (string, error) { // Check if we are in native mode, most likely since there is only native and OpenWRT mode
	if f.fileExists(nativeEOSIndicatorFile) { // In EOS mode, try mapping
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
	// On OpenWRT or no translation needed
	if !f.fileExists(filename) {
		return filename, &NoFileAtPath{name: filename}
	}
	return filename, nil
}

// LocateFile locates the input filename on the filesystem,
// automatically translating it to EOS filenames when needed.
func (f *FilenameLocator) LocateFile(filename string) (string, error) {
	if f.fileExists(filename) {
		return filename, nil
	}
	return f.getPlatformFileName(filename)

}
