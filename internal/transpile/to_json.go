package transpile

import (
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// ToJSON transpiles an UFO-RPC AST schema to it's JSON representation.
//
// The imports, extends and external docstrings of the AST Schema are expected
// to be already resolved, so if there are any unresolved imports or extends, the
// transpiler will ignore them (imports and extends) or threat them as
// literal (docstrings).
//
// All comments and comment blocks will be ignored.
func ToJSON(astSchema ast.Schema) (schema.Schema, error) {
	return schema.Schema{}, nil
}
