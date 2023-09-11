package policy

// Policies are the root of our policy configurations. It includes pointers to substructure.
type Policy struct {
	Defaults    bool     `json:"defaults"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Rules       []string `json:"rules"`

	// DEPRECATED
	Configurations []string `json:"configurations"`
	Flows          []string `json:"flows"`
}
