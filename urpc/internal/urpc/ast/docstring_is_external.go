package ast

import "strings"

// DocstringIsExternal checks if a docstring is an external markdown file.
//
// If it is, it returns the trimmed docstring and true.
// If it is not, it returns an empty string and false.
func DocstringIsExternal(docstring string) (string, bool) {
	trimmed := strings.TrimSpace(docstring)
	if strings.ContainsAny(trimmed, "\r\n") {
		return "", false
	}

	if strings.TrimSuffix(".md", trimmed) == "" {
		return "", false
	}

	if !strings.HasSuffix(trimmed, ".md") {
		return "", false
	}

	return trimmed, true
}
