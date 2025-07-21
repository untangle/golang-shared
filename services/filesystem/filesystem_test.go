package filesystem

import (
	"io"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/untangle/golang-shared/platform"
)

var originalDetectPlatform func(fs.StatFS) platform.HostType

func TestOpen(t *testing.T) {
	tests := []struct {
		name         string
		platformType platform.HostType
		files        fstest.MapFS
		fileName     string
		expectedErr  bool
		expectedData string
	}{
		{
			name:         "Open file on OpenWrt - direct path",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"testfile.txt": {Data: []byte("hello openwrt")},
			},
			fileName:     "testfile.txt",
			expectedErr:  false,
			expectedData: "hello openwrt",
		},
		{
			name:         "Open file on EOS - uniquely mapped file exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"usr/share/bctid/categories.json": {Data: []byte("categories data")},
			},
			fileName:     "/etc/config/categories.json",
			expectedErr:  false,
			expectedData: "categories data",
		},
		{
			name:         "Open file on EOS - uniquely mapped file does not exist",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			fileName:     "/etc/config/categories.json",
			expectedErr:  true,
			expectedData: "",
		},
		{
			name:         "Open file on EOS - settings path exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"mnt/flash/mfw-settings/settings.json": {Data: []byte("setting value")},
			},
			fileName:     "/etc/config/settings.json",
			expectedErr:  false,
			expectedData: "setting value",
		},
		{
			name:         "Open file on EOS - settings path does not exist",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			fileName:     "/etc/config/settings.json",
			expectedErr:  true,
			expectedData: "",
		},
		{
			name:         "Open file on Vittoria - settings path exists",
			platformType: platform.Vittoria,
			files: fstest.MapFS{
				"velocloud/settings.json": {Data: []byte("setting value")},
			},
			fileName:     "/etc/config/settings.json",
			expectedErr:  false,
			expectedData: "setting value",
		},
		{
			name:         "Open non-existent file",
			platformType: platform.OpenWrt,
			files:        fstest.MapFS{},
			fileName:     "nonexistent.txt",
			expectedErr:  true,
			expectedData: "",
		},
		{
			name:         "Open with leading slash in name",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"anotherfile.txt": {Data: []byte("another one")},
			},
			fileName:     "/anotherfile.txt",
			expectedErr:  false,
			expectedData: "another one",
		},
		{
			name:         "Open with trailing slash in name",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"dir/file.txt": {Data: []byte("file in dir")},
			},
			fileName:     "dir/file.txt/",
			expectedErr:  false,
			expectedData: "file in dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pafs := NewPlatformAwareFileSystem(tt.files, tt.platformType)
			file, err := pafs.Open(tt.fileName)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, file)
				return
			}
			require.NoError(t, err)
			defer file.Close()

			stat, err := file.Stat()
			require.NoError(t, err)
			assert.Equal(t, int64(len(tt.expectedData)), stat.Size())

			data, err := io.ReadAll(file)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedData, string(data))

		})
	}
}

func TestStat(t *testing.T) {
	tests := []struct {
		name         string
		platformType platform.HostType
		files        fstest.MapFS
		fileName     string
		expectedErr  bool
		expectedSize int64
	}{
		{
			name:         "Stat file on OpenWrt - direct path",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"testfile.txt": {Data: []byte("some data")},
			},
			fileName:     "testfile.txt",
			expectedErr:  false,
			expectedSize: 9,
		},
		{
			name:         "Stat file on EOS - uniquely mapped file",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"usr/share/bctid/categories.json": {Data: []byte("categories")},
			},
			fileName:     "/etc/config/categories.json",
			expectedErr:  false,
			expectedSize: 10,
		},
		{
			name:         "Stat file on EOS - uniquely mapped file does not exist",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			fileName:     "/etc/config/categories.json",
			expectedErr:  true,
			expectedSize: 0,
		},
		{
			name:         "Stat non-existent file",
			platformType: platform.OpenWrt,
			files:        fstest.MapFS{},
			fileName:     "nonexistent.txt",
			expectedErr:  true,
			expectedSize: 0,
		},
		{
			name:         "Stat directory",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"mydir/":         {},
				"mydir/file.txt": {Data: []byte("content")},
			},
			fileName:     "mydir",
			expectedErr:  false,
			expectedSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pafs := NewPlatformAwareFileSystem(tt.files, tt.platformType)
			info, err := pafs.Stat(tt.fileName)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSize, info.Size())
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no leading or trailing slash", "path/to/file", "path/to/file"},
		{"leading slash", "/path/to/file", "path/to/file"},
		{"trailing slash", "path/to/file/", "path/to/file"},
		{"leading and trailing slash", "/path/to/file/", "path/to/file"},
		{"root slash", "/", ""},
		{"just a file name", "file.txt", "file.txt"},
		{"multiple leading slashes (should only remove first)", "//path/to/file", "/path/to/file"},
		{"multiple trailing slashes (should only remove first)", "path/to/file//", "path/to/file/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, sanitizePath(tt.input))
		})
	}
}

func TestNewPlatformAwareFileSystem(t *testing.T) {
	mockFS := fstest.MapFS{}
	p := platform.OpenWrt
	pafs := NewPlatformAwareFileSystem(mockFS, p)

	assert.NotNil(t, pafs)
	assert.Equal(t, mockFS, pafs.FS)
	assert.Equal(t, p, pafs.platform)
}
