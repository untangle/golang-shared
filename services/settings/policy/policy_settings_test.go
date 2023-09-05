package policy

import (
	"encoding/json"
	"testing"

	"github.com/google/gopacket/layers"
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
			Protocol: uint(layers.IPProtocolTCP),
			Port:     12345,
		},
		{
			Protocol: uint(layers.IPProtocolUDP),
			Port:     12345,
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
				Name:  "someBogus",
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
				Name:  "someBogus",
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
				Name:  "ServiceEndpointTest",
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
                          "items": [{"protocol": 17]}`,
			expectedErr: true,
			expected:    Group{},
		},
		{
			name: "good sg endpoint list",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
                         "type": "ServiceEndpoint",
                          "items": [    
                              {"protocol": 17, "port": 2222},
                              {"protocol": 6, "port": 2223}]}`,
			expectedErr: false,
			expected: Group{
				Name:        "ServiceEndpointTest",
				Description: "Description",
				Type:        ServiceEndpointType,
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items: []ServiceEndpoint{
					{
						Protocol: uint(layers.IPProtocolUDP),
						Port:     2222,
					},
					{
						Protocol: uint(layers.IPProtocolTCP),
						Port:     2223,
					},
				},
			},
		},
		{
			name: "interface list",
			json: `{"name": "InterfaceListTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "description": "description",
                         "type": "Interface",
                          "items": [1, 2, 3]}`,
			expectedErr: false,
			expected: Group{
				Name:        "InterfaceListTest",
				Description: "description",
				Type:        InterfaceType,
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items:       []uint{1, 2, 3},
			},
		},
		{
			name: "bad iface list",
			json: `{"name": "InterfaceListTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "description": "description",
                         "type": "Interface",
                          "items": [1, "boog", 3]}`,
			expectedErr: true,
			expected:    Group{},
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

func TestGroupMarshal(t *testing.T) {
	tests := []struct {
		name         string
		group        Group
		expectedJSON string
	}{
		{
			name: "okay ip list",
			group: Group{
				Name:        "someBogus",
				Description: "Description",
				Type:        "IPAddrList",
				Items:       []net.IPSpecifierString{"132.123.123"},
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
			},
			expectedJSON: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
                         "type": "IPAddrList",
                          "items": ["132.123.123"]}`,
		},
		{
			name: "okay geoip list",
			group: Group{
				Name:        "someBogus",
				Description: "Description",
				Type:        "GeoIPLocation",
				Items:       []string{"AE", "AF"},
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
			},
			expectedJSON: `{"name": "someBogus",
			"id": "702d4c99-9599-455f-8271-215e5680f038",
			"description": "Description",
			"type": "GeoIPLocation",
			"items": ["AE", "AF"]}`,
		},
		{
			name: "good sg endpoint list",
			group: Group{
				Name:        "ServiceEndpointTest",
				Description: "Description",
				Type:        ServiceEndpointType,
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items: []ServiceEndpoint{
					{
						Protocol: uint(layers.IPProtocolUDP),
						Port:     2222,
					},
					{
						Protocol: uint(layers.IPProtocolTCP),
						Port:     2223,
					},
				},
			},
			expectedJSON: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
                         "type": "ServiceEndpoint",
                          "items": [
                              {"protocol": 17, "port": 2222},
                              {"protocol": 6, "port": 2223}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.group)
			assert.Nil(t, err)
			assert.JSONEq(t, tt.expectedJSON, string(data))
		})
	}
}

// Tests unmarshalling the PolicyCondition type with various combos of valid/invalid CIDR addresses and ports
func TestUnmarshalPolicyCondition(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		shouldErr bool
		expected  PolicyCondition
	}{{
		name: "bad type",
		json: `{
				"op": "==",
				"type": "I am not a type",
				"value": ["192.168.5.6/32"]
			}`,
		shouldErr: true,
		expected:  PolicyCondition{},
	},
		{
			name: "ipv4 w mask",
			json: `{
				"op": "==",
				"type": "CLIENT_ADDRESS",
				"value": ["192.168.5.6/32"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    "==",
				CType: "CLIENT_ADDRESS",
				Value: []string{"192.168.5.6/32"},
			},
		},
		{
			name: "ipv4 no mask",
			json: `{
				"op": "==",
				"type": "CLIENT_ADDRESS",
				"value": ["192.168.5.6"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    "==",
				CType: "CLIENT_ADDRESS",
				Value: []string{"192.168.5.6/32"},
			},
		},
		{
			name: "invalid ipv4 w mask",
			json: `{
				"op": "==",
				"type": "SERVER_ADDRESS",
				"value": ["192.168.5.256/32"]
			}`,
			shouldErr: true,
			expected:  PolicyCondition{},
		},
		{
			name: "ipv6 w mask",
			json: `{
				"op": "==",
				"type": "SERVER_ADDRESS",
				"value": ["fd00::1/8"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    "==",
				CType: "SERVER_ADDRESS",
				Value: []string{"fd00::1/8"},
			},
		},
		{
			name: "ipv6 no mask",
			json: `{
				"op": "==",
				"type": "SERVER_ADDRESS",
				"value": ["fd00::1"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    "==",
				CType: "SERVER_ADDRESS",
				Value: []string{"fd00::1/64"},
			},
		},
		{
			name: "valid port",
			json: `{
				"op": "==",
				"type": "CLIENT_PORT",
				"value": ["22"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    "==",
				CType: "CLIENT_PORT",
				Value: []string{"22"},
			},
		},
		{
			name: "invalid port",
			json: `{
				"op": "==",
				"type": "SERVER_PORT",
				"value": ["-1"]
			}`,
			shouldErr: true,
			expected:  PolicyCondition{},
		},
		{
			name: "test time",
			json: `{
				"op": ">=",
				"type": "TIME_OF_DAY",
				"value": ["9:00am"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    ">=",
				CType: "TIME_OF_DAY",
				Value: []string{"9:00am"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualCondition PolicyCondition

			if tt.shouldErr {
				assert.NotNil(t, json.Unmarshal([]byte(tt.json), &actualCondition))
			} else {
				assert.Nil(t, json.Unmarshal([]byte(tt.json), &actualCondition))
				assert.NotNil(t, actualCondition)
				assert.EqualValues(t, tt.expected, actualCondition)
			}
		})
	}
}

func TestPolicyConfigurationJSON(t *testing.T) {
	tests := []struct {
		name                   string
		inputData              string
		wantErr                bool
		wantUnmarshalledConfig PolicyConfiguration
		wantMarshalledJson     string
	}{
		{
			name: "validJsonWithAppSettings",
			inputData: `{
					"id": "A1",
					"name": "B2",
					"description": "C3",
					"key": "value",
					"setting_field": "D4"
				}`,
			wantUnmarshalledConfig: PolicyConfiguration{
				ID:          "A1",
				Name:        "B2",
				Description: "C3",
				AppSettings: map[string]any{
					"setting_field": "D4",
					"key":           "value",
				},
			},
			wantMarshalledJson: `{
				"id": "A1",
				"name": "B2",
				"description": "C3",
				"key": "value",
				"setting_field": "D4"
			}`,
		},
		{
			name: "validJsonWithoutAppSettings",
			inputData: `{
					"id": "A1",
					"name": "B2",
					"description": "C3"
				}`,
			wantUnmarshalledConfig: PolicyConfiguration{
				ID:          "A1",
				Name:        "B2",
				Description: "C3",
				AppSettings: map[string]any{},
			},
			wantMarshalledJson: `{
				"id": "A1",
				"name": "B2",
				"description": "C3"
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pConfig PolicyConfiguration
			err := json.Unmarshal(([]byte)(tt.inputData), &pConfig)
			assert.NoError(t, err)

			assert.EqualValues(t, &tt.wantUnmarshalledConfig, &pConfig)

			jsonValue, err := json.Marshal(tt.wantUnmarshalledConfig)
			assert.NoError(t, err)

			assert.JSONEq(t, tt.wantMarshalledJson, string(jsonValue))
		})
	}
}
