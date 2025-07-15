package types

// HostType is a externally-opaque type for declaring a host
// type.
type Platform struct {
	IndicatorFilename string
	Name              string

	// Path settings files should be looked up in
	SettingsDirPath string

	// Files that have a unique mapping on the platform.
	// The mapping goes from the openWRT equivalent file
	// to it's platform specific file location.
	AdditionalFileMappings map[string]string
}

// Equals checks if two platforms are equal. Simply
// compares their names
func (p *Platform) Equals(other Platform) bool {
	return p.Name == other.Name
}

var (
	EOS = Platform{
		IndicatorFilename: "/etc/Eos-release",
		Name:              "Eos",
		SettingsDirPath:   "/mnt/flash/mfw-settings",
		AdditionalFileMappings: map[string]string{
			"/etc/config/categories.json": "/usr/share/bctid/categories.json",
		},
	}
	OpenWrt = Platform{
		IndicatorFilename: "/etc/openwrt_version",
		Name:              "OpenWrt",
		SettingsDirPath:   "/etc/config",
	}
	Vittoria = Platform{
		IndicatorFilename: "/notknown/",
		Name:              "Unclassified",
		SettingsDirPath:   "/velocloud/",
	}
	Unclassified = Platform{
		IndicatorFilename: "",
		Name:              "Unclassified",
		SettingsDirPath:   "",
	}

	Platforms = []Platform{
		EOS,
		OpenWrt,
		Vittoria,
		Unclassified,
	}
)
