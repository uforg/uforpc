package strutil

import (
	"strings"
	"unicode"
)

// ToPascalCase converts a string to PascalCase, it will interpret all
// space like characters, underscores and dashes as word boundaries.
func ToPascalCase(str string) string {
	result := strings.Builder{}
	newWord := true

	for _, char := range str {
		if unicode.IsSpace(char) || char == '_' || char == '-' {
			newWord = true
			continue
		}
		if newWord {
			char = unicode.ToUpper(char)
			newWord = false
		} else {
			char = unicode.ToLower(char)
		}
		result.WriteRune(char)
	}

	return result.String()
}
