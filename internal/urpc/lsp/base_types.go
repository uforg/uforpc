package lsp

import (
	"encoding/json"
	"fmt"
)

var (
	DefaultMessage = Message{
		JSONRPC: "2.0",
	}
)

// Message is a general message as defined by JSON-RPC. The language server protocol
// always uses “2.0” as the jsonrpc version.
type Message struct {
	JSONRPC string `json:"jsonrpc"`
}

// IntOrString is a type that can be either an int or a string when unmarshalled from JSON.
type IntOrString string

// UnmarshalJSON implements the json.Unmarshaler interface.
// The value is parsed as an int if possible, otherwise it is parsed as a string.
func (i *IntOrString) UnmarshalJSON(b []byte) error {
	var n int
	if err := json.Unmarshal(b, &n); err == nil {
		*i = IntOrString(fmt.Appendf(nil, "%d", n))
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*i = IntOrString(s)
		return nil
	}

	return fmt.Errorf("IntOrString: %s is not an int or a string", string(b))
}

// RequestMessage describes a request between the client and the server. Every processed
// request must send a response back to the sender of the request.
type RequestMessage struct {
	Message
	// The request id.
	ID IntOrString `json:"id"`
	// The method to be invoked.
	Method string `json:"method"`
}

// ResponseMessage sent as a result of a request.
type ResponseMessage struct {
	Message
	// The request id.
	ID IntOrString `json:"id"`
	// The error object.
	Error ResponseError `json:"error,omitzero"`
}

// ResponseError is an error that occurred while processing a request.
type ResponseError struct {
	// A number indicating the error type that occurred.
	Code int `json:"code"`
	// A string providing a short description of the error.
	Message string `json:"message"`
	// Additional error data. Can be omitted.
	Data map[string]any `json:"data"`
}

// NotificationMessage describes a notification between the client and the server.
type NotificationMessage struct {
	Message
	// The method to be invoked.
	Method string `json:"method"`
}
