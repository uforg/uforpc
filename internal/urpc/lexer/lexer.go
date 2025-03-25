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
	nextIndex         int
	nextIndexIsEOF    bool
	nextChar          byte
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

	l.nextIndex = 1
	if l.maxIndex <= 1 {
		l.nextIndexIsEOF = true
	} else {
		l.nextIndexIsEOF = false
	}
	if l.nextIndexIsEOF {
		l.nextChar = 0
	} else {
		l.nextChar = input[l.nextIndex]
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

	l.nextIndex++
	if l.nextIndex > l.maxIndex {
		l.nextIndexIsEOF = true
		l.nextChar = 0
	} else {
		l.nextChar = l.input[l.nextIndex]
	}
}

// readIdentifier reads an identifier from the current index to the next non-letter character.
func (l *Lexer) readIdentifier() string {
	var ident string
	for isLetter(l.currentChar) {
		ident += string(l.currentChar)
		l.readNextChar()
	}
	return ident
}

// skipWhitespace skips whitespace characters from the current index to the next non-whitespace character.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.currentChar) {
		l.readNextChar()
	}
}

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

	var tok token.Token
	l.skipWhitespace()

	// Handle delimiters
	switch l.currentChar {
	case ',':
		tok = token.NewToken(token.COMMA, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case ':':
		tok = token.NewToken(token.COLON, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '(':
		tok = token.NewToken(token.LPAREN, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case ')':
		tok = token.NewToken(token.RPAREN, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '{':
		tok = token.NewToken(token.LBRACE, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '}':
		tok = token.NewToken(token.RBRACE, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '[':
		tok = token.NewToken(token.LBRACKET, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case ']':
		tok = token.NewToken(token.RBRACKET, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '@':
		tok = token.NewToken(token.AT, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '?':
		tok = token.NewToken(token.QUESTION, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	case '\n':
		tok = token.NewToken(token.NEWLINE, string(l.currentChar), l.fileName, l.currentLine, l.currentColumn)
	}

	// Handle identifiers
	if isLetter(l.currentChar) {
		tokenType := token.IDENT
		line := l.currentLine
		column := l.currentColumn
		ident := l.readIdentifier()

		if token.IsKeyword(ident) {
			tokenType = token.GetKeywordTokenType(ident)
		}

		tok = token.Token{
			Type:     tokenType,
			Literal:  ident,
			FileName: l.fileName,
			Line:     line,
			Column:   column,
		}
	}

	l.readNextChar()
	return tok
}
