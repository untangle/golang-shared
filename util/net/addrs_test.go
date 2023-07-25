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

func v6NetIP(a byte, b byte, c byte, d byte) netip.Addr {
	addr := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, a, b, c, d}
	return netip.AddrFrom16(addr)
}

func TestIPSpecString(t *testing.T) {
	tests := []struct {
		name        string
		stringval   string
		shoulderr   bool
		shouldequal any
	}{
		{"ipv4 address", "132.123.123.1", false, net.IPv4(132, 123, 123, 1)},
		{"ipv4 net", "132.123.123.1/24", false,
			func() IPRange {
				_, net, _ := net.ParseCIDR("132.123.123.1/24")
				return NetToRange(net)
			}(),
		},
		{"ipv4 range", "132.123.123.1-132.123.123.3", false,
			IPRange{Start: v6NetIP(132, 123, 123, 1), End: v6NetIP(132, 123, 123, 3)},
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
				assert.True(t, tt.shoulderr)
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
			IPRange{v6NetIP(0, 0, 0, 0), v6NetIP(1, 0, 0, 0)},
			net.IPv4(0, 1, 0, 0),
			true,
		},
		{
			"basic out",
			IPRange{v6NetIP(0, 0, 0, 0), v6NetIP(1, 0, 0, 0)},
			net.IPv4(1, 1, 0, 0),
			false,
		},
		{
			"basic upper border",
			IPRange{v6NetIP(0, 0, 0, 0), v6NetIP(1, 0, 0, 0)},
			net.IPv4(1, 0, 0, 0),
			true,
		},
		{
			"basic lower border",
			IPRange{v6NetIP(0, 0, 0, 0), v6NetIP(1, 0, 0, 0)},
			net.IPv4(0, 0, 0, 0),
			true,
		},
		{
			"basic out, lower",
			IPRange{v6NetIP(1, 0, 0, 0), v6NetIP(2, 0, 0, 0)},
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
				Start: v6NetIP(192, 168, 25, 0),
				End:   v6NetIP(192, 168, 25, 255),
			},
			"192.168.25.0/24",
		},
		{
			"one high bit masked off",
			IPRange{
				Start: v6NetIP(192, 168, 25, 0),
				End:   v6NetIP(192, 168, 25, 127),
			},
			"192.168.25.0/25",
		},
		{
			"two high bits masked off",
			IPRange{
				Start: v6NetIP(192, 168, 25, 0),
				End:   v6NetIP(192, 168, 25, 63),
			},
			"192.168.25.0/26",
		},
		{
			"only low bit allowed",
			IPRange{
				Start: v6NetIP(192, 168, 25, 0),
				End:   v6NetIP(192, 168, 25, 1),
			},
			"192.168.25.0/31",
		},
		{
			"all bits masked",
			IPRange{
				Start: v6NetIP(192, 168, 25, 0),
				End:   v6NetIP(192, 168, 25, 0),
			},
			"192.168.25.0/32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IPSpecifierString(tt.network).Parse()
			switch netRange := result.(type) {
			case IPRange:
				// We are using True() instead of Equal() because the assert
				// library specifically tests for byteslices (which net.IP
				// is), and tests them differently.
				assert.True(t, tt.iprange.Start.Compare(netRange.Start) == 0,
					"net range starts: %s (expected) should equal %s", tt.iprange.Start, netRange.Start)
				assert.True(t, tt.iprange.End.Compare(netRange.End) == 0,
					"net range ends: %s (expected) should equal %s", tt.iprange.End, netRange.End)
			}
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
			fmt.Printf("Loaded %d lines from %s\n", len(lines), filename)
		} else {
			fmt.Printf("Error loading IPs: %v\n", err)
		}
	}
}

// variable used to iterate through the ipIndex
// across repeated calls to the Benchmark
var idx = 0

var netipArray = make([]netip.Addr, 0)

func BenchmarkIP4TestNetIP(b *testing.B) {
	b.StopTimer()

	if len(netipArray) == 0 {
		loadFile("testdata/ip4s.txt")
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					if ipx, err := netip.ParseAddr(line); err == nil {
						netipArray = append(netipArray, ipx)
					}
				}
			}
		}
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		ip := netipArray[idx]
		ipas4 := ip.As4()
		// then set the last octet to 0 and create a range between it and ip
		ipas4[3] = 0
		newip := netip.AddrFrom4(ipas4)
		ipPrefix := netip.PrefixFrom(newip, 24)

		assert.Truef(b, ipPrefix.Contains(ip), "Failed containment of %v\n", ip)

		idx = (idx + 1) % len(netipArray)
	}
}

var ipArray = make([]net.IP, 0)

func BenchmarkIP4Test(b *testing.B) {
	b.StopTimer()

	if len(ipArray) == 0 {
		loadFile("testdata/ip4s.txt")
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					ipx := net.ParseIP(line)
					ipArray = append(ipArray, ipx)
				}
			}
		}
	}
	mask := net.IPv4Mask(255, 255, 255, 0)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		newip := ipArray[idx]
		// then set the last octet to 0 and create a range between it and ip
		newip[3] = 0
		ipNet := net.IPNet{IP: newip, Mask: mask}

		assert.Truef(b, ipNet.Contains(newip), "Failed containment of %v\n", newip)
		idx = (idx + 1) % len(ipArray)
	}
}

func BenchmarkIP4Range(b *testing.B) {
	b.StopTimer()

	if len(ipArray) == 0 {
		loadFile("testdata/ip4s.txt")
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					ipx := net.ParseIP(line)
					ipArray = append(ipArray, ipx.To16())
				}
			}
		}
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		ip := ipArray[idx]
		// then set the last octet to 0 and create a range between it and ip
		ip[15] = 0
		start, _ := netip.AddrFromSlice(ip)
		ip[15] = 255
		end, _ := netip.AddrFromSlice(ip)
		ipRange := IPRange{Start: start, End: end}

		assert.Truef(b, ipRange.ContainsNetIP(end), "Failed containment of %v\n", end)

		idx = (idx + 1) % len(ipArray)
	}
}

func BenchmarkIP6TestNetIP(b *testing.B) {
	b.StopTimer()

	if len(netipArray) == 0 {
		loadFile("testdata/ip6s.txt")
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					if ipx, err := netip.ParseAddr(line); err == nil {
						netipArray = append(netipArray, ipx)
					}
				}
			}
		}
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		ip := netipArray[idx]
		ipas16 := ip.As16()
		// then set the last octet to 0 and create a range between it and ip
		ipas16[15] = 0
		newip := netip.AddrFrom16(ipas16)
		ipPrefix := netip.PrefixFrom(newip, 120)

		assert.Truef(b, ipPrefix.Contains(newip), "Failed containment of %v\n", newip)

		idx = (idx + 1) % len(netipArray)
	}
}

func BenchmarkIP6Test(b *testing.B) {
	b.StopTimer()
	if len(ipArray) == 0 {
		ipArray = make([]net.IP, 0)
		loadFile("testdata/ip6s.txt")
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					ipx := net.ParseIP(line)
					ipArray = append(ipArray, ipx)
				}
			}
		}
	}
	mask := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ip := ipArray[idx]
		newip := ip.To16()
		// then set the last octet to 0 and create a range between it and ip
		newip[15] = 0
		ipNet := net.IPNet{IP: newip, Mask: mask}

		assert.Truef(b, ipNet.Contains(ip), "Failed containment of %v\n", ip)
		idx = (idx + 1) % len(ipArray)
	}
}

func BenchmarkIP6Range(b *testing.B) {
	b.StopTimer()
	if len(ipArray) == 0 {
		loadFile("testdata/ip6s.txt")
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					ipx := net.ParseIP(line)
					ipArray = append(ipArray, ipx)
				}
			}
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ip := ipArray[idx]
		newip := ip.To16()
		// then set the last octet to 0 and create a range between it and ip
		newip[15] = 0
		start, _ := netip.AddrFromSlice(newip)
		newip[15] = 255
		end, _ := netip.AddrFromSlice(newip)
		ipRange := IPRange{Start: start, End: end}

		assert.Truef(b, ipRange.ContainsIP(ip), "Failed containment of %v\n", ip)
		idx = (idx + 1) % len(ipArray)
	}
}

func BenchmarkAll(b *testing.B) {
	// Load lines ahead of benchmark
	loadFile("testdata/ip4s.txt")
	idx = 0
	b.Run("IP4Test with net/netip", BenchmarkIP4TestNetIP)
	idx = 0
	b.Run("IP4Test with net(existing)", BenchmarkIP4Test)
	idx = 0
	b.Run("IP4Test using IPRange", BenchmarkIP4Range)

	// Reset lines for IPv6
	lines = make([]string, 0)
	loadFile("testdata/ip6s.txt")
	idx = 0
	b.Run("IP6Test with net/netip", BenchmarkIP6TestNetIP)
	idx = 0
	b.Run("IP6Test with net(existing)", BenchmarkIP6Test)
	idx = 0
	b.Run("IP6Test using IPRange", BenchmarkIP6Range)
}
