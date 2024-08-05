package database_settings

type Database struct {
	Enabled           bool    `json:"enabled"`
	ID                string  `json:"id"`
	Database          string  `json:"db_name"`
	UserName          *string `json:"db_username,omitempty"`
	Password          *string `json:"db_password,omitempty"`
	PasswordEncrypted *string `json:"db_password_encrypted,omitempty"`
	Server            *string `json:"db_server,omitempty"`
	Port              *int    `json:"db_port,omitempty"`
	Description       string  `json:"description"`
	Name              string  `json:"name"`
	Type              string  `json:"type"`
	ConnectionString  string  `json:"db_connection_string"`
	Default           bool    `json:"default"`
	IsDeletable       bool    `json:"is_deletable"`
}
