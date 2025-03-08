package schema

import (
	"errors"
	"regexp"
	"strings"
)

// Field represents a field within a type or as input/output
type Field struct {
	Type        string           `json:"type"`
	Description string           `json:"description,omitzero"`
	Rules       []RuleCatchAll   `json:"rules,omitzero"`
	Fields      map[string]Field `json:"fields,omitzero"`
}

// Validate checks if the field type is valid, if not it returns an
// error message explaining the issue
func (field Field) IsValid() error {
	if field.Type == "" {
		return errors.New("field type is required")
	}

	if strings.ContainsAny(field.Type, " \t\n") {
		return errors.New("field type should not contain any spaces")
	}

	isValidPrimitiveType := func() bool {
		for _, pt := range PrimitiveTypes {
			if field.Type == pt.Value {
				return true
			}
		}
		return false
	}()

	isValidCustomType := func() bool {
		matched, _ := regexp.MatchString(`^[A-Z][a-zA-Z0-9]*$`, field.Type)
		return matched
	}()

	if !isValidPrimitiveType && !isValidCustomType {
		return errors.New("field type is not valid primitive or custom type")
	}

	return nil
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
