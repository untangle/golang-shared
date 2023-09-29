package settingsutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// CopySettingsToTemp copies the file designated by original to a temp
// directory and returns that file and a function that can be used to
// remove the file and the temp directory that contains it --
// basically, use the returned string to refer to the new copied
// filename and defer calling the second argument.
func CopySettingsToTemp(t testing.TB, original string) (string, func()) {
	dir, err := os.MkdirTemp("/tmp/", "unit-test")
	require.Nil(t, err)
	settingsBytes, err := os.ReadFile(original)
	require.Nil(t, err)
	newName := filepath.Base(original)
	fullPath := filepath.Join(dir, newName)
	err = os.WriteFile(fullPath, settingsBytes, 0666)
	require.Nil(t, err)
	return fullPath, func() {
		os.RemoveAll(dir)
	}
}
