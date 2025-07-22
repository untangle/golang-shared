package platform

// HostType is a externally-opaque type for declaring a host
// type.
type HostType struct {
	IndicatorFilename   string
	Name                string
	SettingsDirPath     string
	UniquelyMappedFiles map[string]string
}

// Equals checks if two platforms are equal
// by simply comparing their names.
func (h *HostType) Equals(o HostType) bool {
	return h.Name == o.Name
}

var (
	EOS = HostType{
		IndicatorFilename: "etc/Eos-release",
		Name:              "Eos",
		SettingsDirPath:   "/mnt/flash/mfw-settings",
		UniquelyMappedFiles: map[string]string{
			"/etc/config/categories.json": "/usr/share/bctid/categories.json",
		},
	}
	OpenWrt = HostType{
		IndicatorFilename:   "etc/openwrt_version",
		Name:                "OpenWrt",
		SettingsDirPath:     "/etc/config",
		UniquelyMappedFiles: make(map[string]string),
	}
	Vittoria = HostType{
		// TODO: Update IndicatorFilename once the version is known
		IndicatorFilename:   "velocloud_version",
		Name:                "Vittoria",
		SettingsDirPath:     "/velocloud",
		UniquelyMappedFiles: make(map[string]string),
	}
	Unclassified = HostType{
		IndicatorFilename:   "",
		Name:                "Unclassified",
		UniquelyMappedFiles: make(map[string]string),
	}
	platforms = []HostType{
		EOS,
		OpenWrt,
		Vittoria,
		Unclassified,
	}
)
