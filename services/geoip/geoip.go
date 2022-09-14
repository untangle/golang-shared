package geoip

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/services/uritranslations"
	"github.com/untangle/golang-shared/util/cache/cacher"
)

// the name of the database file we look for in tarballs we download.
const MaxMindDbFileName = "GeoLite2-Country.mmdb"

// how long a downloaded database is valid for before we should
// download a new one. Thirty days. We figure out how old it is by
// looking at the file timestamp.
const validityOfDbDuration = time.Hour * 24 * 30

// The name of the cache used by geoip, only used for debugging/logging purposes
const cacheName = "geoIpCache"

// Capacity of cache used by GeoIPManager for country code lookups.
const cacheCapacity = 500

// GeoIPDB is an interface that GeoIP databases conform to.
type GeoIPDB interface {
	// LookupCountryCodeOfIP will look up the country code of a
	// given IP address. If it is found, it returns the code and
	// true. If it is not found, it returns "", and false.
	LookupCountryCodeOfIP(ip net.IP) (string, bool)

	// Refresh() will 'refresh' the database -- downloading it
	// from a remote source if necessary.
	Refresh() error
}

// GeoIPClassifier encapsulates all of the logic and responsibilities of
// the GEO IP plugin. It handles:
// 1. Downloading a geoIP database periodically.
//
// 2. Identifying the country code of a dispatch.Session object, when
// called from the main packetd dispatching code as a plugin.
//
// 3. Notifying all Listeners when it determines (2). See the
// Observable object.

// MaxMindGeoIPManager is a GeoIPDB implementation of the geo IP
// database that uses a MaxMind database.
type MaxMindGeoIPManager struct {
	databaseFilename string

	geoDatabaseReader *geoip2.Reader
	databaseCache     cacher.Cacher
	cacheLocker       sync.RWMutex
}

// NewMaxMindGeoIPManager creates a new NewMaxMindGeoIPManager that
// 'points at' the database given by filename. filename need not
// exist, if it does not, you will need to downloadAndExtractDb().
func NewMaxMindGeoIPManager(filename string) *MaxMindGeoIPManager {
	return &MaxMindGeoIPManager{
		databaseCache:    cacher.NewLruCache(cacheCapacity, cacheName),
		databaseFilename: filename}
}

// downloadAndExtractDB will download the MaxMind geoIP database from
// downloads.untangle.com and extract the database itself into the
// filesystem it by calling extractDBFile(). It returns an error if
// not successful. If successful, you will probably want to call
// openDBFile().
func (db *MaxMindGeoIPManager) downloadAndExtractDB() error {
	uid, err := settings.GetUIDOpenwrt()
	if err != nil {
		uid = "00000000-0000-0000-0000-000000000000"
		logger.Warn("Unable to read UID: %s - Using all zeros\n", err.Error())
	}
	target := fmt.Sprintf(
		"https://downloads.untangle.com/download.php?resource=geoipCountry&uid=%s",
		uid)
	translatedTarget, err := uritranslations.GetURI(target)
	if err == nil {
		target = translatedTarget
	}
	resp, err := http.Get(target)
	if err != nil {
		return fmt.Errorf("HTTP GET failure: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP bad return code: %v, %v",
			resp.StatusCode,
			resp.Status)
	}
	return db.extractDBFile(resp.Body)

}
func (db *MaxMindGeoIPManager) extractDBFile(reader io.ReadCloser) error {
	defer reader.Close()

	logger.Info("Starting GeoIP database extraction: %s\n", db.databaseFilename)

	// Make sure the target directory exists
	marker := strings.LastIndex(db.databaseFilename, "/")

	// Get the index of the last slash so we can isolate the path and create the directory
	if marker > 0 {
		os.MkdirAll(db.databaseFilename[0:marker], 0755)
	}

	// Create a reader for the compressed data
	zipReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("error calling gzip.NewReader(): %w", err)
	}
	defer zipReader.Close()

	// Create a tar reader using the uncompressed data stream
	tarReader := tar.NewReader(zipReader)

	// Create the file where we'll store the extracted database
	writer, err := os.Create(db.databaseFilename)

	if err != nil {
		return fmt.Errorf(
			"unable to write database file: %s",
			db.databaseFilename)
	}

	readError := fmt.Errorf(
		"couldn't extract expected db file: %s from tar",
		MaxMindDbFileName)

	for {
		// get the next entry in the archive
		header, err := tarReader.Next()

		// break out of the loop on end of file
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			// log any other errors and break out of the loop
			readError = fmt.Errorf(
				"error while reading database tar archive: %w",
				err)
			break
		}

		// ignore everything that is not a regular file
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// ignore everything except the actual database file
		if !strings.HasSuffix(header.Name, MaxMindDbFileName) {
			continue
		}

		// found the database so write to the output file, set the goodfile flag, and break
		if _, err := io.Copy(writer, tarReader); err != nil {
			readError = fmt.Errorf(
				"Error writing found DB file to disk: %w",
				err)
		} else {
			readError = nil
			logger.Info("Finished GeoIP database download\n")

		}
		break
	}
	writer.Close()
	// If we had an error, delete the created file.
	if readError != nil {
		os.Remove(db.databaseFilename)
		return readError
	}
	return nil
}

// checkForDBFile checks if the MaxMind geoIP database file exists and
// is current. 'current' is defined by the value of
// validityOfDbDuration.
func (db *MaxMindGeoIPManager) checkForDBFile() bool {
	filename := db.databaseFilename
	if fileinfo, err := os.Stat(filename); err != nil {
		return false
	} else if fileinfo.Size() == 0 {
		return false
	} else {

		filetime := fileinfo.ModTime()
		currtime := time.Now()

		// Return true if the time since the file modification time to
		// now is less than the validity period (i.e. return true if
		// the file is not stale yet).
		return currtime.Sub(filetime) < validityOfDbDuration
	}
}

// openDBFile will open the database file, calling the underlying
// MaxMind implementation package. If the database is already open, it
// closes it and re-opens it.
func (db *MaxMindGeoIPManager) openDBFile() error {
	if db.geoDatabaseReader != nil {
		db.geoDatabaseReader.Close()
		db.geoDatabaseReader = nil
	}

	mmDB, err := geoip2.Open(db.databaseFilename)

	if err != nil {
		logger.Warn("Unable to load GeoIP Database: %s\n", err)
		db.geoDatabaseReader = nil
		return fmt.Errorf("couldn't open GeoIP db: %w", err)
	}
	logger.Info("Loading GeoIP Database: %s\n", db.databaseFilename)
	db.geoDatabaseReader = mmDB
	return nil
}

// LookupCountryCodeOfIP looks up the country code of ip. A cache of
// previously looked up countries is checked first. If not in the cache,
// the code is looked up in the database and added to the cache. If the
// country code is found in the cache/database, it
// returns the code and the value true. If not found, it returns ""
// and false.
func (db *MaxMindGeoIPManager) LookupCountryCodeOfIP(ip net.IP) (string, bool) {
	db.cacheLocker.Lock()
	defer db.cacheLocker.Unlock()
	retCountryCode := ""
	retOk := false

	if db.geoDatabaseReader != nil {
		if countryFromCache, ok := db.databaseCache.Get(ip.String()); ok {
			retCountryCode = countryFromCache.(*geoip2.Country).Country.IsoCode
			retOk = true

			// Lookup country code in the database if a cache miss occurs. Update cache
			// with the retrieved value
		} else if countryFromDb, err := db.geoDatabaseReader.Country(ip); err == nil {
			if len(countryFromDb.Country.IsoCode) != 0 {
				db.databaseCache.Put(ip.String(), countryFromDb)

				retCountryCode = countryFromDb.Country.IsoCode
				retOk = true
			}
		}
	} else {
		logger.Warn(
			"LookupCountryCodeOfIP() called with nil MaxMind DB reader!\n")
		retOk = false
	}

	return retCountryCode, retOk
}

// Refresh will:
//
// 1. Check to see if the database file exists and is current. If so,
// it opens it if it hasn't been opened, and exits.
//
// 2. If the database file is not current or does not exist as
// determined by checkForDBFile, it will download it and open it.
//
// 3. Clear the cache storing previously looked up country codes
//
// So unless something goes wrong (i.e. a database file doesn't exist
// on the filesystem and we can't download a new one), at the end of
// this call the database should always be opened in the
// MaxMindGeoIPManager object.
func (db *MaxMindGeoIPManager) Refresh() error {
	if !db.checkForDBFile() {
		err := db.downloadAndExtractDB()
		if err != nil {
			return err
		}
		err = db.openDBFile()
		if err != nil {
			return err
		}
	} else if db.geoDatabaseReader == nil {
		if err := db.openDBFile(); err != nil {
			return err
		}
	}

	db.cacheLocker.Lock()
	db.databaseCache.Clear()
	db.cacheLocker.Unlock()

	return nil
}

// It holds GeoIPDB
var GeoIPManager *LockingGeoIPManager

// LockingGeoIPManager is a wrapper object for a GeoIPManager
// (specifically MaxMindGeoIPManager) that wraps all calls to
// LookupCountryCodeOfIP and Refresh with an RWLock. Refresh() will be
// called with a write lock and LookupCountryCodeOfIP has a read lock.
type LockingGeoIPManager struct {
	lock sync.RWMutex
	ipDB GeoIPDB
}

// NewLockingGeoIPManager creates a new LockingGeoIPManager, which
// wraps the db object given.
func NewLockingGeoIPManager(db GeoIPDB) *LockingGeoIPManager {
	return &LockingGeoIPManager{
		ipDB: db,
	}
}

// LookupCountryCodeOfIP calls the underlying ipDB of the
// LockingGeoIPManager's LookupCountryCodeOfIP after taking out an
// RLock(). It returns whatever the underlying database returns.
func (db *LockingGeoIPManager) LookupCountryCodeOfIP(ip net.IP) (string, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.ipDB.LookupCountryCodeOfIP(ip)
}

// Refresh calls the underling ipDB's Refresh method after calling
// Lock() (i.e. taking out a write lock). It returns whatever the
// Refresh method of the underlying object returns.
func (db *LockingGeoIPManager) Refresh() error {
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.ipDB.Refresh()
}

