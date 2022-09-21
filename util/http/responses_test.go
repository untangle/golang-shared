package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResponseDecode -- test code for decoding and unmarshalling
// responses.
func TestResponseDecode(t *testing.T) {
	type resultType struct {
		Value1 string `json:"key1"`
		Value2 int    `json:"key2"`
	}
	type testResults struct {
		// JSON body of response.
		json []byte

		// http status code of response.
		statusCode int

		// are we expected to get an error on initial JSON decode?
		willErrorOnJSONDecode bool

		// are we expected to get an error when unmarshalling?
		willErrorOnUnmarshal bool

		// the expected value of the unmarshalled result, if
		// we don't expect errors.
		expected *resultType

		// Is the response an api error? If so, we want to
		// test the api error methods.
		isAPIErrorResponse bool

		// messages for either type.
		msgs []*Message
	}

	tests := []testResults{
		{
			json: []byte(`{
"result": {"key1": "value1", "key2": 2},
"messages": [{"message": "hello, world!", "vars": ["a", "b"]}]
}`),
			willErrorOnJSONDecode: false,
			expected:              &resultType{Value1: "value1", Value2: 2},
			statusCode:            http.StatusOK,
			msgs: []*Message{
				{Message: "hello, world!", Variables: []string{"a", "b"}},
			},
		},
		{
			json: []byte(
				`{"type": "generic error", "messages": [{"message": "x", "vars": ["y"]}]}`),
			willErrorOnJSONDecode: true,
			statusCode:            http.StatusBadRequest,
			isAPIErrorResponse:    true,
			msgs:                  []*Message{{Message: "x", Variables: []string{"y"}}},
		},
		{
			json:                  []byte(`{`),
			willErrorOnJSONDecode: true,
			statusCode:            http.StatusOK,
			isAPIErrorResponse:    false,
		},
		{
			json: []byte(`{
"result": {"key1": "abcdefghijkmlnopqrstuvwxyz   abc", "key2": 2000},
"messages": [{"message": "hello, world!", "vars": ["a", "b"]}]
}`),
			willErrorOnJSONDecode: false,
			expected: &resultType{
				Value1: "abcdefghijkmlnopqrstuvwxyz   abc",
				Value2: 2000},
			statusCode: http.StatusOK,
			msgs: []*Message{
				{Message: "hello, world!", Variables: []string{"a", "b"}},
			},
		},
		{
			json: []byte(`{
"result": {"BADKEY": "abcdefghijkmlnopqrstuvwxyz   abc", "key2": 2000},
"messages": [{"message": "hello, world!", "vars": ["a", "b"]}]
}`),
			willErrorOnJSONDecode: false,
			willErrorOnUnmarshal:  true,
			statusCode:            http.StatusOK,
			msgs: []*Message{
				{Message: "hello, world!", Variables: []string{"a", "b"}},
			},
		},
		{
			json: []byte(`{
"result": {"key1": "abcdefghijkmlnopqrstuvwxyz   abc", "key2": 2000}
}`),
			willErrorOnJSONDecode: false,
			expected: &resultType{
				Value1: "abcdefghijkmlnopqrstuvwxyz   abc",
				Value2: 2000},
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		jsonValue := test.json
		resp := &http.Response{
			Status:     "OK",
			StatusCode: test.statusCode,
			Body:       ioutil.NopCloser(bytes.NewBuffer(jsonValue)),
		}
		okay, err := DecodeResponse(resp)

		if !test.willErrorOnJSONDecode {
			require.Nil(t, err)
			assert.Equal(t, test.msgs, okay.Messages)

			// test that we can unmarshall the inner
			// result, which is of variable type depending
			// on the api response.
			result := &resultType{}
			err = okay.UnmarshalResult(result)

			// some of our tests expect this to fail.
			if test.willErrorOnUnmarshal {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, result)
			}
		} else {

			assert.NotNil(t, err)
			if test.isAPIErrorResponse {

				assert.True(t, IsApiError(err))

				assert.Regexp(t, `^http response error, type: .*, messages: .*$`,
					err.Error())
				errorResp := ToApiError(err)
				assert.Equal(t, errorResp.Messages, test.msgs)

			}
		}

	}
}
