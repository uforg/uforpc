package schema

import (
	"github.com/orsinium-labs/enum"
)

// FieldType represents possible field types
type FieldType enum.Member[string]

// Predefined values for FieldType
var (
	FieldTypeString  = FieldType{"string"}
	FieldTypeInt     = FieldType{"int"}
	FieldTypeFloat   = FieldType{"float"}
	FieldTypeBoolean = FieldType{"boolean"}
	FieldTypeObject  = FieldType{"object"}
	FieldTypeArray   = FieldType{"array"}

	// Enum containing all field types
	FieldTypes = enum.New(
		FieldTypeString,
		FieldTypeInt,
		FieldTypeFloat,
		FieldTypeBoolean,
		FieldTypeObject,
		FieldTypeArray,
	)
)

// ProcedureType represents procedure types
type ProcedureType enum.Member[string]

// Predefined values for ProcedureType
var (
	ProcedureTypeQuery    = ProcedureType{"query"}
	ProcedureTypeMutation = ProcedureType{"mutation"}

	// Enum containing all procedure types
	ProcedureTypes = enum.New(
		ProcedureTypeQuery,
		ProcedureTypeMutation,
	)
)

// Schema represents the main schema structure
type Schema struct {
	Version    int                  `json:"version"`
	Types      map[string]Field     `json:"types"`
	Procedures map[string]Procedure `json:"procedures"`
}

// Field represents a schema field
type Field struct {
	Type        string           `json:"type"`
	Description string           `json:"description,omitempty"`
	Rules       Rules            `json:"rules,omitempty"`
	Fields      map[string]Field `json:"fields,omitempty"`
	ArrayType   *Field           `json:"arrayType,omitempty"`
}

// IsBuiltInType checks if the type is one of the built-in types
func (f Field) IsBuiltInType() bool {
	for _, fieldType := range FieldTypes.Members() {
		if fieldType.Value == f.Type {
			return true
		}
	}
	return false
}

// IsCustomType checks if the type is a custom type
func (f Field) IsCustomType() bool {
	return !f.IsBuiltInType()
}

// GetFieldType returns the field type as FieldType if it's a built-in type
func (f Field) GetFieldType() (FieldType, bool) {
	for _, fieldType := range FieldTypes.Members() {
		if fieldType.Value == f.Type {
			return fieldType, true
		}
	}
	return FieldType{}, false
}

// Procedure represents a schema procedure
type Procedure struct {
	Type        string         `json:"type"`
	Description string         `json:"description,omitempty"`
	Input       *Field         `json:"input,omitempty"`
	Output      *Field         `json:"output,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
}

// GetProcedureType returns the procedure type as ProcedureType
func (p Procedure) GetProcedureType() (ProcedureType, bool) {
	for _, procType := range ProcedureTypes.Members() {
		if procType.Value == p.Type {
			return procType, true
		}
	}
	return ProcedureType{}, false
}

// Rules is an interface for all rule types
type Rules any

// StringRules represents rules for string fields
type StringRules struct {
	Optional  bool                `json:"optional,omitzero"`
	Equals    RuleWithStringValue `json:"equals,omitzero"`
	Contains  RuleWithStringValue `json:"contains,omitzero"`
	Regex     RuleWithStringValue `json:"regex,omitzero"`
	MinLen    RuleWithIntValue    `json:"minLen,omitzero"`
	MaxLen    RuleWithIntValue    `json:"maxLen,omitzero"`
	Enum      RuleWithStringArray `json:"enum,omitzero"`
	Email     RuleSimple          `json:"email,omitzero"`
	ISO8601   RuleSimple          `json:"iso8601,omitzero"`
	UUID      RuleSimple          `json:"uuid,omitzero"`
	JSON      RuleSimple          `json:"json,omitzero"`
	Lowercase RuleSimple          `json:"lowercase,omitzero"`
	Uppercase RuleSimple          `json:"uppercase,omitzero"`
}

// IntRules represents rules for int fields
type IntRules struct {
	Optional bool             `json:"optional,omitzero"`
	Equals   RuleWithIntValue `json:"equals,omitzero"`
	Min      RuleWithIntValue `json:"min,omitzero"`
	Max      RuleWithIntValue `json:"max,omitzero"`
	Enum     RuleWithIntArray `json:"enum,omitzero"`
}

// FloatRules represents rules for float fields
type FloatRules struct {
	Optional bool                `json:"optional,omitzero"`
	Equals   RuleWithNumberValue `json:"equals,omitzero"`
	Min      RuleWithNumberValue `json:"min,omitzero"`
	Max      RuleWithNumberValue `json:"max,omitzero"`
	Enum     RuleWithNumberArray `json:"enum,omitzero"`
}

// BooleanRules represents rules for boolean fields
type BooleanRules struct {
	Optional bool                 `json:"optional,omitzero"`
	Equals   RuleWithBooleanValue `json:"equals,omitzero"`
}

// ObjectRules represents rules for object fields
type ObjectRules struct {
	Optional bool `json:"optional,omitzero"`
}

// ArrayRules represents rules for array fields
type ArrayRules struct {
	Optional bool             `json:"optional,omitzero"`
	MinLen   RuleWithIntValue `json:"minLen,omitzero"`
	MaxLen   RuleWithIntValue `json:"maxLen,omitzero"`
}

// CustomTypeRules represents rules for custom type fields
type CustomTypeRules struct {
	Optional bool `json:"optional,omitzero"`
}

// RuleSimple represents a simple rule with an optional error message
type RuleSimple struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithStringValue represents a rule with a string value
type RuleWithStringValue struct {
	Value        string `json:"value"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithIntValue represents a rule with an int value
type RuleWithIntValue struct {
	Value        int    `json:"value"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithNumberValue represents a rule with a float value
type RuleWithNumberValue struct {
	Value        float64 `json:"value"`
	ErrorMessage string  `json:"errorMessage,omitempty"`
}

// RuleWithBooleanValue represents a rule with a boolean value
type RuleWithBooleanValue struct {
	Value        bool   `json:"value"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithStringArray represents a rule with a string array
type RuleWithStringArray struct {
	Values       []string `json:"values"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
}

// RuleWithIntArray represents a rule with an int array
type RuleWithIntArray struct {
	Values       []int  `json:"values"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithNumberArray represents a rule with a float array
type RuleWithNumberArray struct {
	Values       []float64 `json:"values"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

// NewSchema creates a new Schema instance with initialized maps
func NewSchema() *Schema {
	return &Schema{
		Version:    1,
		Types:      make(map[string]Field),
		Procedures: make(map[string]Procedure),
	}
}

// GetRulesForField returns the appropriate rules based on the field type
func GetRulesForField(field Field) Rules {
	fieldType, ok := field.GetFieldType()
	if !ok {
		return &CustomTypeRules{}
	}

	switch fieldType.Value {
	case FieldTypeString.Value:
		return &StringRules{}
	case FieldTypeInt.Value:
		return &IntRules{}
	case FieldTypeFloat.Value:
		return &FloatRules{}
	case FieldTypeBoolean.Value:
		return &BooleanRules{}
	case FieldTypeObject.Value:
		return &ObjectRules{}
	case FieldTypeArray.Value:
		return &ArrayRules{}
	default:
		return nil
	}
}
