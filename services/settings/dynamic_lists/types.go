package dynamic_lists

// DynamicListSettings is the data structure for the dynamic_lists service, including whether it is enabled and its configurations
type DynamicListSettings struct {
	Enabled        bool                         `json:"enabled"`
	Configurations []*DynamicListConfigurations `json:"configurations"`
}

// DynamicListConfigurations is the data strcuture for dynamic_lists configurations
type DynamicListConfigurations struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	Enabled     bool   `json:"enabled"`
	Source      string `json:"source"`
	PullingUnit string `json:"pullingUnit"`
	PullingTime int    `json:"pullingTime"`
	RegexType   string `json:"regexType"`
}
