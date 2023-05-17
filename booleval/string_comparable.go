package booleval

import (
	"fmt"
	"net"
)

type StringComparable struct {
	theString string
}

func NewStringComparable(val string) StringComparable {
	return StringComparable{val}
}

func (s StringComparable) Greater(other any) (bool, error) {
	switch val := other.(type) {
	case string:
		return s.theString > val, nil

	}
	return false, fmt.Errorf(
		"booleval: StringComparable: string does not support ordering with respect to %v(%T)",
		other, other)
}

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
