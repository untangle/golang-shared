package licensemanager

import "time"

// LicenseManagerConfig contains config options used in the license manager
type LicenseManagerConfig struct {
	// ValidServiceHooks is a map of apps and startup/shutdown/enabled hooks
	ValidServiceHooks map[string]ServiceHook

	// LicenseLocation is the location of the license file
	LicenseLocation string

	// ServiceStateLocation is the location of the service state file
	ServiceStateLocation string

	// WatchDogInterval is the watch dog timer interval
	WatchDogInterval time.Duration
}
