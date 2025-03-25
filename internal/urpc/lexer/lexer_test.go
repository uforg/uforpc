package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/token"
)

func TestLexer(t *testing.T) {
	t.Run("TestNextToken", func(t *testing.T) {
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

		lex := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}
	})

	t.Run("TestNextTokenNewLines", func(t *testing.T) {
		input := ",:\n(){\n}\n[]@?\n"

		tests := []token.Token{
			{Type: token.COMMA, Literal: ",", FileName: "test.urpc", Line: 1, Column: 1},
			{Type: token.COLON, Literal: ":", FileName: "test.urpc", Line: 1, Column: 2},
			{Type: token.NEWLINE, Literal: "\n", FileName: "test.urpc", Line: 1, Column: 3},
			{Type: token.LPAREN, Literal: "(", FileName: "test.urpc", Line: 2, Column: 1},
			{Type: token.RPAREN, Literal: ")", FileName: "test.urpc", Line: 2, Column: 2},
			{Type: token.LBRACE, Literal: "{", FileName: "test.urpc", Line: 2, Column: 3},
			{Type: token.NEWLINE, Literal: "\n", FileName: "test.urpc", Line: 2, Column: 4},
			{Type: token.RBRACE, Literal: "}", FileName: "test.urpc", Line: 3, Column: 1},
			{Type: token.NEWLINE, Literal: "\n", FileName: "test.urpc", Line: 3, Column: 2},
			{Type: token.LBRACKET, Literal: "[", FileName: "test.urpc", Line: 4, Column: 1},
			{Type: token.RBRACKET, Literal: "]", FileName: "test.urpc", Line: 4, Column: 2},
			{Type: token.AT, Literal: "@", FileName: "test.urpc", Line: 4, Column: 3},
			{Type: token.QUESTION, Literal: "?", FileName: "test.urpc", Line: 4, Column: 4},
			{Type: token.NEWLINE, Literal: "\n", FileName: "test.urpc", Line: 4, Column: 5},
			{Type: token.EOF, Literal: "", FileName: "test.urpc", Line: 5, Column: 1},
		}

		lex := NewLexer("test.urpc", input)
		for i, test := range tests {
			tok := lex.NextToken()
			require.Equal(t, test.Type, tok.Type, "test %d", i)
			require.Equal(t, test.Literal, tok.Literal, "test %d", i)
			require.Equal(t, test.FileName, tok.FileName, "test %d", i)
			require.Equal(t, test.Line, tok.Line, "test %d", i)
			require.Equal(t, test.Column, tok.Column, "test %d", i)
		}
	})
}
