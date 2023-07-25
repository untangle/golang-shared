package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Test suite for testing path unmarshaller.
type TestJSONUnmarshalSuite struct {
	suite.Suite
}

func (suite *TestJSONUnmarshalSuite) TestUnmarshalling() {
	type dummy struct {
		Value1 int64  `json:"value1"`
		Value2 string `json:"value2"`
	}
	tests := []struct {
		name           string
		jsonString     string
		expectedOutput any
		path           []string
		outputVar      any
	}{
		{
			name:           "simple with one object",
			path:           []string{"one"},
			jsonString:     `{"one": "two"}`,
			expectedOutput: "two",
			outputVar:      nil,
		},
		{
			name:       "simple multi-component object path",
			path:       []string{"one", "two", "three"},
			jsonString: `{"one": {"two": {"three": {"x": 1, "y": 2}}}}`,
			outputVar:  nil,
			expectedOutput: map[string]interface{}{
				"x": json.Number("1"),
				"y": json.Number("2"),
			},
		},
		{
			name:           "struct unmarshal",
			jsonString:     `{"": null, "x": false, "path1": {"path2": {"value1": 222, "value2": "hello"}}}`,
			path:           []string{"path1", "path2"},
			outputVar:      &dummy{},
			expectedOutput: &dummy{Value1: 222, Value2: "hello"},
		},
		{
			name:           "root unmarshal",
			jsonString:     `{"value1": 222, "value2": "hello", "ignored": 3030}`,
			outputVar:      &dummy{},
			expectedOutput: &dummy{Value1: 222, Value2: "hello"},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			unmarshaller := suite.getUnmarshallerForString(tt.jsonString)
			unmarshaller.UseNumber()
			if tt.outputVar != nil {
				suite.NoError(unmarshaller.UnmarshalAtPath(tt.outputVar, tt.path...))
			} else {
				suite.NoError(unmarshaller.UnmarshalAtPath(&tt.outputVar, tt.path...))
			}
			suite.EqualValues(tt.expectedOutput, tt.outputVar)
		})
	}

}

// Test that we catch various error cases and do not panic.
func (suite *TestJSONUnmarshalSuite) TestBadJSON() {
	for _, i := range []string{
		`}}`,
		`{"valid": "json", "invalid": }}}`,
		`[]`,
		`{"some": {"path"}}`,
		`{{{`,
		`{{]]`,
		`{{[[[]]]}}`,
	} {
		var output interface{}
		unmarshaller := suite.getUnmarshallerForString(i)
		err := unmarshaller.UnmarshalAtPath(output, "some", "path")
		suite.NotNil(err)
		fmt.Printf("[ok]: Caught expected error, string value: '%s'\n", err)
	}
}

// Test that deserializing nested structures with slices work.
func (suite *TestJSONUnmarshalSuite) TestNestedStructs() {
	type Inner struct {
		Name   string `json:"name"`
		Status int32  `json:"status"`
		IsCool bool   `json:"is_cool"`
	}
	type Outer struct {
		Key    string  `json:"key"`
		Inners []Inner `json:"inners"`
	}

	output := Outer{}

	json := `
{"ignore": null,
 "large_object": [
1,
2,
3,
4,
5,
6,
7,
1,
2,
3,
4,
5,
6,
7
],
 "otherIgnore": {"key": false, "key2": false},
 "pathComp1": {
   "pathComp2": {
     "pathComp3": {
       "key": "hello!",
       "inners": [
         {"name": "world!", "status": 42, "is_cool": true},
         {"name": "doge", "status": 99, "is_cool": true},
         {"name": "goop", "status": 9000, "is_cool": true},
         {"name": "glop", "status": 0, "is_cool": false},
         {"name": "doge", "status": 998, "is_cool": true},
         {"name": "goop", "status": 9009, "is_cool": true},
         {"name": "glop", "status": 2, "is_cool": false},
         {"name": "doge", "status": 998, "is_cool": true},
         {"name": "goop", "status": 9009, "is_cool": true},
         {"name": "glop", "status": 2, "is_cool": false},
         {"name": "doge", "status": 998, "is_cool": true},
         {"name": "goop", "status": 9009, "is_cool": true},
         {"name": "glop", "status": 2, "is_cool": false}
       ]
     }
   }
 }
}
`
	unm := suite.getUnmarshallerForString(json)
	suite.Nil(unm.UnmarshalAtPath(&output, "pathComp1", "pathComp2", "pathComp3"))
	suite.Equal(output.Key, "hello!")
	suite.Equal(len(output.Inners), 13)
	expected := []Inner{
		{Name: "world!", Status: 42, IsCool: true},
		{Name: "doge", Status: 99, IsCool: true},
		{Name: "goop", Status: 9000, IsCool: true},
		{Name: "glop", Status: 0, IsCool: false},
		{Name: "doge", Status: 998, IsCool: true},
		{Name: "goop", Status: 9009, IsCool: true},
		{Name: "glop", Status: 2, IsCool: false},
		{Name: "doge", Status: 998, IsCool: true},
		{Name: "goop", Status: 9009, IsCool: true},
		{Name: "glop", Status: 2, IsCool: false},
		{Name: "doge", Status: 998, IsCool: true},
		{Name: "goop", Status: 9009, IsCool: true},
		{Name: "glop", Status: 2, IsCool: false},
	}
	suite.Equal(expected, output.Inners)
}

func (suite *TestJSONUnmarshalSuite) getUnmarshallerForString(json string) *PathUnmarshaller {
	reader := bytes.NewReader([]byte(json))
	return NewPathUnmarshaller(reader)
}

func TestUnmarshaller(t *testing.T) {
	testSuite := &TestJSONUnmarshalSuite{}
	suite.Run(t, testSuite)
}
