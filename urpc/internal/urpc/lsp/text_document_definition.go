package lsp

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/urpc/internal/urpc/analyzer"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
)

// RequestMessageTextDocumentDefinition represents a request for the definition of a symbol.
type RequestMessageTextDocumentDefinition struct {
	RequestMessage
	Params RequestMessageTextDocumentDefinitionParams `json:"params"`
}

// RequestMessageTextDocumentDefinitionParams represents the parameters for a definition request.
type RequestMessageTextDocumentDefinitionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The position inside the text document.
	Position TextDocumentPosition `json:"position"`
}

// ResponseMessageTextDocumentDefinition represents a response to a definition request.
type ResponseMessageTextDocumentDefinition struct {
	ResponseMessage
	// The result of the request. Can be a single location or an array of locations.
	Result []Location `json:"result"`
}

// Location represents a location inside a resource, such as a line inside a text file.
type Location struct {
	// The URI of the document.
	URI string `json:"uri"`
	// The range inside the document.
	Range TextDocumentRange `json:"range"`
}

// handleTextDocumentDefinition handles a textDocument/definition request.
func (l *LSP) handleTextDocumentDefinition(rawMessage []byte) (any, error) {
	var request RequestMessageTextDocumentDefinition
	if err := decode(rawMessage, &request); err != nil {
		return nil, fmt.Errorf("failed to decode definition request: %w", err)
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

	// Find the definition
	locations := l.findDefinition(content, astPosition, combinedSchema)

	response := ResponseMessageTextDocumentDefinition{
		ResponseMessage: ResponseMessage{
			Message: DefaultMessage,
			ID:      request.ID,
		},
		Result: locations,
	}

	return response, nil
}

// findDefinition finds the definition of a symbol at the given position.
func (l *LSP) findDefinition(content string, position ast.Position, combinedSchema analyzer.CombinedSchema) []Location {
	// We don't need to parse the document here since we're using the token finder
	// to extract the token at the position and then look it up in the combinedSchema

	// Find the token at the position
	token, err := findTokenAtPosition(content, position)
	if err != nil {
		l.logger.Error("failed to find token at position", "position", position, "error", err)
		return nil
	}

	// Check if the token is a reference to a type
	if location := findTypeDefinition(token, position, combinedSchema); location != nil {
		return []Location{*location}
	}

	// Check if the token is a reference to a rule
	if location := findRuleDefinition(token, position, combinedSchema); location != nil {
		return []Location{*location}
	}

	return nil
}

// findTokenAtPosition finds the token at the given position in the content.
func findTokenAtPosition(content string, position ast.Position) (string, error) {
	// This is a simple implementation that extracts the word at the position
	// A more robust implementation would use the lexer to get the token

	// Convert content to lines
	lines := splitLines(content)
	if position.Line <= 0 || position.Line > len(lines) {
		return "", fmt.Errorf("line out of range: %d", position.Line)
	}

	line := lines[position.Line-1]
	if position.Column <= 0 || position.Column > len(line)+1 {
		return "", fmt.Errorf("column out of range: %d", position.Column)
	}

	// Find the start and end of the word
	start := position.Column - 1
	for start > 0 && isIdentifierChar(line[start-1]) {
		start--
	}

	end := position.Column - 1
	for end < len(line) && isIdentifierChar(line[end]) {
		end++
	}

	if start == end {
		// No token at position (e.g., whitespace or punctuation)
		return "", fmt.Errorf("no token at position")
	}

	return line[start:end], nil
}

// isIdentifierChar returns true if the character is valid in an identifier.
func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// splitLines splits the content into lines.
func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}

	var lines []string
	var line string
	for _, c := range content {
		if c == '\n' {
			lines = append(lines, line)
			line = ""
		} else if c != '\r' {
			line += string(c)
		}
	}
	// Always append the last line, even if it's empty
	lines = append(lines, line)
	return lines
}

// findTypeDefinition finds the definition of a type.
func findTypeDefinition(token string, _ ast.Position, combinedSchema analyzer.CombinedSchema) *Location {
	// Check if the token is a type name
	typeDecl, exists := combinedSchema.TypeDecls[token]
	if !exists {
		return nil
	}

	// Create a location for the type definition
	// Ensure the URI has the file:// prefix
	uri := typeDecl.Pos.Filename
	if !strings.HasPrefix(uri, "file://") {
		uri = "file://" + uri
	}

	return &Location{
		URI: uri,
		Range: TextDocumentRange{
			Start: convertASTPositionToLSPPosition(typeDecl.Pos),
			End:   convertASTPositionToLSPPosition(typeDecl.EndPos),
		},
	}
}

// findRuleDefinition finds the definition of a rule.
func findRuleDefinition(token string, _ ast.Position, combinedSchema analyzer.CombinedSchema) *Location {
	// If the token starts with @, remove it
	if len(token) > 0 && token[0] == '@' {
		token = token[1:]
	}

	// Check if the token is a rule name
	ruleDecl, exists := combinedSchema.RuleDecls[token]
	if !exists {
		return nil
	}

	// Create a location for the rule definition
	// Ensure the URI has the file:// prefix
	uri := ruleDecl.Pos.Filename
	if !strings.HasPrefix(uri, "file://") {
		uri = "file://" + uri
	}

	return &Location{
		URI: uri,
		Range: TextDocumentRange{
			Start: convertASTPositionToLSPPosition(ruleDecl.Pos),
			End:   convertASTPositionToLSPPosition(ruleDecl.EndPos),
		},
	}
}
