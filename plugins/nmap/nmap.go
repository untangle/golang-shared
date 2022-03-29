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
	XMLName xml.Name  `xml:"host"`
	Status  Status    `xml:"status"`
	Address []Address `xml:"address"`
	Ports   Ports     `xml:"ports"`
}

type Status struct {
	XMLName xml.Name `xml:"status"`
	State   string   `xml:"state,attr"` // e.g. up, down
}

type Address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"` // e.g. 192.168.101.1, 60:38:E0:D7:53:4B
	AddrType string   `xml:"addrtype,attr"` // e.g. ipv4, mac
	Vendor   string   `xml:"vendor,attr"` // e.g. Belkin International
}

type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Ports   []Port   `xml:"port"`
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
	// -F = fast mode (fewer ports)
	// -oX = output XML
	// TODO: pass the box IP/prefix subnet to be scanned
	cmd := exec.Command("nmap", "-F", "-oX",  "-", "192.168.101.0/24")
	output, _ := cmd.CombinedOutput()

	// parse xml output data
	var nmap Nmap
	if err := xml.Unmarshal([]byte(output), &nmap); err != nil {
		logger.Err("Unable to unmarshal xml: %s\n", err)
	}

	// iterate hosts
	for _, host := range nmap.Host {
		var mac string
		var ip string
		var info string

		// skip if host is not up
		if (host.Status.State != "up") {
			continue
		}

		// iterate addresses to find mac
		for _, address := range host.Address {
			if (address.AddrType == "mac") {
				mac = address.Addr
			}
			if (address.AddrType == "ipv4") {
				ip = address.Addr
			}
		}

		info += "Found host " + ip + " (" + mac + "), open ports: "

		// initialize the discovery entry
		entry := discovery.DeviceEntry{}
		entry.Init()
		entry.Data.Nmap = &Discoverd.NMAP{}

		// iterate ports
		for _, port := range host.Ports.Ports {
			// lookup only open ports
			if (port.State.State == "open") {
				portNo, _ := strconv.Atoi(port.PortId)
				if (portNo > 0) {
					nmapPort := &Discoverd.NMAPPorts{}
					nmapPort.Port = int32(portNo)
					nmapPort.Protocol = port.Service.Name

					entry.Data.Nmap.OpenPorts = append(entry.Data.Nmap.OpenPorts, nmapPort)

					info += port.PortId + "(" + port.Service.Name + ") "
				}
			}
		}

		// update entry if mac exists
		if (mac != "") {
			discovery.UpdateDiscoveryEntry(mac, entry)
		}

		logger.Info("%s\n", info)
	}
}

