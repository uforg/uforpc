package schema

import "strings"

// Field represents a field within a type or as input/output
type Field struct {
	Type        string           `json:"type"`
	Description string           `json:"description,omitzero"`
	Rules       []RuleCatchAll   `json:"rules,omitzero"`
	Fields      map[string]Field `json:"fields,omitzero"`
}

// IsArray checks if the field type is an array
func (field Field) IsArray() bool {
	field.Type = strings.TrimSpace(field.Type)
	return len(field.Type) >= 2 && field.Type[len(field.Type)-2:] == "[]"
}

// GetArrayDepth returns the depth of an array type
// For example
//   - "string" returns 0
//   - "string[]" returns 1
//   - "string[][]" returns 2
//   - and so on...
func (field Field) GetArrayDepth() int {
	field.Type = strings.TrimSpace(field.Type)
	depth := 0
	for field.IsArray() {
		depth++
		field.Type = field.Type[:len(field.Type)-2]
	}
	return depth
}

// GetBaseType returns the base type of an array type
// For example
//   - "string" returns "string"
//   - "string[]" returns "string"
//   - "string[][]" returns "string"
//   - and so on...
func (field Field) GetBaseType() string {
	field.Type = strings.TrimSpace(field.Type)
	for field.IsArray() {
		field.Type = field.Type[:len(field.Type)-2]
	}
	return field.Type
}
