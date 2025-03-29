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
}

// TypePrimitive handles primitive types like string, int, float, boolean
type TypePrimitive struct {
	Name PrimitiveType
}

func (t TypePrimitive) NodeType() NodeType { return NodeTypeType }
func (t TypePrimitive) TypeName() TypeName { return TypeName(t.Name) }

// TypeArray handles array types
type TypeArray struct {
	ElementsType Type
}

func (t TypeArray) NodeType() NodeType { return NodeTypeType }
func (t TypeArray) TypeName() TypeName { return TypeNameArray }

// TypeObject handles object types
type TypeObject struct {
	Fields []Field
}

func (t TypeObject) NodeType() NodeType { return NodeTypeType }
func (t TypeObject) TypeName() TypeName { return TypeNameObject }

// TypeCustom handles custom types
type TypeCustom struct {
	Name string
}

func (t TypeCustom) NodeType() NodeType { return NodeTypeType }
func (t TypeCustom) TypeName() TypeName { return TypeNameCustom }
