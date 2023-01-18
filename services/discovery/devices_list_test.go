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
	suite.mac4 = "00:AA:BB:BB:DD:EE"
	suite.mac5 = "00:AA:BB:CC:DD:EF"

	suite.devicesTable = map[string]*DeviceEntry{
		suite.mac1: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac1,
			LastUpdate: suite.oneHourAgo.Unix(),
			Neigh: map[string]*disco.NEIGH{
				"192.168.50.4": {Ip: "192.168.50.4"},
				"192.168.50.5": {Ip: "192.168.50.5"}},
		},
			dataTracker: NewDataTracker(defaultBinInterval, defaultTrackDuration),
		},
		suite.mac2: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac2,
			LastUpdate: suite.halfHourAgo.Unix(),
			Nmap: map[string]*disco.NMAP{
				"ee90::a00:37ef:feb8:e927": {Hostname: "TestHostname0", Ip: "ee90::a00:37ef:feb8:e927"},
				"ff90::a00:37ff:feb8:e927": {Hostname: "TestHostname1", Ip: "ff90::a00:37ff:feb8:e927"}}},
			dataTracker: NewDataTracker(defaultBinInterval, defaultTrackDuration)},
		suite.mac3: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac3,
			LastUpdate: suite.halfHourAgo.Unix(),
			Neigh: map[string]*disco.NEIGH{
				"192.168.53.4": {Ip: "192.168.53.4", State: "bad"},
				"192.168.53.5": {Ip: "192.168.53.5", State: "bad"},
			}},
			dataTracker: NewDataTracker(defaultBinInterval, defaultTrackDuration),
		},
		suite.mac4: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac4,
			LastUpdate: suite.aDayago.Unix(),
			Lldp: map[string]*disco.LLDP{
				"ee80::a11:37ee:feb0:e927": {SysName: "Sysname0", Ip: "ee80::a11:37ee:feb0:e927"},
				"ee80::a00:37ff:feb0:e927": {SysName: "Sysname1", Ip: "ee80::a00:37ff:feb0:e927"},
			}},
			dataTracker: NewDataTracker(defaultBinInterval, defaultTrackDuration),
		},
		suite.mac5: {DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: suite.mac5,
			LastUpdate: suite.aDayago.Unix(),
			Neigh: map[string]*disco.NEIGH{
				"192.168.55.4": {Ip: "192.168.55.4"},
				"192.168.55.5": {Ip: "192.168.55.5"}}},
			dataTracker: NewDataTracker(defaultBinInterval, defaultTrackDuration)},
	}
	suite.deviceList = NewDevicesList()
	for _, entry := range suite.devicesTable {
		suite.deviceList.PutDevice(entry)
	}
}

// Tests that the collector components of a device list merge properly.
// It's expected that the new device entry collector fields overwrite the old ones
func (suite *DeviceListTestSuite) TestMergeCollectors() {
	// Sysname to check for to verify the new value overwrote the onld
	expectedSysnameLldp := "SysNameNew"
	expectedStateNeigh := "good"
	expectedHostnameNmap := "new"

	// Since maps have a random order, have to hard code the IPs to match up
	lldpMerge := *suite.devicesTable[suite.mac4]
	neighMerge := *suite.devicesTable[suite.mac3]
	nmapMerge := *suite.devicesTable[suite.mac2]

	ipLldpMerge := lldpMerge.GetDeviceIPs()[0]
	ipNeighMerge := neighMerge.GetDeviceIPs()[0]
	ipNmapMerge := nmapMerge.GetDeviceIPs()[0]

	lldpMerge.Lldp = map[string]*disco.LLDP{ipLldpMerge: {SysName: expectedSysnameLldp, Ip: ipLldpMerge}}
	neighMerge.Neigh = map[string]*disco.NEIGH{ipNeighMerge: {State: expectedStateNeigh, Ip: ipNeighMerge}}
	nmapMerge.Nmap = map[string]*disco.NMAP{ipNmapMerge: {Hostname: expectedHostnameNmap, Ip: ipNmapMerge}}

	suite.deviceList.MergeOrAddDeviceEntry(&lldpMerge, func() {}, false)
	suite.deviceList.MergeOrAddDeviceEntry(&neighMerge, func() {}, false)
	suite.deviceList.MergeOrAddDeviceEntry(&nmapMerge, func() {}, false)

	lldpDevice := suite.deviceList.getDeviceFromIPUnsafe(ipLldpMerge)
	neighDevice := suite.deviceList.getDeviceFromIPUnsafe(ipNeighMerge)
	nmapDevice := suite.deviceList.getDeviceFromIPUnsafe(ipNmapMerge)

	suite.Equal(expectedSysnameLldp, lldpDevice.Lldp[ipLldpMerge].SysName)
	suite.Equal(expectedStateNeigh, neighDevice.Neigh[ipNeighMerge].State)
	suite.Equal(expectedHostnameNmap, nmapDevice.Nmap[ipNmapMerge].Hostname)
}

// Test if MergeOrAddDeviceEntry can merge a DeviceEntry with just an IP address
// if a DeviceEntry is already present in the DeviceList with a matching IP and a MAC Address
func (suite *DeviceListTestSuite) TestMergeByIp() {
	// Try to merge in a new LLDP entry containing an expected field to check for
	// If expectedLldpSysname isn't present in a Device Entry, then it wasn't merged
	// properly
	type expected struct {
		mergeIp             string
		expectedLldpSysName string
	}
	macToExpected := make(map[string]*expected)
	expectedLldpIp := "1.1.1.5"

	// Merge new Device Entries
	macToExpected[suite.mac2] = &expected{
		mergeIp:             suite.devicesTable[suite.mac2].Nmap["ff90::a00:37ff:feb8:e927"].Ip,
		expectedLldpSysName: "Ipv6MergeNmap"}
	macToExpected[suite.mac3] = &expected{
		mergeIp:             suite.devicesTable[suite.mac3].Neigh["192.168.53.4"].Ip,
		expectedLldpSysName: "Ipv4Merge"}
	macToExpected[suite.mac4] = &expected{mergeIp: suite.devicesTable[suite.mac4].Lldp["ee80::a00:37ff:feb0:e927"].Ip, expectedLldpSysName: "Ipv6MergeLldp"}

	suite.deviceList.MergeOrAddDeviceEntry(&DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{
		Neigh: map[string]*disco.NEIGH{
			macToExpected[suite.mac3].mergeIp: {Ip: macToExpected[suite.mac3].mergeIp},
		},
		Lldp: map[string]*disco.LLDP{macToExpected[suite.mac3].mergeIp: {SysName: macToExpected[suite.mac3].expectedLldpSysName}},
	}}, func() {}, false)

	suite.deviceList.MergeOrAddDeviceEntry(&DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{
		Lldp: map[string]*disco.LLDP{macToExpected[suite.mac2].mergeIp: {SysName: macToExpected[suite.mac2].expectedLldpSysName, Ip: macToExpected[suite.mac2].mergeIp}},
	}}, func() {}, false)

	suite.deviceList.MergeOrAddDeviceEntry(&DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{
		Nmap: map[string]*disco.NMAP{macToExpected[suite.mac4].mergeIp: {Ip: macToExpected[suite.mac4].mergeIp}},
		Lldp: map[string]*disco.LLDP{expectedLldpIp: {SysName: macToExpected[suite.mac4].expectedLldpSysName}},
	}}, func() {}, false)

	// Retrieve the updated device entries
	ipv4DevEntry := suite.deviceList.GetDeviceEntryFromIP(macToExpected[suite.mac3].mergeIp)
	ipv6NmapDevEntry := suite.deviceList.GetDeviceEntryFromIP(macToExpected[suite.mac2].mergeIp)
	ipv6LldpDevEntry := suite.deviceList.GetDeviceEntryFromIP(macToExpected[suite.mac4].mergeIp)

	// Get all the lldp entries to use a contains assertion later with the expected SysName
	ipv4SysNames := getLldpSysNameList(ipv4DevEntry.Lldp)
	ipv6LldpSysNames := getLldpSysNameList(ipv6LldpDevEntry.Lldp)
	ipv6NmapSysNames := getLldpSysNameList(ipv6NmapDevEntry.Lldp)

	suite.Equal(1,
		getAppearanceCount(macToExpected[suite.mac3].expectedLldpSysName, ipv4SysNames),
		"Mismatch in the occurrence of an expected LLDP SysName after a merge")

	suite.Equal(1,
		getAppearanceCount(macToExpected[suite.mac4].expectedLldpSysName, ipv6LldpSysNames),
		"Mismatch in the occurrence of an expected LLDP SysName after a merge")

	suite.Equal(1,
		getAppearanceCount(macToExpected[suite.mac2].expectedLldpSysName, ipv6NmapSysNames),
		"Mismatch in the occurrence of an expected LLDP SysName after a merge")

}

// Counts the number of times a given string appears in a list.
// Used to check that duplicates don't happen during a merge
func getAppearanceCount(match string, list []string) int {
	count := 0
	for _, element := range list {
		if element == match {
			count += 1
		}
	}
	return count
}

// Gets a list of lldp sysnames for a device
func getLldpSysNameList(lldpEntries map[string]*disco.LLDP) []string {
	var sysNames []string
	for _, lldp := range lldpEntries {
		sysNames = append(sysNames, lldp.SysName)
	}

	return sysNames
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

// TestMerge tests the merge and add functionality. We test that
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
				MacAddress: mac,
				LastUpdate: v.LastUpdate + 1,
				Lldp:       v.Lldp,
				Neigh:      map[string]*disco.NEIGH{newIp: {Ip: newIp}},
				Nmap:       v.Nmap,
			},
		}
		testSpec := newAndOldPair{
			old:         v,
			new:         newEntry,
			expectedMac: mac,
			expectedIps: append(v.GetDeviceIPs(), newIp),
		}

		// Do the merge multiple times for each device. The invariants should stay the
		// same. Make sure the deviceEntry isn't just a pointer to previously created
		// entry or data races will occur.
		testSpecCopy := testSpec
		testSpecCopy.new = &DeviceEntry{
			DiscoveryEntry: disco.DiscoveryEntry{
				MacAddress: mac,
				LastUpdate: v.LastUpdate + 1,
				Lldp:       v.Lldp,
				Neigh:      map[string]*disco.NEIGH{newIp: {Ip: newIp}},
				Nmap:       v.Nmap,
			},
		}

		i++
		deviceTests = append(deviceTests, testSpec, testSpecCopy)
	}

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
					suite.Equal(mapEntry.LastUpdate, pair.new.LastUpdate)
					suite.ElementsMatch(mapEntry.GetDeviceIPs(), pair.expectedIps)
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

	//device_ip_count is ten since there are ten IPs present in the device entries
	device_ip_count = 10

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
	suite.EqualValues(device_ip_count-4, len(deviceList.devicesByIP), "Cleaned 24 hours older entry")
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
			ClientAddress: suite.devicesTable[suite.mac1].GetDeviceIPs()[0],
		},
		{
			Bytes:         10,
			ClientBytes:   0,
			ServerBytes:   10,
			ClientAddress: suite.devicesTable[suite.mac3].GetDeviceIPs()[1],
		},
		{
			Bytes:       10,
			ClientBytes: 0,
			ServerBytes: 10,
		},
		{
			Bytes:         10,
			ClientBytes:   0,
			ServerBytes:   10,
			ClientAddress: suite.devicesTable[suite.mac4].GetDeviceIPs()[0],
		},
	}

	suite.deviceList.MergeSessions(sessions)
	dev := suite.deviceList.Devices[suite.mac1]
	suite.Equal(len(dev.sessions), 1)

	suite.Equal(1, len(suite.deviceList.Devices[suite.mac3].sessions))

	suite.Equal(1, len(suite.deviceList.Devices[suite.mac4].sessions))
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
			"numSessions":      json.Number("2"),
			"rxTotal":          json.Number("0"),
			"txTotal":          json.Number("0"),
			"dataUsage":        json.Number("0"),
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

	for _, ip := range dev.GetDeviceIPs() {
		foundDev := suite.deviceList.GetDeviceEntryFromIP(ip)
		suite.True(proto.Equal(dev, foundDev))
	}
}

// Test the ApplyToDeviceWithMac function that looks up a mac address
// and allows you to do something with it, if it exists.
func (suite *DeviceListTestSuite) TestApplyToMac() {

	wg := sync.WaitGroup{}
	const nUpdates = 100
	const nNoOps = 20
	noOpsLeft := nNoOps
	macNoexist := "00:00:00:00:00:00"
	originalValue := suite.devicesTable[suite.mac1].LastUpdate
	// In order to test that the concurrent part is done right, we
	// launch many goroutines that do concurrent modifications.
	for i := 0; i < nUpdates+nNoOps; i++ {
		wg.Add(1)
		if i%2 == 0 && noOpsLeft > 0 {
			// Do a 'no op' update -- one where nothing
			// should happen because the device doesn't
			// exist.
			go func() {
				called := false
				val, err := suite.deviceList.ApplyToDeviceWithMac(
					func(entry *DeviceEntry) (interface{}, error) {
						called = true
						entry.LastUpdate = entry.LastUpdate + 1
						return entry.LastUpdate, nil
					},
					macNoexist)
				suite.NotNil(err)
				suite.Nil(val)
				suite.False(called)
				wg.Done()
			}()
			noOpsLeft--
		} else {
			// Update the device's LastUpdate field.
			go func() {
				var oldLastUpdate int64
				val, err := suite.deviceList.ApplyToDeviceWithMac(
					func(entry *DeviceEntry) (interface{}, error) {
						oldLastUpdate = suite.devicesTable[suite.mac1].LastUpdate
						entry.LastUpdate = entry.LastUpdate + 1
						return entry.LastUpdate, nil
					},
					suite.mac1)
				suite.Nil(err)
				suite.EqualValues(val, oldLastUpdate+1)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	suite.EqualValues(originalValue+nUpdates, suite.devicesTable[suite.mac1].LastUpdate)
}

func (suite *DeviceListTestSuite) TestTransformList() {

	_, _ = suite.deviceList.ApplyToTransformedList(
		func(devs []*DeviceEntry) (interface{}, error) {
			suite.Len(devs, 0)
			return nil, nil
		},
		func(deviceEntry *DeviceEntry) *DeviceEntry {
			return nil
		})

	for _, entry := range suite.devicesTable {
		entry.dataTracker.dataUseIntervals = []DataUse{
			{
				Start:   suite.oneHourAgo,
				End:     suite.halfHourAgo,
				RxBytes: 100,
				TxBytes: 100,
			},
		}
	}
	suite.devicesTable[suite.mac3].dataTracker.dataUseIntervals =
		append(suite.devicesTable[suite.mac3].dataTracker.dataUseIntervals,
			DataUse{
				Start:   suite.now,
				RxBytes: 33,
				TxBytes: 33,
			})
	_, _ = suite.deviceList.ApplyToTransformedList(
		func(devs []*DeviceEntry) (interface{}, error) {
			suite.Len(devs, len(suite.devicesTable))
			for _, dev := range devs {
				if dev.GetMacAddress() == suite.mac3 {
					suite.EqualValues(dev.GetDataUse().Total(), 66)
				}
			}

			return nil, nil
		},
		TrimToDataUseSince(time.Minute))

}

func (suite *DeviceListTestSuite) TestDhcpOverlay() {
	// Test ips are removed from the old device.
	oldDevice := suite.buildDhcpOverlayEntry("mac_overlay_1", "192.168.11.1", suite.aDayago.Unix())
	oldDevice.DiscoveryEntry.Neigh["192.168.11.2"] = &disco.NEIGH{Ip: "192.168.11.2"}
	newDevice := suite.buildDhcpOverlayEntry("mac_overlay_2", "192.168.11.1", suite.oneHourAgo.Unix())
	suite.deviceList.MergeOrAddDeviceEntry(oldDevice, func() {}, false)
	suite.deviceList.MergeOrAddDeviceEntry(newDevice, func() {}, false)

	ips := newDevice.GetDeviceIPs()
	assert.False(suite.T(), suite.deviceList.Devices[oldDevice.MacAddress].HasIp(ips[0]))
	assert.True(suite.T(), suite.deviceList.Devices[newDevice.MacAddress].HasIp(ips[0]))

	// Test ips are removed from the old device and old device is removed because it has no more ips.
	oldDevice = suite.buildDhcpOverlayEntry("mac_overlay_3", "192.168.12.1", suite.aDayago.Unix())
	newDevice = suite.buildDhcpOverlayEntry("mac_overlay_4", "192.168.12.1", suite.oneHourAgo.Unix())

	suite.deviceList.MergeOrAddDeviceEntry(newDevice, func() {}, false)
	suite.deviceList.MergeOrAddDeviceEntry(oldDevice, func() {}, false)

	ips = newDevice.GetDeviceIPs()
	_, ok := suite.deviceList.Devices[oldDevice.MacAddress]
	assert.False(suite.T(), ok)
	assert.True(suite.T(), suite.deviceList.Devices[newDevice.MacAddress].HasIp(ips[0]))
}

func (suite *DeviceListTestSuite) buildDhcpOverlayEntry(mac, ip string, lastUpdated int64) *DeviceEntry {
	return &DeviceEntry{
		DiscoveryEntry: disco.DiscoveryEntry{
			MacAddress: mac,
			LastUpdate: lastUpdated,
			Neigh: map[string]*disco.NEIGH{
				ip: {Ip: ip},
			},
			Lldp: map[string]*disco.LLDP{
				ip: {
					SysName: "Sys " + ip,
					Ip:      ip,
				},
			},
			Nmap: map[string]*disco.NMAP{
				ip: {
					Hostname: "Sys host " + ip,
					Ip:       ip,
				},
			},
		},
		dataTracker: NewDataTracker(defaultBinInterval, defaultTrackDuration),
	}
}

func TestDeviceList(t *testing.T) {
	testSuite := &DeviceListTestSuite{}
	suite.Run(t, testSuite)
}
