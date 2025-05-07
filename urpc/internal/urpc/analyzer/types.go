package analyzer

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
)

// FileProvider interface provides a way to get files to the analyzer.
type FileProvider interface {
	// GetFileAndHash returns the content of the file and its sha256 hash.
	//
	//	- relativeTo: is the URI/path of which the path is relative to.
	//	- path: can be relative if relativeTo is provided. Otherwise it must be absolute.
	//
	// Returns:
	//	- content: content of the file
	//  - hash: sha256 hash of the content in hex format
	// 	- error: if the file cannot be found or read. Should wrap os.ErrNotExist if not found.
	GetFileAndHash(relativeTo string, path string) (string, string, error)
}

// Positions represents the start and end positions of a source code range.
//
// Both Pos and EndPos are ast.Position, so they contain the file path, line and column.
type Positions struct {
	Pos    ast.Position
	EndPos ast.Position
}

// Diagnostic represents an error or warning found during analysis.
type Diagnostic struct {
	Positions        // The range of the source code where the diagnostic occurred.
	Message   string // The diagnostic message.
}

// String implements fmt.Stringer interface.
func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: %s", d.Pos.String(), d.Message)
}

// Error implements error interface.
func (d Diagnostic) Error() string {
	return d.String()
}

// CombinedSchema is the result of the process of combining multiple URPC schemas
// based on an entry point and it's recursive import checks.
type CombinedSchema struct {
	Schema    *ast.Schema              // The combined AST of all schemas.
	RuleDecls map[string]*ast.RuleDecl // Map of all rule names -> ast.RuleDecl from the combined Schema.
	TypeDecls map[string]*ast.TypeDecl // Map of all type names -> ast.TypeDecl from the combined Schema.
	ProcDecls map[string]*ast.ProcDecl // Map of all proc names -> ast.ProcDecl from the combined Schema.
}
