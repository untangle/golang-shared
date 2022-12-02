package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDataUse(t *testing.T) {
	binInterval := time.Second / 4
	removalInterval := time.Second
	dev := NewDataTracker(binInterval, removalInterval)
	now := time.Now()

	// Add an interval that is too old, to test that we remove it.
	dev.dataUseIntervals = []DataUse{
		{
			Start:   now.Add(-2 * time.Second),
			RxBytes: 100,
			TxBytes: 100,
		},
	}
	dev.IncrRx(1)
	assert.Len(t, dev.dataUseIntervals, 1)
	dev.IncrTx(1)
	assert.Len(t, dev.dataUseIntervals, 1)
	assert.EqualValues(t, dev.TotalUse().Total(),
		2)
	// Force a new bin opening.
	time.Sleep(binInterval * 3 / 2)
	dev.IncrTx(1)
	assert.EqualValues(t, dev.TotalUse().Total(), 3)
	assert.Len(t, dev.dataUseIntervals, 2)

	dev.RestrictTrackerToInterval(binInterval)
	assert.EqualValues(t, dev.TotalUse(), DataUseAmount{Tx: 1, Rx: 0})
	// If we restrict it to a super short interval, it has to
	// delete all bins, make sure there is not panic.
	dev.RestrictTrackerToInterval(time.Microsecond)
	assert.EqualValues(t, dev.TotalUse(), DataUseAmount{Tx: 0, Rx: 0})
}
