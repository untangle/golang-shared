package discovery

import (
	"sort"
	"time"
)

// DataUse is an interval of data use by the device.
type DataUse struct {
	// Start of this interval.
	Start time.Time

	// End of this interval.
	End time.Time

	// Bytes rx/tx during this interval.
	RxBytes uint
	TxBytes uint
}

// DataTracker is for tracking data use over time, keeping the data
// use in 'bins' which are for some specified time interval. It also
// ages out/trims bins past a specified duration, the
// maxTrackDuration, given during construction.
type DataTracker struct {
	// slice of intervals.
	dataUseIntervals []DataUse

	// How long each bin should represent. It's possible that
	// End-Start for a bin will exceed this amount if no activity
	// was recorded.
	dataUseBinInterval time.Duration

	// maxTrackDuration -- maximum time to keep track of a data
	// use interval. Bins after this interval will be removed when
	// new bins are added. For example if you have it set to 1
	// hour, and you get a data use 'report' somehow (create a
	// bin), that bin will stay in the object until more than an
	// hour later another bin is added.
	maxTrackDuration time.Duration
}

// default interval that a DataUse is in the DataTracker.
const defaultBinInterval = 30 * time.Minute

// default duration to track data in a data tracker, amounts after this are discarded.
const defaultTrackDuration = 24 * time.Hour

// DataUseAmount is a utility struct for dealing with rx/tx pairs.
type DataUseAmount struct {
	Tx uint
	Rx uint
}

// Total returns the total data use, rx + tx.
func (amnt DataUseAmount) Total() uint {
	return amnt.Tx + amnt.Rx
}

// IncrData adds incr to the data use of the current bin or creates a
// new one if the current has expired, and adds it to that.
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

// NewDataTracker creates a new data tracker that keeps bisn with
// binInterval duration and a max track time of maxInterval.
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

// TotalUse gets total data use for all time we keep track of.
func (dataTracker *DataTracker) TotalUse() (output DataUseAmount) {
	for _, i := range dataTracker.dataUseIntervals {
		output.Rx += i.RxBytes
		output.Tx += i.TxBytes
	}
	return
}

// RestrictTrackerToInterval will trim the intervals to those inside
// the duration.
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
