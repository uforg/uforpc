package formatter

import (
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type streamFormatter struct {
	g                 *genkit.GenKit
	streamDecl        *ast.StreamDecl
	children          []*ast.ProcOrStreamDeclChild
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.ProcOrStreamDeclChild
}

func newStreamFormatter(streamDecl *ast.StreamDecl) *streamFormatter {
	if streamDecl == nil {
		streamDecl = &ast.StreamDecl{}
	}

	if streamDecl.Children == nil {
		streamDecl.Children = []*ast.ProcOrStreamDeclChild{}
	}

	maxIndex := max(len(streamDecl.Children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(streamDecl.Children) < 1
	currentIndexChild := ast.ProcOrStreamDeclChild{}

	if !currentIndexEOF {
		currentIndexChild = *streamDecl.Children[0]
	}

	return &streamFormatter{
		g:                 genkit.NewGenKit().WithSpaces(2),
		streamDecl:        streamDecl,
		children:          streamDecl.Children,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *streamFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.ProcOrStreamDeclChild{}

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
//   - A bool indicating if the peeked child is out of bounds (EOL).
func (f *streamFormatter) peekChild(offset int) (ast.ProcOrStreamDeclChild, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.ProcOrStreamDeclChild{}
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
func (f *streamFormatter) format() *genkit.GenKit {
	if f.streamDecl.Docstring != nil {
		f.g.Linef(`"""%s"""`, f.streamDecl.Docstring.Value)
	}

	if f.streamDecl.Deprecated != nil {
		if f.streamDecl.Deprecated.Message == nil {
			f.g.Inline("deprecated ")
		}
		if f.streamDecl.Deprecated.Message != nil {
			f.g.Linef("deprecated(\"%s\")", strutil.EscapeQuotes(*f.streamDecl.Deprecated.Message))
		}
	}

	f.g.Inlinef(`stream %s `, f.streamDecl.Name)

	if len(f.streamDecl.Children) < 1 {
		f.g.Inline("{}")
		return f.g
	}

	hasInlineComment := false
	if f.currentIndexChild.Comment != nil {
		lineDiff := ast.GetLineDiff(f.currentIndexChild, f.streamDecl)
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

func (f *streamFormatter) formatComment() {
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

func (f *streamFormatter) breakBeforeBlock() {
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

func (f *streamFormatter) formatInput() {
	f.breakBeforeBlock()
	f.g.Inline("input ")
	fieldsFormatter := newFieldsFormatter(f.currentIndexChild, f.currentIndexChild.Input.Children)
	f.g.Line(strings.TrimSpace(fieldsFormatter.format().String()))
}

func (f *streamFormatter) formatOutput() {
	f.breakBeforeBlock()
	f.g.Inline("output ")
	fieldsFormatter := newFieldsFormatter(f.currentIndexChild, f.currentIndexChild.Output.Children)
	f.g.Line(strings.TrimSpace(fieldsFormatter.format().String()))
}

func (f *streamFormatter) formatMeta() {
	f.breakBeforeBlock()
	f.g.Inline("meta ")
	metaFormatter := newProcOrStreamMetaFormatter(f.currentIndexChild.Meta)
	f.g.Line(strings.TrimSpace(metaFormatter.format().String()))
}
