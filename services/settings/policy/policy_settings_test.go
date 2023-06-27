package policy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/util/net"
)

func TestGetAllPolicyConfigurationSettings(t *testing.T) {

	var result = map[string]interface{}{
		"enabled": true,
		"passList": []interface{}{
			map[string]interface{}{
				"description": "some test",
				"host":        "3.4.5.6/32",
			},
		},
		"redirect":    false,
		"sensitivity": float64(60),
	}

	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	policySettings, err := getAllPolicyConfigurationSettings(settingsFile)
	assert.Nil(t, err)
	assert.NotNil(t, policySettings)
	assert.Equal(t, 3, len(policySettings["threatprevention"]))
	assert.Equal(t, 1, len(policySettings["webfilter"]))
	assert.Equal(t, 1, len(policySettings["geoip"]))

	teachersUID := "60a9e031-4188-4d06-8083-108ebec63a9e"
	// Spot check a plugin setting.
	assert.EqualValues(t, result, policySettings["threatprevention"][teachersUID])
}

func TestGetPolicyPluginSettings(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	tpPolicies, _ := GetPolicyPluginSettings(settingsFile, "threatprevention")
	assert.Equal(t, 4, len(tpPolicies))
	webFilterPolicies, _ := GetPolicyPluginSettings(settingsFile, "webfilter")
	assert.Equal(t, 2, len(webFilterPolicies))
	geoIPPolicies, _ := GetPolicyPluginSettings(settingsFile, "geoip")
	assert.Equal(t, 2, len(geoIPPolicies))
}

func TestErrorGetPolicyPluginSettings(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	_, err := GetPolicyPluginSettings(settingsFile, "notapolicy")
	assert.NotNil(t, err)
}

func TestGroupUnmarshal(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings_group.json")
	policySettings := PolicySettings{}
	assert.Nil(t, settingsFile.UnmarshalSettingsAtPath(&policySettings, "policy_manager"))
	strlist, ok := policySettings.Groups[0].ItemsIPSpecList()
	assert.True(t, ok)
	assert.Equal(t, []net.IPSpecifierString{
		"1.2.3.4",
		"1.2.3.5/24",
		"1.2.3.4-1.2.3.20"}, strlist)
	endpointList, ok := policySettings.Groups[2].ItemsServiceEndpointList()
	assert.True(t, ok)
	assert.EqualValues(t, []ServiceEndpoint{
		{
			Protocol:    "TCP",
			IPSpecifier: "12.34.56.78",
			Port:        12345,
		},
		{
			Protocol:    "UDP",
			IPSpecifier: "12.34.56.0/24",
			Port:        12345,
		},
		{
			Protocol:    "UDP",
			IPSpecifier: "1.2.3.4-1.2.3.5",
			Port:        12345,
		},
	}, endpointList)
}

func TestGroupUnmarshalEdges(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectedErr bool
		expected    Group
	}{
		{name: "emptyjson", json: ``, expectedErr: true, expected: Group{}},
		{
			name: "Basic bad type test",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "badType",
                          "items:" []}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "okay ip list",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrList",
                          "items": ["132.123.123"]}`,
			expectedErr: false,
			expected: Group{
				Type:  "IPAddrList",
				Items: []net.IPSpecifierString{"132.123.123"},
				ID:    "702d4c99-9599-455f-8271-215e5680f038",
			}},
		{
			name: "okay geoip list",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "GeoIPLocation",
                          "items": ["AE", "AF"]}`,
			expectedErr: false,
			expected: Group{
				Type:  "GeoIPLocation",
				Items: []string{"AE", "AF"},
				ID:    "702d4c99-9599-455f-8271-215e5680f038",
			}},
		{
			name: "malformed JSON",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrList",
                          "items": [{]]}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "bad ip addrlist",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrList",
                          "items": [{}]}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "bad type",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrListBOGUS",
                          "items": []}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "bad items",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrList",
                          "items": false}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "bad items2",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrList",
                          "items": [{}, {}, {}]}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "bad service endpoint",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "ServiceEndpoint",
                          "items": ["googlywoogly"]}`,
			expectedErr: true,
			expected:    Group{}},
		{
			name: "emptylist",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "ServiceEndpoint",
                          "items": []}`,
			expectedErr: false,
			expected: Group{
				Type:  "ServiceEndpoint",
				Items: []ServiceEndpoint{},
				ID:    "702d4c99-9599-455f-8271-215e5680f038",
			},
		},
		{
			name: "bad sg endpoint list",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "ServiceEndpoint",
                          "items": [{"protocol": "UDP", "ipspoocifier": ""]}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "good sg endpoint list",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "ServiceEndpoint",
                          "items": [
                              {"protocol": "UDP", "ipspecifier": "123.123.123.123", "port": "2222"},
                              {"protocol": "TCP", "ipspecifier": "123.123.123.124", "port": "2223"}]}`,
			expectedErr: false,
			expected: Group{
				Type: ServiceEndpointType,
				ID:   "702d4c99-9599-455f-8271-215e5680f038",
				Items: []ServiceEndpoint{
					{Protocol: "UDP",
						IPSpecifier: "123.123.123.123",
						Port:        2222,
					},
					{Protocol: "TCP",
						IPSpecifier: "123.123.123.124",
						Port:        2223,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual Group
			if !tt.expectedErr {
				assert.Nil(t, json.Unmarshal([]byte(tt.json), &actual))
				assert.EqualValues(t, tt.expected, actual)
			} else {
				assert.NotNil(t, json.Unmarshal([]byte(tt.json), &actual))
			}
		})
	}
}
