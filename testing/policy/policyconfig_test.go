package policy

import (
	"fmt"
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
		t.Fail()
	}
	// Good settings file expect it to work
	settingsFile = settings.NewSettingsFile("test_settings.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s", err)
		t.Fail()
	}
	// Bad settings files expect errors
	settingsFile = settings.NewSettingsFile("test_settings_ctype.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		fmt.Printf("LoadPolicyManagerSettings() failed(expected): %s\n", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad ctype\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_badop.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		fmt.Printf("LoadPolicyManagerSettings() failed(expected): %s\n", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad op\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_cfgid.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		fmt.Printf("LoadPolicyManagerSettings() failed(expected): %s\n", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad cfgid\n")
	}
	settingsFile = settings.NewSettingsFile("test_settings_flowid.json")
	policyMgr = policy.NewPolicyManager(settingsFile, logger.GetLoggerInstance())
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		fmt.Printf("LoadPolicyManagerSettings() failed(expected): %s\n", err)
	} else {
		t.Errorf("LoadPolicyManagerSettings() succeeded when it should fail on bad flowid\n")
	}
}
