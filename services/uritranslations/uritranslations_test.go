package uritranslations

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURI(t *testing.T) {
	// Build sample maps for testing
	uriMap = make(map[string]*url.URL)
	hostMap = make(map[string]*url.URL)

	// Add test URIs to the maps
	testURI := "https://example.com/test1"
	testURI1, _ := url.Parse(testURI)
	testURI2, _ := url.Parse("https://example.com/test2")
	uriMap["https://example.com/test1"] = testURI1
	uriMap["https://example.com/test2"] = testURI2
	hostMap["example.com"] = testURI1

	// Test existing URI
	result, err := GetURI("https://example.com/test1")
	assert.NoError(t, err, "Failed to get URI")

	expectedResponse := testURI
	actualResponse := result
	assert.Equal(t, expectedResponse, actualResponse)

	// Test non-existent URI
	_, err = GetURI("https://example.com/test3")
	assert.Error(t, err, "Failed to get error")
}

func TestGetURIWithPath(t *testing.T) {
	// Build some sample maps for testing
	uriMap = make(map[string]*url.URL)
	hostMap = make(map[string]*url.URL)

	// Add some test URIs to the maps
	url1 := "https://example.com/test1/api"
	testURI1, _ := url.Parse(url1)
	testURI2, _ := url.Parse("https://example.com/test2")
	uriMap["https://example.com/test1"] = testURI1
	uriMap["https://example.com/test2"] = testURI2
	hostMap["example.com"] = testURI1

	// Test existing URI with path substitution
	result, err := GetURIWithPath("https://example.com/test1/api")
	assert.NoError(t, err, "Failed to get URI with path")

	expectedResponse := url1
	actualResponse := result
	assert.Equal(t, expectedResponse, actualResponse)

	// Test non-existent URI
	result, err = GetURIWithPath("https://example.com/test3")
	assert.Nil(t, err, "Failed to get error")
}

func TestGetURITranslation(t *testing.T) {
	// Create test URIs
	testURIs := []*url.URL{
		{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/test1",
		},
		{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/test2",
		},
	}

	// Initialize maps
	uriMap = make(map[string]*url.URL)
	hostMap = make(map[string]*url.URL)

	// Populate the maps
	for _, uri := range testURIs {
		uriMap[uri.String()] = uri
		hostMap[uri.Host] = uri
	}

	tests := []struct {
		name         string
		uri          string
		path         bool
		expected     string
		shouldFail   bool
		functionName string
	}{
		{
			name:         "ExistingURI",
			uri:          "https://example.com/test1",
			path:         false,
			expected:     "https://example.com/test1",
			functionName: "GetURI",
		},
		{
			name:         "NonExistentURI",
			uri:          "https://example.com/test3",
			path:         false,
			shouldFail:   true,
			functionName: "GetURI",
		},
		{
			name:         "ExistingURIWithPath",
			uri:          "https://example.com/test1/somepath",
			path:         true,
			expected:     "https://example.com/test1/somepath",
			functionName: "GetURIWithPath",
		},
		{
			name:         "NonExistentURIWithPath",
			uri:          "https://example.com/test3/somepath",
			path:         true,
			shouldFail:   true,
			functionName: "GetURIWithPath",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := getURITranslation(test.uri, test.path)
			if test.shouldFail {
				if test.functionName == "GetURI" {
					assert.Error(t, err)
				} else if test.functionName == "GetURIWithPath" {
					assert.Nil(t, err, "Failed to get error")
				}

			} else {
				assert.NoError(t, err)

				expectedResponse := test.uri
				actualResponse := result
				assert.Equal(t, expectedResponse, actualResponse)
			}
		})
	}
}
