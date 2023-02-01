package discovery

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/alerts"
	"github.com/untangle/golang-shared/structs/protocolbuffers/ActiveSessions"
	protoAlerts "github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

// DeviceEntry represents a device found via discovery and methods on
// it. Mostly used as a key of DevicesList.
type DeviceEntry struct {
	disco.DiscoveryEntry
	sessions    []*ActiveSessions.Session
	dataTracker *DataTracker
}

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
		list.putDevicesByIpUnsafe(ip, entry)
	}
}

func (list *DevicesList) putDevicesByIpUnsafe(ip string, entry *DeviceEntry) {
	if !entry.HasIp(ip) {
		return
	}

	// If the IP is also assigned to another entry, it means
	// that one of the entries IP got changed/removed and reassigned.
	// We need to keep the IP on the entry with the latest LastUpdate.
	ipEntry, ok := list.devicesByIP[ip]
	if !ok || ipEntry.MacAddress == entry.MacAddress {
		list.devicesByIP[ip] = entry
		return
	}

	if entry.LastUpdate > ipEntry.LastUpdate {
		list.cleanEntryIpUnsafe(ip, ipEntry)
		list.devicesByIP[ip] = entry
		return
	}

	list.cleanEntryIpUnsafe(ip, entry)
}

func (list *DevicesList) cleanEntryIpUnsafe(ip string, entry *DeviceEntry) {
	entry.RemoveIp(ip)
	if entry.IsEmpty() {
		delete(list.Devices, entry.MacAddress)
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

// transformDevices returns a list of devices after running each
// element of the list through the transformers. If a transformer
// returns nil, that element is not included in the output list or run
// through transformers after.  This doesn't do anything with locks so
// without the outer function locking appropriately it is unsafe.
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

// PredsToTransformers transforms a list of predicates to transformers
// which return the device pointer if the pred is true else nil. Can
// be used with transformDevices to filter the device list.
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

// ApplyToTransformedList 'transforms' a list by applying the
// transformers in sequence to each device in the list, and calling
// doToList on the result. If any transformer returns nil on a
// particular device, then that device isn't included in the slice
// passed to doToList, allowing you to filter and transform at the
// same time.
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

// MergeOrAddDeviceEntrySilent processes existing entries read from DB reportd
// it will NOT create any new device alerts
func (list *DevicesList) MergeOrAddDeviceEntrySilent(entry *DeviceEntry, callback func()) {
	list.mergeOrAdd(entry, callback)
}

// MergeOrAddDeviceEntry processes entries found by discoverd
// it will create an alert if a new device is discovered
func (list *DevicesList) MergeOrAddDeviceEntry(entry *DeviceEntry, callback func()) {
	list.mergeOrAddWithAlert(entry, callback, alerts.Publisher())
}

// mergeOrAdd merges the new entry if an entry can be found
// that corresponds to the same MAC or IP. If none can be found, we
// put the new entry in the table. The provided callback function is
// called after everything is merged but before the lock is
// released. This can allow you to clone/copy the merged device.
// Make sure to merge new into old.
// returns a pointer to the inserted device (new or existing)
// returns true if the device is a new one or false if it is an existing one
func (list *DevicesList) mergeOrAdd(entry *DeviceEntry, callback func()) (
	processedEntry *DeviceEntry,
	isNewDevice bool,
) {
	// Lock the entry down before reading from it.
	// Otherwise the read in Merge causes a data race
	if entry.MacAddress == "00:00:00:00:00:00" {
		return nil, false
	}

	list.Lock.Lock()
	defer list.Lock.Unlock()

	deviceIps := entry.GetDeviceIPs()
	if entry.MacAddress == "" && len(deviceIps) <= 0 {
		return nil, false
	}

	// deferred functions are called LIFO which means this will be called BEFORE the mutex unlock
	defer callback()

	if oldEntry, ok := list.Devices[entry.MacAddress]; ok {
		oldEntry.Merge(entry)
		list.putDeviceUnsafe(oldEntry)
		return oldEntry, false
	}

	if len(deviceIps) > 0 {
		// See if the IPs of entry correspond to any others
		for _, ip := range deviceIps {
			// Once an old entry is oldEntry and the new entry is merged with it,
			// break out of the loop since any device oldEntry is a pointer that
			// every IP for a device points to
			// We do not merge the new entry into the old entry if the new entry is a new device.
			oldEntry := list.getDeviceFromIPUnsafe(ip)
			if oldEntry != nil && (oldEntry.MacAddress == entry.MacAddress || entry.MacAddress == "" || oldEntry.MacAddress == "") {
				oldEntry.Merge(entry)
				list.putDeviceUnsafe(oldEntry)
				return oldEntry, false
			}
		}

		list.putDeviceUnsafe(entry)
		return entry, true
	}

	list.putDeviceUnsafe(entry)
	return entry, true
}

// mergeOrAddWithAlert cals mergeOrAdd and creates a "new device discovered" alert when necessary
func (list *DevicesList) mergeOrAddWithAlert(entry *DeviceEntry, callback func(), alertsPublisher alerts.AlertPublisher) {
	processedEntry, isNewDevice := list.mergeOrAdd(entry, callback)

	if isNewDevice {
		alertsPublisher.Send(buildNewDeviceDiscoveredAlert(processedEntry))
	}
}

// buildAlert builds the alert object that will be sent to alertd and adds device details
func buildNewDeviceDiscoveredAlert(entry *DeviceEntry) *protoAlerts.Alert {
	deviceIps := entry.GetDeviceIPs()
	return &protoAlerts.Alert{
		Type:     protoAlerts.AlertType_DISCOVERY,
		Severity: protoAlerts.AlertSeverity_INFO,
		Message:  "ALERT_NEW_DEVICE_DISCOVERED",
		Params: map[string]string{
			"ips":        strings.Join(deviceIps, ","),
			"macAddress": entry.MacAddress,
		},
	}
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

func (n *DeviceEntry) HasIp(ip string) bool {
	_, inNamp := n.Nmap[ip]
	_, inNeigh := n.Neigh[ip]
	_, inLldp := n.Lldp[ip]

	return inNamp || inNeigh || inLldp
}

func (n *DeviceEntry) IsEmpty() bool {
	return len(n.Lldp) == 0 && len(n.Neigh) == 0 && len(n.Nmap) == 0
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

func (n *DeviceEntry) RemoveIp(ip string) {
	delete(n.Lldp, ip)
	delete(n.Neigh, ip)
	delete(n.Nmap, ip)
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

// getDataTracker returns the data tracker or creates a new one if it
// doesn't exist, using the default bin interval and track duration.
func (n *DeviceEntry) getDataTracker() *DataTracker {
	if n.dataTracker == nil {
		n.dataTracker = NewDataTracker(
			defaultBinInterval,
			defaultTrackDuration)
	}
	return n.dataTracker
}

// IncrData increments the data use amount for this device, using the
// member data tracker instance.
func (n *DeviceEntry) IncrData(incr DataUseAmount) {
	n.getDataTracker().IncrData(incr)
}

// GetDataUse gets total data use of this device as kept by the
// tracker (which may be trimmed).
func (n *DeviceEntry) GetDataUse() DataUseAmount {
	return n.getDataTracker().TotalUse()
}

// SetMac sets the mac address of the device entry. It 'normalizes' it.
func (n *DeviceEntry) SetMac(mac string) {
	n.MacAddress = strings.ToLower(mac)

}
