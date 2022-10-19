package settingsync

import (
	"syscall"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	pluginName string = "SettingsSync"
)

type SettingsSyncer interface {
	InSync(interface{}) bool

	GetSettingsPath() string

	GetSettingsStruct() interface{}

	SyncSettings(interface{}) error
}

type SettingsSync struct {
	syncers []SettingsSyncer
}

func NewSettingsSync() *SettingsSync {
	return &SettingsSync{syncers: make([]SettingsSyncer, 0)}
}

func (settingsSync *SettingsSync) RegisterManager(manager SettingsSyncer) {
	settingsSync.syncers = append(settingsSync.syncers, manager)
}

func (settingsSync *SettingsSync) Signal(signal syscall.Signal) error {
	switch signal {
	case syscall.SIGHUP:
		for _, syncer := range settingsSync.syncers {
			settingsPath := syncer.GetSettingsPath()

			var updatedSettings interface{}
			if err := settings.UnmarshalSettingsAtPath(updatedSettings, settingsPath); err != nil {
				logger.Err("SettingsSync: %s", err.Error())
				continue
			}

			if !syncer.InSync(updatedSettings) {

				if err := syncer.SyncSettings(settingsPath); err != nil {
					logger.Err("SettingsSync: %s", err.Error())
					continue
				}
			}
		}
	}

	return nil
}

func (settingsSync *SettingsSync) Startup() error {
	return nil
}

func (settingsSync *SettingsSync) Shutdown() error {
	return nil
}

func (settingsSync *SettingsSync) Name() string {
	return pluginName
}
