package settings

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/plugins/types"
	"github.com/untangle/golang-shared/testing/mocks"
)

type fileExistsFake struct {
	fileExists bool
}

func (f *fileExistsFake) doesExist(fname string) bool {
	return f.fileExists
}
func TestFilenameLocator(t *testing.T) {
	existFake := &fileExistsFake{}
	locator := FilenameLocator{
		fileExists: existFake.doesExist}
	tests := []struct {
		filename         string
		outputFileExists bool
		platform         types.Platform
		returnValue      string
		returnErr        error
	}{
		// {
		// 	// File exists, should get the same back.
		// 	filename:         "/etc/config/settings.json",
		// 	outputFileExists: true,
		// 	platform:         types.OpenWrt,
		// 	returnValue:      "/etc/config/settings.json",
		// },
		// {
		// 	// In OpenWRT mode, no translation since the defaults are openwrt paths.
		// 	// returns with no error
		// 	filename:         "/usr/share/geoip",
		// 	outputFileExists: true,
		// 	platform:         types.OpenWrt,
		// 	returnValue:      "/usr/share/geoip",
		// },
		// {
		// 	// In Native mode, do translation
		// 	filename:         "/etc/config/appstate.json",
		// 	outputFileExists: true,
		// 	platform:         types.EOS,
		// 	returnValue:      "/mnt/flash/mfw-settings/appstate.json",
		// },
		// {
		// 	// In Native mode, do translation, file is not there so return error
		// 	filename:         "/etc/config/settings.json",
		// 	outputFileExists: false,
		// 	platform:         types.EOS,
		// 	returnValue:      "/mnt/flash/mfw-settings/settings.json",
		// 	returnErr:        fmt.Errorf("no file at path: /mnt/flash/mfw-settings/settings.json"),
		// },
		{
			// In Native mode, do translation, file exists, not error
			filename:         "/etc/config/appstate.json",
			outputFileExists: false,
			platform:         types.EOS,
			returnValue:      "/mnt/flash/mfw-settings/appstate.json",
		},
		// { // Native mode, no translation, file exists
		// 	filename:         "/usr/share/bctid/categories.json",
		// 	outputFileExists: false,
		// 	Platform:         types.OpenWrt,
		// 	returnValue:      "/usr/share/bctid/categories.json",
		// },
		// { // Native mode, New file not there, return same thing
		// 	filename:         "/tmp/captivesocket",
		// 	outputFileExists: false,
		// 	Platform:         types.EOS,
		// 	returnValue:      "/tmp/captivesocket",
		// 	returnErr:        fmt.Errorf("no file at path: /tmp/captivesocket"),
		// },
		// { // OpenWRT mode, translate. New file not there, return same thing
		// 	filename:         "/tmp/captivesocket",
		// 	outputFileExists: false,
		// 	Platform:         types.OpenWrt,
		// 	returnValue:      "/tmp/captivesocket",
		// 	returnErr:        fmt.Errorf("no file at path: /tmp/captivesocket"),
		// },
		// { // Native mode, translate. New file not there, return error
		// 	filename:         "/etc/config/categories.json",
		// 	outputFileExists: false,
		// 	Platform:         types.EOS,
		// 	returnValue:      "/usr/share/bctid/categories.json",
		// 	returnErr:        fmt.Errorf("no file at path: /usr/share/bctid/categories.json"),
		// },
	}

	for _, test := range tests {
		existFake.fileExists = test.outputFileExists
		locator.platform = test.platform
		result, err := locator.LocateFile(test.filename)
		assert.Equal(t, test.returnValue, result)
		if test.returnErr == nil {
			assert.NoError(t, err)
		} else {
			assert.Regexp(t, test.returnErr.Error(), err.Error(),
				"errors should match")
			matchingError := &NoFileAtPath{}
			assert.True(t, errors.As(err, &matchingError))
		}
	}

}

func TestFileExists(t *testing.T) {
	thisFile, err := os.Executable()
	logger = mocks.NewMockLogger()
	assert.NoError(t, err)
	assert.True(t, FileExists(thisFile))
	assert.False(t,
		FileExists("/some-file/that-should/definitely-not/exist-anywhere"))
}
