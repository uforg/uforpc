package lsp

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/urpc/analyzer"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// RequestMessageTextDocumentHover represents a request for hover information.
type RequestMessageTextDocumentHover struct {
	RequestMessage
	Params RequestMessageTextDocumentHoverParams `json:"params"`
}

// RequestMessageTextDocumentHoverParams represents the parameters for a hover request.
type RequestMessageTextDocumentHoverParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The position inside the text document.
	Position TextDocumentPosition `json:"position"`
}

// ResponseMessageTextDocumentHover represents a response to a hover request.
type ResponseMessageTextDocumentHover struct {
	ResponseMessage
	// The result of the request.
	Result *HoverResult `json:"result"`
}

// HoverResult represents the result of a hover request.
type HoverResult struct {
	// The hover's content.
	Contents MarkupContent `json:"contents"`
	// An optional range that is used to visualize the hover.
	Range *TextDocumentRange `json:"range,omitempty"`
}

// MarkupContent represents a hover content with a specific kind of markup.
type MarkupContent struct {
	// The type of the markup content. Currently only "markdown" is supported.
	Kind string `json:"kind"`
	// The content itself.
	Value string `json:"value"`
}

// handleTextDocumentHover handles a textDocument/hover request.
func (l *LSP) handleTextDocumentHover(rawMessage []byte) (any, error) {
	var request RequestMessageTextDocumentHover
	if err := decode(rawMessage, &request); err != nil {
		return nil, fmt.Errorf("failed to decode hover request: %w", err)
	}

	filePath := request.Params.TextDocument.URI
	position := request.Params.Position

	// Get the document content
	content, _, found, err := l.docstore.GetInMemory("", filePath)
	if !found {
		return nil, fmt.Errorf("text document not found in memory: %s", filePath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get in memory text document: %w", err)
	}

	// Run the analyzer to get the combined schema
	combinedSchema, _, err := l.analyzer.Analyze(filePath)
	if err != nil {
		l.logger.Error("failed to analyze document", "uri", filePath, "error", err)
	}

	// Convert LSP position (0-based) to AST position (1-based)
	astPosition := ast.Position{
		Filename: filePath,
		Line:     position.Line + 1,
		Column:   position.Character + 1,
	}

	// Find the hover information
	hoverResult := l.findHoverInfo(content, astPosition, combinedSchema)

	response := ResponseMessageTextDocumentHover{
		ResponseMessage: ResponseMessage{
			Message: DefaultMessage,
			ID:      request.ID,
		},
		Result: hoverResult,
	}

	return response, nil
}

// findHoverInfo finds hover information for a symbol at the given position.
func (l *LSP) findHoverInfo(content string, position ast.Position, combinedSchema analyzer.CombinedSchema) *HoverResult {
	// Find the token at the position
	token, err := findTokenAtPosition(content, position)
	if err != nil {
		l.logger.Error("failed to find token at position", "position", position, "error", err)
		return nil
	}

	// Check if the token is a reference to a type
	if hoverInfo := l.findTypeHoverInfo(token, combinedSchema); hoverInfo != nil {
		return hoverInfo
	}

	// Check if the token is a reference to a rule
	if hoverInfo := l.findRuleHoverInfo(token, combinedSchema); hoverInfo != nil {
		return hoverInfo
	}

	return nil
}

// findTypeHoverInfo finds hover information for a type.
func (l *LSP) findTypeHoverInfo(token string, combinedSchema analyzer.CombinedSchema) *HoverResult {
	// Check if the token is a type name
	typeDecl, exists := combinedSchema.TypeDecls[token]
	if !exists {
		return nil
	}

	// Get the source code of the type definition
	sourceCode, err := l.getTypeSourceCode(typeDecl)
	if err != nil {
		l.logger.Error("failed to get type source code", "type", token, "error", err)
		return nil
	}

	// Create a hover result with the source code
	return &HoverResult{
		Contents: MarkupContent{
			Kind:  "markdown",
			Value: fmt.Sprintf("```urpc\n%s\n```", sourceCode),
		},
	}
}

// findRuleHoverInfo finds hover information for a rule.
func (l *LSP) findRuleHoverInfo(token string, combinedSchema analyzer.CombinedSchema) *HoverResult {
	// If the token starts with @, remove it
	if len(token) > 0 && token[0] == '@' {
		token = token[1:]
	}

	// Check if the token is a rule name
	ruleDecl, exists := combinedSchema.RuleDecls[token]
	if !exists {
		return nil
	}

	// Get the source code of the rule definition
	sourceCode, err := l.getRuleSourceCode(ruleDecl)
	if err != nil {
		l.logger.Error("failed to get rule source code", "rule", token, "error", err)
		return nil
	}

	// Create a hover result with the source code
	return &HoverResult{
		Contents: MarkupContent{
			Kind:  "markdown",
			Value: fmt.Sprintf("```urpc\n%s\n```", sourceCode),
		},
	}
}

// getTypeSourceCode extracts the source code of a type definition.
func (l *LSP) getTypeSourceCode(typeDecl *ast.TypeDecl) (string, error) {
	// Get the file content
	content, _, found, err := l.docstore.GetInMemory("", typeDecl.Pos.Filename)
	if !found {
		// Try to get from disk if not in memory
		content, _, found, err = l.docstore.GetFromDisk("", typeDecl.Pos.Filename)
		if !found || err != nil {
			return "", fmt.Errorf("failed to get file content: %w", err)
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to get file content: %w", err)
	}

	// Extract the type definition from the content
	return extractCodeFromContent(content, typeDecl.Pos.Line, typeDecl.EndPos.Line)
}

// getRuleSourceCode extracts the source code of a rule definition.
func (l *LSP) getRuleSourceCode(ruleDecl *ast.RuleDecl) (string, error) {
	// Get the file content
	content, _, found, err := l.docstore.GetInMemory("", ruleDecl.Pos.Filename)
	if !found {
		// Try to get from disk if not in memory
		content, _, found, err = l.docstore.GetFromDisk("", ruleDecl.Pos.Filename)
		if !found || err != nil {
			return "", fmt.Errorf("failed to get file content: %w", err)
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to get file content: %w", err)
	}

	// Extract the rule definition from the content
	return extractCodeFromContent(content, ruleDecl.Pos.Line, ruleDecl.EndPos.Line)
}

// extractCodeFromContent extracts a range of lines from the content.
func extractCodeFromContent(content string, startLine, endLine int) (string, error) {
	lines := splitLines(content)

	if startLine <= 0 || startLine > len(lines) {
		return "", fmt.Errorf("start line out of range: %d", startLine)
	}

	if endLine <= 0 || endLine > len(lines) {
		return "", fmt.Errorf("end line out of range: %d", endLine)
	}

	// Extract the lines
	extractedLines := lines[startLine-1 : endLine]

	// Find the minimum indentation
	minIndent := -1
	for _, line := range extractedLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	// Remove the minimum indentation from each line
	if minIndent > 0 {
		for i, line := range extractedLines {
			if len(line) >= minIndent {
				extractedLines[i] = line[minIndent:]
			}
		}
	}

	// Join the lines
	return strings.Join(extractedLines, "\n"), nil
}
