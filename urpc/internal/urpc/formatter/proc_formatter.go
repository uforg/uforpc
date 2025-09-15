package formatter

import (
	"github.com/uforg/ufogenkit"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type procFormatter struct {
	g                 *ufogenkit.GenKit
	procDecl          *ast.ProcDecl
	children          []*ast.ProcOrStreamDeclChild
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.ProcOrStreamDeclChild
}

func newProcFormatter(g *ufogenkit.GenKit, procDecl *ast.ProcDecl) *procFormatter {
	if procDecl == nil {
		procDecl = &ast.ProcDecl{}
	}

	if procDecl.Children == nil {
		procDecl.Children = []*ast.ProcOrStreamDeclChild{}
	}

	maxIndex := max(len(procDecl.Children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(procDecl.Children) < 1
	currentIndexChild := ast.ProcOrStreamDeclChild{}

	if !currentIndexEOF {
		currentIndexChild = *procDecl.Children[0]
	}

	return &procFormatter{
		g:                 g,
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
func (f *procFormatter) peekChild(offset int) (ast.ProcOrStreamDeclChild, ast.LineDiff, bool) {
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
func (f *procFormatter) format() *ufogenkit.GenKit {
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

	// Force strict pascal case
	f.g.Inlinef(`proc %s `, strutil.ToPascalCase(f.procDecl.Name))

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
	fieldsFormatter := newFieldsFormatter(f.g, f.currentIndexChild, f.currentIndexChild.Input.Children)
	fieldsFormatter.format()
	f.g.Break()
}

func (f *procFormatter) formatOutput() {
	f.breakBeforeBlock()
	f.g.Inline("output ")
	fieldsFormatter := newFieldsFormatter(f.g, f.currentIndexChild, f.currentIndexChild.Output.Children)
	fieldsFormatter.format()
	f.g.Break()
}
