package captiveportal

import (
	"context"
	"net"
)

// Resolver is an interface that both net.Resolver and MockResolverService can implement
type CPResolver interface {
	LookupIP(ctx context.Context, host string, domain string) ([]net.IP, error)
}

// Captive portal condition
type CpRuleCondition struct {
	Op    string `json:"Op,omitempty"`
	Type  string `json:"Type,omitempty"`
	Value any    `json:"Value,omitempty"`
}

// Captive portal action

type CpRulesAction struct {
	Type string `json:"Type,omitempty"`
}

//Captive portal rules

type CpRules struct {
	RuleId      string             `json:"RuleId,omitempty"`
	Enabled     bool               `json:"Enabled,omitempty"`
	Description string             `json:"Description,omitempty"`
	Conditions  []*CpRuleCondition `json:"Conditions,omitempty"`
	Action      *CpRulesAction     `json:"Action,omitempty"`
}

type CpSettingType struct {
	Enabled          bool       `json:"Enabled,omitempty"`
	TimeoutValue     int32      `json:"TimeoutValue,omitempty"`
	TimeoutPeriod    string     `json:"TimeoutPeriod,omitempty"`
	AcceptText       string     `json:"AcceptText,omitempty"`
	AcceptButtonText string     `json:"AcceptButtonText,omitempty"`
	MessageHeading   string     `json:"MessageHeading,omitempty"`
	MessageText      string     `json:"MessageText,omitempty"`
	WelcomeText      string     `json:"WelcomeText,omitempty"`
	PageTitle        string     `json:"PageTitle,omitempty"`
	Rules            []*CpRules `json:"Rules,omitempty"`
}

func (x *CpSettingType) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *CpSettingType) GetTimeoutValue() int32 {
	if x != nil {
		return x.TimeoutValue
	}
	return 0
}

func (x *CpSettingType) GetTimeoutPeriod() string {
	if x != nil {
		return x.TimeoutPeriod
	}
	return ""
}

func (x *CpSettingType) GetAcceptText() string {
	if x != nil {
		return x.AcceptText
	}
	return ""
}

func (x *CpSettingType) GetAcceptButtonText() string {
	if x != nil {
		return x.AcceptButtonText
	}
	return ""
}

func (x *CpSettingType) GetMessageHeading() string {
	if x != nil {
		return x.MessageHeading
	}
	return ""
}

func (x *CpSettingType) GetMessageText() string {
	if x != nil {
		return x.MessageText
	}
	return ""
}

func (x *CpSettingType) GetWelcomeText() string {
	if x != nil {
		return x.WelcomeText
	}
	return ""
}

func (x *CpSettingType) GetPageTitle() string {
	if x != nil {
		return x.PageTitle
	}
	return ""
}

func (x *CpSettingType) GetRules() []*CpRules {
	if x != nil {
		return x.Rules
	}
	return nil
}
