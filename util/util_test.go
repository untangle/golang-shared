package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTarGz string = "./testdata/settings_test_unzip.backup.tar.gz"
	testTar   string = "./testdata/settings_test_unzip.backup.tar"
	testJson  string = "./testdata/expected_settings.json"
)

// Tests extracting a tar from an array of bytes
func TestExtractSettingsFromTar(t *testing.T) {
	// Test that the settings can be extracted
	// Read in test tar and make sure the settings can be pulled out of it
	tarDataGz, _ := os.ReadFile(testTarGz)
	tarData, _ := os.ReadFile(testTar)
	expectedSettingsData, _ := os.ReadFile(testJson)

	foundFiles, err := ExtractFilesFromTar(tarData, false, "settings.json")
	assert.NoError(t, err)

	foundFilesGz, err := ExtractFilesFromTar(tarDataGz, true, "settings.json", "fakeFile.fake")
	assert.NoError(t, err)

	assert.Equal(t, string(expectedSettingsData), string(foundFiles["settings.json"]))
	assert.Equal(t, expectedSettingsData, foundFilesGz["settings.json"])

	// Test that a file that didn't exist wasn't found and caused not issues
	assert.NotContains(t, foundFilesGz, "fakeFile.fake")

	// Test that an unzipped tar with isGzip fails expectedly
	_, err = ExtractFilesFromTar(tarData, true, "settings.json")
	assert.Error(t, err)

}
