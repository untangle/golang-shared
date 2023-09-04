package captiveportal

type CpRulesConditions struct {
	Op    string `json:"op"`
	Type  string `json:"type"`
	Value string `json:"value"`
	// CompareValue is used when evaluating conditions,
	// the Value is compared against this field using Op operator.
	CompareValue any `json:"-"`
}

// CpRulesAction
type CpRulesAction struct {
	Type string `json:"type"`
}

var (
	CpRulesEnable  = CpRulesAction{"ENABLE"}
	CpRulesDisable = CpRulesAction{"DISABLE"}
)

// Captive portal rules
type CpRules struct {
	RuleId      string               `json:"rule_id"`
	Enabled     bool                 `json:"enabled"`
	Description string               `json:"description"`
	Conditions  []*CpRulesConditions `json:"conditions"`
	Action      CpRulesAction
}

// captive-portal user settings
type CpSettingType struct {
	Enabled          bool   `json:"enabled"`
	TimeoutValue     int16  `json:"timeoutValue"`
	TimeoutPeriod    string `json:"timeoutPeriod"`
	AcceptText       string `json:"acceptText"`
	AcceptButtonText string `json:"acceptButtonText"`
	MessageText      string `json:"messageText"`
	TosText          string `json:"tosText"`
	WelcomeText      string `json:"welcomeText"`
	Base64ImageData  struct {
		EncodedBase64Image string `json:"imageData"`
		ImageName          string `json:"imageName"`
	} `json:"logo"`
	Rules []*CpRules `json:"rules"`
}
