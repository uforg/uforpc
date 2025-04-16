package transpile

import (
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// ToJSON transpiles an UFO-RPC AST schema to it's JSON representation.
//
// The imports, extends and external docstrings of the AST Schema are expected
// to be already resolved.
//
// If there are any unresolved imports or extends, the transpiler
// will ignore them.
//
// If there are any unresolved external docstrings, the transpiler will
// treat them literally as strings.
//
// All comments and comment blocks will be ignored.
//
// To get the string representation of the JSON Schema, you can use the
// json.Marshal function.
func ToJSON(astSchema ast.Schema) (schema.Schema, error) {
	return schema.Schema{}, nil
}
