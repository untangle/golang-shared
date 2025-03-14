package dpi

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestLoadDpiJson struct {
	suite.Suite
	manager *DpiConfigManager
}

// TestDpiConfigSuite runs the DPI config manager test suite.
func TestLoadDpiJsonSuite(t *testing.T) {
	suite.Run(t, &TestLoadDpiJson{})
}

func (suite *TestLoadDpiJson) SetupTest() {
	// Log the test name.
	suite.T().Logf("Starting test: %s", suite.T().Name())
	// Create a fresh DPI config manager for each test.
	suite.manager = NewDpiConfigManager()
}

// TestLoadConfig_Valid tests loading valid JSON from an io.Reader.
func (suite *TestLoadDpiJson) TestLoadConfig_Valid() {
	sampleJSON := `{
		"0description": "Dpi",
		"0version": "1.0",
		"vendor-attributes": [
			"file",
			"filename"
		],
		"categories": {
			"best-effort": 1,
			"enterprise": 3,
			"general": 2,
			"real-time": 4
    	},
		"services": {
			"audio-video": 40,
			"chat": 20,
			"default": 1,
			"file-transfer": 30,
			"networking": 60,
			"peer-to-peer": 50,
			"software-update": 70
		},
		"applications": {
			"zoom": {
            "family": "Instant Messaging",
            "tag": [
                "aetls",
                "audio_chat",
                "cloud_services",
                "enterprise",
                "im_mc",
                "video_chat",
                "voip"
            ],
            "id": 3,
            "service-category": {
                "audio-video": "real-time",
                "chat": "general",
                "default": "enterprise",
                "file-transfer": "enterprise"
            },
            "vendor-id": 2928,
            "vendor-service-attributes": {
                "service_id": {
                    "id": 300,
                    "type": "uint32",
                    "value-service": {
                        "2": "chat",
                        "5": "file-transfer",
                        "8": "audio-video",
                        "9": "default"
                    }
                }
            }
          }
		}
	}`

	err := suite.manager.LoadConfig(strings.NewReader(sampleJSON))
	suite.NoError(err, "LoadConfig should not return an error for valid JSON")

	// Verify metadata.
	meta := suite.manager.GetMetaData()
	suite.Equal("Dpi", meta.Description)
	suite.Equal("1.0", meta.Version)

	// Verify categories.
	cats := suite.manager.GetCategories()
	suite.Contains(cats, "best-effort")
	suite.Equal(1, cats["best-effort"])
	suite.Contains(cats, "enterprise")
	suite.Equal(3, cats["enterprise"])
	suite.Contains(cats, "general")
	suite.Equal(2, cats["general"])
	suite.Contains(cats, "real-time")
	suite.Equal(4, cats["real-time"])

	// Verify services.
	services := suite.manager.GetServices()
	suite.Contains(services, "audio-video")
	suite.Equal(40, services["audio-video"])
	suite.Contains(services, "chat")
	suite.Equal(20, services["chat"])
	suite.Contains(services, "default")
	suite.Equal(1, services["default"])
	suite.Contains(services, "file-transfer")
	suite.Equal(30, services["file-transfer"])
	suite.Contains(services, "networking")
	suite.Equal(60, services["networking"])
	suite.Contains(services, "peer-to-peer")
	suite.Equal(50, services["peer-to-peer"])
	suite.Contains(services, "software-update")
	suite.Equal(70, services["software-update"])

	// Verify the application.
	applications := suite.manager.GetApplications()
	app, found := applications[3]
	suite.True(found, "Application with ID 3 should be found")
	suite.Equal("zoom", app.Name)
	suite.Equal("Instant Messaging", app.Family)
	expectedTags := []string{"aetls", "audio_chat", "cloud_services", "enterprise", "im_mc", "video_chat", "voip"}
	suite.ElementsMatch(expectedTags, app.Tag)
	suite.Equal(3, app.ID)
	suite.Equal(2928, app.VendorID)

	// Verify the service-category mapping for the application.
	expectedServiceCategory := map[string]string{
		"audio-video":   "real-time",
		"chat":          "general",
		"default":       "enterprise",
		"file-transfer": "enterprise",
	}
	suite.Equal(expectedServiceCategory, app.ServiceCategory)

	// Verify vendor-service-attributes.
	vsa, ok := app.VendorServiceAttributes["service_id"]
	suite.True(ok, "Vendor attribute 'service_id' should exist")
	suite.Equal(300, vsa.ID)
	suite.Equal("uint32", vsa.Type)
	expectedValueService := map[string]string{
		"2": "chat",
		"5": "file-transfer",
		"8": "audio-video",
		"9": "default",
	}
	suite.Equal(expectedValueService, vsa.ValueService)
}

// TestLoadConfigFromFile_Valid tests loading valid JSON from a file.
func (suite *TestLoadDpiJson) TestLoadConfigFromFile_Valid() {
	sampleJSON := `{
		"0description": "Dpi",
		"0version": "2.0",
		"vendor-attributes": ["file"],
		"categories": {"enterprise": 3},
		"services": {"default": 1},
		"applications": {
			"random": {
				"family": "test",
				"tag": ["Tag"],
				"id": 202,
				"service-category": {"default": "enterprise"},
				"vendor-id": 999,
				"vendor-service-attributes": {}
			}
		}
	}`

	// Create a temporary file with the sample JSON.
	tmpFile, err := os.CreateTemp("", "dpi_config_*.json")
	suite.Require().NoError(err, "should create temporary file")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(sampleJSON))
	suite.Require().NoError(err, "should write sample JSON to file")
	suite.Require().NoError(tmpFile.Close(), "should close temporary file")

	err = suite.manager.LoadConfigFromFile(tmpFile.Name())
	suite.NoError(err, "LoadConfigFromFile should not return an error for valid file")
}

// TestLoadConfig_InvalidJSON tests that invalid JSON input returns an error.
func (suite *TestLoadDpiJson) TestLoadConfig_InvalidJSON() {
	invalidJSON := `{
	"0description": "Invalid Config", 
	"0version": "1.0"
	"categories": {"general": 2}
	}` // Missing comma after "1.0"
	err := suite.manager.LoadConfig(strings.NewReader(invalidJSON))
	suite.Error(err, "LoadConfig should return an error for invalid JSON")
}

// Test that we can obtain json data from the application table
func (suite *TestLoadDpiJson) TestLoadApplicationTable() {
	err := suite.manager.LoadConfigFromFile("testdata/DpiDefaultConfig.json")
	suite.NoError(err, "LoadConfigFromFile should not return an error")
	applications, err := suite.manager.GetTable("application")
	suite.NoError(err, "GetTable should not return an error")
	suite.NotEmpty(applications, "Application table should not be empty")

	categories, err := suite.manager.GetTable("category")
	suite.NoError(err, "GetTable should not return an error")
	suite.NotEmpty(categories, "Category table should not be empty")
}
