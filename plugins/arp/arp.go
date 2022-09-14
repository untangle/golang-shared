package arp

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/discoverd/utils"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/structs/interfaces"
	disco_proto "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// Start starts the ARP collector
func Start() {
	logger.Info("Starting ARP collector plugin\n")
	discovery.RegisterCollector(ArpcallBackHandler)
	// Lets do a first run to get the initial data
	ArpcallBackHandler(nil)
}

// Stop stops QoS
func Stop() {
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

func newArpScanner(settings *settings.SettingsFile,
	arpFileName string) *arpScanner {
	scanner := &arpScanner{
		arpFileName: arpFileName,
		settings:    settings,
		wans:        NewStringSet(),
	}

	return scanner
}

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
					Ip:  ip,
					Mac: mac,
				},
			},
		})
}

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
