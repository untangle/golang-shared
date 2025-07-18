package platformfilter

import (
	"testing"

	"testing/fstest"

	"github.com/untangle/golang-shared/plugins"
)

func TestPlatformFilter(t *testing.T) {
	tests := []struct {
		name            string
		files           []string
		metadata        []any
		expected        bool
		currentPlatform HostType
	}{
		{
			name:            "EOS platform, no metadata",
			files:           []string{EOS.indicatorFilename},
			metadata:        nil,
			expected:        true,
			currentPlatform: EOS,
		},
		{
			name:            "OpenWrt platform, no metadata",
			files:           []string{OpenWrt.indicatorFilename},
			metadata:        nil,
			expected:        true,
			currentPlatform: OpenWrt,
		},
		{
			name:            "Unclassified platform, no metadata",
			files:           []string{},
			metadata:        nil,
			expected:        true,
			currentPlatform: Unclassified,
		},
		{
			name:  "EOS platform, only on EOS",
			files: []string{EOS.indicatorFilename},
			metadata: []any{PlatformSpec{
				OnlyOn: []HostType{EOS},
			}},
			expected:        true,
			currentPlatform: EOS,
		},
		{
			name:  "EOS platform, only on OpenWrt",
			files: []string{EOS.indicatorFilename},
			metadata: []any{PlatformSpec{
				OnlyOn: []HostType{OpenWrt},
			}},
			expected:        false,
			currentPlatform: EOS,
		},
		{
			name:  "OpenWrt platform, excludes OpenWrt",
			files: []string{OpenWrt.indicatorFilename},
			metadata: []any{PlatformSpec{
				Excludes: []HostType{OpenWrt},
			}},
			expected:        false,
			currentPlatform: OpenWrt,
		},
		{
			name:  "No platform match",
			files: []string{"/etc/something_else"},
			metadata: []any{PlatformSpec{
				OnlyOn: []HostType{OpenWrt},
			}},
			expected:        false,
			currentPlatform: Unclassified,
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

			filter := NewPlatformFilter(fs)

			if filter.currentPlatform != tt.currentPlatform {
				t.Errorf("Incorrect platform detected. Expected %s, got %s",
					tt.currentPlatform.name, filter.currentPlatform.name)
			}

			var pc plugins.PluginConstructor
			actual := filter.IsRelevant(pc, tt.metadata...)
			if actual != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
