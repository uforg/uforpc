package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocstringGetExternal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{"valid unix path", "./docs/readme.md", "./docs/readme.md", true},
		{"valid absolute unix path", "/usr/local/readme.md", "/usr/local/readme.md", true},
		{"valid windows path", `C:\docs\readme.md`, `C:\docs\readme.md`, true},
		{"invalid extension", "./docs/readme.txt", "", false},
		{"empty string", "", "", false},
		{"only whitespace", "   ", "", false},
		{"newline at end", "./docs/readme.md\n", "./docs/readme.md", true},
		{"newline at start", "\n./docs/readme.md", "./docs/readme.md", true},
		{"newline in middle", "./docs/\nreadme.md", "", false},
		{"carriage return at end", "./docs/readme.md\r", "./docs/readme.md", true},
		{"carriage return in middle", "./docs/\rreadme.md", "", false},
		{"uppercase extension", "./docs/README.MD", "", false},
		{"leading and trailing whitespace", "  ./docs/readme.md  ", "./docs/readme.md", true},
		{"just .md", ".md", "", false},
		{"dotfile but not markdown", ".gitignore", "", false},
		{"directory with dot", "./.config/readme.md", "./.config/readme.md", true},
		{"tricky valid path with newline padding", "\n  ./dir/file.md  \r\n", "./dir/file.md", true},
		{"path with tab in middle", "./docs/\treadme.md", "./docs/\treadme.md", true},
	}

	for _, tt := range tests {
		d := Docstring{Value: tt.input}
		path, ok := d.GetExternal()
		require.Equal(t, tt.ok, ok, tt.name)
		require.Equal(t, tt.expected, path, tt.name)
	}
}
