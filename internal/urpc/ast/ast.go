package ast

type NodeType int

const (
	_ NodeType = iota
	NodeTypeSchema
	NodeTypeVersion
	NodeTypeCustomRuleDeclaration
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

// Schema is the root node of the AST representing an entire URPC schema.
type Schema struct {
	Version     Version
	CustomRules []CustomRuleDeclaration
	Types       []TypeDeclaration
	Procedures  []ProcDeclaration
}

func (s *Schema) NodeType() NodeType { return NodeTypeSchema }

// Version represents the version of the URPC schema.
type Version struct {
	IsSet bool
	Value int
}

func (v *Version) NodeType() NodeType { return NodeTypeVersion }

// CustomRulePrimitiveType represents the primitive types allowed in custom rule parameters.
type CustomRulePrimitiveType int

const (
	_ CustomRulePrimitiveType = iota
	CustomRulePrimitiveTypeString
	CustomRulePrimitiveTypeInt
	CustomRulePrimitiveTypeFloat
	CustomRulePrimitiveTypeBoolean
)

// CustomRuleParamType represents the allowed parameter types for custom rules.
type CustomRuleParamType struct {
	IsArray bool
	Type    CustomRulePrimitiveType // Only primitive types are allowed
}

// CustomRuleDeclaration represents a custom validation rule declaration in the URPC schema.
type CustomRuleDeclaration struct {
	Doc      string
	Name     string
	For      TypeName
	Param    CustomRuleParamType
	ErrorMsg string
}

func (c *CustomRuleDeclaration) NodeType() NodeType { return NodeTypeCustomRuleDeclaration }

// TypeDeclaration represents a type declaration in the URPC schema.
type TypeDeclaration struct {
	Name   string
	Doc    string
	Fields []Field
}

func (t *TypeDeclaration) NodeType() NodeType { return NodeTypeTypeDeclaration }

// Field represents a field in a type declaration or procedure input/output.
type Field struct {
	Name            string
	Type            Type
	Optional        bool
	ValidationRules []ValidationRule
}

func (f *Field) NodeType() NodeType { return NodeTypeField }

// ValidationRuleShape represents the shape of a validation rule.
type ValidationRuleShape int

const (
	_ ValidationRuleShape = iota
	ValidationRuleShapeSimple
	ValidationRuleShapeWithValue
	ValidationRuleShapeWithArray
)

// ValidationRuleValueType represents the type of a validation rule value.
type ValidationRuleValueType int

const (
	ValidationRuleValueTypeNone ValidationRuleValueType = iota
	ValidationRuleValueTypeString
	ValidationRuleValueTypeInt
	ValidationRuleValueTypeFloat
	ValidationRuleValueTypeBoolean
)

// ValidationRule represents a validation rule for a field.
type ValidationRule interface {
	NodeType() NodeType
	ValidationRuleShape() ValidationRuleShape
	ValidationRuleValueType() ValidationRuleValueType
}

// ValidationRuleSimple represents a simple validation rule that contains an error message.
type ValidationRuleSimple struct {
	RuleName     string
	ErrorMessage string
}

func (v *ValidationRuleSimple) NodeType() NodeType { return NodeTypeValidationRule }
func (v *ValidationRuleSimple) ValidationRuleShape() ValidationRuleShape {
	return ValidationRuleShapeSimple
}
func (v *ValidationRuleSimple) ValidationRuleValueType() ValidationRuleValueType {
	return ValidationRuleValueTypeNone
}

// ValidationRuleWithValue represents a validation rule that contains a string value
// and an error message.
type ValidationRuleWithValue struct {
	RuleName     string
	Value        string
	ValueType    ValidationRuleValueType
	ErrorMessage string
}

func (v *ValidationRuleWithValue) NodeType() NodeType { return NodeTypeValidationRule }
func (v *ValidationRuleWithValue) ValidationRuleShape() ValidationRuleShape {
	return ValidationRuleShapeWithValue
}
func (v *ValidationRuleWithValue) ValidationRuleValueType() ValidationRuleValueType {
	return v.ValueType
}

// ValidationRuleWithArray represents a validation rule that contains an array of values
// and an error message.
type ValidationRuleWithArray struct {
	RuleName     string
	Values       []string
	ValueType    ValidationRuleValueType
	ErrorMessage string
}

func (v *ValidationRuleWithArray) NodeType() NodeType { return NodeTypeValidationRule }
func (v *ValidationRuleWithArray) ValidationRuleShape() ValidationRuleShape {
	return ValidationRuleShapeWithArray
}
func (v *ValidationRuleWithArray) ValidationRuleValueType() ValidationRuleValueType {
	return v.ValueType
}

// ProcDeclaration represents a procedure declaration in the URPC schema.
type ProcDeclaration struct {
	Name     string
	Doc      string
	Input    ProcInput
	Output   ProcOutput
	Metadata ProcMeta
}

func (p *ProcDeclaration) NodeType() NodeType { return NodeTypeProcDeclaration }

// ProcInput represents the input of a procedure.
type ProcInput struct {
	Fields []Field
}

func (i *ProcInput) NodeType() NodeType { return NodeTypeInput }

// ProcOutput represents the output of a procedure.
type ProcOutput struct {
	Fields []Field
}

func (o *ProcOutput) NodeType() NodeType { return NodeTypeOutput }

// ProcMeta represents the metadata of a procedure.
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

// Type represents a type in the URPC schema, either a primitive type or a custom type.
type Type interface {
	TypeName() TypeName
}

// TypeString represents the string type.
type TypeString struct{}

func (t *TypeString) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeString) TypeName() TypeName { return TypeNameString }

// TypeInt represents the int type.
type TypeInt struct{}

func (t *TypeInt) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeInt) TypeName() TypeName { return TypeNameInt }

// TypeFloat represents the float type.
type TypeFloat struct{}

func (t *TypeFloat) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeFloat) TypeName() TypeName { return TypeNameFloat }

// TypeBoolean represents the boolean type.
type TypeBoolean struct{}

func (t *TypeBoolean) NodeType() NodeType { return NodeTypePrimitiveType }
func (t *TypeBoolean) TypeName() TypeName { return TypeNameBoolean }

// TypeObject represents an inline object type.
type TypeObject struct {
	Fields []Field
}

func (o *TypeObject) NodeType() NodeType { return NodeTypeInlineObjectType }
func (o *TypeObject) TypeName() TypeName { return TypeNameObject }

// TypeArray represents an array type.
type TypeArray struct {
	ArrayType Type
}

func (a *TypeArray) NodeType() NodeType { return NodeTypeArrayType }
func (a *TypeArray) TypeName() TypeName { return TypeNameArray }

// TypeCustom represents a custom type reference.
type TypeCustom struct {
	Name string
}

func (t *TypeCustom) NodeType() NodeType { return NodeTypeTypeReference }
func (t *TypeCustom) TypeName() TypeName { return TypeNameCustom }
