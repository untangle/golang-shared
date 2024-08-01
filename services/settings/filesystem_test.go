package settings

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/testing/mocks"
)

type fileExistsFake struct {
	rvals []bool
}

func (f *fileExistsFake) doesExist(fname string) bool {
	rval := f.rvals[0]
	f.rvals = f.rvals[1:]
	return rval
}
func TestFilenameLocator(t *testing.T) {
	existFake := &fileExistsFake{}
	locator := FilenameLocator{
		fileExists: existFake.doesExist}
	tests := []struct {
		filename     string
		existResults []bool // Indicates (file initially exists, In kernel/openwrt mode , new file exists)
		returnValue  string
		returnErr    error
	}{
		{
			filename:     "/etc/config/settings.json",
			existResults: []bool{false, false, true},
			returnValue:  "/mnt/flash/mfw-settings/settings.json",
		},
		{
			filename:     "/usr/share/geoip",
			existResults: []bool{false, false, true},
			returnValue:  "/mfw/usr/share/geoip",
		},
		{
			filename:     "/etc/config/appstate.json",
			existResults: []bool{false, false, true},
			returnValue:  "/mnt/flash/mfw-settings/appstate.json",
		},
		{
			filename:     "/etc/config/settings.json",
			existResults: []bool{true, true, true},
			returnValue:  "/etc/config/settings.json",
		},
		{
			filename:     "/etc/config/appstate.json",
			existResults: []bool{true, true, true},
			returnValue:  "/etc/config/appstate.json",
		},
		{
			filename:     "/etc/config/appstate.json",
			existResults: []bool{false, false, false},
			returnValue:  "/mnt/flash/mfw-settings/appstate.json",
			returnErr:    fmt.Errorf("unable to find config file: /mnt/flash/mfw-settings/appstate.json"),
		},
	}

	for _, test := range tests {
		existFake.rvals = test.existResults
		result, err := locator.LocateFile(test.filename)
		assert.Equal(t, result, test.returnValue)
		if test.returnErr == nil {
			assert.NoError(t, err)
		} else {
			assert.Regexp(t, test.returnErr.Error(), err.Error(),
				"errors should match")
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
