package platform

import (
	"io/fs"
	"os"
)

// DetectPlatform detect the current platform with the provided FS
func DetectPlatform() HostType {
	unmodifiedFS := os.DirFS("/").(fs.StatFS)
	return DetectPlatformFromFS(unmodifiedFS)
}

// DetectPlatform detect the current platform with the provided FS
func DetectPlatformFromFS(fs fs.StatFS) HostType {
	for _, p := range platforms {
		_, err := fs.Stat(p.IndicatorFilename)
		if err == nil {
			return p
		}
	}
	return Unclassified
}
