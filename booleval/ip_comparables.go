package booleval

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/untangle/golang-shared/services/logger"
)

// IPComparable is a Comparable for net.IPs. It cannot be ordered.
type IPComparable struct {
	GreaterNotApplicable
	ipaddr netip.Addr
}

// NewIPComparable returns a new IPComparable, given the string. It
// calls net.ParseIP to get the IP.
func NewIPComparable(val string) IPComparable {
	if result, err := netip.ParseAddr(val); err != nil {
		logger.Warn("Error parsing IP: %s %v\n", val, err)
		return IPComparable{}
	} else {
		return IPComparable{ipaddr: result}
	}
}

// Equal -- returns true if other is an IP and is the same as i,
// or if other is an IPNet and contains i.
func (i IPComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case netip.Addr:
		return i.ipaddr.Compare(val) == 0, nil
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
			result, err := netip.ParseAddr(val)
			return i.ipaddr.Compare(result) == 0, err
		}
	default:
		return false, fmt.Errorf(
			"booleval IPComparable.Equal: cannot coerce  %v(type %T) to net.IPaddr",
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
		return IPNetComparable{ipnet: ipnet.Masked()}, nil
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
	case netip.Prefix:
		return i.ipnet.Contains(val.Addr()), nil
	case string:
		result, err := netip.ParseAddr(val)
		return i.ipnet.Contains(result), err
	default:
		return false, fmt.Errorf(
			"booleval IPNetComparable.Equal: cannot coerce %v(type %T) to net.IP", val, val)
	}
}

var _ Comparable = IPNetComparable{}
