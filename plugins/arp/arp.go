package arp

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"syscall"
	"time"

	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/discoverd/utils"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/structs/interfaces"
	disco_proto "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

const (
	pluginName          string = "arp"
	enabledDefault      bool   = false
	autoIntervalDefault uint   = math.MaxUint32
)

var autoArpCollectionChan chan bool

type arpSettingType struct {
	Enabled      bool `json:"enabled"`
	AutoInterval uint `json:"autoInterval"`
}

var arpSettings *arpSettingType

func init() {
	autoArpCollectionChan = make(chan bool)
}

// Start starts the ARP collector IF enabled in the settings file
func Start() {
	logger.Info("Starting ARP collector plugin\n")

	// Load in settings. Plugin will be started by syncCallbackHandler() if
	// the plugin is enabled
	syncCallbackHandler()
}

func startArp() {
	discovery.RegisterCollector(pluginName, ArpcallBackHandler)

	// Lets do a first run to get the initial data
	ArpcallBackHandler(nil)

	startAutoArpCollection()
}

func autoArpCollection() {
	logger.Debug("Starting automatic collection from ARP plugin with an interval of %d seconds\n", arpSettings.AutoInterval)
	for {
		select {
		case <-autoArpCollectionChan:
			logger.Debug("Stopping automatic collection from ARP plugin\n")
			autoArpCollectionChan <- true
			return
		case <-time.After(time.Duration(arpSettings.AutoInterval) * time.Second):
			ArpcallBackHandler(nil)
		}
	}
}

func startAutoArpCollection() {
	go autoArpCollection()
}

func stopAutoArpCollection() {
	autoArpCollectionChan <- true

	select {
	case <-autoArpCollectionChan:
		logger.Info("Successful shutdown of the automatic ARP collector\n")
	case <-time.After(10 * time.Second):
		logger.Warn("Failed to shutdown automatic ARP collector\n")
	}
}

// PluginSignal handles a sighup seen, turning WF on or off
func PluginSignal(message syscall.Signal) {
	logger.Info("PluginSignal(%s) has been called\n", pluginName)
	switch message {
	case syscall.SIGHUP:
		syncCallbackHandler()
	}
}

func syncCallbackHandler() {
	logger.Debug("Syncing %s settings\n", pluginName)
	var systemArpsettings interface{}

	// Get current state of plugin so the automatic running of
	// the arp plugin can be restarted if necessary. The plugin should only
	// be restarted if it was running previously and the settings were altered.
	restartPlugin := false
	if (arpSettings != nil) && (arpSettings.Enabled) {
		if arpSettings.Enabled {
			restartPlugin = true
		}
	}

	// Avoid failing on errors if the plugin settings weren't read in correctly, just use the defaults.
	// All the logging for a bad settings file is done in createSettings()
	systemArpsettings, _ = getPluginSettings("discovery", pluginName)
	createSettings(systemArpsettings.(map[string]interface{}))

	if arpSettings.Enabled {
		if restartPlugin {
			Stop()
		}

		logger.Debug("Starting %s Plugin", pluginName)
		startArp()
	} else {
		logger.Debug("Stopping %s Plugin", pluginName)
		Stop()
	}
}

// Get a plugin's settings json as a map. Takes daemonName to search for a plugin in
// in a daemon's list of plugins and pluginType to search for a plugin's settings
func getPluginSettings(daemonName string, pluginType string) (interface{}, error) {
	daemonSettings, err := settings.GetSettings([]string{daemonName})

	if err != nil {
		// Return the JSON error object returned by GetSettings
		return daemonSettings, err
	}

	pluginSettings, ok := daemonSettings.(map[string]interface{})["plugins"]
	if !ok {
		return nil, fmt.Errorf("no plugin settings found for the %s deamon", daemonName)
	}

	for _, pluginSetting := range pluginSettings.([]interface{}) {
		plugType, ok := pluginSetting.(map[string]interface{})["type"]
		if ok && (plugType == pluginType) {
			return pluginSetting, nil
		}
	}

	return nil, fmt.Errorf("he plugin settings for %s could not be found for the %s daemon", pluginType, daemonName)
}

// Sets settings for the ARP Plugin. If no settings json was present, just uses defaults
func createSettings(m map[string]interface{}) {
	arpSettings = &arpSettingType{Enabled: enabledDefault, AutoInterval: autoIntervalDefault}

	if m == nil {
		logger.Warn("No enabled setting for ARP provided, using the defaults\n")
	}

	if m["autoInterval"] != nil {
		arpSettings.AutoInterval = uint(m["autoInterval"].(float64))
	} else {
		logger.Warn("No autointerval setting for ARP provided, using the default of %d sec\n", autoIntervalDefault)
	}

	if m["enabled"] != nil {
		arpSettings.Enabled = m["enabled"].(bool)
	} else {
		logger.Warn("No enabled setting for ARP provided, using the default of %t sec\n", enabledDefault)
	}
}

// Stop stops QoS
func Stop() {
	logger.Info("Stopping ARP collector plugin\n")
	discovery.UnregisterCollector(pluginName)
	stopAutoArpCollection()
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
