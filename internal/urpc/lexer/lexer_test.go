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
			{Type: token.COMMA, Literal: ",", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", Line: 1, Column: 2},
			{Type: token.LPAREN, Literal: "(", FileName: "test.urpc", Line: 1, Column: 3},
			{Type: token.RPAREN, Literal: ")", FileName: "test.urpc", Line: 1, Column: 4},
			{Type: token.LBRACE, Literal: "{", FileName: "test.urpc", Line: 1, Column: 5},
			{Type: token.RBRACE, Literal: "}", FileName: "test.urpc", Line: 1, Column: 6},
			{Type: token.LBRACKET, Literal: "[", FileName: "test.urpc", Line: 1, Column: 7},
			{Type: token.RBRACKET, Literal: "]", FileName: "test.urpc", Line: 1, Column: 8},
			{Type: token.AT, Literal: "@", FileName: "test.urpc", Line: 1, Column: 9},
			{Type: token.QUESTION, Literal: "?", FileName: "test.urpc", Line: 1, Column: 10},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 11},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerNewLines", func(t *testing.T) {
		input := ",:\n(){\n}\n[]@?\n"

		tests := []token.Token{
			{Type: token.COMMA, Literal: ",", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", Line: 1, Column: 2},
			{Type: token.LPAREN, Literal: "(", FileName: "test.urpc", Line: 2, Column: 1},
			{Type: token.RPAREN, Literal: ")", FileName: "test.urpc", Line: 2, Column: 2},
			{Type: token.LBRACE, Literal: "{", FileName: "test.urpc", Line: 2, Column: 3},
			{Type: token.RBRACE, Literal: "}", FileName: "test.urpc", Line: 3, Column: 1},
			{Type: token.LBRACKET, Literal: "[", FileName: "test.urpc", Line: 4, Column: 1},
			{Type: token.RBRACKET, Literal: "]", FileName: "test.urpc", Line: 4, Column: 2},
			{Type: token.AT, Literal: "@", FileName: "test.urpc", Line: 4, Column: 3},
			{Type: token.QUESTION, Literal: "?", FileName: "test.urpc", Line: 4, Column: 4},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 5, Column: 1},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerKeywords", func(t *testing.T) {
		input := "version type proc input output meta error true false"

		tests := []token.Token{
			{Type: token.VERSION, Literal: "version", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.TYPE, Literal: "type", FileName: "test.urpc", Line: 1, Column: 9},
			{Type: token.PROC, Literal: "proc", FileName: "test.urpc", Line: 1, Column: 14},
			{Type: token.INPUT, Literal: "input", FileName: "test.urpc", Line: 1, Column: 19},
			{Type: token.OUTPUT, Literal: "output", FileName: "test.urpc", Line: 1, Column: 25},
			{Type: token.META, Literal: "meta", FileName: "test.urpc", Line: 1, Column: 32},
			{Type: token.ERROR, Literal: "error", FileName: "test.urpc", Line: 1, Column: 37},
			{Type: token.TRUE, Literal: "true", FileName: "test.urpc", Line: 1, Column: 43},
			{Type: token.FALSE, Literal: "false", FileName: "test.urpc", Line: 1, Column: 48},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 53},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerIdentifiers", func(t *testing.T) {
		input := "hello world someIdentifier"

		tests := []token.Token{
			{Type: token.IDENT, Literal: "hello", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.IDENT, Literal: "world", FileName: "test.urpc", Line: 1, Column: 7},
			{Type: token.IDENT, Literal: "someIdentifier", FileName: "test.urpc", Line: 1, Column: 13},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 27},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerIdentifiersWithNumbers", func(t *testing.T) {
		input := "hello123 world456 someIdentifier789"

		tests := []token.Token{
			{Type: token.IDENT, Literal: "hello123"},
			{Type: token.IDENT, Literal: "world456"},
			{Type: token.IDENT, Literal: "someIdentifier789"},
			{Type: token.EOF, Literal: ""},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
		}
	})

	t.Run("TestLexerNumbers", func(t *testing.T) {
		input := "1 2 3 456 789"

		tests := []token.Token{
			{Type: token.INT, Literal: "1", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.INT, Literal: "2", FileName: "test.urpc", Line: 1, Column: 3},
			{Type: token.INT, Literal: "3", FileName: "test.urpc", Line: 1, Column: 5},
			{Type: token.INT, Literal: "456", FileName: "test.urpc", Line: 1, Column: 7},
			{Type: token.INT, Literal: "789", FileName: "test.urpc", Line: 1, Column: 11},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 14},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerFloats", func(t *testing.T) {
		input := "1.2 3.45 67.89 1.2.3.4"

		tests := []token.Token{
			{Type: token.FLOAT, Literal: "1.2", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.FLOAT, Literal: "3.45", FileName: "test.urpc", Line: 1, Column: 5},
			{Type: token.FLOAT, Literal: "67.89", FileName: "test.urpc", Line: 1, Column: 10},
			{Type: token.FLOAT, Literal: "1.2", FileName: "test.urpc", Line: 1, Column: 16},
			{Type: token.ILLEGAL, Literal: ".", FileName: "test.urpc", Line: 1, Column: 19},
			{Type: token.FLOAT, Literal: "3.4", FileName: "test.urpc", Line: 1, Column: 20},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 23},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerStrings", func(t *testing.T) {
		input := `"hello" "world" "hello world!"test`

		tests := []token.Token{
			{Type: token.STRING, Literal: "hello", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.STRING, Literal: "world", FileName: "test.urpc", Line: 1, Column: 9},
			{Type: token.STRING, Literal: "hello world!", FileName: "test.urpc", Line: 1, Column: 17},
			{Type: token.IDENT, Literal: "test", FileName: "test.urpc", Line: 1, Column: 31},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 35},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerIllegal", func(t *testing.T) {
		input := "$ % ^ & ."

		tests := []token.Token{
			{Type: token.ILLEGAL, Literal: "$", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.ILLEGAL, Literal: "%", FileName: "test.urpc", Line: 1, Column: 3},
			{Type: token.ILLEGAL, Literal: "^", FileName: "test.urpc", Line: 1, Column: 5},
			{Type: token.ILLEGAL, Literal: "&", FileName: "test.urpc", Line: 1, Column: 7},
			{Type: token.ILLEGAL, Literal: ".", FileName: "test.urpc", Line: 1, Column: 9},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 10},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerComments", func(t *testing.T) {
		input := "// This is a comment\nversion: 1"

		tests := []token.Token{
			{Type: token.COMMENT, Literal: "This is a comment", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.VERSION, Literal: "version", FileName: "test.urpc", Line: 2, Column: 1},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", Line: 2, Column: 8},
			{Type: token.INT, Literal: "1", FileName: "test.urpc", Line: 2, Column: 10},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 2, Column: 11},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerUnterminatedString", func(t *testing.T) {
		input := `"hello`

		tests := []token.Token{
			{Type: token.ILLEGAL, Literal: "\"hello", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 7},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerDocstrings", func(t *testing.T) {
		input := `""" This is a docstring """`

		tests := []token.Token{
			{Type: token.DOCSTRING, Literal: "This is a docstring", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 1, Column: 28},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}

		lex2 := NewLexer("test.urpc", input)
		tokens := lex2.ReadTokens()
		require.Equal(t, tests, tokens)
	})

	t.Run("TestLexerMultilineDocstrings", func(t *testing.T) {
		input := "\"\"\" This is a multiline docstring\nwith multiple lines \"\"\"\n"

		tests := []token.Token{
			{Type: token.DOCSTRING, Literal: "This is a multiline docstring\nwith multiple lines"},
			{Type: token.EOF, Literal: ""},
		}

		lex1 := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex1.NextToken()
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
