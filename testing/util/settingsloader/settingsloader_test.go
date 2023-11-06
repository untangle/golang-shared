package settingsloader

import (
	"testing"
)

func TestPolicyManagerSettingsFile(t *testing.T) {
	// Skip test, until port ranges have been added to loader. See MFW-3775
	t.Skip("Skipping test until tasks listed in MFW-3775 are completed")

	// settingskey := []string{"policy_manager"}
	// t.Run("load_global_settings_file", func(t *testing.T) {
	// 	result := policy.PolicySettings{}
	// 	err := LoadSettingsFromURL(&result, settingskey,
	// 		"https://raw.githubusercontent.com/untangle/mfw_schema/master/v1/policy_manager/test_settings.json",
	// 	)
	// }
}
