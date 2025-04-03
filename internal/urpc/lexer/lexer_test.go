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
			{Type: token.Comma, Literal: ",", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.Colon, Literal: ":", FileName: "test.urpc", LineStart: 1, ColumnStart: 2, LineEnd: 1, ColumnEnd: 2},
			{Type: token.LParen, Literal: "(", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.RParen, Literal: ")", FileName: "test.urpc", LineStart: 1, ColumnStart: 4, LineEnd: 1, ColumnEnd: 4},
			{Type: token.LBrace, Literal: "{", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.RBrace, Literal: "}", FileName: "test.urpc", LineStart: 1, ColumnStart: 6, LineEnd: 1, ColumnEnd: 6},
			{Type: token.LBracket, Literal: "[", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
			{Type: token.RBracket, Literal: "]", FileName: "test.urpc", LineStart: 1, ColumnStart: 8, LineEnd: 1, ColumnEnd: 8},
			{Type: token.At, Literal: "@", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 9},
			{Type: token.Question, Literal: "?", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 10},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 11, LineEnd: 1, ColumnEnd: 11},
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
			{Type: token.Comma, Literal: ",", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.Colon, Literal: ":", FileName: "test.urpc", LineStart: 1, ColumnStart: 2, LineEnd: 1, ColumnEnd: 2},
			{Type: token.LParen, Literal: "(", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 1},
			{Type: token.RParen, Literal: ")", FileName: "test.urpc", LineStart: 2, ColumnStart: 2, LineEnd: 2, ColumnEnd: 2},
			{Type: token.LBrace, Literal: "{", FileName: "test.urpc", LineStart: 2, ColumnStart: 3, LineEnd: 2, ColumnEnd: 3},
			{Type: token.RBrace, Literal: "}", FileName: "test.urpc", LineStart: 3, ColumnStart: 1, LineEnd: 3, ColumnEnd: 1},
			{Type: token.LBracket, Literal: "[", FileName: "test.urpc", LineStart: 4, ColumnStart: 1, LineEnd: 4, ColumnEnd: 1},
			{Type: token.RBracket, Literal: "]", FileName: "test.urpc", LineStart: 4, ColumnStart: 2, LineEnd: 4, ColumnEnd: 2},
			{Type: token.At, Literal: "@", FileName: "test.urpc", LineStart: 4, ColumnStart: 3, LineEnd: 4, ColumnEnd: 3},
			{Type: token.Question, Literal: "?", FileName: "test.urpc", LineStart: 4, ColumnStart: 4, LineEnd: 4, ColumnEnd: 4},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 5, ColumnStart: 1, LineEnd: 5, ColumnEnd: 1},
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
		input := "version rule type proc input output meta error true false for param extends string int float boolean import datetime"

		tests := []token.Token{
			{Type: token.Version, Literal: "version", FileName: "test.urpc", LineStart: 1, LineEnd: 1, ColumnStart: 1, ColumnEnd: 7},
			{Type: token.Rule, Literal: "rule", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 12},
			{Type: token.Type, Literal: "type", FileName: "test.urpc", LineStart: 1, ColumnStart: 14, LineEnd: 1, ColumnEnd: 17},
			{Type: token.Proc, Literal: "proc", FileName: "test.urpc", LineStart: 1, ColumnStart: 19, LineEnd: 1, ColumnEnd: 22},
			{Type: token.Input, Literal: "input", FileName: "test.urpc", LineStart: 1, ColumnStart: 24, LineEnd: 1, ColumnEnd: 28},
			{Type: token.Output, Literal: "output", FileName: "test.urpc", LineStart: 1, ColumnStart: 30, LineEnd: 1, ColumnEnd: 35},
			{Type: token.Meta, Literal: "meta", FileName: "test.urpc", LineStart: 1, ColumnStart: 37, LineEnd: 1, ColumnEnd: 40},
			{Type: token.Error, Literal: "error", FileName: "test.urpc", LineStart: 1, ColumnStart: 42, LineEnd: 1, ColumnEnd: 46},
			{Type: token.TrueLiteral, Literal: "true", FileName: "test.urpc", LineStart: 1, ColumnStart: 48, LineEnd: 1, ColumnEnd: 51},
			{Type: token.FalseLiteral, Literal: "false", FileName: "test.urpc", LineStart: 1, ColumnStart: 53, LineEnd: 1, ColumnEnd: 57},
			{Type: token.For, Literal: "for", FileName: "test.urpc", LineStart: 1, ColumnStart: 59, LineEnd: 1, ColumnEnd: 61},
			{Type: token.Param, Literal: "param", FileName: "test.urpc", LineStart: 1, ColumnStart: 63, LineEnd: 1, ColumnEnd: 67},
			{Type: token.Extends, Literal: "extends", FileName: "test.urpc", LineStart: 1, ColumnStart: 69, LineEnd: 1, ColumnEnd: 75},
			{Type: token.String, Literal: "string", FileName: "test.urpc", LineStart: 1, ColumnStart: 77, LineEnd: 1, ColumnEnd: 82},
			{Type: token.Int, Literal: "int", FileName: "test.urpc", LineStart: 1, ColumnStart: 84, LineEnd: 1, ColumnEnd: 86},
			{Type: token.Float, Literal: "float", FileName: "test.urpc", LineStart: 1, ColumnStart: 88, LineEnd: 1, ColumnEnd: 92},
			{Type: token.Boolean, Literal: "boolean", FileName: "test.urpc", LineStart: 1, ColumnStart: 94, LineEnd: 1, ColumnEnd: 100},
			{Type: token.Import, Literal: "import", FileName: "test.urpc", LineStart: 1, ColumnStart: 102, LineEnd: 1, ColumnEnd: 107},
			{Type: token.Datetime, Literal: "datetime", FileName: "test.urpc", LineStart: 1, ColumnStart: 109, LineEnd: 1, ColumnEnd: 116},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 117, LineEnd: 1, ColumnEnd: 117},
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
			{Type: token.Ident, Literal: "hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 5},
			{Type: token.Ident, Literal: "world", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 11},
			{Type: token.Ident, Literal: "someIdentifier", FileName: "test.urpc", LineStart: 1, ColumnStart: 13, LineEnd: 1, ColumnEnd: 26},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 27, LineEnd: 1, ColumnEnd: 27},
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
			{Type: token.Ident, Literal: "hello123", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 8},
			{Type: token.Ident, Literal: "world456", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 17},
			{Type: token.Ident, Literal: "someIdentifier789", FileName: "test.urpc", LineStart: 1, ColumnStart: 19, LineEnd: 1, ColumnEnd: 35},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 36, LineEnd: 1, ColumnEnd: 36},
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
			{Type: token.IntLiteral, Literal: "1", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.IntLiteral, Literal: "2", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.IntLiteral, Literal: "3", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.IntLiteral, Literal: "456", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 9},
			{Type: token.IntLiteral, Literal: "789", FileName: "test.urpc", LineStart: 1, ColumnStart: 11, LineEnd: 1, ColumnEnd: 13},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 14, LineEnd: 1, ColumnEnd: 14},
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
			{Type: token.FloatLiteral, Literal: "1.2", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 3},
			{Type: token.FloatLiteral, Literal: "3.45", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 8},
			{Type: token.FloatLiteral, Literal: "67.89", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 14},
			{Type: token.FloatLiteral, Literal: "1.2", FileName: "test.urpc", LineStart: 1, ColumnStart: 16, LineEnd: 1, ColumnEnd: 18},
			{Type: token.Illegal, Literal: ".", FileName: "test.urpc", LineStart: 1, ColumnStart: 19, LineEnd: 1, ColumnEnd: 19},
			{Type: token.FloatLiteral, Literal: "3.4", FileName: "test.urpc", LineStart: 1, ColumnStart: 20, LineEnd: 1, ColumnEnd: 22},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 23, LineEnd: 1, ColumnEnd: 23},
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
		input := `"hello" "world" "hello world!"test "hello \"quotes\" \\"`

		tests := []token.Token{
			{Type: token.StringLiteral, Literal: "hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 7},
			{Type: token.StringLiteral, Literal: "world", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 15},
			{Type: token.StringLiteral, Literal: "hello world!", FileName: "test.urpc", LineStart: 1, ColumnStart: 17, LineEnd: 1, ColumnEnd: 30},
			{Type: token.Ident, Literal: "test", FileName: "test.urpc", LineStart: 1, ColumnStart: 31, LineEnd: 1, ColumnEnd: 34},
			{Type: token.StringLiteral, Literal: "hello \"quotes\" \\", FileName: "test.urpc", LineStart: 1, ColumnStart: 36, LineEnd: 1, ColumnEnd: 56},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 57, LineEnd: 1, ColumnEnd: 57},
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
			{Type: token.Illegal, Literal: "$", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 1},
			{Type: token.Illegal, Literal: "%", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.Illegal, Literal: "^", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.Illegal, Literal: "&", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
			{Type: token.Illegal, Literal: ".", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 9},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 10},
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
			{Type: token.Comment, Literal: "This is a comment", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 20},
			{Type: token.Version, Literal: "version", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 7},
			{Type: token.Colon, Literal: ":", FileName: "test.urpc", LineStart: 2, ColumnStart: 8, LineEnd: 2, ColumnEnd: 8},
			{Type: token.IntLiteral, Literal: "1", FileName: "test.urpc", LineStart: 2, ColumnStart: 10, LineEnd: 2, ColumnEnd: 10},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 2, ColumnStart: 11, LineEnd: 2, ColumnEnd: 11},
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

	t.Run("TestLexerCommentBlocks", func(t *testing.T) {
		input := "/* This is a multiline comment\nwith multiple lines */"

		tests := []token.Token{
			{Type: token.CommentBlock, Literal: "This is a multiline comment\nwith multiple lines", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 2, ColumnEnd: 22},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 2, ColumnStart: 23, LineEnd: 2, ColumnEnd: 23},
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
			{Type: token.Illegal, Literal: "\"hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 6},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
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
			{Type: token.Docstring, Literal: "This is a docstring", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 27},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 28, LineEnd: 1, ColumnEnd: 28},
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

	t.Run("TestLexerDocstringsWithQuotes", func(t *testing.T) {
		input := `""" This is a docstring with "quotes" inside """`

		tests := []token.Token{
			{Type: token.Docstring, Literal: "This is a docstring with \"quotes\" inside", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 48},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 1, ColumnStart: 49, LineEnd: 1, ColumnEnd: 49},
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
			{Type: token.Docstring, Literal: "This is a multiline docstring\nwith multiple lines", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 2, ColumnEnd: 23},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 3, ColumnStart: 1, LineEnd: 3, ColumnEnd: 1},
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
			{Type: token.Comment, Literal: "This test evaluates the lexer with a full URPC file."},
			{Type: token.Comment, Literal: "It is used to ensure that the lexer is working correctly."},
			{Type: token.Version, Literal: "version"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.IntLiteral, Literal: "1"},
			{Type: token.Docstring, Literal: "Product is a type that represents a product."},
			{Type: token.Type, Literal: "type"},
			{Type: token.Ident, Literal: "Product"},
			{Type: token.Extends, Literal: "extends"},
			{Type: token.Ident, Literal: "OtherType"},
			{Type: token.Comma, Literal: ","},
			{Type: token.Ident, Literal: "AnotherType"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "id"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.String, Literal: "string"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "uuid"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "minLen"},
			{Type: token.LParen, Literal: "("},
			{Type: token.IntLiteral, Literal: "36"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.Ident, Literal: "name"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.String, Literal: "string"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "minLen"},
			{Type: token.LParen, Literal: "("},
			{Type: token.IntLiteral, Literal: "3"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "maxLen"},
			{Type: token.LParen, Literal: "("},
			{Type: token.IntLiteral, Literal: "100"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.Ident, Literal: "price"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.Float, Literal: "float"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "min"},
			{Type: token.LParen, Literal: "("},
			{Type: token.FloatLiteral, Literal: "0.01"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.Ident, Literal: "tags"},
			{Type: token.Question, Literal: "?"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.String, Literal: "string"},
			{Type: token.LBracket, Literal: "["},
			{Type: token.RBracket, Literal: "]"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "maxItems"},
			{Type: token.LParen, Literal: "("},
			{Type: token.IntLiteral, Literal: "5"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Docstring, Literal: "Creates a product and returns the product id."},
			{Type: token.Proc, Literal: "proc"},
			{Type: token.Ident, Literal: "CreateProduct"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Input, Literal: "input"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "product"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.Ident, Literal: "Product"},
			{Type: token.Ident, Literal: "priority"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.Int, Literal: "int"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "enum"},
			{Type: token.LParen, Literal: "("},
			{Type: token.LBracket, Literal: "["},
			{Type: token.IntLiteral, Literal: "1"},
			{Type: token.Comma, Literal: ","},
			{Type: token.IntLiteral, Literal: "2"},
			{Type: token.Comma, Literal: ","},
			{Type: token.IntLiteral, Literal: "3"},
			{Type: token.RBracket, Literal: "]"},
			{Type: token.Comma, Literal: ","},
			{Type: token.Error, Literal: "error"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.StringLiteral, Literal: "Priority must be 1, 2, or 3"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Output, Literal: "output"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "success"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.Boolean, Literal: "boolean"},
			{Type: token.Ident, Literal: "productId"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.String, Literal: "string"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "uuid"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Meta, Literal: "meta"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "audit"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.TrueLiteral, Literal: "true"},
			{Type: token.Ident, Literal: "maxRetries"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.IntLiteral, Literal: "3"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Eof, Literal: ""},
		}

		lex := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})
}
