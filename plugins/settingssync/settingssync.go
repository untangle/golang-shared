package settingssync

import (
	"github.com/untangle/golang-shared/services/logger"
)

// A Consumer for handling the syncing of plugins/services within a daemon

type SettingsSyncer interface {
	// Method used to check if a plugin's settings are in sync by
	// passing in the latest settings of a plugin
	InSync(interface{}) bool

	// Gets the current settings for a plugin as a json
	// by reading from the settings file
	GetCurrentSettingsStruct() (interface{}, error)

	// Updates the plugins settings with the settings provided
	SyncSettings(interface{}) error
}

type SettingsSync struct {
	// List of plugins implementing the SettingsSyncer interface
	syncers []SettingsSyncer
}

// Returns a new instance of Settings Sync.
func NewSettingsSyncHandler() *SettingsSync {
	return &SettingsSync{}
}

// Registers a plugin as being managed by Settings Sync
func (settingsSync *SettingsSync) RegisterPlugin(plug SettingsSyncer) {
	settingsSync.syncers = append(settingsSync.syncers, plug)
}

// Syncs the settings of all plugins registered with Settings Sync
func (settingsSync *SettingsSync) SyncSettings() {
	logger.Info("Syncing Plugin Settings\n")

	for _, syncer := range settingsSync.syncers {

		updatedSettings, err := syncer.GetCurrentSettingsStruct()
		if err != nil {
			logger.Err("An error occurred and could not sync settings: %s", err.Error())
			continue
		}

		if !syncer.InSync(updatedSettings) {

			if err := syncer.SyncSettings(updatedSettings); err != nil {
				logger.Err("SettingsSync: %s", err.Error())
				continue
			}
		}
	}
}
