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

// Node is the interface that all AST nodes implement.
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
	Input    ProcInput
	Output   ProcOutput
	Metadata ProcMeta
}

func (p *ProcDeclaration) NodeType() NodeType { return NodeTypeProcDeclaration }

type ProcInput struct {
	Fields []Field
}

func (i *ProcInput) NodeType() NodeType { return NodeTypeInput }

type ProcOutput struct {
	Fields []Field
}

func (o *ProcOutput) NodeType() NodeType { return NodeTypeOutput }

// ProcMeta
type ProcMeta struct {
	Entries []ProcMetaKV
}

func (m *ProcMeta) NodeType() NodeType { return NodeTypeMetadata }

type ProcMetaValueTypeName string

const (
	ProcMetaValueTypeString  ProcMetaValueTypeName = "string"
	ProcMetaValueTypeInt     ProcMetaValueTypeName = "int"
	ProcMetaValueTypeFloat   ProcMetaValueTypeName = "float"
	ProcMetaValueTypeBoolean ProcMetaValueTypeName = "boolean"
)

type ProcMetaKV struct {
	Type  ProcMetaValueTypeName
	Key   string
	Value string
}

type TypeName string

const (
	TypeNameString  TypeName = "string"
	TypeNameInt     TypeName = "int"
	TypeNameFloat   TypeName = "float"
	TypeNameBoolean TypeName = "boolean"
	TypeNameObject  TypeName = "object"
	TypeNameArray   TypeName = "array"
	TypeNameCustom  TypeName = "custom"
)

// Type system implementations
type Type interface {
	TypeName() TypeName
}

type TypeString struct{}

func (t *TypeString) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeString) TypeName() TypeName { return TypeNameString }

type TypeInt struct{}

func (t *TypeInt) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeInt) TypeName() TypeName { return TypeNameInt }

type TypeFloat struct{}

func (t *TypeFloat) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeFloat) TypeName() TypeName { return TypeNameFloat }

type TypeBoolean struct{}

func (t *TypeBoolean) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeBoolean) TypeName() TypeName { return TypeNameBoolean }

type TypeObject struct {
	Fields []Field
}

func (o *TypeObject) NodeType() NodeType { return NodeTypeInlineObjectType }
func (o *TypeObject) TypeName() TypeName { return TypeNameObject }

type TypeArray struct {
	ArrayType Type
}

func (a *TypeArray) NodeType() NodeType { return NodeTypeArrayType }
func (a *TypeArray) TypeName() TypeName { return TypeNameArray }

type TypeCustom struct {
	Name string
}

func (t *TypeCustom) NodeType() NodeType { return NodeTypeTypeReference }
func (t *TypeCustom) TypeName() TypeName { return TypeNameCustom }
