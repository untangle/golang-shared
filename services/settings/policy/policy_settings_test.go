package policy

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/util/net"
)

func TestGetAllPolicyConfigs(t *testing.T) {

	var result = PolicyConfiguration{
		Description: "TP students",
		Type:        "mfw-config-threatprevention",
		Name:        "TP for students",
		ID:          "d9b27e4a-2b8b-4500-a64a-51e7ee5777d5",
		Enabled:     false,
		Settings: map[string]interface{}{
			"enabled": true,
			"passList": []interface{}{
				map[string]interface{}{
					"description": "some test",
					"host":        "3.4.5.6/32",
				},
			},
			"redirect":    false,
			"sensitivity": float64(60),
		},
	}

	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	policySettings, err := GetAllPolicyConfigs(settingsFile)
	assert.NoError(t, err)
	assert.NotNil(t, policySettings)
	assert.Len(t, policySettings["mfw-config-threatprevention"], 2)
	assert.Len(t, policySettings["mfw-config-webfilter"], 1)
	assert.Len(t, policySettings["mfw-config-geoipfilter"], 3)

	teachersUID := "d9b27e4a-2b8b-4500-a64a-51e7ee5777d5"
	// Spot check a plugin setting.
	assert.EqualValues(t, &result, policySettings["mfw-config-threatprevention"][teachersUID])
}

func TestGetPolicyPluginSettings(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	tpPolicies, _ := GetPolicyPluginSettings(settingsFile, "threatprevention")
	assert.Len(t, tpPolicies, 3)
	webFilterPolicies, _ := GetPolicyPluginSettings(settingsFile, "webfilter")
	assert.Len(t, webFilterPolicies, 2)
	geoIPPolicies, _ := GetPolicyPluginSettings(settingsFile, "geoip")
	assert.Len(t, geoIPPolicies, 4)

	// Get the default and make sure it matches the expected object
	var defaultObj = PolicyConfiguration{
		Name:        "",
		ID:          "00000000-0000-0000-0000-000000000000",
		Description: "",
		Type:        "mfw-config-threatprevention",
		Settings: map[string]interface{}{
			"enabled":     false,
			"passList":    []interface{}{},
			"redirect":    false,
			"sensitivity": (float64)(20),
		}}
	assert.Equal(t, &defaultObj, tpPolicies["00000000-0000-0000-0000-000000000000"])
}

func TestErrorGetPolicyPluginSettings(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	_, err := GetPolicyPluginSettings(settingsFile, "notapolicy")
	assert.NotNil(t, err)
}

type unmarshalTest struct {
	name        string
	json        string
	expectedErr bool
	expected    Object
}

// runUnmarshalTest runs the unmarshal test.
func runUnmarshalTest(t *testing.T, tests []unmarshalTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual Object
			if !tt.expectedErr {
				assert.NoError(t, json.Unmarshal([]byte(tt.json), &actual))
				assert.EqualValues(t, actual, tt.expected)
			} else {
				assert.Error(t, actual.UnmarshalJSON([]byte(tt.json)))
			}
		})
	}
}

func TestRulesUnmarshal(t *testing.T) {
	tests := []unmarshalTest{
		{
			name: "Geo Rule Tester",
			json: `{"name": "GeoipRuleObject Name",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "mfw-rule-geoip",
                         "description": "GeoipRuleObject Description",
                         "conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
                         "action": {
                            "type": "SET_CONFIGURATION",
                            "configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
                            "key": "mfw-rule-geoip"
                            }
                          }`,
			expectedErr: false,
			expected: Object{
				Name:        "GeoipRuleObject Name",
				Type:        GeoipRuleObject,
				Description: "GeoipRuleObject Description",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "SET_CONFIGURATION",
					UUID: "1202b42e-2f21-49e9-b42c-5614e04d0031",
					Key:  "mfw-rule-geoip",
				},
				ID: "c2428365-65be-4901-bfc0-bde2b310fedf",
			},
		},

		{
			name: "bad rule object type",
			json: `{"name": "Geo Rule Tester",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "asdfasdf",
                         "description": "Whatever",
                         "conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
                         "action": {
                            "type": "SET_CONFIGURATION",
                            "configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
                            "key": "GeoipRuleObject"
                            }
                          }`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "rule object without action",
			json: `{"name": "Geo Rule Tester",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "mfw-rule-geoip",
                         "description": "Whatever",
                         "conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
                          }`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "ApplicationControlRuleObject test",
			json: `{"name": "ApplicationControlRuleObject Tester",
										"id": "c2428365-65be-4902-bfc0-bde2b310fedf",
										"type": "mfw-rule-applicationcontrol",
										"description": "ApplicationControlRuleObject",
										"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
										"action": {
										"type": "SET_CONFIGURATION",
										"configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
										"key": "mfw-rule-applicationcontrol"
										}
										}`,
			expectedErr: false,
			expected: Object{
				Name:        "ApplicationControlRuleObject Tester",
				Type:        ApplicationControlRuleObject,
				Description: "ApplicationControlRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "SET_CONFIGURATION",
					UUID: "1202b42e-2f21-49e9-b42c-5614e04d0031",
					Key:  "mfw-rule-applicationcontrol",
				},
				ID: "c2428365-65be-4902-bfc0-bde2b310fedf",
			},
		},
		{
			name: "CaptivePortalRuleObject test",
			json: `{"name": "CaptivePortalRuleObject Tester",
									"id": "c2428365-65be-4903-bfc0-bde2b310fedf",
									"type": "mfw-rule-captiveportal",
									"description": "CaptivePortalRuleObject",
									"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
									"action": {
									"type": "SET_CONFIGURATION",
									"configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
									"key": "mfw-rule-captiveportal"
									}
									}`,
			expectedErr: false,
			expected: Object{
				Name:        "CaptivePortalRuleObject Tester",
				Type:        CaptivePortalRuleObject,
				Description: "CaptivePortalRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "SET_CONFIGURATION",
					UUID: "1202b42e-2f21-49e9-b42c-5614e04d0031",
					Key:  "mfw-rule-captiveportal",
				},
				ID: "c2428365-65be-4903-bfc0-bde2b310fedf",
			},
		},
		{
			name: "NATRuleObject test",
			json: `{"name": "NATRuleObject Tester",
							"id": "c2428365-65be-4904-bfc0-bde2b310fedf",
							"type": "mfw-rule-nat",
							"description": "NATRuleObject",
							"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
							"action": {
							"type": "SNAT",
                                                        "snat_address": "192.168.56.2"
							}
							}`,
			expectedErr: false,
			expected: Object{
				Name:        "NATRuleObject Tester",
				Type:        NATRuleObject,
				Description: "NATRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type:        "SNAT",
					SNATAddress: "192.168.56.2",
				},
				ID: "c2428365-65be-4904-bfc0-bde2b310fedf",
			},
		},
		{
			name: "PortForwardRuleObject test",
			json: `{"name": "PortForwardRuleObject Tester",
							"id": "c2428365-65be-4905-bfc0-bde2b310fedf",
							"type": "mfw-rule-portforward",
							"description": "PortForwardRuleObject",
							"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
							"action": {
							    "type": "DNAT",
                                                            "dnat_address": "192.168.100.3",
                                                            "dnat_port": "81"
                                                         }

							}`,
			expectedErr: false,
			expected: Object{
				Name:        "PortForwardRuleObject Tester",
				Type:        PortForwardRuleObject,
				Description: "PortForwardRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type:        "DNAT",
					DNATAddress: "192.168.100.3",
					DNATPort:    "81",
				},
				ID: "c2428365-65be-4905-bfc0-bde2b310fedf",
			},
		},
		{
			name: "SecurityRuleObject Accept test",
			json: `{"name": "SecurityRuleObject Accept Tester",
			                         "id": "c2428365-65be-4906-bfc0-bde2b310fedf",
			                         "type": "mfw-rule-security",
			                         "description": "SecurityRuleObject",
			                         "conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
			                         "action": {
			                            "type": "ACCEPT",
			                            "configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
			                            "key": "mfw-rule-security"
			                            }
			                          }`,
			expectedErr: false,
			expected: Object{
				Name:        "SecurityRuleObject Accept Tester",
				Type:        SecurityRuleObject,
				Description: "SecurityRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "ACCEPT",
					UUID: "1202b42e-2f21-49e9-b42c-5614e04d0031",
					Key:  "mfw-rule-security",
				},
				ID: "c2428365-65be-4906-bfc0-bde2b310fedf",
			},
		},
		{
			name: "SecurityRuleObject Reject test",
			json: `{"name": "SecurityRuleObject Reject Tester",
			                         "id": "c2428365-65be-4916-bfc0-bde2b310fedf",
			                         "type": "mfw-rule-security",
			                         "description": "SecurityRuleObject",
			                         "conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
			                         "action": {
			                            "type": "REJECT",
			                            "configuration_id": "1202b42e-2f21-49ea-b42c-5614e04d0031",
			                            "key": "mfw-rule-security"
			                            }
			                          }`,
			expectedErr: false,
			expected: Object{
				Name:        "SecurityRuleObject Reject Tester",
				Type:        SecurityRuleObject,
				Description: "SecurityRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "REJECT",
					UUID: "1202b42e-2f21-49ea-b42c-5614e04d0031",
					Key:  "mfw-rule-security",
				},
				ID: "c2428365-65be-4916-bfc0-bde2b310fedf",
			},
		},
		{
			name: "ShapingRuleObject test",
			json: `{"name": "ShapingRuleObject Tester",
							"id": "c2428365-65be-4906-bfc0-bde2b310fedf",
							"type": "mfw-rule-shaping",
							"description": "ShapingRuleObject",
							"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
							"action": {
							"type": "SET_CONFIGURATION",
							"configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
							"key": "mfw-rule-shaping"
							}
							}`,
			expectedErr: false,
			expected: Object{
				Name:        "ShapingRuleObject Tester",
				Type:        ShapingRuleObject,
				Description: "ShapingRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "SET_CONFIGURATION",
					UUID: "1202b42e-2f21-49e9-b42c-5614e04d0031",
					Key:  "mfw-rule-shaping",
				},
				ID: "c2428365-65be-4906-bfc0-bde2b310fedf",
			},
		},
		{
			name: "WANPolicyRuleObject test",
			json: `{"name": "WANPolicyRuleObject Tester",
							"id": "c2428365-65be-4907-bfc0-bde2b310fedf",
							"type": "mfw-rule-wanpolicy",
							"description": "WANPolicyRuleObject",
							"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
							"action": {
							"type": "SET_CONFIGURATION",
							"configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
							"key": "mfw-rule-wanpolicy"
							}
							}`,
			expectedErr: false,
			expected: Object{
				Name:        "WANPolicyRuleObject Tester",
				Type:        WANPolicyRuleObject,
				Description: "WANPolicyRuleObject",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "SET_CONFIGURATION",
					UUID: "1202b42e-2f21-49e9-b42c-5614e04d0031",
					Key:  "mfw-rule-wanpolicy",
				},
				ID: "c2428365-65be-4907-bfc0-bde2b310fedf",
			},
		},
		{
			name: "quota rule test",
			json: `{"name": "quota rule test",
         			"id": "c2428365-65be-4907-bfc0-bde2b310fedf",
                   		"type": "mfw-rule-quota",
                   		"description": "QUOTAMAN",
                   		"conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
                   		"action": {
                            		"type": "APPLY_QUOTA",
                                        "configuration_id": "1458dc12-a9c2-4d0c-8203-1340c61c2c3e"
                         	 }
			}`,
			expectedErr: false,
			expected: Object{
				Name:        "quota rule test",
				Type:        QuotaRuleObject,
				Description: "QUOTAMAN",
				Conditions:  []string{"1458dc12-a9c2-4d0c-8203-1340c61c2c3b"},
				Action: &Action{
					Type: "APPLY_QUOTA",
					UUID: "1458dc12-a9c2-4d0c-8203-1340c61c2c3e",
				},
				ID: "c2428365-65be-4907-bfc0-bde2b310fedf",
			},
		},
	}
	runUnmarshalTest(t, tests)

}

func TestUnmarshalQuotas(t *testing.T) {
	tests := []unmarshalTest{
		{
			name: "Quota test",
			json: `{"name": "Quota",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "mfw-quota",
                         "description": "My quota description",
                         "settings": {
                               "amount_bytes": 100000,
                               "refresh": "1h"
                          },
                          "action": {
                            "type": "REJECT"
                            }
                          }`,
			expectedErr: false,
			expected: Object{
				Name:        "Quota",
				Type:        QuotaType,
				Description: "My quota description",
				Action: &Action{
					Type: "REJECT",
				},
				ID: "c2428365-65be-4901-bfc0-bde2b310fedf",
				Settings: &QuotaSettings{
					AmountBytes:     100000,
					RefreshInterval: QuotaRefreshTime(time.Hour),
				},
			},
		},
		{
			name: "Quota test",
			json: `{"name": "Quota",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "mfw-quota",
                         "description": "My quota description",
                         "settings": {
                               "amount_bytes": "10g000",
                               "refresh": "1h"
                          },
                          "action": {
                            "type": "REJECT"
                            }
                          }`,
			expectedErr: true,
		},
		{
			name: "Quota test",
			json: `{"name": "Quota",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "mfw-quota",
                         "description": "My quota description",
                         "settings": {
                               "amount_bytes": 100000,
                               "refresh": "1h1m2s"
                          },
                          "action": {
                            "type": "REJECT"
                            }
                          }`,
			expectedErr: false,
			expected: Object{
				Name:        "Quota",
				Type:        QuotaType,
				Description: "My quota description",
				Action: &Action{
					Type: "REJECT",
				},
				ID: "c2428365-65be-4901-bfc0-bde2b310fedf",
				Settings: &QuotaSettings{
					AmountBytes: 100000,
					RefreshInterval: QuotaRefreshTime(time.Hour +
						time.Minute +
						2*time.Second),
				},
			},
		},
		{
			name: "Quota test",
			json: `{"name": "Quota",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "mfw-quota",
                         "description": "My quota description",
                         "settings": {
                               "amount_bytes": 100000,
                               "refresh": "1googly"
                          },
                          "action": {
                            "type": "REJECT"
                            }
                          }`,
			expectedErr: true,
		},
	}

	runUnmarshalTest(t, tests)
}

// TestUnmarshalQuotaSettingsJSON test unmarshalling the settings.json
// with quotas. Since we test quotas more thoroughly in
// TestUnmarshalQuotas, this is just to make sure the thing works
// together from files.
func TestUnmarshalQuotaSettingsJSON(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	var quotas []Object
	err := settingsFile.UnmarshalSettingsAtPath(&quotas, "policy_manager", "quotas")
	assert.NoError(t, err)
	assert.Greater(t, len(quotas), 0, "There should be at least one quota in the settings.json")
}

func TestObjectUnmarshal(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings_group.json")
	var objects []Object
	err := settingsFile.UnmarshalSettingsAtPath(&objects, "policy_manager", "objects")
	assert.NoError(t, err)
	strlist, ok := objects[0].ItemsIPSpecList()
	assert.True(t, ok)

	assert.Equal(t, []net.IPSpecifierString{
		"1.2.3.4",
		"1.2.3.5/24",
		"1.2.3.4-1.2.3.20"}, strlist)
	endpointList, ok := objects[2].ItemsServiceEndpointList()
	assert.True(t, ok)
	assert.EqualValues(t, []ServiceEndpoint{
		{
			Protocol: []uint{uint(layers.IPProtocolTCP), uint(layers.IPProtocolUDP)},
			Port:     []net.PortSpecifierString{"12345", "80", "53"},
		},
		{
			Protocol: []uint{uint(layers.IPProtocolUDP)},
			Port:     []net.PortSpecifierString{"12345", "11", "22", "67", "66"},
		},
	}, endpointList)
}

// Test Unmarshalling an ApplicationObject from test_settings.json
func TestApplicationObjectUnmarshal(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	var objects []Object
	assert.Nil(t, settingsFile.UnmarshalSettingsAtPath(&objects, "policy_manager", "objects"))
	for i := range objects {
		if objects[i].Type == ApplicationType {
			if applicationObject, ok := objects[i].ItemsApplicationObject(); ok {
				if len(applicationObject.Port) > 0 && len(applicationObject.IPAddrList) > 0 {
					assert.EqualValues(t, ApplicationObject{
						Port:       []net.PortSpecifierString{"80", "8088", "443"},
						IPAddrList: []net.IPSpecifierString{"1.2.3.4", "2.3.4.5-3.4.5.6", "4.5.6.7/32"},
					}, applicationObject)
				} else if len(applicationObject.Port) > 0 {
					assert.EqualValues(t, ApplicationObject{
						Port:       []net.PortSpecifierString{"80", "8088", "443"},
						IPAddrList: nil,
					}, applicationObject)

				} else if len(applicationObject.IPAddrList) > 0 {
					assert.EqualValues(t, ApplicationObject{
						Port:       nil,
						IPAddrList: []net.IPSpecifierString{"1.2.3.4", "2.3.4.5-3.4.5.6", "4.5.6.7/32"},
					}, applicationObject)
				}
			} else {
				// Empty ApplicationObject is returned if anything goes wrong
				// Returning an empty object rather than nil prevents the objects loading from failing
				assert.Zero(t, len(applicationObject.Port)+len(applicationObject.IPAddrList))
			}
		}
	}
}

func TestGroupUnmarshalEdges(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectedErr bool
		expected    Object
	}{
		{name: "emptyjson", json: ``, expectedErr: true, expected: Object{}},
		{
			name: "Basic bad type test",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "badType",
                          "items:" []}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "okay ip list",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-ipaddress",
                          "items": ["132.123.123"]}`,
			expectedErr: false,
			expected: Object{
				Name:  "someBogus",
				Type:  "mfw-object-ipaddress",
				Items: []net.IPSpecifierString{"132.123.123"},
				ID:    "702d4c99-9599-455f-8271-215e5680f038",
			}},
		{
			name: "okay geoip list",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-geoip",
                          "items": ["AE", "AF"]}`,
			expectedErr: false,
			expected: Object{
				Name:  "someBogus",
				Type:  "mfw-object-geoip",
				Items: []string{"AE", "AF"},
				ID:    "702d4c99-9599-455f-8271-215e5680f038",
			}},
		{
			name: "malformed JSON",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-ipaddress",
                          "items": [{]]}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "bad ip addrlist",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-ipaddress",
                          "items": [{}]}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "bad type",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "IPAddrListBOGUS",
                          "items": []}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "bad items",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-ipaddress",
                          "items": false}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "bad items2",
			json: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-ipaddress",
                          "items": [{}, {}, {}]}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "bad service endpoint",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-service",
                          "items": ["googlywoogly"]}`,
			expectedErr: true,
			expected:    Object{}},
		{
			name: "emptylist",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "mfw-object-service",
                          "items": []}`,
			expectedErr: false,
			expected: Object{
				Name:  "ServiceEndpointTest",
				Type:  "mfw-object-service",
				Items: []ServiceEndpoint{},
				ID:    "702d4c99-9599-455f-8271-215e5680f038",
			},
		},
		{
			name: "bad sg endpoint list",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "type": "ServiceEndpoint",
                          "items": [{"protocol": [17]]}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "good sg endpoint list",
			json: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
                         "type": "mfw-object-service",
                          "items": [
                              {"protocol": [17,6,1], "port": ["2222", "80", "88"]},
                              {"protocol": [6], "port": ["2223", "11", "53"]}
							  ]}`,
			expectedErr: false,
			expected: Object{
				Name:        "ServiceEndpointTest",
				Description: "Description",
				Type:        "mfw-object-service",
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items: []ServiceEndpoint{
					{
						Protocol: []uint{uint(layers.IPProtocolUDP), uint(layers.IPProtocolTCP), uint(layers.IPProtocolICMPv4)},
						Port:     []net.PortSpecifierString{"2222", "80", "88"},
					},
					{
						Protocol: []uint{uint(layers.IPProtocolTCP)},
						Port:     []net.PortSpecifierString{"2223", "11", "53"},
					},
				},
			},
		},
		{
			name: "good ApplicationObject",
			json: `{"name": "ApplicationObject Test 1",
					"id": "702d4c99-9599-455f-dead-215e5680f038",
					"description": "Description",
					"type": "mfw-object-application",
					"items": [
						{
							"port": ["2222", "80", "88"],
							"ips": ["1.2.3.4", "2.3.4.5-3.4.5.6", "7.8.9.0/24"]
						}
					]}`,
			expectedErr: false,
			expected: Object{
				Name:        "ApplicationObject Test 1",
				Description: "Description",
				Type:        "mfw-object-application",
				ID:          "702d4c99-9599-455f-dead-215e5680f038",
				Items: []ApplicationObject{
					{
						Port:       []net.PortSpecifierString{"2222", "80", "88"},
						IPAddrList: []net.IPSpecifierString{"1.2.3.4", "2.3.4.5-3.4.5.6", "7.8.9.0/24"},
					},
				},
			},
		},
		{
			name: "bad ApplicationObject",
			json: `{"name": "Bad ApplicationObject Test 1",
					"id": "702d4c99-9599-455f-deac-215e5680f038",
					"description": "Description",
					"type": "mfw-object-application",
					"items": [
						{ 
							"port": "gobus",
							"ips": ["1.2.3.4", "2.3.4.5-3.4.5.6", "7.8.9.0/24"]
						}
					]}`,
			expectedErr: true,
		},
		{
			name: "interface list",
			json: `{"name": "InterfaceListTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "description": "description",
                         "type": "Interface",
                          "items": [1, 2, 3]}`,
			expectedErr: false,
			expected: Object{
				Name:        "InterfaceListTest",
				Description: "description",
				Type:        InterfaceType,
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items:       []uint{1, 2, 3},
			},
		},
		{
			name: "good ApplicationObjectGroup",
			json: `{"name": "ApplicationObjectGroup Test 1",
					"id": "702d4c99-959a-455f-dead-215e5680f038",
					"description": "Description",
					"type": "mfw-object-application-group",
					"items": [
						"8105f355-cb98-43eb-deaf-74542a524abb",
						"8105f355-cb98-43eb-dead-74542a524abb"
					]}`,
			expectedErr: false,
			expected: Object{
				Name:        "ApplicationObjectGroup Test 1",
				Description: "Description",
				Type:        "mfw-object-application-group",
				ID:          "702d4c99-959a-455f-dead-215e5680f038",
				Items: []string{
					"8105f355-cb98-43eb-deaf-74542a524abb",
					"8105f355-cb98-43eb-dead-74542a524abb",
				},
			},
		},
		{
			name: "bad ApplicationObjectGroup",
			json: `{"name": "Bad ApplicationObjectGroup Test 1",
					"id": "702d4c99-959a-455f-dead-215e5680f038",
					"description": "Description",
					"type": "mfw-object-application-group",
					"items": [
						12345
					]}`,
			expectedErr: true,
		},
		{
			name: "bad iface list",
			json: `{"name": "InterfaceListTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
                         "description": "description",
                         "type": "Interface",
                          "items": [1, "boog", 3]}`,
			expectedErr: true,
			expected:    Object{},
		},
		{
			name: "condition object",
			json: `{
		                    "name": "blooblah",
		                    "id": "702d4c99-9599-455f-8271-215e5680f039",
		                    "type": "mfw-object-condition",
		                    "items": [
		                        {
		                         "op": "==",
		                         "type": "SERVER_ADDRESS",
		                         "value": []
		                        }
		                    ]
		                }`,
			expectedErr: false,
			expected: Object{
				Name: "blooblah",
				ID:   "702d4c99-9599-455f-8271-215e5680f039",
				Type: ConditionType,
				Items: []*PolicyCondition{
					{
						Op:    "==",
						CType: "SERVER_ADDRESS",
						Value: []string{},
					},
				},
			},
		},
		{
			name: "condition group object",
			json: `{
                            "name": "blooblah",
                            "id": "702d4c99-9599-455f-8271-215e5680f039",
                            "type": "mfw-object-condition-group",
                            "items": [
                                 "702d4c99-9599-455f-8271-215e5680f038"
                            ]
                        }`,
			expectedErr: false,
			expected: Object{
				Name: "blooblah",
				ID:   "702d4c99-9599-455f-8271-215e5680f039",
				Type: ConditionGroupType,
				Items: []string{
					"702d4c99-9599-455f-8271-215e5680f038",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual Object
			if !tt.expectedErr {
				assert.NoError(t, json.Unmarshal([]byte(tt.json), &actual))
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
		object       Object
		expectedJSON string
	}{
		{
			name: "okay ip list",
			object: Object{
				Name:        "someBogus",
				Description: "Description",
				Type:        IPObjectType,
				Items:       []net.IPSpecifierString{"132.123.123"},
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
			},
			expectedJSON: `{"name": "someBogus",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
                         "type": "mfw-object-ipaddress",
                          "items": ["132.123.123"]}`,
		},
		{
			name: "okay geoip list",
			object: Object{
				Name:        "someBogus",
				Description: "Description",
				Type:        GeoIPObjectType,
				Items:       []string{"AE", "AF"},
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
			},
			expectedJSON: `{"name": "someBogus",
			"id": "702d4c99-9599-455f-8271-215e5680f038",
			"description": "Description",
			"type": "mfw-object-geoip",
			"items": ["AE", "AF"]}`,
		},
		{
			name: "good sg endpoint list",
			object: Object{
				Name:        "ServiceEndpointTest",
				Description: "Description",
				Type:        ServiceEndpointObjectType,
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items: []ServiceEndpoint{
					{
						Protocol: []uint{uint(layers.IPProtocolUDP)},
						Port:     []net.PortSpecifierString{"2222"},
					},
					{
						Protocol: []uint{uint(layers.IPProtocolTCP)},
						Port:     []net.PortSpecifierString{"2223"},
					},
				},
			},
			expectedJSON: `{"name": "ServiceEndpointTest",
                         "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
                         "type": "mfw-object-service",
                          "items": [
                              {"protocol": [17], "port": ["2222"]},
                              {"protocol": [6], "port": ["2223"]}]}`,
		},
		{
			name: "ServiceEndpointTest with port ranges",
			object: Object{
				Name:        "ServiceEndpointTest with port ranges",
				Description: "Description",
				Type:        ServiceEndpointObjectType,
				ID:          "702d4c99-9599-455f-8271-215e5680f038",
				Items: []ServiceEndpoint{
					{
						Protocol: []uint{uint(layers.IPProtocolUDP)},
						Port:     []net.PortSpecifierString{"2222", "2223-2225"},
					},
				},
			},
			expectedJSON: `{"name": "ServiceEndpointTest with port ranges",
						 "id": "702d4c99-9599-455f-8271-215e5680f038",
						 "description": "Description",
						 "type": "mfw-object-service",
						 "items": [
							 {"protocol": [17], "port": ["2222", "2223-2225"]}
						 ]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.object)
			assert.NoError(t, err)
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
				Value: []string{"fd00::1/128"},
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
		{
			name: "test APPLICATION",
			json: `{
				"op": "==",
				"type": "APPLICATION",
				"value": ["8105f355-cb98-43eb-dead-74542a524abb"]
			}`,
			shouldErr: false,
			expected: PolicyCondition{
				Op:    "==",
				CType: "APPLICATION",
				Value: []string{"8105f355-cb98-43eb-dead-74542a524abb"},
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
			name: "validJsonWithSettings",
			inputData: `{
					"id": "A1",
					"name": "B2",
					"description": "C3",
					"type": "",
					"settings": {
						"setting_field": "D4",
						"key": "value"
					}
				}`,
			wantUnmarshalledConfig: PolicyConfiguration{
				ID:          "A1",
				Name:        "B2",
				Description: "C3",
				Type:        "",
				Settings: map[string]any{
					"setting_field": "D4",
					"key":           "value",
				},
			},
			wantMarshalledJson: `{
				"id": "A1",
				"name": "B2",
				"description": "C3",
				"type": "",
				"settings": {
					"key": "value",
					"setting_field": "D4"
				}
			}`,
		},
		{
			name: "validJsonWithoutSettings",
			inputData: `{
					"id": "A1",
					"name": "B2",
					"description": "C3",
					"type": ""
				}`,
			wantUnmarshalledConfig: PolicyConfiguration{
				ID:          "A1",
				Name:        "B2",
				Description: "C3",
				Type:        "",
			},
			wantMarshalledJson: `{
				"id": "A1",
				"name": "B2",
				"description": "C3",
				"type": ""
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
