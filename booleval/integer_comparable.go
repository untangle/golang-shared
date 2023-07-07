package booleval

import (
	"fmt"
	"math"
	"strconv"
)

// IntegerComparable is a Comparable for any integer or integer-like
// object.
type IntegerComparable struct {
	theInteger int64
}

// NewIntegerComparableFromIntType returns an IntegerComparable for any integer type.
func NewIntegerComparableFromIntType[T int | int64 | uint | int16 | uint16 | uint32 | int32 | uint64 | int8 | uint8](val T) IntegerComparable {
	newIntVal, _ := getInt(val)
	return IntegerComparable{newIntVal}
}

// NewIntegerComparableFromAny returns an IntegerComparable for
//
// a.) any int type or
// b.) it will parse a string if given one, and try to make that an int.
//
// it returns an error if val was a string an could not be parsed, or
// some other value we don't know what to do with.
func NewIntegerComparableFromAny(val any) (IntegerComparable, error) {
	intVal, err := getInt(val)
	return IntegerComparable{intVal}, err
}

func getInt(other any) (int64, error) {
	// At least with the go1.19 compiler, I can't just put all
	// these into a single case of int64, int, int32..., the
	// compiler complains.
	switch val := other.(type) {
	case int64:
		return val, nil
	case int32:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		if intValue, err := strconv.ParseInt(val, 10, 64); err == nil {
			return intValue, nil
		} else if floatValue, err := strconv.ParseFloat(val, 32); err == nil {
			if math.Floor(floatValue) == floatValue {
				return int64(floatValue), nil
			}
		}
	}
	return 0, fmt.Errorf(
		"booleval getInt(): value: %v(%T) is not an integer and cannot be made one", other, other)
}

// Greater -- tue if other can be coerced to an integer and i is greater than other.
func (i IntegerComparable) Greater(other any) (bool, error) {
	if intval, err := getInt(other); err != nil {
		return false, err
	} else {
		return i.theInteger > intval, nil
	}
}

// Equal -- true of other can be coerced to an integer and is equal to i.
func (i IntegerComparable) Equal(other any) (bool, error) {
	if intval, err := getInt(other); err != nil {
		return false, err
	} else {
		return i.theInteger == intval, nil
	}
}
