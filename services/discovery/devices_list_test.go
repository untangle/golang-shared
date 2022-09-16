package discovery

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

type DeviceListTestSuite struct {
	suite.Suite

	mac1         string
	mac2         string
	mac3         string
	mac4         string
	devicesTable *DevicesList
	now          time.Time
	oneHourAgo   time.Time
	halfHourAgo  time.Time
	aDayago      time.Time
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
	suite.devicesTable = &DevicesList{
		Devices: map[string]*DeviceEntry{
			suite.mac1: {disco.DiscoveryEntry{
				MacAddress:  suite.mac1,
				LastUpdate:  suite.oneHourAgo.Unix(),
				IPv4Address: "192.168.56.1",
			}},
			suite.mac2: {disco.DiscoveryEntry{
				MacAddress:  suite.mac2,
				LastUpdate:  suite.halfHourAgo.Unix(),
				IPv4Address: "192.168.56.2",
			}},
			suite.mac3: {disco.DiscoveryEntry{
				MacAddress:  suite.mac3,
				LastUpdate:  suite.halfHourAgo.Unix(),
				IPv4Address: "192.168.56.3",
			}},
			suite.mac4: {disco.DiscoveryEntry{
				MacAddress: suite.mac4,
				LastUpdate: suite.aDayago.Unix(),
			}},
		},
	}
}

//TestListing tests that applying various predicates to the list
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
			theListIface, err := suite.devicesTable.ApplyToDeviceList(cloneDevs, localTest.predicates...)
			suite.Nil(err)
			theList := theListIface.([]*DeviceEntry)
			assert.ElementsMatch(suite.T(), localTest.expectedListMacs, getMacs(theList), localTest.description)
			// then re-put all devices for more assurance that this works, creating race conditions.
			for _, entry := range theList {
				suite.devicesTable.PutDevice(entry)
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
	for _, v := range suite.devicesTable.Devices {
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
			suite.devicesTable.MergeOrAddDeviceEntry(pair.new,
				func() {
					// assert that we put this entry in the table.
					mapEntry := suite.devicesTable.Devices[pair.expectedMac]
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

// TestMarshallingList tests that we can marshal a list of devices
// obtained via the ApplyToDeviceList function to JSON without getting
// an exception.
func (suite *DeviceListTestSuite) TestMarshallingList() {
	output, err := suite.devicesTable.ApplyToDeviceList(
		func(list []*DeviceEntry) (interface{}, error) {
			return json.Marshal(list)
		})
	suite.Nil(err)
	bytes := output.([]byte)
	fmt.Printf("JSON string output: %s\n", string(bytes))
}

// TestGetDevFromIP just tests that GetDeviceEntryFromIP functions as
// indended.
func (suite *DeviceListTestSuite) TestGetDevFromIP() {
	dev := suite.devicesTable.Devices[suite.mac1]
	foundDev := suite.devicesTable.GetDeviceEntryFromIP(dev.IPv4Address)
	suite.True(proto.Equal(dev, foundDev))
}

func TestDeviceList(t *testing.T) {
	testSuite := &DeviceListTestSuite{}
	suite.Run(t, testSuite)
}
