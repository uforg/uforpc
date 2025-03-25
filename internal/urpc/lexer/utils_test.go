package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLetter(t *testing.T) {
	require.True(t, isLetter('a'))
	require.True(t, isLetter('z'))
	require.True(t, isLetter('A'))
	require.True(t, isLetter('Z'))

	require.False(t, isLetter('1'))
	require.False(t, isLetter('@'))
	require.False(t, isLetter('.'))
	require.False(t, isLetter(' '))
	require.False(t, isLetter('\n'))
}

func TestIsDigit(t *testing.T) {
	require.True(t, isDigit('0'))
	require.True(t, isDigit('9'))

	require.False(t, isDigit('a'))
	require.False(t, isDigit('z'))
	require.False(t, isDigit('A'))
	require.False(t, isDigit('Z'))
	require.False(t, isDigit('@'))
	require.False(t, isDigit('.'))
	require.False(t, isDigit(' '))
	require.False(t, isDigit('\n'))
}

func TestIsWhitespace(t *testing.T) {
	require.True(t, isWhitespace(' '))
	require.True(t, isWhitespace('\t'))
	require.True(t, isWhitespace('\r'))

	require.False(t, isWhitespace('\n'))
	require.False(t, isWhitespace('a'))
	require.False(t, isWhitespace('z'))
	require.False(t, isWhitespace('A'))
	require.False(t, isWhitespace('Z'))
	require.False(t, isWhitespace('1'))
	require.False(t, isWhitespace('@'))
	require.False(t, isWhitespace('.'))
}

func TestIsNewline(t *testing.T) {
	require.True(t, isNewline('\n'))

	require.False(t, isNewline(' '))
	require.False(t, isNewline('a'))
	require.False(t, isNewline('z'))
	require.False(t, isNewline('A'))
	require.False(t, isNewline('Z'))
	require.False(t, isNewline('1'))
	require.False(t, isNewline('@'))
	require.False(t, isNewline('.'))
}

func TestContainsDecimalPoint(t *testing.T) {
	require.True(t, containsDecimalPoint("123.456"))
	require.True(t, containsDecimalPoint("123.456"))
	require.True(t, containsDecimalPoint("123.456.789"))

	require.False(t, containsDecimalPoint("123"))
	require.False(t, containsDecimalPoint("abcde"))
}
