package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/testing/util/settingsutil"
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

func TestSetSettingsNoSync(t *testing.T) {
	tempfile, cleanup := settingsutil.CopySettingsToTemp(t, "testdata/writes.json")
	defer cleanup()
	sf := NewSettingsFile(tempfile)
	type example struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	x := example{
		Name:  "hello",
		Value: 1,
	}
	assert.Nil(t, sf.SetSettingsNoSync([]string{"some", "path"}, x))
	y := example{}
	assert.Nil(t, sf.UnmarshalSettingsAtPath(&y, "some", "path"))
	assert.EqualValues(t,
		example{
			Name:  "hello",
			Value: 1,
		},
		y)

}

// Uses a script that just points at the testdata/settings.json to generate a backup
func TestGenerateBackup(t *testing.T) {
	s := NewSettingsFile(testSettingsFilePath)

	// Since the pipeline runs this code in in containers that are both ash/bash
	// just hack GenerateBackupFile to run a command they can both use.
	// The echo command simulates a scripts output
	name, data, err := s.GenerateBackupFile("echo", "Backup location: ./testdata/settings.json")
	assert.NoError(t, err)
	assert.Equal(t, name, "./testdata/settings.json")

	assert.Contains(t, string(data), "\"a\":", "\"b\":", "\"foo\": \"hello\"", "\"bar\": 1")
}
