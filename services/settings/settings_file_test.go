package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSettingsFilePath = "./testdata/settings.json"

// Test basic settings -- we do not need to go in depth here since
// this object just wraps the PathUnmarshaller object.
func TestSettings(t *testing.T) {
	type settingsObject struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}
	s := NewSettingsFile(testSettingsFilePath)
	value := settingsObject{}
	err := s.UnmarshalSettingsAtPath(&value, "a", "b")
	assert.Nil(t, err)
	assert.Equal(
		t,
		value,
		settingsObject{
			Foo: "hello",
			Bar: 1})
}

// Tests retrieving all settings
func TestGetAllSettings(t *testing.T) {
	s := NewSettingsFile(testSettingsFilePath)
	settings, err := s.GetAllSettings()
	assert.NoError(t, err)

	if !assert.Contains(t, settings, "a") {
		return
	}

	aNoType := settings["a"]
	a := aNoType.(map[string]interface{})
	if !assert.Contains(t, a, "b") {
		return
	}

	bNoType := a["b"]
	b := bNoType.(map[string]interface{})
	if !assert.Contains(t, b, "foo") || !assert.Contains(t, b, "bar") {
		return
	}
}

// Uses a script that just points at the testdata/settings.json to generate a backup
func TestGenerateBackup(t *testing.T) {
	s := NewSettingsFile(testSettingsFilePath)
	name, data, err := s.GenerateBackupFile("./testdata/testBackupGeneration.sh")
	assert.NoError(t, err)
	assert.Equal(t, name, "./testdata/settings.json")

	assert.Contains(t, string(data), "\"a\":", "\"b\":", "\"foo\": \"hello\"", "\"bar\": 1")
}
