package ast

import plexer "github.com/alecthomas/participle/v2/lexer"

// Any node in the AST containing a field Pos lexer.Position
// will be automatically populated from the nearest matching token.
//
// Any node in the AST containing a field EndPos lexer.Position
// will be automatically populated from the token at the end of the node.
//
// https://github.com/alecthomas/participle/blob/master/README.md#error-reporting

// Position is an alias for the participle.Position type.
type Position = plexer.Position

// Positions is a struct that contains the start and end positions of a node.
//
// Used to embed in structs that contain a start and end position and
// automatically populate the Pos field, EndPos field, and the
// GetPositions method.
type Positions struct {
	Pos    Position
	EndPos Position
}

// GetPositions returns the start and end positions of the node.
func (p Positions) GetPositions() Positions {
	return p
}

// WithPositions is an interface that can be implemented by any type
// that has a GetPositions method.
type WithPositions interface {
	GetPositions() Positions
}
