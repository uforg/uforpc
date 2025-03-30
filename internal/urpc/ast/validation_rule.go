package ast

// ValidationRuleShape represents the shape of a validation rule.
type ValidationRuleShape string

const (
	ValidationRuleShapeSimple    ValidationRuleShape = "simple"
	ValidationRuleShapeWithValue ValidationRuleShape = "value"
	ValidationRuleShapeWithArray ValidationRuleShape = "array"
)

// ValidationRuleValueType represents the type of a validation rule value.
type ValidationRuleValueType string

const (
	ValidationRuleValueTypeNone    ValidationRuleValueType = "none"
	ValidationRuleValueTypeString  ValidationRuleValueType = "string"
	ValidationRuleValueTypeInt     ValidationRuleValueType = "int"
	ValidationRuleValueTypeFloat   ValidationRuleValueType = "float"
	ValidationRuleValueTypeBoolean ValidationRuleValueType = "boolean"
)

// ValidationRule represents a validation rule for a field.
type ValidationRule interface {
	NodeType() NodeType
	ValidationRuleShape() ValidationRuleShape
	ValidationRuleValueType() ValidationRuleValueType
	GetPosition() Position
}

// ValidationRuleSimple represents a simple validation rule that contains an error message.
type ValidationRuleSimple struct {
	Pos   Position
	Name  string
	Error string
}

func (v *ValidationRuleSimple) NodeType() NodeType { return NodeTypeValidationRule }
func (v *ValidationRuleSimple) ValidationRuleShape() ValidationRuleShape {
	return ValidationRuleShapeSimple
}
func (v *ValidationRuleSimple) ValidationRuleValueType() ValidationRuleValueType {
	return ValidationRuleValueTypeNone
}
func (v *ValidationRuleSimple) GetPosition() Position {
	return v.Pos
}

// ValidationRuleWithValue represents a validation rule that contains a string value
// and an error message.
type ValidationRuleWithValue struct {
	Pos       Position
	Name      string
	Value     string
	ValueType ValidationRuleValueType
	Error     string
}

func (v *ValidationRuleWithValue) NodeType() NodeType { return NodeTypeValidationRule }
func (v *ValidationRuleWithValue) ValidationRuleShape() ValidationRuleShape {
	return ValidationRuleShapeWithValue
}
func (v *ValidationRuleWithValue) ValidationRuleValueType() ValidationRuleValueType {
	return v.ValueType
}
func (v *ValidationRuleWithValue) GetPosition() Position {
	return v.Pos
}

// ValidationRuleWithArray represents a validation rule that contains an array of values
// and an error message.
type ValidationRuleWithArray struct {
	Pos       Position
	Name      string
	Values    []string
	ValueType ValidationRuleValueType
	Error     string
}

func (v *ValidationRuleWithArray) NodeType() NodeType { return NodeTypeValidationRule }
func (v *ValidationRuleWithArray) ValidationRuleShape() ValidationRuleShape {
	return ValidationRuleShapeWithArray
}
func (v *ValidationRuleWithArray) ValidationRuleValueType() ValidationRuleValueType {
	return v.ValueType
}
func (v *ValidationRuleWithArray) GetPosition() Position {
	return v.Pos
}
