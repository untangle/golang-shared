package discovery

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/structs/protocolbuffers/ActiveSessions"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

// DataUse is an interval of data use by the device.
type DataUse struct {
	Start   time.Time
	End     time.Time
	RxBytes uint
	TxBytes uint
}

type DataTracker struct {
	dataUseIntervals   []DataUse
	dataUseBinInterval time.Duration
	maxTrackDuration   time.Duration
}

// DeviceEntry represents a device found via discovery and methods on
// it. Mostly used as a key of DevicesList.
type DeviceEntry struct {
	disco.DiscoveryEntry
	sessions    []*ActiveSessions.Session
	dataTracker *DataTracker
}

const defaultBinInterval = 30 * time.Minute
const defaultTrackDuration = 24 * time.Hour

// DevicesList is an in-memory 'list' of all known devices (stored as
// a map from mac address string to device).
type DevicesList struct {
	Devices map[string]*DeviceEntry

	// TODO: Prior to 1.18, net.HardwareAddr and net.IpAddr cannot be used as a map key.
	// Currently we have to be careful about making our Mac and IPV6 string uppercase
	// before using them as a key in devicesByIp and Devices maps to avoid duplicates. Once upgraded
	// to 1.18+, swap these string type keys to net.HardwareAddr for Devices and net.IpAddr for devicesByIp
	devicesByIP map[string]*DeviceEntry
	Lock        sync.RWMutex
}

// NewDevicesList returns a new DevicesList which is ready for use.
func NewDevicesList() *DevicesList {
	return &DevicesList{
		Devices:     map[string]*DeviceEntry{},
		devicesByIP: map[string]*DeviceEntry{},
	}
}

// ListPredicate is a function that when applied to an entry returns
// true if it 'accepts' the entry.
type ListPredicate func(entry *DeviceEntry) bool

// ListElementTransformer is a function that transforms its input in
// some way.
type ListElementTransformer func(entry *DeviceEntry) *DeviceEntry

// WrapPredicateAsTransformer returns a ListElementTransformer that
// returns nil if the predicate fails, else the entry.
func WrapPredicateAsTransformer(pred ListPredicate) ListElementTransformer {
	return func(entry *DeviceEntry) *DeviceEntry {
		if pred(entry) {
			return entry
		}
		return nil
	}
}

// WithUpdatesWithinDuration returns a predicate ensuring that the
// LastUpdate member was within the period. So if period is an hour
// for example, the returned predicate will return true when the
// device has a LastUpdate time within the last hour.
func WithUpdatesWithinDuration(period time.Duration) ListPredicate {
	return func(entry *DeviceEntry) bool {
		lastUpdated := time.Unix(entry.LastUpdate, 0)
		now := time.Now()
		return (now.Sub(lastUpdated) <= period)
	}
}

// TrimToDataUseSince trims the tracked data use to be within the
// specified period.
func TrimToDataUseSince(period time.Duration) ListElementTransformer {
	return func(entry *DeviceEntry) *DeviceEntry {
		entry.getDataTracker().RestrictTrackerToInterval(period)
		return entry
	}
}

// putDeviceUnsafe puts the device in the list without locking it.
func (list *DevicesList) putDeviceUnsafe(entry *DeviceEntry) {
	list.Devices[entry.MacAddress] = entry

	for _, ip := range entry.GetDeviceIPs() {
		list.devicesByIP[ip] = entry
	}
}

func (list *DevicesList) PutDevice(entry *DeviceEntry) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	list.putDeviceUnsafe(entry)
}

// Get 24hours older device discovery entry from device list and
// delete the entry from device list
func (list *DevicesList) CleanOldDeviceEntry(preds ...ListPredicate) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	listOfDevs := list.transformDevices(PredsToTransformers(preds)...)
	list.CleanDevices(listOfDevs)
}

// Clean device discovery entry from devices list if the entry
// lastUpdate is 24 hours older
func (list *DevicesList) CleanDevices(devices []*DeviceEntry) {

	for _, device := range devices {
		delete(list.Devices, device.MacAddress)
		for _, ip := range device.GetDeviceIPs() {
			delete(list.devicesByIP, ip)
		}
		logger.Debug("Deleted entry %s:%s\n", device.MacAddress, device.MacAddress)
	}
}

// Get the device discovery entries lastUpdate time is older than the duration(24 hours)
func LastUpdateOlderThanDuration(period time.Duration) ListPredicate {
	return func(entry *DeviceEntry) bool {
		lastUpdated := time.Unix(entry.LastUpdate, 0)
		now := time.Now()
		return (now.Sub(lastUpdated) >= period)
	}

}

// listDevices returns a list of devices matching all predicates. It
// doesn't do anything with locks so without the outer function
// locking appropriately is unsafe.
func (list *DevicesList) transformDevices(transformers ...ListElementTransformer) (returns []*DeviceEntry) {

	returns = []*DeviceEntry{}
search:
	for _, device := range list.Devices {
		for _, trans := range transformers {
			if device = trans(device); device == nil {
				continue search
			}
		}
		returns = append(returns, device)
	}
	return
}

func PredsToTransformers(preds []ListPredicate) []ListElementTransformer {
	output := make([]ListElementTransformer, 0, len(preds))
	for _, pred := range preds {
		output = append(output, WrapPredicateAsTransformer(pred))
	}
	return output
}

// ApplyToDeviceList applys doToList to the list of device entries in
// all the DeviceList device entries that match all the predicates. It
// does this with the lock taken.
func (list *DevicesList) ApplyToDeviceList(
	doToList func([]*DeviceEntry) (interface{}, error),
	preds ...ListPredicate) (interface{}, error) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	listOfDevs := list.transformDevices(PredsToTransformers(preds)...)
	return doToList(listOfDevs)
}

func (list *DevicesList) ApplyToTransformedList(
	doToList func([]*DeviceEntry) (interface{}, error),
	trans ...ListElementTransformer) (interface{}, error) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	listOfDevs := list.transformDevices(trans...)
	return doToList(listOfDevs)
}

// ApplyToDeviceWithMac will apply doToDev to the device with the
// given mac, if present, and return the result, otherwise it will
// return an error.
func (list *DevicesList) ApplyToDeviceWithMac(
	doToDev func(*DeviceEntry) (interface{}, error),
	mac string) (interface{}, error) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	if dev, wasFound := list.Devices[mac]; wasFound {
		return doToDev(dev)
	}
	return nil, fmt.Errorf("device with mac: %s not found", mac)
}

// getDeviceFromIPUnsafe gets a device in the table by IP address.
func (list *DevicesList) getDeviceFromIPUnsafe(ip string) *DeviceEntry {
	return list.devicesByIP[ip]
}

// GetDeviceEntryFromIP returns a copy of a device entry in the list with IP
// address ip. *Currently only works with ipv4*.
func (list *DevicesList) GetDeviceEntryFromIP(ip string) *disco.DiscoveryEntry {
	list.Lock.RLock()
	defer list.Lock.RUnlock()
	if entry := list.getDeviceFromIPUnsafe(ip); entry != nil {
		cloned := proto.Clone(&entry.DiscoveryEntry)
		return cloned.(*disco.DiscoveryEntry)
	}
	return nil
}

// MergeOrAddDeviceEntry merges the new entry if an entry can be found
// that corresponds to the same MAC or IP. If none can be found, we
// put the new entry in the table. The provided callback function is
// called after everything is merged but before the lock is
// released. This can allow you to clone/copy the merged device.
// Make sure to merge new into old.
func (list *DevicesList) MergeOrAddDeviceEntry(entry *DeviceEntry, callback func()) {
	// Lock the entry down before reading from it.
	// Otherwise the read in Merge causes a data race
	if entry.MacAddress == "00:00:00:00:00:00" {
		return
	}

	list.Lock.Lock()
	defer list.Lock.Unlock()

	deviceIps := entry.GetDeviceIPs()
	if entry.MacAddress == "" && len(deviceIps) <= 0 {
		return
	} else if oldEntry, ok := list.Devices[entry.MacAddress]; ok {
		oldEntry.Merge(entry)
		list.putDeviceUnsafe(oldEntry)
	} else if len(deviceIps) > 0 {
		// See if the IPs of entry correspond to any others
		found := false
		for _, ip := range deviceIps {
			// Once an old entry is oldEntry and the new entry is merged with it,
			// break out of the loop since any device oldEntry is a pointer that
			// every IP for a device points to
			if oldEntry := list.getDeviceFromIPUnsafe(ip); oldEntry != nil {
				oldEntry.Merge(entry)
				list.putDeviceUnsafe(oldEntry)
				found = true
				break
			}
		}
		if !found {
			list.putDeviceUnsafe(entry)
		}
	} else {
		list.putDeviceUnsafe(entry)
	}

	callback()
}

// MergeSessions merges the sessions into the devices.
func (list *DevicesList) MergeSessions(sessions []*ActiveSessions.Session) {
	for _, device := range list.Devices {
		device.sessions = nil
	}
	for _, session := range sessions {
		if entry, ok := list.devicesByIP[session.ClientAddress]; ok {
			entry.sessions = append(entry.sessions, session)
		}
	}
}

// Init initialize a new DeviceEntry
func (n *DeviceEntry) Init() {
	n.MacAddress = ""
	n.Lldp = nil
	n.Neigh = nil
	n.Nmap = nil
}

// GetDeviceIPSet returns the set of IPs as a hashmap.
func (n *DeviceEntry) GetDeviceIPSet() map[string]struct{} {
	// Use a set to easily get the list of unique IPs assigned to a device
	ipSet := make(map[string]struct{})

	for ip := range n.Neigh {
		ipSet[ip] = struct{}{}
	}

	for ip := range n.Lldp {
		ipSet[ip] = struct{}{}
	}

	for ip := range n.Nmap {
		ipSet[ip] = struct{}{}
	}
	return ipSet
}

// GetDeviceIPs returns the list of IPs being used by a device. Does
// not acquire any locks before accessing device list elements. The
// IPs are fetched by going through each collector entry and adding
// any IPs found to a set.
func (n *DeviceEntry) GetDeviceIPs() []string {

	var ipList []string
	ipSet := n.GetDeviceIPSet()
	for ip := range ipSet {
		ipList = append(ipList, ip)
	}

	return ipList
}

// Merge fills the relevant fields of n that are not present with ones
// of newEntry that are.
func (n *DeviceEntry) Merge(newEntry *DeviceEntry) {
	// The protobuf library has a merge function that merges exactly as needed,
	// except for the case where the LastUpdated time coming in is less than
	// The current LastUpdated time. Take a snapshot of the original before merging
	lastUpdate := n.LastUpdate

	proto.Merge(n, newEntry)

	if lastUpdate > n.LastUpdate {
		n.LastUpdate = lastUpdate
	}
}

// SessionDetail is a summary of active session details for a device.
type SessionDetail struct {
	// Total byte transfer rate of all sessions.
	ByteTransferRate int64 `json:"byteTransferRate"`

	// Total number of active sessions on this device.
	NumSessions int64 `json:"numSessions"`

	// Total number of bytes used by a device
	DataUsage int64 `json:"dataUsage"`

	RxTotal int64 `json:"rxTotal"`

	TxTotal int64 `json:"txTotal"`
}

func (n *DeviceEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		*disco.DiscoveryEntry
		SessionDetail SessionDetail `json:"sessionDetail"`
	}{
		DiscoveryEntry: &n.DiscoveryEntry,
		SessionDetail:  n.calcSessionDetails(),
	})
}

func (n *DeviceEntry) calcSessionDetails() (output SessionDetail) {
	for _, session := range n.sessions {
		output.ByteTransferRate += int64(session.ByteRate)
		output.NumSessions++
	}
	data := n.GetDataUse()
	output.RxTotal = int64(data.Rx)
	output.TxTotal = int64(data.Tx)
	output.DataUsage = int64(data.Total())
	return
}

func (n *DeviceEntry) getDataTracker() *DataTracker {
	if n.dataTracker == nil {
		n.dataTracker = NewDataTracker(
			defaultBinInterval,
			defaultTrackDuration)
	}
	return n.dataTracker
}

func (n *DeviceEntry) IncrData(incr DataUseAmount) {
	n.getDataTracker().IncrData(incr)
}

func (n *DeviceEntry) GetDataUse() DataUseAmount {
	return n.getDataTracker().TotalUse()
}

type DataUseAmount struct {
	Tx uint
	Rx uint
}

func (amnt DataUseAmount) Total() uint {
	return amnt.Tx + amnt.Rx
}
func (dataTracker *DataTracker) IncrData(incr DataUseAmount) {
	last := len(dataTracker.dataUseIntervals) - 1
	lastInterval := &dataTracker.dataUseIntervals[last]

	firstInterval := &dataTracker.dataUseIntervals[0]
	if time.Since(firstInterval.Start) > dataTracker.maxTrackDuration {
		dataTracker.RestrictTrackerToInterval(dataTracker.maxTrackDuration)
		dataTracker.IncrData(incr)
		return
	} else if time.Since(lastInterval.Start) > dataTracker.dataUseBinInterval {
		now := time.Now()
		dataTracker.dataUseIntervals = append(
			dataTracker.dataUseIntervals,
			DataUse{
				Start: now,
			})
		lastInterval.End = now
		dataTracker.IncrData(incr)
		return
	}
	lastInterval.RxBytes += incr.Rx
	lastInterval.TxBytes += incr.Tx
}

func NewDataTracker(
	binInterval time.Duration,
	maxInterval time.Duration) *DataTracker {
	return &DataTracker{
		dataUseIntervals: []DataUse{
			{
				Start: time.Now(),
			},
		},
		maxTrackDuration:   maxInterval,
		dataUseBinInterval: binInterval,
	}
}

// IncrTx increments total tx bytes by tx, returns updated total.
func (dataTracker *DataTracker) IncrTx(tx uint) {
	dataTracker.IncrData(DataUseAmount{Tx: tx})
}

// IncrRx increments total rx bytes by rx, returns updated total.
func (dataTracker *DataTracker) IncrRx(rx uint) {
	dataTracker.IncrData(DataUseAmount{Rx: rx})
}

func (dataTracker *DataTracker) DataUseInInterval(before time.Duration) (output DataUseAmount) {
	for i := len(dataTracker.dataUseIntervals) - 1; i >= 0; i-- {
		interval := &dataTracker.dataUseIntervals[i]
		if time.Since(interval.Start) > before {
			break
		}
		output.Rx += interval.RxBytes
		output.Tx += interval.TxBytes
	}
	return
}

func (dataTracker *DataTracker) TotalUse() (output DataUseAmount) {
	for _, i := range dataTracker.dataUseIntervals {
		output.Rx += i.RxBytes
		output.Tx += i.TxBytes
	}
	return
}

func (dataTracker *DataTracker) RestrictTrackerToInterval(before time.Duration) {
	// Find the first entry that is within the interval, and use
	// the slice after that.  This allows us to get a new data
	// tracker with up to that interval.
	begin := sort.Search(
		len(dataTracker.dataUseIntervals),
		func(idx int) bool {
			return time.Since(dataTracker.dataUseIntervals[idx].Start) <= before
		})
	dataTracker.dataUseIntervals = dataTracker.dataUseIntervals[begin:]
	// If there were only old bins, create a new one starting now.
	if len(dataTracker.dataUseIntervals) == 0 {
		dataTracker.dataUseIntervals = append(dataTracker.dataUseIntervals,
			DataUse{Start: time.Now()})
	}
}

// SetMac sets the mac address of the device entry. It 'normalizes' it.
func (n *DeviceEntry) SetMac(mac string) {
	n.MacAddress = strings.ToLower(mac)

}
