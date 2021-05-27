package licensemanager

import "time"

// LicenseManagerConfig contains config options used in the license manager
type LicenseManagerConfig struct {
	// ValidApps is a map of apps and startup/shutdown/enabled hooks
	ValidApps map[string]AppHook

	// LicenseLocation is the location of the license file
	LicenseLocation string

	// AppStateLocation is the location of the app state file
	AppStateLocation string

	// WatchDogInterval is the watch dog timer interval
	WatchDogInterval time.Duration
}
