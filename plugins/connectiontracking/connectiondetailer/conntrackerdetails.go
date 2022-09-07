package connectiondetailer

import (
	"encoding/xml"
	"os/exec"
)

type ConnTrackerDetails struct {
	connectionDetails *ConnectionDetails
}

func (connTrackerDetails *ConnTrackerDetails) InitializeConnectionDetails() error {
	// run conntrack command
	cmd := exec.Command("conntrack", "--dump", "--output", "extended", "--output", "xml")
	conntrackXml, _ := cmd.CombinedOutput()

	var err error
	connTrackerDetails.connectionDetails, err = parseConntrackXml(conntrackXml)

	return err
}

// Parses XML provided by the conntrackXml byte slice
func parseConntrackXml(conntrackXml []byte) (*ConnectionDetails, error) {
	var connTrackerDetails *ConnectionDetails
	err := xml.Unmarshal(conntrackXml, &connTrackerDetails)

	return connTrackerDetails, err
}

func (connTrackerDetails *ConnTrackerDetails) GetConnectionDetails() *ConnectionDetails {
	return connTrackerDetails.connectionDetails
}
