package lexer

import (
	"strings"

	"github.com/uforg/uforpc/internal/urpc/token"
)

type Lexer struct {
	FileName      string
	CurrentLine   int
	CurrentColumn int

	input             string
	maxIndex          int
	currentIndex      int
	currentIndexIsEOF bool
	currentChar       byte
}

// NewLexer creates a new Lexer and initializes it with the given file
// name and input string.
func NewLexer(fileName, input string) *Lexer {
	l := &Lexer{}

	l.FileName = fileName
	l.CurrentLine = 1
	l.CurrentColumn = 1
	l.input = input
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
		l.CurrentLine++
		l.CurrentColumn = 1
	} else {
		l.CurrentColumn++
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

// readDocstring reads a docstring from the current index to the next triple quote.
//
// Returns the docstring and a boolean indicating if the docstring is unterminated.
func (l *Lexer) readDocstring() (string, bool) {
	if l.currentChar != '"' {
		return "", false
	}

	nextChar, eofReached := l.peekChar(1)
	if eofReached || nextChar != '"' {
		return "", false
	}

	nextChar2, eofReached2 := l.peekChar(2)
	if eofReached2 || nextChar2 != '"' {
		return "", false
	}

	// Skip the opening quotes
	l.readNextChar()
	l.readNextChar()
	l.readNextChar()

	var docstring string
	for {
		isEOF := func() bool {
			if l.currentIndexIsEOF {
				return true
			}

			_, eofReached := l.peekChar(1)
			if eofReached || eofReached2 {
				return true
			}

			_, eofReached2 := l.peekChar(2)
			if eofReached2 || nextChar2 != '"' {
				return true
			}

			return false
		}()
		if isEOF {
			break
		}

		isEndOfDocstring := func() bool {
			if l.currentChar != '"' {
				return false
			}

			nextChar, eofReached := l.peekChar(1)
			if eofReached || nextChar != '"' {
				return false
			}

			nextChar2, eofReached2 := l.peekChar(2)
			if eofReached2 || nextChar2 != '"' {
				return false
			}

			return true
		}()
		if isEndOfDocstring {
			break
		}

		docstring += string(l.currentChar)
		l.readNextChar()
	}

	// Trim beginning and ending space characters
	docstring = strings.TrimSpace(docstring)

	if l.currentIndexIsEOF {
		return docstring, true
	}

	// Skip the 2 remaining closing quotes
	l.readNextChar()
	l.readNextChar()

	return docstring, false
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

	return strings.TrimSpace(comment)
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
			FileName: l.FileName,
			Line:     l.CurrentLine,
			Column:   l.CurrentColumn,
		}
	}

	defer l.readNextChar()
	l.skipWhitespace()

	// Handle delimiters
	if token.IsDelimiter(l.currentChar) {
		return token.Token{
			Type:     token.GetDelimiterTokenType(l.currentChar),
			Literal:  string(l.currentChar),
			FileName: l.FileName,
			Line:     l.CurrentLine,
			Column:   l.CurrentColumn,
		}
	}

	// Handle strings and docstrings
	if l.currentChar == '"' {
		startLine := l.CurrentLine
		startColumn := l.CurrentColumn

		isDocstring := func() bool {
			nextChar, eofReached := l.peekChar(1)
			if eofReached || nextChar != '"' {
				return false
			}

			nextChar2, eofReached2 := l.peekChar(2)
			if eofReached2 || nextChar2 != '"' {
				return false
			}

			return true
		}()

		if isDocstring {
			docstring, unterminated := l.readDocstring()
			if unterminated {
				return token.Token{
					Type:     token.ILLEGAL,
					Literal:  `"""` + docstring,
					FileName: l.FileName,
					Line:     startLine,
					Column:   startColumn,
				}
			}

			return token.Token{
				Type:     token.DOCSTRING,
				Literal:  docstring,
				FileName: l.FileName,
				Line:     startLine,
				Column:   startColumn,
			}
		}

		str, unterminated := l.readString()

		if unterminated {
			return token.Token{
				Type:     token.ILLEGAL,
				Literal:  `"` + str,
				FileName: l.FileName,
				Line:     startLine,
				Column:   startColumn,
			}
		}

		return token.Token{
			Type:     token.STRING,
			Literal:  str,
			FileName: l.FileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Handle ints and floats
	if isNumber(l.currentChar) {
		startLine := l.CurrentLine
		startColumn := l.CurrentColumn
		num := l.readNumber()

		tok := token.Token{
			Type:     token.INT,
			Literal:  num,
			FileName: l.FileName,
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
			FileName: l.FileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Handle identifiers
	if isLetter(l.currentChar) {
		startLine := l.CurrentLine
		startColumn := l.CurrentColumn
		ident := l.readIdentifier()

		if token.IsKeyword(ident) {
			return token.Token{
				Type:     token.GetKeywordTokenType(ident),
				Literal:  ident,
				FileName: l.FileName,
				Line:     startLine,
				Column:   startColumn,
			}
		}

		return token.Token{
			Type:     token.IDENT,
			Literal:  ident,
			FileName: l.FileName,
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
				FileName: l.FileName,
				Line:     l.CurrentLine,
				Column:   l.CurrentColumn,
			}
		}

		startLine := l.CurrentLine
		startColumn := l.CurrentColumn
		comment := l.readComment()
		return token.Token{
			Type:     token.COMMENT,
			Literal:  comment,
			FileName: l.FileName,
			Line:     startLine,
			Column:   startColumn,
		}
	}

	// Everything else is illegal
	return token.Token{
		Type:     token.ILLEGAL,
		Literal:  string(l.currentChar),
		FileName: l.FileName,
		Line:     l.CurrentLine,
		Column:   l.CurrentColumn,
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
