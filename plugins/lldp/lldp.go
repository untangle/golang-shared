package lldp

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"sync"
	"time"

	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type jsonData struct {
	Lldp []lldp `json:"lldp"`
}

type lldp struct {
	Intf []intf `json:"interface"`
}

type intf struct {
	Chassis []chassis `json:"chassis"`
	LldpMed []lldpMed `json:"lldp-med"`
}

type chassis struct {
	ID         []identity          `json:"id"`
	Name       []value             `json:"name"`
	Desc       []value             `json:"descr"`
	Capability []chassisCapability `json:"capability"`
}

type identity struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type chassisCapability struct {
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type lldpMed struct {
	DeviceType []value             `json:"device-type"`
	Capability []lldpMedCapability `json:"capability"`
	Inventory  []lldpMedInventory  `json:"inventory"`
}

type lldpMedCapability struct {
	Type      string `json:"type,omitempty"`
	Available bool   `json:"available,omitempty"`
}

type lldpMedInventory struct {
	Hardware     []value `json:"hardware"`
	Software     []value `json:"software"`
	Serial       []value `json:"serial"`
	Manufacturer []value `json:"manufacturer"`
	Model        []value `json:"model"`
}

type value struct {
	Value string `json:"value,omitempty"`
}

const (
	pluginName string = "lldp"
)

var (
	lldpSingleton *Lldp
	once          sync.Once

	settingsPath []string = []string{"discovery", "plugins"}
)

type lldpPluginType struct {
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	AutoInterval uint   `json:"autoInterval"`
}

// Setup the Lldp struct as a singleton
type Lldp struct {
	autoLldpCollectionChan chan bool
	lldpSettings           lldpPluginType
}

// Gets a singleton instance of the Lldp plugin
func NewLldp() *Lldp {
	once.Do(func() {
		lldpSingleton = &Lldp{autoLldpCollectionChan: make(chan bool)}
	})

	return lldpSingleton
}

func (lldp *Lldp) InSync(settings interface{}) bool {
	newSettings, ok := settings.(lldpPluginType)
	if !ok {
		logger.Warn("LLDP: Could not compare the settings file provided to the current plugin settings. The settings cannot be updated.")
		return false
	}

	if newSettings == lldp.lldpSettings {
		logger.Debug("Settings remain unchanged for the LLDP plugin\n")
		return true
	}

	logger.Info("Updating LLDP plugin settings")
	return false
}

func (lldp *Lldp) GetSettingsStruct() (interface{}, error) {
	var fileSettings []lldpPluginType
	if err := settings.UnmarshalSettingsAtPath(&fileSettings, settingsPath...); err != nil {
		return nil, fmt.Errorf("LLDP: %s", err.Error())
	}

	// Plugins are in an array in the settings.json. Have to go through all of them
	// to find the desired settings struct
	for _, pluginSetting := range fileSettings {
		if pluginSetting.Type == pluginName {
			return pluginSetting, nil
		}
	}

	return nil, fmt.Errorf("no settings could be found for %s plugin", pluginName)
}

// Returns name of the plugin.
// The function is not static to satisfy the SettingsSyncer interface requirements
func (lldp *Lldp) Name() string {
	return pluginName
}

// Returns name of the plugin.
// The function is not static to satisfy the SettingsSyncer interface requirements
func (lldp *Lldp) SyncSettings(settings interface{}) error {

	originalSettings := lldp.lldpSettings
	newSettings, ok := settings.(lldpPluginType)
	if !ok {
		return fmt.Errorf("LLDP: Settings provided were %s but expected %s",
			reflect.TypeOf(settings).String(), reflect.TypeOf(lldp.lldpSettings).String())
	}

	lldp.lldpSettings = newSettings

	// If settings changed but the plugin was previously enabled, restart the plugin
	// for changes to take effect
	var shutdownError error
	if originalSettings.Enabled && lldp.lldpSettings.Enabled {
		shutdownError = lldp.Shutdown()
	}

	if lldp.lldpSettings.Enabled {
		lldp.startLldp()
	} else {
		shutdownError = lldp.Shutdown()
	}

	return shutdownError
}

// Stop stops LLDP collector
func (lldp *Lldp) Shutdown() error {
	logger.Info("Stopping LLDP collector plugin\n")
	discovery.UnregisterCollector(pluginName)
	lldp.stopAutoLldpCollection()

	return nil
}

// Start starts the LLDP collector
func (lldp *Lldp) Startup() error {
	logger.Info("Starting LLDP collector plugin\n")

	// Grab the initial settings on startup
	settings, err := lldp.GetSettingsStruct()
	if err != nil {
		return err
	}

	// SyncSettings will start the plugin if it's enabled
	err = lldp.SyncSettings(settings)
	if err != nil {
		return err
	}

	return nil
}

func (lldp *Lldp) startLldp() {
	discovery.RegisterCollector(pluginName, LldpcallBackHandler)

	LldpcallBackHandler(nil)

	lldp.startAutoLldpCollection()
}

func (lldp *Lldp) startAutoLldpCollection() {
	go lldp.autoLldpCollection()
}

func (lldp *Lldp) stopAutoLldpCollection() {
	// The send to kill the AutoLldpCollection goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case lldp.autoLldpCollectionChan <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-lldp.autoLldpCollectionChan:
		logger.Info("Successful shutdown of the automatic LLDP collector\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown automatic LLDP collector. It may never have been started\n")
	}
}

func (lldp *Lldp) autoLldpCollection() {
	logger.Debug("Starting automatic collection from LLDP plugin with an interval of %d seconds\n", lldp.lldpSettings.AutoInterval)
	for {
		select {
		case <-lldp.autoLldpCollectionChan:
			logger.Debug("Stopping automatic collection from LLDP plugin\n")
			lldp.autoLldpCollectionChan <- true
			return
		case <-time.After(time.Duration(lldp.lldpSettings.AutoInterval) * time.Second):
			LldpcallBackHandler(nil)
		}
	}
}

// LldpcallBackHandler is the callback handler for the LLDP collector
func LldpcallBackHandler(commands []discovery.Command) {
	logger.Debug("LLDP neighbors callback handler: Received %d commands\n", len(commands))

	// run neighbors command
	cmd := exec.Command("lldpcli", "-f", "json0", "show", "n", "details")
	output, _ := cmd.CombinedOutput()

	// logger.Info("LLDP output: %s\n", string(output))

	// parse json output data
	var result jsonData
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		logger.Err("Unable to unmarshal json: %s\n", err)
		return
	}

	// return on empty data
	if len(result.Lldp[0].Intf) == 0 {
		logger.Debug("No LLDP neighbors found!\n")
		return
	}

	// iterate over interface items
	for _, intf := range result.Lldp[0].Intf {
		// initialize the discovery entry
		entry := &disc.DeviceEntry{}
		entry.Init()
		entry.Lldp = &Discoverd.LLDP{}

		// mac is used as id for discovery entry update
		var mac = ""

		if len(intf.Chassis) > 0 {
			chassis := intf.Chassis[0]
			// mac discovery
			if len(chassis.ID) > 0 {
				if chassis.ID[0].Type == "mac" {
					mac = chassis.ID[0].Value
				}
			}

			if len(chassis.Name) > 0 {
				entry.Lldp.SysName = chassis.Name[0].Value
			}
			if len(chassis.Desc) > 0 {
				entry.Lldp.SysDesc = chassis.Desc[0].Value
			}

			// chasis capabilities
			if len(chassis.Capability) > 0 {
				for _, val := range chassis.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Enabled
					entry.Lldp.ChassisCapabilities = append(entry.Lldp.ChassisCapabilities, cap)
				}
			}
		}

		if len(intf.LldpMed) > 0 {
			lldpmed := intf.LldpMed[0]

			// LLDP-MED inventory
			if len(lldpmed.Inventory) > 0 {
				inv := lldpmed.Inventory[0]

				if len(inv.Hardware) > 0 {
					entry.Lldp.InventoryHWRev = inv.Hardware[0].Value
				}
				if len(inv.Software) > 0 {
					entry.Lldp.InventorySoftRev = inv.Software[0].Value
				}
				if len(inv.Serial) > 0 {
					entry.Lldp.InventorySerial = inv.Serial[0].Value
				}
				if len(inv.Model) > 0 {
					entry.Lldp.InventoryModel = inv.Model[0].Value
				}
				if len(inv.Manufacturer) > 0 {
					entry.Lldp.InventoryVendor = inv.Manufacturer[0].Value
				}
			}

			// LLDP-MED capabilities
			if len(lldpmed.Capability) > 0 {
				for _, val := range lldpmed.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Available
					entry.Lldp.MedCapabilities = append(entry.Lldp.MedCapabilities, cap)
				}
			}
		}
		entry.Lldp.LastUpdate = time.Now().Unix()

		if mac != "" {
			entry.MacAddress = mac
			discovery.UpdateDiscoveryEntry(mac, entry)
		}
	}
}
