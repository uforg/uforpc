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

	if l.currentChar == '\n' {
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

	switch l.currentChar {
	// Delimiters
	case ',':
		tok = token.NewToken(token.COMMA, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case ':':
		tok = token.NewToken(token.COLON, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '(':
		tok = token.NewToken(token.LPAREN, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case ')':
		tok = token.NewToken(token.RPAREN, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '{':
		tok = token.NewToken(token.LBRACE, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '}':
		tok = token.NewToken(token.RBRACE, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '[':
		tok = token.NewToken(token.LBRACKET, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case ']':
		tok = token.NewToken(token.RBRACKET, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '@':
		tok = token.NewToken(token.AT, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '?':
		tok = token.NewToken(token.QUESTION, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	case '\n':
		tok = token.NewToken(token.NEWLINE, l.currentChar, l.fileName, l.currentLine, l.currentColumn)
	}

	l.readNextChar()
	return tok
}
