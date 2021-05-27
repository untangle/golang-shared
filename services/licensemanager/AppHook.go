package licensemanager

// AppHook struct is used by the product to populate the ValidApps array, which is a kvp of licensed apps
type AppHook struct {
	Start   func()
	Stop    func()
	Enabled func() bool
}
