package arp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/settings"
	disc_pb "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

func TestArpDevices(t *testing.T) {

	testParams := []struct {
		arpfile      string
		settingsfile string
		entries      []*disc.DeviceEntry
		error        string
	}{
		{
			arpfile:      "./testdata/mock_arp_proc",
			settingsfile: "./testdata/settings.json",
			entries: []*disc.DeviceEntry{
				{DiscoveryEntry: disc_pb.DiscoveryEntry{
					MacAddress:  "00:11:22:33:44:55",
					IPv4Address: "192.168.66.1",
					Arp: &disc_pb.ARP{
						Ip:  "192.168.66.1",
						Mac: "00:11:22:33:44:55",
					}}},
			},
		},
		{
			arpfile:      "./testdata/mock_arp_proc_3wans",
			settingsfile: "./testdata/settings-3wans.json",
			entries: []*disc.DeviceEntry{
				{DiscoveryEntry: disc_pb.DiscoveryEntry{
					MacAddress:  "00:11:22:33:44:55",
					IPv4Address: "192.168.66.3",
					Arp: &disc_pb.ARP{
						Ip:  "192.168.66.3",
						Mac: "00:11:22:33:44:55",
					}}},
				{DiscoveryEntry: disc_pb.DiscoveryEntry{
					MacAddress:  "00:11:22:33:44:88",
					IPv4Address: "192.168.66.4",
					Arp: &disc_pb.ARP{
						Ip:  "192.168.66.4",
						Mac: "00:11:22:33:44:88",
					}}},
			},
		},
		{
			arpfile:      "./testdata/bogusarp",
			settingsfile: "./testdata/settings-3wans.json",
			entries:      []*disc.DeviceEntry{},
			error:        "couldn't open arp file",
		},
		{
			arpfile:      "./testdata/mock_arp_proc",
			settingsfile: "./testdata/settings-bogus.json",
			entries:      []*disc.DeviceEntry{},
			error:        "couldn't unmarshall network settings",
		},
		{
			arpfile:      "./testdata/emptyarp",
			settingsfile: "./testdata/settings.json",
			entries:      []*disc.DeviceEntry{},
		},
	}

	for _, params := range testParams {
		settings := settings.NewSettingsFile(params.settingsfile)
		scanner := newArpScanner(
			settings,
			params.arpfile)
		entries, err := scanner.getArpEntriesFromFile()
		if params.error != "" {
			assert.True(t, strings.Contains(err.Error(), params.error),
				fmt.Sprintf("The returned error: %s should contain the given substring: %s",
					err.Error(),
					params.error))

		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, len(entries), len(params.entries))
		for i := range entries {
			assert.True(t, proto.Equal(entries[i], params.entries[i]))
		}
	}
}
