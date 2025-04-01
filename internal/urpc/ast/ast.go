package ast

import (
	plexer "github.com/alecthomas/participle/v2/lexer"
)

// This AST is used for parsing the URPC schema and it uses the
// participle library for parsing.
//
// It includes Pos and EndPos fields for each node to track the
// position of the node in the original source code.
//
// Any node in the AST containing a field Pos lexer.Position
// will be automatically populated from the nearest matching token.
//
// Any node in the AST containing a field EndPos lexer.Position
// will be automatically populated from the token at the end of the node.
//
// https://github.com/alecthomas/participle/blob/master/README.md#error-reporting

// The Pos and EndPos fields can be added only to nodes that will be referenced from
// other nodes to use them later in the analyzer and LSP and give useful error
// messages and auto-completion.
//
// Other inner nodes can skip the Pos and EndPos fields to make the AST more
// compact and easier to work with.

type Position plexer.Position

// URPCSchema is the root of the URPC schema AST.
type URPCSchema struct {
	Version *Version    `parser:"@@?"`
	Imports []*Import   `parser:"@@*"`
	Rules   []*RuleDecl `parser:"@@*"`
	Types   []*TypeDecl `parser:"@@*"`
}

// Version represents the version of the URPC schema.
type Version struct {
	Pos    Position
	EndPos Position
	Number int `parser:"Version @IntLiteral"`
}

// Import represents an import statement.
type Import struct {
	Pos    Position
	EndPos Position
	Path   string `parser:"Import @StringLiteral"`
}

// RuleDecl represents a custom validation rule declaration.
type RuleDecl struct {
	Pos       Position
	EndPos    Position
	Docstring string       `parser:"@Docstring?"`
	Name      string       `parser:"Rule At @Ident"`
	Body      RuleDeclBody `parser:"LBrace @@ RBrace"`
}

// RuleDeclBody represents the body of a custom validation rule declaration.
type RuleDeclBody struct {
	For          string `parser:"For Colon @(Ident | String | Int | Float | Boolean | Datetime)"`
	Param        string `parser:"(Param Colon @(String | Int | Float | Boolean))?"`
	ParamIsArray bool   `parser:"@(LBracket RBracket)?"`
	Error        string `parser:"(Error Colon @StringLiteral)?"`
}

// TypeDecl represents a custom type declaration.
type TypeDecl struct {
	Pos       Position
	EndPos    Position
	Docstring string   `parser:"@Docstring?"`
	Name      string   `parser:"Type @Ident"`
	Extends   []string `parser:"(Extends @Ident (Comma @Ident)*)?"`
	Fields    []*Field `parser:"LBrace @@+ RBrace"`
}
