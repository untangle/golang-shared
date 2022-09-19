package data

import (
        "path"
        "runtime"
)

// GetCommonTestDataDirectory returns data directory for tests.
func GetCommonTestDataDirectory() string {
        _, ourPath, _, _ := runtime.Caller(0)
        return path.Join(path.Dir(ourPath), "testdata")
}

// GetTestFileLocation returns the full path of a specific test file.
func GetTestFileLocation(testFile string) string {
        return path.Join(GetCommonTestDataDirectory(), testFile)
}
