package settingsloader

import (
	"bytes"
	"io"
	"net/http"

	"github.com/untangle/golang-shared/services/settings"
)

const (
	// DefaultSettingsURL is the default URL to load settings from
	DefaultSettingsURL = "https://raw.githubusercontent.com/untangle/mfw_schema/master/v1/policy_manager/test_settings.json"
)

// Load a global settings file via URL, return the settings object
func LoadSettingsFromURL(output any, key []string, url string) error {

	// Get settings from URL
	if url == "" {
		url = DefaultSettingsURL
	}

	/*
		// Load from file for now.
		content := policy.PolicySettings{}
		sFile := settings.NewSettingsFile("settings_test.json")
		err := sFile.UnmarshalSettingsAtPath(&content, "policy_manager")
	*/
	// Load settings file from URL

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if buf, err := io.ReadAll(resp.Body); err == nil {
			unmarshaller := settings.NewPathUnmarshaller(bytes.NewReader(buf))
			err = unmarshaller.UnmarshalAtPath(output, key...)
			return err
		} else {
			return err
		}
	} else {

		return err
	}
}
