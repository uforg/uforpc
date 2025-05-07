package formatter

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type ruleFormatter struct {
	g                 *genkit.GenKit
	ruleDecl          *ast.RuleDecl
	maxIndex          int
	currentIndex      int
	currentIndexEOF   bool
	currentIndexChild ast.RuleDeclChild
}

func newRuleFormatter(rule *ast.RuleDecl) *ruleFormatter {
	if rule == nil {
		rule = &ast.RuleDecl{}
	}

	maxIndex := max(len(rule.Children)-1, 0)
	currentIndex := 0
	currentIndexEOF := len(rule.Children) < 1
	currentIndexChild := ast.RuleDeclChild{}

	if !currentIndexEOF {
		currentIndexChild = *rule.Children[0]
	}

	return &ruleFormatter{
		g:                 genkit.NewGenKit().WithSpaces(2),
		ruleDecl:          rule,
		maxIndex:          maxIndex,
		currentIndex:      currentIndex,
		currentIndexEOF:   currentIndexEOF,
		currentIndexChild: currentIndexChild,
	}
}

// loadNextChild moves the current index to the next child.
func (f *ruleFormatter) loadNextChild() {
	currentIndex := f.currentIndex + 1
	currentIndexEOF := currentIndex > f.maxIndex
	currentIndexChild := ast.RuleDeclChild{}

	if !currentIndexEOF {
		currentIndexChild = *f.ruleDecl.Children[currentIndex]
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
func (f *ruleFormatter) peekChild(offset int) (ast.RuleDeclChild, ast.LineDiff, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.RuleDeclChild{}
	lineDiff := ast.LineDiff{}

	if !peekIndexEOF {
		peekIndexChild = *f.ruleDecl.Children[peekIndex]
		lineDiff = ast.GetLineDiff(peekIndexChild, f.currentIndexChild)
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// LineAndComment writes a line of content to the formatter. It also handles inline comments.
func (f *ruleFormatter) LineAndComment(content string) {
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
func (f *ruleFormatter) LineAndCommentf(format string, args ...any) {
	f.LineAndComment(fmt.Sprintf(format, args...))
}

// format formats the entire rule, handling spacing and EOL comments.
//
// Returns the formatted genkit.GenKit.
func (f *ruleFormatter) format() *genkit.GenKit {
	if f.ruleDecl.Docstring != nil {
		f.g.Linef(`"""%s"""`, f.ruleDecl.Docstring.Value)
	}

	if f.ruleDecl.Deprecated != nil {
		if f.ruleDecl.Deprecated.Message == nil {
			f.g.Inline("deprecated ")
		}
		if f.ruleDecl.Deprecated.Message != nil {
			f.g.Linef("deprecated(\"%s\")", strutil.EscapeQuotes(*f.ruleDecl.Deprecated.Message))
		}
	}

	f.g.Inlinef(`rule @%s `, f.ruleDecl.Name)

	if f.currentIndexEOF {
		f.g.Inline("{}")
		return f.g
	}

	hasInlineComment := false
	if f.currentIndexChild.Comment != nil {
		lineDiff := ast.GetLineDiff(f.currentIndexChild, f.ruleDecl)
		if lineDiff.StartToStart == 0 {
			hasInlineComment = true
		}
	}

	if hasInlineComment {
		f.g.Inline("{ ")
	} else {
		f.g.Line("{")
	}

	f.g.Indent()

	for !f.currentIndexEOF {
		if f.currentIndexChild.Comment != nil {
			f.formatComment()
		}

		if f.currentIndexChild.For != nil {
			f.formatFor()
		}

		if f.currentIndexChild.Param != nil {
			f.formatParam()
		}

		if f.currentIndexChild.Error != nil {
			f.formatError()
		}

		f.loadNextChild()
	}

	f.g.Dedent()
	f.g.Inline("}")

	return f.g
}

func (f *ruleFormatter) formatComment() {
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

func (f *ruleFormatter) formatFor() {
	if f.currentIndexChild.For.IsArray {
		f.LineAndCommentf("for: %s[]", f.currentIndexChild.For.Type)
	} else {
		f.LineAndCommentf("for: %s", f.currentIndexChild.For.Type)
	}
}

func (f *ruleFormatter) formatParam() {
	if f.currentIndexChild.Param.IsArray {
		f.LineAndCommentf("param: %s[]", f.currentIndexChild.Param.Param)
	} else {
		f.LineAndCommentf("param: %s", f.currentIndexChild.Param.Param)
	}
}

func (f *ruleFormatter) formatError() {
	f.LineAndCommentf("error: \"%s\"", strutil.EscapeQuotes(f.currentIndexChild.Error.Error))
}
