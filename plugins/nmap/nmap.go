package nmap

import (
	"encoding/xml"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
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

var defaultNetwork string
var serviceShutdown = make(chan bool)

var nmapProcesses = make(map[string]nmapProcess)
var nmapProcessesMutex sync.RWMutex = sync.RWMutex{}

// Start starts the NMAP collector
func Start() {
	logger.Info("Starting NMAP collector plugin\n")
	discovery.RegisterCollector(NmapcallBackHandler)

	// Do an initial scan
	NmapcallBackHandler(nil)
}

// Stop stops NMAP collector
func Stop() {
}

// NmapcallBackHandler is the callback handler for the NMAP collector
func NmapcallBackHandler(commands []discovery.Command) {
	logger.Debug("NMap scan handler: Received %d commands\n", len(commands))

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
