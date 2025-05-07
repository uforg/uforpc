package formatter

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type procFormatter struct {
	g                 *genkit.GenKit
	procDecl          *ast.ProcDecl
	children          []*ast.ProcDeclChild
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.ProcDeclChild
}

func newProcFormatter(procDecl *ast.ProcDecl) *procFormatter {
	if procDecl == nil {
		procDecl = &ast.ProcDecl{}
	}

	if procDecl.Children == nil {
		procDecl.Children = []*ast.ProcDeclChild{}
	}

	maxIndex := max(len(procDecl.Children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(procDecl.Children) < 1
	currentIndexChild := ast.ProcDeclChild{}

	if !currentIndexEOF {
		currentIndexChild = *procDecl.Children[0]
	}

	return &procFormatter{
		g:                 genkit.NewGenKit().WithSpaces(2),
		procDecl:          procDecl,
		children:          procDecl.Children,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *procFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.ProcDeclChild{}

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
func (f *procFormatter) peekChild(offset int) (ast.ProcDeclChild, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.ProcDeclChild{}
	lineDiff := ast.LineDiff{}

	if !peekIndexEOF {
		peekIndexChild = *f.children[peekIndex]
		lineDiff = ast.GetLineDiff(peekIndexChild, f.currentIndexChild)
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// format formats the entire procDecl, handling spacing and EOL comments.
//
// Returns the formatted genkit.GenKit.
func (f *procFormatter) format() *genkit.GenKit {
	if f.procDecl.Docstring != nil {
		f.g.Linef(`"""%s"""`, f.procDecl.Docstring.Value)
	}

	if f.procDecl.Deprecated != nil {
		if f.procDecl.Deprecated.Message == nil {
			f.g.Inline("deprecated ")
		}
		if f.procDecl.Deprecated.Message != nil {
			f.g.Linef("deprecated(\"%s\")", strutil.EscapeQuotes(*f.procDecl.Deprecated.Message))
		}
	}

	f.g.Inlinef(`proc %s `, f.procDecl.Name)

	if len(f.procDecl.Children) < 1 {
		f.g.Inline("{}")
		return f.g
	}

	hasInlineComment := false
	if f.currentIndexChild.Comment != nil {
		lineDiff := ast.GetLineDiff(f.currentIndexChild, f.procDecl)
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

			if f.currentIndexChild.Input != nil {
				f.formatInput()
			}

			if f.currentIndexChild.Output != nil {
				f.formatOutput()
			}

			if f.currentIndexChild.Meta != nil {
				f.formatMeta()
			}

			f.loadNextChild()
		}
	})

	f.g.Inline("}")

	return f.g
}

func (f *procFormatter) formatComment() {
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

func (f *procFormatter) breakBeforeBlock() {
	prev, prevLineDiff, prevEOF := f.peekChild(-1)
	prevWasComment := prev.Comment != nil

	if prevEOF {
		return
	}

	if prevWasComment {
		if prevLineDiff.StartToStart < -1 {
			f.g.Break()
			return
		}
		return
	}

	f.g.Break()
}

func (f *procFormatter) formatInput() {
	f.breakBeforeBlock()
	f.g.Inline("input ")
	fieldsFormatter := newFieldsFormatter(f.currentIndexChild, f.currentIndexChild.Input.Children)
	f.g.Line(strings.TrimSpace(fieldsFormatter.format().String()))
}

func (f *procFormatter) formatOutput() {
	f.breakBeforeBlock()
	f.g.Inline("output ")
	fieldsFormatter := newFieldsFormatter(f.currentIndexChild, f.currentIndexChild.Output.Children)
	f.g.Line(strings.TrimSpace(fieldsFormatter.format().String()))
}

func (f *procFormatter) formatMeta() {
	f.breakBeforeBlock()
	f.g.Inline("meta ")
	metaFormatter := newProcMetaFormatter(f.currentIndexChild.Meta)
	f.g.Line(strings.TrimSpace(metaFormatter.format().String()))
}

////////////////////
////////////////////
////////////////////

type procMetaFormatter struct {
	g                 *genkit.GenKit
	procMeta          *ast.ProcDeclChildMeta
	children          []*ast.ProcDeclChildMetaChild
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.ProcDeclChildMetaChild
}

func newProcMetaFormatter(procMeta *ast.ProcDeclChildMeta) *procMetaFormatter {
	if procMeta == nil {
		procMeta = &ast.ProcDeclChildMeta{}
	}

	if procMeta.Children == nil {
		procMeta.Children = []*ast.ProcDeclChildMetaChild{}
	}

	maxIndex := max(len(procMeta.Children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(procMeta.Children) < 1
	currentIndexChild := ast.ProcDeclChildMetaChild{}

	if !currentIndexEOF {
		currentIndexChild = *procMeta.Children[0]
	}

	return &procMetaFormatter{
		g:                 genkit.NewGenKit().WithSpaces(2),
		procMeta:          procMeta,
		children:          procMeta.Children,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *procMetaFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.ProcDeclChildMetaChild{}

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
func (f *procMetaFormatter) peekChild(offset int) (ast.ProcDeclChildMetaChild, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.ProcDeclChildMetaChild{}
	lineDiff := ast.LineDiff{}

	if !peekIndexEOF {
		peekIndexChild = *f.children[peekIndex]
		lineDiff = ast.GetLineDiff(peekIndexChild, f.currentIndexChild)
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// LineAndComment writes a line of content to the formatter. It also handles inline comments.
func (f *procMetaFormatter) LineAndComment(content string) {
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
func (f *procMetaFormatter) LineAndCommentf(format string, args ...any) {
	f.LineAndComment(fmt.Sprintf(format, args...))
}

// format formats the entire rule, handling spacing and EOL comments.
//
// Returns the formatted genkit.GenKit.
func (f *procMetaFormatter) format() *genkit.GenKit {
	if f.currentIndexEOF {
		f.g.Inline("{}")
		return f.g
	}

	hasInlineComment := false
	if f.currentIndexChild.Comment != nil {
		lineDiff := ast.GetLineDiff(f.currentIndexChild, f.procMeta)
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

			if f.currentIndexChild.KV != nil {
				f.formatKV()
			}

			f.loadNextChild()
		}
	})

	f.g.Inline("}")

	return f.g
}

func (f *procMetaFormatter) formatComment() {
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

func (f *procMetaFormatter) formatKV() {
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

	key := f.currentIndexChild.KV.Key
	value := f.currentIndexChild.KV.Value.String()

	f.LineAndCommentf("%s: %s", key, value)
}
