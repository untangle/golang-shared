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
type PlatformAwareFileSystem struct {
	fs.FS
	platform platform.HostType
}

// Interface for an object that
// can lookup files actual location
// in the Filesystem. Should not be used
// to grab a filename, then os.Open().
// Use the FS for that. The same struct
// will be doing both
type PathOnPlatformGetter interface {
	GetPathOnPlatform(string) (string, error)
}

// NewPlatformAwareFileSystem returns a new Filesystem object
func NewPlatformAwareFileSystem(fs fs.FS, p platform.HostType) *PlatformAwareFileSystem {
	return &PlatformAwareFileSystem{
		FS:       fs,
		platform: p,
	}
}

// Open opens a file on the filesystem
func (f *PlatformAwareFileSystem) Open(n string) (fs.File, error) {
	nameOnPlat, err := f.GetPathOnPlatform(n)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	nameOnPlat = sanitizePath(nameOnPlat)

	return f.FS.Open(nameOnPlat)
}

// Stat stats a file. An FS isn't required to implement Stat,
// but it does always implement Open. Use stat if it's implemented
// otherwise use open
func (f *PlatformAwareFileSystem) Stat(n string) (fs.FileInfo, error) {
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
	nameOnPlat = sanitizePath(nameOnPlat)

	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	return statFS.Stat(nameOnPlat)
}

// sanitizePath sanitizes a path by stripping off
// any leading /. Absolute paths will cause fs.Open() to fail.
// It appends the directory the Filesystem is created
// with to the path provided. This results in a bad path.
func sanitizePath(p string) string {
	return strings.Trim(p, "/")
}

func (f *PlatformAwareFileSystem) FileExists(n string) bool {
	return FileExistsInFs(n, f.FS.(fs.StatFS))
}

func (f *PlatformAwareFileSystem) GetPathOnPlatform(p string) (string, error) {
	if f.platform.Equals(platform.OpenWrt) {
		return p, nil
	}

	if nativePath, ok := f.platform.UniquelyMappedFiles[p]; ok {
		if !f.FileExists(nativePath) {
			return nativePath, &NoFileAtPath{name: nativePath}
		} else {
			return nativePath, nil
		}
	}

	nativePath := p
	if strings.Contains(p, platform.OpenWrt.SettingsDirPath) {
		nativePath = filepath.Join(f.platform.SettingsDirPath, p[strings.LastIndex(p, "/")+1:])
	}

	if !f.FileExists(nativePath) {
		return nativePath, &NoFileAtPath{name: nativePath}
	}

	return nativePath, nil
}

// GetPathOnPlatformBad is a temporary function to be used by the settings package before
// an FS can be provided to it. Should not be used outside of that
func GetPathOnPlatformBad(p string) (string, error) {
	unmodifiedFS := os.DirFS("/")
	fs := NewPlatformAwareFileSystem(unmodifiedFS, platform.DetectPlatform())
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

// fileExistsInFs takes an fs.StatFS and returns true if the file
// exists at the path fname.
func FileExistsInFs(fname string, fs fs.StatFS) bool {

	fname = sanitizePath(fname)
	if fname == "" {
		return false
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
