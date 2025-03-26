package ast

type NodeType int

const (
	_ NodeType = iota
	NodeTypeSchema
	NodeTypeVersion
	NodeTypeTypeDeclaration
	NodeTypeProcDeclaration
	NodeTypeField
	NodeTypeValidationRule
	NodeTypeInput
	NodeTypeOutput
	NodeTypeMetadata
	NodeTypePrimitiveType
	NodeTypeArrayType
	NodeTypeInlineObjectType
	NodeTypeTypeReference
)

type Node interface {
	NodeType() NodeType
}

// Root node
type Schema struct {
	Version    Version
	Types      []TypeDeclaration
	Procedures []ProcDeclaration
}

func (s *Schema) NodeType() NodeType { return NodeTypeSchema }

// Version
type Version struct {
	IsSet bool
	Value int
}

func (v *Version) NodeType() NodeType { return NodeTypeVersion }

// Type system
type TypeDeclaration struct {
	Name   string
	Doc    string
	Fields []Field
}

func (t *TypeDeclaration) NodeType() NodeType { return NodeTypeTypeDeclaration }

type Field struct {
	Name        string
	Type        Type
	Optional    bool
	Validations []ValidationRule
}

func (f *Field) NodeType() NodeType { return NodeTypeField }

type ValidationRule struct {
	Name     string
	Params   []any
	ErrorMsg string
}

func (v *ValidationRule) NodeType() NodeType { return NodeTypeValidationRule }

// Procedures
type ProcDeclaration struct {
	Name     string
	Doc      string
	Input    Input
	Output   Output
	Metadata Metadata
}

func (p *ProcDeclaration) NodeType() NodeType { return NodeTypeProcDeclaration }

type Input struct {
	Fields []Field
}

func (i *Input) NodeType() NodeType { return NodeTypeInput }

type Output struct {
	Fields []Field
}

func (o *Output) NodeType() NodeType { return NodeTypeOutput }

// Metadata
type Metadata struct {
	Entries []KeyValue
}

func (m *Metadata) NodeType() NodeType { return NodeTypeMetadata }

type KeyValue struct {
	Key   string
	Value any // string|number|boolean
}

// Type system implementations
type Type interface {
	Node
	TypeName() string
}

type PrimitiveType struct {
	Name string // "string", "int", etc.
}

func (p *PrimitiveType) NodeType() NodeType { return NodeTypePrimitiveType }
func (p *PrimitiveType) TypeName() string   { return p.Name }

type ArrayType struct {
	Element Type
}

func (a *ArrayType) NodeType() NodeType { return NodeTypeArrayType }
func (a *ArrayType) TypeName() string   { return "array" }

type InlineObjectType struct {
	Fields []Field
}

func (o *InlineObjectType) NodeType() NodeType { return NodeTypeInlineObjectType }
func (o *InlineObjectType) TypeName() string   { return "object" }

type TypeReference struct {
	Name string
}

func (t *TypeReference) NodeType() NodeType { return NodeTypeTypeReference }
func (t *TypeReference) TypeName() string   { return t.Name }
