package net

import (
	"errors"
	"fmt"
	"math/big"
	"net/netip"
	"strings"
)

// CheckIPAddressType returns the type of the IP address.
func CheckIPAddressTypeNetIP(ip netip.Addr) (string, error) {
	if ip.Is4() {
		return IPv4Str, nil
	} else if ip.Is6() {
		return IPv6Str, nil
	} else {
		return "", fmt.Errorf("InvalidIPStr")
	}
}

// IPRange is a range of IPs, from Start to End inclusive.
type IPRangeNetIP struct {
	Start  netip.Addr
	End    netip.Addr
	Prefix netip.Prefix
}

// Determine a netip.Prefix based on start and end
// This is only implemented for IPv4 for now
func GetPrefix(start netip.Addr, end netip.Addr) netip.Prefix {
	ipas16 := start.As16()
	bigTop := big.NetInt(0)
	bigTop.SetBytes(ipas16[0:8])
	bigStart := big.NewInt(0)
	bigStart.SetBytes(ipas16[8:16])
	ipas16 = end.As16()
	bigEnd := big.NewInt(0)
	bigEnd.SetBytes(ipas16[8:16])

	bits := 128
	bigMask := big.NewInt(0xFFFFFFFFFFFFFFFF)
	bigTwo := big.NewInt(2)
	for bits = 128; bits > 0; bits-- {
		if bigStart.And(bigStart, bigMask).Cmp(bigEnd.And(bigEnd, bigMask)) == 0 {
			break
		}
		bigMask.Mul(bigMask, bigTwo)
	}
	base16 := append(bigTop, bigStart.And(bigStart, bigMask).Bytes())
	return netip.PrefixFrom(netip.AddrFrom16(base16), bits)
}

func NewIPRangeNetIP(start netip.Addr, end, netip.Addr) {
	return IPRangeNetIP{
		Start: start,
		End: end,
		Prefix: GetPrefix(start, end),
	}
}

// Contains returns true if the ip is between the Start and End of r,
// inclusive.
func (r IPRangeNetIP) Contains(ip netip.Addr) bool {
	// Could potentially optimize this using the Prefix
	return r.Start.Compare(ip) <= 0 && r.End.Compare(ip) >= 0
}

// Parse returns the parsed specifier as one of:
// -- net.IP : regular IP.
// -- *net.IPNet: CIDR net, the specifier contained a slash (/).
// -- IPRange -- IPRange, if the specifier was a range <IP>-<IP>.
// -- error -- if the ip specifier was none of these we return an error object.
func (ss IPSpecifierString) ParseNetIP() any {
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

		return IPRangeNetIP{Start: start, End: end, Prefix: GetPrefix(start, end)}
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
func NetToRangeNetIP(prefix netip.Prefix) (IPRangeNetIP, error) {

	if !prefix.IsValid() {
		return IPRangeNetIP{}, errors.New("invalid prefix")
	}
	maskBits := prefix.Bits()
	if prefix.Addr().Is4In6() && maskBits < 96 {
		return IPRangeNetIP{}, errors.New("prefix with 4in6 address must have mask >= 96")
	}
	base := prefix.Masked().Addr()

	// the internal 128bit representation is private
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

	return IPRangeNetIP{base, last, prefix}, nil
}
