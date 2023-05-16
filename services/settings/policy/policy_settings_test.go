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
	tpPolicies := GetPolicyPluginSettings(settingsFile, "threatprevention")
	assert.Equal(t, 3, len(tpPolicies))
	webFilterPolicies := GetPolicyPluginSettings(settingsFile, "webfilter")
	assert.Equal(t, 1, len(webFilterPolicies))
	geoIPPolicies := GetPolicyPluginSettings(settingsFile, "geoip")
	assert.Equal(t, 1, len(geoIPPolicies))
}
