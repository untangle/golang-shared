package settings

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/untangle/golang-shared/services/logger"
)

// LicenseSub type, the main engine for licenses
type LicenseSub struct {
	enabledServices map[string]bool
	licenseFilename string
	uidFile         string
}

// License is the struct representing each individual license
type License struct {
	UID         string `json:UID`
	Type        string `json:type`
	End         int    `json:end`
	Start       int    `json:start`
	Seats       int    `json:seats default:-1`
	DisplayName string `json:"displayName"`
	Key         string `json:key`
	KeyVersion  int    `json:keyVersion`
	Name        string `json:name`
	javaClass   string `json:javaClass`
}

const (
	licenseReadRetries = 5 //Number of retries to try on reading the license file
)

// NewLicenseSub creates new license
func NewLicenseSub(licenseFilename string, uidFile string) *LicenseSub {
	logger.Info("Starting license sub...\n")

	l := new(LicenseSub)
	l.enabledServices = l.GetDefaults()
	l.licenseFilename = licenseFilename
	l.uidFile = uidFile

	return l
}

// GetDefaults gets the default enabled services file, where everything is disabled
// TODO: flag for product
func (l *LicenseSub) GetDefaults() map[string]bool {
	return map[string]bool{
		"loadBalancing":    false,
		"sslCertUpload":    false,
		"advancedLogging":  false,
		"manualRuleConfig": false,
		"ruleException":    false,
	}
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
		l.enabledServices = l.GetDefaults()
		return l.enabledServices, errors.New("Failed to read license file")
	}

	// unmarshal license inforamtion
	var licenses []License
	jsonErr := json.Unmarshal(fileBytes, &licenses)
	if jsonErr != nil {
		l.enabledServices = l.GetDefaults()
		return l.enabledServices, jsonErr
	}

	// determine if valid information
	for _, license := range licenses {
		isValid := l.isLicenseValid(license)
		if isValid {
			l.enabledServices[license.Name] = true
		}
	}

	return l.enabledServices, nil
}

// isLicenseValid checks if a given license is valid
func (l *LicenseSub) isLicenseValid(license License) bool {
	timeNow := time.Now()

	// Check if a valid type of service for the license
	_, ok := l.enabledServices[license.Name]
	if !ok {
		logger.Warn("Invalid license " + license.Name + ": Could not find in enabledServices map")
		return false
	}

	// start date in future
	if license.Start > int(timeNow.Unix()) {
		logger.Warn("Invalid license " + license.Name + ": Invalid start date (After today)\n")
		return false
	}

	// expired
	if license.End < int(timeNow.Unix()) {
		logger.Warn("Invalid license " + license.Name + ": Expired\n")
		return false
	}

	// check the uid
	uid, uidErr := GetUID(l.uidFile)
	if uidErr == nil {
		if len(license.UID) == 0 || license.UID != uid {
			logger.Warn("Invalid license " + license.Name + ": Invalid (UID mistmatch) \n")
			return false
		}
	} else {
		logger.Warn("Can't get uid check for license\n")
		return false
	}

	// get key
	var input string
	salt := "the meaning of life is 42"
	if license.KeyVersion == 1 || license.Seats == -1 {
		input = fmt.Sprintf("%d%s%s%s%d%d%s", license.KeyVersion, license.UID, license.Name, license.Type, license.Start, license.End, salt)
	} else if license.KeyVersion == 3 {
		input = fmt.Sprintf("%d%s%s%s%d%d%d%s", license.KeyVersion, license.UID, license.Name, license.Type, license.Start, license.End, license.Seats, salt)
	} else {
		logger.Warn("Invalid license " + license.Name + ": Invalid key version\n")
		return false
	}

	hasher := md5.New()
	hasher.Write([]byte(input))
	hex := hex.EncodeToString(hasher.Sum(nil))

	if hex != license.Key {
		logger.Warn("Invalid license " + license.Name + ": Invalid (key mistmatch) \n")
		return false
	}

	return true
}
