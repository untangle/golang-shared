package interfaces

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type interfacesTestStruct struct {
	testInterface Interface
	passIps       []string
	failIps       []string
}

var (
	testMap     map[string]interfacesTestStruct
	ipv4Mask24  string
	ipv4Mask16  string
	ipv6Mask64  string
	ipv4Aliases string
	ipv6Aliases string
	ipBridged   string
)

func interfacesTestSetup(t *testing.T) func(t *testing.T) {
	testMap = make(map[string]interfacesTestStruct)
	ipv4Mask24 = "ipv4Mask24"
	testMap[ipv4Mask24] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  24,
			IsWAN:           false,
		},
		passIps: []string{"192.168.0.10/24"},
		failIps: []string{"192.168.1.1/24"},
	}
	ipv4Mask16 = "ipv4Mask16"
	testMap[ipv4Mask16] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  16,
			IsWAN:           false,
		},
		passIps: []string{"192.168.56.10/16"},
		failIps: []string{"192.1.56.10/24"},
	}
	ipv6Mask64 = "ipv6Mask64"
	testMap[ipv6Mask64] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V6StaticAddress: "2001:DB8:0000:0000:244:17FF:FEB6:D37D",
			V6StaticPrefix:  64,
			IsWAN:           false,
		},
		passIps: []string{"2001:DB8::/64"},
		failIps: []string{"2002:DB8::/16"},
	}

	ipv4Aliases = "ipv4Aliases"
	testMap[ipv4Aliases] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  16,
			V4Aliases: []V4IpAliases{
				{
					V4Address: "172.16.0.2",
					V4Prefix:  24,
				},
				{
					V4Address: "172.18.0.2",
					V4Prefix:  24,
				},
			},

			IsWAN: false,
		},
		passIps: []string{"192.168.0.1/16", "172.16.0.2/24", "172.18.0.2/24"},
		failIps: []string{"172.19.0.2/24", "172.192.0.2/24"},
	}

	ipBridged = "ipBridged"
	testMap[ipBridged] = interfacesTestStruct{
		testInterface: Interface{
			ConfigType:      ConfigTypeBridged,
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  16,
			V4Aliases: []V4IpAliases{
				{
					V4Address: "172.16.0.2",
					V4Prefix:  24,
				},
				{
					V4Address: "172.18.0.2",
					V4Prefix:  24,
				},
			},

			IsWAN: false,
		},
		passIps: []string{},
		failIps: []string{},
	}

	ipv6Aliases = "ipv6Aliases"
	testMap[ipv6Aliases] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V6StaticAddress: "2001:DB8:0000:0000:244:17FF:FEB6:D37D",
			V6StaticPrefix:  64,
			V6Aliases: []V6IpAliases{
				{
					V6Address: "FD2F:C9C2:3EE9:35E5::1",
					V6Prefix:  "64",
				},
				{
					V6Address: "FD2B:055E:1ED4:8DC3::1",
					V6Prefix:  "64",
				},
			},

			IsWAN: false,
		},
		passIps: []string{"2001:DB8:0000:0000:244:17FF:FEB6:D37D/64", "FD2F:C9C2:3EE9:35E5::1/64", "FD2B:055E:1ED4:8DC3::1/64"},
		failIps: []string{"FD2F:C9C2:3EE9:3566::1/64", "FD2B:055E:1ED4:8888::1/64"},
	}

	return func(t *testing.T) {
		// shutdown
	}
}

// Test GetNetworks for static IPV4 and IPV6 addresses and test IPV4 and IPV6 aliases addresses
func TestGetNetworks(t *testing.T) {
	tearDownSuite := interfacesTestSetup(t)
	defer tearDownSuite(t)

	runTest := func(testName string, shouldPass bool, ipStrs []string, intf Interface) {
		var found bool
		for _, ipStr := range ipStrs {
			ip, _, err := net.ParseCIDR(ipStr)
			if err == nil {
				for _, network := range intf.GetNetworks() {
					found = false
					if network.Contains(ip) {
						found = true
						break
					} else {
						continue
					}
				}
				if found != shouldPass {
					var printStr string
					if shouldPass {
						printStr = "Test %s: IP '%s' is not on the network '%s'"
					} else {
						printStr = "Test %s: IP '%s' is on the network '%s'"
					}
					t.Errorf(fmt.Sprintf(printStr, testName, ip.String(), ipStr))
				}

			} else {
				t.Errorf(fmt.Sprintf("Test %s: IP '%s' did not parse correctly from %s", testName, ip.String(), ipStr))
			}
		}
	}

	testArr := []string{
		ipv4Mask24,
		ipv4Mask16,
		ipv6Mask64,
		ipv4Aliases,
		ipv6Aliases,
		ipBridged,
	}

	for _, testName := range testArr {
		test := testMap[testName]
		runTest(testName, true, test.passIps, test.testInterface)
		runTest(testName, false, test.failIps, test.testInterface)
	}
}

func cidr2Net(cidr []string) []*net.IPNet {
	nets := []*net.IPNet{}
	for _, str := range cidr {
		_, net, _ := net.ParseCIDR(str)
		nets = append(nets, net)
	}
	return nets
}
func TestGetMostSpecificPrefix(t *testing.T) {

	type testParams struct {
		nets        []string
		ip          string
		expectedNet string
	}
	tests := []testParams{
		{
			nets:        []string{"192.168.128.0/17", "192.168.0.0/16"},
			ip:          "192.168.128.1",
			expectedNet: "192.168.128.0/17",
		},
		{
			nets:        []string{"10.10.192.0/18", "10.10.192.0/24"},
			ip:          "10.10.192.1",
			expectedNet: "10.10.192.0/24",
		},
	}

	for _, test := range tests {
		outputNet := MostSpecificPrefixMatch(cidr2Net(test.nets), net.ParseIP(test.ip))
		assert.NotNil(t, outputNet)
		_, expectedNet, _ := net.ParseCIDR(test.expectedNet)
		assert.EqualValues(t, outputNet, expectedNet)
	}
}
