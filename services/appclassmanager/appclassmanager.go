package appclassmanager

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"

	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/util/environments"
)

var logger = logService.GetLoggerInstance()

// ApplicationInfo stores the details for each know application
// this structure matches the csv fields
type ApplicationInfo struct {
	GUID         string `json:"guid"`
	Index        int    `json:"index"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	Productivity uint   `json:"productivity"`
	Risk         uint   `json:"risk"`
	Flags        uint64 `json:"flags"`
	Reference    string `json:"reference"`
	Plugin       string `json:"plugin"`
}

// CategoryInfo contains details about a category (used when converting ApplicationTable to categories for extjs store)
type CategoryInfo struct {
	Name string `json:"name"`
}

// MetaDataTable stores global metadata
type MetaDataTable struct {
	Description      string   `json:"0description"`
	Version          string   `json:"0version"`
	VendorAttributes []string `json:"vendor-attributes"`
}

// QosmosCategoryInfo stores categories mapped to integers
type QosmosCategoryInfo struct {
	Categories map[string]int `json:"categories"`
}

// QosmosServiceInfo stores services mapped to integers
type QosmosServiceInfo struct {
	Services map[string]int `json:"services"`
}

// QosmosInfo stores the complete application data
type QosmosInfo struct {
	Name                    string                     `json:"name"` // Storing the application name here to access via ID
	ID                      int                        `json:"id"`
	ServiceCategory         map[string]string          `json:"service-category"`
	VendorID                int                        `json:"vendor-id"`
	VendorServiceAttributes map[string]VendorAttribute `json:"vendor-service-attributes"`
}

// Vendor Service Attributes
type VendorAttribute struct {
	ID           int               `json:"id"`
	Type         string            `json:"type"`
	ValueService map[string]string `json:"value-service"`
}

const guidInfoFile = "/usr/share/untangle-classd/protolist.csv"
const QosmosFile = "/usr/share/veos/DpiDefaultConfig.json"

// MetaDataTable instance for storing metadata
var MetaData MetaDataTable

// QosmosCategoryInfo instance for storing categories
var QosmosCategory QosmosCategoryInfo

// QosmosServiceInfo instance for storing services
var QosmosService QosmosServiceInfo

// QosmosTable maps application IDs to QosmosInfo structs
var QosmosTable map[int]*QosmosInfo

// ApplicationTable stores the details for each known application
var ApplicationTable map[string]*ApplicationInfo

// Startup is called when the packetd service starts
func Startup() {
	logger.Info("Starting up the Application Classification Table manager service\n")
	if !environments.IsEOS() {
		logger.Info("Loading Procera Application Table")
		loadApplicationTable()
	} else {
		logger.Info("Loading Qosmos Table")
		loadQosmosTable()
	}
}

// Shutdown is called when the packetd service stops
func Shutdown() {
	logger.Info("Shutting down the Application Classification Table manager service\n")
}

// GetTable gets the classd table specified by the table param
func GetTable(table string) (string, error) {
	logger.Debug("Getting %s table...\n", table)

	var data string
	var err error
	switch table {
	case "application":
		data, err = getApplicationTable()
	case "category":
		data, err = getCategoryTable()
	default:
		return data, errors.New("failed_to_get_table")
	}

	if err != nil {
		logger.Err("Unable to get ClassD %s table: %s\n", table, err.Error())
		return "", err
	}

	return data, nil
}

// GetApplicationTable returns the application table as JSON
func getApplicationTable() (string, error) {
	logger.Debug("Getting application table...\n")

	// convert it to a slice first
	appTable := []*ApplicationInfo{}

	for _, val := range ApplicationTable {
		appTable = append(appTable, val)
	}

	jsonData, err := json.Marshal(appTable)

	if err != nil {
		logger.Err("Unable to get ClassD application table: %s\n", err.Error())
		return "", err
	}

	return string(jsonData), nil
}

// GetCategoryTable returns a distinct list of the categories we currently have in the ApplicationTable
func getCategoryTable() (string, error) {
	logger.Debug("Getting Category table...\n")

	// Instead of two loops, create a map that indicates if the items exist in the slice
	catMap := make(map[string]bool)
	catSlice := []*CategoryInfo{}

	// Iterate the table, if the map contains the slice then continue, otherwise add it to the map
	for _, val := range ApplicationTable {
		if catMap[val.Category] {
			continue
		}

		catSlice = append(catSlice, &CategoryInfo{Name: val.Category})
		catMap[val.Category] = true
	}

	// Now convert the slice of CategoryInfo into JSON
	jsonData, err := json.Marshal(catSlice)

	if err != nil {
		logger.Err("Unable to get ClassD Category table: %s\n", err.Error())
		return "", err
	}

	return string(jsonData), nil

}

// loadApplicationTable loads the details for each application
func loadApplicationTable() {
	var file *os.File
	var linecount int
	var infocount int
	var list []string
	var err error

	ApplicationTable = make(map[string]*ApplicationInfo)
	filename, err := settings.LocateFile(guidInfoFile)
	if err != nil {
		logger.Warn("Unable to locate GUID info file: %s\n",
			guidInfoFile)
		return
	}
	// open the guid info file provided by Sandvine
	file, err = os.Open(filename)

	// if there was an error log and return
	if err != nil {
		logger.Warn("Unable to load application details: %s\n", guidInfoFile)
		return
	}

	// create a new CSV reader
	reader := csv.NewReader(bufio.NewReader(file))
	for {
		list, err = reader.Read()

		if err == io.EOF {
			// on end of file just break out of the read loop
			break
		} else if err != nil {
			// for anything else log the error and break
			logger.Err("Unable to parse application details: %v\n", err)
			break
		}

		// count the number of lines read so we can compare with
		// the number successfully parsed when we finish loading
		linecount++

		// skip the first line that holds the file format description
		if linecount == 1 {
			continue
		}

		// if we did not parse exactly 10 fields skip the line
		if len(list) != 10 {
			logger.Warn("Invalid line length: %d\n", len(list))
			continue
		}

		// create a object to store the details
		info := new(ApplicationInfo)

		info.GUID = list[0]
		info.Index, err = strconv.Atoi(list[1])
		if err != nil {
			logger.Warn("Invalid index: %s\n", list[1])
		}
		info.Name = list[2]
		info.Description = list[3]
		info.Category = list[4]
		tempProd, err := strconv.ParseUint(list[5], 10, 8)
		if err != nil {
			logger.Warn("Invalid productivity: %s\n", list[5])
		}
		info.Productivity = uint(tempProd)
		tempRisk, err := strconv.ParseUint(list[6], 10, 8)
		if err != nil {
			logger.Warn("Invalid risk: %s\n", list[6])
		}
		info.Risk = uint(tempRisk)
		info.Flags, err = strconv.ParseUint(list[7], 10, 64)
		if err != nil {
			logger.Warn("Invalid flags: %s %s\n", list[7], err)
		}
		info.Reference = list[8]
		info.Plugin = list[9]

		// store the object in the table using the guid as the index
		ApplicationTable[info.GUID] = info
		infocount++
	}

	file.Close()
	logger.Info("Loaded classification details for %d applications\n", infocount)

	// if there were any bad lines in the file log a warning
	if infocount != linecount-1 {
		logger.Warn("Detected garbage in the application info file: %s\n", guidInfoFile)
	}
}

func loadQosmosTable() {

	var err error

	QosmosTable = make(map[int]*QosmosInfo)

	filename, err := settings.LocateFile(QosmosFile)
	if err != nil {
		logger.Warn("Unable to locate Qosmos info file: %s\n", QosmosFile)
		return
	}

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		logger.Warn("Unable to load Qosmos application details: %s\n", QosmosFile)
		return
	}

	// Temporary struct for unmarshaling the entire JSON file
	// this struct will be used to populate the metadatatable, and allow for mapping from appid to appname
	var rawData struct {
		Description      string                `json:"0description"`
		Version          string                `json:"0version"`
		Categories       map[string]int        `json:"categories"`
		Services         map[string]int        `json:"services"`
		VendorAttributes []string              `json:"vendor-attributes"`
		Applications     map[string]QosmosInfo `json:"applications"`
	}

	// Unmarshal the JSON into the temporary struct
	err = json.Unmarshal(jsonData, &rawData)
	if err != nil {
		logger.Err("Failed to unmarshal Qosmos JSON file: %s\n", err)
		return
	}

	if len(rawData.Applications) == 0 {
		logger.Warn("No applications found in Qosmos JSON file")
		return
	}

	logger.Info("Qosmos Json data successfuly extracted")

	// Populate MetaDataTable with the parsed metadata
	MetaData = MetaDataTable{
		Description:      rawData.Description,
		Version:          rawData.Version,
		VendorAttributes: rawData.VendorAttributes,
	}

	QosmosCategory.Categories = rawData.Categories
	QosmosService.Services = rawData.Services

	// Populate QosmosTable with applications
	for appName, app := range rawData.Applications {
		// Add application name to the struct
		app.Name = appName
		// Store the application in QosmosTable keyed by its ID
		QosmosTable[app.ID] = &app
	}
	logger.Info("Qosmos name for ID-949 : %s", QosmosTable[949].Name)
	logger.Info("Successfully loaded Qosmos metadata and application details")
}
