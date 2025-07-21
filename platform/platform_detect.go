package platform

import (
	"io/fs"
	"os"
)

// DetectPlatform detect the current platform with the provided FS
func DetectPlatform() HostType {
	unmodifiedFS := os.DirFS("/").(fs.StatFS)
	for _, p := range platforms {
		_, err := unmodifiedFS.Stat(p.IndicatorFilename)
		if err == nil {
			return p
		}
	}
	return Unclassified
}
