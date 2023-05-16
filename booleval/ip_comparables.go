package booleval

import (
	"fmt"
	"net"
	"strings"
)

type IPComparable struct {
	GreaterNotApplicable
	ipaddr net.IP
}

func NewIPComparable(val string) IPComparable {
	return IPComparable{ipaddr: net.ParseIP(val)}
}

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

type IPNetComparable struct {
	GreaterNotApplicable
	ipnet net.IPNet
}

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
