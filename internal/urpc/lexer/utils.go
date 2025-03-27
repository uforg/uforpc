package lexer

// isLetter returns true if the character is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isNumber returns true if the character is a number.
func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isWhitespace returns true if the character is a whitespace.
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

// isNewline returns true if the character is a newline.
func isNewline(ch byte) bool {
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
