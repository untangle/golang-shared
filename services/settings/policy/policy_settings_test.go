package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
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

	// Spot check a plugin setting.
	assert.EqualValues(t, result, policySettings["threatprevention"]["Teachers"])
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
