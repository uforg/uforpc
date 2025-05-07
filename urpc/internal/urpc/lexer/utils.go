package lexer

// isLetter returns true if the character is a letter.
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isNumber returns true if the character is a number.
func isNumber(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isWhitespace returns true if the character is a whitespace. Line breaks
// are not considered whitespace.
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

// isNewline returns true if the character is a newline.
func isNewline(ch rune) bool {
	return ch == '\n'
}

// containsDecimalPoint returns true if the string contains a decimal point.
func containsDecimalPoint(s string) bool {
	for i := range len(s) {
		if s[i] == '.' {
			return true
		}
	}
	return false
}
