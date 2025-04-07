package formatter

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/util/strutil"
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

func newFieldsFormatter(parent ast.WithPositions, fields []*ast.FieldOrComment) *fieldsFormatter {
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
		g:                 genkit.NewGenKit().WithSpaces(2),
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
//   - A boolean indicating if the peeked child is out of bounds (EOL).
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
		for {
			if f.currentIndexEOF {
				break
			}

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
	_, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := false
	if !prevEOF {
		if prevLineDiff.EndToStart < -1 {
			shouldBreakBefore = true
		}
	}

	if shouldBreakBefore {
		f.g.Break()
	}

	if f.currentIndexChild.Field.Optional {
		f.g.Inlinef("%s?: ", f.currentIndexChild.Field.Name)
	} else {
		f.g.Inlinef("%s: ", f.currentIndexChild.Field.Name)
	}

	typeLiteral := ""

	if f.currentIndexChild.Field.Type.Base.Named != nil {
		typeLiteral = *f.currentIndexChild.Field.Type.Base.Named
	}

	if f.currentIndexChild.Field.Type.Base.Object != nil {
		children := f.currentIndexChild.Field.Type.Base.Object.Children
		nestedFormatter := newFieldsFormatter(f.currentIndexChild, children)
		typeLiteral = strings.TrimSpace(nestedFormatter.format().String())
	}

	for range f.currentIndexChild.Field.Type.Depth {
		typeLiteral = typeLiteral + "[]"
	}

	if f.currentIndexChild.Field.Children != nil {
		rulesFormatter := newFieldRulesFormatter(f.currentIndexChild, f.currentIndexChild.Field.Children)
		children := rulesFormatter.format()
		typeLiteral = typeLiteral + children.String()
	}

	f.LineAndComment(typeLiteral)
}

////////////////////
////////////////////
////////////////////

type fieldRulesFormatter struct {
	g                 *genkit.GenKit
	parent            ast.WithPositions
	children          []*ast.FieldChild
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.FieldChild
}

func newFieldRulesFormatter(parent ast.WithPositions, children []*ast.FieldChild) *fieldRulesFormatter {
	if children == nil {
		children = []*ast.FieldChild{}
	}

	maxIndex := max(len(children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(children) < 1
	currentIndexChild := ast.FieldChild{}

	if !currentIndexEOF {
		currentIndexChild = *children[0]
	}

	return &fieldRulesFormatter{
		g:                 genkit.NewGenKit().WithSpaces(2),
		parent:            parent,
		children:          children,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *fieldRulesFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.FieldChild{}

	if !currentIndexEOF {
		currentIndexChild = *f.children[currentIndex]
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
func (f *fieldRulesFormatter) peekChild(offset int) (ast.FieldChild, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.FieldChild{}
	lineDiff := ast.LineDiff{}

	if !peekIndexEOF {
		peekIndexChild = *f.children[peekIndex]
		lineDiff = ast.GetLineDiff(peekIndexChild, f.currentIndexChild)
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// format formats the entire rule, handling spacing and EOL comments.
func (f *fieldRulesFormatter) format() *genkit.GenKit {
	if f.currentIndexEOF {
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
		f.g.Inline(" ")
		f.formatComment()
		f.loadNextChild()
	}

	f.g.Block(func() {
		for {
			if f.currentIndexEOF {
				break
			}

			if f.currentIndexChild.Comment != nil {
				f.formatComment()
			}

			if f.currentIndexChild.Rule != nil {
				f.formatRule()
			}

			f.loadNextChild()
		}
	})

	return f.g
}

func (f *fieldRulesFormatter) formatComment() {
	_, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := !prevEOF && prevLineDiff.EndToStart < -1
	if shouldBreakBefore {
		f.g.Break()
	}

	shouldSpaceBefore := !prevEOF && prevLineDiff.EndToStart == 0
	if shouldSpaceBefore {
		f.g.Inline(" ")
	}

	if f.currentIndexChild.Comment.Simple != nil {
		f.g.Inlinef("//%s", *f.currentIndexChild.Comment.Simple)
	}

	if f.currentIndexChild.Comment.Block != nil {
		f.g.Inlinef("/*%s*/", *f.currentIndexChild.Comment.Block)
	}
}

func (f *fieldRulesFormatter) formatRule() {
	_, prevLineDiff, prevEOF := f.peekChild(-1)

	shouldBreakBefore := !prevEOF && prevLineDiff.EndToStart < -1
	if shouldBreakBefore {
		f.g.Break()
	}

	f.g.Break()
	f.g.Inlinef("@%s", f.currentIndexChild.Rule.Name)

	if f.currentIndexChild.Rule.Body != nil {
		f.g.Inline("(")
		hasParam := false

		if f.currentIndexChild.Rule.Body.ParamSingle != nil {
			f.g.Inlinef("%s", f.currentIndexChild.Rule.Body.ParamSingle.String())
			hasParam = true
		}

		if f.currentIndexChild.Rule.Body.ParamListString != nil {
			f.g.Inline("[")
			for i, param := range f.currentIndexChild.Rule.Body.ParamListString {
				if i > 0 {
					f.g.Inline(", ")
				}
				f.g.Inlinef(`"%s"`, strutil.EscapeQuotes(param))
			}
			f.g.Inline("]")
			hasParam = true
		}

		if f.currentIndexChild.Rule.Body.ParamListInt != nil {
			f.g.Inlinef("[%s]", strings.Join(f.currentIndexChild.Rule.Body.ParamListInt, ", "))
			hasParam = true
		}

		if f.currentIndexChild.Rule.Body.ParamListFloat != nil {
			f.g.Inlinef("[%s]", strings.Join(f.currentIndexChild.Rule.Body.ParamListFloat, ", "))
			hasParam = true
		}

		if f.currentIndexChild.Rule.Body.ParamListBoolean != nil {
			f.g.Inlinef("[%s]", strings.Join(f.currentIndexChild.Rule.Body.ParamListBoolean, ", "))
			hasParam = true
		}

		if f.currentIndexChild.Rule.Body.Error != nil {
			if hasParam {
				f.g.Inline(", ")
			}
			f.g.Inlinef(`error: "%s"`, *f.currentIndexChild.Rule.Body.Error)
		}

		f.g.Inline(")")
	}
}
