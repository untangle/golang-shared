package interfaces

import (
	"fmt"
	"net"
	"testing"
)

type interfacesTestStruct struct {
	testInterface Interface
	passIp        string
	failIp        string
}

var (
	testMap    map[string]interfacesTestStruct
	ipv4Mask24 string
	ipv4Mask16 string
	ipv6Mask64 string
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
		passIp: "192.168.0.10/24",
		failIp: "192.168.1.1/24",
	}
	ipv4Mask16 = "ipv4Mask16"
	testMap[ipv4Mask16] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  16,
			IsWAN:           false,
		},
		passIp: "192.168.56.10/16",
		failIp: "192.1.56.10/24",
	}
	ipv6Mask64 = "ipv6Mask64"
	testMap[ipv6Mask64] = interfacesTestStruct{
		testInterface: Interface{
			Enabled:         true,
			V6StaticAddress: "2001:DB8:0000:0000:244:17FF:FEB6:D37D",
			V6StaticPrefix:  64,
			IsWAN:           false,
		},
		passIp: "2001:DB8::/64",
		failIp: "2002:DB8::/16",
	}

	return func(t *testing.T) {
		// shutdown
	}
}

func TestGetNetworks(t *testing.T) {
	tearDownSuite := interfacesTestSetup(t)
	defer tearDownSuite(t)

	runTest := func(testName string, shouldPass bool, ipStr string, intf Interface) {
		ip, _, err := net.ParseCIDR(ipStr)
		if err == nil {
			network := intf.GetNetworks()[0]
			if network.Contains(ip) != shouldPass {
				var printStr string
				if shouldPass {
					printStr = "Test %s: IP '%s' is not on the network '%s'"
				} else {
					printStr = "Test %s: IP '%s' is on the network '%s'"
				}
				t.Errorf(fmt.Sprintf(printStr, testName, ip.String(), network.String()))
			}
		} else {
			t.Errorf(fmt.Sprintf("Test %s: IP '%s' did not parse correctly", testName, ip.String()))
		}
	}

	testArr := []string{
		ipv4Mask24,
		ipv4Mask16,
		ipv6Mask64,
	}

	for _, testName := range testArr {
		test := testMap[testName]
		runTest(testName, true, test.passIp, test.testInterface)
		runTest(testName, false, test.failIp, test.testInterface)
	}
}
