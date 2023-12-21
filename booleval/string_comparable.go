package booleval

import (
	"fmt"
	"net"
)

// StringComparable is a Comparable for strings.
type StringComparable struct {
	theString string
}

// NewStringComparable returns a Comparable for a string value.
func NewStringComparable(val string) StringComparable {
	return StringComparable{val}
}

// Greater returns true only if other is a string, and i is
// lexicographically greater than other.
func (s StringComparable) Greater(other any) (bool, error) {
	switch val := other.(type) {
	case string:
		return s.theString > val, nil

	}
	return false, fmt.Errorf(
		"booleval: StringComparable: string does not support ordering with respect to %v(%T)",
		other, other)
}

// Equal returns true if other is a string and the two strings
// i.theString and other are equal, or if other is and integer or IP
// address, and its string representation is equal to s.theString.
func (s StringComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case string:
		return s.theString == val, nil
	case uint32, uint64, uint, uint8, uint16, int, int32, int64,
		bool, float32, float64, net.IP, net.IPNet:
		return s.theString == fmt.Sprintf("%v", val), nil
	}
	return false, nil
}

var _ Comparable = StringComparable{}

// NewStringArrayComparable returns a Comparable for an array of strings.
func NewStringArrayComparable(val []string) ArrayComparable {
	return NewArrayComparable(val)
}
