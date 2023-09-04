package captiveportal

// captive-portal user settings
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
}
