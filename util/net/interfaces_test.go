package interfaces

import (
	"fmt"
	"net"
	"testing"
	//logService "github.com/untangle/golang-shared/services/logger"
)

type testStruct struct {
	testInterface    Interface
	passNetworkHasIp string
	failNetworkHasIp string
}

var (
	testMap    map[string]testStruct
	ipv4Mask24 string
	ipv4Mask16 string
	ipv6Mask64 string
)

func setupSuite(t *testing.T) func(t *testing.T) {
	//var logger := logService.GetLoggerInstance()
	testMap = make(map[string]testStruct)
	ipv4Mask24 = "ipv4Mask24"
	testMap[ipv4Mask24] = testStruct{
		testInterface: Interface{
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  24,
			IsWAN:           false,
		},
		passNetworkHasIp: "192.168.0.10/24",
		failNetworkHasIp: "192.168.1.1/24",
	}
	ipv4Mask16 = "ipv4Mask16"
	testMap[ipv4Mask16] = testStruct{
		testInterface: Interface{
			Enabled:         true,
			V4StaticAddress: "192.168.0.1",
			V4StaticPrefix:  16,
			IsWAN:           false,
		},
		passNetworkHasIp: "192.168.56.10/16",
		failNetworkHasIp: "192.168.56.10/24",
	}
	ipv6Mask64 = "ipv6Mask64"
	testMap[ipv6Mask64] = testStruct{
		testInterface: Interface{
			Enabled:         true,
			V6StaticAddress: "2001:DB8:0000:0000:244:17FF:FEB6:D37D",
			v6StaticPrefix:  64,
			IsWAN:           false,
		},
		passNetworkHasIp: "2001:DB8:0000:0001:244:17FF:FEB6:D37D/64",
		failNetworkHasIp: "2001:DB8:0000:0000:244:17FF:FEB6:D37D/16",
	}

	return func(t *testing.T) {
		// shutdown
	}
}

func TestNetworkHasIp(t *testing.T) {
	tearDownSuite := setupSuite(t)
	defer tearDownSuite(t)

	runTest := func(cidrString string, intf Interface, expected bool, testName string) {
		ip, _, err := net.ParseCIDR(cidrString)
		if err != nil {
			actual := intf.NetworkHasIP(ip)
			if actual != expected {
				t.Errorf(fmt.Sprintf("Test %s: Expected '%t', got '%t'", testName, expected, actual))
			}
		}
	}

	testArr := []string{
		ipv4Mask24,
		ipv4Mask16,
		ipv6Mask64,
	}

	for _, testName := range testArr {
		test := testMap[testName]
		runTest(test.passNetworkHasIp, test.testInterface, true, testName)
		runTest(test.failNetworkHasIp, test.testInterface, false, testName)
	}
}

func TestGetNetwork(t *testing.T) {
	tearDownSuite := setupSuite(t)
	defer tearDownSuite(t)

	runTest := func(intf Interface, expectedStr string, testName string) {
		actual, err0 := intf.GetNetwork()
		if err0 != nil {
			_, expected, err1 := net.ParseCIDR(expectedStr)
			if err1 != nil {
				// if the networks contain the other's IP, then we got the correct network
				if !actual.Contains(expected.IP) || !expected.Contains(actual.IP) {
					t.Errorf(fmt.Sprintf("Test %s: IPs '%s' and '%s' are not on the same network", testName, actual.String(), expected.String()))
				}
			}

		}
	}

	testArr := []string{
		ipv4Mask24,
		ipv4Mask16,
		ipv6Mask64,
	}

	for _, testName := range testArr {
		test := testMap[testName]
		runTest(test.testInterface, test.passNetworkHasIp, testName)
	}
}
