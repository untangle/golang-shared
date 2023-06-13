package net

import (
	"bytes"
	"fmt"
	"net"
	"strings"
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

// IPSpecifierString is a string in the form of an IP range, CIDR address, or regular IP.
type IPSpecifierString string

// IPRange is a range of IPs, from Start to End inclusive.
type IPRange struct {
	Start net.IP
	End   net.IP
}

// Contains returns true if the ip is between the Start and End of r,
// inclusive.
func (r IPRange) Contains(ip net.IP) bool {
	return bytes.Compare(r.Start, ip) <= 0 &&
		bytes.Compare(r.End, ip) >= 0

}

// Parse returns the parsed specifier as one of:
// -- net.IP : regular IP.
// -- *net.IPNet: CIDR net, the specifier contained a slash (/).
// -- IPRange -- IPRange, if the specifier was a range <IP>-<IP>.
// -- error -- if the ip specifier was none of these we return an error object.
func (ss IPSpecifierString) Parse() any {
	if strings.Contains(string(ss), "-") {
		parts := strings.Split(string(ss), "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid ip specifier string range, contains too many -: %s",
				ss)
		}
		start := net.ParseIP(parts[0])
		end := net.ParseIP(parts[1])

		if start == nil || end == nil {
			return fmt.Errorf("invalid ip specifier string range, contains bad IPs: %s",
				ss)
		}

		if bytes.Compare(start, end) > 0 {
			return fmt.Errorf("invalid IP range, start > end: %s", ss)
		}

		return IPRange{Start: start, End: end}
	} else if strings.Contains(string(ss), "/") {
		if _, network, err := net.ParseCIDR(string(ss)); err != nil {
			return err
		} else {
			return network
		}

	} else if ip := net.ParseIP(string(ss)); ip != nil {
		return ip
	} else {
		return fmt.Errorf("invalid ip specifier: %s", ss)
	}
}

// NetToRange converts a *net.IPNet to an IPRange.
func NetToRange(network *net.IPNet) IPRange {
	masked := network.IP.Mask(network.Mask)
	lower := make(net.IP, len(network.IP), len(network.IP))
	upper := make(net.IP, len(network.IP), len(network.IP))
	copy(lower, masked)
	copy(upper, masked)
	ones, bits := network.Mask.Size()
	maskedBytes := (bits - ones) / 8
	remainderBits := (bits - ones) % 8

	for i := 1; i <= maskedBytes; i++ {
		upper[len(masked)-i] = 0xff
	}

	if remainderBits != 0 {
		remainderMask := (1 << (remainderBits)) - 1
		upper[len(masked)-(maskedBytes+1)] |= byte(remainderMask)
	}
	return IPRange{Start: lower, End: upper}
}
