package interface_settings

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
	expected    Interfaces
}

// runUnmarshalTest runs the unmarshal test.
func runUnmarshalTest(t *testing.T, tests []unmarshalTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual Interfaces
			if !tt.expectedErr {
				assert.NoError(t, json.Unmarshal([]byte(tt.json), &actual))
				assert.EqualValues(t, actual, tt.expected)
			} else {
				assert.Error(t, json.Unmarshal([]byte(tt.json), &actual))
			}
		})
	}
}

// TestInterfaceUnmarshal tests unmarshalling the Interface settings
func TestInterfaceUnmarshal(t *testing.T) {
	tests := []unmarshalTest{
		{
			name: "Generic Interface settings unmarshal test",
			json: `{"Interfaces": [{"v4PPPoEPassword": "password", "name": "internal"
					}]}`,
			expectedErr: false,
			expected: Interfaces{
				[]Interface{
					{
						V4PPPoEPassword: "password",
						Name:            "internal",
					},
				},
			},
		},
	}
	runUnmarshalTest(t, tests)
}
