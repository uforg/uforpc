package lexer

// isLetter returns true if the character is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isDigit returns true if the character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isWhitespace returns true if the character is a whitespace.
func isWhitespace(ch byte) bool {
	return ch == ' '
}

// isNewline returns true if the character is a newline.
func isNewline(ch byte) bool {
	return ch == '\n'
}
