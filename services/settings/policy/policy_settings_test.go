package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
)

func TestGetAllPolicyConfigurationSettings(t *testing.T) {
	// Good settings file expect it to work
	settingsFile := settings.NewSettingsFile("./testdata/test_settings.json")
	policySettings, err := getAllPolicyConfigurationSettings(settingsFile)
	assert.Nil(t, err)
	assert.NotNil(t, policySettings)
	assert.Equal(t, 3, len(policySettings["threatprevention"]))
	assert.Equal(t, 1, len(policySettings["webfilter"]))
	assert.Equal(t, 1, len(policySettings["geoip"]))
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
