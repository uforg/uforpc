package ast

import (
	"strings"

	"github.com/uforg/uforpc/urpc/internal/util/strutil"
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

// GetRulesMap returns a map of rule names to rule declarations.
func (s *Schema) GetRulesMap() map[string]*RuleDecl {
	rulesMap := make(map[string]*RuleDecl)
	for _, rule := range s.GetRules() {
		rulesMap[rule.Name] = rule
	}
	return rulesMap
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

// GetTypesMap returns a map of type names to type declarations.
func (s *Schema) GetTypesMap() map[string]*TypeDecl {
	typesMap := make(map[string]*TypeDecl)
	for _, typeDecl := range s.GetTypes() {
		typesMap[typeDecl.Name] = typeDecl
	}
	return typesMap
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

// GetProcsMap returns a map of procedure names to procedure declarations.
func (s *Schema) GetProcsMap() map[string]*ProcDecl {
	procsMap := make(map[string]*ProcDecl)
	for _, proc := range s.GetProcs() {
		procsMap[proc.Name] = proc
	}
	return procsMap
}

// GetStreams returns all streams in the URPC schema.
func (s *Schema) GetStreams() []*StreamDecl {
	streams := []*StreamDecl{}
	for _, node := range s.Children {
		if node.Kind() == SchemaChildKindStream {
			streams = append(streams, node.Stream)
		}
	}
	return streams
}

// GetStreamsMap returns a map of stream names to stream declarations.
func (s *Schema) GetStreamsMap() map[string]*StreamDecl {
	streamsMap := make(map[string]*StreamDecl)
	for _, stream := range s.GetStreams() {
		streamsMap[stream.Name] = stream
	}
	return streamsMap
}

// SchemaChildKind represents the kind of a schema child node.
type SchemaChildKind string

const (
	SchemaChildKindVersion   SchemaChildKind = "Version"
	SchemaChildKindComment   SchemaChildKind = "Comment"
	SchemaChildKindDocstring SchemaChildKind = "Docstring"
	SchemaChildKindRule      SchemaChildKind = "Rule"
	SchemaChildKindType      SchemaChildKind = "Type"
	SchemaChildKindProc      SchemaChildKind = "Proc"
	SchemaChildKindStream    SchemaChildKind = "Stream"
)

// SchemaChild represents a child node of the Schema root node.
type SchemaChild struct {
	Positions
	Version   *Version    `parser:"  @@"`
	Comment   *Comment    `parser:"| @@"`
	Rule      *RuleDecl   `parser:"| @@"`
	Type      *TypeDecl   `parser:"| @@"`
	Proc      *ProcDecl   `parser:"| @@"`
	Stream    *StreamDecl `parser:"| @@"`
	Docstring *Docstring  `parser:"| @@"`
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
	if n.Rule != nil {
		return SchemaChildKindRule
	}
	if n.Type != nil {
		return SchemaChildKindType
	}
	if n.Proc != nil {
		return SchemaChildKindProc
	}
	if n.Stream != nil {
		return SchemaChildKindStream
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
	Type    string `parser:"For Colon @(Ident | String | Int | Float | Bool | Datetime)"`
	IsArray bool   `parser:"@(LBracket RBracket)?"`
}

// RuleDeclChildParam represents the "param" clause within a RuleDecl block.
type RuleDeclChildParam struct {
	Positions
	Param   string `parser:"Param Colon @(String | Int | Float | Bool)"`
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
	Docstring  *Docstring               `parser:"(@@ (?! Newline Newline))?"`
	Deprecated *Deprecated              `parser:"(@@ (?= Proc))?"`
	Name       string                   `parser:"Proc @Ident"`
	Children   []*ProcOrStreamDeclChild `parser:"LBrace @@* RBrace"`
}

// StreamDecl represents a stream declaration.
type StreamDecl struct {
	Positions
	Docstring  *Docstring               `parser:"(@@ (?! Newline Newline))?"`
	Deprecated *Deprecated              `parser:"(@@ (?= Stream))?"`
	Name       string                   `parser:"Stream @Ident"`
	Children   []*ProcOrStreamDeclChild `parser:"LBrace @@* RBrace"`
}

// ProcOrStreamDeclChild represents a child node within a ProcDecl or StreamDecl block (Comment, Input, Output, or Meta).
type ProcOrStreamDeclChild struct {
	Positions
	Comment *Comment                     `parser:"  @@"`
	Input   *ProcOrStreamDeclChildInput  `parser:"| @@"`
	Output  *ProcOrStreamDeclChildOutput `parser:"| @@"`
	Meta    *ProcOrStreamDeclChildMeta   `parser:"| @@"`
}

// ProcOrStreamDeclChildInput represents the Input{...} block within a ProcDecl or StreamDecl.
type ProcOrStreamDeclChildInput struct {
	Positions
	Children []*FieldOrComment `parser:"Input LBrace @@* RBrace"`
}

// ProcOrStreamDeclChildOutput represents the Output{...} block within a ProcDecl or StreamDecl.
type ProcOrStreamDeclChildOutput struct {
	Positions
	Children []*FieldOrComment `parser:"Output LBrace @@* RBrace"`
}

// ProcOrStreamDeclChildMeta represents the Meta{...} block within a ProcDecl or StreamDecl.
type ProcOrStreamDeclChildMeta struct {
	Positions
	Children []*ProcOrStreamDeclChildMetaChild `parser:"Meta LBrace @@* RBrace"`
}

// ProcOrStreamDeclChildMetaChild represents a child node within a MetaBlock (either a Comment or a Key-Value pair).
type ProcOrStreamDeclChildMetaChild struct {
	Positions
	Comment *Comment                     `parser:"  @@"`
	KV      *ProcOrStreamDeclChildMetaKV `parser:"| @@"`
}

// ProcOrStreamDeclChildMetaKV represents a key-value pair within a MetaBlock.
type ProcOrStreamDeclChildMetaKV struct {
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

// GetExternal returns a path and a bool indicating if the docstring
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
	Named  *string          `parser:"@(Ident | String | Int | Float | Bool | Datetime)"`
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
	ParamSingle     *AnyLiteral `parser:"@@?"`
	ParamListString []string    `parser:"(LBracket @StringLiteral (Comma @StringLiteral)* RBracket)?"`
	ParamListInt    []string    `parser:"(LBracket @IntLiteral (Comma @IntLiteral)* RBracket)?"`
	ParamListFloat  []string    `parser:"(LBracket @FloatLiteral (Comma @FloatLiteral)* RBracket)?"`
	ParamListBool   []string    `parser:"(LBracket @(TrueLiteral | FalseLiteral) (Comma @(TrueLiteral | FalseLiteral))* RBracket)?"`
	// Error clause, if present, must appear after parameters.
	Error *string `parser:"(Comma? Error Colon @StringLiteral)?"`
}
