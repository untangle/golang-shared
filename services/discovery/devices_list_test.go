package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/untangle/golang-shared/structs/protocolbuffers/ActiveSessions"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

type DeviceListTestSuite struct {
	suite.Suite

	mac1         string
	mac2         string
	mac3         string
	mac4         string
	mac5         string
	devicesTable map[string]*DeviceEntry
	now          time.Time
	oneHourAgo   time.Time
	halfHourAgo  time.Time
	aDayago      time.Time

	deviceList *DevicesList
}

// SetupTest just populates the device table.
func (suite *DeviceListTestSuite) SetupTest() {
	suite.now = time.Now()
	suite.oneHourAgo = suite.now.Add(-1 * time.Hour)
	suite.halfHourAgo = suite.now.Add(-30 * time.Minute)
	suite.aDayago = suite.now.Add(-24 * time.Hour)
	suite.mac1 = "00:11:22:33:44:55"
	suite.mac2 = "00:11:22:33:44:66"
	suite.mac3 = "00:11:33:44:55:66"
	suite.mac4 = "00:aa:bb:cc:dd:ee"
	suite.mac5 = "00:aa:bb:cc:dd:ef"

	suite.devicesTable = map[string]*DeviceEntry{
		suite.mac1: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac1,
			LastUpdate: suite.oneHourAgo.Unix(),
			Neigh:      []*disco.NEIGH{{Ip: "192.168.50.4"}, {Ip: "192.168.50.5"}},
		}},
		suite.mac2: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac2,
			LastUpdate: suite.halfHourAgo.Unix(),
			Nmap:       []*disco.NMAP{{Hostname: "TestHostname0"}, {Hostname: "TestHostname1"}},
		}},
		suite.mac3: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac3,
			LastUpdate: suite.halfHourAgo.Unix(),
			Neigh:      []*disco.NEIGH{{Ip: "192.168.53.4"}, {Ip: "192.168.53.5"}},
		}},
		suite.mac4: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac4,
			LastUpdate: suite.aDayago.Unix(),
			Lldp:       []*disco.LLDP{{SysName: "Sysname0"}, {SysName: "Sysname1"}},
		}},
		suite.mac5: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac5,
			LastUpdate: suite.aDayago.Unix(),
			Neigh:      []*disco.NEIGH{{Ip: "192.168.55.4"}, {Ip: "192.168.55.5"}},
		}},
	}
	suite.deviceList = NewDevicesList()
	for _, entry := range suite.devicesTable {
		suite.deviceList.PutDevice(entry)
	}
}

// Test if MergeOrAddDeviceEntry can merge a DeviceEntry with just an IP address
// if a DeviceEntry is already present in the DeviceList with a matching IP and a MAC Address
func (suite *DeviceListTestSuite) TestMergeByIp() {
	expectedSysName := "testIfThisFieldPresent"

	suite.deviceList.MergeOrAddDeviceEntry(&DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{
		Neigh: []*disco.NEIGH{{Ip: suite.devicesTable[suite.mac3].Neigh[0].Ip}},
		Lldp:  []*disco.LLDP{{SysName: expectedSysName}},
	}}, func() {})

	dev := suite.deviceList.listDevices(func(d *DeviceEntry) bool {
		if len(d.Lldp) > 0 && d.Lldp[0].SysName == expectedSysName {
			return true
		}
		return false
	})

	suite.True(len(dev) > 0)
}

// TestListing tests that applying various predicates to the list
// works.
func (suite *DeviceListTestSuite) TestListing() {

	getMacs := func(devs []*DeviceEntry) (output []string) {
		output = []string{}
		for _, dev := range devs {
			output = append(output, dev.MacAddress)
		}
		return
	}

	type testSpec struct {
		predicates       []ListPredicate
		expectedListMacs []string
		description      string
	}
	tests := []testSpec{
		{
			predicates:       []ListPredicate{},
			expectedListMacs: []string{suite.mac1, suite.mac2, suite.mac3, suite.mac4, suite.mac5},
			description:      "Empty predicate list should return all MACs.",
		},
		{
			predicates:       []ListPredicate{WithUpdatesWithinDuration(time.Hour - time.Second)},
			expectedListMacs: []string{suite.mac2, suite.mac3},
			description:      "List of macs with an update time within the hour should be mac2, mac3.",
		},
		{
			predicates: []ListPredicate{
				WithUpdatesWithinDuration(time.Minute * 20),
			},
			expectedListMacs: []string{},
			description:      "There are no devices in the list that have an update time within 20 minutes and are not local.",
		},
		{
			predicates: []ListPredicate{
				LastUpdateOlderThanDuration(time.Hour * 24),
			},
			expectedListMacs: []string{suite.mac4, suite.mac5},
			description:      "List of macs with an update time older than 24 hours should be mac4",
		},
	}
	// duplicate the tests for concurrent testing.
	testsDuplicated := append(tests, tests...)
	testsDuplicated = append(testsDuplicated, tests...)

	// clone a device list, which will be given to us by
	// ApplyToDeviceList. We can't just return a copy because the
	// entries are locked.
	cloneDevs := func(inputs []*DeviceEntry) (interface{}, error) {
		output := []*DeviceEntry{}
		for _, dev := range inputs {
			new := &DeviceEntry{}
			proto.Merge(
				&new.DiscoveryEntry,
				&dev.DiscoveryEntry)
			output = append(output, new)
		}
		return output, nil
	}

	wg := sync.WaitGroup{}
	for i := range testsDuplicated {
		wg.Add(1)
		// localTest is in it's own scope and can be captured
		// in a close reliably, but the test loop variable is
		// in an outer scope and repeatedly overwritten.
		localTest := &testsDuplicated[i]
		go func() {
			theListIface, err := suite.deviceList.ApplyToDeviceList(cloneDevs, localTest.predicates...)
			suite.Nil(err)
			theList := theListIface.([]*DeviceEntry)
			assert.ElementsMatch(suite.T(), localTest.expectedListMacs, getMacs(theList), localTest.description)
			// then re-put all devices for more assurance that this works, creating race conditions.
			for _, entry := range theList {
				suite.deviceList.PutDevice(entry)
			}
			wg.Done()
		}()

	}
	wg.Wait()

}

// TestMerge tests the merge and add funcionality. We test that
// concurrent merges are okay and that we actually do merge undefined
// or newer fields.
func (suite *DeviceListTestSuite) TestMerge() {
	type newAndOldPair struct {
		old         *DeviceEntry
		new         *DeviceEntry
		expectedMac string
		expectedIps []string
	}
	deviceTests := []newAndOldPair{}
	i := 0

	for mac, v := range suite.devicesTable {
		newIp := fmt.Sprintf("192.168.1.%d", i)
		newEntry := &DeviceEntry{
			DiscoveryEntry: disco.DiscoveryEntry{
				// The mac address has to be a deep copy to not cause a race condition
				MacAddress: mac,
				LastUpdate: v.LastUpdate + 1,
				Lldp:       v.Lldp,
				Neigh:      []*disco.NEIGH{{Ip: newIp}},
				Nmap:       v.Nmap,
			},
		}
		testSpec := newAndOldPair{
			old:         v,
			new:         newEntry,
			expectedMac: v.MacAddress,
			expectedIps: append(v.getDeviceIpsUnsafe(), newIp),
		}

		i++
		deviceTests = append(deviceTests, testSpec)
	}
	suite.deviceList.Lock.Unlock()

	// Do the merge multiple times for each device. The invariants should stay the
	// same.
	deviceTests = append(deviceTests, deviceTests...)
	wg := sync.WaitGroup{}
	var count uint32
	for _, devTest := range deviceTests {
		wg.Add(1)
		go func(pair newAndOldPair) {
			defer wg.Done()
			suite.deviceList.MergeOrAddDeviceEntry(pair.new,
				func() {
					// assert that we put this entry in the table.
					mapEntry := suite.deviceList.Devices[pair.expectedMac]
					suite.Same(mapEntry, pair.new)
					suite.Equal(mapEntry.LastUpdate, pair.new.LastUpdate)
					suite.ElementsMatch(mapEntry.getDeviceIpsUnsafe(), pair.expectedIps)
					suite.Equal(mapEntry.MacAddress, pair.expectedMac)
					atomic.AddUint32(&count, 1)
				})

		}(devTest)
	}
	wg.Wait()
	suite.Equal(count, uint32(len(deviceTests)))
}

// TestBroadcastInsertion tests that we do not add in a broadcast entry.
func (suite *DeviceListTestSuite) TestBroadcastInsertion() {

	var deviceList DevicesList
	var count uint32
	deviceList.Devices = map[string]*DeviceEntry{}

	for mac, device := range suite.devicesTable {
		deviceList.Devices[mac] = device
		count++
	}

	newBroadcastDiscovery := DeviceEntry{
		DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: "00:00:00:00:00:00",
			LastUpdate: suite.halfHourAgo.Unix(),
		},
	}

	deviceList.MergeOrAddDeviceEntry(&newBroadcastDiscovery,
		func() {
		})

	// Asssert that broadcast entry was not added.
	suite.EqualValues(count, len(deviceList.Devices), "Adding broadcast discovery entry.")
	suite.Equal(suite.devicesTable, deviceList.Devices, "Adding broadcast discovery entry.")
}

// Test clean device entry
func (suite *DeviceListTestSuite) TestCleanDeviceEntry() {

	var deviceList DevicesList
	var device_count, device_ip_count uint32
	deviceList.Devices = map[string]*DeviceEntry{}
	deviceList.devicesByIP = map[string]*DeviceEntry{}

	for _, entry := range suite.devicesTable {
		deviceList.PutDevice(entry)

		device_count++

	}

	//device_ip_count is six since there are six IPs present in the device entries
	device_ip_count = 6

	//Clean entries which are 48 hours older
	predicates1 := []ListPredicate{
		LastUpdateOlderThanDuration(time.Hour * 48),
	}
	deviceList.CleanOldDeviceEntry(predicates1...)
	//No entries should be deleted from the list, because no entries with LastUpdate older than 48 hours
	suite.EqualValues(device_count, len(deviceList.Devices), "Cleaned 48 hour and older entry")
	suite.EqualValues(device_ip_count, len(deviceList.devicesByIP), "Cleaned 48 hour and older entry")

	//Clean entries which are 24 hours older
	predicates2 := []ListPredicate{
		LastUpdateOlderThanDuration(time.Hour * 24),
	}
	deviceList.CleanOldDeviceEntry(predicates2...)
	//two device entries should be deleted from the device list which entry has LastUpdate with >24 hours
	suite.EqualValues(device_count-2, len(deviceList.Devices), "Cleaned 24 hours older entry")
	suite.EqualValues(device_ip_count-2, len(deviceList.devicesByIP), "Cleaned 24 hours older entry")

}

// TestMarshallingList tests that we can marshal a list of devices
// obtained via the ApplyToDeviceList function to JSON without getting
// an exception.
func (suite *DeviceListTestSuite) TestMarshallingList() {
	output, err := suite.deviceList.ApplyToDeviceList(
		func(list []*DeviceEntry) (interface{}, error) {
			return json.Marshal(list)
		})
	suite.Nil(err)
	bytes := output.([]byte)
	fmt.Printf("JSON string output: %s\n", string(bytes))
}

func (suite *DeviceListTestSuite) TestMergeSessions() {
	sessions := []*ActiveSessions.Session{
		{
			Bytes:         10,
			ClientBytes:   0,
			ServerBytes:   10,
			ClientAddress: suite.devicesTable[suite.mac1].getDeviceIpsUnsafe()[0],
		},
		{
			Bytes:         10,
			ClientBytes:   0,
			ServerBytes:   10,
			ClientAddress: suite.devicesTable[suite.mac3].getDeviceIpsUnsafe()[1],
		},
		{
			Bytes:       10,
			ClientBytes: 0,
			ServerBytes: 10,
		},
	}

	suite.deviceList.MergeSessions(sessions)
	dev := suite.deviceList.Devices[suite.mac1]
	suite.Equal(len(dev.sessions), 1)

	suite.Equal(1, len(suite.deviceList.Devices[suite.mac3].sessions))
}

func (suite *DeviceListTestSuite) TestDeviceMarshal() {
	macaddr := "00:11:22:33:44:55:66"
	update := 123456
	dev := DeviceEntry{
		DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: macaddr,
			LastUpdate: int64(update),
		},
		sessions: []*ActiveSessions.Session{
			{
				ByteRate: 11,
				Bytes:    22,
			},
			{
				ByteRate: 13,
				Bytes:    33,
			},
		},
	}
	output, err := json.Marshal(&dev)
	suite.Nil(err)
	dictMarshal := map[string]interface{}{
		"macAddress": macaddr,
		"LastUpdate": json.Number(fmt.Sprintf("%d", update)),
		"sessionDetail": map[string]interface{}{
			"byteTransferRate": json.Number("24"),
			"dataUsage":        json.Number("55"),
			"numSessions":      json.Number("2"),
		},
	}
	testDict := map[string]interface{}{}
	decoder := json.NewDecoder(bytes.NewBuffer(output))
	decoder.UseNumber()
	suite.Nil(decoder.Decode(&testDict))
	suite.Equal(dictMarshal, testDict)
}

// TestGetDevFromIP just tests that GetDeviceEntryFromIP functions as
// indended.
func (suite *DeviceListTestSuite) TestGetDevFromIP() {
	dev := suite.devicesTable[suite.mac1]

	for _, ip := range dev.getDeviceIpsUnsafe() {
		foundDev := suite.deviceList.GetDeviceEntryFromIP(ip)
		suite.True(proto.Equal(dev, foundDev))
	}
}

func TestDeviceList(t *testing.T) {
	testSuite := &DeviceListTestSuite{}
	suite.Run(t, testSuite)
}
