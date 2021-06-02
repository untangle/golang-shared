package licensemanager

// ServiceCommand is used for setting the service state
type ServiceCommand struct {
	Name     string `json:"name"`
	NewState State  `json:"command"`
}
