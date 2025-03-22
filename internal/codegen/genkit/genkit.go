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

// Raw writes a raw content without any indentation or line breaks.
func (g *GenKit) Raw(content string) *GenKit {
	g.sb.WriteString(content)
	return g
}

// Rawf writes a formatted raw content without any indentation or line breaks.
func (g *GenKit) Rawf(format string, args ...any) *GenKit {
	return g.Raw(fmt.Sprintf(format, args...))
}

// Break literally writes a line break.
func (g *GenKit) Break() *GenKit {
	return g.Raw("\n")
}

// Inline writes a line with the current indentation and does not add a line break
// before the line content. If the line contains newlines, each line will be properly indented.
func (g *GenKit) Inline(line string) *GenKit {
	if line != "" {
		sublines := strings.Split(line, "\n")
		for idx, subline := range sublines {
			if idx > 0 {
				g.Raw("\n")
			}
			if subline != "" {
				g.Raw(strings.Repeat(g.indentString, g.indentLevel))
				g.Raw(subline)
			}
		}
	}
	return g
}

// Inlinef writes a formatted line with the current indentation and does not add a line break
// before the line content.
func (g *GenKit) Inlinef(format string, args ...any) *GenKit {
	return g.Inline(fmt.Sprintf(format, args...))
}

// Line writes a line with the current indentation and adds a line break
// before the line content.
func (g *GenKit) Line(line string) *GenKit {
	g.Break()
	return g.Inline(line)
}

// Linef writes a formatted line with the current indentation and adds a line break
// before the line content.
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
