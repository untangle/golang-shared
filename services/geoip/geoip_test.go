package geoip

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/untangle/golang-shared/testing/data"
)

// TestGeoIP is a Test suite for testing the geoip plugin. We use the
// testify suite package for this.
type TestGeoIP struct {
	suite.Suite

	// geoip plugin under test.
	//GeoIPClassifier *GeoIPClassifier

	// list of files to delete after test cases.
	deleteFiles []string
}

// get a filename to extract the database to. This filename will be in
// a temporary directory created for the test case. It will contain a
// non-existent directory inside that directory.
func (suite *TestGeoIP) getDBFilename() string {
	tmpDir, err := os.MkdirTemp("", "GeoIPUnitTest")
	suite.failIferror(err, "Can't open tmpDir")
	fullFileName := path.Join(tmpDir, "extraComponent", MaxMindDbFileName)
	suite.addDeleteFile(tmpDir)
	return fullFileName
}

// addDeleteFile adds file to a list of files that will be deleted (as
// if with rm -r) after the end of each test case.
func (suite *TestGeoIP) addDeleteFile(file string) {
	suite.deleteFiles = append(suite.deleteFiles, file)
}

// generic function for testing that we fail gracefully for various
// types of bad/invalid database files.
//
// filename -- filename to try to extract with extractDBFile
//
// description -- string description of the expected error so we can
// log it for visibility.
//
// shouldExtractFail -- if set to true we assert extractDBFile returns
// an error, else we don't.
func (suite *TestGeoIP) testFailure(filename string, description string, shouldExtractFail bool) {
	extractedDbFileName := suite.getDBFilename()
	readerCloser := suite.getReaderCloserForFile(filename)
	geoIP := NewMaxMindGeoIPManager(extractedDbFileName)
	result := geoIP.extractDBFile(readerCloser)
	if shouldExtractFail {
		suite.NotNil(result)
	} else {
		suite.Nil(result)
	}

	result = geoIP.openDBFile()
	suite.NotNil(result)
	log.Printf("Caught error from %s (expected): %v",
		description,
		result)
	countryResult, found := geoIP.LookupCountryCodeOfIP(net.ParseIP("2000:ff::1"))
	suite.Equal(countryResult, "")
	suite.False(found)
}

// Check that Refresh() works with a pre-existing file.
func (suite *TestGeoIP) TestRefreshWithGoodDbFile() {
	fullFileName := suite.getDBFilename()
	geoIP := NewMaxMindGeoIPManager(fullFileName)
	result, found := geoIP.LookupCountryCodeOfIP(net.ParseIP("3.3.3.3"))
	suite.False(found)
	suite.Equal(result, "")
	is_file_valid, is_file_stale := geoIP.checkForDBFile()
	suite.False(is_file_valid)
	suite.False(is_file_stale)
	suite.Nil(geoIP.extractDBFile(suite.getReaderCloserForTestTarball()))
	is_file_valid, is_file_stale = geoIP.checkForDBFile()
	suite.True(is_file_valid)
	suite.False(is_file_stale)
	suite.Nil(geoIP.geoDatabaseReader)
	suite.Nil(geoIP.Refresh())
	is_file_valid, is_file_stale = geoIP.checkForDBFile()
	suite.True(is_file_valid)
	suite.False(is_file_stale)
	suite.NotNil(geoIP.geoDatabaseReader)
	result, found = geoIP.LookupCountryCodeOfIP(net.ParseIP("3.3.3.3"))
	suite.True(found)
	suite.NotEqual(result, "")
}
func (suite *TestGeoIP) TestDBExtract() {
	fullFileName := suite.getDBFilename()
	geoIP := NewMaxMindGeoIPManager(fullFileName)

	result := geoIP.extractDBFile(suite.getReaderCloserForTestTarball())
	suite.Nil(result)
	if _, err := os.Stat(fullFileName); err != nil {
		suite.Failf("File wasn't created", fullFileName)
	}

	// Manually computed from terminal.
	expectedFileSha256Hex :=
		"85f9fef478a5366daac20c71b5d6784b90bce61fafc90502ff974f083c09563c"
	openedFile, err := os.Open(fullFileName)
	suite.failIferror(err, "Can't open produced file")
	bytes, err := io.ReadAll(openedFile)
	suite.failIferror(err, "Couldn't read produced DB file")

	hash := sha256.Sum256(bytes)
	suite.Equal(
		expectedFileSha256Hex,
		hex.EncodeToString(hash[:]))
}

// Check that the checkForDBFile method works:
//
//  1. It returns false for both is_file_valid and is_file_stale
//     if the file doesn't exist.
//
//  2. It returns false for is_file_valid and true for is_file_stale
//     if the file exists and too old when the download fails.
//
//  3. It returns true for is_file_valid and false for is_file_stale
//     if the file exists and not too old.
func (suite *TestGeoIP) TestDBStatusChecker() {
	fullFileName := suite.getDBFilename()
	geoIP := NewMaxMindGeoIPManager(fullFileName)
	is_file_valid, is_file_stale := geoIP.checkForDBFile()
	suite.False(is_file_valid)
	suite.False(is_file_stale)
	extractResult := geoIP.extractDBFile(
		suite.getReaderCloserForTestTarball())
	suite.Nil(extractResult)
	is_file_valid, is_file_stale = geoIP.checkForDBFile()
	suite.True(is_file_valid)
	suite.False(is_file_stale)

	// Here we use os.Chtimes to backdate the file's timestamps
	// and make it appear old. It is one second older than
	// validityOfDbDuration, so we expect checkForDBFile to return
	// false.
	now := time.Now()
	invalidTime := now.Add(-(validityOfDbDuration + time.Second))
	err := os.Chtimes(
		fullFileName,
		invalidTime,
		invalidTime)
	suite.failIferror(err, "Couldn't backdate file timestamp")
	is_file_valid, is_file_stale = geoIP.checkForDBFile()
	suite.False(is_file_valid)
	suite.True(is_file_stale)
	suite.NotNil(geoIP.geoDatabaseReader)
	result, found := geoIP.LookupCountryCodeOfIP(net.ParseIP("3.3.3.3"))
	suite.True(found)
	suite.NotEqual(result, "")
}

// Test that we call the MaxMind database reader correctly and do not
// get errors. Currently we do not make assertions about returned
// country codes, since we do not 'own' the databas, but merely query
// it. This is to make sure that we don't get unexpected errors, and
// that the basic sequence of functions for extracting and opening the
// file work as expected.
func (suite *TestGeoIP) TestDBReader() {
	fullFileName := suite.getDBFilename()
	geoIP := NewMaxMindGeoIPManager(fullFileName)
	extractResult := geoIP.extractDBFile(
		suite.getReaderCloserForTestTarball())
	suite.Nil(extractResult)
	suite.failIferror(geoIP.openDBFile(), "Couldn't open DB file")
	cc, didSucceed := geoIP.LookupCountryCodeOfIP(
		net.IPv4(3, 3, 3, 3))
	suite.True(didSucceed)

	// In this case, don't make assertions on a database we don't
	// own. We just want to make sure that we can look things up
	// with no errors.
	fmt.Printf("Country Code of 3.3.3.3: %v\n", cc)

	// Look up a google IP.
	googleIP := "2001:4860:4860::8888"
	cc, didSucceed = geoIP.LookupCountryCodeOfIP(
		net.ParseIP(googleIP))
	suite.True(didSucceed)
	fmt.Printf("Country Code of %s: %v\n", googleIP, cc)

	// Look up a google IP.
	googleIP = "2001:4860:4860::8888"
	cc, didSucceed = geoIP.LookupCountryCodeOfIP(
		net.ParseIP(googleIP))
	suite.True(didSucceed)
	fmt.Printf("Country Code of %s: %v\n", googleIP, cc)
}

// Test that the downloadAndExtractDB() method works -- this will do a
// 'real' download of the MaxMind geoIP country database from
// downloads.untangle.com with an all-zero UID.
func (suite *TestGeoIP) TestDownload() {
	fullFileName := suite.getDBFilename()
	geoIP := NewMaxMindGeoIPManager(fullFileName)
	suite.Nil(geoIP.downloadAndExtractDB())
	suite.failIferror(geoIP.openDBFile(),
		"Couldn't open DB file after download.")
	// Test that we can re-open, which involves closing.
	suite.failIferror(geoIP.openDBFile(),
		"Couldn't open DB file after download (second open).")
}

// Test that if we have the geoIP manager object in a bad state, it
// doesn't panic or throw errors but just acts as if it can't find the
// IP.
func (suite *TestGeoIP) TestNilLookup() {
	fullFileName := "/garbage/garbage2/"
	geoIP := NewMaxMindGeoIPManager(fullFileName)
	result, found := geoIP.LookupCountryCodeOfIP(net.ParseIP("111.111.111.111"))
	suite.Equal(result, "")
	suite.False(found)
}

// Test that we fail gracefully if the tarball is completely
// bad/wrong.
func (suite *TestGeoIP) TestBadTar() {
	suite.testFailure(
		"FakeGEOIP.tar.gz", "bad tar file", true)
}

// Test that we fail gracefully if the file isn't in the tarball at
// all.
func (suite *TestGeoIP) TestMissingFile() {
	suite.testFailure(
		"GeoIP2Missing.tar.gz",
		"tar had missing database file",
		true)
}

// Test that we fail gracefully if we have a database file that exists
// in the tar but it is bad.
func (suite *TestGeoIP) TestBadDBFile() {
	suite.testFailure("GeoIP2BadDB.tar.gz", "tar bad database file", false)
}

// Fail the test if err != nil with msg as a message.
func (suite *TestGeoIP) failIferror(err error, msg string) {
	if err != nil {
		suite.Fail(msg, err)
	}
}

// Get a io.ReadCloser object that wraps a valid test tarball with a
// real database.
func (suite *TestGeoIP) getReaderCloserForTestTarball() io.ReadCloser {
	return suite.getReaderCloserForFile("GeoIP2.tar.gz")
}

// Get a io.ReadCloser for a given named file in test data.
func (suite *TestGeoIP) getReaderCloserForFile(fname string) io.ReadCloser {
	goodTarFile := data.GetTestFileLocation(fname)
	file, err := os.Open(goodTarFile)
	suite.failIferror(err, "Can't open GEOIP tar file.")
	return io.NopCloser(file)

}
func TestGeoIPSuite(t *testing.T) {
	suite.Run(t, &TestGeoIP{})
}
