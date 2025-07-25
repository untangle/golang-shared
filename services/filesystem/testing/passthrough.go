package testing

import (
	"io/fs"
	"os"
)

type PassthroughPathGetter struct {
}

// GetPathOnPlatform is a passthrough/identity function for testing.
func (p *PassthroughPathGetter) GetPathOnPlatform(
	path string) (string, error) {
	return path, nil
}

// CurDirStatFS returns a fs.StatFS rooted in the current directory.
func CurDirStatFS() fs.StatFS {
	dir, _ := os.Getwd()
	return os.DirFS(dir).(fs.StatFS)
}
