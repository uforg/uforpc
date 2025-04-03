package lexer

import (
	"bytes"
	"io"
	"strings"

	plexer "github.com/alecthomas/participle/v2/lexer"
	"github.com/uforg/uforpc/internal/urpc/token"
)

////////////////////////////////////////////////////////////////////////
// This file contains the lexer adapter to make the handwritted lexer //
// compatible with the participle parser.                             //
////////////////////////////////////////////////////////////////////////

// getLexerSymbols returns a map of token types to their corresponding
// participle compatible token types.
func getLexerSymbols() map[string]plexer.TokenType {
	symbols := map[string]plexer.TokenType{}

	for i, tokTypeName := range token.TokenTypes {
		// Should be negative and start from -1 (EOF is -1)
		// https://pkg.go.dev/github.com/alecthomas/participle/v2/lexer#TokenType
		tokenTypeIndex := (i + 1) * -1
		symbols[string(tokTypeName)] = plexer.TokenType(tokenTypeIndex)
	}

	return symbols
}

// lexerSymbols is a map of token types to their corresponding
// participle compatible token types.
var lexerSymbols = getLexerSymbols()

// ParticipleLexer is an adapter that makes the handwritted lexer
// compatible with the participle parser.
//
// Documentation:
//   - https://github.com/alecthomas/participle/blob/master/README.md#lexing
//
// Implements:
//   - https://pkg.go.dev/github.com/alecthomas/participle/v2/lexer#Definition
//   - https://pkg.go.dev/github.com/alecthomas/participle/v2/lexer#StringDefinition
//   - https://pkg.go.dev/github.com/alecthomas/participle/v2/lexer#BytesDefinition
//   - https://pkg.go.dev/github.com/alecthomas/participle/v2/lexer#Lexer
type ParticipleLexer struct {
	customLexer *Lexer
}

func (l *ParticipleLexer) Symbols() map[string]plexer.TokenType {
	return lexerSymbols
}

// Lex creates a new lexer from a reader. Should create a new lexer each
// time it is called to avoid sharing the same state across multiple
// parsings.
func (l *ParticipleLexer) Lex(filename string, r io.Reader) (plexer.Lexer, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lexer := NewLexer(filename, string(content))
	return &ParticipleLexer{customLexer: lexer}, nil
}

func (l *ParticipleLexer) LexString(filename, input string) (plexer.Lexer, error) {
	return l.Lex(filename, strings.NewReader(input))
}

func (l *ParticipleLexer) LexBytes(filename string, content []byte) (plexer.Lexer, error) {
	return l.Lex(filename, bytes.NewReader(content))
}

func (l *ParticipleLexer) Next() (plexer.Token, error) {
	tok := l.customLexer.NextToken()

	tokType := lexerSymbols[string(tok.Type)]
	tokValue := tok.Literal

	tokPos := plexer.Position{
		Filename: l.customLexer.FileName,
		Offset:   max(tok.ColumnStart-1, 0), // Is the offset of the token from the start of the line
		Line:     tok.LineStart,
		Column:   tok.ColumnStart,
	}

	return plexer.Token{
		Type:  tokType,
		Value: tokValue,
		Pos:   tokPos,
	}, nil
}
