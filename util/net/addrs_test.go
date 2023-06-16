package net

import (
	"net"
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