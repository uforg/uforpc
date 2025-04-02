package lsp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// encode encodes the given data into a valid LSP JSON-RPC message and returns
// the encoded message as a byte slice.
func encode(data any) ([]byte, error) {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	contentLength := len(marshaled)
	content := string(marshaled)

	return fmt.Appendf(nil, "Content-Length: %d\r\n\r\n%s", contentLength, content), nil
}

// decode decodes the given data into the given value. It expects the data to be
// a valid LSP JSON-RPC message.
func decode(data []byte, v any) error {
	if !bytes.HasPrefix(data, []byte("Content-Length: ")) {
		return fmt.Errorf("invalid LSP JSON-RPC message")
	}

	segments := bytes.Split(data, []byte("\r\n\r\n"))
	if len(segments) != 2 {
		return fmt.Errorf("invalid LSP JSON-RPC message")
	}

	if err := json.Unmarshal(segments[1], &v); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return nil
}
