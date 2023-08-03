package settings

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseJsonOutput(t *testing.T) {
	t.Run("get_json_output_from_success_output", func(t *testing.T) {
		output, err := os.ReadFile("testdata/sync_settings_output_success")
		assert.Nil(t, err)

		jsonOutput, mapOutput, err := parseJsonOutput(output)
		expectedString := `{"success": true, "logLines": "Syncing to system...\nNo changed files.\nDeleted files:\n/etc/config/nftables-rules.d/207-wf-rules\nDeleting files...\nCopying files...", "message": null, "traceback": null, "raisedException": null}`
		expectedMap := map[string]any{
			"success":         true,
			"logLines":        "Syncing to system...\nNo changed files.\nDeleted files:\n/etc/config/nftables-rules.d/207-wf-rules\nDeleting files...\nCopying files...",
			"message":         nil,
			"raisedException": nil,
			"traceback":       nil,
		}

		assert.Equal(t, expectedString, jsonOutput)
		assert.Equal(t, expectedMap, mapOutput)
		assert.Nil(t, err)
	})

	t.Run("get_json_output_from_fail_output", func(t *testing.T) {
		output, err := os.ReadFile("testdata/sync_settings_output_failure")
		assert.Nil(t, err)

		jsonOutput, mapOutput, err := parseJsonOutput(output)
		expectedString := `{"success": false, "logLines": "Sanitization changed settings. Saving new settings...", "message": "CaptivePortalManager: Settings verification failed!", "traceback": "  File /usr/bin/sync-settings, line 676, in <module>", "raisedException": "Missing or invalid type of captive portal setting: enabled"}`
		expectedMap := map[string]any{
			"success":         false,
			"logLines":        "Sanitization changed settings. Saving new settings...",
			"message":         "CaptivePortalManager: Settings verification failed!",
			"raisedException": "Missing or invalid type of captive portal setting: enabled",
			"traceback":       "  File /usr/bin/sync-settings, line 676, in <module>",
		}

		assert.Equal(t, expectedString, jsonOutput)
		assert.Equal(t, expectedMap, mapOutput)
		assert.Nil(t, err)
	})

	t.Run("get_json_output_from_fail_output_no_json", func(t *testing.T) {
		output, err := os.ReadFile("testdata/sync_settings_output_failure_no_json")
		assert.Nil(t, err)

		jsonOutput, mapOutput, err := parseJsonOutput(output)
		expectedString := ""
		var expectedMap map[string]any

		assert.Equal(t, expectedString, jsonOutput)
		assert.Equal(t, expectedMap, mapOutput)
		assert.Equal(t, "invalid character 'K' looking for beginning of value", err.Error())
	})
}
