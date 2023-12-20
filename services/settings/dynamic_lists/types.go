package dynamic_lists

// Example for this struct is present on mfw_schema

// Configs is the data structure for JSON marshalling and unamrshalling Dynamic Lists configurations under dynamic_lists package
type Config struct {
	Name          string `json:"name"`
	ID            string `json:"id"`
	Type          string `json:"type"`
	Enabled       bool   `json:"enabled"`
	Source        string `json:"source"`
	PullingUnit   string `json:"pullingUnit"`
	PullingTime   int    `json:"pullingTime"`
	SkipCertCheck bool   `json:"skipCertCheck"`
	ParsingMethod string `json:"parsingMethod"`
}
