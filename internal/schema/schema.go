package schema

import (
	"encoding/json"

	"github.com/orsinium-labs/enum"
)

// FieldType represents possible field types in the schema
type FieldType enum.Member[string]

// Field types enum values
var (
	FieldTypeString  = FieldType{"string"}
	FieldTypeInt     = FieldType{"int"}
	FieldTypeFloat   = FieldType{"float"}
	FieldTypeBoolean = FieldType{"boolean"}
	FieldTypeObject  = FieldType{"object"}
	FieldTypeArray   = FieldType{"array"}

	// FieldTypes is the enum containing all built-in field types
	FieldTypes = enum.New(
		FieldTypeString,
		FieldTypeInt,
		FieldTypeFloat,
		FieldTypeBoolean,
		FieldTypeObject,
		FieldTypeArray,
	)
)

// MarshalJSON implements json.Marshaler for FieldType
func (f FieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

// UnmarshalJSON implements json.Unmarshaler for FieldType
func (f *FieldType) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.Value)
}

// Schema represents the main schema structure containing types and procedures
type Schema struct {
	Version    int                  `json:"version"`
	Types      map[string]Field     `json:"types"`
	Procedures map[string]Procedure `json:"procedures"`
}

// NewSchema creates a new Schema instance with initialized maps
func NewSchema() Schema {
	return Schema{
		Version:    1,
		Types:      make(map[string]Field),
		Procedures: make(map[string]Procedure),
	}
}

// Field represents a schema field with its type, description, validation rules and subelements
type Field struct {
	Type        string           `json:"type"`
	Description string           `json:"description,omitempty"`
	Optional    bool             `json:"optional,omitempty"`
	Rules       json.RawMessage  `json:"rules,omitempty"`
	Fields      map[string]Field `json:"fields,omitempty"`
	ArrayType   *Field           `json:"arrayType,omitempty"`

	// Processed rules - not directly marshaled/unmarshaled
	ProcessedRules Rules `json:"-"`
}

// IsBuiltInType checks if the field has a built-in type (string, int, etc.)
func (f Field) IsBuiltInType() bool {
	for _, fieldType := range FieldTypes.Members() {
		if fieldType.Value == f.Type {
			return true
		}
	}
	return false
}

// IsCustomType checks if the field has a custom type defined in the schema
func (f Field) IsCustomType() bool {
	return !f.IsBuiltInType()
}

// GetFieldType returns the FieldType enum for built-in types
func (f Field) GetFieldType() (FieldType, bool) {
	for _, fieldType := range FieldTypes.Members() {
		if fieldType.Value == f.Type {
			return fieldType, true
		}
	}
	return FieldType{}, false
}

// Procedure represents a procedure definition with its inputs, outputs and metadata
type Procedure struct {
	Description string         `json:"description,omitempty"`
	Input       Field          `json:"input,omitzero"`
	Output      Field          `json:"output,omitzero"`
	Meta        map[string]any `json:"meta,omitempty"`
}

// Rules is the interface implemented by all rule types
type Rules any

// StringRules defines validation rules for string fields
type StringRules struct {
	Equals    RuleWithStringValue `json:"equals,omitzero"`
	Contains  RuleWithStringValue `json:"contains,omitzero"`
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

// IntRules defines validation rules for integer fields
type IntRules struct {
	Equals RuleWithIntValue `json:"equals,omitzero"`
	Min    RuleWithIntValue `json:"min,omitzero"`
	Max    RuleWithIntValue `json:"max,omitzero"`
	Enum   RuleWithIntArray `json:"enum,omitzero"`
}

// FloatRules defines validation rules for float fields
type FloatRules struct {
	Equals RuleWithNumberValue `json:"equals,omitzero"`
	Min    RuleWithNumberValue `json:"min,omitzero"`
	Max    RuleWithNumberValue `json:"max,omitzero"`
	Enum   RuleWithNumberArray `json:"enum,omitzero"`
}

// BooleanRules defines validation rules for boolean fields
type BooleanRules struct {
	Equals RuleWithBooleanValue `json:"equals,omitzero"`
}

// ObjectRules defines validation rules for object fields
type ObjectRules struct{}

// ArrayRules defines validation rules for array fields
type ArrayRules struct {
	MinLen RuleWithIntValue `json:"minLen,omitzero"`
	MaxLen RuleWithIntValue `json:"maxLen,omitzero"`
}

// CustomTypeRules defines validation rules for custom type fields
type CustomTypeRules struct{}

// RuleSimple represents a simple validation rule with an optional error message
type RuleSimple struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithStringValue represents a validation rule with a string value
type RuleWithStringValue struct {
	Value        string `json:"value"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithIntValue represents a validation rule with an integer value
type RuleWithIntValue struct {
	Value        int    `json:"value"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithNumberValue represents a validation rule with a float value
type RuleWithNumberValue struct {
	Value        float64 `json:"value"`
	ErrorMessage string  `json:"errorMessage,omitempty"`
}

// RuleWithBooleanValue represents a validation rule with a boolean value
type RuleWithBooleanValue struct {
	Value        bool   `json:"value"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithStringArray represents a validation rule with an array of strings
type RuleWithStringArray struct {
	Values       []string `json:"values"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
}

// RuleWithIntArray represents a validation rule with an array of integers
type RuleWithIntArray struct {
	Values       []int  `json:"values"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithNumberArray represents a validation rule with an array of floats
type RuleWithNumberArray struct {
	Values       []float64 `json:"values"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}
