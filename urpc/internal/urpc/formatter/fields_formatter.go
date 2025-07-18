package formatter

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type fieldsFormatter struct {
	g                 *genkit.GenKit
	parent            ast.WithPositions
	fields            []*ast.FieldOrComment
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.FieldOrComment
}

func newFieldsFormatter(g *genkit.GenKit, parent ast.WithPositions, fields []*ast.FieldOrComment) *fieldsFormatter {
	if fields == nil {
		fields = []*ast.FieldOrComment{}
	}

	maxIndex := max(len(fields)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(fields) < 1
	currentIndexChild := ast.FieldOrComment{}

	if !currentIndexEOF {
		currentIndexChild = *fields[0]
	}

	return &fieldsFormatter{
		g:                 g,
		parent:            parent,
		fields:            fields,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *fieldsFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.FieldOrComment{}

	if !currentIndexEOF {
		currentIndexChild = *f.fields[currentIndex]
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
//   - A bool indicating if the peeked child is out of bounds (EOL).
func (f *fieldsFormatter) peekChild(offset int) (ast.FieldOrComment, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.FieldOrComment{}
	lineDiff := ast.LineDiff{}

	if !peekIndexEOF {
		peekIndexChild = *f.fields[peekIndex]
		lineDiff = ast.GetLineDiff(peekIndexChild, f.currentIndexChild)
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// LineAndComment writes a line of content to the formatter. It also handles inline comments.
func (f *fieldsFormatter) LineAndComment(content string) {
	next, nextLineDiff, nextEOF := f.peekChild(1)

	// If next is an inline comment
	if !nextEOF && next.Comment != nil && nextLineDiff.StartToEnd == 0 {
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
func (f *fieldsFormatter) LineAndCommentf(format string, args ...any) {
	f.LineAndComment(fmt.Sprintf(format, args...))
}

// format formats the entire rule, handling spacing and EOL comments.
//
// Returns the formatted genkit.GenKit.
func (f *fieldsFormatter) format() *genkit.GenKit {
	if f.currentIndexEOF {
		f.g.Inline("{}")
		return f.g
	}

	hasInlineComment := false
	if f.currentIndexChild.Comment != nil {
		lineDiff := ast.GetLineDiff(f.currentIndexChild, f.parent)
		if lineDiff.StartToStart == 0 {
			hasInlineComment = true
		}
	}

	if hasInlineComment {
		f.g.Inline("{ ")
	} else {
		f.g.Line("{")
	}

	f.g.Block(func() {
		for !f.currentIndexEOF {
			if f.currentIndexChild.Comment != nil {
				f.formatComment()
			}

			if f.currentIndexChild.Field != nil {
				f.formatField()
			}

			f.loadNextChild()
		}
	})

	f.g.Inline("}")

	return f.g
}

func (f *fieldsFormatter) formatComment() {
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

func (f *fieldsFormatter) formatField() {
	prev, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := false
	if !prevEOF {
		if prevLineDiff.EndToStart < -1 {
			shouldBreakBefore = true
		}
	}

	if shouldBreakBefore {
		f.g.Break()
	}

	if f.currentIndexChild.Field.Docstring != nil {
		// Add a break before the docstring if it's not the first field
		// and the previous element is a field
		if !prevEOF && prev.Field != nil {
			f.g.Break()
		}

		f.g.Inline(`"""`)
		f.g.Raw(f.currentIndexChild.Field.Docstring.Value)
		f.g.Raw(`"""`)
		f.g.Break()
	}

	// Force strict camel case
	if f.currentIndexChild.Field.Optional {
		f.g.Inlinef("%s?: ", strutil.ToCamelCase(f.currentIndexChild.Field.Name))
	} else {
		f.g.Inlinef("%s: ", strutil.ToCamelCase(f.currentIndexChild.Field.Name))
	}

	if f.currentIndexChild.Field.Type.Base.Named != nil {
		typeLiteral := *f.currentIndexChild.Field.Type.Base.Named
		// Force strict pascal case for non primitive types
		if !ast.IsPrimitiveType(typeLiteral) {
			typeLiteral = strutil.ToPascalCase(typeLiteral)
		}
		f.g.Inline(typeLiteral)
	}

	if f.currentIndexChild.Field.Type.Base.Object != nil {
		children := f.currentIndexChild.Field.Type.Base.Object.Children
		nestedFormatter := newFieldsFormatter(f.g, f.currentIndexChild, children)
		nestedFormatter.format()
	}

	if f.currentIndexChild.Field.Type.IsArray {
		f.g.Inline("[]")
	}

	f.LineAndComment("")
}
