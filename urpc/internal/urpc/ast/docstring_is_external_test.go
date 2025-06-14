package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocstringIsExternal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{
			name:     "Valid file",
			input:    "external_doc.md",
			expected: "external_doc.md",
			ok:       true,
		},
		{
			name:     "With spaces",
			input:    "   doc.md   ",
			expected: "doc.md",
			ok:       true,
		},
		{
			name:     "With newline",
			input:    "some\ndoc.md",
			expected: "",
			ok:       false,
		},
		{
			name:     "Does not end with .md",
			input:    "doc.txt",
			expected: "",
			ok:       false,
		},
		{
			name:     "Incorrect suffix",
			input:    ".md",
			expected: "",
			ok:       false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := DocstringIsExternal(tt.input)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.expected, result)
		})
	}
}
