package schema

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema.json
var jsonSchemaRaw string

// compileJSONSchema compiles the JSON schema to be used for schema input validation
func compileJSONSchema() (*jsonschema.Schema, error) {
	dummySchemaURL := "https://raw.githubusercontent.com/uforg/uforpc/refs/heads/main/internal/schema/schema.json"

	unmarshaled, err := jsonschema.UnmarshalJSON(strings.NewReader(jsonSchemaRaw))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json schema: %w", err)
	}

	c := jsonschema.NewCompiler()

	if err := c.AddResource(dummySchemaURL, unmarshaled); err != nil {
		return nil, err
	}

	return c.Compile(dummySchemaURL)
}

// ValidateSchema validates the input schema against the defined JSON schema
func ValidateSchema(schema string) error {
	unmarshaled, err := jsonschema.UnmarshalJSON(strings.NewReader(schema))
	if err != nil {
		return fmt.Errorf("failed to unmarshal input schema: %w", err)
	}

	jsonSchema, err := compileJSONSchema()
	if err != nil {
		return fmt.Errorf("failed to compile json schema: %w", err)
	}

	if err := jsonSchema.Validate(unmarshaled); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	return nil
}

// validateFieldType validates that a field's type exists if it's a custom type
func validateFieldType(path string, field Field, definedTypes map[string]bool) error {
	// Check if the custom type is defined (for built-in types, it's always valid)
	if !field.IsBuiltInType() && !definedTypes[field.Type] {
		return fmt.Errorf("undefined custom type: %s in field %s", field.Type, path)
	}

	// Recursively check fields of objects
	if field.Type == FieldTypeObject.Value && len(field.Fields) > 0 {
		for fieldName, subField := range field.Fields {
			if err := validateFieldType(path+"."+fieldName, subField, definedTypes); err != nil {
				return err
			}
		}
	}

	// Check array type
	if field.Type == FieldTypeArray.Value && field.ArrayType != nil {
		if err := validateFieldType(path+".arrayType", *field.ArrayType, definedTypes); err != nil {
			return err
		}
	}

	return nil
}

// validateSchemaTypes ensures all custom types referenced in the schema are defined
func validateSchemaTypes(schema Schema) error {
	// Set of all defined types
	definedTypes := make(map[string]bool)
	for typeName := range schema.Types {
		definedTypes[typeName] = true
	}

	// Check types referenced in custom type fields
	for typeName, typeField := range schema.Types {
		for fieldName, field := range typeField.Fields {
			if err := validateFieldType(typeName+"."+fieldName, field, definedTypes); err != nil {
				return err
			}
		}
	}

	// Check types referenced in procedure input/output
	for procName, proc := range schema.Procedures {
		// Check input fields
		for fieldName, field := range proc.Input {
			if err := validateFieldType(procName+".input."+fieldName, field, definedTypes); err != nil {
				return err
			}
		}

		// Check output fields
		for fieldName, field := range proc.Output {
			if err := validateFieldType(procName+".output."+fieldName, field, definedTypes); err != nil {
				return err
			}
		}
	}

	return nil
}
