package licensemanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	// LicenseFileDoesNotExistStr is the string to check if licenses should be reloaded when status is returned
	LicenseFileDoesNotExistStr string = "RELOAD_LICENSES"
)

var config Config
var services map[string]*Service

var errServiceNotFound error = errors.New("service_not_found")
var shutdownChannelLicense chan bool
var wg sync.WaitGroup
var watchDog *time.Timer

// Startup the license manager service.
// @param configOptions LicenseManagerConfig - a license manager config object used for configuring the service
func Startup(configOptions Config) {
	shutdownChannelLicense = make(chan bool)
	services = make(map[string]*Service)

	config = configOptions

	logger.Info("Starting the license service\n")

	serviceStates, err := loadServiceStates(config.ServiceStateLocation)
	if err != nil {
		logger.Warn("Unable to retrieve previous service state. %v\n", err)
	}

	if serviceStates == nil {
		// Gen a new state file with services set to StateDisable
		blankServiceStates := make([]ServiceState, 0)
		for name := range config.ValidServiceHooks {
			newServiceState := ServiceState{Name: name, AllowedState: StateDisable}
			blankServiceStates = append(blankServiceStates, newServiceState)
		}
		err = saveServiceStates(config.ServiceStateLocation, blankServiceStates)
		if err != nil {
			logger.Warn("Unable to initialize service states file. %v\n", err)
			return
		}
	}

	// Create each service
	logger.Debug("States %+v\n", serviceStates)
	for name, o := range config.ValidServiceHooks {
		var serviceState ServiceState
		var found bool
		serviceState, found = findServiceState(name, serviceStates)
		if !found {
			serviceState = ServiceState{Name: name, AllowedState: StateDisable}
		}
		service := Service{Name: name, Hook: o, State: serviceState}
		services[name] = &service
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
					shutdownServices(config.LicenseLocation, services)
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
	shutdownServices(config.LicenseLocation, services)
}

// GetLicenseDefaults gets the default validServiceStates
// @return []string - string array of service keys for CLS to use
func GetLicenseDefaults() []string {
	logger.Debug("GetLicenseDefaults()\n")
	keys := make([]string, len(config.ValidServiceHooks))
	i := 0
	for k := range config.ValidServiceHooks {
		keys[i] = k
		i++
	}
	watchDog.Reset(config.WatchDogInterval)
	return keys
}

// ClsIsAlive resets the watchdog interval for license <> service synchronization
func ClsIsAlive() {
	watchDog.Reset(config.WatchDogInterval)
}

// GetServices gets the current Service
// @return []Service - array of current services
func GetServices() map[string]*Service {
	return services
}

// RefreshLicenses restart the client licence service
func RefreshLicenses() error {
	// pkill is not installed on MFW.
	output, err := exec.Command("pgrep", "client-license-service").CombinedOutput()
	if err != nil {
		spid := strings.TrimSuffix(string(output), "\n")
		npid, err := strconv.Atoi(spid)
		if err != nil {
			logger.Warn("Not able to get pid of CLS: %v\n", err)
			return err
		}
		syscall.Kill(npid, syscall.SIGUSR1)
		return nil
	}
	logger.Warn("Not able to refresh CLS: %v\n", err)
	return err
}

// IsLicenseEnabled is called from API to see if service is currently enabled.
// @param serviceName string - the name of the service to check Enabled status of
func IsLicenseEnabled(serviceName string) (bool, error) {
	var serv *Service
	var err error
	if serv, err = findService(serviceName); err != nil {
		return false, errServiceNotFound
	}
	return serv.State.getAllowedState() == StateEnable, nil
}

// GetLicenseDetails will use the current license location to load and return the license file
// @return LicenseInfo - the license info, containing license details
// @return error - associated errors
func GetLicenseDetails() (LicenseInfo, error) {

	var retLicense LicenseInfo

	// Load file
	licenseFileExists := licenseFileExists(config.LicenseLocation)
	if !licenseFileExists {
		logger.Warn("License file does not exist\n")
		return retLicense, errors.New(LicenseFileDoesNotExistStr)
	}

	jsonLicense, err := ioutil.ReadFile(config.LicenseLocation)
	if err != nil {
		logger.Warn("Error opening license file: %s\n", err.Error())
		return retLicense, err
	}

	// Unmarshal
	err = json.Unmarshal(jsonLicense, &retLicense)
	if err != nil {
		logger.Warn("Error unmarshalling licenseInfo: %s\n", err.Error())
		return retLicense, err
	}

	// Return
	return retLicense, nil
}

// GetLicenseFileDoesNotExistStr returns the error string for license file does not exist for comparison reasons
// @return string of the license file does not exist error
func GetLicenseFileDoesNotExistStr() string {
	return LicenseFileDoesNotExistStr
}

// SetServices will disable any disabled services to un-enabled in settings
func SetServices(enabledServices map[string]bool) error {
	var err error = nil
	for serviceName, valid := range enabledServices {
		if !valid {
			// find service, get disabled hook, and run it
			var service *Service
			logger.Debug("Set %s to invalid\n", serviceName)
			service, err = findService(serviceName)
			if err != nil {
				logger.Warn("Failed to set un-enabled for service %s\n", serviceName)
				continue
			}

			// get the new disabled settings, the segments to set, and any errors
			newSettings, settingsSegments, disableErr := service.Hook.Disabled()
			if disableErr != nil {
				logger.Warn("Failed to get disabled settings for service %s\n", serviceName)
				err = disableErr
				continue
			}

			// Set settings
			_, err = settings.SetSettings(settingsSegments, newSettings, true)
			if err != nil {
				logger.Warn("Failed to set disabled settings for service %s\n", serviceName)
			}
		}
	}

	for serviceName, valid := range enabledServices {
		cmd := "disable"
		if valid {
			cmd = "enable"
		}
		err = setServiceState(serviceName, cmd, true)

		if err != nil {
			logger.Warn("Failed to set service: %s: %s\n", serviceName, err.Error())
			continue
		}
	}

	return err
}

// setServiceState sets the given serviceName to the given allowedState
// @param string serviceName - service to set
// @param string newAllowedState - new allowed state such as enabled or disabled
// @param bool saveStates - whether ServiceState file should be saved
// @return any error
func setServiceState(serviceName string, newAllowedState string, saveStates bool) error {
	service, err := findService(serviceName)
	if err != nil {
		logger.Warn("Failure to find service: %s\n", err.Error())
		return err
	}

	var newState State
	err = newState.FromString(newAllowedState)
	if err != nil {
		logger.Warn("Failure getting newAllowedState: %s\n", err.Error())
		return err
	}

	err = service.setServiceState(newState)
	if err != nil {
		logger.Warn("Failure setting service state: %s\n", err.Error())
		return err
	}

	if saveStates {
		err = saveServiceStatesFromServices(config.ServiceStateLocation, services)
	}

	return nil

}

// licenseFileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func licenseFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// shutdownServices iterates servicesToShutdown and calls the shutdown hook on them, and also removes the license file
// @param licenseFile string - the license file location
// @param servicesToShutdown map[string]Service - the services we want to shutdown
func shutdownServices(licenseFile string, servicesToShutdown map[string]*Service) {
	err := os.Remove(licenseFile)
	err = ioutil.WriteFile(licenseFile, []byte("{\"list\": []}"), 0444)
	if err != nil {
		logger.Warn("Failure to write non-license file: %v\n", err)
	}
	for _, service := range servicesToShutdown {
		service.setServiceState(StateDisable)
	}
}

// findService finds the service in the services map
// @param string serviceName - service to find
// @return *Service - the service found, nil if not found
// @return error such as errServiceNotFound
func findService(serviceName string) (*Service, error) {
	service, ok := services[serviceName]
	if !ok {
		return nil, errServiceNotFound
	}
	return service, nil
}

// findServiceState finds a given service in the passed ServiceState array
// @param string serviceName - service to find
// @param []ServiceState serviceStates - array of ServiceState to search through
// @return ServiceState of the state found, blank ServiceState if not found
// @return bool on if found
func findServiceState(serviceName string, serviceStates []ServiceState) (ServiceState, bool) {
	for _, o := range serviceStates {
		if o.Name == serviceName {
			return o, true
		}
	}
	return ServiceState{}, false
}

// saveServiceStatesFromServices saves the services states in an array of services
// @param string fileLocation - location to save service states to
// @param map[string]*Service services - services to save
// @return any error from saving, returned from saveServiceStates
func saveServiceStatesFromServices(fileLocation string, services map[string]*Service) error {
	serviceStates := make([]ServiceState, 0)
	for _, o := range services {
		serviceStates = append(serviceStates, o.State)
	}
	return saveServiceStates(fileLocation, serviceStates)
}

// saveServiceStates stores the services in the service state file
// @param string fileLocation - the location of the service state file
// @param []ServiceState serviceStates - map of service hooks to create ServiceStates for
// @return error - associated errors
func saveServiceStates(fileLocation string, serviceStates []ServiceState) error {
	data, err := json.Marshal(serviceStates)
	if err != nil {
		logger.Warn("Failure to marshal states: %s\n", err.Error())
		return err
	}
	err = ioutil.WriteFile(fileLocation, data, 0644)
	if err != nil {
		logger.Warn("Failure to write state file: %s\n", err.Error())
		return err
	}

	return nil

}

// loadServiceStates retrieves the previously saved service state
// @param fileLocation - the location of the service states file
// @return []ServiceState - an array of service states, loaded from the file
// @return error - associated errors
func loadServiceStates(fileLocation string) ([]ServiceState, error) {
	var serviceStates = make([]ServiceState, 0)
	fileContent, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		logger.Warn("Not able to find service state file.\n", err)
		return nil, nil
	}

	err = json.Unmarshal(fileContent, &serviceStates)
	if err != nil {
		logger.Warn("Not able to read content of service state file.%v \n ", err)
		return nil, err
	}
	return serviceStates, nil
}
