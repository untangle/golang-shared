package lldp

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"sync"
	"time"

	"github.com/untangle/discoverd/plugins/discovery"
	"github.com/untangle/discoverd/utils"
	"github.com/untangle/golang-shared/plugins"
	"github.com/untangle/golang-shared/plugins/zmqmsg"
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
	MgmtIp     []value             `json:"mgmt-ip"`
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

func init() {
	plugins.GlobalPluginControl().RegisterPlugin(NewLldp)
}

type lldpPluginSettings struct {
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	AutoInterval uint   `json:"autoInterval"`
}

type Lldp struct {
	autoLldpCollectionShutdown    chan bool
	autoLldpCollectionShutdownAck chan bool

	lldpSettings lldpPluginSettings
}

// Gets a singleton instance of the Lldp plugin
func NewLldp() *Lldp {
	once.Do(func() {
		lldpSingleton = &Lldp{autoLldpCollectionShutdown: make(chan bool),
			autoLldpCollectionShutdownAck: make(chan bool)}
	})

	return lldpSingleton
}

// Returns true if the current settings match the 'new' settings Provided, otherwise false
func (lldp *Lldp) InSync(settings interface{}) bool {
	newSettings, ok := settings.(lldpPluginSettings)
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

// Returns a struct containing the plugins settings of type lldpPluginSettings
func (lldp *Lldp) GetCurrentSettingsStruct() (interface{}, error) {
	var fileSettings []lldpPluginSettings
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

// Updates the current settings with the settings passed in. If the plugin was already running
// but the settings changed, the plugin is restarted.
// An error is returned if the settings can't be synced
func (lldp *Lldp) SyncSettings(settings interface{}) error {

	originalSettings := lldp.lldpSettings
	newSettings, ok := settings.(lldpPluginSettings)
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
	discovery.NewDiscovery().DeregisterCollector(pluginName)
	lldp.stopAutoLldpCollection()

	return nil
}

// Start starts the LLDP collector
// Meant to only be run once
func (lldp *Lldp) Startup() error {
	logger.Info("Starting LLDP collector plugin\n")

	// Grab the initial settings on startup
	settings, err := lldp.GetCurrentSettingsStruct()
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

// Start method of the plugin. Meant to be used in a restart of the plugin
func (lldp *Lldp) startLldp() {
	discovery.NewDiscovery().RegisterCollector(pluginName, LldpcallBackHandler)

	LldpcallBackHandler(nil)

	go lldp.autoLldpCollection()
}

func (lldp *Lldp) stopAutoLldpCollection() {
	// The send to kill the AutoLldpCollection goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case lldp.autoLldpCollectionShutdown <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-lldp.autoLldpCollectionShutdownAck:
		logger.Info("Successful shutdown of the automatic LLDP collector\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown automatic LLDP collector. It may never have been started\n")
	}
}

// Runs the plugin's handler on a timer. Meant to be run as a goroutine
func (lldp *Lldp) autoLldpCollection() {
	logger.Debug("Starting automatic collection from LLDP plugin with an interval of %d seconds\n", lldp.lldpSettings.AutoInterval)
	for {
		select {
		case <-lldp.autoLldpCollectionShutdown:
			logger.Debug("Stopping automatic collection from LLDP plugin\n")
			lldp.autoLldpCollectionShutdownAck <- true
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
		entry.Lldp = []*Discoverd.LLDP{}
		lldp := &Discoverd.LLDP{}

		if len(intf.Chassis) > 0 {
			chassis := intf.Chassis[0]
			// mac discovery
			if len(chassis.ID) > 0 {
				if chassis.ID[0].Type == "mac" {

					if !utils.IsMacAddress(chassis.ID[0].Value) {
						continue
					}
					lldp.Mac = chassis.ID[0].Value
				}
			}

			if len(chassis.Name) > 0 {
				lldp.SysName = chassis.Name[0].Value
			}
			if len(chassis.Desc) > 0 {
				lldp.SysDesc = chassis.Desc[0].Value
			}
			if len(chassis.MgmtIp) > 0 {
				lldp.Ip = chassis.MgmtIp[0].Value
			}

			// chasis capabilities
			if len(chassis.Capability) > 0 {
				for _, val := range chassis.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Enabled
					lldp.ChassisCapabilities = append(lldp.ChassisCapabilities, cap)
				}
			}
		}

		if len(intf.LldpMed) > 0 {
			lldpmed := intf.LldpMed[0]

			// LLDP-MED inventory
			if len(lldpmed.Inventory) > 0 {
				inv := lldpmed.Inventory[0]

				if len(inv.Hardware) > 0 {
					lldp.InventoryHWRev = inv.Hardware[0].Value
				}
				if len(inv.Software) > 0 {
					lldp.InventorySoftRev = inv.Software[0].Value
				}
				if len(inv.Serial) > 0 {
					lldp.InventorySerial = inv.Serial[0].Value
				}
				if len(inv.Model) > 0 {
					lldp.InventoryModel = inv.Model[0].Value
				}
				if len(inv.Manufacturer) > 0 {
					lldp.InventoryVendor = inv.Manufacturer[0].Value
				}
			}

			// LLDP-MED capabilities
			if len(lldpmed.Capability) > 0 {
				for _, val := range lldpmed.Capability {
					cap := &Discoverd.LLDPCapabilities{}
					cap.Capability = val.Type
					cap.Enabled = val.Available
					lldp.MedCapabilities = append(lldp.MedCapabilities, cap)
				}
			}
		}
		lldp.LastUpdate = time.Now().Unix()
		logger.Info("lldp %v\n", lldp)
		entry.Lldp = append(entry.Lldp, lldp)

		entry.MacAddress = lldp.Mac
		entry.LastUpdate = time.Now().Unix()
		discovery.ZmqpublishEntry(entry, zmqmsg.NMAPDeviceZMQTopic)
		logger.Info("lldp entry%v\n", entry)
		//discovery.UpdateDiscoveryEntry(mac, entry)

	}
}