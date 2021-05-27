package licensemanager

// Used to load and save AppState. Need to set each app to its previous state when starting up.
type AppState struct {
	Name      string `json: "appname`
	IsEnabled bool   `json: "enabled"`
}
