package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/token"
)

func TestLexer(t *testing.T) {
	t.Run("TestLexerBasic", func(t *testing.T) {
		input := ",:(){}[]@?"

		tests := []token.Token{
			{Type: token.COMMA, Literal: ",", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", LineStart: 1, ColumnStart: 2, LineEnd: 1, ColumnEnd: 2},
			{Type: token.LPAREN, Literal: "(", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.RPAREN, Literal: ")", FileName: "test.urpc", LineStart: 1, ColumnStart: 4, LineEnd: 1, ColumnEnd: 4},
			{Type: token.LBRACE, Literal: "{", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.RBRACE, Literal: "}", FileName: "test.urpc", LineStart: 1, ColumnStart: 6, LineEnd: 1, ColumnEnd: 6},
			{Type: token.LBRACKET, Literal: "[", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
			{Type: token.RBRACKET, Literal: "]", FileName: "test.urpc", LineStart: 1, ColumnStart: 8, LineEnd: 1, ColumnEnd: 8},
			{Type: token.AT, Literal: "@", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 9},
			{Type: token.QUESTION, Literal: "?", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 10},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 11, LineEnd: 1, ColumnEnd: 11},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerNewLines", func(t *testing.T) {
		input := ",:\n(){\n}\n[]@?\n"

		tests := []token.Token{
			{Type: token.COMMA, Literal: ",", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", LineStart: 1, ColumnStart: 2, LineEnd: 1, ColumnEnd: 2},
			{Type: token.LPAREN, Literal: "(", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 1},
			{Type: token.RPAREN, Literal: ")", FileName: "test.urpc", LineStart: 2, ColumnStart: 2, LineEnd: 2, ColumnEnd: 2},
			{Type: token.LBRACE, Literal: "{", FileName: "test.urpc", LineStart: 2, ColumnStart: 3, LineEnd: 2, ColumnEnd: 3},
			{Type: token.RBRACE, Literal: "}", FileName: "test.urpc", LineStart: 3, ColumnStart: 1, LineEnd: 3, ColumnEnd: 1},
			{Type: token.LBRACKET, Literal: "[", FileName: "test.urpc", LineStart: 4, ColumnStart: 1, LineEnd: 4, ColumnEnd: 1},
			{Type: token.RBRACKET, Literal: "]", FileName: "test.urpc", LineStart: 4, ColumnStart: 2, LineEnd: 4, ColumnEnd: 2},
			{Type: token.AT, Literal: "@", FileName: "test.urpc", LineStart: 4, ColumnStart: 3, LineEnd: 4, ColumnEnd: 3},
			{Type: token.QUESTION, Literal: "?", FileName: "test.urpc", LineStart: 4, ColumnStart: 4, LineEnd: 4, ColumnEnd: 4},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 5, ColumnStart: 1, LineEnd: 5, ColumnEnd: 1},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerKeywords", func(t *testing.T) {
		input := "version rule type proc input output meta error true false for param extends"

		tests := []token.Token{
			{Type: token.VERSION, Literal: "version", FileName: "test.urpc", LineStart: 1, LineEnd: 1, ColumnStart: 1, ColumnEnd: 7},
			{Type: token.RULE, Literal: "rule", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 12},
			{Type: token.TYPE, Literal: "type", FileName: "test.urpc", LineStart: 1, ColumnStart: 14, LineEnd: 1, ColumnEnd: 17},
			{Type: token.PROC, Literal: "proc", FileName: "test.urpc", LineStart: 1, ColumnStart: 19, LineEnd: 1, ColumnEnd: 22},
			{Type: token.INPUT, Literal: "input", FileName: "test.urpc", LineStart: 1, ColumnStart: 24, LineEnd: 1, ColumnEnd: 28},
			{Type: token.OUTPUT, Literal: "output", FileName: "test.urpc", LineStart: 1, ColumnStart: 30, LineEnd: 1, ColumnEnd: 35},
			{Type: token.META, Literal: "meta", FileName: "test.urpc", LineStart: 1, ColumnStart: 37, LineEnd: 1, ColumnEnd: 40},
			{Type: token.ERROR, Literal: "error", FileName: "test.urpc", LineStart: 1, ColumnStart: 42, LineEnd: 1, ColumnEnd: 46},
			{Type: token.TRUE, Literal: "true", FileName: "test.urpc", LineStart: 1, ColumnStart: 48, LineEnd: 1, ColumnEnd: 51},
			{Type: token.FALSE, Literal: "false", FileName: "test.urpc", LineStart: 1, ColumnStart: 53, LineEnd: 1, ColumnEnd: 57},
			{Type: token.FOR, Literal: "for", FileName: "test.urpc", LineStart: 1, ColumnStart: 59, LineEnd: 1, ColumnEnd: 61},
			{Type: token.PARAM, Literal: "param", FileName: "test.urpc", LineStart: 1, ColumnStart: 63, LineEnd: 1, ColumnEnd: 67},
			{Type: token.EXTENDS, Literal: "extends", FileName: "test.urpc", LineStart: 1, ColumnStart: 69, LineEnd: 1, ColumnEnd: 75},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 76, LineEnd: 1, ColumnEnd: 76},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerIdentifiers", func(t *testing.T) {
		input := "hello world someIdentifier"

		tests := []token.Token{
			{Type: token.IDENT, Literal: "hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 5},
			{Type: token.IDENT, Literal: "world", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 11},
			{Type: token.IDENT, Literal: "someIdentifier", FileName: "test.urpc", LineStart: 1, ColumnStart: 13, LineEnd: 1, ColumnEnd: 26},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 27, LineEnd: 1, ColumnEnd: 27},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerIdentifiersWithNumbers", func(t *testing.T) {
		input := "hello123 world456 someIdentifier789"

		tests := []token.Token{
			{Type: token.IDENT, Literal: "hello123", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 8},
			{Type: token.IDENT, Literal: "world456", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 17},
			{Type: token.IDENT, Literal: "someIdentifier789", FileName: "test.urpc", LineStart: 1, ColumnStart: 19, LineEnd: 1, ColumnEnd: 35},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 36, LineEnd: 1, ColumnEnd: 36},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}
	})

	t.Run("TestLexerNumbers", func(t *testing.T) {
		input := "1 2 3 456 789"

		tests := []token.Token{
			{Type: token.INT, Literal: "1", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.INT, Literal: "2", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.INT, Literal: "3", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.INT, Literal: "456", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 9},
			{Type: token.INT, Literal: "789", FileName: "test.urpc", LineStart: 1, ColumnStart: 11, LineEnd: 1, ColumnEnd: 13},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 14, LineEnd: 1, ColumnEnd: 14},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerFloats", func(t *testing.T) {
		input := "1.2 3.45 67.89 1.2.3.4"

		tests := []token.Token{
			{Type: token.FLOAT, Literal: "1.2", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 3},
			{Type: token.FLOAT, Literal: "3.45", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 8},
			{Type: token.FLOAT, Literal: "67.89", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 14},
			{Type: token.FLOAT, Literal: "1.2", FileName: "test.urpc", LineStart: 1, ColumnStart: 16, LineEnd: 1, ColumnEnd: 18},
			{Type: token.ILLEGAL, Literal: ".", FileName: "test.urpc", LineStart: 1, ColumnStart: 19, LineEnd: 1, ColumnEnd: 19},
			{Type: token.FLOAT, Literal: "3.4", FileName: "test.urpc", LineStart: 1, ColumnStart: 20, LineEnd: 1, ColumnEnd: 22},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 23, LineEnd: 1, ColumnEnd: 23},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerStrings", func(t *testing.T) {
		input := `"hello" "world" "hello world!"test`

		tests := []token.Token{
			{Type: token.STRING, Literal: "hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 7},
			{Type: token.STRING, Literal: "world", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 15},
			{Type: token.STRING, Literal: "hello world!", FileName: "test.urpc", LineStart: 1, ColumnStart: 17, LineEnd: 1, ColumnEnd: 30},
			{Type: token.IDENT, Literal: "test", FileName: "test.urpc", LineStart: 1, ColumnStart: 31, LineEnd: 1, ColumnEnd: 34},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 35, LineEnd: 1, ColumnEnd: 35},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerIllegal", func(t *testing.T) {
		input := "$ % ^ & ."

		tests := []token.Token{
			{Type: token.ILLEGAL, Literal: "$", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.ILLEGAL, Literal: "%", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.ILLEGAL, Literal: "^", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.ILLEGAL, Literal: "&", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
			{Type: token.ILLEGAL, Literal: ".", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 9},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 10},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerComments", func(t *testing.T) {
		input := "// This is a comment\nversion: 1"

		tests := []token.Token{
			{Type: token.COMMENT, Literal: "This is a comment", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 20},
			{Type: token.VERSION, Literal: "version", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 7},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", LineStart: 2, ColumnStart: 8, LineEnd: 2, ColumnEnd: 8},
			{Type: token.INT, Literal: "1", FileName: "test.urpc", LineStart: 2, ColumnStart: 10, LineEnd: 2, ColumnEnd: 10},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 2, ColumnStart: 11, LineEnd: 2, ColumnEnd: 11},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestMultilineComments", func(t *testing.T) {
		input := "/* This is a multiline comment\nwith multiple lines */"

		tests := []token.Token{
			{Type: token.COMMENT, Literal: "This is a multiline comment\nwith multiple lines", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 2, ColumnEnd: 22},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 2, ColumnStart: 23, LineEnd: 2, ColumnEnd: 23},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerUnterminatedString", func(t *testing.T) {
		input := `"hello`

		tests := []token.Token{
			{Type: token.ILLEGAL, Literal: "\"hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 6},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerDocstrings", func(t *testing.T) {
		input := `""" This is a docstring """`

		tests := []token.Token{
			{Type: token.DOCSTRING, Literal: "This is a docstring", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 27},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 28, LineEnd: 1, ColumnEnd: 28},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerMultilineDocstrings", func(t *testing.T) {
		input := "\"\"\" This is a multiline docstring\nwith multiple lines \"\"\"\n"

		tests := []token.Token{
			{Type: token.DOCSTRING, Literal: "This is a multiline docstring\nwith multiple lines", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 2, ColumnEnd: 23},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", LineStart: 3, ColumnStart: 1, LineEnd: 3, ColumnEnd: 1},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.LineStart, tok.LineStart, "test %d", i)
			require.Equal(t, test.ColumnStart, tok.ColumnStart, "test %d", i)
			require.Equal(t, test.LineEnd, tok.LineEnd, "test %d", i)
			require.Equal(t, test.ColumnEnd, tok.ColumnEnd, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerURPC", func(t *testing.T) {
		input := `
			// This test evaluates the lexer with a full URPC file.
			// It is used to ensure that the lexer is working correctly.

			version: 1

			""" Product is a type that represents a product. """
			type Product extends OtherType, AnotherType {
				id: string
					@uuid
					@minLen(36)
				
				name: string
					@minLen(3)
					@maxLen(100)
				
				price: float
					@min(0.01)
				
				tags?: string[]
					@maxItems(5)
			}

			""" Creates a product and returns the product id. """
			proc CreateProduct {
				input {
					product: Product
					priority: int
						@enum([1, 2, 3], error: "Priority must be 1, 2, or 3")
				}
				
				output {
					success: boolean
					productId: string
						@uuid
				}
				
				meta {
					audit: true
					maxRetries: 3
				}
			}`

		tests := []token.Token{
			{Type: token.COMMENT, Literal: "This test evaluates the lexer with a full URPC file."},
			{Type: token.COMMENT, Literal: "It is used to ensure that the lexer is working correctly."},
			{Type: token.VERSION, Literal: "version"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.INT, Literal: "1"},
			{Type: token.DOCSTRING, Literal: "Product is a type that represents a product."},
			{Type: token.TYPE, Literal: "type"},
			{Type: token.IDENT, Literal: "Product"},
			{Type: token.EXTENDS, Literal: "extends"},
			{Type: token.IDENT, Literal: "OtherType"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.IDENT, Literal: "AnotherType"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.IDENT, Literal: "id"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "string"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "uuid"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "minLen"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.INT, Literal: "36"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.IDENT, Literal: "name"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "string"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "minLen"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.INT, Literal: "3"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "maxLen"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.INT, Literal: "100"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.IDENT, Literal: "price"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "float"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "min"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.FLOAT, Literal: "0.01"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.IDENT, Literal: "tags"},
			{Type: token.QUESTION, Literal: "?"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "string"},
			{Type: token.LBRACKET, Literal: "["},
			{Type: token.RBRACKET, Literal: "]"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "maxItems"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.INT, Literal: "5"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.DOCSTRING, Literal: "Creates a product and returns the product id."},
			{Type: token.PROC, Literal: "proc"},
			{Type: token.IDENT, Literal: "CreateProduct"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.INPUT, Literal: "input"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.IDENT, Literal: "product"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "Product"},
			{Type: token.IDENT, Literal: "priority"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "int"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "enum"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.LBRACKET, Literal: "["},
			{Type: token.INT, Literal: "1"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.INT, Literal: "2"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.INT, Literal: "3"},
			{Type: token.RBRACKET, Literal: "]"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.ERROR, Literal: "error"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.STRING, Literal: "Priority must be 1, 2, or 3"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.OUTPUT, Literal: "output"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.IDENT, Literal: "success"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "boolean"},
			{Type: token.IDENT, Literal: "productId"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.IDENT, Literal: "string"},
			{Type: token.AT, Literal: "@"},
			{Type: token.IDENT, Literal: "uuid"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.META, Literal: "meta"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.IDENT, Literal: "audit"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.TRUE, Literal: "true"},
			{Type: token.IDENT, Literal: "maxRetries"},
			{Type: token.COLON, Literal: ":"},
			{Type: token.INT, Literal: "3"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.EOF, Literal: ""},
		}

		lex := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})
}
