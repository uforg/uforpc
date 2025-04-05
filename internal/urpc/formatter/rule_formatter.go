package formatter

import (
	"fmt"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

type ruleFormatter struct {
	g                 *genkit.GenKit
	rule              *ast.RuleDecl
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
		rule:              rule,
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
		currentIndexChild = *f.rule.Children[currentIndex]
	}

	f.currentIndex = currentIndex
	f.currentIndexEOF = currentIndexEOF
	f.currentIndexChild = currentIndexChild
}

// peekChild returns information about the child at the current index +- offset.
//
// Returns:
//   - The child at the current index +- offset.
//   - The line difference between the current child and the peeked child.
//   - A boolean indicating if the peeked child is out of bounds (EOL).
func (f *ruleFormatter) peekChild(offset int) (ast.RuleDeclChild, int, bool) {
	peekIndex := f.currentIndex + offset
	peekIndexEOF := peekIndex < 0 || peekIndex > f.maxIndex
	peekIndexChild := ast.RuleDeclChild{}
	lineDiff := 0

	if !peekIndexEOF {
		peekIndexChild = *f.rule.Children[peekIndex]
		lineDiff = peekIndexChild.Pos.Line - f.currentIndexChild.Pos.Line
	}

	return peekIndexChild, lineDiff, peekIndexEOF
}

// LineAndComment writes a line of content to the formatter. It also handles inline comments.
func (f *ruleFormatter) LineAndComment(content string) {
	next, nextLineDiff, nextEOF := f.peekChild(1)

	// If next is an inline comment
	if !nextEOF && next.Comment != nil && nextLineDiff == 0 {
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
	if f.rule.Docstring != "" {
		f.g.Linef(`"""%s"""`, f.rule.Docstring)
	}
	f.g.Inlinef(`rule @%s `, f.rule.Name)

	if f.currentIndexEOF {
		f.g.Inline("{}")
		return f.g
	}

	hasInlineComment := false
	if f.currentIndexChild.Comment != nil {
		lineDiff := f.currentIndexChild.Pos.Line - f.rule.Pos.Line
		if lineDiff == 0 {
			hasInlineComment = true
		}
	}

	if hasInlineComment {
		f.g.Inline("{ ")
	} else {
		f.g.Line("{")
	}

	f.g.Indent()

	for {
		if f.currentIndexEOF {
			break
		}

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
		if prevLineDiff < -1 {
			shouldBreakBefore = true
		}
	}

	if shouldBreakBefore {
		f.g.Break()
	}

	if f.currentIndexChild.Comment.Simple != nil {
		f.LineAndCommentf("//%s", *f.currentIndexChild.Comment.Simple)
	}

	if f.currentIndexChild.Comment.Block != nil {
		f.LineAndCommentf("/*%s*/", *f.currentIndexChild.Comment.Block)
	}
}

func (f *ruleFormatter) formatFor() {
	f.LineAndCommentf("for: %s", f.currentIndexChild.For.For)
}

func (f *ruleFormatter) formatParam() {
	if f.currentIndexChild.Param.IsArray {
		f.LineAndCommentf("param: %s[]", f.currentIndexChild.Param.Param)
	} else {
		f.LineAndCommentf("param: %s", f.currentIndexChild.Param.Param)
	}
}

func (f *ruleFormatter) formatError() {
	f.LineAndCommentf("error: \"%s\"", f.currentIndexChild.Error.Error)
}
