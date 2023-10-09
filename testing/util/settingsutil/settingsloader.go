package settingsutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/untangle/golang-shared/services/settings"
)

// Load a global settings file via URL, return the settings object
func LoadSettingsFromURL(output interface{}, url string, key []string) error {

	// Get settings from URL
	if url == "" {
		return errors.New("No URL provided")
	}
	// Load settings file
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
