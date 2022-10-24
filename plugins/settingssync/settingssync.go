package settingsync

import (
	"github.com/untangle/golang-shared/services/logger"
)

type SettingsSyncer interface {
	InSync(interface{}) bool

	GetCurrentSettingsStruct() (interface{}, error)

	SyncSettings(interface{}) error
}

type SettingsSync struct {
	syncers []SettingsSyncer
}

// Returns a new instance of Settings Sync.
func NewSettingsSyncHandler() *SettingsSync {
	return &SettingsSync{}
}

// Register's a plugin as being managed by Settings Sync
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
