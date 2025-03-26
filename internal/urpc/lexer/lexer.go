package lexer

import "github.com/uforg/uforpc/internal/urpc/token"

type Lexer struct {
	input             string
	fileName          string
	currentLine       int
	currentColumn     int
	maxIndex          int
	currentIndex      int
	currentIndexIsEOF bool
	currentChar       byte
}

// NewLexer creates a new Lexer and initializes it with the given file
// name and input string.
func NewLexer(fileName, input string) *Lexer {
	l := &Lexer{}

	l.input = input
	l.fileName = fileName
	l.currentLine = 1
	l.currentColumn = 1
	l.maxIndex = len(input) - 1

	l.currentIndex = 0
	if l.maxIndex <= 0 {
		l.currentIndexIsEOF = true
	} else {
		l.currentIndexIsEOF = false
	}
	if l.currentIndexIsEOF {
		l.currentChar = 0
	} else {
		l.currentChar = input[l.currentIndex]
	}

	return l
}

// readNextChar reads the next character from the input and
// updates the current Lexer state
func (l *Lexer) readNextChar() {
	if l.currentIndexIsEOF {
		return
	}

	if isNewline(l.currentChar) {
		l.currentLine++
		l.currentColumn = 1
	} else {
		l.currentColumn++
	}

	l.currentIndex++
	if l.currentIndex > l.maxIndex {
		l.currentIndexIsEOF = true
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.currentIndex]
	}
}

// peekChar peeks the character at the next index + depth without moving the current index.
//
// Returns the character and a boolean indicating if the EOF was reached.
func (l *Lexer) peekChar(depth int) (byte, bool) {
	indexToPeek := l.currentIndex + depth
	if indexToPeek > l.maxIndex {
		return 0, true
	}
	return l.input[indexToPeek], false
}

// readIdentifier reads an identifier from the current index to the next non-letter character.
func (l *Lexer) readIdentifier() string {
	var ident string
	for isLetter(l.currentChar) {
		ident += string(l.currentChar)

		nextChar, eofReached := l.peekChar(1)
		if eofReached || !isLetter(nextChar) {
			break
		}

		l.readNextChar()
	}
	return ident
}

// readNumber reads a number from the current index to the next non-digit character.
func (l *Lexer) readNumber() string {
	var num string
	for isNumber(l.currentChar) {
		num += string(l.currentChar)

		nextChar, eofReached := l.peekChar(1)
		if eofReached || !isNumber(nextChar) {
			break
		}

		l.readNextChar()
	}
	return num
}

// readString reads a string from the current index to the next double quote.
//
// Returns the string and a boolean indicating if the string is unterminated.
func (l *Lexer) readString() (string, bool) {
	if l.currentChar != '"' {
		return "", false
	}

	// Skip the opening quote
	l.readNextChar()

	var str string
	for !l.currentIndexIsEOF && l.currentChar != '"' {
		str += string(l.currentChar)
		l.readNextChar()
	}

	if l.currentIndexIsEOF {
		return str, true
	}

	return str, false
}

// readComment reads a comment from the current index to the next newline or EOF.
// it does not skip the end newline character.
func (l *Lexer) readComment() string {
	if l.currentChar != '/' {
		return ""
	}

	nextChar, eofReached := l.peekChar(1)
	if eofReached || nextChar != '/' {
		return ""
	}

	// Skip the opening slashes
	l.readNextChar()
	l.readNextChar()

	var comment string
	for {
		comment += string(l.currentChar)

		nextChar, eofReached := l.peekChar(1)
		if eofReached || nextChar == '\n' {
			break
		}

		l.readNextChar()
	}

	return comment
}

// skipWhitespace skips whitespace characters from the current index to the next non-whitespace character.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.currentChar) {
		l.readNextChar()
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	if l.currentIndexIsEOF {
		return token.Token{
			Type:     token.EOF,
			Literal:  "",
			FileName: l.fileName,
			Line:     l.currentLine,
			Column:   l.currentColumn,
		}
	}

	defer l.readNextChar()
	l.skipWhitespace()

	// Handle delimiters
	if token.IsDelimiter(l.currentChar) {
		return token.Token{
			Type:     token.GetDelimiterTokenType(l.currentChar),
			Literal:  string(l.currentChar),
			FileName: l.fileName,
			Line:     l.currentLine,
			Column:   l.currentColumn,
		}
	}

	// Handle strings
	if l.currentChar == '"' {
		startLine := l.currentLine
		startColumn := l.currentColumn
		str, unterminated := l.readString()

		if unterminated {
			return token.Token{
				Type:     token.ILLEGAL,
				Literal:  "\"" + str,
				FileName: l.fileName,
				Line:     startLine,
				Column:   startColumn,
			}
		}

		return token.Token{
			Type:     token.STRING,
			Literal:  str,
			FileName: l.fileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Handle ints and floats
	if isNumber(l.currentChar) {
		startLine := l.currentLine
		startColumn := l.currentColumn
		num := l.readNumber()

		tok := token.Token{
			Type:     token.INT,
			Literal:  num,
			FileName: l.fileName,
			Line:     startLine,
			Column:   startColumn,
		}

		nextChar, eofReached := l.peekChar(1)
		if eofReached || nextChar != '.' {
			return tok
		}

		nextChar, eofReached = l.peekChar(2)
		if eofReached || !isNumber(nextChar) {
			return tok
		}

		// Double read, one for the dot (should be ignored) and one for the start
		// of the next number
		l.readNextChar()
		l.readNextChar()

		num += "." + l.readNumber()
		return token.Token{
			Type:     token.FLOAT,
			Literal:  num,
			FileName: l.fileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Handle identifiers
	if isLetter(l.currentChar) {
		startLine := l.currentLine
		startColumn := l.currentColumn
		ident := l.readIdentifier()

		if token.IsKeyword(ident) {
			return token.Token{
				Type:     token.GetKeywordTokenType(ident),
				Literal:  ident,
				FileName: l.fileName,
				Line:     startLine,
				Column:   startColumn,
			}
		}

		return token.Token{
			Type:     token.IDENT,
			Literal:  ident,
			FileName: l.fileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Handle comments
	if l.currentChar == '/' {
		nextChar, eofReached := l.peekChar(1)
		if eofReached || nextChar != '/' {
			return token.Token{
				Type:     token.ILLEGAL,
				Literal:  string(l.currentChar) + string(nextChar),
				FileName: l.fileName,
				Line:     l.currentLine,
				Column:   l.currentColumn,
			}
		}

		startLine := l.currentLine
		startColumn := l.currentColumn
		comment := l.readComment()
		return token.Token{
			Type:     token.COMMENT,
			Literal:  comment,
			FileName: l.fileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Everything else is illegal
	return token.Token{
		Type:     token.ILLEGAL,
		Literal:  string(l.currentChar),
		FileName: l.fileName,
		Line:     l.currentLine,
		Column:   l.currentColumn,
	}
}

// ReadTokens reads all tokens from the input until the EOF is reached.
func (l *Lexer) ReadTokens() []token.Token {
	var tokens []token.Token
	for {
		nextToken := l.NextToken()
		tokens = append(tokens, nextToken)

		if nextToken.Type == token.EOF {
			break
		}
	}
	return tokens
}
