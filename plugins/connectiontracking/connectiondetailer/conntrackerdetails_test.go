package connectiondetailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnectionList(t *testing.T) {
	// read in captured conntrack XML output from file
	connTracker := new(ConnTrackerDetails)

	xmlBytes, err := os.ReadFile("./testdata/conntrack_output.xml")

	assert.Nil(t, err, "The test data could not be read in")

	connTracker.SetConnectionsXml(xmlBytes)

	connections, _ := connTracker.GetConnectionList()

	// Go through and make sure all the fields were read in
	for i := 0; i < len(connections); i++ {
		assert.NotNil(t, connections[i].Independent.Id)
		assert.NotNil(t, connections[i].Independent.Mark)
		assert.NotNil(t, connections[i].Independent.Timeout)
		assert.NotNil(t, connections[i].Independent.Use)

		assert.NotNil(t, connections[i].Reply.Layer3.Protoname)
		assert.NotNil(t, connections[i].Reply.Layer3.Protonum)
		assert.NotNil(t, connections[i].Reply.Layer3.Src)
		assert.NotNil(t, connections[i].Reply.Layer3.Dst)

		assert.NotNil(t, connections[i].Reply.Layer3.Protoname)
		assert.NotNil(t, connections[i].Reply.Layer3.Protonum)
		assert.NotNil(t, connections[i].Reply.Layer3.Src)
		assert.NotNil(t, connections[i].Reply.Layer3.Dst)

	}
}
