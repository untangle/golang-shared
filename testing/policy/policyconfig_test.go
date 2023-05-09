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
	settingsFile := settings.NewSettingsFile("test_settings_empty.json")
	logger := logger.GetLoggerInstance()
	policyMgr := policy.NewPolicyManager(settingsFile, *logger)
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s", err)
		t.Fail()
	}

	settingsFile = settings.NewSettingsFile("test_settings.json")
	policyMgr = policy.NewPolicyManager(settingsFile, *logger)
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Errorf("LoadPolicyManagerSettings() failed: %s", err)
		t.Fail()
	}
}
