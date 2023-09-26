package licensemanager

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/untangle/golang-shared/testing/mocks"
)

var (
	// When a licensemanager instance shutdowns, it deletes
	// the licenses.json file. Point licensemanager instances
	// to the same empty file, which they can delete whenever.
	// Use licenseWithData when testing functionality around
	// the actual file
	licenseLocation = "./testdata/licenses/dummy_licenses.json"
	licenseWithData = "./testdata/licenses/licenses.json"

	// Service states files starting at an expected state
	servicesStatesMixedFile    = "./testdata/appstates/appstate_mixed.json"
	servicesStatesEnabledFile  = "./testdata/appstates/appstate_enabled.json"
	servicesStatesDisabledFile = "./testdata/appstates/appstate_disabled.json"
	servicesStatesInvalidFile  = "./testdata/appstates/appstate_invalid.json"

	// Service states file to test updates on. Its state is
	// unknown when it's read in. Use the mutex when testing with it
	// to prevent concurrent reads/writes from tests
	serviceStatesTestUpdates     = "./testdata/appstates/test_updates.json"
	serviceStatesTestUpdatesLock sync.RWMutex

	// Some tests will be altering files in the allowedstate directory
	// Make sure they don't interfere with each other with a lock.
	allowedStateDirLock sync.RWMutex
)

func TestName(t *testing.T) {
	lm := &LicenseManager{}

	assert.Equal(t, pluginName, lm.Name())
}

func TestFindServiceState(t *testing.T) {
	serviceStates, _ := LoadServiceStates(servicesStatesEnabledFile)
	actualServiceState, found := findServiceState("untangle-node-discovery", serviceStates)
	assert.True(t, found)
	assert.NotNil(t, actualServiceState)

}

// Test that on a time out, the watchdog refreshes the CLI
func TestClsWatchdog(t *testing.T) {
	lm := NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})
	lm.config.WatchDogInterval = time.Second * 1

	watchdogAlerted := false
	lm.RefreshLicenses = func() error {
		watchdogAlerted = true
		return nil
	}

	if startupErr := lm.Startup(); startupErr != nil {
		lm.logger.Warn("Failed to to start Licence manager service : %v \n", startupErr.Error())
	}

	defer func() {
		_ = lm.Shutdown()
	}()

	go lm.clsWatchdog()
	time.Sleep(time.Second * 2)

	assert.True(t, watchdogAlerted)
}

func TestLoadServiceStates(t *testing.T) {
	// Test reading in
	actualEnabled, err := LoadServiceStates(servicesStatesEnabledFile)
	assert.NoError(t, err)
	if len(actualEnabled) < 1 {
		t.Fatalf("No services read from %s could be used to conduct the test.\n", servicesStatesEnabledFile)
	}

	for _, serviceState := range actualEnabled {
		assert.Equal(t, StateEnable, serviceState.getAllowedState())
		assert.NotEmpty(t, serviceState.Name)
	}

	actualDisabled, err := LoadServiceStates(servicesStatesDisabledFile)
	assert.NoError(t, err)
	if len(actualDisabled) < 1 {
		t.Fatalf("No services read from %s could be used to conduct the test.\n", servicesStatesDisabledFile)
	}

	for _, serviceState := range actualDisabled {
		assert.Equal(t, StateDisable, serviceState.getAllowedState())
		assert.NotEmpty(t, serviceState.Name)
	}

	// Test a bad path
	// Function isn't considered to be in an error state if the appstate file can't be found
	_, err = LoadServiceStates("bogusPath")
	assert.NoError(t, err)

	// Test a valid path but invalid file
	_, err = LoadServiceStates(servicesStatesInvalidFile)
	assert.Error(t, err)
}

func TestGetLicenseDefaults(t *testing.T) {
	lm := NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})
	serviceKeys := lm.GetLicenseDefaults()

	expectedKeys := []string{
		"untangle-node-discovery",
		"untangle-node-classd",
		"untangle-node-threat-prevention",
		"untangle-node-sitefilter",
		"untangle-node-geoip",
		"untangle-node-captiveportal",
	}

	assert.ElementsMatch(t, expectedKeys, serviceKeys)
}

func TestGetLicenseDetails(t *testing.T) {
	lm := &LicenseManager{
		config: &Config{
			LicenseLocation: licenseWithData,
		},
		logger: mocks.NewMockLogger(),
	}
	licenses, err := lm.GetLicenseDetails()

	assert.NoError(t, err)
	assert.Equal(t, "java.util.LinkedList", licenses.JavaClass)

	// Spot check values in the large data structure returned
	expectedDisplayNames := []string{
		"Throughput",
		"Threat Prevention",
		"Web Filter",
		"GeoIP Fencing",
		"Application Control",
		"Database Services",
		"Device Discovery",
		"Captive Portal",
	}

	actualDisplayNames := []string{}
	for _, l := range licenses.List {
		actualDisplayNames = append(actualDisplayNames, l.DisplayName)
	}

	assert.ElementsMatch(t, expectedDisplayNames, actualDisplayNames)

	// Test that an error is returned when the license file cannot be read
	lm.config.LicenseLocation = "badPath"
	_, err = lm.GetLicenseDetails()
	assert.Error(t, err)
}

func TestGetLicenseFileDoesNotExistStr(t *testing.T) {
	assert.Equal(t, LicenseFileDoesNotExistStr, GetLicenseFileDoesNotExistStr())
}

func TestLicenseFileExists(t *testing.T) {
	assert.True(t, licenseFileExists(servicesStatesEnabledFile))

	// Give it a directory
	assert.False(t, licenseFileExists("./"))

	assert.False(t, licenseFileExists("./doesNotExists.txt"))
}

func TestShutdownServices(t *testing.T) {
	lm := NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})

	if startupErr := lm.Startup(); startupErr != nil {
		lm.logger.Warn("Failed to to start Licence manager service : %v \n", startupErr.Error())
	}

	defer func() {
		_ = lm.Shutdown()
	}()

	lm.shutdownServices()

	for _, v := range lm.services {
		assert.Equal(t, StateDisable, v.State.AllowedState)
	}
}

func TestSaveServiceStatesFrom(t *testing.T) {
	lm := NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})
	lm.config.ServiceStateLocation = serviceStatesTestUpdates

	serviceStatesTestUpdatesLock.Lock()
	defer serviceStatesTestUpdatesLock.Unlock()

	if startupErr := lm.Startup(); startupErr != nil {
		lm.logger.Warn("Failed to to start Licence manager service : %v \n", startupErr.Error())
	}

	defer func() {
		_ = lm.Shutdown()
	}()

	// Setup the test by setting the service to a known state
	_ = lm.services["untangle-node-discovery"].setServiceState(StateEnable)

	lm.services["untangle-node-discovery"].State.AllowedState = StateDisable

	_ = saveServiceStatesFromServices(serviceStatesTestUpdates, lm.services)

	// Read back in the states
	statesFromFile, err := LoadServiceStates(serviceStatesTestUpdates)
	assert.NoError(t, err)
	var actualServiceState *ServiceState

	// Get state of the service from the file
	for _, v := range statesFromFile {
		if v.Name == "untangle-node-discovery" {
			actualServiceState = &v
			break
		}
		continue
	}

	// Fail the test if the service couldn't be found. Nothing to
	if !assert.NotNil(t, actualServiceState) {
		t.Fatalf("The actual service could not be found in the state file %v\n", actualServiceState)
	}

	assert.Equal(t, actualServiceState.getAllowedState(), actualServiceState.getAllowedState())
}

func TestSaveServiceStates(t *testing.T) {
	lm := NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})
	lm.config.ServiceStateLocation = serviceStatesTestUpdates

	serviceStatesTestUpdatesLock.Lock()
	defer serviceStatesTestUpdatesLock.Unlock()

	if startupErr := lm.Startup(); startupErr != nil {
		lm.logger.Warn("Failed to to start Licence manager service : %v \n", startupErr.Error())
	}

	defer func() {
		_ = lm.Shutdown()
	}()

	// Setup the test by setting the service to a known state
	_ = lm.services["untangle-node-discovery"].setServiceState(StateEnable)
	lm.services["untangle-node-discovery"].State.AllowedState = StateDisable

	serviceStates := []ServiceState{}
	for _, v := range lm.services {
		serviceStates = append(serviceStates, v.State)
	}

	_ = saveServiceStates(serviceStatesTestUpdates, serviceStates)

	// Read back in the states
	statesFromFile, err := LoadServiceStates(serviceStatesTestUpdates)
	assert.NoError(t, err)
	var actualServiceState *ServiceState

	// Get state of the service from the file
	for _, v := range statesFromFile {
		if v.Name == "untangle-node-discovery" {
			actualServiceState = &v
			break
		}
		continue
	}

	// Fail the test if the service couldn't be found. Nothing to
	if !assert.NotNil(t, actualServiceState) {
		t.Fatalf("The actual service could not be found in the state file %v\n", actualServiceState)
	}

	assert.Equal(t, actualServiceState.getAllowedState(), actualServiceState.getAllowedState())
}

// Requires the service's Startup func to be run, so it would normally be
// added to the test suite. However, it needs to be able to write
// to the licenses service state file without altering expected state for other tests.
func TestSetServiceStateLicenseManager(t *testing.T) {
	lm := NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})
	lm.config.ServiceStateLocation = serviceStatesTestUpdates

	serviceStatesTestUpdatesLock.Lock()
	defer serviceStatesTestUpdatesLock.Unlock()

	if startupErr := lm.Startup(); startupErr != nil {
		lm.logger.Warn("Failed to to start Licence manager service : %v \n", startupErr.Error())
	}

	defer func() {
		_ = lm.Shutdown()
	}()

	// Service can be found, and its state is updated only in the internal data structures
	err := lm.services["untangle-node-discovery"].setServiceState(StateEnable)
	assert.NoError(t, err)
	assert.Equal(t, StateEnable, lm.services["untangle-node-discovery"].State.getAllowedState())

	// Verify the actual file was updated
	err = lm.services["untangle-node-discovery"].setServiceState(StateDisable)
	assert.NoError(t, err)
	assert.Equal(t, StateDisable, lm.services["untangle-node-discovery"].State.getAllowedState())

	statesFromFile, err := LoadServiceStates(serviceStatesTestUpdates)
	assert.NoError(t, err)
	var actualServiceState *ServiceState

	// Get state of the service from the file
	for _, v := range statesFromFile {
		if v.Name == "untangle-node-discovery" {
			actualServiceState = &v
			break
		}
		continue
	}

	// Fail the test if the service couldn't be found. Nothing to
	if !assert.NotNil(t, actualServiceState) {
		t.Fatalf("The actual service could not be found in the state file %v\n", actualServiceState)
	}

	assert.Equal(t, StateDisable, actualServiceState.getAllowedState())
}

// Licensemanager configs are read in at startup. Any function
// that needs the state of licensemanager after startup is tested
// using the suite
type LicenseManagerTestSuite struct {
	suite.Suite
	lm *LicenseManager

	// Expected service states when the suite starts
	// the licensemanager
	expectedServiceStates map[string]*Service
}

func (suite *LicenseManagerTestSuite) SetupSuite() {
	ServicesAllowedStatesLocation = "./testdata/allowedstates/"
	suite.lm = NewLicenseManager(getTestConfig(), &mocks.LoggerHelper{})

	// Swap out funciton used to interact with the CLI
	suite.lm.RefreshLicenses = func() error { return nil }

	suite.expectedServiceStates = map[string]*Service{
		"untangle-node-discovery":         {Name: "untangle-node-discovery", State: ServiceState{AllowedState: 0}},
		"untangle-node-classd":            {Name: "untangle-node-classd", State: ServiceState{AllowedState: 1}},
		"untangle-node-threat-prevention": {Name: "untangle-node-threat-prevention", State: ServiceState{AllowedState: 1}},
		"untangle-node-sitefilter":        {Name: "untangle-node-sitefilter", State: ServiceState{AllowedState: 0}},
		"untangle-node-geoip":             {Name: "untangle-node-geoip", State: ServiceState{AllowedState: 0}},
		"untangle-node-captiveportal":     {Name: "untangle-node-captiveportal", State: ServiceState{AllowedState: 0}},
	}

	if startupErr := suite.lm.Startup(); startupErr != nil {
		suite.lm.logger.Warn("Failed to to start Licence manager service : %v \n", startupErr.Error())
	}
}

func (suite *LicenseManagerTestSuite) TearDownSuite() {
	_ = suite.lm.Shutdown()
}

func TestLicenseManagerTestSuite(t *testing.T) {
	suite.Run(t, new(LicenseManagerTestSuite))
}

// Verify the state of a licensemanager after startup is ran
func (suite *LicenseManagerTestSuite) TestStartup() {
	for k, v := range suite.expectedServiceStates {
		if !(suite.Contains(suite.lm.services, k)) {
			continue
		}
		suite.Equal(v.Name, suite.lm.services[k].Name, "Name mismatch for Service named: %v\n", k)
		suite.Equal(v.State.AllowedState, suite.lm.services[k].State.AllowedState, "State mismatch for Service named: %v\n", k)
	}
}

func (suite *LicenseManagerTestSuite) TestFindService() {
	actualService, err := suite.lm.findService("untangle-node-discovery")
	suite.NoError(err)
	suite.NotNil(actualService)
}

func (suite *LicenseManagerTestSuite) TestIsLicenseEnabled() {
	// Test getting a service that is known to be enabled
	enabled, err := suite.lm.IsLicenseEnabled("untangle-node-discovery")
	suite.NoError(err)
	suite.True(enabled)

	// Test getting a service that is known to be disabled
	enabled, err = suite.lm.IsLicenseEnabled("untangle-node-threat-prevention")
	suite.NoError(err)
	suite.False(enabled)

	// Test getting a service that does not exist
	enabled, err = suite.lm.IsLicenseEnabled("bogus-service")
	suite.Error(err)
	suite.False(enabled)
}

func (suite *LicenseManagerTestSuite) TestGetServices() {
	allowedStateDirLock.Lock()
	defer allowedStateDirLock.Unlock()
	actual := suite.lm.GetServices()

	for k, v := range suite.expectedServiceStates {
		if !(suite.Contains(actual, k)) {
			continue
		}
		suite.Equal(v.Name, actual[k].Name, "Name mismatch for Service named: %v\n", k)
		suite.Equal(v.State.AllowedState, actual[k].State.AllowedState, "State mismatch for Service named: %v\n", k)
	}
}

// Returns a licensemanager config that can by default
// be used by most tests.
func getTestConfig() *Config {
	apps := map[string]ServiceHook{
		// Setting Star/Stop to non-nil
		// prevents sighups from being sent
		// out to non-existent binaries
		"untangle-node-captiveportal": {
			Start:    func() {},
			Stop:     func() {},
			Enabled:  nil,
			Disabled: disableCaptivePortal,
		},
		"untangle-node-threat-prevention": {
			Start:    func() {},
			Stop:     func() {},
			Enabled:  nil,
			Disabled: disableThreatPrevention,
		},
		"untangle-node-sitefilter": {
			Start:    func() {},
			Stop:     func() {},
			Enabled:  nil,
			Disabled: disableWebFilter,
		},
		"untangle-node-geoip": {
			Start:    func() {},
			Stop:     func() {},
			Enabled:  nil,
			Disabled: disableGeoipFilter,
		},
		"untangle-node-discovery": {
			Start:    func() {},
			Stop:     func() {},
			Enabled:  nil,
			Disabled: disableDiscovery,
		},
		"untangle-node-classd": {
			Start:    func() {},
			Stop:     func() {},
			Enabled:  nil,
			Disabled: disableApplicationControl,
		},
	}

	return &Config{
		LicenseLocation:      licenseLocation,
		ServiceStateLocation: servicesStatesMixedFile,
		WatchDogInterval:     (6*time.Hour + 5*time.Minute),
		ValidServiceHooks:    apps,
		Executable:           "testd",
	}
}

// DisableCaptivePortal
func disableCaptivePortal() (interface{}, []string, error) {
	return false, []string{"captiveportal", "enabled"}, nil
}

// DisableThreatPrevention
func disableThreatPrevention() (interface{}, []string, error) {
	return false, []string{"threatprevention", "enabled"}, nil
}

// DisableWebFilter
func disableWebFilter() (interface{}, []string, error) {
	return false, []string{"webfilter", "enabled"}, nil
}

// DisableGeoipFilter
func disableGeoipFilter() (interface{}, []string, error) {
	return false, []string{"geoip", "enabled"}, nil
}

// DisableDiscovery
func disableDiscovery() (interface{}, []string, error) {
	return false, []string{"discovery", "enabled"}, nil
}

// DisableApplicationControl
func disableApplicationControl() (interface{}, []string, error) {
	return false, []string{"application_control", "enabled"}, nil
}
