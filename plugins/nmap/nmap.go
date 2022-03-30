package nmap

import (
	"encoding/xml"
	"os/exec"
	"strconv"

	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type Nmap struct {
	Nmaprun xml.Name `xml:"nmaprun"`
	Host    []Host   `xml:"host"`
}

type Host struct {
	XMLName   xml.Name  `xml:"host"`
	Status    Status    `xml:"status"`
	Address   []Address `xml:"address"`
	Hostnames Hostnames `xml:"hostnames"`
	Ports     Ports     `xml:"ports"`
	Os        Os        `xml:"os"`
	Uptime    Uptime    `xml:"uptime,omitempty"`
}

type Status struct {
	XMLName xml.Name `xml:"status"`
	State   string   `xml:"state,attr"` // e.g. up, down
}

type Hostnames struct {
	XMLName xml.Name `xml:"hostnames"`
	Hostname []Hostname `xml:"hostname"`
}

type Hostname struct {
	Name     string   `xml:"name,attr"` // e.g. home.lan
}

type Address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"` // e.g. 192.168.101.1, 60:38:E0:D7:53:4B
	AddrType string   `xml:"addrtype,attr"` // e.g. ipv4, mac
	Vendor   string   `xml:"vendor,attr"` // e.g. Belkin International
}

type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Port   []Port   `xml:"port"`
}

type Port struct {
	XMLName  xml.Name `xml:"port"`
	Protocol string   `xml:"protocol,attr"` // e.g. tcp
	PortId   string   `xml:"portid,attr"` // e.g. 53
	State    State    `xml:"state"`
	Service  Service  `xml:"service"`
}

type State struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr"` // e.g. open
}

type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"` // e.g. ssh, http
	Method  string   `xml:"method,attr"` // e.g. table
}

type Os struct {
	XMLName xml.Name  `xml:"os"`
	OsMatch []OsMatch `xml:"osmatch"`
}

type OsMatch struct {
	XMLName xml.Name `xml:"osmatch"`
	Name    string   `xml:"name,attr"` // e.g. Windows, Linux
}


type Uptime struct {
	XMLName  xml.Name `xml:"uptime"`
	Seconds  string   `xml:"seconds,attr"` // e.g. 1212665
	LastBoot string   `xml:"lastboot,attr"` // e.g. Wed Mar 16 09:43:39 2022
}

// Start starts the NMAP collector
func Start() {
	logger.Info("Starting NMAP collector plugin\n")
	discovery.RegisterCollector(NmapcallBackHandler)

	// initial run
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
	// TODO: pass the box IP/prefix subnet to be scanned
	cmd := exec.Command("nmap", "-sT", "-O", "-F", "-oX",  "-", "192.168.101.0/24")
	output, _ := cmd.CombinedOutput()

	// parse xml output data
	var nmap Nmap
	if err := xml.Unmarshal([]byte(output), &nmap); err != nil {
		logger.Err("Unable to unmarshal xml: %s\n", err)
	}

	// iterate hosts
	for _, host := range nmap.Host {
		var mac string
		var macVendor string
		var ip string

		// skip if host is not up
		if (host.Status.State != "up") {
			continue
		}

		// initialize the discovery entry
		entry := discovery.DeviceEntry{}
		entry.Init()
		entry.Data.Nmap = &Discoverd.NMAP{}

		// iterate addresses to find mac
		for _, address := range host.Address {
			if (address.AddrType == "mac") {
				mac = address.Addr
				macVendor = address.Vendor
			}
			if (address.AddrType == "ipv4") {
				ip = address.Addr
			}
		}

		logger.Info("--- nmap discovery ---\n")

		if (mac != "") {
			logger.Info("> MAC: %s, Vendor: %s\n", mac, macVendor)
			entry.Data.Nmap.MacVendor = macVendor
		} else {
			logger.Info("> MAC: n/a\n")
		}

		logger.Info("> IPv4: %s\n", ip)

 	  // hostname
		if (len(host.Hostnames.Hostname) > 0) {
			var hostname = host.Hostnames.Hostname[0].Name
			logger.Info("> Hostname: %s\n", hostname)
			entry.Data.Nmap.Hostname = hostname
		} else {
			logger.Info("> Hostname: n/a\n")
		}

		// os
		if (len(host.Os.OsMatch) > 0) {
			var osname = host.Os.OsMatch[0].Name
			logger.Info("> OS: %s\n", osname)
			entry.Data.Nmap.Os = osname
		} else {
			logger.Info("> OS: n/a\n")
		}

		// uptime
		if (host.Uptime.Seconds != "") {
			logger.Info("> Uptime: %s seconds\n", host.Uptime.Seconds)
			logger.Info("> Last boot: %s\n", host.Uptime.LastBoot)
			entry.Data.Nmap.Uptime = host.Uptime.Seconds
			entry.Data.Nmap.LastBoot = host.Uptime.LastBoot
		} else {
			logger.Info("> Uptime: n/a\n")
			logger.Info("> Last boot: n/a\n")
		}

		// ports
		if (len(host.Ports.Port) > 0) {
			var portInfo string
			for _, port := range host.Ports.Port {
				// lookup only open ports
				if (port.State.State == "open") {
					portNo, _ := strconv.Atoi(port.PortId)
					if (portNo > 0) {
						nmapPort := &Discoverd.NMAPPorts{}
						nmapPort.Port = int32(portNo)
						nmapPort.Protocol = port.Service.Name

						entry.Data.Nmap.OpenPorts = append(entry.Data.Nmap.OpenPorts, nmapPort)

						portInfo += port.PortId + "(" + port.Service.Name + ") "
					}
				}
			}
			logger.Info("> Open Ports: %s\n", portInfo)
		} else {
			logger.Info("> Open Ports: n/a\n")
		}

		// update entry if mac exists
		if (mac != "") {
			discovery.UpdateDiscoveryEntry(mac, entry)
		}
	}
}

