package net

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"net"
	"net/netip"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPSpecString(t *testing.T) {
	tests := []struct {
		name        string
		stringval   string
		shoudlerr   bool
		shouldequal any
	}{
		{"ipv4 address", "132.123.123.1", false, netip.AddrFrom4([4]byte{132, 123, 123, 1})},
		{"ipv4 net", "132.123.123.1/24", false,
			func() netip.Prefix {
				net, _ := netip.ParsePrefix("132.123.123.1/24")
				return net
			}(),
		},
		{"ipv4 range", "132.123.123.1-132.123.123.3", false,
			IPRange{Start: netip.AddrFrom4([4]byte{132, 123, 123, 1}),
				End: netip.AddrFrom4([4]byte{132, 123, 123, 3})},
		},
		{"ipv4 range, start less than end", "132.123.123.1-132.123.123.0", true, nil},
		{"ipv4 range, too many dashes", "132.123.123.1--132.123.123.20", true, nil},
		{"ipv4 CIDR net, badly formatted", "132.123.123.1//20", true, nil},
		{"bogus string", "booga", true, nil},
		{"empty string", "", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := IPSpecifierString(tt.stringval).Parse()
			switch typed := val.(type) {
			case netip.Addr, netip.Prefix, IPRange:
				assert.EqualValues(t, typed, tt.shouldequal)
			case error:
				assert.True(t, tt.shoudlerr)
			default:
				assert.FailNow(t, "invalid type: %T", typed)
			}
		})
	}
}

func TestIPRange(t *testing.T) {
	tests := []struct {
		name          string
		ipRange       IPRange
		ipAddr        netip.Addr
		shouldContain bool
	}{
		{
			"basic in",
			IPRange{netip.AddrFrom4([4]byte{0, 0, 0, 0}), netip.AddrFrom4([4]byte{1, 0, 0, 0})},
			netip.AddrFrom4([4]byte{0, 1, 0, 0}),
			true,
		},
		{
			"basic out",
			IPRange{netip.AddrFrom4([4]byte{0, 0, 0, 0}), netip.AddrFrom4([4]byte{1, 0, 0, 0})},
			netip.AddrFrom4([4]byte{1, 1, 0, 0}),
			false,
		},
		{
			"basic upper border",
			IPRange{netip.AddrFrom4([4]byte{0, 0, 0, 0}), netip.AddrFrom4([4]byte{1, 0, 0, 0})},
			netip.AddrFrom4([4]byte{1, 0, 0, 0}),
			true,
		},
		{
			"basic lower border",
			IPRange{netip.AddrFrom4([4]byte{0, 0, 0, 0}), netip.AddrFrom4([4]byte{1, 0, 0, 0})},
			netip.AddrFrom4([4]byte{0, 0, 0, 0}),
			true,
		},
		{
			"basic out, lower",
			IPRange{netip.AddrFrom4([4]byte{1, 0, 0, 0}), netip.AddrFrom4([4]byte{2, 0, 0, 0})},
			netip.AddrFrom4([4]byte{0, 0, 0, 1}),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldContain, tt.ipRange.Contains(tt.ipAddr))
		})
	}
}

func TestIPRangeFromCIDR(t *testing.T) {
	tests := []struct {
		name    string
		iprange IPRange
		network string
	}{
		{
			"simple, no bits",
			IPRange{
				Start: netip.AddrFrom4([4]byte{192, 168, 25, 0}),
				End:   netip.AddrFrom4([4]byte{192, 168, 25, 255}),
			},

			"192.168.25.0/24",
		},
		{
			"one high bit masked off",
			IPRange{
				Start: netip.AddrFrom4([4]byte{192, 168, 25, 0}),
				End:   netip.AddrFrom4([4]byte{192, 168, 25, 127}),
			},
			"192.168.25.0/25",
		},
		{
			"two high bits masked off",
			IPRange{
				Start: netip.AddrFrom4([4]byte{192, 168, 25, 0}),
				End:   netip.AddrFrom4([4]byte{192, 168, 25, 63}),
			},
			"192.168.25.0/26",
		},
		{
			"only low bit allowed",
			IPRange{
				Start: netip.AddrFrom4([4]byte{192, 168, 25, 0}),
				End:   netip.AddrFrom4([4]byte{192, 168, 25, 1}),
			},
			"192.168.25.0/31",
		},
		{
			"all bits masked",
			IPRange{
				Start: netip.AddrFrom4([4]byte{192, 168, 25, 0}),
				End:   netip.AddrFrom4([4]byte{192, 168, 25, 0}),
			},
			"192.168.25.0/32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, _ := netip.ParsePrefix(tt.network)
			netRange, _ := NetToRange(network)

			assert.True(t, tt.iprange.Start.Compare(netRange.Start) == 0,
				"net range starts: %s (expected) should equal %s", tt.iprange.Start, netRange.Start)
			assert.True(t, tt.iprange.End.Compare(netRange.End) == 0,
				"net range ends: %s (expected) should equal %s", tt.iprange.End, netRange.End)
		})
	}
}

var lines []string

func loadFile(filename string) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	if len(lines) == 0 {
		if f, err := os.OpenFile(filename, 0, fs.FileMode(os.O_RDONLY)); err == nil {
			defer f.Close()
			fileScanner := bufio.NewScanner(f)
			for fileScanner.Scan() {
				lines = append(lines, fileScanner.Text())
			}
		} else {
			fmt.Printf("Error loading IPs: %v\n", err)
		}
	}
}

func BenchmarkIPTest(b *testing.B) {
	ipArray := make([]netip.Addr, 0)
	loadFile("ips.txt")
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != '#' {
				if ipx, err := netip.ParseAddr(line); err == nil {
					ipArray = append(ipArray, ipx)
				}
			}
		}
	}
	for n := 0; n < b.N; n++ {
		for _, ip := range ipArray {
			ipas4 := ip.As4()
			limit := ipas4[3]
			if limit != 0 {
				ipas4[3] = 0
				newip := netip.AddrFrom4(ipas4)
				ipPrefix := netip.PrefixFrom(newip, 24)
				assert.Truef(b, ipPrefix.Contains(ip), "Failed containment of %v\n", ip)
				// then set the last octet to 0 and create a range between it and ip
				for octet := byte(0); octet < limit; octet++ {
					ipas4[3] = octet
					newip = netip.AddrFrom4(ipas4)
					assert.Truef(b, ipPrefix.Contains(newip), "Failed containment of %v\n", newip)
				}
			}
		}
	}
}

// These structures and functions are resurrected from the old addrs.go
// to support a comparative benchmark of netip versus net.ip.
// They are not referenced anywhere else
type IPRangeOld struct {
	Start net.IP
	End   net.IP
}

func (r IPRangeOld) Contains(ip net.IP) bool {
	return bytes.Compare(r.Start, ip) <= 0 &&
		bytes.Compare(r.End, ip) >= 0
}

func BenchmarkIPTestOld(b *testing.B) {
	ipArray := make([]net.IP, 0)
	loadFile("ips.txt")
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != '#' {
				ipx := net.ParseIP(line)
				ipArray = append(ipArray, ipx)
			}
		}
	}
	for n := 0; n < b.N; n++ {
		for _, ip := range ipArray {
			limit := ip[15]
			if limit != 0 {
				// then set the last octet to 0 and create a range between it and ip
				newip := net.IP{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5],
					ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13],
					ip[14], 0}
				ipNet := net.IPNet{IP: newip, Mask: net.IPv4Mask(255, 255, 255, 0)}
				assert.Truef(b, ipNet.Contains(ip), "Failed containment of %v\n", ip)
				intLimit := int(limit)
				for octet := 0; octet <= intLimit; octet++ {
					newip[15] = byte(octet)
					assert.Truef(b, ipNet.Contains(newip), "Failed containment of %v\n", newip)
				}
			}
		}
	}
}

func BenchmarkAll(b *testing.B) {
	loadFile("ips.txt")
	b.Run("IPTest", BenchmarkIPTest)
	b.Run("IPTestOld", BenchmarkIPTestOld)
}
