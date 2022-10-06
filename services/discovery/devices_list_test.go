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
	suite.devicesTable = map[string]*DeviceEntry{
		suite.mac1: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress:  suite.mac1,
			LastUpdate:  suite.oneHourAgo.Unix(),
			IPv4Address: "192.168.56.1",
		}},
		suite.mac2: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress:  suite.mac2,
			LastUpdate:  suite.halfHourAgo.Unix(),
			IPv4Address: "192.168.56.2",
		}},
		suite.mac3: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress:  suite.mac3,
			LastUpdate:  suite.halfHourAgo.Unix(),
			IPv4Address: "192.168.56.3",
		}},
		suite.mac4: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac4,
			LastUpdate: suite.aDayago.Unix(),
		}},
	}
	suite.deviceList = NewDevicesList()
	for _, entry := range suite.devicesTable {
		suite.deviceList.PutDevice(entry)
	}
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
			expectedListMacs: []string{suite.mac1, suite.mac2, suite.mac3, suite.mac4},
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
			expectedListMacs: []string{suite.mac4},
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
		expectedIP  string
	}
	deviceTests := []newAndOldPair{}
	i := 0
	for _, v := range suite.devicesTable {
		newEntry := &DeviceEntry{
			DiscoveryEntry: disco.DiscoveryEntry{
				IPv4Address: v.IPv4Address,
				MacAddress:  v.MacAddress,
				LastUpdate:  v.LastUpdate + 1,
			},
		}
		testSpec := newAndOldPair{
			old:         v,
			new:         newEntry,
			expectedMac: v.MacAddress,
			expectedIP:  v.IPv4Address,
		}

		if i%2 == 0 {
			suite.createEmptyFieldsForMerge(v, newEntry)
		}
		i++
		deviceTests = append(deviceTests, testSpec)
	}

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
					suite.Equal(mapEntry.IPv4Address, pair.expectedIP)
					suite.Equal(mapEntry.MacAddress, pair.expectedMac)
					atomic.AddUint32(&count, 1)
				})

		}(devTest)
	}
	wg.Wait()
	suite.Equal(count, uint32(len(deviceTests)))
}

// Here we make sure the merge logic is tested by deleting some fields
// to make sure that they get filled in after the merge.
func (suite *DeviceListTestSuite) createEmptyFieldsForMerge(oldDevice *DeviceEntry, newDevice *DeviceEntry) {
	if newDevice.MacAddress != "" {
		newDevice.IPv4Address = ""
	} else if newDevice.IPv4Address != "" {
		oldDevice.IPv4Address = ""
		newDevice.MacAddress = ""
	}
}

// TestBroadcastInsertion tests that we do not add in a broadcast entry.
func (suite *DeviceListTestSuite) TestBroadcastInsertion() {

	var deviceList DevicesList
	var count uint32
	deviceList.Devices = map[string]*DeviceEntry{}

	for _, entry := range suite.devicesTable {
		deviceList.Devices[entry.MacAddress] = &DeviceEntry{
			DiscoveryEntry: disco.DiscoveryEntry{
				IPv4Address: entry.IPv4Address,
				MacAddress:  entry.MacAddress,
				LastUpdate:  entry.LastUpdate,
			},
		}
		count++
	}

	newBroadcastDiscovery := DeviceEntry{
		DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress:  "00:00:00:00:00:00",
			LastUpdate:  suite.halfHourAgo.Unix(),
			IPv4Address: "192.168.56.4",
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
	var count uint32
	deviceList.Devices = map[string]*DeviceEntry{}

	for _, entry := range suite.devicesTable {
		deviceList.Devices[entry.MacAddress] = &DeviceEntry{
			DiscoveryEntry: disco.DiscoveryEntry{
				IPv4Address: entry.IPv4Address,
				MacAddress:  entry.MacAddress,
				LastUpdate:  entry.LastUpdate,
			},
		}
		count++
	}
	//Clean entries which are 48 hours older
	predicates1 := []ListPredicate{
		LastUpdateOlderThanDuration(time.Hour * 48),
	}
	deviceList.CleanOldDeviceEntry(predicates1...)
	//No entries should be deleted from the list, because no entries with LastUpdate older than 48 hours
	suite.EqualValues(count, len(deviceList.Devices), "Cleaned 48 hours older entry")

	//Clean entries which are 24 hours older
	predicates2 := []ListPredicate{
		LastUpdateOlderThanDuration(time.Hour * 24),
	}
	deviceList.CleanOldDeviceEntry(predicates2...)
	//One device entry should be deleted from the device list which entry has LastUpdate with >24 hours
	suite.EqualValues(count-1, len(deviceList.Devices), "Cleaned 24 hours older entry")
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
			ClientAddress: suite.getDevIP(suite.mac1),
		},
		{
			Bytes:         10,
			ClientBytes:   0,
			ServerBytes:   10,
			ClientAddress: suite.getDevIP(suite.mac2),
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

	suite.Equal(len(suite.deviceList.Devices[suite.mac2].sessions), 1)
}

func (suite *DeviceListTestSuite) TestDeviceMarshal() {
	macaddr := "00:11:22:33:44:55:66"
	ipaddr := "192.168.55.22"
	update := 123456
	dev := DeviceEntry{
		DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress:  macaddr,
			IPv4Address: ipaddr,
			LastUpdate:  int64(update),
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
		"macAddress":  macaddr,
		"IPv4Address": ipaddr,
		"LastUpdate":  json.Number(fmt.Sprintf("%d", update)),
		"sessionDetail": map[string]interface{}{
			"byteTransferRate": json.Number("24"),
			"numSessions":      json.Number("2"),
		},
	}
	testDict := map[string]interface{}{}
	decoder := json.NewDecoder(bytes.NewBuffer(output))
	decoder.UseNumber()
	suite.Nil(decoder.Decode(&testDict))
	suite.Equal(dictMarshal, testDict)
}

func (suite *DeviceListTestSuite) getDevIP(mac string) string {
	return suite.deviceList.Devices[mac].IPv4Address
}

// TestGetDevFromIP just tests that GetDeviceEntryFromIP functions as
// indended.
func (suite *DeviceListTestSuite) TestGetDevFromIP() {
	dev := suite.devicesTable[suite.mac1]
	foundDev := suite.deviceList.GetDeviceEntryFromIP(dev.IPv4Address)
	suite.True(proto.Equal(dev, foundDev))
}

func TestDeviceList(t *testing.T) {
	testSuite := &DeviceListTestSuite{}
	suite.Run(t, testSuite)
}
