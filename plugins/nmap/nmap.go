package nmap

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/untangle/discoverd/plugins/discovery"
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
}

const (
	pluginName   string = "nmap"
	randStartMin int    = 5
	randStartMax int    = 10
)

var (
	nmapSingleton         *Nmap
	once                  sync.Once
	RandStartScanNetTimer *time.Ticker

	defaultNetwork     string
	nmapProcesses                   = make(map[string]nmapProcess)
	nmapProcessesMutex sync.RWMutex = sync.RWMutex{}

	settingsPath []string = []string{"discovery", "plugins"}
)

func init() {
	// Start network scan at random interval between randStartMin to randStartMax
	// to avoid network load during packetd startup
	randStartTime := rand.Intn(randStartMax-randStartMin) + randStartMin
	RandStartScanNetTimer = time.NewTicker(time.Duration(randStartTime) * time.Minute)
}

type nmapPluginType struct {
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	AutoInterval uint   `json:"autoInterval"`
}

// Setup the Nmap struct as a singleton
type Nmap struct {
	autoNmapCollectionChan chan bool
	nmapSettings           nmapPluginType
}

// Gets a singleton instance of the Nmap plugin
func NewNmap() *Nmap {
	once.Do(func() {
		nmapSingleton = &Nmap{autoNmapCollectionChan: make(chan bool)}
	})

	return nmapSingleton
}

func (nmap *Nmap) InSync(settings interface{}) bool {
	newSettings, ok := settings.(nmapPluginType)
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

func (nmap *Nmap) GetSettingsStruct() (interface{}, error) {
	var fileSettings []nmapPluginType
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

func (nmap *Nmap) Name() string {
	return pluginName
}

func (nmap *Nmap) SyncSettings(settings interface{}) error {
	originalSettings := nmap.nmapSettings
	newSettings, ok := settings.(nmapPluginType)
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
func (nmap *Nmap) Startup() error {
	logger.Info("Starting NMAP collector plugin\n")

	// Grab the initial settings on startup
	settings, err := nmap.GetSettingsStruct()
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

	discovery.NewDiscovery().UnregisterCollector(pluginName)

	return nil
}

func (nmap *Nmap) startNmap() {
	discovery.NewDiscovery().RegisterCollector(pluginName, NmapcallBackHandler)

	// Lets do a first run to get the initial data
	NmapcallBackHandler(nil)

	nmap.startAutoNmapCollection()
}

func (nmap *Nmap) autoNmapCollection() {
	logger.Debug("Starting automatic collection from NMAP plugin with an interval of %d seconds\n", nmap.nmapSettings.AutoInterval)
	for {
		select {
		case <-nmap.autoNmapCollectionChan:
			logger.Debug("Stopping automatic collection from NMAP plugin\n")
			nmap.autoNmapCollectionChan <- true
			return
		case <-time.After(time.Duration(nmap.nmapSettings.AutoInterval) * time.Second):
			scanLanNetworks()

		case <-RandStartScanNetTimer.C:
			scanLanNetworks()

			RandStartScanNetTimer.Stop()
		}
	}
}

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

func (nmap *Nmap) startAutoNmapCollection() {
	go nmap.autoNmapCollection()
}

func (nmap *Nmap) stopAutoNmapCollection() {
	// The send to kill the AutoNmapCollection goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case nmap.autoNmapCollectionChan <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-nmap.autoNmapCollectionChan:
		logger.Info("Successful shutdown of the automatic NMAP collector\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown automatic NMAP collector. It may never have been started\n")
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
		// Network Scan
		if command.Command == discovery.CmdScanNet {
			for _, network := range command.Arguments {
				args := []string{"nmap", "-sT", "-O", "-F", "-oX", "-", network}
				go runNMAPcmd(args)
			}
		}
		// Host Scan
		if command.Command == discovery.CmdScanHost {
			for _, host := range command.Arguments {
				args := []string{"nmap", "-sT", "-O", "-F", "-oX", "-", host}
				go runNMAPcmd(args)
			}
		}
	}
}

func runNMAPcmd(args []string) {
	if cmdAllreadyRunning(args) {
		logger.Warn("NMap scan already running for %v, running since %s\n", args, time.Unix(nmapProcesses[strings.Join(args, " ")].timeStarted, 0))
		return
	}

	cmd := createCmd(args)
	output, _ := cmd.CombinedOutput()
	removeProcess(args)
	processScan(output)
}

func createCmd(args []string) *exec.Cmd {
	cmd := exec.Command("nmap", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}

func cmdAllreadyRunning(args []string) bool {
	argsStr := strings.Join(args, " ")
	nmapProcessesMutex.Lock()
	defer nmapProcessesMutex.Unlock()
	if _, ok := nmapProcesses[argsStr]; ok {
		return true
	} else {
		nmapProcesses[argsStr] = nmapProcess{time.Now().Unix(), 0}
		return false
	}
}

func removeProcess(args []string) {
	argsStr := strings.Join(args, " ")
	nmapProcessesMutex.Lock()
	delete(nmapProcesses, argsStr)
	nmapProcessesMutex.Unlock()
}

func processScan(output []byte) {
	// parse xml output data
	var nmap nmap
	if err := xml.Unmarshal([]byte(output), &nmap); err != nil {
		logger.Err("Unable to unmarshal xml: %s\n", err)
	}

	// iterate hosts
	for _, host := range nmap.Host {
		var mac string
		var macVendor string
		var ip string

		// skip if host is not up
		if host.Status.State != "up" {
			continue
		}

		// initialize the discovery entry
		entry := &disc.DeviceEntry{}
		entry.Init()
		entry.Nmap = &Discoverd.NMAP{}

		// iterate addresses to find mac
		for _, address := range host.Address {
			if address.AddrType == "mac" {
				mac = address.Addr
				macVendor = address.Vendor
			}
			if address.AddrType == "ipv4" {
				ip = address.Addr
			}
		}

		logger.Debug("--- nmap discovery ---\n")

		if mac != "" {
			logger.Debug("--- nmap discovery ---\n")
			logger.Debug("> MAC: %s, Vendor: %s\n", mac, macVendor)
			entry.Nmap.MacVendor = macVendor
		} else {
			logger.Debug("> MAC: n/a\n")
		}

		logger.Debug("> IPv4: %s\n", ip)

		// hostname
		if len(host.Hostnames.Hostname) > 0 {
			var hostname = host.Hostnames.Hostname[0].Name
			logger.Debug("> Hostname: %s\n", hostname)
			entry.Nmap.Hostname = hostname
		} else {
			logger.Debug("> Hostname: n/a\n")
		}

		// os
		if len(host.Os.OsMatch) > 0 {
			var osname = host.Os.OsMatch[0].Name
			logger.Debug("> OS: %s\n", osname)
			entry.Nmap.Os = osname
		} else {
			logger.Debug("> OS: n/a\n")
		}

		// uptime
		if host.Uptime.Seconds != "" {
			logger.Debug("> Uptime: %s seconds\n", host.Uptime.Seconds)
			logger.Debug("> Last boot: %s\n", host.Uptime.LastBoot)
			entry.Nmap.Uptime = host.Uptime.Seconds
			entry.Nmap.LastBoot = host.Uptime.LastBoot
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

						entry.Nmap.OpenPorts = append(entry.Nmap.OpenPorts, nmapPort)

						portInfo += port.PortID + "(" + port.Service.Name + ") "
					}
				}
			}
			logger.Debug("> Open Ports: %s\n", portInfo)
		} else {
			logger.Debug("> Open Ports: n/a\n")
		}
		entry.MacAddress = mac
		entry.Nmap.LastUpdate = time.Now().Unix()
		discovery.UpdateDiscoveryEntry(mac, entry)
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
