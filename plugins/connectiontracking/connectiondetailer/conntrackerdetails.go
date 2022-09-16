package connectiondetailer

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"

	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// Structs used to read in the XML generated by the conntrack command
type ConntrackConnectionDetails struct {
	XMLName xml.Name `xml:"conntrack"`
	Flows   []Flow   `xml:"flow"`
}

type Flow struct {
	XMLName xml.Name `xml:"flow"`
	Metas   []Meta   `xml:"meta"`
}

type Meta struct {
	XMLName   xml.Name `xml:"meta"`
	Direction string   `xml:"direction,attr"`

	// Direction is original or reply
	LayerThree ConntrackLayerThree `xml:"layer3"`
	LayerFour  ConntrackLayerFour  `xml:"layer4"`

	// Direction is independent
	Timeout int32 `xml:"timeout"`
	Mark    int64 `xml:"mark"`
	Use     int32 `xml:"use"`
	Id      int64 `xml:"id"`
}

type ConntrackLayerThree struct {
	XMLName   xml.Name `xml:"layer3"`
	Protonum  int32    `xml:"protonum,attr"`
	Protoname string   `xml:"protoname,attr"`
	Src       string   `xml:"src"`
	Dst       string   `xml:"dst"`
}

type ConntrackLayerFour struct {
	XMLName   xml.Name `xml:"layer4"`
	Protonum  int32    `xml:"protonum,attr"`
	Protoname string   `xml:"protoname,attr"`
	SPort     int32    `xml:"sport"`
	DPort     int32    `xml:"dport"`
}

type ConnTrackerDetails struct {
	connectionsXml []byte
}

// Unit tests don't run as root, but the conntrack command requires it.
// To make this code testable, separate retrieving the connection XML
// from getting the ConnectionInfo list
// FetchSystemConnections() runs the conntrack command to retrieve its XML output
// Returns an error if something went wrong while running the command
func (connTrackerDetails *ConnTrackerDetails) FetchSystemConnections() error {
	var retError error

	// run conntrack command
	cmd := exec.Command("conntrack", "--dump", "--output", "extended", "--output", "xml")
	connTrackerDetails.connectionsXml, retError = cmd.CombinedOutput()

	return retError
}

// Sets connectionsXml, the XML file containing all system connections. Only used for testing
func (connTrackerDetails *ConnTrackerDetails) SetConnectionsXml(connectionXml []byte) {
	connTrackerDetails.connectionsXml = connectionXml
}

func (connTrackerDetails *ConnTrackerDetails) GetConnectionsXml() []byte {
	return connTrackerDetails.connectionsXml
}

// Gets the list of connections on the system. Make sure to run FetchSystemConnections
// before running GetConnectionList()
func (connTrackerDetails *ConnTrackerDetails) GetDeviceToConnections() (map[string][]*Discoverd.Connection, error) {
	if connTrackerDetails.connectionsXml == nil {
		return nil, errors.New("ConnTrackerDetails requires that FetchSystemConnections is run before GetConnectionList()")
	}

	if connections, err := connTrackerDetails.getConnections(); err == nil {
		return getDeviceToConnections(connections), nil
	} else {
		return nil, err
	}

}

// Get the list of connections from the conntrack command
func (connTrackerDetails *ConnTrackerDetails) getConnections() ([]*Discoverd.Connection, error) {
	// Unmarshal XML output of conntrack command and get it into a useful data structure
	connTracker, err := parseConntrackXml(connTrackerDetails.connectionsXml)

	// Fail early if the XML provided by the conntrack command could not be read
	if err != nil {
		return nil, fmt.Errorf("conntracker: failed to run conntrack command: %w", err)
	}

	//connections := make([]*Discoverd.ConnectionTracking, len(connTracker.Flows))
	var connections []*Discoverd.Connection

	// XML structure is pretty awkward, pull out it's data and put it in a more friendly data structure
	for _, flow := range connTracker.Flows {
		connectionInfo := new(Discoverd.Connection)
		for _, meta := range flow.Metas {

			if meta.Direction == "independent" {
				connectionInfo.Independent = new(Discoverd.Independent)
				connectionInfo.Independent.Timeout = meta.Timeout
				connectionInfo.Independent.Mark = meta.Mark
				connectionInfo.Independent.Use = meta.Use
				connectionInfo.Independent.Id = meta.Id

			} else if meta.Direction == "reply" {
				connectionInfo.Reply = new(Discoverd.Reply)

				connectionInfo.Reply.LayerThree = new(Discoverd.LayerThree)
				connectionInfo.Reply.LayerThree.Protonum = meta.LayerThree.Protonum
				connectionInfo.Reply.LayerThree.Protoname = meta.LayerThree.Protoname
				connectionInfo.Reply.LayerThree.Src = meta.LayerThree.Src
				connectionInfo.Reply.LayerThree.Dst = meta.LayerThree.Dst

				connectionInfo.Reply.LayerFour = new(Discoverd.LayerFour)
				connectionInfo.Reply.LayerFour.Protonum = meta.LayerFour.Protonum
				connectionInfo.Reply.LayerFour.Protoname = meta.LayerFour.Protoname
				connectionInfo.Reply.LayerFour.SPort = meta.LayerFour.SPort
				connectionInfo.Reply.LayerFour.DPort = meta.LayerFour.DPort
			} else if meta.Direction == "original" {
				connectionInfo.Original = new(Discoverd.Original)

				connectionInfo.Original.LayerThree = new(Discoverd.LayerThree)
				connectionInfo.Original.LayerThree.Protonum = meta.LayerThree.Protonum
				connectionInfo.Original.LayerThree.Protoname = meta.LayerThree.Protoname
				connectionInfo.Original.LayerThree.Src = meta.LayerThree.Src
				connectionInfo.Original.LayerThree.Dst = meta.LayerThree.Dst

				connectionInfo.Original.LayerFour = new(Discoverd.LayerFour)
				connectionInfo.Original.LayerFour.Protonum = meta.LayerFour.Protonum
				connectionInfo.Original.LayerFour.Protoname = meta.LayerFour.Protoname
				connectionInfo.Original.LayerFour.SPort = meta.LayerFour.SPort
				connectionInfo.Original.LayerFour.DPort = meta.LayerFour.DPort
			}
		}
		connections = append(connections, connectionInfo)
	}

	return connections, nil
}

// From the list of connections, build a map of a device's IPv4 address to it's connections
func getDeviceToConnections(connections []*Discoverd.Connection) map[string][]*Discoverd.Connection {
	deviceToConnections := make(map[string][]*Discoverd.Connection)

	for _, connection := range connections {
		originalIp := connection.Original.LayerThree.Src
		replyIp := connection.Reply.LayerThree.Src

		// Check that the connection has IPv4s to be used as a key
		if originalIp != "" {
			deviceToConnections[originalIp] = append(deviceToConnections[originalIp], connection)
		}

		if replyIp != "" {
			deviceToConnections[replyIp] = append(deviceToConnections[replyIp], connection)
		}
	}

	return deviceToConnections
}

// Parses XML provided by the conntrackXml byte slice
func parseConntrackXml(conntrackXml []byte) (*ConntrackConnectionDetails, error) {
	var connTrackerDetails *ConntrackConnectionDetails
	err := xml.Unmarshal(conntrackXml, &connTrackerDetails)

	return connTrackerDetails, err
}
