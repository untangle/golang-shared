package net

import (
	"bufio"
	"io/fs"
	"net/netip"
	"os"
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

func BenchmarkIPTest(b *testing.B) {
	ipArray := make([]netip.Addr, 0)
	if f, err := os.OpenFile("../../util/net/ips.txt", 0, fs.FileMode(os.O_RDONLY)); err == nil {
		defer f.Close()
		fileScanner := bufio.NewScanner(f)
		var lines []string
		for fileScanner.Scan() {
			lines = append(lines, fileScanner.Text())
		}
		for _, line := range lines {
			if len(line) > 0 {
				if line[0] != '#' {
					if ipx, err := netip.ParseAddr(line); err == nil {
						ipArray = append(ipArray, ipx)
					}
				} else {
					b.Logf("%s\n", line)
				}
			}
		}
	} else {
		b.Errorf("Error loading IPs: %v\n", err)
	}
	var zeroIP = netip.AddrFrom4([4]byte{0, 0, 0, 0})

	b.Logf("Loaded %d IPs\n", len(ipArray))

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, ip := range ipArray {
			ipRange := IPRange{Start: zeroIP, End: ip}
			assert.Truef(b, ipRange.Contains(ip), "Failed containment of %v\n", ip)
			assert.Truef(b, ipRange.Start.Compare(zeroIP) == 0, "Failed start check of %v\n", ip)
			assert.Truef(b, ipRange.End.Compare(ip) == 0, "Failed end check of %v\n", ip)
			assert.Truef(b, !ipRange.Contains(ip.Next()), "Failed containment of next: $v\n", ip.Next())

			ipRange = IPRange{Start: ip, End: ip}
			assert.Truef(b, ipRange.Contains(ip), "Failed containment of %v\n", ip)
			assert.Truef(b, ipRange.Start.Compare(ip) == 0, "Failed start check of %v\n", ip)
			assert.Truef(b, ipRange.End.Compare(ip) == 0, "Failed end check of %v\n", ip)
			assert.Truef(b, !ipRange.Contains(ip.Prev()), "Failed containment of previous: $v\n", ip.Prev())
			assert.Truef(b, !ipRange.Contains(ip.Next()), "Failed containment of next: $v\n", ip.Next())
		}
	}
}
