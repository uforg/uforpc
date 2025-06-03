package lsp

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/urpc/internal/urpc/analyzer"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
)

func TestFindTokenAtPosition(t *testing.T) {
	content := `version 1

type FooType {
  firstField: string
  secondField: int[]
}

proc BarProc {
  input {
    foo: FooType
  }

  output {
    baz: bool
  }
}`

	tests := []struct {
		name     string
		position ast.Position
		want     string
		wantErr  bool
	}{
		{
			name: "Find type name",
			position: ast.Position{
				Line:   3,
				Column: 7,
			},
			want:    "FooType",
			wantErr: false,
		},
		{
			name: "Find field name",
			position: ast.Position{
				Line:   4,
				Column: 5,
			},
			want:    "firstField",
			wantErr: false,
		},
		{
			name: "Find field type",
			position: ast.Position{
				Line:   4,
				Column: 18,
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "Find proc name",
			position: ast.Position{
				Line:   8,
				Column: 7,
			},
			want:    "BarProc",
			wantErr: false,
		},
		{
			name: "Find type reference",
			position: ast.Position{
				Line:   10,
				Column: 12,
			},
			want:    "FooType",
			wantErr: false,
		},
		{
			name: "Position out of range",
			position: ast.Position{
				Line:   100,
				Column: 1,
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findTokenAtPosition(content, tt.position)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHandleTextDocumentDefinition(t *testing.T) {
	// Create a mock reader and writer for the LSP
	mockReader := &bytes.Buffer{}
	mockWriter := &bytes.Buffer{}

	// Create an LSP instance
	lsp := New(mockReader, mockWriter)

	// Create a test schema
	schemaContent := `version 1

type FooType {
  firstField: string
  secondField: int[]
}

proc BarProc {
  input {
    foo: FooType
  }

  output {
    baz: bool
  }
}`

	// Add the schema to the docstore
	filePath := "file:///test.urpc"
	err := lsp.docstore.OpenInMem(filePath, schemaContent)
	require.NoError(t, err)

	// Create a definition request
	request := RequestMessageTextDocumentDefinition{
		RequestMessage: RequestMessage{
			Message: Message{
				JSONRPC: "2.0",
				Method:  "textDocument/definition",
				ID:      "1",
			},
		},
		Params: RequestMessageTextDocumentDefinitionParams{
			TextDocument: TextDocumentIdentifier{
				URI: filePath,
			},
			Position: TextDocumentPosition{
				Line:      9,
				Character: 10,
			},
		},
	}

	// Analyze the file to populate the combined schema
	_, _, err = lsp.analyzer.Analyze(filePath)
	require.NoError(t, err)

	// Encode the request
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)

	// Create a mock analyzer
	mockAnalyzer, err := analyzer.NewAnalyzer(lsp.docstore)
	require.NoError(t, err)
	lsp.analyzer = mockAnalyzer

	// Handle the request
	response, err := lsp.handleTextDocumentDefinition(requestBytes)
	require.NoError(t, err)

	// Check the response
	defResponse, ok := response.(ResponseMessageTextDocumentDefinition)
	require.True(t, ok)
	require.NotNil(t, defResponse.Result)
	require.Len(t, defResponse.Result, 1)
	assert.Equal(t, filePath, defResponse.Result[0].URI)
	assert.Equal(t, 2, defResponse.Result[0].Range.Start.Line) // Line 3 in 0-based indexing
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "Empty string",
			content: "",
			want:    []string{},
		},
		{
			name:    "Single line",
			content: "Hello, world!",
			want:    []string{"Hello, world!"},
		},
		{
			name:    "Multiple lines",
			content: "Line 1\nLine 2\nLine 3",
			want:    []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:    "Windows line endings",
			content: "Line 1\r\nLine 2\r\nLine 3",
			want:    []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:    "Mixed line endings",
			content: "Line 1\nLine 2\r\nLine 3",
			want:    []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:    "Trailing newline",
			content: "Line 1\nLine 2\n",
			want:    []string{"Line 1", "Line 2", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLines(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsIdentifierChar(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want bool
	}{
		{
			name: "Lowercase letter",
			c:    'a',
			want: true,
		},
		{
			name: "Uppercase letter",
			c:    'Z',
			want: true,
		},
		{
			name: "Digit",
			c:    '9',
			want: true,
		},
		{
			name: "Underscore",
			c:    '_',
			want: true,
		},
		{
			name: "Space",
			c:    ' ',
			want: false,
		},
		{
			name: "Symbol",
			c:    '@',
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIdentifierChar(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}
