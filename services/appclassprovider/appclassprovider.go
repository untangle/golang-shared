package appclassprovider

import (
	"errors"
	"io/fs"

	"github.com/untangle/golang-shared/services/appclassmanager"
	"github.com/untangle/golang-shared/services/dpi"
	"github.com/untangle/golang-shared/util"
)

// ApplicationClassProvider is the interface for providing application class information
// There are currently two providers, dpi and appclassmanager
type ApplicationClassProvider interface {
	Startup() error
	Shutdown() error
	GetTable(table string) (string, error)
	Name() string
}

// Generic SetProvider function.
func GetApplicationClassProvider(fs fs.FS) (ApplicationClassProvider, error) {
	var provider ApplicationClassProvider
	var err error
	platform := util.GetPlatform()
	switch platform {
	case util.EOS:
		provider = dpi.NewDpiConfigManager(fs)
		err = provider.Startup()
	case util.OpenWRT:
		provider = appclassmanager.NewAppClassManager(fs)
		err = provider.Startup()
	default:
		err = errors.New("unknown_platform")
	}
	return provider, err
}
