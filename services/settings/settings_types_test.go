package settings

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePluginSettings(t *testing.T) {

	type testCase struct {
		SettingsObj interface{} `json:"settingsObj"`
		Valid       bool        `json:"valid"`
		Description string      `json:"description"`
	}

	raw, err := ioutil.ReadFile("./testdata/settings_types.json")
	assert.Nil(t, err, "error reading test file")

	testObject := []testCase{}

	err = json.Unmarshal(raw, &testObject)
	assert.Nil(t, err, "error unmarshalling test file")

	for testIndex, test := range testObject {
		bodyBytes, err := json.Marshal(test.SettingsObj)
		assert.Nil(t, err, "error marshalling test case %v:%v", testIndex, test.Description)

		assert.Equal(t, test.Valid, ValidateDiscoverySettings(bodyBytes), "wrong result for test case %v:%v", testIndex, test.Description)
	}
}
