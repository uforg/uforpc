package strutil

import (
	"strings"
	"unicode"
)

// ToPascalCase converts a string to PascalCase, it will interpret all
// space like characters, underscores and dashes as word boundaries.
//
// Example:
//
//	"hello world" -> "HelloWorld"
//	"hello_world" -> "HelloWorld"
//	"hello-world" -> "HelloWorld"
//	"hello world" -> "HelloWorld"
//	"hello WORLD" -> "HelloWorld"
//	"helloWORLD"  -> "HelloWorld"
func ToPascalCase(str string) string {
	if str == "" {
		return ""
	}

	// First, split the string into words
	var words []string
	word := strings.Builder{}
	for i, char := range str {
		if unicode.IsSpace(char) || char == '_' || char == '-' {
			if word.Len() > 0 {
				words = append(words, word.String())
				word.Reset()
			}
		} else {
			word.WriteRune(char)
		}

		// Handle the last word
		if i == len(str)-1 && word.Len() > 0 {
			words = append(words, word.String())
		}
	}

	// Then, convert each word to PascalCase
	result := strings.Builder{}
	for _, w := range words {
		if len(w) == 0 {
			continue
		}

		// Always capitalize the first character of each word
		firstChar := unicode.ToUpper(rune(w[0]))
		result.WriteRune(firstChar)

		// For the rest of the word, preserve existing uppercase letters
		// that are not at word boundaries
		if len(w) > 1 {
			for i := 1; i < len(w); i++ {
				char := rune(w[i])
				prevChar := rune(w[i-1])

				// If current char is uppercase and previous char is lowercase,
				// keep it uppercase (camelCase pattern)
				if unicode.IsUpper(char) && unicode.IsLower(prevChar) {
					result.WriteRune(char)
				} else {
					result.WriteRune(unicode.ToLower(char))
				}
			}
		}
	}

	return result.String()
}
