package plugins

// SettingsInjectablePlugin is a plugin that supports 'settings
// injection'.
//
// Settings Injection is super simple -- return a pointer to a
// settings object you want 'filled out' (e.g. there is for example
// some JSON somewhere, and we want to load it into our config object
// that has the proper tags).  Then later, the object in charge of
// actually finding and loading settings calls SetSettings() with that
// same object as a notification that the new settings are ready.
type SettingsInjectablePlugin interface {
	Plugin
	GetNewSettings() any
	SetSettings(any)
}
