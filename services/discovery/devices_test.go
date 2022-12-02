package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test the DataTracker object.
func TestDataUse(t *testing.T) {
	binInterval := time.Second / 4
	removalInterval := time.Second
	dataTracker := NewDataTracker(binInterval, removalInterval)
	now := time.Now()

	// Add an interval that is too old, to test that we remove it.
	dataTracker.dataUseIntervals = []DataUse{
		{
			Start:   now.Add(-2 * time.Second),
			RxBytes: 100,
			TxBytes: 100,
		},
	}
	dataTracker.IncrRx(1)
	assert.Len(t, dataTracker.dataUseIntervals, 1)
	dataTracker.IncrTx(1)
	assert.Len(t, dataTracker.dataUseIntervals, 1)
	assert.EqualValues(t, dataTracker.TotalUse().Total(),
		2)
	// Force a new bin opening.
	time.Sleep(binInterval * 3 / 2)
	dataTracker.IncrTx(1)
	assert.EqualValues(t, dataTracker.TotalUse().Total(), 3)
	assert.Len(t, dataTracker.dataUseIntervals, 2)

	dataTracker.RestrictTrackerToInterval(binInterval)
	assert.EqualValues(t, dataTracker.TotalUse(), DataUseAmount{Tx: 1, Rx: 0})
	// If we restrict it to a super short interval, it has to
	// delete all bins, make sure there is not panic.
	dataTracker.RestrictTrackerToInterval(time.Microsecond)
	assert.EqualValues(t, dataTracker.TotalUse(), DataUseAmount{Tx: 0, Rx: 0})
}
