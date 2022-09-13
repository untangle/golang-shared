package connectiondetailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The conntrack command must be run as root. Unit tests don't run as root,
//  so read in captured xml data from the conntrack command
func TestGetConnectionList(t *testing.T) {
	// read in captured conntrack XML output from file
	connTracker := new(ConnTrackerDetails)

	xmlBytes, err := os.ReadFile("./testdata/conntrack_output.xml")

	assert.Nil(t, err, "The test data could not be read in")

	connTracker.SetConnectionsXml(xmlBytes)

	connections, _ := connTracker.GetConnectionList()

	// Go through and make sure all the fields were read in
	for _, connection := range connections {
		assert.NotNil(t, connection.Independent.Id)
		assert.NotNil(t, connection.Independent.Mark)
		assert.NotNil(t, connection.Independent.Timeout)
		assert.NotNil(t, connection.Independent.Use)

		assert.NotNil(t, connection.Reply.LayerThree.Protoname)
		assert.NotNil(t, connection.Reply.LayerThree.Protonum)
		assert.NotNil(t, connection.Reply.LayerThree.Src)
		assert.NotNil(t, connection.Reply.LayerThree.Dst)

		assert.NotNil(t, connection.Reply.LayerFour.Protoname)
		assert.NotNil(t, connection.Reply.LayerFour.Protonum)
		assert.NotNil(t, connection.Reply.LayerFour.DPort)
		assert.NotNil(t, connection.Reply.LayerFour.SPort)

		assert.NotNil(t, connection.Original.LayerThree.Protoname)
		assert.NotNil(t, connection.Original.LayerThree.Protonum)
		assert.NotNil(t, connection.Original.LayerThree.Src)
		assert.NotNil(t, connection.Original.LayerThree.Dst)

		assert.NotNil(t, connection.Original.LayerFour.Protoname)
		assert.NotNil(t, connection.Original.LayerFour.Protonum)
		assert.NotNil(t, connection.Original.LayerFour.DPort)
		assert.NotNil(t, connection.Original.LayerFour.SPort)
	}
}
