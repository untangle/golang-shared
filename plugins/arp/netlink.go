package arp

import (
	"fmt"
	"time"

	"github.com/untangle/discoverd/plugins/discovery"
	"github.com/untangle/discoverd/utils"
	"github.com/untangle/golang-shared/plugins/zmqmsg"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/structs/interfaces"
	disc_pb "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"github.com/vishvananda/netlink"
)

func NetlinkNeighbourCallbackController(commands []discovery.Command) {
	logger.Debug("Arp Callback handler: Received %d commands\n", len(commands))
	scanner, err := newNetlinkScanner(settings.GetSettingsFileSingleton())

	if err != nil {
		logger.Warn("Couldn't initiate netlink scanner: %s\n", err)
		return
	}
	defer scanner.Close()
	entries, err := scanner.getIpNeighbourEntries()
	if err != nil {
		logger.Warn("Couldn't scan ip neigh for devices: %s\n", err)
		return
	}

	logger.Debug("Discovered entries:\n")
	for i, entry := range entries {
		logger.Debug("Entry nr: %d. Data: %+v\n", i, entry)
		entry.LastUpdate = time.Now().Unix()
		discovery.ZmqpublishEntry(entry, zmqmsg.ARPDeviceZMQTopic)
		//discovery.UpdateDiscoveryEntry(entry.MacAddress, entry)
	}
}

type NetlinkHandler interface {
	LinkByName(name string) (netlink.Link, error)
	NeighList(linkIndex, family int) ([]netlink.Neigh, error)
	Delete()
}

// netlinkScanner scans linux host devices related information using the netlink package.
type netlinkScanner struct {
	handler  NetlinkHandler
	settings *settings.SettingsFile
	// helps with testing
	timestamp func() int64

	deviceEntries []*disc.DeviceEntry
	skippedStates map[int]bool
}

// newNetlinkScanner
// SettingsFile are the system settings that contain the definition of network interfaces.
func newNetlinkScanner(settings *settings.SettingsFile) (*netlinkScanner, error) {
	netlinkHandler, err := netlink.NewHandle()
	if err != nil {
		return nil, err
	}

	skippedStatuses := make(map[int]bool)
	skippedStatuses[netlink.NUD_INCOMPLETE] = true

	scanner := &netlinkScanner{
		settings: settings,
		handler:  netlinkHandler,
		timestamp: func() int64 {
			return time.Now().Unix()
		},
	}

	return scanner, nil
}

func (s *netlinkScanner) Close() error {
	s.deviceEntries = nil
	s.handler.Delete()

	return nil
}

// getIpNeighbourEntries returns the list of hosts connected to network devices using the netlink package
// It takes all non WAN network interface devices and searches their hosts (neighbours) using `ip neigh` command (via netlink package).
func (s *netlinkScanner) getIpNeighbourEntries() ([]*disc.DeviceEntry, error) {
	s.deviceEntries = []*disc.DeviceEntry{}

	networkInterfaces := []*interfaces.Interface{}
	if err := s.settings.UnmarshalSettingsAtPath(&networkInterfaces, "network", "interfaces"); err != nil {
		return nil, fmt.Errorf("getIpNeighEntries: couldn't unmarshall network settings: %w", err)
	}

	for _, networkInterface := range networkInterfaces {
		if networkInterface.IsWAN {
			continue
		}

		neighbours, err := s.findNeighbours(networkInterface)
		if err != nil {
			return nil, err
		}

		for _, neighbour := range neighbours {
			s.addNeighbour(neighbour)
		}
	}

	return s.deviceEntries, nil
}

func (s *netlinkScanner) findNeighbours(networkInterface *interfaces.Interface) ([]netlink.Neigh, error) {
	// Get interface details: we need the internal index in order to get the list of neighbours.
	interfaceDetails, err := s.handler.LinkByName(networkInterface.Device)
	if err != nil {
		return nil, fmt.Errorf("findNeighbours: couldn't identify device: %w", err)
	}

	// Search neighbours. This is equalent to calling `ip neigh`.
	neighbours, err := s.handler.NeighList(interfaceDetails.Attrs().Index, 0)
	if err != nil {
		return nil, fmt.Errorf("findNeighbours: unable to find neighbours: %w", err)
	}

	return neighbours, nil
}

// addNeighbour saves a device host (neighbour) to the internal structure.
func (s *netlinkScanner) addNeighbour(neighbour netlink.Neigh) {
	if s.skippNeighbour(neighbour) {
		return
	}

	/*ipv4 := ""
	if neighbour.Family == unix.AF_INET {
		ipv4 = neighbour.IP.String()
	}*/
	if !utils.IsMacAddress(neighbour.HardwareAddr.String()) {
		logger.Warn("Invalid Mac\n", neighbour.HardwareAddr.String())
		return
	}

	s.deviceEntries = append(s.deviceEntries, &disc.DeviceEntry{
		DiscoveryEntry: disc_pb.DiscoveryEntry{
			MacAddress: neighbour.HardwareAddr.String(),
			//IPv4Address: ipv4,
			Arp: []*disc_pb.ARP{{
				Ip:         neighbour.IP.String(),
				Mac:        neighbour.HardwareAddr.String(),
				LastUpdate: s.timestamp(),
				State:      translateHostState(neighbour.State),
			},
			}},
	})
}

func (s *netlinkScanner) skippNeighbour(neighbour netlink.Neigh) bool {
	if skipp, ok := s.skippedStates[neighbour.State]; skipp && ok {
		return true
	}

	return false
}

// translateHostState returns the string equalent for a netlink host status.
func translateHostState(state int) string {
	switch state {
	case netlink.NUD_NONE:
		return "NONE"
	case netlink.NUD_INCOMPLETE:
		return "INCOMPLETE"
	case netlink.NUD_REACHABLE:
		return "REACHABLE"
	case netlink.NUD_STALE:
		return "STALE"
	case netlink.NUD_DELAY:
		return "DELAY"
	case netlink.NUD_PROBE:
		return "PROBE"
	case netlink.NUD_FAILED:
		return "FAILED"
	case netlink.NUD_NOARP:
		return "NOARP"
	case netlink.NUD_PERMANENT:
		return "PERMANENT"
	default:
		return ""
	}
}
