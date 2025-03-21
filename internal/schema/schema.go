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

// ProcedureType represents the type of a procedure (query or mutation)
type ProcedureType enum.Member[string]

// Procedure types enum values
var (
	ProcedureTypeQuery    = ProcedureType{"query"}
	ProcedureTypeMutation = ProcedureType{"mutation"}

	// ProcedureTypes is the enum containing all procedure types
	ProcedureTypes = enum.New(
		ProcedureTypeQuery,
		ProcedureTypeMutation,
	)
)

// MarshalJSON implements json.Marshaler for ProcedureType
func (p ProcedureType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

// UnmarshalJSON implements json.Unmarshaler for ProcedureType
func (p *ProcedureType) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Value)
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

// Procedure represents a procedure definition with its type, inputs, outputs and metadata
type Procedure struct {
	Type        ProcedureType  `json:"type"`
	Description string         `json:"description,omitempty"`
	Input       Field          `json:"input,omitzero"`
	Output      Field          `json:"output,omitzero"`
	Meta        map[string]any `json:"meta,omitempty"`
}

// Rules is the interface implemented by all rule types
type Rules interface {
	// IsOptional returns whether the field is optional
	IsOptional() bool
}

// StringRules defines validation rules for string fields
type StringRules struct {
	Optional  bool                `json:"optional,omitempty"`
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

// IsOptional implements Rules interface
func (r StringRules) IsOptional() bool {
	return r.Optional
}

// IntRules defines validation rules for integer fields
type IntRules struct {
	Optional bool             `json:"optional,omitempty"`
	Equals   RuleWithIntValue `json:"equals,omitzero"`
	Min      RuleWithIntValue `json:"min,omitzero"`
	Max      RuleWithIntValue `json:"max,omitzero"`
	Enum     RuleWithIntArray `json:"enum,omitzero"`
}

// IsOptional implements Rules interface
func (r IntRules) IsOptional() bool {
	return r.Optional
}

// FloatRules defines validation rules for float fields
type FloatRules struct {
	Optional bool                `json:"optional,omitempty"`
	Equals   RuleWithNumberValue `json:"equals,omitzero"`
	Min      RuleWithNumberValue `json:"min,omitzero"`
	Max      RuleWithNumberValue `json:"max,omitzero"`
	Enum     RuleWithNumberArray `json:"enum,omitzero"`
}

// IsOptional implements Rules interface
func (r FloatRules) IsOptional() bool {
	return r.Optional
}

// BooleanRules defines validation rules for boolean fields
type BooleanRules struct {
	Optional bool                 `json:"optional,omitempty"`
	Equals   RuleWithBooleanValue `json:"equals,omitzero"`
}

// IsOptional implements Rules interface
func (r BooleanRules) IsOptional() bool {
	return r.Optional
}

// ObjectRules defines validation rules for object fields
type ObjectRules struct {
	Optional bool `json:"optional,omitempty"`
}

// IsOptional implements Rules interface
func (r ObjectRules) IsOptional() bool {
	return r.Optional
}

// ArrayRules defines validation rules for array fields
type ArrayRules struct {
	Optional bool             `json:"optional,omitempty"`
	MinLen   RuleWithIntValue `json:"minLen,omitzero"`
	MaxLen   RuleWithIntValue `json:"maxLen,omitzero"`
}

// IsOptional implements Rules interface
func (r ArrayRules) IsOptional() bool {
	return r.Optional
}

// CustomTypeRules defines validation rules for custom type fields
type CustomTypeRules struct {
	Optional bool `json:"optional,omitempty"`
}

// IsOptional implements Rules interface
func (r CustomTypeRules) IsOptional() bool {
	return r.Optional
}

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

// RuleWithStringArray represents a validation rule with a string array
type RuleWithStringArray struct {
	Values       []string `json:"values"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
}

// RuleWithIntArray represents a validation rule with an integer array
type RuleWithIntArray struct {
	Values       []int  `json:"values"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RuleWithNumberArray represents a validation rule with a float array
type RuleWithNumberArray struct {
	Values       []float64 `json:"values"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}
