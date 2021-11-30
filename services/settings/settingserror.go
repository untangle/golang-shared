package settings

type SetSettingsError struct {
	Confirm Confirmation `json:"CONFIRM"`
}

type Confirmation struct {
	Rules    map[string]InvalidType `json:"RULES"`
	Policies map[string]InvalidType `json:"POLICIES"`
}

type InvalidType struct {
	AffectedValue      string `json:"affectedValue"`
	InvalidReason      string `json:"invalidReason"`
	invalidReasonType  string `json:"invalidReasonType"`
	InvalidReasonValue string `json:"invalidReasonValue"`
	Parent             string `json:"parent"`
}

type SetSettingsErrorUI struct {
	MainTranslationString string          `json:"mainTranslationString"`
	InvalidReason         string          `json:"invalidReason"`
	AffectedValues        []AffectedValue `json:"affectedValues"`
}

type AffectedValue struct {
	AffectedType  string `json:"affectedType"`
	AffectedValue string `json:"affectedValue"`
}
