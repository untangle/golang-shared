package settingsloader

import (
	"testing"

	"github.com/untangle/golang-shared/services/settings/policy"

	"github.com/stretchr/testify/assert"
)

func TestPolicyManagerSettingsFile(t *testing.T) {

	settingskey := []string{"policy_manager"}
	t.Run("load_global_settings_file", func(t *testing.T) {
		result := policy.PolicySettings{}
		err := LoadSettingsFromURL(&result, settingskey,
			"https://raw.githubusercontent.com/untangle/mfw_schema/master/v1/policy_manager/test_settings.json",
		)

		assert.Nil(t, err, "error should be nil, but was %v", err)

		// This may change as the example settings file changes.
		assert.Equal(t, len(result.Rules), 13)
	})
}