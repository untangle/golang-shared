package net

import (
	"fmt"
	"net"
	"net/netip"
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

// CheckIPAddressType returns the type of the IP address.
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
	Start netip.Addr
	End   netip.Addr
}

// Contains returns true if the ip is between the Start and End of r,
// inclusive. This is retained for backwards compatibility.
func (r IPRange) Contains(ip net.IP) bool {
	ipNetIP, _ := netip.AddrFromSlice(ip.To16())
	return r.ContainsNetIP(ipNetIP)
}

// ContainsNetIP returns true if the given netip.ADdr is between the Start and End of r,
// inclusive.
func (r IPRange) ContainsNetIP(ip netip.Addr) bool {
	return r.Start.Compare(ip) <= 0 && r.End.Compare(ip) >= 0
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
		if start, err := netip.ParseAddr(parts[0]); err != nil {
			return fmt.Errorf("invalid ip specifier string range, contains bad IPs: %s",
				ss)
		} else if end, err := netip.ParseAddr(parts[1]); err != nil {
			return fmt.Errorf("invalid ip specifier string range, contains bad IPs: %s",
				ss)
		} else if start.Compare(end) > 0 {
			return fmt.Errorf("invalid IP range, start > end: %s", ss)
		} else {
			// This is required to coerce the addresses into full IPV6 format
			// If this isn't done then comparisons fail because of bit count difference.
			return IPRange{Start: netip.AddrFrom16(start.As16()), End: netip.AddrFrom16(end.As16())}
		}
	} else if strings.Contains(string(ss), "/") {
		if _, network, err := net.ParseCIDR(string(ss)); err != nil {
			return err
		} else {
			return NetToRange(network)
		}
	} else if ip := net.ParseIP(string(ss)); ip != nil {
		return ip
	} else {
		return fmt.Errorf("invalid ip specifier: %s", ss)
	}
}

// NetToRange converts a *net.IPNet to an IPRange.
func NetToRange(network *net.IPNet) IPRange {
	// This is required to coerce the addresses into full IPV6 format
	// If this isn't done then comparisons fail because of bit count difference.
	masked := network.IP.Mask(network.Mask).To16()
	len := len(masked)
	lower := make(net.IP, len)
	upper := make(net.IP, len)
	copy(lower, masked)
	copy(upper, masked)
	ones, bits := network.Mask.Size()
	maskedBytes := (bits - ones) / 8
	remainderBits := (bits - ones) % 8

	// Maybe this can be optimized
	for i := 1; i <= maskedBytes; i++ {
		upper[len-i] = 0xff
	}
	if remainderBits != 0 {
		remainderMask := (1 << (remainderBits)) - 1
		upper[len-(maskedBytes+1)] |= byte(remainderMask)
	}
	lowerNetIP, _ := netip.AddrFromSlice(lower)
	upperNetIP, _ := netip.AddrFromSlice(upper)
	return IPRange{Start: lowerNetIP, End: upperNetIP}
}
