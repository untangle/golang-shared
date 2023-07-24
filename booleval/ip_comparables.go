package booleval

import (
	"fmt"
	"net"
	"strings"

	utilnet "github.com/untangle/golang-shared/util/net"
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

// Return the next IP based on a given IPComparable or error if the IP is invalid
func (ip IPComparable) Next() (net.IP, error) {
	addr := ip.ipaddr
	len := len(addr)
	result := make([]byte, len)
	copy(result, addr)
	for i := len - 1; i >= 0; i-- {
		if result[i]++; result[i] != 0 {
			break
		}
	}
	if len == net.IPv4len {
		return net.IPv4(result[0], result[1], result[2], result[3]), nil
	} else if len == net.IPv6len {
		return result, nil
	}
	return nil, fmt.Errorf("could not handle IPv6 %v", addr)
}

// Return the next IP based on the  end of a given subnet
func (ipnet IPNetComparable) Next() net.IP {
	next := utilnet.NetToRange(&ipnet.ipnet).End.Next()

	return next.AsSlice()
}
