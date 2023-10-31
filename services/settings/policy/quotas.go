package policy

import (
	"fmt"
	"strings"
	"time"
)

// QuotaRefreshTime is a refresh time (time.Duration)
type QuotaRefreshTime time.Duration

// QuotaSettings are settings for a quota, the amount of bytes and the refresh interval.
type QuotaSettings struct {
	AmountBytes     uint64           `json:"amount_bytes"`
	RefreshInterval QuotaRefreshTime `json:"refresh"`
}

// Quota is an object with Type = QuotaType
type Quota Object

// GetSettings gets the QuotaSettings from the SettingsField.
func (q *Quota) GetSettings() *QuotaSettings {
	settings, ok := q.Settings.(*QuotaSettings)
	if !ok {
		return nil
	}
	return settings
}

// UnmarshalJSON: unmarshal a quota.
func (q *QuotaRefreshTime) UnmarshalJSON(b []byte) error {
	// remove surrounding double quotes, treat as string.
	str := strings.Trim(string(b), "\"")
	if t, err := time.ParseDuration(str); err != nil {
		return fmt.Errorf("unable to parse quota time: %w", err)
	} else {
		*q = QuotaRefreshTime(t)
	}
	return nil
}
