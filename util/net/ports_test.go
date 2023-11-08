package net

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPortsSpecifierStringParse tests the Parse function of PortSpecifierString.
func TestPortsSpecifierStringParse(t *testing.T) {
	tests := []struct {
		name      string
		ps        PortSpecifierString
		shouldErr bool
		expected  any
	}{
		{
			"single port",
			PortSpecifierString("1234"),
			false,
			Port(1234),
		},
		{
			"port range",
			PortSpecifierString("1234-5678"),
			false,
			PortRange{Start: 1234, End: 5678},
		},
		{
			"invalid port range",
			PortSpecifierString("1234-567"),
			true,
			fmt.Errorf("invalid port range, start > end: 1234-567"),
		},
		{
			"Multiple dashes",
			PortSpecifierString("1234-567-890"),
			true,
			fmt.Errorf("invalid port specifier string range, contains too many -: 1234-567-890"),
		},
		{
			"Invalid start port",
			PortSpecifierString("1234a-5678"),
			true,
			fmt.Errorf("invalid port specifier string range, contains bad start port: 1234a"),
		},
		{
			"Invalid end port",
			PortSpecifierString("1234-5678a"),
			true,
			fmt.Errorf("invalid port specifier string range, contains bad end port: 5678a"),
		},
		{
			"Invalid port",
			PortSpecifierString("1234a"),
			true,
			fmt.Errorf("invalid port specifier: 1234a"),
		},
		{
			"Empty string",
			PortSpecifierString(""),
			true,
			fmt.Errorf("invalid port specifier: "),
		},
		{
			"Out of range port",
			PortSpecifierString("65536"),
			true,
			fmt.Errorf("invalid port specifier: 65536"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.ps.Parse()
			switch typed := val.(type) {
			case Port:
				assert.EqualValues(t, typed, tt.expected)
			case PortRange:
				assert.EqualValues(t, typed, tt.expected)
			case error:
				assert.True(t, tt.shouldErr)
			default:
				assert.FailNow(t, "invalid type: %T", typed)
			}
		})
	}
}

// TestPortRangeContainsPort tests the ContainsPort function of PortRange.
func TestPortRangeContainsPort(t *testing.T) {
	tests := []struct {
		name          string
		portRange     PortRange
		port          Port
		shouldContain bool
	}{
		{
			"basic in",
			PortRange{Start: 0, End: 100},
			50,
			true,
		},
		{
			"basic out",
			PortRange{Start: 0, End: 100},
			101,
			false,
		},
		{
			"basic upper limit",
			PortRange{Start: 0, End: 100},
			100,
			true,
		},
		{
			"basic lower limit",
			PortRange{Start: 0, End: 100},
			0,
			true,
		},
		{
			"basic out, lower",
			PortRange{Start: 100, End: 200},
			50,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldContain, tt.portRange.Contains(tt.port))
		})
	}
}
