// Package genkit provides a simple and powerful code generation toolkit.
package genkit

import (
	"fmt"
	"strings"
)

// GenKit provides a fluent interface for generating code in any programming language.
// It handles indentation and line breaks automatically.
type GenKit struct {
	sb           strings.Builder
	indentLevel  int
	indentString string
}

// NewGenKit creates a new GenKit instance with default settings (2 spaces indentation).
func NewGenKit() *GenKit {
	return &GenKit{
		sb:           strings.Builder{},
		indentLevel:  0,
		indentString: "  ",
	}
}

// WithSpaces sets the indentation to use the specified number of spaces.
func (g *GenKit) WithSpaces(spaces int) *GenKit {
	g.indentString = strings.Repeat(" ", spaces)
	return g
}

// WithTabs sets the indentation to use tabs.
func (g *GenKit) WithTabs() *GenKit {
	g.indentString = "\t"
	return g
}

// Indent increases the indentation level by 1.
func (g *GenKit) Indent() *GenKit {
	g.indentLevel++
	return g
}

// Dedent decreases the indentation level by 1.
func (g *GenKit) Dedent() *GenKit {
	if g.indentLevel > 0 {
		g.indentLevel--
	}
	return g
}

// Line writes a line with the current indentation.
func (g *GenKit) Line(line ...string) *GenKit {
	pickedLine := ""
	if len(line) > 0 {
		pickedLine = line[0]
	}

	if pickedLine != "" {
		g.sb.WriteString(strings.Repeat(g.indentString, g.indentLevel))
		g.sb.WriteString(pickedLine)
	}
	g.sb.WriteString("\n")
	return g
}

// Linef writes a formatted line with the current indentation.
func (g *GenKit) Linef(format string, args ...any) *GenKit {
	return g.Line(fmt.Sprintf(format, args...))
}

// Block executes a function within an indented block. In other words
// all the lines written within the function will be indented one
// more than the current indentation level.
func (g *GenKit) Block(fn func()) *GenKit {
	g.Indent()
	fn()
	g.Dedent()
	return g
}

// String returns the generated code as a string.
func (g *GenKit) String() string {
	return g.sb.String()
}
