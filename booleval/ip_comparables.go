package booleval

import (
	"fmt"
	"net"
	"strings"
)

// IPComparable is a Comparable for net.IPs. It cannot be ordered.
type IPComparable struct {
	GreaterNotApplicable
	ipaddr net.IP
}

// NewIPComparable returns a new IPComparable, given the string. It
// calls net.ParseIP to get the IP.
func NewIPComparable(val string) IPComparable {
	return IPComparable{ipaddr: net.ParseIP(val)}
}

// Equal -- returns true if other is an IP and is the same as i,
// or if other is an IPNet and contains i.
func (i IPComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case net.IP:
		return i.ipaddr.Equal(val), nil
	case string:
		if strings.Contains(val, "/") {
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return false, fmt.Errorf(
					"value %s contained a /, but is not a CIDR address",
					val)
			}
			return ipnet.Contains(i.ipaddr), nil
		} else {
			return i.ipaddr.Equal(net.ParseIP(val)), nil
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
	ipnet net.IPNet
}

// NewIPOrIPNetComparable is a generic constructor for either an
// IPComparable or an IPNetComparable, given the format of the
// string. It may return an error if the string is malformed.
func NewIPOrIPNetComparable(val string) (Comparable, error) {
	if strings.Contains(val, "/") {
		_, ipnet, err := net.ParseCIDR(val)
		if err != nil {
			return nil, fmt.Errorf(
				"value %s contained a /, but is not a CIDR address",
				val)
		}
		return IPNetComparable{ipnet: *ipnet}, nil
	} else {
		return IPComparable{ipaddr: net.ParseIP(val)}, nil
	}
}

// Equal returns true if other is an IP address that is contained
// within the subnet of i, or a string representing an IP contained in
// i.
func (i IPNetComparable) Equal(other any) (bool, error) {
	switch val := other.(type) {
	case net.IP:
		return i.ipnet.Contains(val), nil
	case string:
		return i.ipnet.Contains(net.ParseIP(val)), nil
	default:
		return false, fmt.Errorf(
			"booleval IPNetComparable.Equal: cannot coerce %v(type %T) to net.IP", val, val)
	}
}

var _ Comparable = IPNetComparable{}
