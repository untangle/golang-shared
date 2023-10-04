package settings

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSyncSettingsJsonOutput(t *testing.T) {
	t.Run("get_json_output_from_success_output", func(t *testing.T) {
		output, err := os.ReadFile("testdata/sync_settings_output_success")
		assert.Nil(t, err)

		jsonOutput, mapOutput, err := parseSyncSettingsJsonOutput(output)
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

		jsonOutput, mapOutput, err := parseSyncSettingsJsonOutput(output)
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

		jsonOutput, mapOutput, err := parseSyncSettingsJsonOutput(output)
		expectedString := ""
		var expectedMap map[string]any

		assert.Equal(t, expectedString, jsonOutput)
		assert.Equal(t, expectedMap, mapOutput)
		assert.Contains(t, err.Error(), "parse sync-settings output error:")
	})
}

func TestWriteSettingsFileJSON(t *testing.T) {
	// Create a temporary test file for writing
	testFile, err := os.CreateTemp("", "test_settings.json")
	assert.NoError(t, err, "Failed to create temporary file test_settings.json")
	defer os.Remove(testFile.Name()) // Remove the temporary file after the test

	// Define the JSON data to write to the file
	jsonData := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	// Call the function to write JSON data to the test file
	success, writeErr := writeSettingsFileJSON(jsonData, testFile)
	assert.NoError(t, writeErr, "Error writing JSON to file")

	// Ensure that the file was written successfully
	assert.True(t, success, "writeSettingsFileJSON returned false for success")

	// Read the content of the written file to verify
	fileContent, readErr := os.ReadFile(testFile.Name())
	assert.NoError(t, readErr, "Error reading file content")

	// Unmarshal the file content to compare with the original JSON data
	var parsedData map[string]interface{}
	unmarshalErr := json.Unmarshal(fileContent, &parsedData)
	assert.NoError(t, unmarshalErr, "Error unmarshalling file content")

	// Compare the parsed data with the original JSON data
	if !jsonEqual(parsedData, jsonData) {
		t.Fatalf("Written JSON data does not match expected data")
	}
}

// jsonEqual checks if two JSON objects are equal
func jsonEqual(a, b map[string]interface{}) bool {
	aJSON, errA := json.Marshal(a)
	bJSON, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return string(aJSON) == string(bJSON)
}
