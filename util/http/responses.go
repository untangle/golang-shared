package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// Message is a message from the API.
type Message struct {
	Message   string   `json:"message"`
	Variables []string `json:"vars"`
}

// OkayResponse is a response with no error (not necessarily just 200).
type OkayResponse struct {
	Result   interface{} `json:"result"`
	Messages []*Message  `json:"messages"`
}

// ErrorResponse is only for errors.
type ErrorResponse struct {
	Type     string     `json:"type"`
	Messages []*Message `json:"messages"`
}

func (e *ErrorResponse) Error() string {
	msgString := ""
	for _, msg := range e.Messages {
		msgString += fmt.Sprintf("%+v", msg)
	}
	return fmt.Sprintf("http response error, type: %s, messages: %+v", e.Type, msgString)
}

func ToApiError(err error) *ErrorResponse {
	switch e := err.(type) {
	case *ErrorResponse:
		return e
	default:
		return nil
	}
}

func IsApiError(err error) bool {
	return ToApiError(err) != nil
}

func unmarshalResponse(
	resp *http.Response,
	value interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(
			"unable to read http response: %w",
			err)
	}
	err = json.Unmarshal(body, value)
	if err != nil {
		return fmt.Errorf(
			"unable to unmarshal JSON for http response: %w",
			err)
	}
	return nil
}

// DecodeResponse decodes the HTTP response and returns the response,
// or error. If the returned response is an http error, we return the
// ErrorResponse, which implements the error interface.
func DecodeResponse(resp *http.Response) (*OkayResponse, error) {
	if resp.StatusCode < http.StatusMultipleChoices && resp.StatusCode >= http.StatusOK {
		ok := &OkayResponse{}
		err := unmarshalResponse(resp, ok)
		if err != nil {
			return nil, fmt.Errorf("decoding successful http response failed: %w", err)
		}
		return ok, nil
	}
	errResp := &ErrorResponse{}
	err := unmarshalResponse(resp, errResp)
	if err != nil {
		return nil, fmt.Errorf("decoding error http response failed: %w", err)
	}
	return nil, errResp

}

// UnmarshalResult will unmarshall the Result member of the
// OkayResponse using mapstructure into the result interface. It will
// use the JSON struct tags to unmarshal.
func (okay *OkayResponse) UnmarshalResult(result interface{}) error {
	decoderConfig := mapstructure.DecoderConfig{
		TagName:     "json",
		Result:      result,
		ErrorUnused: true,
		Squash:      true,
	}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return fmt.Errorf("responses: unable to create decoder: %w", err)
	}
	if err := decoder.Decode(okay.Result); err != nil {
		return fmt.Errorf("unable to decode http okay response: %w", err)
	}
	return nil
}
