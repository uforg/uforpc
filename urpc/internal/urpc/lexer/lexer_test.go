package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/urpc/internal/urpc/token"
)

func TestLexer(t *testing.T) {
	// TODO: Add more tests specifically for the token positions

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
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.LParen, Literal: "(", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 1},
			{Type: token.RParen, Literal: ")", FileName: "test.urpc", LineStart: 2, ColumnStart: 2, LineEnd: 2, ColumnEnd: 2},
			{Type: token.LBrace, Literal: "{", FileName: "test.urpc", LineStart: 2, ColumnStart: 3, LineEnd: 2, ColumnEnd: 3},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 2, ColumnStart: 4, LineEnd: 2, ColumnEnd: 4},
			{Type: token.RBrace, Literal: "}", FileName: "test.urpc", LineStart: 3, ColumnStart: 1, LineEnd: 3, ColumnEnd: 1},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 3, ColumnStart: 2, LineEnd: 3, ColumnEnd: 2},
			{Type: token.LBracket, Literal: "[", FileName: "test.urpc", LineStart: 4, ColumnStart: 1, LineEnd: 4, ColumnEnd: 1},
			{Type: token.RBracket, Literal: "]", FileName: "test.urpc", LineStart: 4, ColumnStart: 2, LineEnd: 4, ColumnEnd: 2},
			{Type: token.At, Literal: "@", FileName: "test.urpc", LineStart: 4, ColumnStart: 3, LineEnd: 4, ColumnEnd: 3},
			{Type: token.Question, Literal: "?", FileName: "test.urpc", LineStart: 4, ColumnStart: 4, LineEnd: 4, ColumnEnd: 4},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 4, ColumnStart: 5, LineEnd: 4, ColumnEnd: 5},
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
		input := "version type proc input output true false string int float bool datetime deprecated stream"

		tests := []token.Token{
			{Type: token.Version, Literal: "version"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Type, Literal: "type"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Proc, Literal: "proc"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Input, Literal: "input"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Output, Literal: "output"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.TrueLiteral, Literal: "true"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.FalseLiteral, Literal: "false"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.String, Literal: "string"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Int, Literal: "int"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Float, Literal: "float"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Bool, Literal: "bool"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Datetime, Literal: "datetime"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Deprecated, Literal: "deprecated"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.Stream, Literal: "stream"},
			{Type: token.Eof, Literal: ""},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})

	t.Run("TestLexerIdentifiers", func(t *testing.T) {
		input := "hello world someIdentifier"

		tests := []token.Token{
			{Type: token.Ident, Literal: "hello", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 5},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 6, LineEnd: 1, ColumnEnd: 6},
			{Type: token.Ident, Literal: "world", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 11},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 12, LineEnd: 1, ColumnEnd: 12},
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
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 9},
			{Type: token.Ident, Literal: "world456", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 17},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 18, LineEnd: 1, ColumnEnd: 18},
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
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 2, LineEnd: 1, ColumnEnd: 2},
			{Type: token.IntLiteral, Literal: "2", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 4, LineEnd: 1, ColumnEnd: 4},
			{Type: token.IntLiteral, Literal: "3", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 6, LineEnd: 1, ColumnEnd: 6},
			{Type: token.IntLiteral, Literal: "456", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 9},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 10},
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
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 4, LineEnd: 1, ColumnEnd: 4},
			{Type: token.FloatLiteral, Literal: "3.45", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 8},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 9},
			{Type: token.FloatLiteral, Literal: "67.89", FileName: "test.urpc", LineStart: 1, ColumnStart: 10, LineEnd: 1, ColumnEnd: 14},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 15, LineEnd: 1, ColumnEnd: 15},
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
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 8, LineEnd: 1, ColumnEnd: 8},
			{Type: token.StringLiteral, Literal: "world", FileName: "test.urpc", LineStart: 1, ColumnStart: 9, LineEnd: 1, ColumnEnd: 15},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 16, LineEnd: 1, ColumnEnd: 16},
			{Type: token.StringLiteral, Literal: "hello world!", FileName: "test.urpc", LineStart: 1, ColumnStart: 17, LineEnd: 1, ColumnEnd: 30},
			{Type: token.Ident, Literal: "test", FileName: "test.urpc", LineStart: 1, ColumnStart: 31, LineEnd: 1, ColumnEnd: 34},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 35, LineEnd: 1, ColumnEnd: 35},
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
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 2, LineEnd: 1, ColumnEnd: 2},
			{Type: token.Illegal, Literal: "%", FileName: "test.urpc", LineStart: 1, ColumnStart: 3, LineEnd: 1, ColumnEnd: 3},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 4, LineEnd: 1, ColumnEnd: 4},
			{Type: token.Illegal, Literal: "^", FileName: "test.urpc", LineStart: 1, ColumnStart: 5, LineEnd: 1, ColumnEnd: 5},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 6, LineEnd: 1, ColumnEnd: 6},
			{Type: token.Illegal, Literal: "&", FileName: "test.urpc", LineStart: 1, ColumnStart: 7, LineEnd: 1, ColumnEnd: 7},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 1, ColumnStart: 8, LineEnd: 1, ColumnEnd: 8},
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
			{Type: token.Comment, Literal: " This is a comment", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 20},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 1, ColumnStart: 21, LineEnd: 1, ColumnEnd: 21},
			{Type: token.Version, Literal: "version", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 7},
			{Type: token.Colon, Literal: ":", FileName: "test.urpc", LineStart: 2, ColumnStart: 8, LineEnd: 2, ColumnEnd: 8},
			{Type: token.Whitespace, Literal: " ", FileName: "test.urpc", LineStart: 2, ColumnStart: 9, LineEnd: 2, ColumnEnd: 9},
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

	t.Run("TestLexerEmptyComments", func(t *testing.T) {
		input := "//\nversion 1/**/version 1"

		tests := []token.Token{
			{Type: token.Comment, Literal: ""},
			{Type: token.Newline, Literal: "\n"},
			{Type: token.Version, Literal: "version"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.IntLiteral, Literal: "1"},
			{Type: token.CommentBlock, Literal: ""},
			{Type: token.Version, Literal: "version"},
			{Type: token.Whitespace, Literal: " "},
			{Type: token.IntLiteral, Literal: "1"},
		}

		lex := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})

	t.Run("TestLexerCommentsPositions", func(t *testing.T) {
		input := "// This is a comment\n\n/* This is a comment */\n\n// This is a comment"

		tests := []token.Token{
			{Type: token.Comment, Literal: " This is a comment", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 20},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 1, ColumnStart: 21, LineEnd: 1, ColumnEnd: 21},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 2, ColumnStart: 1, LineEnd: 2, ColumnEnd: 1},
			{Type: token.CommentBlock, Literal: " This is a comment ", FileName: "test.urpc", LineStart: 3, ColumnStart: 1, LineEnd: 3, ColumnEnd: 23},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 3, ColumnStart: 24, LineEnd: 3, ColumnEnd: 24},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 4, ColumnStart: 1, LineEnd: 4, ColumnEnd: 1},
			{Type: token.Comment, Literal: " This is a comment", FileName: "test.urpc", LineStart: 5, ColumnStart: 1, LineEnd: 5, ColumnEnd: 20},
			{Type: token.Eof, Literal: "", FileName: "test.urpc", LineStart: 5, ColumnStart: 21, LineEnd: 5, ColumnEnd: 21},
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
			{Type: token.CommentBlock, Literal: " This is a multiline comment\nwith multiple lines ", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 2, ColumnEnd: 22},
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
			{Type: token.Docstring, Literal: " This is a docstring ", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 27},
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
			{Type: token.Docstring, Literal: " This is a docstring with \"quotes\" inside ", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 1, ColumnEnd: 48},
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
			{Type: token.Docstring, Literal: " This is a multiline docstring\nwith multiple lines ", FileName: "test.urpc", LineStart: 1, ColumnStart: 1, LineEnd: 2, ColumnEnd: 23},
			{Type: token.Newline, Literal: "\n", FileName: "test.urpc", LineStart: 2, ColumnStart: 24, LineEnd: 2, ColumnEnd: 24},
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

	t.Run("TestLexerUnicode", func(t *testing.T) {
		input := `
			// áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍

			/* áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍 */

			""" áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍 """

			"áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍"
		`

		tests := []token.Token{
			{Type: token.Comment, Literal: " áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍"},
			{Type: token.CommentBlock, Literal: " áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍 "},
			{Type: token.Docstring, Literal: " áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍 "},
			{Type: token.StringLiteral, Literal: "áéíóúàèìòùäëïöüâêîôûåæœçðñµßøœ你好こんにちは안녕하세요Привет👍"},
		}

		lex := NewLexer("test.urpc", input)
		tokens := lex.ReadTokens()

		// This test does not test newline and whitespace tokens to avoid verbosity
		// They're already tested in previous tests

		testableTokens := []token.Token{}
		for _, tok := range tokens {
			if tok.Type == token.Newline || tok.Type == token.Whitespace {
				continue
			}
			testableTokens = append(testableTokens, tok)
		}

		for i, test := range tests {
			tok := testableTokens[i]

			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})

	t.Run("TestLexerURPC", func(t *testing.T) {
		input := `
			// This test evaluates the lexer with a full URPC file.
			// It is used to ensure that the lexer is working correctly.

			version: 1

			""" Product is a type that represents a product. """
			deprecated("This type will be removed in v2.0")
			type Product {
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
					success: bool
					productId: string
						@uuid
				}
				
				meta {
					audit: true
					maxRetries: 3
				}
			}

			stream NewProduct {
				input {
					product: Product
				}
			}`

		tests := []token.Token{
			{Type: token.Comment, Literal: " This test evaluates the lexer with a full URPC file."},
			{Type: token.Comment, Literal: " It is used to ensure that the lexer is working correctly."},
			{Type: token.Version, Literal: "version"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.IntLiteral, Literal: "1"},
			{Type: token.Docstring, Literal: " Product is a type that represents a product. "},
			{Type: token.Deprecated, Literal: "deprecated"},
			{Type: token.LParen, Literal: "("},
			{Type: token.StringLiteral, Literal: "This type will be removed in v2.0"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.Type, Literal: "type"},
			{Type: token.Ident, Literal: "Product"},
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
			{Type: token.Docstring, Literal: " Creates a product and returns the product id. "},
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
			{Type: token.Ident, Literal: "error"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.StringLiteral, Literal: "Priority must be 1, 2, or 3"},
			{Type: token.RParen, Literal: ")"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Output, Literal: "output"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "success"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.Bool, Literal: "bool"},
			{Type: token.Ident, Literal: "productId"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.String, Literal: "string"},
			{Type: token.At, Literal: "@"},
			{Type: token.Ident, Literal: "uuid"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Ident, Literal: "meta"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "audit"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.TrueLiteral, Literal: "true"},
			{Type: token.Ident, Literal: "maxRetries"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.IntLiteral, Literal: "3"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Stream, Literal: "stream"},
			{Type: token.Ident, Literal: "NewProduct"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Input, Literal: "input"},
			{Type: token.LBrace, Literal: "{"},
			{Type: token.Ident, Literal: "product"},
			{Type: token.Colon, Literal: ":"},
			{Type: token.Ident, Literal: "Product"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.RBrace, Literal: "}"},
			{Type: token.Eof, Literal: ""},
		}

		lex := NewLexer("test.urpc", input)
		tokens := lex.ReadTokens()

		// This test does not test newline and whitespace tokens to avoid verbosity
		// They're already tested in previous tests

		testableTokens := []token.Token{}
		for _, tok := range tokens {
			if tok.Type == token.Newline || tok.Type == token.Whitespace {
				continue
			}
			testableTokens = append(testableTokens, tok)
		}

		for i, test := range tests {
			tok := testableTokens[i]

			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})
}
