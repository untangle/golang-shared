package database_settings

type Database struct {
	Enabled          bool   `json:"enabled"`
	ID               string `json:"id"`
	Database         string `json:"db_name"`
	UserName         string `json:"db_username"`
	Password         string `json:"db_password"`
	Server           string `json:"db_server"`
	Port             int    `json:"db_port"`
	Description      string `json:"description"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"db_connection_string"`
	Default          bool   `json:"default"`
}
