package platformfilter

import (
	"testing"

	"testing/fstest"

	"github.com/untangle/golang-shared/platform"
	"github.com/untangle/golang-shared/plugins"
)

func TestPlatformFilter(t *testing.T) {
	tests := []struct {
		name            string
		files           []string
		metadata        []any
		expected        bool
		currentPlatform platform.HostType
	}{
		{
			name:            "EOS platform, no metadata",
			files:           []string{platform.EOS.IndicatorFilename},
			metadata:        nil,
			expected:        true,
			currentPlatform: platform.EOS,
		},
		{
			name:            "OpenWrt platform, no metadata",
			files:           []string{platform.OpenWrt.IndicatorFilename},
			metadata:        nil,
			expected:        true,
			currentPlatform: platform.OpenWrt,
		},
		{
			name:            "Unclassified platform, no metadata",
			files:           []string{},
			metadata:        nil,
			expected:        true,
			currentPlatform: platform.Unclassified,
		},
		{
			name:  "EOS platform, only on EOS",
			files: []string{platform.EOS.IndicatorFilename},
			metadata: []any{PlatformSpec{
				OnlyOn: []platform.HostType{platform.EOS},
			}},
			expected:        true,
			currentPlatform: platform.EOS,
		},
		{
			name:  "EOS platform, only on OpenWrt",
			files: []string{platform.EOS.IndicatorFilename},
			metadata: []any{PlatformSpec{
				OnlyOn: []platform.HostType{platform.OpenWrt},
			}},
			expected:        false,
			currentPlatform: platform.EOS,
		},
		{
			name:  "OpenWrt platform, excludes OpenWrt",
			files: []string{platform.OpenWrt.IndicatorFilename},
			metadata: []any{PlatformSpec{
				Excludes: []platform.HostType{platform.OpenWrt},
			}},
			expected:        false,
			currentPlatform: platform.OpenWrt,
		},
		{
			name:  "No platform match",
			files: []string{"/etc/something_else"},
			metadata: []any{PlatformSpec{
				OnlyOn: []platform.HostType{platform.OpenWrt},
			}},
			expected:        false,
			currentPlatform: platform.Unclassified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := fstest.MapFS{}

			for _, f := range tt.files {
				// We slice of the leading / because mapfs just looks
				// literally at the string as the filename, and we
				// need to strip off leading "/" for regular fs.FS
				// objects.
				fs[f[1:]] = &fstest.MapFile{}
			}

			filter := NewPlatformFilter(tt.currentPlatform)

			if filter.currentPlatform.Name != tt.currentPlatform.Name {
				t.Errorf("Incorrect platform detected. Expected %s, got %s",
					tt.currentPlatform.Name, filter.currentPlatform.Name)
			}

			var pc plugins.PluginConstructor
			actual := filter.IsRelevant(pc, tt.metadata...)
			if actual != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
