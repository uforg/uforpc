package lsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	request := ResponseMessage{
		Message: DefaultMessage,
		ID:      IntOrString("1"),
	}
	encoded, err := encode(request)
	require.NoError(t, err)
	require.Equal(t, "Content-Length: 26\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":\"1\"}", string(encoded))
}

func TestDecode(t *testing.T) {
	encoded := []byte("Content-Length: 26\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":\"1\"}")
	expected := ResponseMessage{
		Message: DefaultMessage,
		ID:      IntOrString("1"),
	}

	var decoded ResponseMessage
	require.NoError(t, decode(encoded, &decoded))
	require.Equal(t, expected, decoded)
}
