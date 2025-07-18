package settings

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/untangle/golang-shared/plugins/types"
)

// FilenameLocator finds files on the local filesystem, allowing the
// system to be EOS or OpenWrt mode and concealing the
// differences.
type FilenameLocator struct {
	fileExists func(filename string) bool
	platform   types.Platform
}

func NewFilenameLocator(platform types.Platform) *FilenameLocator {
	return &FilenameLocator{
		fileExists: FileExists,
	}
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
	return FileExistsInFS(
		fname,
		os.DirFS("/").(fs.StatFS))
}

// FileExistsInFS takes an fs.StatFS and returns true if the file
// exists at the path fname.
func FileExistsInFS(fname string, fs fs.StatFS) bool {
	if len(fname) == 0 {
		return false
	} else if fname[0] == '/' {
		fname = fname[1:]
	}

	lenOfFname := len(fname)
	if fname[lenOfFname-1] == '/' {
		fname = fname[0 : lenOfFname-1]
	}

	if _, err := fs.Stat(fname); err != nil {
		if !os.IsNotExist(err) {
			// Use fmt.Fprintf here because the logger may or may not
			// exist at this time.
			fmt.Fprintf(os.Stderr,
				"Unexpected error code from os.Stat: %v\n",
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
	// File lookups are mapped from the OpenWRT filenames. No additional translation needed.
	fmt.Printf("\n wat \n")
	if f.platform.Equals(types.OpenWrt) {
		return filename, nil
	}

	if nativePath, ok := f.platform.AdditionalFileMappings[filename]; ok {
		fmt.Printf("\nadditional mappings\n")
		if !f.fileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	} else if strings.Contains(filename, types.OpenWrt.SettingsDirPath) {
		nativePath := filepath.Join(f.platform.SettingsDirPath, filename[strings.LastIndex(filename, "/")+1:])

		if !f.fileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	}

	return filename, nil
}

// LocateFile locates the input filename on the filesystem,
// automatically translating it to EOS filenames when needed.
func (f *FilenameLocator) LocateFile(filename string) (string, error) {
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
