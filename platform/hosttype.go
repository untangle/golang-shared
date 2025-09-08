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
	Vittoria = HostType{
		IndicatorFilename:   "velocloud",
		Name:                "Vittoria",
		SettingsDirPath:     "/opt/mfw/etc",
		UniquelyMappedFiles: make(map[string]string),
	}
	OpenWrt = HostType{
		IndicatorFilename:   "etc/openwrt_version",
		Name:                "OpenWrt",
		SettingsDirPath:     "/etc/config",
		UniquelyMappedFiles: make(map[string]string),
	}
	EOS = HostType{
		IndicatorFilename: "etc/Eos-release",
		Name:              "Eos",
		SettingsDirPath:   "/mnt/flash/mfw-settings",
		UniquelyMappedFiles: map[string]string{
			"/etc/config/categories.json": "/usr/share/bctid/categories.json",
		},
	}
	Unclassified = HostType{
		IndicatorFilename:   "",
		Name:                "Unclassified",
		UniquelyMappedFiles: make(map[string]string),
	}
	platforms = []HostType{
		Vittoria,
		OpenWrt,
		EOS,
		Unclassified,
	}
)
