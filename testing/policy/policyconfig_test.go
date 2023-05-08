package policy

import (
	"testing"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	policy "github.com/untangle/golang-shared/services/settings/policymanager"
)

func TestPolicyManager(t *testing.T) {
	settingsFile := settings.NewSettingsFile("test_settings_empty.json")
	logger := logger.GetLoggerInstance()
	policyMgr := policy.NewPolicyManager(settingsFile, *logger)
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Fail()
	}

	settingsFile = settings.NewSettingsFile("test_settings.json")
	policyMgr = policy.NewPolicyManager(settingsFile, *logger)
	if err := policyMgr.LoadPolicyManagerSettings(); err != nil {
		t.Fail()
	}
}
