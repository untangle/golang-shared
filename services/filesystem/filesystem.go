package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/untangle/golang-shared/services/platformdetect"
)

// Custom FS embedding golang's standard file system object.
type Filesystem struct {
	fs.FS
	platform platformdetect.HostType
}

// NewFileSystem returns a new Filesystem object
func NewFileSystem(fs fs.FS, p platformdetect.HostType) *Filesystem {
	return &Filesystem{
		FS:       fs,
		platform: p,
	}
}

// Open opens a file on the filesystem
func (fs *Filesystem) Open(n string) (fs.File, error) {
	nameOnPlat, err := fs.getPathOnPlatform(n)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	return fs.FS.Open(nameOnPlat)
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

	nameOnPlat, err := f.getPathOnPlatform(n)
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	return statFS.Stat(nameOnPlat)
}

func (fs *Filesystem) getPathOnPlatform(p string) (string, error) {
	if fs.platform.Equals(platformdetect.OpenWrt) {
		return p, nil
	}

	if nativePath, ok := fs.platform.UniquelyMappedFiles[p]; ok {
		if !fileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	}

	if strings.Contains(p, platformdetect.OpenWrt.SettingsDirPath) {
		nativePath := filepath.Join(fs.platform.SettingsDirPath, p[strings.LastIndex(p, "/")+1:])

		if !fileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	}

	return p, nil
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
