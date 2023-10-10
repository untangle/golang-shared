package settingsloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGlobalSettingsFile(t *testing.T) {

	settingskey := []string{"policy_manager"}
	t.Run("load_global_settings_file", func(t *testing.T) {
		var result interface{}
		err := LoadSettingsFromURL("result",
			"https://raw.githubusercontent.com/untangle/mfw_schema/master/v1/policy_manager/test_settings.json",
			settingskey)
		assert.Nil(t, err, "error should be nil, but was %v", err)
		assert.Equal(t, result, "result")
	})
}
