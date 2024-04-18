package util

import (
	"math/rand"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testTarGz string = "./testdata/settings_test_unzip.backup.tar.gz"
	testTar   string = "./testdata/settings_test_unzip.backup.tar"
	testJson  string = "./testdata/expected_settings.json"
)

// Test RandomizeSlice
func TestRandomizeSlice(t *testing.T) {
	seeds := []int64{10, 20, 30}
	tests := []struct {
		name   string
		actual []interface{}
	}{
		{
			name:   "Randomize empty slice",
			actual: []interface{}{},
		},
		{
			name:   "Randomize single-element slice",
			actual: []interface{}{1},
		},
		{
			name:   "Randomize multiple-element slice",
			actual: []interface{}{1, 2, 3, 4, 5},
		},
	}

	for _, seed := range seeds {
		rand.Seed(seed)

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Make a copy of the input slice to compare with the result
				inputCopy := make([]interface{}, len(test.actual))
				copy(inputCopy, test.actual)

				RandomizeSlice(test.actual)

				// Check if the input slice is not equal to the expected slice.
				// lists of size <= 1 excluded.
				if len(test.actual) > 1 && reflect.DeepEqual(test.actual, inputCopy) {
					t.Errorf("Expected slice to be randomized, but it's the same as the original: %v", test.actual)
				}
			})
		}
	}
}

func TestWaitTimeout(t *testing.T) {
	tests := []struct {
		name      string
		timeout   time.Duration
		expectErr bool
	}{
		{
			name:      "No timeout",
			timeout:   2 * time.Second,
			expectErr: false,
		},
		{
			name:      "Timeout",
			timeout:   500 * time.Millisecond,
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var wg sync.WaitGroup

			// Simulate some work
			wg.Add(1)
			go func() {
				time.Sleep(1 * time.Second)
				wg.Done()
			}()

			start := time.Now()
			err := WaitGroupDoneOrTimeout(&wg, test.timeout)
			elapsed := time.Since(start)

			if test.expectErr && !err {
				t.Error("Expected timeout error, but got none")
			} else if !test.expectErr && err {
				t.Error("Expected no timeout error, but got one")
			}

			if test.expectErr && elapsed < test.timeout {
				t.Errorf("Expected elapsed time (%v) to be greater than or equal to the timeout (%v)", elapsed, test.timeout)
			}
		})
	}
}

// Tests extracting a tar from an array of bytes
func TestExtractSettingsFromTar(t *testing.T) {
	// Test that the settings can be extracted
	// Read in test tar and make sure the settings can be pulled out of it
	tarDataGz, _ := os.ReadFile(testTarGz)
	tarData, _ := os.ReadFile(testTar)
	expectedSettingsData, _ := os.ReadFile(testJson)

	foundFiles, err := ExtractFilesFromTar(tarData, false, "settings.json")
	assert.NoError(t, err)

	foundFilesGz, err := ExtractFilesFromTar(tarDataGz, true, "settings.json", "fakeFile.fake")
	assert.NoError(t, err)

	assert.Equal(t, string(expectedSettingsData), string(foundFiles["settings.json"]))
	assert.Equal(t, expectedSettingsData, foundFilesGz["settings.json"])

	// Test that a file that didn't exist wasn't found and caused not issues
	assert.NotContains(t, foundFilesGz, "fakeFile.fake")

	// Test that an unzipped tar with isGzip fails expectedly
	_, err = ExtractFilesFromTar(tarData, true, "settings.json")
	assert.Error(t, err)

}

// Tests StringArrayToDB
func TestStringArrayToDB(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Empty input",
			input:    []string{},
			expected: "",
		},
		{
			name:     "Single element",
			input:    []string{"test"},
			expected: "test",
		},
		{
			name:     "Multiple elements",
			input:    []string{"test", "test2", "test3"},
			expected: "test|test2|test3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := StringArrayToDB(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

// Test DecodeAttribute functions
func TestDecodeAttribute(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		errorExpected bool
	}{
		{
			name:          "Null string input",
			input:         "",
			expected:      "",
			errorExpected: false,
		},
		{
			name:          "Actual value test",
			input:         "Testingval",
			expected:      "",
			errorExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.errorExpected {
				actual, err := DecodeAttribute(test.input)
				assert.Equal(t, test.expected, actual)
				assert.NoError(t, err)
			} else {
				_, err := DecodeAttribute(test.input)
				assert.Error(t, err)
			}
		})
	}
}
