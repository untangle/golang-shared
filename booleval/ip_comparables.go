package booleval

import (
	"fmt"
	"net/netip"
	"strings"
)

// IPComparable is a Comparable for net.IPs. It cannot be ordered.
type IPComparable struct {
	GreaterNotApplicable
	ipaddr netip.Addr
}

// NewIPComparable returns a new IPComparable, given the string. It
// calls net.ParseIP to get the IP.
func NewIPComparable(val string) IPComparable {
	testAddr, err := netip.ParseAddr(val)
	if err == nil {
		return IPComparable{ipaddr: testAddr}
	}
	// Maybe need to be returning an err here
	return IPComparable{}
}

// Equal -- returns true if other is an IP and is the same as i,
// or if other is an IPNet and contains i.
func (i IPComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case netip.Addr:
		return 0 == i.ipaddr.Compare(val), nil
	case string:
		if strings.Contains(val, "/") {
			ipprefix, err := netip.ParsePrefix(val)
			if err != nil {
				return false, fmt.Errorf(
					"value %s contained a /, but is not a CIDR address",
					val)
			}
			return ipprefix.Contains(i.ipaddr), nil
		} else {
			ip, err := netip.ParseAddr(val)
			return 0 == i.ipaddr.Compare(ip), err
		}
	default:
		return false, fmt.Errorf(
			"booleval IPComparable.Equal: cannot coerce  %v(type %T) to netip.Addr",
			val,
			val)
	}
}

var _ Comparable = IPComparable{}

// IPNetComparable is a comparable representing an IP network (subnet).
type IPNetComparable struct {
	GreaterNotApplicable
	ipnet netip.Prefix
}

// NewIPOrIPNetComparable is a generic constructor for either an
// IPComparable or an IPNetComparable, given the format of the
// string. It may return an error if the string is malformed.
func NewIPOrIPNetComparable(val string) (Comparable, error) {
	if strings.Contains(val, "/") {
		ipnet, err := netip.ParsePrefix(val)
		if err != nil {
			return nil, fmt.Errorf(
				"value %s contained a /, but is not a CIDR address",
				val)
		}
		return IPNetComparable{ipnet: ipnet}, nil
	} else {
		result, err := netip.ParseAddr(val)
		return IPComparable{ipaddr: result}, err
	}
}

// Equal returns true if other is an IP address that is contained
// within the subnet of i, or a string representing an IP contained in
// i.
func (i IPNetComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case netip.Addr:
		return i.ipnet.Contains(val), nil
	case string:false