package database_settings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// unnmarshalTest stolen from policy tests - we can probably generalize this to be reused across areas better
type unmarshalTest struct {
	name        string
	json        string
	expectedErr bool
	expected    Databases
}

// runUnmarshalTest runs the unmarshal test.
func runUnmarshalTest(t *testing.T, tests []unmarshalTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual Databases
			if !tt.expectedErr {
				assert.NoError(t, json.Unmarshal([]byte(tt.json), &actual))
				assert.EqualValues(t, actual, tt.expected)
			} else {
				assert.Error(t, json.Unmarshal([]byte(tt.json), &actual))
			}
		})
	}
}

// TestDatabaseUnmarshal tests unmarshalling the database settings
func TestDatabaseUnmarshal(t *testing.T) {
	var userNames  = []string{"test_user", ""}
    var passwords  = []string{"test_pw", ""}
    var servers  = []string{"test_server", ""}
    var portValues = []int{5,0}
	tests := []unmarshalTest{
		{
			name: "Generic database settings unmarshal test",
			json: `{"databases": [{"enabled": true,
					"db_name": "testingdb",
					"db_username": "test_user",
					"db_password": "test_pw",
					"db_server": "test_server",
					"db_port": 5,
					"description": "Some desc",
					"name": "New DB Source",
					"type": "DB Type 5",
					"default": false,
					"db_connection_string": "asdfasdfasdf"
					}]}`,
			expectedErr: false,
			expected: Databases{
				[]Database{
					{
						Enabled:          true,
						Database:         "testingdb",
						UserName:         &userNames[0],
						Password:         &passwords[0],
						Server:           &servers[0],
						Port:             &portValues[0],
						Description:      "Some desc",
						Name:             "New DB Source",
						Type:             "DB Type 5",
						Default:          false,
						ConnectionString: "asdfasdfasdf",
					},
				},
			},
		},
		{
			name: "Default sqlite DB test",
			json: `{"databases": [
				{
				"db_connection_string": "sqlite:///tmp/reports.db",
				"db_name": "reports.db",
				"db_username": "",
				"db_password": "",
				"db_server": "",
				"db_port": 0,
				"description": "Local reports database",
				"enabled": true,
				"id": "66a6bc90-2f5e-4dc3-8180-a7cf4133daf2",
				"name": "Local DB",
				"type": "sqlite",
				"default": true
				}]}`,
			expectedErr: false,
			expected: Databases{
				[]Database{
					{
						Enabled:          true,
						ConnectionString: "sqlite:///tmp/reports.db",
						Database:         "reports.db",
						UserName:         &userNames[1],
						Password:         &passwords[1],
						Server:           &servers[1],
						Port:             &portValues[1],
						Description:      "Local reports database",
						Name:             "Local DB",
						Type:             "sqlite",
						Default:          true,
						ID:               "66a6bc90-2f5e-4dc3-8180-a7cf4133daf2",
					},
				},
			}},
		{
			name: "bad rule object type",
			json: `{"name": "Geo Rule Tester",
                         "id": "c2428365-65be-4901-bfc0-bde2b310fedf",
                         "type": "asdfasdf",
                         "description": "Whatever",
                         "conditions": ["1458dc12-a9c2-4d0c-8203-1340c61c2c3b"],
                         "action": {
                            "type": "SET_CONFIGURATION",
                            "configuration_id": "1202b42e-2f21-49e9-b42c-5614e04d0031",
                            "key": "GeoipRuleObject"
                            }
                          }`,
			expectedErr: false,
			expected:    Databases{},
		},
	}
	runUnmarshalTest(t, tests)
}
