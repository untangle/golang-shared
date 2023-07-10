package net

import (
	"errors"
	"fmt"
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
func CheckIPAddressType(ip netip.Addr) (string, error) {
	if ip.Is4() {
		return IPv4Str, nil
	} else if ip.Is6() {
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
// inclusive.
func (r IPRange) Contains(ip netip.Addr) bool {
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
		start, err := netip.ParseAddr(parts[0])
		if err != nil {
			return fmt.Errorf("invalid ip specifier string range, contains bad IPs: %s",
				ss)
		}
		end, err := netip.ParseAddr(parts[1])
		if err != nil {
			return fmt.Errorf("invalid ip specifier string range, contains bad IPs: %s",
				ss)
		}
		if start.Compare(end) > 0 {
			return fmt.Errorf("invalid IP range, start > end: %s", ss)
		}

		return IPRange{Start: start, End: end}
	} else if strings.Contains(string(ss), "/") {
		if network, err := netip.ParsePrefix(string(ss)); err != nil {
			return err
		} else {
			return network
		}

	} else if ip, err := netip.ParseAddr(string(ss)); err == nil {
		return ip
	} else {
		return fmt.Errorf("invalid ip specifier: %s", ss)
	}
}

// NetToRange converts a *netip.Prefix to an IPRange.
func NetToRange(prefix netip.Prefix) (IPRange, error) {

	if !prefix.IsValid() {
		return IPRange{}, errors.New("invalid prefix")
	}
	maskBits := prefix.Bits()
	if prefix.Addr().Is4In6() && maskBits < 96 {
		return IPRange{}, errors.New("prefix with 4in6 address must have mask >= 96")
	}
	base := prefix.Masked().Addr()

	// the internal 128bit representation is privat
	// all calculations must be done in the bytes representation
	a16 := base.As16()

	if base.Is4() {
		maskBits += 96
	}

	// set host bits to 1
	for b := maskBits; b < 128; b++ {
		byteNum, bitInByte := b/8, 7-(b%8)
		a16[byteNum] |= 1 << uint(bitInByte)
	}

	// back to internal 128bit representation
	last := netip.AddrFrom16(a16)

	// unmap last to v4 if base is v4
	if base.Is4() {
		last = last.Unmap()
	}

	return IPRange{base, last}, nil
}
