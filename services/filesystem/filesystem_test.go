package filesystem

import (
	"io"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/untangle/golang-shared/platform"
)

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
			name:         "Open file on EOS - settings path with subdirectory exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"mnt/flash/mfw-settings/sub/dir/file.txt": {Data: []byte("deep file")},
			},
			fileName:     "/etc/config/sub/dir/file.txt",
			expectedErr:  false,
			expectedData: "deep file",
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
		{
			name:         "Open file with no path modifications and doesn't exist",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			fileName:     "dir/file.txt/",
			expectedErr:  true,
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
			name:         "Stat file on EOS - settings path with subdirectory exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"mnt/flash/mfw-settings/sub/dir/file.txt": {Data: []byte("deep file")},
			},
			fileName:     "/etc/config/sub/dir/file.txt",
			expectedErr:  false,
			expectedSize: 9,
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
		{
			name:         "Stat file with no path modifications and doesn't exist",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			fileName:     "dir/file.txt/",
			expectedErr:  true,
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

func TestGetPathOnPlatform(t *testing.T) {
	tests := []struct {
		name         string
		platformType platform.HostType
		files        fstest.MapFS
		path         string
		prefix       string
		expectedPath string
		expectedErr  bool
	}{
		{
			name:         "OpenWrt, no prefix, file exists",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"file.txt": {Data: []byte("content")},
			},
			path:         "file.txt",
			prefix:       "",
			expectedPath: "file.txt",
			expectedErr:  false,
		},
		{
			name:         "OpenWrt, with prefix, file exists",
			platformType: platform.OpenWrt,
			files: fstest.MapFS{
				"file.txt": {Data: []byte("content")},
			},
			path:         "file.txt",
			prefix:       "/tmp",
			expectedPath: "/tmp/file.txt",
			expectedErr:  false,
		},
		{
			name:         "EOS, mapped file, with prefix, file exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"root/usr/share/bctid/categories.json": {Data: []byte("categories data")},
			},
			path:         "/etc/config/categories.json",
			prefix:       "/root",
			expectedPath: "/root/usr/share/bctid/categories.json",
			expectedErr:  false,
		},
		{
			name:         "EOS, settings file, with prefix, file exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"root/mnt/flash/mfw-settings/settings.json": {Data: []byte("setting value")},
			},
			path:         "/etc/config/settings.json",
			prefix:       "/root",
			expectedPath: "/root/mnt/flash/mfw-settings/settings.json",
			expectedErr:  false,
		},
		{
			name:         "EOS, settings file, no prefix, file exists",
			platformType: platform.EOS,
			files: fstest.MapFS{
				"mnt/flash/mfw-settings/settings.json": {Data: []byte("setting value")},
			},
			path:         "/etc/config/settings.json",
			prefix:       "",
			expectedPath: "/mnt/flash/mfw-settings/settings.json",
			expectedErr:  false,
		},
		{
			name:         "File does not exist, with prefix",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			path:         "/etc/config/not-there.json",
			prefix:       "/root",
			expectedErr:  true,
		},
		{
			name:         "File does not exist, no prefix",
			platformType: platform.EOS,
			files:        fstest.MapFS{},
			path:         "/etc/config/not-there.json",
			prefix:       "",
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []FileSystemOption
			if tt.prefix != "" {
				opts = append(opts, WithPrefix(tt.prefix))
			}
			pafs := NewPlatformAwareFileSystem(tt.files, tt.platformType, opts...)
			path, err := pafs.GetPathOnPlatform(tt.path)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)
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
		{"multiple leading slashes (should remove all)", "//path/to/file", "path/to/file"},
		{"multiple trailing slashes (should only remove all)", "path/to/file//", "path/to/file"},
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

	t.Run("without options", func(t *testing.T) {
		pafs := NewPlatformAwareFileSystem(mockFS, p)

		assert.NotNil(t, pafs)
		assert.Equal(t, mockFS, pafs.FS)
		assert.Equal(t, p, pafs.platform)
		assert.Empty(t, pafs.prefix)
	})

	t.Run("with prefix option", func(t *testing.T) {
		prefix := "/tmp/root"
		pafs := NewPlatformAwareFileSystem(mockFS, p, WithPrefix(prefix))

		assert.NotNil(t, pafs)
		assert.Equal(t, mockFS, pafs.FS)
		assert.Equal(t, p, pafs.platform)
		assert.Equal(t, prefix, pafs.prefix)
	})
}
