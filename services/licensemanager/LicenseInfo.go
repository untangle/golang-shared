package licensemanager

// LicenseInfo represents the json returned from license server
type LicenseInfo struct {
	JavaClass string    `json:"javaClass"`
	List      []License `json:"list"`
}
