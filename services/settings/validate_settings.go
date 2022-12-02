package settings

import (
	"strings"

	"github.com/untangle/golang-shared/services/settings/discovery_settings"
)

// the key for each validator is made of the segments array's elements joined by the "/" character
var validationMap = map[string]func(bytes []byte) bool{
	getValidationKey("accounts", "credentials"): exampleValidation,
	getValidationKey("discovery"):               discovery_settings.ValidateDiscoverySettings,
}

// contcatenates the segments elements and returns the key at which the validator for received settings should be found
func getValidationKey(segments ...string) string {
	return strings.Join(segments, "/")
}

// example validation function
func exampleValidation(settingsObjBytes []byte) bool {
	return true
}

// ValidateSettings - ensures the settings we received are in a valid format
//  returns true if the object is valid, false otherwise
func ValidateSettings(settingsObjBytes []byte, segments []string) bool {
	settingsPath := getValidationKey(segments...)

	validateFunc, ok := validationMap[settingsPath]
	if !ok {
		// if we have no validation, we consider anythin valid
		return true
	}

	return validateFunc(settingsObjBytes)
}
