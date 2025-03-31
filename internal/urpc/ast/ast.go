package ast

import plexer "github.com/alecthomas/participle/v2/lexer"

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

type Position plexer.Position

// URPCSchema is the root of the URPC schema AST.
type URPCSchema struct {
	Pos     Position
	EndPos  Position
	Version *Version    `parser:"@@?"`
	Rules   []*RuleDecl `parser:"@@*"`
}

// Version represents the version of the URPC schema.
type Version struct {
	Pos    Position
	EndPos Position
	Number int `parser:"VERSION COLON @INT"`
}

// Docstring represents the documentation for a rule, type or procedure declaration.
type Docstring struct {
	Pos     Position
	EndPos  Position
	Content string `parser:"@DOCSTRING"`
}

// RuleDecl represents a custom validation rule declaration.
type RuleDecl struct {
	Pos       Position
	EndPos    Position
	Docstring *Docstring   `parser:"@@?"`
	Name      string       `parser:"RULE AT @IDENT"`
	Body      RuleDeclBody `parser:"LBRACE @@ RBRACE"`
}

// RuleDeclBody represents the body of a custom validation rule declaration.
type RuleDeclBody struct {
	Pos    Position
	EndPos Position
	For    string `parser:"FOR COLON @IDENT"`
	Param  string `parser:"(PARAM COLON @IDENT)?"`
	Error  string `parser:"(ERROR COLON @STRING)?"`
}
