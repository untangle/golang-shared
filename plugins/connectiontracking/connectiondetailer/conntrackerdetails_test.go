package connectiondetailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type ConntrackerDetailsTestSuite struct {
	suite.Suite

	Conntracker *ConnTrackerDetails

	// List of connections expected to be read in
	ExpectedConnectionList []*Discoverd.ConnectionTracking
}

func TestConntrackerDetailsTestSuite(t *testing.T) {
	suite.Run(t, &ConntrackerDetailsTestSuite{Conntracker: &ConnTrackerDetails{}})
}

// The conntrack command must be run as root. Unit tests don't run as root,
//  so read in captured xml data from the conntrack command
func (suite *ConntrackerDetailsTestSuite) TestgetConnections() {
	// read in captured conntrack XML output from file

	xmlBytes, err := os.ReadFile("./testdata/conntrack_output.xml")

	suite.Nil(err, "The test data could not be read in")

	suite.Conntracker.SetConnectionsXml(xmlBytes)

	connections, _ := suite.Conntracker.getConnections()

	for i, connection := range connections {
		suite.Equal(suite.ExpectedConnectionList[i].Independent.Id, connection.Independent.Id)
		suite.Equal(suite.ExpectedConnectionList[i].Independent.Mark, connection.Independent.Mark)
		suite.Equal(suite.ExpectedConnectionList[i].Independent.Timeout, connection.Independent.Timeout)
		suite.Equal(suite.ExpectedConnectionList[i].Independent.Use, connection.Independent.Use)

		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerThree.Protoname, connection.Reply.LayerThree.Protoname)
		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerThree.Protonum, connection.Reply.LayerThree.Protonum)
		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerThree.Src, connection.Reply.LayerThree.Src)
		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerThree.Dst, connection.Reply.LayerThree.Dst)

		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerFour.Protoname, connection.Reply.LayerFour.Protoname)
		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerFour.Protonum, connection.Reply.LayerFour.Protonum)
		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerFour.DPort, connection.Reply.LayerFour.DPort)
		suite.Equal(suite.ExpectedConnectionList[i].Reply.LayerFour.SPort, connection.Reply.LayerFour.SPort)

		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerThree.Protoname, connection.Original.LayerThree.Protoname)
		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerThree.Protonum, connection.Original.LayerThree.Protonum)
		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerThree.Src, connection.Original.LayerThree.Src)
		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerThree.Dst, connection.Original.LayerThree.Dst)

		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerFour.Protoname, connection.Original.LayerFour.Protoname)
		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerFour.Protonum, connection.Original.LayerFour.Protonum)
		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerFour.DPort, connection.Original.LayerFour.DPort)
		suite.Equal(suite.ExpectedConnectionList[i].Original.LayerFour.SPort, connection.Original.LayerFour.SPort)
	}
}

func (suite *ConntrackerDetailsTestSuite) TestgetDeviceToConnections() {
	actualDeviceToConnections := getDeviceToConnections(suite.ExpectedConnectionList)

	for ip, details := range actualDeviceToConnections {
		// Catch a common problems
		suite.NotEqual(ip, "", "An IP with no value was used as a key!")
		suite.NotNil(details, "An IP(a device) has no connections associated with it!")
	}

	// Just make sure the size of connection list per device is what's expected
	suite.Equal(len(actualDeviceToConnections["192.168.0.117"]), 1)
	suite.Equal(len(actualDeviceToConnections["192.168.0.246"]), 3)
	suite.Equal(len(actualDeviceToConnections["224.0.0.251"]), 1)
	suite.Equal(len(actualDeviceToConnections["142.250.217.78"]), 1)
	suite.Equal(len(actualDeviceToConnections["142.251.33.99"]), 2)

	// Make sure all there aren't any unexpected keys in the map
	suite.Equal(len(actualDeviceToConnections), 5)
}

// Setup expected values and initialize suite class members
func (suite *ConntrackerDetailsTestSuite) SetupTest() {
	// read in captured conntrack XML output from file
	suite.Conntracker = new(ConnTrackerDetails)

	xmlBytes, err := os.ReadFile("./testdata/conntrack_output.xml")

	suite.Nil(err, "The test data could not be read in")

	suite.Conntracker.SetConnectionsXml(xmlBytes)
	// Create list of expected connections
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

	connectionFour := &Discoverd.ConnectionTracking{
		Independent: &Discoverd.Independent{
			Id:      2978983895,
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
				SPort:     11300,
				DPort:     10300,
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
				SPort:     10300,
				DPort:     11300,
			},
		},
	}

	connectionFive := &Discoverd.ConnectionTracking{
		Independent: &Discoverd.Independent{
			Id:      63545631,
			Mark:    0,
			Timeout: 28,
			Use:     1,
		},
		Reply: &Discoverd.Reply{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "",
				Dst:       "",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "udp",
				Protonum:  17,
				SPort:     11300,
				DPort:     10300,
			},
		},
		Original: &Discoverd.Original{
			LayerThree: &Discoverd.LayerThree{
				Protoname: "ipv4",
				Protonum:  2,
				Src:       "",
				Dst:       "",
			},
			LayerFour: &Discoverd.LayerFour{
				Protoname: "udp",
				Protonum:  17,
				SPort:     10300,
				DPort:     11300,
			},
		},
	}

	suite.ExpectedConnectionList = []*Discoverd.ConnectionTracking{connectionOne, connectionTwo, connectionThree, connectionFour, connectionFive}
}
