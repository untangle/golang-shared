package licensemanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/logger"
)

var config LicenseManagerConfig
var appStates []AppState

var errAppNotFoundError error = errors.New("app_not_found")
var shutdownChannelLicense chan bool
var wg sync.WaitGroup
var watchDog *time.Timer

// Startup the license manager service.
// @param configOptions LicenseManagerConfig - a license manager config object used for configuring the service
func Startup(configOptions LicenseManagerConfig) {
	shutdownChannelLicense = make(chan bool)

	config = configOptions

	logger.Info("Starting the license service\n")

	appStates, err := loadAppState(config.AppStateLocation)
	if err != nil {
		logger.Warn("Unable to retrieve previous app state. %v\n", err)
	}

	if appStates == nil {
		// Gen a new state file
		appStates, err = saveAppState(config.AppStateLocation, config.ValidApps)
		if err != nil {
			logger.Warn("Unable to initialize app states file. %v\n", err)
			return
		}
	}

	// Set each app to its previous state.
	logger.Debug("appstate %+v\n", appStates)
	for _, o := range appStates {
		if _, err = findApp(o.Name); err != nil {
			logger.Debug("App %s not found. Err: %v", o.Name, err)
			continue
		}
		cmd := AppCommand{Name: o.Name}
		if o.IsEnabled {
			cmd.NewState = StateEnable
		} else {
			cmd.NewState = StateDisable
		}
		SetAppState(cmd, true)
	}

	// restart licenses
	err = RefreshLicenses()
	if err != nil {
		logger.Warn("Not able to restart CLS: %v\n", err)
	}

	// watchdog for if CLS is alive
	wg.Add(1)
	go func() {
		defer wg.Done()
		watchDog = time.NewTimer(config.WatchDogInterval)
		defer watchDog.Stop()
		for {
			select {
			case <-shutdownChannelLicense:
				logger.Info("Shutdown CLS watchdog\n")
				return
			case <-watchDog.C:
				// on watch dog seen, restart license server
				// shutdown license items if restart did not work
				logger.Warn("Watch seen\n")
				refreshErr := RefreshLicenses()
				if refreshErr != nil {
					logger.Warn("Couldn't restart CLS: %s\n", refreshErr)
					shutdownApps(config.LicenseLocation, config.ValidApps)
				} else {
					logger.Info("Restarted CLS from watchdog\n")
				}
				watchDog.Reset(config.WatchDogInterval)
			}
		}
	}()
}

// Shutdown is called when the service stops
func Shutdown() {
	logger.Info("Shutting down the license service\n")
	if shutdownChannelLicense != nil {
		close(shutdownChannelLicense)
		wg.Wait()
	}
	shutdownApps(config.LicenseLocation, config.ValidApps)
}

// GetLicenseDefaults gets the default validApps
// @return []string - string array of app keys for CLS to use
func GetLicenseDefaults() []string {
	logger.Debug("GetLicenseDefaults()\n")
	keys := make([]string, len(config.ValidApps))
	i := 0
	for k := range config.ValidApps {
		keys[i] = k
		i++
	}
	watchDog.Reset(config.WatchDogInterval)
	return keys
}

// ClsIsAlive resets the watchdog interval for license <> app synchronization
func ClsIsAlive() {
	watchDog.Reset(config.WatchDogInterval)
}

// GetAppStates gets the current Service States
// @return []AppState - array of current app states
func GetAppStates() []AppState {
	return appStates
}

// SetAppState sets the desired state of an app
// @param cmd AppCommand - the command to run on the app
// @param save bool - if we should store the app state or not
// @return error - associated errors
func SetAppState(cmd AppCommand, save bool) error {
	var err error
	var app AppHook
	logger.Debug("Setting state for app %s to %v\n", cmd.Name, cmd.NewState)
	if app, err = findApp(cmd.Name); err != nil {
		return errAppNotFoundError
	}

	switch cmd.NewState {
	case StateEnable:
		app.Start()
	case StateDisable:
		app.Stop()
	}
	if save {
		appStates, err = saveAppState(config.LicenseLocation, config.ValidApps)
	}
	return err
}

// RefreshLicences restart the client licence service
func RefreshLicenses() error {
	output, err := exec.Command("/etc/init.d/clientlic", "restart").CombinedOutput()
	if err != nil {
		logger.Warn("license fetch failed: %s\n", err.Error())
		return err
	}
	if strings.Contains(string(output), "Command failed") {
		logger.Warn("license fetch failed: %s\n", string(output))
		err = errors.New(string(output))
		return err
	}
	return nil
}

// IsEnabled is called from API to see if app is currently enabled.
// @param appName string - the name of the app to check Enabled status of
func IsEnabled(appName string) (bool, error) {
	var app AppHook
	var err error
	if app, err = findApp(appName); err != nil {
		return false, errAppNotFoundError
	}
	return app.Enabled(), nil
}

// shutdownApps iterates appsToShutdown and calls the shutdown hook on them, and also removes the license file
// @param licenseFile string - the license file location
// @param appsToShutdown map[string]AppHook - the apps we want to shutdown
func shutdownApps(licenseFile string, appsToShutdown map[string]AppHook) {
	err := os.Remove(licenseFile)
	err = ioutil.WriteFile(licenseFile, []byte("{\"list\": []}"), 0444)
	if err != nil {
		logger.Warn("Failure to write non-license file: %v\n", err)
	}
	for name, _ := range appsToShutdown {
		cmd := AppCommand{Name: name, NewState: StateDisable}
		SetAppState(cmd, false)
	}
}

// findApp is used to check if app is valid and return its hooks
// @param appName string - the name of the app
func findApp(appName string) (AppHook, error) {
	app, ok := config.ValidApps[appName]
	if !ok {
		return AppHook{}, errAppNotFoundError
	}
	return app, nil
}

// saveAppState stores the apps in the appstate file
// @param string fileLocation - the location of the appstate file
// @param map[string]AppHook appsToStore - kvp of apps to store
// @return []appState - the app state array
// @return error - associated errors
func saveAppState(fileLocation string, appsToStore map[string]AppHook) ([]AppState, error) {
	var retAppStates = make([]AppState, 0)
	for name, o := range appsToStore {
		retAppStates = append(retAppStates, AppState{Name: name, IsEnabled: o.Enabled()})
	}
	data, err := json.Marshal(retAppStates)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(fileLocation, data, 0644)
	if err != nil {
		return nil, err
	}
	return retAppStates, nil

}

// loadAppState retrieves the previously saved app state
// @param fileLocation - the location of the app states file
// @return []appState - an array of app states, loaded from the file
// @return error - associated errors
func loadAppState(fileLocation string) ([]AppState, error) {
	var retAppStates = make([]AppState, 0)
	appStateContent, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		logger.Warn("Not able to find appstate file.\n", err)
		return nil, nil
	}

	err = json.Unmarshal(appStateContent, &retAppStates)
	if err != nil {
		logger.Warn("Not able to read content of app state file.%v \n ", err)
		return nil, err
	}
	return retAppStates, nil
}
