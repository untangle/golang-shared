package arp

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sync"
	"time"

	"github.com/untangle/discoverd/plugins/discovery"
	"github.com/untangle/discoverd/utils"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/structs/interfaces"
	disco_proto "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

const (
	pluginName string = "arp"
)

var (
	arpSingleton *Arp
	once         sync.Once

	settingsPath []string = []string{"discovery", "plugins"}
)

type arpPluginType struct {
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	AutoInterval uint   `json:"autoInterval"`
}

// Setup the Arp struct as a singleton
type Arp struct {
	autoArpCollectionChan chan bool
	arpSettings           arpPluginType
}

// Gets a singleton instance of the Arp plugin
func NewArp() *Arp {
	once.Do(func() {
		arpSingleton = &Arp{autoArpCollectionChan: make(chan bool)}
	})

	return arpSingleton
}

func (arp *Arp) InSync(settings interface{}) bool {
	newSettings, ok := settings.(arpPluginType)
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

func (arp *Arp) GetSettingsStruct() (interface{}, error) {
	var fileSettings []arpPluginType
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

func (arp *Arp) SyncSettings(settings interface{}) error {

	originalSettings := arp.arpSettings
	newSettings, ok := settings.(arpPluginType)
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
func (arp *Arp) Startup() error {
	logger.Info("Starting ARP collector plugin\n")

	// Grab the initial settings on startup
	settings, err := arp.GetSettingsStruct()
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

	discovery.NewDiscovery().UnregisterCollector(pluginName)

	return nil
}

func (arp *Arp) startArp() {
	discovery.NewDiscovery().RegisterCollector(pluginName, ArpcallBackHandler)

	// Lets do a first run to get the initial data
	ArpcallBackHandler(nil)

	arp.startAutoArpCollection()
}

func (arp *Arp) autoArpCollection() {
	logger.Debug("Starting automatic collection from ARP plugin with an interval of %d seconds\n", arp.arpSettings.AutoInterval)
	for {
		select {
		case <-arp.autoArpCollectionChan:
			logger.Debug("Stopping automatic collection from ARP plugin\n")
			arp.autoArpCollectionChan <- true
			return
		case <-time.After(time.Duration(arp.arpSettings.AutoInterval) * time.Second):
			ArpcallBackHandler(nil)
		}
	}
}

func (arp *Arp) startAutoArpCollection() {
	go arp.autoArpCollection()
}

func (arp *Arp) stopAutoArpCollection() {
	// The send to kill the AutoNmapCollection goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case arp.autoArpCollectionChan <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-arp.autoArpCollectionChan:
		logger.Info("Successful shutdown of the automatic ARP collector\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown automatic ARP collector. It may never have been started\n")
	}
}

// StringSet is a set of strings, backed by a map.
type StringSet struct {
	StringMap map[string]struct{}
}

// NewStringSet creates a new StringSet with no members.
func NewStringSet() *StringSet {
	return &StringSet{
		StringMap: map[string]struct{}{},
	}
}

// Add adds an item to the StringSet.
func (set *StringSet) Add(item string) {
	set.StringMap[item] = struct{}{}
}

// Contains tests for membership of item in StringSet.
func (set *StringSet) Contains(item string) (contains bool) {
	_, contains = set.StringMap[item]
	return
}

// arpScanner scans a linux /proc/net/arp -format file.
type arpScanner struct {
	settings    *settings.SettingsFile
	arpFileName string
	entryList   []*disc.DeviceEntry
	wans        *StringSet
}

func (scanner *arpScanner) buildWANList() error {
	intfList := []*interfaces.Interface{}
	if err := scanner.settings.UnmarshalSettingsAtPath(&intfList, "network", "interfaces"); err != nil {
		return fmt.Errorf("buildWANList: couldn't unmarshall network settings: %w", err)
	}

	for _, intf := range intfList {
		if intf.IsWAN {
			scanner.wans.Add(intf.Device)
		}
	}
	return nil

}

var arpPattern = (utils.IPv4Regex + `\s+` + // ipv4 is capture groups 2-3
	utils.HexRegex + `\s+` + // HW type is 4-5
	utils.HexRegex + `\s+` + // Flags is 6-7
	utils.MacRegex + `\s+` + // mac is 8-9
	utils.MaskRegex + `\s+` + // mask is 10-11
	utils.DeviceRegex) // device is 12-13
var arpLineRegex *regexp.Regexp = regexp.MustCompile(arpPattern)

const (
	arpIPGroupBegin     = 2
	arpMacGroupBegin    = 8
	arpDeviceGroupBegin = 12
)

// newArpScanner creates an arp scanner that uses the settings
// SettingsFile to figure out what WAN devices there are, and reads
// arp entries frorm arpFilename, (use /proc/net/arp).
func newArpScanner(settings *settings.SettingsFile,
	arpFileName string) *arpScanner {
	scanner := &arpScanner{
		arpFileName: arpFileName,
		settings:    settings,
		wans:        NewStringSet(),
	}

	return scanner
}

// scanLineForEntries scans a single line of the arp file for an arp
// entry and parses it. If it is not on a WAN interface, it's added to
// the internal device list.
func (scanner *arpScanner) scanLineForEntries(line []byte) {
	indices := arpLineRegex.FindSubmatchIndex(line)
	if len(indices) <= arpDeviceGroupBegin+1 {
		return
	}
	ip := string(line[indices[arpIPGroupBegin]:indices[arpIPGroupBegin+1]])
	mac := string(line[indices[arpMacGroupBegin]:indices[arpMacGroupBegin+1]])
	if scanner.wans.Contains(string(line[indices[arpDeviceGroupBegin]:indices[arpDeviceGroupBegin+1]])) {
		return
	}

	scanner.entryList = append(scanner.entryList,
		&disc.DeviceEntry{
			DiscoveryEntry: disco_proto.DiscoveryEntry{
				MacAddress:  mac,
				IPv4Address: ip,
				Arp: &disco_proto.ARP{
					Ip:         ip,
					Mac:        mac,
					LastUpdate: time.Now().Unix(),
				},
			},
		})
}

// getArpEntriesFromFile gets all arp entries from the file given in
// the constructor that are not on WAN interfaces. It returns these as
// device entries, or if an error occurs, nil and an error.
func (scanner *arpScanner) getArpEntriesFromFile() ([]*disc.DeviceEntry, error) {
	scanner.entryList = []*disc.DeviceEntry{}
	arp, err := os.Open(scanner.arpFileName)

	if err != nil {
		return nil, fmt.Errorf("couldn't open arp file %s: %w", scanner.arpFileName, err)
	}

	if err := scanner.buildWANList(); err != nil {
		return nil, fmt.Errorf("unable to build WAN list: %w", err)
	}
	fileScanner := bufio.NewScanner(arp)

	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		scanner.scanLineForEntries(line)
	}

	if fileScanner.Err() != nil {
		return nil, fmt.Errorf("couldn't scan arp file: %w", fileScanner.Err())
	}
	if err := arp.Close(); err != nil {
		logger.Debug("Couldn't close arp file: %s\n", err)
	}
	return scanner.entryList, nil
}

// ArpcallBackHandler is the callback handler for the ARP collector
func ArpcallBackHandler(commands []discovery.Command) {
	logger.Debug("Arp Callback handler: Received %d commands\n", len(commands))
	scanner := newArpScanner(
		settings.GetSettingsFileSingleton(),
		"/proc/net/arp")
	entries, err := scanner.getArpEntriesFromFile()
	if err != nil {
		logger.Warn("Couldn't scan /proc/net/arp for devices: %s", err)
		return
	}
	for _, entry := range entries {
		discovery.UpdateDiscoveryEntry(entry.MacAddress, entry)
	}
}
