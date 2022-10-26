package arp

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/untangle/discoverd/plugins/discovery"
	"github.com/untangle/golang-shared/plugins"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	pluginName string = "arp"
)

var (
	arpSingleton *Arp
	once         sync.Once

	settingsPath []string = []string{"discovery", "plugins"}
)

func init() {
	plugins.GlobalPluginControl().RegisterPlugin(NewArp)
}

type arpPluginSettings struct {
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	AutoInterval uint   `json:"autoInterval"`
}

// Setup the Arp struct as a singleton
type Arp struct {
	// During shutdown, goroutines are stopped. If they are already
	// stopped, and shutdown gets called again the process could block.
	// Use two channels for a non-blocking shutdown and ack
	autoArpCollectionShutdown    chan bool
	autoArpCollectionShutdownAck chan bool

	arpSettings arpPluginSettings
}

// Gets a singleton instance of the Arp plugin
func NewArp() *Arp {
	once.Do(func() {
		arpSingleton = &Arp{autoArpCollectionShutdown: make(chan bool),
			autoArpCollectionShutdownAck: make(chan bool)}
	})

	return arpSingleton
}

// Returns true if the current settings match the 'new' settings Provided, otherwise false
func (arp *Arp) InSync(settings interface{}) bool {
	newSettings, ok := settings.(arpPluginSettings)
	if !ok {
		logger.Warn("Arp: Could not compare the settings file provided to the current plugin settings. The settings cannot be updated.")
		return false
	}

	if newSettings == arp.arpSettings {
		logger.Debug("Settings remain unchanged for the ARP plugin\n")
		return true
	}

	logger.Info("Updating ARP plugin settings\n")
	return false
}

// Returns a struct containing the plugins settings of type arpPluginSettings
func (arp *Arp) GetCurrentSettingsStruct() (interface{}, error) {
	var fileSettings []arpPluginSettings
	if err := settings.UnmarshalSettingsAtPath(&fileSettings, settingsPath...); err != nil {
		return nil, fmt.Errorf("ARP: %s", err.Error())
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
func (arp *Arp) Name() string {
	return pluginName
}

// Updates the current settings with the settings passed in. If the plugin was already running
// but the settings changed, the plugin is restarted.
// An error is returned if the settings can't be synced
func (arp *Arp) SyncSettings(settings interface{}) error {

	originalSettings := arp.arpSettings
	newSettings, ok := settings.(arpPluginSettings)
	if !ok {
		return fmt.Errorf("ARP: Settings provided were %s but expected %s",
			reflect.TypeOf(settings).String(), reflect.TypeOf(arp.arpSettings).String())
	}

	arp.arpSettings = newSettings

	// If settings changed but the plugin was previously enabled, restart the plugin
	// for changes to take effect
	var shutdownError error
	if originalSettings.Enabled && arp.arpSettings.Enabled {
		shutdownError = arp.Shutdown()
	}

	if arp.arpSettings.Enabled {
		arp.startArp()
	} else {
		shutdownError = arp.Shutdown()
	}

	return shutdownError
}

// Startup starts the ARP collector IF enabled in the settings file
// Meant to only be run once
func (arp *Arp) Startup() error {
	logger.Info("Starting ARP collector plugin\n")

	// Grab the initial settings on startup
	settings, err := arp.GetCurrentSettingsStruct()
	if err != nil {
		return err
	}

	// SyncSettings will start the plugin if it's enabled
	err = arp.SyncSettings(settings)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown stops QoS
func (arp *Arp) Shutdown() error {
	logger.Info("Stopping ARP collector plugin\n")

	arp.stopAutoArpCollection()

	discovery.NewDiscovery().DeregisterCollector(pluginName)

	return nil
}

// Start method of the plugin. Meant to be used in a restart of the plugin
func (arp *Arp) startArp() {
	discovery.NewDiscovery().RegisterCollector(pluginName, NetlinkNeighbourCallbackController)

	// Lets do a first run to get the initial data
	NetlinkNeighbourCallbackController(nil)

	go arp.autoArpCollection()
}

// Runs the plugin's handler on a timer. Meant to be run as a goroutine
func (arp *Arp) autoArpCollection() {
	logger.Debug("Starting automatic collection from ARP plugin with an interval of %d seconds\n", arp.arpSettings.AutoInterval)
	for {
		select {
		case <-arp.autoArpCollectionShutdown:
			logger.Debug("Stopping automatic collection from ARP plugin\n")
			arp.autoArpCollectionShutdownAck <- true
			return
		case <-time.After(time.Duration(arp.arpSettings.AutoInterval) * time.Second):
			NetlinkNeighbourCallbackController(nil)
		}
	}
}

// Stops automatically running the ARP callback handler on a timer
func (arp *Arp) stopAutoArpCollection() {
	// The send to kill the AutoNmapCollection goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case arp.autoArpCollectionShutdown <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-arp.autoArpCollectionShutdownAck:
		logger.Info("Successful shutdown of the automatic ARP collector\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown automatic ARP collector. It may never have been started\n")
	}
}
