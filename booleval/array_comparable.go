package booleval

import "net"

// ArrayComparable is a Comparable for arrays of Comparables.
type ArrayComparable struct {
	theThings []Comparable
}

var _ Comparable = ArrayComparable{}

// NewArrayComparable returns a Comparable for a array of value.
func NewArrayComparable(value any) ArrayComparable {
	var comparables []Comparable
	switch val := value.(type) {
	case []string:
		for _, v := range val {
			comparables = append(comparables, NewStringComparable(v))
		}
	// We loose the type if we have more than one case listed.
	case []int:
		for _, v := range val {
			if iComp, err := NewIntegerComparableFromAny(v); err == nil {
				comparables = append(comparables, iComp)
			}
		}
	case []net.IP:
		for _, v := range val {
			if ipComp, err := NewIPOrIPNetComparable(v.String()); err == nil {
				comparables = append(comparables, ipComp)
			}
		}
	case []net.IPNet:
		for _, v := range val {
			if ipComp, err := NewIPOrIPNetComparable(v.String()); err == nil {
				comparables = append(comparables, ipComp)
			}
		}
	}
	return ArrayComparable{comparables}
}

// NewArrayComparableFromComparables creates a new ArrayComparable
// from an existing array of Comparables
func NewArrayComparableFromComparables(comparables []Comparable) ArrayComparable {
	return ArrayComparable{comparables}
}

// Equal returns true if other is is equal to any of the comparables in the array
func (s ArrayComparable) Equal(other any) (bool, error) {
	for _, v := range s.theThings {
		if equal, err := v.Equal(other); err == nil && equal {
			return true, nil
		}
	}
	return false, nil
}

// Greater returns true if other is is greater than any of the things in the array
func (s ArrayComparable) Greater(other any) (bool, error) {
	for _, v := range s.theThings {
		if greater, err := v.Greater(other); err == nil && greater {
			return true, nil
		}
	}
	return false, nil
}
