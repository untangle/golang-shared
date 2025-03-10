package dpiclassmanager

import (
	"encoding/json"
	"errors"
	"os"

	logService "github.com/untangle/golang-shared/services/logger"
)

var logger = logService.GetLoggerInstance()

// Application Info, not every field from config file is represented here
type ApplicationInfo struct {
	family    string
	id        int
	tag       []string
	attribute interface{}
}

type ApplicationTable map[string]*ApplicationInfo

// Provide an interface to load and query Qosmos application classification data
type DpiClassManager interface {
	// Load the application classification data
	LoadApplicationTable() error
	// Get a specific table
	GetTable(table string) (string, error)
}

type DpiClassManagerImpl struct {
	DpiConfigFile    string
	ApplicationTable ApplicationTable
	logger           *logService.Logger
}

const DpiConfigFile = "/usr/shared/veos/DpiDefaultConfig.json"

func Startup() {
	logger.Info("Starting up the DPI class manager service\n")
	// Create the DPI class manager
	dpi := &DpiClassManagerImpl{}
	dpi.DpiConfigFile = DpiConfigFile
	dpi.ApplicationTable = make(ApplicationTable)
	dpi.logger = logService.GetLoggerInstance()
	// Load the application table
	err := dpi.LoadApplicationTable()
	if err != nil {
		logger.Err("Failed to load DPI application table: %s\n", err.Error())
	}
}

func Shutdown() {
	logger.Info("Shutting down the DPI class manager service\n")
}

func (d *DpiClassManagerImpl) LoadApplicationTable() error {
	logger.Debug("Loading application table...\n")

	// Open the file
	file, err := os.Open(d.DpiConfigFile)
	if err != nil {
		logger.Err("Error opening file: %s\n", err)
		return err
	}
	defer file.Close()

	// Read file, its a json file.
	decoder := json.NewDecoder(file)
	// Just create generic map for now and we will convert to ApplicationTable
	var appConfig map[string]interface{}
	err = decoder.Decode(&appConfig)
	if err != nil {
		logger.Err("Error decoding json: %s\n", err)
		return err
	}

	// Convert to ApplicationTable, look for the key "applications"
	apps := appConfig["applications"]
	if apps == nil {
		logger.Err("No applications found in json\n")
		return errors.New("no applications found in json")
	}

	// Assert apps to the correct type
	appsMap, ok := apps.(map[string]interface{})
	if !ok {
		logger.Err("invalid type for applications\n")
		return errors.New("invalid type for applications")
	}

	// Iterate over the applications
	for key, value := range appsMap {
		// Convert value to map
		valueMap, ok := value.(map[string]interface{})
		if !ok {
			logger.Err("invalid type for application value\n")
			return errors.New("invalid type for application value")
		}

		// Convert to ApplicationInfo
		appInfo := &ApplicationInfo{}
		appInfo.family, _ = valueMap["family"].(string)
		appInfo.id, _ = valueMap["id"].(int)
		appInfo.tag, _ = valueMap["tag"].([]string)
		appInfo.id, _ = valueMap["vendor-id"].(int)
		appInfo.attribute = valueMap["vendor-service-attribute"]
		// Add to ApplicationTable
		d.ApplicationTable[key] = appInfo
	}

	return nil
}

func (d *DpiClassManagerImpl) GetTable(table string) (string, error) {
	logger.Debug("Getting %s table...\n", table)

	var data string
	var err error
	switch table {
	case "application":
		data, err = getApplicationTable(d.ApplicationTable)
	case "category":
		data, err = getCategoryTable(d.ApplicationTable)
	default:
		return data, errors.New("failed_to_get_table")
	}

	if err != nil {
		logger.Err("Unable to get DPI %s table: %s\n", table, err.Error())
		return "", err
	}

	return data, nil
}

// GetApplicationTable returns the application table as JSON
func getApplicationTable(at ApplicationTable) (string, error) {
	logger.Debug("Getting application table...\n")

	// convert it to a slice first
	appTable := []*ApplicationInfo{}

	for _, val := range at {
		appTable = append(appTable, val)
	}

	jsonData, err := json.Marshal(appTable)

	if err != nil {
		logger.Err("Unable to get DPI application table: %s\n", err.Error())
		return "", err
	}

	return string(jsonData), nil
}

// GetCategoryTable returns a distinct list of the family list we currently have in the ApplicationTable
// This is a better representation of categories, than any other field.
func getCategoryTable(at ApplicationTable) (string, error) {
	logger.Debug("Getting Category table...\n")
	// List of strings
	cats := make(map[string]bool)
	for _, val := range at {
		cats[val.family] = true
	}
	jsonData, err := json.Marshal(cats)
	if err != nil {
		logger.Err("Unable to get DPI category table: %s\n", err.Error())
		return "", err
	}
	return string(jsonData), nil
}
