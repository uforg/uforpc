package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"HeLLO", "HeLLO"},
		{"HELLO", "HELLO"},
		{"hello world", "Hello world"},
		{"", ""},
		{"123", "123"},
		{"123abc", "123abc"},
		{"123abc123", "123abc123"},
		{"123abc123", "123abc123"},
		{"123abc123", "123abc123"},
	}

	for _, test := range tests {
		result := Capitalize(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestCapitalizeStrict(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"HeLLO", "Hello"},
		{"HELLO", "Hello"},
		{"hello world", "Hello world"},
		{"Hello World", "Hello world"},
		{"", ""},
		{"123", "123"},
		{"123abc", "123abc"},
		{"123abc123", "123abc123"},
		{"123abc123", "123abc123"},
		{"123abc123", "123abc123"},
	}

	for _, test := range tests {
		result := CapitalizeStrict(test.input)
		assert.Equal(t, test.expected, result)
	}
}
