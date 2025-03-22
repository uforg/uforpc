package validator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

/** START FROM HERE **/

// schemaValidatorType represents the available types for validation
type schemaValidatorType string

const (
	SchemaTypeString  schemaValidatorType = "string"
	SchemaTypeNumber  schemaValidatorType = "number"
	SchemaTypeInt     schemaValidatorType = "int"
	SchemaTypeFloat   schemaValidatorType = "float"
	SchemaTypeBoolean schemaValidatorType = "boolean"
	SchemaTypeArray   schemaValidatorType = "array"
	SchemaTypeObject  schemaValidatorType = "object"
)

// schemaValidatorResult contains the result of a validation
type schemaValidatorResult struct {
	IsValid bool   `json:"isValid"`
	Error   string `json:"error,omitempty"`
}

// schemaValidator represents a validation schema
type schemaValidator struct {
	Type           schemaValidatorType
	ErrorMessage   string
	IsRequired     bool
	Pattern        *regexp.Regexp
	EqualsValue    any
	ContainsValue  *string
	LengthValue    *int
	MinLengthValue *int
	MaxLengthValue *int
	EnumValues     []any
	MinValue       *float64
	MaxValue       *float64
	IsEmail        bool
	IsIso8601      bool
	IsUUID         bool
	IsJSON         bool
	IsLowercase    bool
	IsUppercase    bool
	ArraySchema    *schemaValidator
	ObjectSchema   map[string]*schemaValidator
	LazySchemaFunc func() *schemaValidator
}

var regexes = struct {
	email   *regexp.Regexp
	iso8601 *regexp.Regexp
	uuid    *regexp.Regexp
}{
	email:   regexp.MustCompile(`^([^\.][\w\-+_.]*[^\.])(@\w+)(\.\w+(\.\w+)?[^.\W])$`),
	iso8601: regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[\+\-]\d{2}:\d{2})?$`),
	uuid:    regexp.MustCompile(`(?i)^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`),
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schemaType schemaValidatorType, errorMessage string) *schemaValidator {
	return &schemaValidator{
		Type:         schemaType,
		ErrorMessage: errorMessage,
	}
}

// Required marks the schema as required
func (sv *schemaValidator) Required(errorMessage string) *schemaValidator {
	sv.IsRequired = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Regex adds regex pattern validation
func (sv *schemaValidator) Regex(pattern string, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.Pattern = regexp.MustCompile(pattern)
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Equals adds equality validation
func (sv *schemaValidator) Equals(value any, errorMessage string) *schemaValidator {
	sv.EqualsValue = value
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Contains adds substring validation
func (sv *schemaValidator) Contains(value string, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.ContainsValue = &value
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Length adds exact length validation
func (sv *schemaValidator) Length(length int, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.LengthValue = &length
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// MinLength adds minimum length validation
func (sv *schemaValidator) MinLength(length int, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.MinLengthValue = &length
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// MaxLength adds maximum length validation
func (sv *schemaValidator) MaxLength(length int, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.MaxLengthValue = &length
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Enum adds enumeration validation
func (sv *schemaValidator) Enum(values []any, errorMessage string) *schemaValidator {
	sv.EnumValues = values
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Email adds email validation
func (sv *schemaValidator) Email(errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.IsEmail = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Iso8601 adds ISO8601 date validation
func (sv *schemaValidator) Iso8601(errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.IsIso8601 = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// UUID adds UUID validation
func (sv *schemaValidator) UUID(errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.IsUUID = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// JSON adds JSON validation
func (sv *schemaValidator) JSON(errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.IsJSON = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Lowercase adds lowercase validation
func (sv *schemaValidator) Lowercase(errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.IsLowercase = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Uppercase adds uppercase validation
func (sv *schemaValidator) Uppercase(errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeString {
		return sv
	}
	sv.IsUppercase = true
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Min adds minimum value validation
func (sv *schemaValidator) Min(value float64, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeNumber && sv.Type != SchemaTypeInt && sv.Type != SchemaTypeFloat {
		return sv
	}
	sv.MinValue = &value
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Max adds maximum value validation
func (sv *schemaValidator) Max(value float64, errorMessage string) *schemaValidator {
	if sv.Type != SchemaTypeNumber && sv.Type != SchemaTypeInt && sv.Type != SchemaTypeFloat {
		return sv
	}
	sv.MaxValue = &value
	if errorMessage != "" {
		sv.ErrorMessage = errorMessage
	}
	return sv
}

// Array creates an array schema
func (sv *schemaValidator) Array(schema *schemaValidator, errorMessage string) *schemaValidator {
	newSchema := NewSchemaValidator(SchemaTypeArray, errorMessage)
	newSchema.ArraySchema = schema
	return newSchema
}

// Object creates an object schema
func (sv *schemaValidator) Object(schema map[string]*schemaValidator, errorMessage string) *schemaValidator {
	newSchema := NewSchemaValidator(SchemaTypeObject, errorMessage)
	newSchema.ObjectSchema = schema
	return newSchema
}

// Lazy creates a lazy schema for recursive validation
func Lazy(schema func() *schemaValidator, errorMessage string) *schemaValidator {
	lazySchema := NewSchemaValidator(SchemaTypeObject, errorMessage)
	lazySchema.LazySchemaFunc = schema
	return lazySchema
}

// validateType checks if a value matches the expected type
func (sv *schemaValidator) validateType(value any) bool {
	if value == nil {
		return false
	}

	switch sv.Type {
	case SchemaTypeString:
		_, ok := value.(string)
		return ok
	case SchemaTypeNumber, SchemaTypeFloat:
		_, ok := value.(float64)
		return ok
	case SchemaTypeInt:
		num, ok := value.(float64)
		return ok && float64(int(num)) == num
	case SchemaTypeBoolean:
		_, ok := value.(bool)
		return ok
	case SchemaTypeArray:
		_, ok := value.([]any)
		return ok
	case SchemaTypeObject:
		_, ok := value.(map[string]any)
		return ok
	default:
		return false
	}
}

// Validate performs the validation according to the schema
func (sv *schemaValidator) Validate(value any) schemaValidatorResult {
	// Handle required validation
	if value == nil {
		if sv.IsRequired {
			return schemaValidatorResult{
				IsValid: false,
				Error:   sv.ErrorMessage,
			}
		}
		return schemaValidatorResult{IsValid: true}
	}

	// Type validation
	if !sv.validateType(value) {
		return schemaValidatorResult{
			IsValid: false,
			Error:   fmt.Sprintf("Invalid type, expected %s", sv.Type),
		}
	}

	// Equals validation
	if sv.EqualsValue != nil {
		if value != sv.EqualsValue {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "Value does not equal expected value",
			}
		}
	}

	// Enum validation
	if sv.EnumValues != nil {
		found := false
		for _, enumVal := range sv.EnumValues {
			if value == enumVal {
				found = true
				break
			}
		}
		if !found {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "Value is not in the allowed enumeration",
			}
		}
	}

	// String validations
	if sv.Type == SchemaTypeString {
		strValue := value.(string)

		if sv.Pattern != nil && !sv.Pattern.MatchString(strValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String does not match pattern",
			}
		}

		if sv.ContainsValue != nil && !strings.Contains(strValue, *sv.ContainsValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String does not contain the required substring",
			}
		}

		if sv.LengthValue != nil && len(strValue) != *sv.LengthValue {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String length does not match expected length",
			}
		}

		if sv.MinLengthValue != nil && len(strValue) < *sv.MinLengthValue {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is shorter than minimum length",
			}
		}

		if sv.MaxLengthValue != nil && len(strValue) > *sv.MaxLengthValue {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is longer than maximum length",
			}
		}

		if sv.IsEmail && !regexes.email.MatchString(strValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is not a valid email address",
			}
		}

		if sv.IsIso8601 && !regexes.iso8601.MatchString(strValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is not a valid ISO8601 date",
			}
		}

		if sv.IsUUID && !regexes.uuid.MatchString(strValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is not a valid UUID",
			}
		}

		if sv.IsJSON {
			var js any
			if err := json.Unmarshal([]byte(strValue), &js); err != nil {
				return schemaValidatorResult{
					IsValid: false,
					Error:   "String is not valid JSON",
				}
			}
		}

		if sv.IsLowercase && strValue != strings.ToLower(strValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is not in lowercase",
			}
		}

		if sv.IsUppercase && strValue != strings.ToUpper(strValue) {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "String is not in uppercase",
			}
		}
	}

	// Number validations
	if sv.Type == SchemaTypeNumber || sv.Type == SchemaTypeInt || sv.Type == SchemaTypeFloat {
		numValue := value.(float64)

		if sv.MinValue != nil && numValue < *sv.MinValue {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "Number is less than the minimum allowed value",
			}
		}

		if sv.MaxValue != nil && numValue > *sv.MaxValue {
			return schemaValidatorResult{
				IsValid: false,
				Error:   "Number is greater than the maximum allowed value",
			}
		}
	}

	// Array validation
	if sv.Type == SchemaTypeArray {
		arrValue := value.([]any)
		if sv.ArraySchema != nil {
			for _, item := range arrValue {
				result := sv.ArraySchema.Validate(item)
				if !result.IsValid {
					return result
				}
			}
		}
	}

	// Object validation
	if sv.Type == SchemaTypeObject {
		if sv.LazySchemaFunc != nil {
			return sv.LazySchemaFunc().Validate(value)
		}

		mapValue := value.(map[string]any)
		if sv.ObjectSchema != nil {
			for key, schema := range sv.ObjectSchema {
				val, exists := mapValue[key]
				if !exists {
					if schema.IsRequired {
						return schemaValidatorResult{
							IsValid: false,
							Error:   "Field is required",
						}
					}
					continue
				}
				result := schema.Validate(val)
				if !result.IsValid {
					return result
				}
			}
		}
	}

	return schemaValidatorResult{IsValid: true}
}

// schValidatorFactory is a factory for creating validation schemas
type schValidatorFactory struct {
	String  func(errorMessage string) *schemaValidator
	Number  func(errorMessage string) *schemaValidator
	Int     func(errorMessage string) *schemaValidator
	Float   func(errorMessage string) *schemaValidator
	Boolean func(errorMessage string) *schemaValidator
	Array   func(schema *schemaValidator, errorMessage string) *schemaValidator
	Object  func(schema map[string]*schemaValidator, errorMessage string) *schemaValidator
	Lazy    func(schema func() *schemaValidator, errorMessage string) *schemaValidator
}

// schValidator is the global validator factory instance
var schValidator = schValidatorFactory{
	String: func(errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeString, errorMessage)
	},
	Number: func(errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeNumber, errorMessage)
	},
	Int: func(errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeInt, errorMessage)
	},
	Float: func(errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeFloat, errorMessage)
	},
	Boolean: func(errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeBoolean, errorMessage)
	},
	Array: func(schema *schemaValidator, errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeArray, errorMessage).Array(schema, errorMessage)
	},
	Object: func(schema map[string]*schemaValidator, errorMessage string) *schemaValidator {
		return NewSchemaValidator(SchemaTypeObject, errorMessage).Object(schema, errorMessage)
	},
	Lazy: func(schema func() *schemaValidator, errorMessage string) *schemaValidator {
		return Lazy(schema, errorMessage)
	},
}
