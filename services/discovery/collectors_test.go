package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

func TestNormalizeCollectorEntry(t *testing.T) {
	expectedIp := "2345:0425:2CA1:0000:0000:0567:5673:23B5"
	expectedMac := "FF:FF"

	lldp := &disco.LLDP{Mac: "ff:ff", Ip: "2345:0425:2ca1:0000:0000:0567:5673:23b5"}
	nmap := &disco.NMAP{Mac: "ff:ff", Ip: "2345:0425:2ca1:0000:0000:0567:5673:23b5"}
	neighbor := &disco.NEIGH{Mac: "ff:ff", Ip: "2345:0425:2ca1:0000:0000:0567:5673:23b5"}

	err := NormalizeCollectorEntry(lldp)
	assert.NoError(t, err)

	err = NormalizeCollectorEntry(nmap)
	assert.NoError(t, err)

	err = NormalizeCollectorEntry(neighbor)
	assert.NoError(t, err)

	err = NormalizeCollectorEntry(nil)
	assert.Error(t, err)

	assert.Equal(t, expectedIp, lldp.Ip)
	assert.Equal(t, expectedMac, lldp.Mac)

	assert.Equal(t, expectedIp, nmap.Ip)
	assert.Equal(t, expectedMac, nmap.Mac)

	assert.Equal(t, expectedIp, neighbor.Ip)
	assert.Equal(t, expectedMac, neighbor.Mac)
}

func TestWrapCollectorInDeviceEntry(t *testing.T) {
	testIp := "2345:0425:2ca1:0000:0000:0567:5673:23b5"
	testMac := "ff:ff"

	lldp := &disco.LLDP{Mac: testMac, Ip: testIp}
	nmap := &disco.NMAP{Mac: testMac, Ip: testIp}
	neighbor := &disco.NEIGH{Mac: testMac, Ip: testIp}

	expectedLldp := &DeviceEntry{}
	expectedLldp.MacAddress = testMac
	expectedLldp.Lldp = map[string]*disco.LLDP{testIp: lldp}

	expectedNmap := &DeviceEntry{}
	expectedNmap.MacAddress = testMac
	expectedNmap.Nmap = map[string]*disco.NMAP{testIp: nmap}

	expectedNeighbor := &DeviceEntry{}
	expectedNeighbor.MacAddress = testMac
	expectedNeighbor.Neigh = map[string]*disco.NEIGH{testIp: neighbor}

	actualLldp, err := WrapCollectorInDeviceEntry(lldp)
	assert.NoError(t, err)
	assert.Equal(t, expectedLldp, actualLldp)

	actualNmap, err := WrapCollectorInDeviceEntry(nmap)
	assert.NoError(t, err)
	assert.Equal(t, expectedNmap, actualNmap)

	actualNeighbor, err := WrapCollectorInDeviceEntry(neighbor)
	assert.NoError(t, err)
	assert.Equal(t, expectedNeighbor, actualNeighbor)

	badCollector := &disco.LLDP{Mac: "", Ip: ""}
	actualBad, err := WrapCollectorInDeviceEntry(badCollector)
	assert.Error(t, err)
	assert.Nil(t, actualBad)
}
