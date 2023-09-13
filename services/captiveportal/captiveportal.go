package captiveportal

// Captive portal rule conditions
type CpRulesConditions struct {
	Op    string `json:"op"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Captive portal rule actions
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

// Captive portal user settings
type CpSettingType struct {
	Enabled          bool   `json:"enabled"`
	TimeoutValue     int16  `json:"timeoutValue"`
	TimeoutPeriod    string `json:"timeoutPeriod"`
	AcceptText       string `json:"acceptText"`
	AcceptButtonText string `json:"acceptButtonText"`
	MessageHeading   string `json:"messageHeading"`
	MessageText      string `json:"messageText"`
	WelcomeText      string `json:"welcomeText"`
	PageTitle        string `json:"pageTitle"`
	Base64ImageData  struct {
		EncodedBase64Image string `json:"imageData"`
		ImageName          string `json:"imageName"`
	} `json:"logo"`
	Rules []*CpRules `json:"rules"`
}