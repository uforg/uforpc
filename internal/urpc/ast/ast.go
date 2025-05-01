package ast

import (
	"strings"

	"github.com/uforg/uforpc/internal/util/strutil"
)

// This AST is used for parsing the URPC schema and it uses the
// participle library for parsing.
//
// It includes embedded Positions fields for each node to track the
// position of the node in the original source code, it is used
// later in the analyzer and LSP to give useful error messages
// and auto-completion. Those positions are automatically populated
// by the participle library.

// Schema is the root of the URPC schema AST.
type Schema struct {
	Positions
	Children []*SchemaChild `parser:"@@*"`
}

// GetVersions returns all version declarations in the URPC schema.
func (s *Schema) GetVersions() []*Version {
	versions := []*Version{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindVersion {
			versions = append(versions, node.Version)
		}
	}
	return versions
}

// GetComments returns all comments in the URPC schema.
func (s *Schema) GetComments() []*Comment {
	comments := []*Comment{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindComment {
			comments = append(comments, node.Comment)
		}
	}
	return comments
}

// GetDocstrings returns all docstrings in the URPC schema.
func (s *Schema) GetDocstrings() []*Docstring {
	docstrings := []*Docstring{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindDocstring {
			docstrings = append(docstrings, node.Docstring)
		}
	}
	return docstrings
}

// GetImports returns all import statements in the URPC schema.
func (s *Schema) GetImports() []*Import {
	imports := []*Import{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindImport {
			imports = append(imports, node.Import)
		}
	}
	return imports
}

// GetRules returns all custom validation rules in the URPC schema.
func (s *Schema) GetRules() []*RuleDecl {
	rules := []*RuleDecl{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindRule {
			rules = append(rules, node.Rule)
		}
	}
	return rules
}

// GetTypes returns all custom types in the URPC schema.
func (s *Schema) GetTypes() []*TypeDecl {
	types := []*TypeDecl{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindType {
			types = append(types, node.Type)
		}
	}
	return types
}

// GetProcs returns all procedures in the URPC schema.
func (s *Schema) GetProcs() []*ProcDecl {
	procs := []*ProcDecl{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindProc {
			procs = append(procs, node.Proc)
		}
	}
	return procs
}

// SchemaChildKind represents the kind of a schema child node.
type SchemaChildKind string

const (
	SchemaChildKindVersion   SchemaChildKind = "Version"
	SchemaChildKindComment   SchemaChildKind = "Comment"
	SchemaChildKindDocstring SchemaChildKind = "Docstring"
	SchemaChildKindImport    SchemaChildKind = "Import"
	SchemaChildKindRule      SchemaChildKind = "Rule"
	SchemaChildKindType      SchemaChildKind = "Type"
	SchemaChildKindProc      SchemaChildKind = "Proc"
)

// SchemaChild represents a child node of the Schema root node.
type SchemaChild struct {
	Positions
	Version   *Version   `parser:"  @@"`
	Comment   *Comment   `parser:"| @@"`
	Import    *Import    `parser:"| @@"`
	Rule      *RuleDecl  `parser:"| @@"`
	Type      *TypeDecl  `parser:"| @@"`
	Proc      *ProcDecl  `parser:"| @@"`
	Docstring *Docstring `parser:"| @@"`
}

func (n *SchemaChild) Kind() SchemaChildKind {
	if n.Version != nil {
		return SchemaChildKindVersion
	}
	if n.Comment != nil {
		return SchemaChildKindComment
	}
	if n.Docstring != nil {
		return SchemaChildKindDocstring
	}
	if n.Import != nil {
		return SchemaChildKindImport
	}
	if n.Rule != nil {
		return SchemaChildKindRule
	}
	if n.Type != nil {
		return SchemaChildKindType
	}
	if n.Proc != nil {
		return SchemaChildKindProc
	}
	return ""
}

// Version represents the version of the URPC schema.
type Version struct {
	Positions
	Number int `parser:"Version @IntLiteral"`
}

// Comment represents both simple and block comments in the URPC schema.
type Comment struct {
	Positions
	Simple *string `parser:"  @Comment"`
	Block  *string `parser:"| @CommentBlock"`
}

// Import represents an import statement.
type Import struct {
	Positions
	Path string `parser:"Import @StringLiteral"`
}

// RuleDecl represents a custom validation rule declaration.
type RuleDecl struct {
	Positions
	Docstring  *Docstring       `parser:"(@@ (?! Newline Newline))?"`
	Deprecated *Deprecated      `parser:"(@@ (?= Rule))?"`
	Name       string           `parser:"Rule At @Ident"`
	Children   []*RuleDeclChild `parser:"LBrace @@* RBrace"`
}

// RuleDeclChild represents a child node within a RuleDecl block.
type RuleDeclChild struct {
	Positions
	Comment *Comment            `parser:"  @@"`
	For     *RuleDeclChildFor   `parser:"| @@"`
	Param   *RuleDeclChildParam `parser:"| @@"`
	Error   *RuleDeclChildError `parser:"| @@"`
}

// RuleDeclChildFor represents the "for" clause within a RuleDecl block.
type RuleDeclChildFor struct {
	Positions
	For     string `parser:"For Colon @(Ident | String | Int | Float | Boolean | Datetime)"`
	IsArray bool   `parser:"@(LBracket RBracket)?"`
}

// RuleDeclChildParam represents the "param" clause within a RuleDecl block.
type RuleDeclChildParam struct {
	Positions
	Param   string `parser:"Param Colon @(String | Int | Float | Boolean)"`
	IsArray bool   `parser:"@(LBracket RBracket)?"`
}

// RuleDeclChildError represents the "error" clause within a RuleDecl block.
type RuleDeclChildError struct {
	Positions
	Error string `parser:"Error Colon @StringLiteral"`
}

// TypeDecl represents a custom type declaration.
type TypeDecl struct {
	Positions
	Docstring  *Docstring        `parser:"(@@ (?! Newline Newline))?"`
	Deprecated *Deprecated       `parser:"(@@ (?= Type))?"`
	Name       string            `parser:"Type @Ident"`
	Children   []*FieldOrComment `parser:"LBrace @@* RBrace"`
}

// ProcDecl represents a procedure declaration.
type ProcDecl struct {
	Positions
	Docstring  *Docstring       `parser:"(@@ (?! Newline Newline))?"`
	Deprecated *Deprecated      `parser:"(@@ (?= Proc))?"`
	Name       string           `parser:"Proc @Ident"`
	Children   []*ProcDeclChild `parser:"LBrace @@* RBrace"`
}

// ProcDeclChild represents a child node within a ProcDecl block (Comment, Input, Output, or Meta).
type ProcDeclChild struct {
	Positions
	Comment *Comment             `parser:"  @@"`
	Input   *ProcDeclChildInput  `parser:"| @@"`
	Output  *ProcDeclChildOutput `parser:"| @@"`
	Meta    *ProcDeclChildMeta   `parser:"| @@"`
}

// ProcDeclChildInput represents the Input{...} block within a ProcDecl.
type ProcDeclChildInput struct {
	Positions
	Children []*FieldOrComment `parser:"Input LBrace @@* RBrace"`
}

// ProcDeclChildOutput represents the Output{...} block within a ProcDecl.
type ProcDeclChildOutput struct {
	Positions
	Children []*FieldOrComment `parser:"Output LBrace @@* RBrace"`
}

// ProcDeclChildMeta represents the Meta{...} block within a ProcDecl.
type ProcDeclChildMeta struct {
	Positions
	Children []*ProcDeclChildMetaChild `parser:"Meta LBrace @@* RBrace"`
}

// ProcDeclChildMetaChild represents a child node within a MetaBlock (either a Comment or a Key-Value pair).
type ProcDeclChildMetaChild struct {
	Positions
	Comment *Comment             `parser:"  @@"`
	KV      *ProcDeclChildMetaKV `parser:"| @@"`
}

// ProcDeclChildMetaKV represents a key-value pair within a MetaBlock.
type ProcDeclChildMetaKV struct {
	Positions
	Key   string     `parser:"@Ident"`
	Value AnyLiteral `parser:"Colon @@"`
}

//////////////////
// SHARED TYPES //
//////////////////

// Docstring represents a docstring in the URPC schema.
type Docstring struct {
	Positions
	Value string `parser:"@Docstring"`
}

// GetExternal returns a path and a boolean indicating if the docstring
// references an external Markdown file.
func (d Docstring) GetExternal() (string, bool) {
	trimmed := strings.TrimSpace(d.Value)
	if strings.ContainsAny(trimmed, "\r\n") {
		return "", false
	}

	if strings.TrimSuffix(".md", trimmed) == "" {
		return "", false
	}

	if !strings.HasSuffix(trimmed, ".md") {
		return "", false
	}

	return trimmed, true
}

// Deprecated represents a deprecated declaration.
type Deprecated struct {
	Positions
	Message *string `parser:"Deprecated (LParen @StringLiteral RParen)?"`
}

// AnyLiteral represents any of the built-in literal types.
type AnyLiteral struct {
	Positions
	Str   *string `parser:"  @StringLiteral"`
	Int   *string `parser:"| @IntLiteral"`
	Float *string `parser:"| @FloatLiteral"`
	True  *string `parser:"| @TrueLiteral"`
	False *string `parser:"| @FalseLiteral"`
}

// String returns the string representation of the value of the literal.
func (al AnyLiteral) String() string {
	if al.Str != nil {
		return `"` + strutil.EscapeQuotes(*al.Str) + `"`
	}
	if al.Int != nil {
		return *al.Int
	}
	if al.Float != nil {
		return *al.Float
	}
	if al.True != nil {
		return "true"
	}
	if al.False != nil {
		return "false"
	}
	return ""
}

// FieldOrComment represents a child node within blocks that contain fields,
// such as TypeDecl, ProcDeclChildInput, ProcDeclChildOutput, and FieldTypeObject.
type FieldOrComment struct {
	Positions
	Comment *Comment `parser:"  @@"`
	Field   *Field   `parser:"| @@"`
}

// Field represents a field definition. It can contain comments and rules after the type definition.
type Field struct {
	Positions
	Name     string        `parser:"@Ident"`
	Optional bool          `parser:"@(Question)?"`
	Type     FieldType     `parser:"Colon @@"`
	Children []*FieldChild `parser:"@@*"` // Captures rules and comments following the type
}

// FieldChild represents a child node following a Field's type definition (either a Comment or a FieldRule).
type FieldChild struct {
	Positions
	// Field comments are only captured if they are followed by a rule
	// otherwise they are captured by the parent FieldOrComment
	Comment *Comment   `parser:"  @@ (?= At Ident )"`
	Rule    *FieldRule `parser:"| @@"`
}

// FieldType represents the type of a field.
type FieldType struct {
	Positions
	Base    *FieldTypeBase `parser:"@@"`
	IsArray bool           `parser:"@(LBracket RBracket)?"`
}

// FieldTypeBase represents the base type of a field (primitive, named, or inline object).
type FieldTypeBase struct {
	Positions
	Named  *string          `parser:"@(Ident | String | Int | Float | Boolean | Datetime)"`
	Object *FieldTypeObject `parser:"| @@"`
}

// FieldTypeObject represents an inline object type definition.
type FieldTypeObject struct {
	Positions
	Children []*FieldOrComment `parser:"LBrace @@* RBrace"`
}

// FieldRule represents a validation rule applied to a field. It does not support comments within its body.
type FieldRule struct {
	Positions
	Name string         `parser:"At @Ident"`
	Body *FieldRuleBody `parser:"(LParen @@ RParen)?"` // Body is optional
}

// FieldRuleBody represents the body of a validation rule applied to a field.
// It captures parameters and an optional error message. Comments are not supported within this structure.
type FieldRuleBody struct {
	Positions
	// Parameters are captured positionally; validation must ensure correct number/type.
	ParamSingle      *AnyLiteral `parser:"@@?"`
	ParamListString  []string    `parser:"(LBracket @StringLiteral (Comma @StringLiteral)* RBracket)?"`
	ParamListInt     []string    `parser:"(LBracket @IntLiteral (Comma @IntLiteral)* RBracket)?"`
	ParamListFloat   []string    `parser:"(LBracket @FloatLiteral (Comma @FloatLiteral)* RBracket)?"`
	ParamListBoolean []string    `parser:"(LBracket @(TrueLiteral | FalseLiteral) (Comma @(TrueLiteral | FalseLiteral))* RBracket)?"`
	// Error clause, if present, must appear after parameters.
	Error *string `parser:"(Comma? Error Colon @StringLiteral)?"`
}
