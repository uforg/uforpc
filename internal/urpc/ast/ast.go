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

type Position struct {
	FileName string
	Line     int
	Column   int
}

type Node interface {
	NodeType() NodeType
	GetPosition() Position
}

// Root node
type Schema struct {
	Position   Position
	Version    Version
	Types      []TypeDeclaration
	Procedures []ProcDeclaration
}

func (s *Schema) NodeType() NodeType    { return NodeTypeSchema }
func (s *Schema) GetPosition() Position { return s.Position }

// Version
type Version struct {
	Position Position
	Value    int
}

func (v *Version) NodeType() NodeType    { return NodeTypeVersion }
func (v *Version) GetPosition() Position { return v.Position }

// Type system
type TypeDeclaration struct {
	Position Position
	Name     string
	Doc      string
	Fields   []Field
}

func (t *TypeDeclaration) NodeType() NodeType    { return NodeTypeTypeDeclaration }
func (t *TypeDeclaration) GetPosition() Position { return t.Position }

type Field struct {
	Position    Position
	Name        string
	Type        Type
	Optional    bool
	Validations []ValidationRule
}

func (f *Field) NodeType() NodeType    { return NodeTypeField }
func (f *Field) GetPosition() Position { return f.Position }

type ValidationRule struct {
	Position Position
	Name     string
	Params   []any
	ErrorMsg string
}

func (v *ValidationRule) NodeType() NodeType    { return NodeTypeValidationRule }
func (v *ValidationRule) GetPosition() Position { return v.Position }

// Procedures
type ProcDeclaration struct {
	Position Position
	Name     string
	Doc      string
	Input    Input
	Output   Output
	Metadata Metadata
}

func (p *ProcDeclaration) NodeType() NodeType    { return NodeTypeProcDeclaration }
func (p *ProcDeclaration) GetPosition() Position { return p.Position }

type Input struct {
	Position Position
	Fields   []Field
}

func (i *Input) NodeType() NodeType    { return NodeTypeInput }
func (i *Input) GetPosition() Position { return i.Position }

type Output struct {
	Position Position
	Fields   []Field
}

func (o *Output) NodeType() NodeType    { return NodeTypeOutput }
func (o *Output) GetPosition() Position { return o.Position }

// Metadata
type Metadata struct {
	Position Position
	Entries  []KeyValue
}

func (m *Metadata) NodeType() NodeType    { return NodeTypeMetadata }
func (m *Metadata) GetPosition() Position { return m.Position }

type KeyValue struct {
	Position Position
	Key      string
	Value    any // string|number|boolean
}

// Type system implementations
type Type interface {
	Node
	TypeName() string
}

type PrimitiveType struct {
	Position Position
	Name     string // "string", "int", etc.
}

func (p *PrimitiveType) NodeType() NodeType    { return NodeTypePrimitiveType }
func (p *PrimitiveType) GetPosition() Position { return p.Position }
func (p *PrimitiveType) TypeName() string      { return p.Name }

type ArrayType struct {
	Position Position
	Element  Type
}

func (a *ArrayType) NodeType() NodeType    { return NodeTypeArrayType }
func (a *ArrayType) GetPosition() Position { return a.Position }
func (a *ArrayType) TypeName() string      { return "array" }

type InlineObjectType struct {
	Position Position
	Fields   []Field
}

func (o *InlineObjectType) NodeType() NodeType    { return NodeTypeInlineObjectType }
func (o *InlineObjectType) GetPosition() Position { return o.Position }
func (o *InlineObjectType) TypeName() string      { return "object" }

type TypeReference struct {
	Position Position
	Name     string
}

func (t *TypeReference) NodeType() NodeType    { return NodeTypeTypeReference }
func (t *TypeReference) GetPosition() Position { return t.Position }
func (t *TypeReference) TypeName() string      { return t.Name }
