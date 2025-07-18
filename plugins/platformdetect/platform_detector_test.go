package platformdetect

import (
	"testing"

	"testing/fstest"

	"github.com/untangle/golang-shared/plugins"
	"github.com/untangle/golang-shared/plugins/types"
)

func TestPlatformFilter(t *testing.T) {
	tests := []struct {
		name            string
		files           []string
		metadata        []any
		expected        bool
		currentPlatform types.Platform
	}{
		{
			name:            "EOS platform, no metadata",
			files:           []string{types.EOS.IndicatorFilename},
			metadata:        nil,
			expected:        true,
			currentPlatform: types.EOS,
		},
		{
			name:            "OpenWrt platform, no metadata",
			files:           []string{types.OpenWrt.IndicatorFilename},
			metadata:        nil,
			expected:        true,
			currentPlatform: types.OpenWrt,
		},
		{
			name:            "Unclassified platform, no metadata",
			files:           []string{},
			metadata:        nil,
			expected:        true,
			currentPlatform: types.Unclassified,
		},
		{
			name:  "EOS platform, only on EOS",
			files: []string{types.EOS.IndicatorFilename},
			metadata: []any{PlatformSpec{
				OnlyOn: []types.Platform{types.EOS},
			}},
			expected:        true,
			currentPlatform: types.EOS,
		},
		{
			name:  "EOS platform, only on OpenWrt",
			files: []string{types.EOS.IndicatorFilename},
			metadata: []any{PlatformSpec{
				OnlyOn: []types.Platform{types.OpenWrt},
			}},
			expected:        false,
			currentPlatform: types.EOS,
		},
		{
			name:  "OpenWrt platform, excludes OpenWrt",
			files: []string{types.OpenWrt.IndicatorFilename},
			metadata: []any{PlatformSpec{
				Excludes: []types.Platform{types.OpenWrt},
			}},
			expected:        false,
			currentPlatform: types.OpenWrt,
		},
		{
			name:  "No platform match",
			files: []string{"/etc/something_else"},
			metadata: []any{PlatformSpec{
				OnlyOn: []types.Platform{types.OpenWrt},
			}},
			expected:        false,
			currentPlatform: types.Unclassified,
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

			filter := NewPlatformFilter(GetCurrentPlatform(fs))

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
