package platformdetect

import (
	"io/fs"
	"slices"

	"github.com/untangle/golang-shared/plugins"
	"github.com/untangle/golang-shared/services/settings"
)

// HostType is a externally-opaque type for declaring a host
// type.
type HostType struct {
	indicatorFilename string
	name              string
}

var (
	EOS = HostType{
		indicatorFilename: "/etc/Eos-release",
		name:              "Eos",
	}
	OpenWrt = HostType{
		indicatorFilename: "/etc/openwrt_version",
		name:              "OpenWrt",
	}
	Unclassified = HostType{
		indicatorFilename: "",
		name:              "Unclassified",
	}
)

// PlatformFilter is a plugin predicate (declared in
// golang-shared/plugins) which filters out plugins that have
// PluginSpec metadata and don't apply to the current platform.
type PlatformFilter struct {
	currentPlatform HostType
}

// PlatformSpec is a specification of the platforms that apply to a
// plugin and should be listed in the registration metadata.
//
// If the OnlyOn list is nonempty and we are on a platform in that
// list, the plugin will be allowed to start and run.  If the Excludes
// list is nonempty and we are running on a platform in that list the
// plugin will not be allowed to start and run.
//
// Do not specify both OnlyOn and Excludes. Both lists should not be
// empty (just don't supply a PlatformSpec in this case) but if they
// are, the plugin will be run.
type PlatformSpec struct {
	OnlyOn   []HostType
	Excludes []HostType
}

// NewPlatformFilter creates a new platform filter, during
// construction we determine the platform from the filesystem.
func NewPlatformFilter(fs fs.StatFS) *PlatformFilter {
	platforms := []HostType{
		EOS,
		OpenWrt,
		Unclassified,
	}
	filter := &PlatformFilter{
		currentPlatform: Unclassified,
	}
	for _, plat := range platforms {
		if settings.FileExistsInFS(plat.indicatorFilename, fs) {
			filter.currentPlatform = plat
			return filter
		}
	}
	return filter
}

// IsRelevant implements the golang-shared plugins.PluginPredicate
// interface and only returns true when the current platform is supported.
func (pf *PlatformFilter) IsRelevant(pc plugins.PluginConstructor, metadata ...any) bool {
	for _, i := range metadata {
		if spec, ok := i.(PlatformSpec); ok {
			if len(spec.OnlyOn) > 0 && !slices.Contains(spec.OnlyOn, pf.currentPlatform) {
				return false
			} else if len(spec.Excludes) > 0 && slices.Contains(spec.Excludes, pf.currentPlatform) {
				return false
			}
		}
	}
	return true
}
