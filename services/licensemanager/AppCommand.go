package licensemanager

// Command is used for specific commands (SetState)
type AppCommand struct {
	Name     string `json:"name"`
	NewState State  `json:"command"`
}
