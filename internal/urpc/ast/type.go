package ast

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

type Type interface {
	NodeType() NodeType
	TypeName() TypeName
	GetPosition() Position
}

// TypePrimitive handles primitive types like string, int, float, boolean
type TypePrimitive struct {
	Pos  Position
	Name PrimitiveType
}

func (t TypePrimitive) NodeType() NodeType    { return NodeTypeType }
func (t TypePrimitive) TypeName() TypeName    { return TypeName(t.Name) }
func (t TypePrimitive) GetPosition() Position { return t.Pos }

// TypeArray handles array types
type TypeArray struct {
	Pos          Position
	ElementsType Type
}

func (t TypeArray) NodeType() NodeType    { return NodeTypeType }
func (t TypeArray) TypeName() TypeName    { return TypeNameArray }
func (t TypeArray) GetPosition() Position { return t.Pos }

// TypeObject handles object types
type TypeObject struct {
	Pos    Position
	Fields []Field
}

func (t TypeObject) NodeType() NodeType    { return NodeTypeType }
func (t TypeObject) TypeName() TypeName    { return TypeNameObject }
func (t TypeObject) GetPosition() Position { return t.Pos }

// TypeCustom handles custom types
type TypeCustom struct {
	Pos  Position
	Name string
}

func (t TypeCustom) NodeType() NodeType    { return NodeTypeType }
func (t TypeCustom) TypeName() TypeName    { return TypeNameCustom }
func (t TypeCustom) GetPosition() Position { return t.Pos }
