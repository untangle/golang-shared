package policy

import (
	"testing"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	policy "github.com/untangle/golang-shared/services/settings/policymanager"
)

// // Hard coded policy settings
// var Policies = []*policySettingsType{
// 	&policySettingsType{
// 		Enabled: true,
// 		Name:    "Teachers",
// 		Source: []string{
// 			"192.168.56.30/32", "192.168.56.31/32",
// 		},
// 	},
// 	&policySettingsType{
// 		Enabled: true,
// 		Name:    "Students",
// 		Source: []string{
// 			"192.168.56.20/32", "192.168.56.21/32",
// 		},
// 	},
// }

func TestPolicyManager(t *testing.T) {
	// Empty settings file expect it to work
	settingsFile := settings.NewSettingsFile("test_settings_empty.json")
	policyMgr := policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s", err)
	}
	if err := policyMgr.ValidatePolicies(); err != nil {
		t.Errorf("ValidatePolicies() failed: %s", err)
	}
	// Good settings file expect it to work
	settingsFile = settings.NewSettingsFile("test_settings.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s", err)
	}
	if err := policyMgr.ValidatePolicies(); err != nil {
		t.Errorf("ValidatePolicies() failed: %s", err)
	}
	// Bad settings files expect errors
	settingsFile = settings.NewSettingsFile("test_settings_ctype.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s\n", err)
	} else if err := policyMgr.ValidatePolicies(); err != nil {
		t.Logf("ValidatePolicies() failed(expected): %s", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad ctype\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_badop.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s\n", err)
	} else if err := policyMgr.ValidatePolicies(); err != nil {
		t.Logf("ValidatePolicies() failed(expected): %s", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad op\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_cfgid.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s\n", err)
	} else if err := policyMgr.ValidatePolicies(); err != nil {
		t.Logf("ValidatePolicies() failed(expected): %s", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad cfgid\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_flowid.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s\n", err)
	} else if err := policyMgr.ValidatePolicies(); err != nil {
		t.Logf("ValidatePolicies() failed(expected): %s", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad flowid\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_ctype.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s\n", err)
	} else if err := policyMgr.ValidatePolicies(); err != nil {
		t.Logf("ValidatePolicies() failed(expected): %s", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad ctype\n")
	}
}

// Concurrency testing using single globalPolicyManager
var globalSettingsFile = settings.NewSettingsFile("test_settings.json")
var globalPolicyMgr = policy.NewPolicyManager(globalSettingsFile, logger.GetLoggerInstance())

func TestLoad(t *testing.T) {
	t.Parallel()
	for i := 0; i < 99; i++ {
		if err := globalPolicyMgr.LoadPolicyManagerSettings(); err != nil {
			t.Errorf("LoadPolicyManagerSettings() failed: %s\n", err)
		}
	}
}

func TestRead1(t *testing.T) {
	t.Parallel()
	for i := 0; i < 99; i++ {
		if err := globalPolicyMgr.ValidatePolicies(); err != nil {
			t.Errorf("ValidatePolicies() failed: %s\n", err)
		}
	}
}

func TestRead2(t *testing.T) {
	t.Parallel()
	for i := 0; i < 99; i++ {
		if err := globalPolicyMgr.ValidatePolicies(); err != nil {
			t.Errorf("ValidatePolicies() failed: %s\n", err)
		}
	}
}
