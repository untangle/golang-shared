package settingsloader

import (
	"testing"
)

func TestPolicyManagerSettingsFile(t *testing.T) {
	t.Skip("broken by upstream schema.")
	// commented because the linter doesn't like this.

	// settingskey := []string{"policy_manager"}

	// t.Run("load_global_settings_file", func(t *testing.T) {
	// 	result := policy.PolicySettings{}
	// 	err := LoadSettingsFromURL(&result, settingskey,
	// 		"https://raw.githubusercontent.com/untangle/mfw_schema/master/v1/policy_manager/test_settings.json",
	// 	)

	// 	assert.NoError(t, err)

	// 	// This may change, when the example settings file changes.
	// 	assert.Equal(t, len(result.Rules), 13)
	// })
}
