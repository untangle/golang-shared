package settings

import (
	"crypto/md5"
	"encoding/hex"
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
	Hash            string
	HashSalt        string
	LicenseInfo     LicenseInfo `json:"licenseInfo"`
}

// LicenseInfo represents the json returned from license server
type LicenseInfo struct {
	JavaClass string    `json:"javaClass"`
	List      []License `json:"list"`
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
	Valid       bool   `json:"valid" default:"false"`
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

// CheckHash checks if the license file has changed, returns true if it has
func CheckHash(filename string, currentHash string, hashSalt string) (bool, error) {
	hex, hexErr := CalculateHash(filename, hashSalt)
	if hexErr != nil {
		logger.Warn("Failed to calculate hash: %s\n", hexErr.Error())
		return true, hexErr
	}
	if hex != currentHash {
		logger.Warn("Hex does not match current hash\n")
		return true, nil
	}
	return false, nil
}

// CalculateHash calculates a hash given a file bytes
func CalculateHash(filename string, salt string) (string, error) {
	// read license file
	retries := licenseReadRetries
	var fileBytes []byte
	for retries > 0 {
		var readErr error
		fileBytes, readErr = ioutil.ReadFile(filename)
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
		return "", errors.New("Failed to read hashable file")
	}
	// create new hash
	hasher := md5.New()
	hasher.Write(fileBytes)
	hasher.Write([]byte(salt))
	hex := hex.EncodeToString(hasher.Sum(nil))

	return hex, nil
}

// NewLicenseSub creates new license
func NewLicenseSub(licenseFilename string, uidFile string, product string, hashSalt string) *LicenseSub {
	logger.Info("Starting license sub...\n")

	l := new(LicenseSub)
	l.enabledServices = GetLicenseDefaults(product)
	l.licenseFilename = licenseFilename
	l.uidFile = uidFile
	l.product = product
	l.HashSalt = hashSalt

	return l
}

// CleanUp cleans up the contexts of the licenseSub
func (l *LicenseSub) CleanUp() {
	logger.Info("Shutting down license sub...\n")
}

// GetEnabledServices gets the enabled services
func (l *LicenseSub) GetEnabledServices() (map[string]bool, error) {
	l.enabledServices = GetLicenseDefaults(l.product)
	// get hash of new license file
	hash, hashErr := CalculateHash(l.licenseFilename, l.HashSalt)
	if hashErr != nil {
		logger.Warn("Failure generating hash: %s\n", hashErr.Error())
		return nil, hashErr
	}

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
		return l.enabledServices, errors.New("Failed to read license file")
	}

	// unmarshal license inforamtion
	var licenseInfo LicenseInfo
	jsonErr := json.Unmarshal(fileBytes, &licenseInfo)
	if jsonErr != nil {
		return l.enabledServices, jsonErr
	}
	l.LicenseInfo = licenseInfo
	l.determineEnabledServices()

	//set hash
	l.Hash = hash

	return l.enabledServices, nil
}

// determineEnabledServices determines which service are enabled for the license handler
func (l *LicenseSub) determineEnabledServices() {
	for _, license := range l.LicenseInfo.List {
		_, ok := l.enabledServices[license.Name]
		if ok {
			l.enabledServices[license.Name] = license.Valid
		} else {
			logger.Warn("Saw a license name that's unknown: %s\n", license.Name)
		}
	}
}
