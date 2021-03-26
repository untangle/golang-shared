package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/untangle/golang-shared/services/logger"
)

// LicenseSub type, the main engine for licenses
type LicenseSub struct {
	enabledServices map[string]bool
	licenseFilename string
	uidFile         string
	product         string
}

// License is the struct representing each individual license
type License struct {
	UID         string `json:"UID"`
	Type        string `json:"type"`
	End         int    `json:"end"`
	Start       int    `json:"start"`
	Seats       int    `json:"seats" default:"-1"`
	DisplayName string `json:"displayName"`
	Key         string `json:"key"`
	KeyVersion  int    `json:"keyVersion"`
	Name        string `json:"name"`
	JavaClass   string `json:"javaClass"`
}

const (
	licenseReadRetries = 5 //Number of retries to try on reading the license file
)

// GetLicenseDefaults gets the default enabled services file, where everything is disabled
func GetLicenseDefaults(product string) map[string]bool {
	var defaults map[string]bool
	switch product {
	case "WAF":
		defaults = map[string]bool{
			"loadBalancing":    false,
			"sslCertUpload":    false,
			"advancedLogging":  false,
			"manualRuleConfig": false,
			"ruleException":    false,
		}
	}

	return defaults
}

// NewLicenseSub creates new license
func NewLicenseSub(licenseFilename string, uidFile string, product string) *LicenseSub {
	logger.Info("Starting license sub...\n")

	l := new(LicenseSub)
	l.enabledServices = GetLicenseDefaults(product)
	l.licenseFilename = licenseFilename
	l.uidFile = uidFile
	l.product = product

	return l
}

// CleanUp cleans up the contexts of the licenseSub
func (l *LicenseSub) CleanUp() {
	logger.Info("Shutting down license sub...\n")
}

// GetLicenses gets the enabled services
func (l *LicenseSub) GetLicenses() (map[string]bool, error) {
	// read license file
	retries := licenseReadRetries
	var fileBytes []byte
	for retries > 0 {
		var readErr error
		fileBytes, readErr = ioutil.ReadFile(l.licenseFilename)
		if readErr != nil {
			retries = retries - 1
			// sleep one second, perhaps file is being written
			time.Sleep(time.Second)
		} else {
			retries = 1
			break
		}
	}

	if retries <= 0 {
		l.enabledServices = GetLicenseDefaults(l.product)
		return l.enabledServices, errors.New("Failed to read license file")
	}

	// unmarshal license inforamtion
	var licenses map[string]bool
	jsonErr := json.Unmarshal(fileBytes, &licenses)
	if jsonErr != nil {
		l.enabledServices = GetLicenseDefaults(l.product)
		return l.enabledServices, jsonErr
	}

	l.enabledServices = licenses
	return l.enabledServices, nil
}
