package booleval

import (
	"fmt"
	"time"
)

// TimeComparable is a comparable for time objects.
type TimeComparable struct {
	time time.Time
}

// Equal returns true if other is an integer and, if interpreted as a
// unix timestamp, it is equal to t.time, or if other is a time object
// and is equal to t.time
func (t TimeComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case int, uint32, uint64, int64, uint, int32:
		return t.time.Unix() == val, nil
	case time.Time:
		return t.time.Equal(val), nil
	}
	return false, fmt.Errorf("booleval TimeComparable.Equal: can't convert %v to a time", other)
}

// Greater is true if other is an integer and t.time is greater than its value interpreted as
// a unix timestamp, or if other is a time.Time
// and t.time is greater than it.
func (t TimeComparable) Greater(other any) (bool, error) {
	switch val := other.(type) {
	case int64, uint, uint32, int, int32, uint64:
		// ignore error because that func returns error only
		// for strings.
		int64Val, _ := getInt(val)
		return t.time.Unix() > int64Val, nil
	case time.Time:
		return t.time.After(val), nil
	}
	return false, fmt.Errorf("booleval TimeComparable.Greater: can't convert %v to a time", other)
}
