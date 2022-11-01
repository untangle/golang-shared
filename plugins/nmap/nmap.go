package nmap

import (
	"context"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
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
	interfaces "github.com/untangle/golang-shared/util/net"
)

type nmap struct {
	Nmaprun xml.Name `xml:"nmaprun"`
	Host    []host   `xml:"host"`
}

type host struct {
	XMLName   xml.Name        `xml:"host"`
	Status    status          `xml:"status"`
	Address   []address       `xml:"address"`
	Hostnames hostnames       `xml:"hostnames"`
	Ports     ports           `xml:"ports"`
	Os        operatingSystem `xml:"os"`
	Uptime    uptime          `xml:"uptime,omitempty"`
}

type status struct {
	XMLName xml.Name `xml:"status"`
	State   string   `xml:"state,attr"` // e.g. up, down
}

type hostnames struct {
	XMLName  xml.Name   `xml:"hostnames"`
	Hostname []hostname `xml:"hostname"`
}

type hostname struct {
	Name string `xml:"name,attr"` // e.g. home.lan
}

type address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"`     // e.g. 192.168.101.1, 60:38:E0:D7:53:4B
	AddrType string   `xml:"addrtype,attr"` // e.g. ipv4, mac
	Vendor   string   `xml:"vendor,attr"`   // e.g. Belkin International
}

type ports struct {
	XMLName xml.Name `xml:"ports"`
	Port    []port   `xml:"port"`
}

type port struct {
	XMLName  xml.Name `xml:"port"`
	Protocol string   `xml:"protocol,attr"` // e.g. tcp
	PortID   string   `xml:"portid,attr"`   // e.g. 53
	State    state    `xml:"state"`
	Service  service  `xml:"service"`
}

type state struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr"` // e.g. open
}

type service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`   // e.g. ssh, http
	Method  string   `xml:"method,attr"` // e.g. table
}

type operatingSystem struct {
	XMLName xml.Name  `xml:"os"`
	OsMatch []osMatch `xml:"osmatch"`
}

type osMatch struct {
	XMLName xml.Name `xml:"osmatch"`
	Name    string   `xml:"name,attr"` // e.g. Windows, Linux
}

type uptime struct {
	XMLName  xml.Name `xml:"uptime"`
	Seconds  string   `xml:"seconds,attr"`  // e.g. 1212665
	LastBoot string   `xml:"lastboot,attr"` // e.g. Wed Mar 16 09:43:39 2022
}

type nmapProcess struct {
	timeStarted int64
	pid         int
	retry       int
}

const (
	pluginName      string = "nmap"
	randStartMin    int    = 5
	randStartMax    int    = 10
	NmapScanTimeout        = 30 * time.Second
	NmapScanRetry          = 3
	Ipv4Str                = "IPV4"
	Ipv6Str                = "IPV6"
	InvalidIPStr           = "Invalid IP"
)

var (
	nmapSingleton         *Nmap
	once                  sync.Once
	randStartScanNetTimer *time.Ticker

	defaultNetwork     string
	nmapProcesses                   = make(map[string]nmapProcess)
	nmapProcessesMutex sync.RWMutex = sync.RWMutex{}

	settingsPath []string = []string{"discovery", "plugins"}
	// Nmap requests published to NmapScan routine
	nmapPublisherChannel = make(chan map[string]nmapProcess, 1000)
	serviceShutdown      = make(chan bool)
)

func init() {
	// Start network scan at random interval between randStartMin to randStartMax
	// to avoid network load during packetd startup
	randStartTime := rand.Intn(randStartMax-randStartMin) + randStartMin
	randStartScanNetTimer = time.NewTicker(time.Duration(randStartTime) * time.Minute)

	plugins.GlobalPluginControl().RegisterPlugin(NewNmap)
}

type nmapPluginSettings struct {
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	AutoInterval uint   `json:"autoInterval"`
}

// Setup the Nmap struct as a singleton
type Nmap struct {
	// During shutdown, goroutines are stopped. If they are already
	// stopped, and shutdown gets called again the process could block.
	// Use two channels for a non-blocking shutdown and ack
	autoNmapCollectionShutdown    chan bool
	autoNmapCollectionShutdownAck chan bool

	nmapSettings nmapPluginSettings
}

// Gets a singleton instance of the Nmap plugin
func NewNmap() *Nmap {
	once.Do(func() {
		nmapSingleton = &Nmap{autoNmapCollectionShutdown: make(chan bool),
			autoNmapCollectionShutdownAck: make(chan bool)}
	})

	return nmapSingleton
}

// Returns true if the current settings match the 'new' settings Provided, otherwise false
func (nmap *Nmap) InSync(settings interface{}) bool {
	newSettings, ok := settings.(nmapPluginSettings)
	if !ok {
		logger.Warn("NMAP: Could not compare the settings file provided to the current plugin settings. The settings cannot be updated.")
		return false
	}

	if newSettings == nmap.nmapSettings {
		logger.Debug("Settings remain unchanged for the NMAP plugin\n")
		return true
	}

	logger.Info("Updating NMAP plugin settings\n")
	return false
}

// Returns a struct containing the plugins settings of type nmapPluginSettings
func (nmap *Nmap) GetCurrentSettingsStruct() (interface{}, error) {
	var fileSettings []nmapPluginSettings
	if err := settings.UnmarshalSettingsAtPath(&fileSettings, settingsPath...); err != nil {
		return nil, fmt.Errorf("NMAP: %s", err.Error())
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
func (nmap *Nmap) Name() string {
	return pluginName
}

// Updates the current settings with the settings passed in. If the plugin was already running
// but the settings changed, the plugin is restarted.
// An error is returned if the settings can't be synced
func (nmap *Nmap) SyncSettings(settings interface{}) error {
	originalSettings := nmap.nmapSettings
	newSettings, ok := settings.(nmapPluginSettings)
	if !ok {
		return fmt.Errorf("NMAP: Settings provided were %s but expected %s",
			reflect.TypeOf(settings).String(), reflect.TypeOf(nmap.nmapSettings).String())
	}

	nmap.nmapSettings = newSettings

	// If settings changed but the plugin was previously enabled, restart the plugin
	// for changes to take effect
	var shutdownError error
	if originalSettings.Enabled && nmap.nmapSettings.Enabled {
		shutdownError = nmap.Shutdown()
	}

	if nmap.nmapSettings.Enabled {
		nmap.startNmap()
	} else {
		shutdownError = nmap.Shutdown()
	}

	return shutdownError
}

// Start starts the NMAP collector
// Meant to only be run once
func (nmap *Nmap) Startup() error {
	logger.Info("Starting NMAP collector plugin\n")

	// Grab the initial settings on startup
	settings, err := nmap.GetCurrentSettingsStruct()
	if err != nil {
		return err
	}

	// SyncSettings will start the plugin if it's enabled
	err = nmap.SyncSettings(settings)
	if err != nil {
		return err
	}

	return nil
}

// Stop stops NMAP collector
func (nmap *Nmap) Shutdown() error {
	logger.Info("Stopping NMAP collector plugin\n")

	nmap.stopAutoNmapCollection()

	discovery.NewDiscovery().DeregisterCollector(pluginName)

	return nil
}

// Start method of the plugin. Meant to be used in a restart of the plugin
func (nmap *Nmap) startNmap() {
	discovery.NewDiscovery().RegisterCollector(pluginName, NmapcallBackHandler)

	go NmapScan()
	// Lets do a first run to get the initial data
	NmapcallBackHandler(nil)

	go nmap.autoNmapCollection()
}

// Stop stops NMAP collector
func Stop() {
	serviceShutdown <- true
}

func (nmap *Nmap) autoNmapCollection() {
	logger.Debug("Starting automatic collection from NMAP plugin with an interval of %d seconds\n", nmap.nmapSettings.AutoInterval)
	for {
		select {
		case <-nmap.autoNmapCollectionShutdown:
			logger.Debug("Stopping automatic collection from NMAP plugin\n")
			nmap.autoNmapCollectionShutdownAck <- true
			return
		case <-time.After(time.Duration(nmap.nmapSettings.AutoInterval) * time.Second):
			scanLanNetworks()

		case <-randStartScanNetTimer.C:
			scanLanNetworks()

			randStartScanNetTimer.Stop()
		}
	}
}

// It process nmapProcess map and run nmap scan one by one
func NmapScan() {
	for {
		select {
		case nmapScanReqs := <-nmapPublisherChannel:
			for k, v := range nmapScanReqs {
				if v.retry <= NmapScanRetry {
					output, err := runNampScan(k)
					if err != nil {
						updateRetry(k, v)
						continue
					}
					removeProcess(k)
					processScan(output)
				} else {
					logger.Debug("Nmap scan retry limit exceeded for %v so remove the processing the request  retry count is %v\n", k, v.retry-1)
					removeProcess(k)
				}
			}
		case <-serviceShutdown:
			return
		}
	}
}

// Update the retry value in nmap scan key
func updateRetry(key string, value nmapProcess) {
	nmapProcessesMutex.Lock()
	defer nmapProcessesMutex.Unlock()
	_, exists := nmapProcesses[key]
	if exists {
		nmapProcesses[key] = nmapProcess{value.timeStarted, value.pid, value.retry + 1}
	}

}

// Calls the Nmap callback handler after getting a list of LANs from the settings file
func scanLanNetworks() {
	// Get list of interfaces from settings file
	localIntfs := interfaces.GetInterfaces(func(intf interfaces.Interface) bool {
		return !intf.IsWAN && intf.Enabled && intf.V4StaticAddress != ""
	})

	var localNetworksCidr []string
	for _, intf := range localIntfs {
		localNetworksCidr = append(localNetworksCidr, intf.GetCidrNotation())
	}

	logger.Debug("Scanning LAN networks: %v\n", localNetworksCidr)
	NmapcallBackHandler([]discovery.Command{{Command: discovery.CmdScanHost, Arguments: localNetworksCidr}})
}

// Stops running the back handler automatically
func (nmap *Nmap) stopAutoNmapCollection() {
	// The send to kill the AutoNmapCollection goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case nmap.autoNmapCollectionShutdown <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-nmap.autoNmapCollectionShutdownAck:
		logger.Info("Successful shutdown of the automatic NMAP collector\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown automatic NMAP collector. It may never have been started\n")
	}
}

// Run Nmap scan for host/network
func runNampScan(args string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), NmapScanTimeout)
	defer cancel()
	argSlice := strings.Split(args, " ")
	cmd := exec.CommandContext(ctx, "nmap", argSlice...)
	output, err := cmd.CombinedOutput()
	return output, err
}

// Check host and network addresstype and return address type
func CheckIPAddressType(ip string) (string, error) {
	var IP string = ip
	if strings.Contains(ip, "/") {
		addr, _, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Warn("Invalid network %v\n", ip)
			return "", fmt.Errorf("InvalidIPStr err:%w", err)
		}
		IP = addr.String()
	}
	if net.ParseIP(IP).To4() != nil {
		return Ipv4Str, nil
	} else if net.ParseIP(IP).To16() != nil {
		return Ipv6Str, nil
	} else {
		return "", fmt.Errorf("InvalidIPStr")
	}

}

// NmapcallBackHandler is the callback handler for the NMAP collector
func NmapcallBackHandler(commands []discovery.Command) {
	logger.Debug("NMAP scan handler: Received %d commands\n", len(commands))

	// run nmap subnet scan
	// -sT scan ports
	// -O scan OS
	// -F = fast mode (fewer ports)
	// -oX = output XML

	if commands == nil && defaultNetwork != "" {
		logger.Debug("NMap scan handler: Running default scan on %s\n", defaultNetwork)
		commands = []discovery.Command{
			{
				Command:   discovery.CmdScanNet,
				Arguments: []string{defaultNetwork},
			},
		}
	}

	for _, command := range commands {
		var args []string
		if command.Command != discovery.CmdScanNet && command.Command != discovery.CmdScanHost {
			logger.Debug("It is not a namp request %v\n", command.Command)
			continue
		} else {
			for _, network := range command.Arguments {
				addressType, err := CheckIPAddressType(network)
				if err != nil {
					continue
				}
				logger.Debug("network:%v addressType=%v \n", network, addressType)
				if addressType == Ipv4Str {
					args = []string{"nmap", "-sT", "-O", "-F", "-oX", "-", network}
				} else if addressType == Ipv6Str {
					args = []string{"nmap", "-sT", "-O", "-F", "-6", "-oX", "-", network}
				}
				if isNewNmapReq(args) {
					nmapPublisherChannel <- nmapProcesses
				}
			}
		}
	}
}

// Check and update Nmap request in the map table with nmap command as a key
func isNewNmapReq(args []string) bool {
	argsStr := strings.Join(args, " ")
	nmapProcessesMutex.Lock()
	defer nmapProcessesMutex.Unlock()
	if _, ok := nmapProcesses[argsStr]; ok {
		return false
	} else {
		nmapProcesses[argsStr] = nmapProcess{time.Now().Unix(), 0, 0}
		return true
	}
}

// Remove Nmap scan command from map table if it is compelted
func removeProcess(args string) {
	nmapProcessesMutex.Lock()
	delete(nmapProcesses, args)
	nmapProcessesMutex.Unlock()
}

// Process the nmap scan output and update device discovery entry
func processScan(output []byte) {
	// parse xml output data
	var nmap nmap
	if err := xml.Unmarshal([]byte(output), &nmap); err != nil {
		logger.Err("Unable to unmarshal xml: %s\n", err)
	}

	// iterate hosts
	for _, host := range nmap.Host {

		// skip if host is not up
		if host.Status.State != "up" {
			continue
		}

		// initialize the discovery entry
		entry := &disc.DeviceEntry{}
		entry.Init()
		entry.Nmap = []*Discoverd.NMAP{}
		nmap := &Discoverd.NMAP{}

		// iterate addresses to find mac
		for _, address := range host.Address {
			if address.AddrType == "mac" {
				if address.Addr != "" {
					if !utils.IsMacAddress(address.Addr) {
						continue
					}
					logger.Debug("--- nmap discovery ---\n")
					nmap.Mac = address.Addr
					nmap.MacVendor = address.Vendor
					logger.Debug("> MAC: %s, Vendor: %s\n", nmap.Mac, nmap.MacVendor)
				}
			}
			if address.AddrType == "ipv4" || address.AddrType == "ipv6" {
				if address.Addr == "" {
					continue
				}
				nmap.Ip = address.Addr
				logger.Debug("> IPv4: %s\n", nmap.Ip)

			}

			// hostname
			if len(host.Hostnames.Hostname) > 0 {
				var hostname = host.Hostnames.Hostname[0].Name
				logger.Debug("> Hostname: %s\n", hostname)
				nmap.Hostname = hostname
			} else {
				logger.Debug("> Hostname: n/a\n")
			}

			// os
			if len(host.Os.OsMatch) > 0 {
				var osname = host.Os.OsMatch[0].Name
				logger.Debug("> OS: %s\n", osname)
				nmap.Os = osname
			} else {
				logger.Debug("> OS: n/a\n")
			}

			// uptime
			if host.Uptime.Seconds != "" {
				logger.Debug("> Uptime: %s seconds\n", host.Uptime.Seconds)
				logger.Debug("> Last boot: %s\n", host.Uptime.LastBoot)
				nmap.Uptime = host.Uptime.Seconds
				nmap.LastBoot = host.Uptime.LastBoot
			} else {
				logger.Debug("> Uptime: n/a\n")
				logger.Debug("> Last boot: n/a\n")
			}

			// ports
			if len(host.Ports.Port) > 0 {
				var portInfo string
				for _, port := range host.Ports.Port {
					// lookup only open ports
					if port.State.State == "open" {
						portNo, _ := strconv.Atoi(port.PortID)
						if portNo > 0 {
							nmapPort := &Discoverd.NMAPPorts{}
							nmapPort.Port = int32(portNo)
							nmapPort.Protocol = port.Service.Name

							nmap.OpenPorts = append(nmap.OpenPorts, nmapPort)

							portInfo += port.PortID + "(" + port.Service.Name + ") "
						}
					}
				}
				logger.Debug("> Open Ports: %s\n", portInfo)
			} else {
				logger.Debug("> Open Ports: n/a\n")
			}
			entry.MacAddress = nmap.Mac
			nmap.LastUpdate = time.Now().Unix()
			entry.Nmap = append(entry.Nmap, nmap)
			entry.LastUpdate = time.Now().Unix()
			logger.Info("nmap - %v \n", nmap)

			logger.Info("NMAP entry - %v \n", entry)
			discovery.ZmqpublishEntry(entry, zmqmsg.NMAPDeviceZMQTopic)
			//discovery.UpdateDiscoveryEntry(mac, entry)
		}
	}
}

// SetNetwork, sets the default network for nmap to scan.
func SetNetwork(network string) {
	logger.Debug("Setting network to %s\n", network)
	_, _, err := net.ParseCIDR(network)
	if err != nil {
		logger.Err("Invalid network: %s\n", err)
		return
	}
	defaultNetwork = network
}