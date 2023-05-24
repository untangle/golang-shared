package plugins

type SettingsInjectablePlugin interface {
	Plugin
	GetNewSettings() any
	SetSettings(any)
}
