package uritranslations

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

// Settings fields for uriTranslation
type URITranslation struct {
	URI	    string `json:"uri"`
	Scheme  string `json:"scheme"`
	Host    string `json:"host"`
	Path    string `json:"path"`
	Port    int    `json:"port"`
}

// Map of full uri key net/url/URL value
var uriMap map[string]*url.URL
// Map of host key net/url/URL value
var hostMap map[string]*url.URL
// Read/write lock for maps.
var mapMutex sync.RWMutex

// Startup is called when the service starts
func Startup() {
	buildMaps()
}

// Shutdown is called when the service stops
func Shutdown() {
}

// Build maps with uri an hosts as keys.
func buildMaps() {
	var entry *URITranslation
	var hostKey = ""

	path := []string{"uris", "uriTranslations"}

	// Read settings into list of uriTranslations variables.
	jsonResult, err := settings.GetSettings(path)
	if err != nil {
		logger.Info("Failed to read settings for path %v, %v\n", path, err.Error());
		return
	}
	if jsonResult == nil{
		logger.Info("Failed to read settings for path %v, %v\n", strings.Join(path,"/"), jsonResult);
		return
	}

	newURIMap := make(map[string]*url.URL)
	newHostMap := make(map[string]*url.URL)

	for _, b := range jsonResult.([]interface{}) {
		entry = new(URITranslation)

		jsonbody, err := json.Marshal(b)
		if err != nil {
			logger.Info("Error marshalling entry=%v with error %v\n", b, err.Error())
			continue
		}

		err = json.Unmarshal(jsonbody, &entry)
		if err != nil {
			logger.Info("Error unmarshalling entry=%v with error=%v\n", b, err.Error())
			continue
		}
		if entry.URI == ""{
			logger.Info("URI key field empty for entry=%v\n", b)
			continue
		}

		// Parse URITranslation fields into URL
		uri, err := url.Parse(entry.URI)
		hostKey = ""
		if err != nil{
			logger.Info("Unable to parse URI=%s with error=%v\n", entry.URI, err.Error())
			continue
		}
		if entry.Scheme != "" {
			uri.Scheme = entry.Scheme
		}
		hostKey = uri.Host
		if entry.Host != "" {
			uri.Host = entry.Host
		}
		if entry.Port > 0 {
			// Append post to host.
			uri.Host += ":" + strconv.Itoa(entry.Port)
		}
		if entry.Path != "" {
			uri.Path = entry.Path
		}

		// Add variables to new maps
		newURIMap[entry.URI] = uri
		newHostMap[hostKey] = uri
	}

	// Update active maps
	mapMutex.Lock()
	uriMap = newURIMap
	hostMap = newHostMap
	mapMutex.Unlock()
}

// getUriTranslation parses the URI (ignoring the query) and looks for the translated URL and returns the string of the found URL.
// If the path boolean is true, the lookup is performed on only the host instead of the full URL and
// the incoming path is used to build the returned URL.  This is useful for cases where the path is 
// possibly variable and dedicating a translation to the URL would require additional housekeeping
// than just changing a path.  For example, shop description links like:
// https://www.untangle.com/shop/virus-blocker
// https://www.untangle.com/shop/Live-Support
// ...
// Work best as with a single URI for https://www.untangle.com/shop
//
// If a match is not found, an error will be returned.
func getURITranslation(uri string, path bool) (string, error) {
	var err error = nil
	var ok bool
	var translatedUrl *url.URL

	if uriMap == nil {
		buildMaps()
	}

	parsedUrl, err := url.Parse(uri)
	if err != nil {
		// Unable to parse uri.
		logger.Info("Unable to parse uri=%s with error=%v\n", uri, err.Error())
		err = fmt.Errorf("Unable to parse uri=%v", uri)
	}else{
		// Get and clear query from parsed
		rawQuery := parsedUrl.RawQuery
		parsedUrl.RawQuery = ""
	
		mapMutex.RLock()
		if path == true{
			translatedUrl, ok = hostMap[parsedUrl.Host]
		}else{
			translatedUrl, ok = uriMap[parsedUrl.String()]
		}
		mapMutex.RUnlock()

		if ok != true {
			// Translation not found
			logger.Err("Unable to find url=%v\n", uri)
			err = fmt.Errorf("Unable to find url=%v", uri)
		}else{
			// Update the parsedUrl based URL with translated values.
			if translatedUrl.Scheme != "" {
				parsedUrl.Scheme = translatedUrl.Scheme
			}
			if translatedUrl.Host != "" {
				parsedUrl.Host = translatedUrl.Host
			}
			// Only add path if we're not explicitly overwrititng it.
			if path == false && translatedUrl.Path != "" {
				parsedUrl.Path = translatedUrl.Path
			}
			// Add query back.
			parsedUrl.RawQuery = rawQuery
			uri = parsedUrl.String()
		}
	}

	return uri, err
}

// GetURI looks up the specified URI (ignoring query) and returns the appropriate match.
// If returns an error if not found.
func GetURI(uri string) (string, error) {
	return getURITranslation(uri, false);
}

// GetURI looks up the host from the specified URI (ignoring query) and returns the appropriate match with the lookup URI's 
// path substituted for the translated value.
// It returns an error if not found.
func GetURIWithPath(uri string) (string, error) {
	return getURITranslation(uri, true);
}
