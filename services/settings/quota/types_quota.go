package quota

import (
    "time"
)

type QuotaConditions struct {
	Type  string `json:"type"`
	Op    string `json:"op"`
	Value string `json:"value"`
	Proto string `json:"proto"`
}

type QuotaExceedActions struct {
	ID          string        `json:"id"`
	Name        string        `json:name`
	Description string        `json:description`
	Action      string        `json:action`
	Priority    string        `json:priroity`
	Default     bool          `json:enabled`
}

type QuotaAction struct {
	Type             string    `json:"type"`
	QuotaID          string    `json:"quota_id"`
	ExceedActionID   string    `json:"exceed_action_id"`
}

type QuotaConfiguration struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Period       string `json:"period"`
	DataSize     string `json:"datasize"`
}

type QuotaRules struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Conditions  []*QuotaConditions `json:"conditions"`
	Actions     *QuotaAction       `json:"action"`

	Command     string
	MatchCmd    string
	Timer       *time.Ticker
	ChSignal    chan bool
}

// QuotaSettings is the main data structure for Quota Management.
// It contains an array of QuotaConfigurations, an array of QuotaRules and an array
// of exceed_actions.
type QuotaSettings struct {
	Enabled           bool                     `json:"enabled"`
	Configurations    []*QuotaConfiguration    `json:"configuration"`
	Rules             []*QuotaRules            `json:"rules"`
	ExceedActions     []*QuotaExceedActions    `json:"exceed_actions"`
}
