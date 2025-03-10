package dpiclassmanager

import (
	"testing"

	logService "github.com/untangle/golang-shared/services/logger"
)

// Test that we load the application table.
func TestLoadApplicationTable(t *testing.T) {
	dpi := DpiClassManagerImpl{}
	dpi.DpiConfigFile = "testdata/DpiDefaultConfig.json"
	dpi.ApplicationTable = make(ApplicationTable)
	dpi.logger = logService.GetLoggerInstance()
	err := dpi.LoadApplicationTable()
	if err != nil {
		t.Errorf("LoadApplicationTable failed: %s", err)
	}
	// Check that we have some data
	if len(dpi.ApplicationTable) == 0 {
		t.Errorf("No data loaded")
	}

	// Test GetTable -- application
	data, err := dpi.GetTable("application")
	if err != nil {
		t.Errorf("GetTable failed: %s", err)
	}
	if len(data) == 0 {
		t.Errorf("No data returned")
	}

	// Test GetTable -- category
	data, err = dpi.GetTable("category")
	if err != nil {
		t.Errorf("GetTable failed: %s", err)
	}
	if len(data) == 0 {
		t.Errorf("No data returned")
	}

	// Test GetTable with invalid table
	data, err = dpi.GetTable("invalid")
	if err == nil {
		t.Errorf("GetTable should have failed")
	}
	if len(data) != 0 {
		t.Errorf("Data should be empty")
	}
}
