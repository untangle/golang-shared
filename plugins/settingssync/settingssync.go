package settingsync

import (
	"github.com/untangle/golang-shared/services/logger"
)

type SettingsSyncer interface {
	InSync(interface{}) bool

	GetSettingsStruct() (interface{}, error)

	SyncSettings(interface{}) error
}

type SettingsSync struct {
	syncers []SettingsSyncer
}

func NewSettingsSyncHandler() *SettingsSync {
	return &SettingsSync{}
}

func (settingsSync *SettingsSync) RegisterPlugin(plug SettingsSyncer) {
	// Have to strip the type off of plugin to check if it implements an interface
	// Ignore the golang linter complaining about this, it has to happen
	settingsSync.syncers = append(settingsSync.syncers, plug)
}

func (settingsSync *SettingsSync) SyncSettings() {
	logger.Info("Syncing Plugin Settings\n")

	for _, syncer := range settingsSync.syncers {

		updatedSettings, err := syncer.GetSettingsStruct()
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
