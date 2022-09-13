package connectiondetailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// The conntrack command must be run as root. Unit tests don't run as root,
//  so read in captured xml data from the conntrack command
func TestGetConnectionList(t *testing.T) {
	// Create ConnectionTracking structs to compare GetConnectionList() output to

	connectionOne := &Discoverd.ConnectionTracking{
		Independent: &Discoverd.Independent{
			Id:      4245641327,
			Mark:    0,
			Timeout: 367,
			Use:     1,
		},
		Reply: &Discoverd.Reply{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "224.0.0.251",
				Dst:       "192.168.0.117",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "unknown",
				Protonum:  2,
			},
		},
		Original: &Discoverd.Original{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "192.168.0.117",
				Dst:       "224.0.0.251",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "unknown",
				Protonum:  2,
			},
		},
	}

	connectionTwo := &Discoverd.ConnectionTracking{
		Independent: &Discoverd.Independent{
			Id:      3523745495,
			Mark:    0,
			Timeout: 118,
			Use:     1,
		},
		Reply: &Discoverd.Reply{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "142.250.217.78",
				Dst:       "192.168.0.246",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "udp",
				Protonum:  17,
				SPort:     443,
				DPort:     53629,
			},
		},
		Original: &Discoverd.Original{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "192.168.0.246",
				Dst:       "142.250.217.78",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "udp",
				Protonum:  17,
				SPort:     53629,
				DPort:     443,
			},
		},
	}

	connectionThree := &Discoverd.ConnectionTracking{
		Independent: &Discoverd.Independent{
			Id:      2944662851,
			Mark:    0,
			Timeout: 28,
			Use:     1,
		},
		Reply: &Discoverd.Reply{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "142.251.33.99",
				Dst:       "192.168.0.246",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "udp",
				Protonum:  17,
				SPort:     443,
				DPort:     53063,
			},
		},
		Original: &Discoverd.Original{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "192.168.0.246",
				Dst:       "142.251.33.99",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "udp",
				Protonum:  17,
				SPort:     53063,
				DPort:     443,
			},
		},
	}

	exepectedConnections := []*Discoverd.ConnectionTracking{connectionOne, connectionTwo, connectionThree}

	// read in captured conntrack XML output from file
	connTracker := new(ConnTrackerDetails)

	xmlBytes, err := os.ReadFile("./testdata/conntrack_output.xml")

	assert.Nil(t, err, "The test data could not be read in")

	connTracker.SetConnectionsXml(xmlBytes)

	connections, _ := connTracker.GetConnectionList()

	for i, connection := range connections {
		assert.Equal(t, exepectedConnections[i].Independent.Id, connection.Independent.Id)
		assert.Equal(t, exepectedConnections[i].Independent.Mark, connection.Independent.Mark)
		assert.Equal(t, exepectedConnections[i].Independent.Timeout, connection.Independent.Timeout)
		assert.Equal(t, exepectedConnections[i].Independent.Use, connection.Independent.Use)

		assert.Equal(t, exepectedConnections[i].Reply.LayerThree.Protoname, connection.Reply.LayerThree.Protoname)
		assert.Equal(t, exepectedConnections[i].Reply.LayerThree.Protonum, connection.Reply.LayerThree.Protonum)
		assert.Equal(t, exepectedConnections[i].Reply.LayerThree.Src, connection.Reply.LayerThree.Src)
		assert.Equal(t, exepectedConnections[i].Reply.LayerThree.Dst, connection.Reply.LayerThree.Dst)

		assert.Equal(t, exepectedConnections[i].Reply.LayerFour.Protoname, connection.Reply.LayerFour.Protoname)
		assert.Equal(t, exepectedConnections[i].Reply.LayerFour.Protonum, connection.Reply.LayerFour.Protonum)
		assert.Equal(t, exepectedConnections[i].Reply.LayerFour.DPort, connection.Reply.LayerFour.DPort)
		assert.Equal(t, exepectedConnections[i].Reply.LayerFour.SPort, connection.Reply.LayerFour.SPort)

		assert.Equal(t, exepectedConnections[i].Original.LayerThree.Protoname, connection.Original.LayerThree.Protoname)
		assert.Equal(t, exepectedConnections[i].Original.LayerThree.Protonum, connection.Original.LayerThree.Protonum)
		assert.Equal(t, exepectedConnections[i].Original.LayerThree.Src, connection.Original.LayerThree.Src)
		assert.Equal(t, exepectedConnections[i].Original.LayerThree.Dst, connection.Original.LayerThree.Dst)

		assert.Equal(t, exepectedConnections[i].Original.LayerFour.Protoname, connection.Original.LayerFour.Protoname)
		assert.Equal(t, exepectedConnections[i].Original.LayerFour.Protonum, connection.Original.LayerFour.Protonum)
		assert.Equal(t, exepectedConnections[i].Original.LayerFour.DPort, connection.Original.LayerFour.DPort)
		assert.Equal(t, exepectedConnections[i].Original.LayerFour.SPort, connection.Original.LayerFour.SPort)

	}
}
