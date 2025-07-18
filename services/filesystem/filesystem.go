package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/untangle/golang-shared/platform"
)

// Custom FS embedding golang's standard file system object.
type Filesystem struct {
	fs.FS
	platform platform.HostType
}

// Interface for an object that
// can lookup files actual location
// in the Filesystem. Should not be used
// to grab a filename, then os.Open().
// Use the FS for that. The same struct
// will be doing both
type FileSeeker interface {
	GetPathOnPlatform(string) (string, error)
}

// NewFileSystem returns a new Filesystem object
func NewFileSystem(fs fs.FS, p platform.HostType) *Filesystem {
	return &Filesystem{
		FS:       fs,
		platform: p,
	}
}

// Open opens a file on the filesystem
func (f *Filesystem) Open(n string) (fs.File, error) {
	nameOnPlat, err := f.GetPathOnPlatform(n)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	return f.FS.Open(nameOnPlat)
}

// Stat stats a file. An FS isn't required to implement Stat,
// but it does always implement Open. Use stat if it's implemented
// otherwise use open
func (f *Filesystem) Stat(n string) (fs.FileInfo, error) {
	statFS, ok := f.FS.(fs.StatFS)
	if !ok {
		// Underlying FS doesn't implement Stat.
		// Use open and grab the stat off the whole file
		// returned
		file, err := f.Open(n)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return file.Stat()
	}

	nameOnPlat, err := f.GetPathOnPlatform(n)
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	return statFS.Stat(nameOnPlat)
}

func (f *Filesystem) FileExists(n string) bool {
	return fileExists(n)
}

func (f *Filesystem) GetPathOnPlatform(p string) (string, error) {
	if f.platform.Equals(platform.OpenWrt) {
		return p, nil
	}

	if nativePath, ok := f.platform.UniquelyMappedFiles[p]; ok {
		if !fileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	}

	if strings.Contains(p, platform.OpenWrt.SettingsDirPath) {
		nativePath := filepath.Join(f.platform.SettingsDirPath, p[strings.LastIndex(p, "/")+1:])

		if !fileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	}

	return p, nil
}

// GetPathOnPlatformBad is a temporary function to be used by the settings package before
// an FS can be provided to it. Should not be used outside of that
func GetPathOnPlatformBad(p string) (string, error) {
	unmodifiedFS := os.DirFS("/")
	fs := NewFileSystem(unmodifiedFS, platform.DetectPlatform(unmodifiedFS.(fs.StatFS)))
	return fs.GetPathOnPlatform(p)

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

// fileExists returns true if we can Stat the filename. We don't
// distinguish between various kinds of errors, but do log them, on
// the theory that if you can't Stat the filename, for most purposes,
// that is the same as it not existing, and isn't a common case.
func fileExists(fname string) bool {
	return fileExistsInFs(
		fname,
		os.DirFS("/").(fs.StatFS))
}

// fileExistsInFs takes an fs.StatFS and returns true if the file
// exists at the path fname.
func fileExistsInFs(fname string, fs fs.StatFS) bool {
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
