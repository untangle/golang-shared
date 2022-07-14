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

// Test the most basic schenario -- get a value for a key.
func (suite *TestJSONUnmarshalSuite) TestSimpleUnmarshalPath() {
	jsonSimple := `{"one": "two"}`
	unmarshaller := suite.getUnmarshallerForString(jsonSimple)
	var output interface{}
	suite.Nil(unmarshaller.UnmarshalAtPath(&output, "one"))
	stringOutput := output.(string)
	suite.Equal(stringOutput, "two")
}

// Test a more nested path.
func (suite *TestJSONUnmarshalSuite) TestBiggerPath() {
	jsonAdvanced := `{"one": {"two": {"three": {"x": 1, "y": 2}}}}`
	var output interface{} = nil
	unmarshaller := suite.getUnmarshallerForString(jsonAdvanced)
	unmarshaller.UseNumber()
	suite.Nil(unmarshaller.UnmarshalAtPath(&output, "one", "two", "three"))
	mapOutput := output.(map[string]interface{})
	xValue, err := mapOutput["x"].(json.Number).Int64()
	suite.Nil(err)
	yValue, err := mapOutput["y"].(json.Number).Int64()
	suite.Nil(err)
	suite.Equal(int(xValue), 1)
	suite.Equal(int(yValue), 2)
}

// Test that basic structure unmarshalling works.
func (suite *TestJSONUnmarshalSuite) TestStructUnmarshal() {

	type Dummy struct {
		Value1 int64  `json:"value1"`
		Value2 string `json:"value2"`
	}

	jsonDummy := `{"": null, "x": false, "path1": {"path2": {"value1": 222, "value2": "hello"}}}`
	output := &Dummy{}
	unmarshaller := suite.getUnmarshallerForString(jsonDummy)
	suite.Nil(unmarshaller.UnmarshalAtPath(output, "path1", "path2"))
	suite.Equal(output.Value1, int64(222))
	suite.Equal(output.Value2, "hello")
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
 "otherIgnore": {"key": false, "key2": false},
 "pathComp1": {
   "pathComp2": {
     "pathComp3": {
       "key": "hello!",
       "inners": [
         {"name": "world!", "status": 42, "is_cool": true},
         {"name": "doge", "status": 99, "is_cool": true},
         {"name": "goop", "status": 9000, "is_cool": true},
         {"name": "glop", "status": 0, "is_cool": false}
       ]
     }
   }
 }
}
`
	unm := suite.getUnmarshallerForString(json)
	suite.Nil(unm.UnmarshalAtPath(&output, "pathComp1", "pathComp2", "pathComp3"))
	suite.Equal(output.Key, "hello!")
	suite.Equal(len(output.Inners), 4)
	expected := []Inner{
		{Name: "world!", Status: 42, IsCool: true},
		{Name: "doge", Status: 99, IsCool: true},
		{Name: "goop", Status: 9000, IsCool: true},
		{Name: "glop", Status: 0, IsCool: false},
	}
	suite.Equal(expected, output.Inners)
}

func (suite *TestJSONUnmarshalSuite) getUnmarshallerForString(json string) *PathUnmarshaller {
	reader := bytes.NewBuffer([]byte(json))
	return NewPathUnmarshaller(reader)
}

func TestUnmarshaller(t *testing.T) {
	testSuite := &TestJSONUnmarshalSuite{}
	suite.Run(t, testSuite)
}
