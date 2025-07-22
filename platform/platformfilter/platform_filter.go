package platformfilter

import (
	"fmt"
	"slices"

	"github.com/untangle/golang-shared/platform"
	"github.com/untangle/golang-shared/plugins"
)

// PlatformFilter is a plugin predicate (declared in
// golang-shared/plugins) which filters out plugins that have
// PluginSpec metadata and don't apply to the current platform.
type PlatformFilter struct {
	currentPlatform platform.HostType
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
	OnlyOn   []platform.HostType
	Excludes []platform.HostType
}

// NewPlatformFilter creates a new platform filter, during
// construction we determine the platform from the filesystem.
func NewPlatformFilter(platform platform.HostType) *PlatformFilter {
	return &PlatformFilter{
		currentPlatform: platform,
	}
}

// IsRelevant implements the golang-shared plugins.PluginPredicate
// interface and only returns true when the current platform is supported.
func (pf *PlatformFilter) IsRelevant(pc plugins.PluginConstructor, metadata ...any) bool {
	for _, i := range metadata {
		if spec, ok := i.(PlatformSpec); ok {
			if len(spec.OnlyOn) > 0 && len(spec.Excludes) > 0 {
				panic(
					fmt.Sprintf(
						"PlatformFilter.IsRelevant: "+
							"both OnlyOn and Excludes are filled out, only one should be specified."+
							"(OnlyOn: %v, Excludes: %v)",
						spec.OnlyOn,
						spec.Excludes))

			}
			platformMatch := func(p platform.HostType) bool { return p.Equals(pf.currentPlatform) }
			if len(spec.OnlyOn) > 0 && !slices.ContainsFunc(spec.OnlyOn, platformMatch) {
				return false
			} else if len(spec.Excludes) > 0 && slices.ContainsFunc(spec.Excludes, platformMatch) {
				return false
			}
		}
	}
	return true
}
