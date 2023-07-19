package net

import (
	"bufio"
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
		{"ipv4 address", "132.123.123.1", false, net.IPv4(132, 123, 123, 1)},
		{"ipv4 net", "132.123.123.1/24", false,
			func() *net.IPNet {
				_, net, _ := net.ParseCIDR("132.123.123.1/24")
				return net
			}(),
		},
		{"ipv4 range", "132.123.123.1-132.123.123.3", false,
			IPRange{Start: net.IPv4(132, 123, 123, 1), End: net.IPv4(132, 123, 123, 3)},
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
			case net.IP, *net.IPNet, IPRange:
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
		ipAddr        net.IP
		shouldContain bool
	}{
		{
			"basic in",
			IPRange{net.IPv4(0, 0, 0, 0), net.IPv4(1, 0, 0, 0)},
			net.IPv4(0, 1, 0, 0),
			true,
		},
		{
			"basic out",
			IPRange{net.IPv4(0, 0, 0, 0), net.IPv4(1, 0, 0, 0)},
			net.IPv4(1, 1, 0, 0),
			false,
		},
		{
			"basic upper border",
			IPRange{net.IPv4(0, 0, 0, 0), net.IPv4(1, 0, 0, 0)},
			net.IPv4(1, 0, 0, 0),
			true,
		},
		{
			"basic lower border",
			IPRange{net.IPv4(0, 0, 0, 0), net.IPv4(1, 0, 0, 0)},
			net.IPv4(0, 0, 0, 0),
			true,
		},
		{
			"basic out, lower",
			IPRange{net.IPv4(1, 0, 0, 0), net.IPv4(2, 0, 0, 0)},
			net.IPv4(0, 0, 0, 1),
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
				Start: net.IPv4(192, 168, 25, 0),
				End:   net.IPv4(192, 168, 25, 255),
			},

			"192.168.25.0/24",
		},
		{
			"one high bit masked off",
			IPRange{
				Start: net.IPv4(192, 168, 25, 0),
				End:   net.IPv4(192, 168, 25, 127),
			},
			"192.168.25.0/25",
		},
		{
			"two high bits masked off",
			IPRange{
				Start: net.IPv4(192, 168, 25, 0),
				End:   net.IPv4(192, 168, 25, 63),
			},
			"192.168.25.0/26",
		},
		{
			"only low bit allowed",
			IPRange{
				Start: net.IPv4(192, 168, 25, 0),
				End:   net.IPv4(192, 168, 25, 1),
			},
			"192.168.25.0/31",
		},
		{
			"all bits masked",
			IPRange{
				Start: net.IPv4(192, 168, 25, 0),
				End:   net.IPv4(192, 168, 25, 0),
			},
			"192.168.25.0/32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, network, _ := net.ParseCIDR(tt.network)
			netRange := NetToRange(network)
			// We are using True() instead of Equal() because the assert
			// library specifically tests for byteslices (which net.IP
			// is), and tests them differently.
			assert.True(t, tt.iprange.Start.Equal(netRange.Start),
				"net range starts: %s (expected) should equal %s", tt.iprange.Start, netRange.Start)
			assert.True(t, tt.iprange.End.Equal(netRange.End),
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

func TestGetPrefixFromIPs(t *testing.T) {
	x := GetPrefixFromNetIPs(netip.AddrFrom4([4]byte{1, 1, 1, 1}), netip.AddrFrom4([4]byte{1, 1, 1, 253}))
	assert.Truef(t, x.Bits() == 120, "GetPrefixFromNetIPs should have returned 120")
}

func BenchmarkIP4TestNetIP(b *testing.B) {
	ipArray := make([]netip.Addr, 0)
	loadFile("testdata/ip4s.txt")
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
			ipas4[3] = 0
			newip := netip.AddrFrom4(ipas4)
			ipPrefix := netip.PrefixFrom(newip, 24)
			assert.Truef(b, ipPrefix.Contains(ip), "Failed containment of %v\n", ip)
			// then set the last octet to 0 and create a range between it and ip
			limit := byte(255)
			for octet := byte(0); octet < limit; octet++ {
				ipas4[3] = octet
				newip = netip.AddrFrom4(ipas4)
				assert.Truef(b, ipPrefix.Contains(newip), "Failed containment of %v\n", newip)
			}
		}
	}
}

func BenchmarkIP4Test(b *testing.B) {
	ipArray := make([]net.IP, 0)
	loadFile("testdata/ip4s.txt")
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
			// then set the last octet to 0 and create a range between it and ip
			newip := net.IP{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5],
				ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13],
				ip[14], 0}
			ipNet := net.IPNet{IP: newip, Mask: net.IPv4Mask(255, 255, 255, 0)}
			assert.Truef(b, ipNet.Contains(ip), "Failed containment of %v\n", ip)
			limit := byte(255)
			for octet := byte(0); octet < limit; octet++ {
				newip[15] = octet
				assert.Truef(b, ipNet.Contains(newip), "Failed containment of %v\n", newip)
			}
		}
	}
}

func BenchmarkIP6TestNetIP(b *testing.B) {
	ipArray := make([]netip.Addr, 0)
	loadFile("testdata/ip6s.txt")
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != '#' {
				if ipx, err := netip.ParseAddr(line); err == nil {
					ipArray = append(ipArray, ipx)
				}
			}
		}
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, ip := range ipArray {
			ipas16 := ip.As16()
			ipas16[15] = 0
			newip := netip.AddrFrom16(ipas16)
			ipPrefix := netip.PrefixFrom(newip, 120)
			assert.Truef(b, ipPrefix.Contains(ip), "Failed containment of %v\n", ip)
			// then set the last octet to 0 and create a range between it and ip
			limit := byte(255)
			for octet := byte(0); octet < limit; octet++ {
				ipas16[15] = octet
				newip = netip.AddrFrom16(ipas16)
				assert.Truef(b, ipPrefix.Contains(newip), "Failed containment of %v\n", newip)
			}
		}
	}
}

func BenchmarkIP6Test(b *testing.B) {
	ipArray := make([]net.IP, 0)
	loadFile("testdata/ip6s.txt")
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != '#' {
				ipx := net.ParseIP(line)
				ipArray = append(ipArray, ipx)
			}
		}
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, ip := range ipArray {
			// then set the last octet to 0 and create a range between it and ip
			newip := net.IP{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5],
				ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13],
				ip[14], 0}
			ipNet := net.IPNet{
				IP:   newip,
				Mask: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0},
			}
			assert.Truef(b, ipNet.Contains(ip), "Failed containment of %v\n", ip)
			limit := byte(255)
			for octet := byte(0); octet < limit; octet++ {
				newip[15] = octet
				assert.Truef(b, ipNet.Contains(newip), "Failed containment of %v\n", newip)
			}
		}
	}
}

func BenchmarkAll(b *testing.B) {
	// Load lines ahead of benchmark
	loadFile("testdata/ip4s.txt")
	b.Run("IP4Test with net/netip", BenchmarkIP4TestNetIP)
	b.Run("IP4Test with net(existing)", BenchmarkIP4Test)

	// Reset lines for IPv6
	lines = make([]string, 0)
	loadFile("testdata/ip6s.txt")
	b.Run("IP6Test with net/netip", BenchmarkIP6TestNetIP)
	b.Run("IP6Test with net(existing)", BenchmarkIP6Test)
}
