package dpi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"

	logger "github.com/untangle/golang-shared/services/logger"
)

const pluginName string = "dpi"

var QosmosFile = "/usr/share/veos/DpiDefaultConfig.json"

// MetaDataTable stores global metadata.
type MetaDataTable struct {
	Description      string   `json:"0description"`
	Version          string   `json:"0version"`
	VendorAttributes []string `json:"vendor-attributes"`
}

// QosmosInfo stores the complete application data.
type QosmosInfo struct {
	Name                    string                     `json:"name"` // Storing the application name here to access via ID.
	ID                      int                        `json:"id"`
	Description             string                     `json:"description"`
	Family                  string                     `json:"family"`
	Tag                     []string                   `json:"tag"`
	ServiceCategory         map[string]string          `json:"service-category"`
	VendorID                int                        `json:"vendor-id"`
	VendorServiceAttributes map[string]VendorAttribute `json:"vendor-service-attributes"`
}

// VendorAttribute holds vendor service attribute details.
type VendorAttribute struct {
	ID           int               `json:"id"`
	Type         string            `json:"type"`
	ValueService map[string]string `json:"value-service"`
}

// DpiConfig holds the entire DPI configuration.
type DpiConfig struct {
	MetaData     MetaDataTable       `json:"-"`
	Categories   map[string]int      `json:"categories"`
	Services     map[string]int      `json:"services"`
	Applications map[int]*QosmosInfo `json:"-"`
}

// rawConfig is used for JSON unmarshaling.
type rawConfig struct {
	Description      string                 `json:"0description"`
	Version          string                 `json:"0version"`
	Categories       map[string]int         `json:"categories"`
	Services         map[string]int         `json:"services"`
	VendorAttributes []string               `json:"vendor-attributes"`
	Applications     map[string]*QosmosInfo `json:"applications"`
}

// DpiConfigManager is the object that encapsulates the DPI information.
// Unexported config to prevent direct access to the configuration.
// Implements Plugin interface and PacketProcessorPlugin interface
type DpiConfigManager struct {
	config     DpiConfig
	fileSystem fs.FS
}

// returns DpiConfigManager instance
// provided as a constructor to the DI container
func NewDpiConfigManager(fs fs.FS) *DpiConfigManager {
	return &DpiConfigManager{
		config: DpiConfig{
			Categories:   make(map[string]int),
			Services:     make(map[string]int),
			Applications: make(map[int]*QosmosInfo),
		},
		fileSystem: fs,
	}
}

// LoadConfig reads the configuration JSON from the provided reader.
// reader is expected to be a JSON file.
func (m *DpiConfigManager) LoadConfig(r io.Reader) error {
	logger.Info("Loading DPI json configuration\n") // Read the JSON data from the reader
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read JSON data: %w", err)
	}
	logger.Debug("LoadConfig: Read %d bytes of JSON data\n", len(data))

	//raw contains the raw data from Dpi json
	var raw rawConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	logger.Info("LoadConfig: Successfully unmarshalled DPI JSON data\n")

	// Populate metadata object from Dpi json
	m.config.MetaData = MetaDataTable{
		Description:      raw.Description,
		Version:          raw.Version,
		VendorAttributes: raw.VendorAttributes,
	}

	// Populate categories and services from Dpi json
	m.config.Categories = raw.Categories
	m.config.Services = raw.Services

	logger.Debug("LoadConfig: Loaded metadata: %+v\n", m.config.MetaData)
	logger.Debug("LoadConfig: Loaded categories: %+v\n", m.config.Categories)
	logger.Debug("LoadConfig: Loaded services: %+v\n", m.config.Services)

	// Populate applications, mapping by ID to the application struct.
	for appName, app := range raw.Applications {
		app.Name = appName                  // Save the name within the struct.
		m.config.Applications[app.ID] = app // Map by ID.
	}
	logger.Info("Successfully loaded DPI configuration.\n")
	return nil
}

// LoadConfigFromFile loads the configuration from a file.
// 'filename' is the path to the JSON file containing the Dpi info.
func (m *DpiConfigManager) LoadConfigFromFile(filename string) error {
	file, err := m.fileSystem.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file %v: %w", filename, err)
	}
	defer file.Close()

	return m.LoadConfig(file)
}

// GetApplications retrieves all applications from the json.
func (m *DpiConfigManager) GetApplications() map[int]*QosmosInfo {
	return m.config.Applications
}

// GetMetaData returns a copy of the metadata. Used for testing
func (m *DpiConfigManager) GetMetaData() MetaDataTable {
	return m.config.MetaData
}

// GetCategories returns a copy of the categories map. Used for testing
func (m *DpiConfigManager) GetCategories() map[string]int {
	categoriesMap := make(map[string]int, len(m.config.Categories))
	for k, v := range m.config.Categories {
		categoriesMap[k] = v
	}
	return categoriesMap
}

// GetServices returns a copy of the services map. Used for testing
func (m *DpiConfigManager) GetServices() map[string]int {
	servicesMap := make(map[string]int, len(m.config.Services))
	for k, v := range m.config.Services {
		servicesMap[k] = v
	}
	return servicesMap
}

// Startup() is called once when the plugin is loaded.
func (m *DpiConfigManager) Startup() error {
	if err := m.LoadConfigFromFile(QosmosFile); err != nil {
		logger.Err("Failed to load DPI data from json file: %v\n", err)
		return err
	}
	return nil
}

// Name() returns the name of the plugin.
func (m *DpiConfigManager) Name() string {
	return pluginName
}

// Shutdown() is called once before plugin stops.
func (m *DpiConfigManager) Shutdown() error {
	logger.Info("Stopping the dpi config manager service\n")
	return nil
}

// GetTable returns the requested table as JSON
func (d *DpiConfigManager) GetTable(table string) (string, error) {
	logger.Debug("Getting %s table...\n", table)

	var data string
	var err error
	switch table {
	case "application":
		data, err = getApplicationTable(d.config.Applications)
	case "category":
		data, err = getCategoryTable(d.config.Applications)
	default:
		return data, errors.New("failed to get application classification table")
	}

	if err != nil {
		logger.Err("Unable to get DPI %s table: %s\n", table, err.Error())
		return "", err
	}

	return data, nil
}

// GetApplicationTable returns the application table as JSON
func getApplicationTable(qosinfo map[int]*QosmosInfo) (string, error) {
	// The format is a JSON array of information, we will fill in what we have.
	type result struct {
		Guid         string `json:"guid"`
		Index        int    `json:"index"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Category     string `json:"category"`
		Productivity int    `json:"productivity"`
		Risk         int    `json:"risk"`
		Flags        int    `json:"flags"`
		Reference    string `json:"reference"`
		Plugin       string `json:"plugin"`
	}
	logger.Debug("Getting application table...\n")
	// Populate the result array
	var results []result
	for _, app := range qosinfo {
		results = append(results, result{
			Guid:         fmt.Sprintf("%d", app.ID),
			Index:        app.ID,
			Name:         app.Name,
			Description:  app.Description,
			Category:     app.Family,
			Productivity: 0,
			Risk:         0,
			Flags:        0,
			Reference:    "",
			Plugin:       "",
		})
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Err("Unable to get DPI application table: %s\n", err.Error())
		return "", err
	}
	return string(jsonData), nil
}

// GetCategoryTable returns a distinct list of the family list we currently have in the ApplicationTable
// This is a better representation of categories, than any other field.
func getCategoryTable(qosinfo map[int]*QosmosInfo) (string, error) {
	// The format is a JSON array of "name": "string" pairs
	type result struct {
		Name string `json:"name"`
	}
	logger.Debug("Getting Category table...\n")

	// Create a map to hold the distinct family names
	familyMap := make(map[string]bool)
	for _, app := range qosinfo {
		familyMap[app.Family] = true
	}
	// Create result array
	var results []result
	for family := range familyMap {
		results = append(results, result{Name: family})
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Err("Unable to get DPI category table: %s\n", err.Error())
		return "", err
	}
	return string(jsonData), nil
}

// SetQosmosFile is used to modify path to DPI json file, used for testing purposes only.
func SetQosmosFile(path string) {
	QosmosFile = path
}
