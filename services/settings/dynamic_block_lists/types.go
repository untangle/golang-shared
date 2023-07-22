package dynamic_block_lists

// Example for this struct is present on mfw_schema

// DynamicBlockListsConfig is the data structure for JSON marshalling and unamrshalling Dynamic Block Lists configurations.
type DynamicBlockListsConfig struct {
	Name           string   `json:"name"`
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Enable         bool     `json:"enable"`
	Source         string   `json:"source"`
	RegexType      string   `json:"regexType"`
	UpdateInterval Interval `json:"interval"`
}

// This is the Interval data structure, which stores the time interval between two update attempts.
type Interval struct {
	Enabled bool `json:"enabled"`
	DayOfWeek int `json:"dayOfWeek"`
	HourOfDay int `json:"hourOfDay"`
	MinuteOfHour int `json:"minuteOfHour"`
}
