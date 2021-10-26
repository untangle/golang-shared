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
		// Gen a new state file
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
		service := Service{Hook: o, State: serviceState}
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
	jsonLicense, err := ioutil.ReadFile(config.LicenseLocation)

	// Unmarshal
	err = json.Unmarshal(jsonLicense, &retLicense)
	if err != nil {
		logger.Warn("Error unmarshalling licenseInfo: %s\n", err.Error())
		return retLicense, err
	}

	// Return
	return retLicense, nil
}

// SetServiceState TODO
func SetServiceState(serviceName string, newAllowedState string, saveStates bool) error {
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

// findServiceHook is used to check if service is valid and return its hooks
// @param serviceName string - the name of the service
func findServiceHook(serviceName string) (ServiceHook, error) {
	service, ok := config.ValidServiceHooks[serviceName]
	if !ok {
		return ServiceHook{}, errServiceNotFound
	}
	return service, nil
}

// TODO
func findService(serviceName string) (*Service, error) {
	service, ok := services[serviceName]
	if !ok {
		return nil, errServiceNotFound
	}
	return service, nil
}

// TODO
func findServiceState(serviceName string, serviceStates []ServiceState) (ServiceState, bool) {
	for _, o := range serviceStates {
		if o.Name == serviceName {
			return o, true
		}
	}
	return ServiceState{}, false
}

// TODO
func saveServiceStatesFromServices(fileLocation string, services map[string]*Service) error {
	serviceStates := make([]ServiceState, 0)
	for _, o := range services {
		serviceStates = append(serviceStates, o.State)
	}
	return saveServiceStates(fileLocation, serviceStates)
}

// saveServiceStates stores the services in the service state file
// @param string fileLocation - the location of the service state file
// @param map[string]ServiceHook serviceHooks - map of service hooks to create ServiceStates for
// @param runInterrupt TODO
// @return []ServiceState - the service state array
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

func runSighup() error {
	logger.Info("Running interrupt\n")
	// write out service commands

	// TODO make generic here
	pidStr, err := exec.Command("pgrep", "packetd").CombinedOutput()
	if err != nil {
		logger.Err("Failure to get packetd pid: %s\n", err.Error())
		return err
	}
	logger.Info("Pid: %s\n", pidStr)

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidStr)))
	if err != nil {
		logger.Err("Failure converting pid for packetd: %s\n", err.Error())
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Err("Failure to get packetd process: %s\n", err.Error())
		return err
	}
	return process.Signal(syscall.SIGHUP)
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
