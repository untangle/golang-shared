package licensemanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"

	loggerModel "github.com/untangle/golang-shared/logger"
	"github.com/untangle/golang-shared/plugins/util"
)

const (
	// LicenseFileDoesNotExistStr is the string to check if licenses should be reloaded when status is returned
	LicenseFileDoesNotExistStr string = "RELOAD_LICENSES"

	pluginName string = "licensemanager"

	clientLicenseService string = "bin/client-license-service"
)

var (
	errServiceNotFound error = errors.New("service_not_found")
)

type LicenseManager struct {
	config   *Config
	watchDog *time.Timer
	services map[string]*Service
	logger   loggerModel.LoggerLevels

	ctx       context.Context
	cancelCtx context.CancelFunc

	// Function to use when refreshing the CLI
	// Can be swapped out for unit testing
	RefreshLicenses func() error
}

// Returns a new LicenseManager instance pointer. If the provided config
// is not valid, the function will return nil. Can't return an error since
// this is being used by the GlobalPluginManager
func NewLicenseManager(config *Config, logger loggerModel.LoggerLevels) *LicenseManager {
	if config == nil {
		logger.Err("Invalid config used when creating the License Manager\n")
		return nil
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	return &LicenseManager{
		config:   config,
		watchDog: time.NewTimer(config.WatchDogInterval),
		services: make(map[string]*Service),
		logger:   logger,

		ctx:       ctx,
		cancelCtx: cancelCtx,

		RefreshLicenses: RefreshLicenses,
	}
}

// Returns name of the service
func (lm *LicenseManager) Name() string {
	return pluginName
}

// Startup the license manager service.
func (lm *LicenseManager) Startup() error {
	lm.logger.Info("Starting the license service\n")

	serviceStates, err := LoadServiceStates(lm.config.ServiceStateLocation)
	if err != nil {
		lm.logger.Warn("Unable to retrieve previous service state. %v\n", err)
	}

	if serviceStates == nil {
		// Gen a new state file with services set to StateDisable
		blankServiceStates := make([]ServiceState, 0)
		for name := range lm.config.ValidServiceHooks {
			newServiceState := ServiceState{Name: name, AllowedState: StateDisable}
			blankServiceStates = append(blankServiceStates, newServiceState)
		}
		if err = saveServiceStates(lm.config.ServiceStateLocation, blankServiceStates); err != nil {
			return fmt.Errorf("unable to initialize service states file. %w", err)
		}
	}

	// Create each service
	lm.logger.Debug("States %+v\n", serviceStates)
	for name, o := range lm.config.ValidServiceHooks {
		var serviceState ServiceState
		var found bool
		serviceState, found = findServiceState(name, serviceStates)
		if !found {
			serviceState = ServiceState{Name: name, AllowedState: StateDisable}
		}
		service := Service{Name: name, Hook: o, State: serviceState}
		lm.services[name] = &service
	}

	// restart licenses
	if err = lm.RefreshLicenses(); err != nil {
		lm.logger.Warn("Not able to restart CLS: %v\n", err)
	}

	go lm.clsWatchdog()
	return nil
}

// Watchdog checking if CLS is alive
func (lm *LicenseManager) clsWatchdog() {
	defer lm.watchDog.Stop()
	for {
		select {
		case <-lm.ctx.Done():
			lm.logger.Info("Shutdown CLS watchdog\n")
			return
		case <-lm.watchDog.C:
			// on watch dog seen, restart license server
			// shutdown license items if restart did not work
			lm.logger.Warn("Watch seen\n")
			if refreshErr := lm.RefreshLicenses(); refreshErr != nil {
				lm.logger.Warn("Couldn't restart CLS: %s\n", refreshErr)
				lm.shutdownServices()
			} else {
				lm.logger.Info("Restarted CLS from watchdog\n")
			}
			lm.watchDog.Reset(lm.config.WatchDogInterval)
		}
	}
}

// Shutdown is called when the service stops
func (lm *LicenseManager) Shutdown() error {
	lm.logger.Info("Shutting down the license service\n")
	lm.cancelCtx()

	lm.shutdownServices()
	return nil
}

// GetLicenseDefaults gets the default validServiceStates
// @return []string - string array of service keys for CLS to use
func (lm *LicenseManager) GetLicenseDefaults() []string {
	lm.logger.Debug("GetLicenseDefaults()\n")
	keys := make([]string, len(lm.config.ValidServiceHooks))
	i := 0
	for k := range lm.config.ValidServiceHooks {
		keys[i] = k
		i++
	}
	lm.watchDog.Reset(lm.config.WatchDogInterval)
	return keys
}

// ClsIsAlive resets the watchdog interval for license <> service synchronization
func (lm *LicenseManager) ClsIsAlive() {
	lm.watchDog.Reset(lm.config.WatchDogInterval)
}

// GetServices gets the current Service
// @return []Service - array of current services
func (lm *LicenseManager) GetServices() map[string]*Service {
	return lm.services
}

// RefreshLicenses restart the client licence service
func RefreshLicenses() error {
	// do not bail when license refresh fails
	_ = util.RunSigusr1(clientLicenseService)
	return nil
}

// IsLicenseEnabled is called from API to see if service is currently enabled.
// @param serviceName string - the name of the service to check Enabled status of
func (lm *LicenseManager) IsLicenseEnabled(serviceName string) (bool, error) {
	var serv *Service
	var err error
	if serv, err = lm.findService(serviceName); err != nil {
		return false, errServiceNotFound
	}
	return serv.State.getAllowedState() == StateEnable, nil
}

// GetLicenseDetails will use the current license location to load and return the license file
// @return LicenseInfo - the license info, containing license details
// @return error - associated errors
func (lm *LicenseManager) GetLicenseDetails() (LicenseInfo, error) {

	var retLicense LicenseInfo

	// Load file
	licenseFileExists := licenseFileExists(lm.config.LicenseLocation)
	if !licenseFileExists {
		lm.logger.Warn("License file does not exist\n")
		return retLicense, errors.New(LicenseFileDoesNotExistStr)
	}

	jsonLicense, err := os.ReadFile(lm.config.LicenseLocation)
	if err != nil {
		lm.logger.Warn("Error opening license file: %s\n", err.Error())
		return retLicense, err
	}

	// Unmarshal
	err = json.Unmarshal(jsonLicense, &retLicense)
	if err != nil {
		lm.logger.Warn("Error unmarshalling licenseInfo: %s\n", err.Error())
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

// SetServices will disable any disabled services to un-enabled in
// settings, and the appstate file.
func (lm *LicenseManager) SetServices(enabledServices map[string]bool) error {
	var err error
	for serviceName, isEnabled := range enabledServices {
		if service, err := lm.findService(serviceName); err != nil {
			lm.logger.Warn("LicenseManager: when updating services, given nonexistent service: %s\n",
				serviceName)
		} else if isEnabled {
			if setServiceStateErr := service.setServiceState(StateEnable); setServiceStateErr != nil {
				lm.logger.Warn("Failed to set the desired state for service %v with error %v\n", serviceName, setServiceStateErr.Error())
			}
		} else {
			lm.disableService(service)
		}
	}

	err = saveServiceStatesFromServices(lm.config.ServiceStateLocation, lm.services)

	_, HybridSignalErr := os.Stat("/usr/bin/updateSysdbSignal")
	if HybridSignalErr == nil || !os.IsNotExist(HybridSignalErr) {
		if HybridSignalErr := exec.Command("/usr/bin/updateSysdbSignal", "--sighup").Run(); HybridSignalErr != nil {
			logger.Warn("Failed to run EOS-MFW script `updateSysdbSignal` command with error: %+v\n", HybridSignalErr)
		}
	}

	if RunSighupErr := util.RunSighup(lm.config.Executable); RunSighupErr != nil {
		lm.logger.Warn("Failed to run RunSighup on executable %v with an error %v\n", lm.config.Executable, RunSighupErr.Error())
	}

	return err
}

// disableService disables a service
func (lm *LicenseManager) disableService(service *Service) {
	if err := service.setServiceState(StateDisable); err != nil {
		lm.logger.Warn("Failed to set the desired state for service %v with error %v\n", service, err.Error())
	}
	if service.Hook.Disabled == nil {
		return
	}
	newSettings, settingsSegs, err := service.Hook.Disabled()
	if err != nil {
		lm.logger.Warn("Failed to get disabled settings for service %s\n", service.Name)
	}

	// Set settings
	if _, err = settings.SetSettings(settingsSegs, newSettings, true, false); err != nil {
		lm.logger.Warn("Failed to set disabled settings for service %s\n", service.Name)
	}
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
func (lm *LicenseManager) shutdownServices() {
	if err := os.Remove(lm.config.LicenseLocation); err != nil {
		lm.logger.Err("Could not remove the license file when shutting down services: %v\n", err)
	}

	if err := os.WriteFile(lm.config.LicenseLocation, []byte("{\"list\": []}"), 0444); err != nil {
		lm.logger.Warn("Failure to write non-license file: %v\n", err)
	}

	for _, service := range lm.services {
		if err := service.setServiceState(StateDisable); err != nil {
			lm.logger.Warn("Failed to set the desired state for service %v with error %v\n", service, err.Error())
		}
	}

	_, HybridSignalErr := os.Stat("/usr/bin/updateSysdbSignal")
	if HybridSignalErr == nil || !os.IsNotExist(HybridSignalErr) {
		if HybridSignalErr := exec.Command("/usr/bin/updateSysdbSignal", "--sighup").Run(); HybridSignalErr != nil {
			logger.Warn("Failed to run EOS-MFW script `updateSysdbSignal` command with error: %+v\n", HybridSignalErr)
		}
	}

	if RunSighupErr := util.RunSighup(lm.config.Executable); RunSighupErr != nil {
		lm.logger.Warn("Failed to run RunSighup given executable: %v\n", RunSighupErr.Error())
	}
}

// findService finds the service in the services map
// @param string serviceName - service to find
// @return *Service - the service found, nil if not found
// @return error such as errServiceNotFound
func (lm *LicenseManager) findService(serviceName string) (*Service, error) {
	service, ok := lm.services[serviceName]
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
func findServiceState(
	serviceName string,
	serviceStates []ServiceState) (ServiceState, bool) {
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
	if err = os.WriteFile(fileLocation, data, 0644); err != nil {
		logger.Warn("Failure to write state file: %s\n", err.Error())
		return err
	}

	return nil
}

// LoadServiceStates retrieves the previously saved service state
// @param fileLocation - the location of the service states file
// @return []ServiceState - an array of service states, loaded from the file
// @return error - associated errors
func LoadServiceStates(fileLocation string) ([]ServiceState, error) {
	var serviceStates = make([]ServiceState, 0)
	fileContent, err := os.ReadFile(fileLocation)
	if err != nil {
		logger.Warn("Not able to find service state file.\n", err)
		return nil, nil
	}

	if err = json.Unmarshal(fileContent, &serviceStates); err != nil {
		logger.Warn("Not able to read content of service state file.%v \n ", err)
		return nil, err
	}
	return serviceStates, nil
}
