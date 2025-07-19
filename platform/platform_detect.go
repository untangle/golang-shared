package platform

import (
	"io/fs"
)

// DetectPlatform detect the current platform with the provided FS
func DetectPlatform(fs fs.StatFS) HostType {
	for _, p := range platforms {
		_, err := fs.Stat(p.IndicatorFilename)
		if err == nil {
			return p
		}
	}
	return Unclassified
}
