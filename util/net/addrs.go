package net

import (
	"fmt"
	"net"
)

const (
	// IPv4Str -- String denoting IPV4, used like an enum.
	IPv4Str = "IPV4"

	// IPv6Str -- string denoting IPv6, used like an enum.
	IPv6Str = "IPV6"

	// InvalidIPStr --  Error message for invalid IP.
	InvalidIPStr = "Invalid IP"
)

// CheckIpAddressType returns the type of the IP address.
func CheckIPAddressType(ip net.IP) (string, error) {
	if ip.To4() != nil {
		return IPv4Str, nil
	} else if ip.To16() != nil {
		return IPv6Str, nil
	} else {
		return "", fmt.Errorf("InvalidIPStr")
	}
}
