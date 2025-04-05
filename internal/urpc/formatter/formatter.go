package formatter

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/parser"
)

// Format formats URPC code according to the spec, using 2 spaces for indentation.
func Format(filename, content string) (string, error) {
	if strings.TrimSpace(content) == "" {
		return "", nil
	}

	schema, err := parser.Parser.ParseString(filename, content)
	if err != nil {
		return "", fmt.Errorf("error parsing URPC: %w", err)
	}

	schFormatter := newSchemaFormatter(schema)
	formatted := schFormatter.format().String()

	// Ensure the formatted string ends with exactly one newline
	formatted = strings.TrimSpace(formatted)
	formatted += "\n"
	return formatted, nil
}

// schemaFormatter is a formatter for a schema.
type schemaFormatter struct {
	g                 *genkit.GenKit
	sch               *ast.Schema
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.SchemaChild
}

// newSchemaFormatter creates a new schema formatter and initializes all the necessary fields.
func newSchemaFormatter(sch *ast.Schema) *schemaFormatter {
	maxIndex := max(len(sch.Children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(sch.Children) < 1
	currentIndexChild := ast.SchemaChild{}

	if !currentIndexEOF {
		currentIndexChild = *sch.Children[0]
	}

	return &schemaFormatter{
		g:                 genkit.NewGenKit().WithSpaces(2),
		sch:               sch,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *schemaFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.SchemaChild{}

	if !currentIndexEOF {
		currentIndexChild = *f.sch.Children[currentIndex]
	}

	f.currentIndex = currentIndex
	f.currentIndexEOF = currentIndexEOF
	f.currentIndexChild = currentIndexChild
}

// peekChild returns information about the child at the current index +- offset.
//
// Returns:
//   - The child at the current index +- offset.
//   - The line diff between the peeked child and the current child.
//   - A boolean indicating if the peeked child is out of bounds (EOL).
func (f *schemaFormatter) peekChild(offset int) (ast.SchemaChild, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.SchemaChild{}
	lineDiff := ast.LineDiff{}

	if !peekIndexEOF {
		peekIndexChild = *f.sch.Children[peekIndex]
		lineDiff = ast.GetLineDiff(peekIndexChild, f.currentIndexChild)
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// format formats the entire schema, handling spacing and EOL comments.
//
// Returns the formatted genkit.GenKit.
func (f *schemaFormatter) format() *genkit.GenKit {
	for {
		if f.currentIndexEOF {
			break
		}

		switch f.currentIndexChild.Kind() {
		case ast.SchemaChildKindComment:
			f.formatComment()
		case ast.SchemaChildKindVersion:
			f.formatVersion()
		case ast.SchemaChildKindImport:
			f.formatImport()
		case ast.SchemaChildKindRule:
			f.formatRule()
		case ast.SchemaChildKindType:
			f.formatType()
		case ast.SchemaChildKindProc:
			f.formatProc()
		}

		f.loadNextChild()
	}

	return f.g
}

// LineAndComment writes a line of content to the formatter. It also handles inline comments.
func (f *schemaFormatter) LineAndComment(content string) {
	next, nextLineDiff, nextEOF := f.peekChild(1)

	// If next is an inline comment
	if !nextEOF && next.Kind() == ast.SchemaChildKindComment && nextLineDiff.StartToEnd == 0 {
		f.g.Inline(content)

		if next.Comment.Simple != nil {
			f.g.Linef(" //%s", *next.Comment.Simple)
		}

		if next.Comment.Block != nil {
			f.g.Linef(" /*%s*/", *next.Comment.Block)
		}

		// Skip the inline comment because it's already written
		f.loadNextChild()
		return
	}

	f.g.Line(content)
}

// LineAndCommentf is the same as Line but with a formatted string.
func (f *schemaFormatter) LineAndCommentf(format string, args ...any) {
	f.LineAndComment(fmt.Sprintf(format, args...))
}

func (f *schemaFormatter) formatComment() {
	_, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := false
	if !prevEOF {
		if prevLineDiff.StartToStart < -1 {
			shouldBreakBefore = true
		}
	}

	if shouldBreakBefore {
		f.g.Break()
	}

	if f.currentIndexChild.Comment.Simple != nil {
		f.g.Linef("//%s", *f.currentIndexChild.Comment.Simple)
	}

	if f.currentIndexChild.Comment.Block != nil {
		f.g.Linef("/*%s*/", *f.currentIndexChild.Comment.Block)
	}
}

func (f *schemaFormatter) formatVersion() {
	f.LineAndCommentf("version %d", f.currentIndexChild.Version.Number)
}

func (f *schemaFormatter) formatImport() {
	prev, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := false
	if !prevEOF {
		if prev.Kind() != ast.SchemaChildKindImport && prev.Kind() != ast.SchemaChildKindComment {
			shouldBreakBefore = true
		}

		if prevLineDiff.StartToStart < -1 {
			shouldBreakBefore = true
		}
	}

	if shouldBreakBefore {
		f.g.Break()
	}

	f.LineAndCommentf(`import "%s"`, f.currentIndexChild.Import.Path)
}

func (f *schemaFormatter) formatRule() {
	prev, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := false
	if !prevEOF {
		if prev.Kind() != ast.SchemaChildKindComment {
			shouldBreakBefore = true
		}

		if prevLineDiff.StartToStart < -1 {
			shouldBreakBefore = true
		}
	}

	if shouldBreakBefore {
		f.g.Break()
	}

	ruleFormatter := newRuleFormatter(f.currentIndexChild.Rule)
	f.LineAndComment(strings.TrimSpace(ruleFormatter.format().String()))
}

func (f *schemaFormatter) formatType() {
}

func (f *schemaFormatter) formatProc() {
}
