package transpile

import (
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// ToURPC transpiles an UFO-RPC JSON schema to it's AST representation.
//
// The resulting AST Schema will not include any imports, extends, external
// docstrings, comments nor comment blocks.
//
// To get the string representation of the AST Schema, you can use the
// formatter package.
func ToURPC(jsonSchema schema.Schema) (ast.Schema, error) {
	return ast.Schema{}, nil
}
