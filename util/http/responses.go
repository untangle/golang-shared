package http

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
