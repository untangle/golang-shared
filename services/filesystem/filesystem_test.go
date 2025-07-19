package filesystem

import (
	"testing"
)

// type fileExistsFake struct {
// 	rvals []bool
// }

//	func (f *fileExistsFake) doesExist(fname string) bool {
//		rval := f.rvals[0]
//		f.rvals = f.rvals[1:]
//		return rval
//	}
func TestFilenameLocator(t *testing.T) {
	// existFake := &fileExistsFake{}
	// locator := FilenameLocator{
	// 	fileExists: existFake.doesExist}
	// tests := []struct {
	// 	filename     string
	// 	existResults []bool // Indicates (input file exists, OnEOSPlatform, File exists after translation)
	// 	returnValue  string
	// 	returnErr    error
	// }{
	// 	{
	// 		// File exists, should get the same back.
	// 		filename:     "/etc/config/settings.json",
	// 		existResults: []bool{true},
	// 		returnValue:  "/etc/config/settings.json",
	// 	},
	// 	{
	// 		// In OpenWRT mode, no translation since the defaults are openwrt paths.
	// 		// returns with no error
	// 		filename:     "/usr/share/geoip",
	// 		existResults: []bool{false, false, true},
	// 		returnValue:  "/usr/share/geoip",
	// 	},
	// 	{
	// 		// In Native mode, do translation
	// 		filename:     "/etc/config/appstate.json",
	// 		existResults: []bool{false, true, true},
	// 		returnValue:  "/mnt/flash/mfw-settings/appstate.json",
	// 	},
	// 	{
	// 		// In Native mode, do translation, file is not there so return error
	// 		filename:     "/etc/config/settings.json",
	// 		existResults: []bool{false, true, false},
	// 		returnValue:  "/mnt/flash/mfw-settings/settings.json",
	// 		returnErr:    fmt.Errorf("no file at path: /mnt/flash/mfw-settings/settings.json"),
	// 	},
	// 	{
	// 		// In Native mode, do translation, file exists, not error
	// 		filename:     "/etc/config/appstate.json",
	// 		existResults: []bool{false, true, true},
	// 		returnValue:  "/mnt/flash/mfw-settings/appstate.json",
	// 	},
	// 	{ // Native mode, no translation, file exists
	// 		filename:     "/usr/share/bctid/categories.json",
	// 		existResults: []bool{false, false, true},
	// 		returnValue:  "/usr/share/bctid/categories.json",
	// 	},
	// 	{ // Native mode, New file not there, return same thing
	// 		filename:     "/tmp/captivesocket",
	// 		existResults: []bool{false, true, false},
	// 		returnValue:  "/tmp/captivesocket",
	// 		returnErr:    fmt.Errorf("no file at path: /tmp/captivesocket"),
	// 	},
	// 	{ // OpenWRT mode, translate. New file not there, return same thing
	// 		filename:     "/tmp/captivesocket",
	// 		existResults: []bool{false, false, false},
	// 		returnValue:  "/tmp/captivesocket",
	// 		returnErr:    fmt.Errorf("no file at path: /tmp/captivesocket"),
	// 	},
	// 	{ // Native mode, translate. New file not there, return error
	// 		filename:     "/etc/config/categories.json",
	// 		existResults: []bool{false, true, false},
	// 		returnValue:  "/usr/share/bctid/categories.json",
	// 		returnErr:    fmt.Errorf("no file at path: /usr/share/bctid/categories.json"),
	// 	},
	// }

	// for _, test := range tests {
	// 	existFake.rvals = test.existResults
	// 	result, err := locator.LocateFile(test.filename)
	// 	assert.Equal(t, result, test.returnValue)
	// 	if test.returnErr == nil {
	// 		assert.NoError(t, err)
	// 	} else {
	// 		assert.Regexp(t, test.returnErr.Error(), err.Error(),
	// 			"errors should match")
	// 		matchingError := &NoFileAtPath{}
	// 		assert.True(t, errors.As(err, &matchingError))
	// 	}
	// }

}

// func TestFileExists(t *testing.T) {
// 	thisFile, err := os.Executable()
// 	logger = mocks.NewMockLogger()
// 	assert.NoError(t, err)
// 	assert.True(t, FileExists(thisFile))
// 	assert.False(t,
// 		FileExists("/some-file/that-should/definitely-not/exist-anywhere"))
// }
