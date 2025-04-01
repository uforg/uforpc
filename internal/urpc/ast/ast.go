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
	Procs   []*ProcDecl `parser:"@@*"`
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

// ProcDecl represents a procedure declaration.
type ProcDecl struct {
	Pos       Position
	EndPos    Position
	Docstring string       `parser:"@Docstring?"`
	Name      string       `parser:"Proc @Ident"`
	Body      ProcDeclBody `parser:"LBrace @@? RBrace"`
}

// ProcDeclBody represents the body of a procedure declaration.
type ProcDeclBody struct {
	Input  []*Field              `parser:"(Input LBrace @@+ RBrace)?"`
	Output []*Field              `parser:"(Output LBrace @@+ RBrace)?"`
	Meta   []*ProcDeclBodyMetaKV `parser:"(Meta LBrace @@+ RBrace)?"`
}

// ProcDeclBodyMetaKV represents a key-value pair in the meta information of a procedure declaration.
type ProcDeclBodyMetaKV struct {
	Key   string `parser:"@Ident"`
	Value string `parser:"Colon @(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral)"`
}

//////////////////
// SHARED TYPES //
//////////////////

// Field represents a field in a custom type or procedure input/output.
type Field struct {
	Name     string       `parser:"@Ident"`
	Optional bool         `parser:"@(Question)?"`
	Type     FieldType    `parser:"Colon @@"`
	Rules    []*FieldRule `parser:"@@*"`
}

// FieldType represents the type of a field. If the field is an array, the Depth
// represents the depth of the array otherwise it is 0.
type FieldType struct {
	Base  *FieldTypeBase `parser:"@@"`
	Depth FieldTypeDepth `parser:"@((LBracket RBracket)*)"`
}

// FieldTypeDepth represents the depth of an array.
type FieldTypeDepth int

func (a *FieldTypeDepth) Capture(values []string) error {
	count := 0
	for i := range len(values) {
		if values[i] == "[" && values[i+1] == "]" {
			count++
		}
	}

	*a = FieldTypeDepth(count)
	return nil
}

// FieldTypeBase represents the base type of a field. If the field is a primitive
// or custom type, the Named field is set. If the field is an inline object, the Object
// field is set.
type FieldTypeBase struct {
	Named  *string          `parser:"@(Ident | String | Int | Float | Boolean | Datetime)"`
	Object *FieldTypeObject `parser:"| @@"`
}

// FieldTypeObject represents an inline object type.
type FieldTypeObject struct {
	Fields []*Field `parser:"LBrace @@+ RBrace"`
}

// FieldRule represents a rule applied to a field.
type FieldRule struct {
	Name string        `parser:"At @Ident"`
	Body FieldRuleBody `parser:"(LParen @@ RParen)?"`
}

// FieldRuleBody represents the body of a rule applied to a field.
type FieldRuleBody struct {
	ParamSingle *string  `parser:"@(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral)?"`
	ParamList   []string `parser:"(LBracket @(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral) (Comma @(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral))* RBracket)?"`
	Error       string   `parser:"(Comma? Error Colon @StringLiteral)?"`
}
