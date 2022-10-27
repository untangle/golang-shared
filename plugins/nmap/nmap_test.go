package nmap

import (
	"fmt"
	"os"
	"testing"

	"github.com/untangle/discoverd/plugins/discovery"
)

// Test Setwork function with valid and invalid network IP
func TestSetNetwork(t *testing.T) {
	type args struct {
		network string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-valid-network",
			args: args{network: "172.16.0.2/24"},
		},
		{
			name: "test-invalid-network",
			args: args{network: "172.16.0.2/244"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetNetwork(tt.args.network)
		})
	}
}

// Test NmapcallBackHandler() with Valid and invalid IPV4, IPV6 host and network IPs. And test non nmap requests.
// It covers updateRetry() CheckIPAddressType() and isNewNmapReq() functions and NmapScan go routine.
func TestNmapcallBackHandler(t *testing.T) {
	type args struct {
		commands []discovery.Command
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-Namp-IPV4Net-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"192.168.0.1/24"}}}},
		},
		{
			name: "test-Namp-IPV4Host-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanHost, Arguments: []string{"192.168.0.1"}}}},
		},
		{
			name: "test-Namp-IPV6Host-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanHost, Arguments: []string{"fdf1:ab86:f5ab:1234::1"}}}},
		},
		{
			name: "test-Namp-IPV6Net-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"fdf1:ab86:f5ab:1234::1/64"}}}},
		},
		{
			name: "test-Namp-IPV6Net-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"fdf1:ab86:f5ab:1234::1/64"}}}},
		},
		{
			name: "test-Namp-IPV6Net-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"fdf1:ab86:f5ab:1234::1/64"}}}},
		},
		{
			name: "test-Namp-IPV6Net-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"fdf1:ab86:f5ab:1234::1/64"}}}},
		},
		{
			name: "test-Namp-IPV6Net-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"fdf1:ab86:f5ab:1234::1/64"}}}},
		},
		{
			name: "test-Non-Namp-Request",
			args: args{[]discovery.Command{{Command: 5, Arguments: []string{"fdf1:ab86:f5ab:1234::1/64"}}}},
		},
		{
			name: "test-Namp-Invalid_IPV6Net-Request",
			args: args{[]discovery.Command{{Command: discovery.CmdScanNet, Arguments: []string{"fdf1:ab86:f5ab:1234::1/644"}}}},
		},
		{
			name: "test-Namp-Invalid-IPV4Host-scan",
			args: args{[]discovery.Command{{Command: discovery.CmdScanHost, Arguments: []string{"192.168.0.256"}}}},
		},
		{
			name: "test-Namp-Nil",
			args: args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NmapcallBackHandler(tt.args.commands)
		})
	}
}

/*
// Test nmap plugin startup
func TestStart(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start()
		})
	}
}
*/
// Test removeProcess with existing and non existing entries
func Test_removeProcess(t *testing.T) {
	type args struct {
		args string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-existing-entry",
			args: args{"nmap -sT -O -F -6 -oX - fdf1:ab86:f5ab:1234::1/64"},
		},
		{
			name: "test-non-existing-entry",
			args: args{"nmap -sT -O -F -oX - fdf1:ab86:f5ab:1234::1/64"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeProcess(tt.args.args)
		})
	}
}

// Test processScan with test nmap scan output xml file
func Test_processScan(t *testing.T) {
	nampTestFile := "./testdata/nmap_output"
	nmapOutput, err := os.ReadFile(nampTestFile) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	type args struct {
		output []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-nmap-scan-output",
			args: args{output: nmapOutput},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processScan(tt.args.output)
		})
	}
}

func TestNmap_startNmap(t *testing.T) {
	type fields struct {
		autoNmapCollectionShutdown    chan bool
		autoNmapCollectionShutdownAck chan bool
		nmapSettings                  nmapPluginSettings
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test-nmap",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nmap := &Nmap{
				autoNmapCollectionShutdown:    tt.fields.autoNmapCollectionShutdown,
				autoNmapCollectionShutdownAck: tt.fields.autoNmapCollectionShutdownAck,
				nmapSettings:                  tt.fields.nmapSettings,
			}
			nmap.startNmap()
		})
	}
}
