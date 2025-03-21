package schema

import (
	"encoding/json"
	"fmt"
)

// ParseSchema parses and validates a JSON schema string into a Schema struct
func ParseSchema(schemaStr string) (Schema, error) {
	if err := ValidateSchema(schemaStr); err != nil {
		return Schema{}, fmt.Errorf("invalid schema: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return Schema{}, fmt.Errorf("error decoding schema: %w", err)
	}

	if err := processSchema(&schema); err != nil {
		return Schema{}, fmt.Errorf("error post-processing schema: %w", err)
	}

	if err := validateCustomTypeReferences(schema); err != nil {
		return Schema{}, fmt.Errorf("error validating custom type references: %w", err)
	}

	return schema, nil
}

// validateCustomTypeReferences ensures all custom types referenced in the schema are defined
func validateCustomTypeReferences(schema Schema) error {
	// Set of all defined types
	definedTypes := make(map[string]bool)
	for typeName := range schema.Types {
		definedTypes[typeName] = true
	}

	// Check types referenced in custom type fields
	for typeName, typeField := range schema.Types {
		if err := validateFieldCustomTypes(typeName, typeField, definedTypes); err != nil {
			return err
		}
	}

	// Check types referenced in procedure input/output
	for procName, proc := range schema.Procedures {
		// Check input if defined
		if proc.Input.Type != "" {
			if err := validateFieldType(procName+".input", proc.Input, definedTypes); err != nil {
				return err
			}
		}

		// Check output if defined
		if proc.Output.Type != "" {
			if err := validateFieldType(procName+".output", proc.Output, definedTypes); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateFieldCustomTypes recursively validates that all custom types referenced in a field are defined
func validateFieldCustomTypes(path string, field Field, definedTypes map[string]bool) error {
	if err := validateFieldType(path, field, definedTypes); err != nil {
		return err
	}

	// Recursively check fields of objects
	if field.Type == FieldTypeObject.Value && len(field.Fields) > 0 {
		for fieldName, subField := range field.Fields {
			if err := validateFieldCustomTypes(path+"."+fieldName, subField, definedTypes); err != nil {
				return err
			}
		}
	}

	// Check array type
	if field.Type == FieldTypeArray.Value && field.ArrayType != nil {
		if err := validateFieldCustomTypes(path+".arrayType", *field.ArrayType, definedTypes); err != nil {
			return err
		}
	}

	return nil
}

// validateFieldType validates that a field's type exists if it's a custom type
func validateFieldType(path string, field Field, definedTypes map[string]bool) error {
	// Skip built-in types
	if field.IsBuiltInType() {
		return nil
	}

	// Check if the custom type is defined
	if !definedTypes[field.Type] {
		return fmt.Errorf("undefined custom type: %s in field %s", field.Type, path)
	}

	return nil
}

// processSchema processes the Schema structure after decoding from JSON
func processSchema(schema *Schema) error {
	if schema.Version != 1 {
		return fmt.Errorf("unsupported schema version: %d", schema.Version)
	}

	// Process all type definitions
	for typeName, field := range schema.Types {
		if err := processField(typeName, &field, schema); err != nil {
			return fmt.Errorf("error processing type %s: %w", typeName, err)
		}
		schema.Types[typeName] = field
	}

	// Process all procedure definitions
	for procName, proc := range schema.Procedures {
		// Check if input is provided, otherwise skip
		if proc.Input.Type != "" {
			if err := processField(procName+".input", &proc.Input, schema); err != nil {
				return fmt.Errorf("error processing input for procedure %s: %w", procName, err)
			}
		}

		// Check if output is provided, otherwise skip
		if proc.Output.Type != "" {
			if err := processField(procName+".output", &proc.Output, schema); err != nil {
				return fmt.Errorf("error processing output for procedure %s: %w", procName, err)
			}
		}

		schema.Procedures[procName] = proc
	}

	return nil
}

// processField processes a field and its subfields recursively
func processField(path string, field *Field, schema *Schema) error {
	// Skip processing if field has no type defined
	if field.Type == "" {
		return nil
	}

	// Process built-in types
	if field.IsBuiltInType() {
		fieldType, _ := field.GetFieldType()

		switch fieldType {
		case FieldTypeObject:
			if len(field.Fields) == 0 {
				return fmt.Errorf("object field %s must have defined fields", path)
			}

			for fieldName, subField := range field.Fields {
				if err := processField(path+"."+fieldName, &subField, schema); err != nil {
					return fmt.Errorf("error processing field %s: %w", fieldName, err)
				}
				field.Fields[fieldName] = subField
			}

		case FieldTypeArray:
			if field.ArrayType == nil {
				return fmt.Errorf("array field %s must have a defined array type", path)
			}

			if err := processField(path+".arrayType", field.ArrayType, schema); err != nil {
				return fmt.Errorf("error processing array type for field %s: %w", path, err)
			}
		}

		// Process rules for the field
		rules, err := parseRules(field)
		if err != nil {
			return fmt.Errorf("error parsing rules for field %s: %w", path, err)
		}
		field.ProcessedRules = rules

	} else { // Process custom types
		if _, exists := schema.Types[field.Type]; !exists {
			return fmt.Errorf("undefined custom type: %s in field %s", field.Type, path)
		}

		// Process rules for custom type
		rules, err := parseRules(field)
		if err != nil {
			return fmt.Errorf("error parsing rules for field %s: %w", path, err)
		}
		field.ProcessedRules = rules
	}

	return nil
}

// parseRules parses field rules from JSON into their typed Go structs
func parseRules(field *Field) (Rules, error) {
	// If no rules are defined, create default rules based on field type
	if field.Rules == nil {
		return createDefaultRules(field)
	}

	// Parse rules based on field type
	if field.IsBuiltInType() {
		fieldType, _ := field.GetFieldType()
		return parseBuiltInTypeRules(fieldType, field.Rules)
	}

	// Parse custom type rules
	return parseCustomTypeRules(field.Rules)
}

// createDefaultRules creates default rules based on field type
func createDefaultRules(field *Field) (Rules, error) {
	if field.IsCustomType() {
		return CustomTypeRules{}, nil
	}

	fieldType, ok := field.GetFieldType()
	if !ok {
		return nil, fmt.Errorf("unknown field type: %s", field.Type)
	}

	switch fieldType {
	case FieldTypeString:
		return StringRules{}, nil
	case FieldTypeInt:
		return IntRules{}, nil
	case FieldTypeFloat:
		return FloatRules{}, nil
	case FieldTypeBoolean:
		return BooleanRules{}, nil
	case FieldTypeObject:
		return ObjectRules{}, nil
	case FieldTypeArray:
		return ArrayRules{}, nil
	default:
		return nil, fmt.Errorf("unsupported field type: %s", fieldType.Value)
	}
}

// parseBuiltInTypeRules parses rules for built-in field types
func parseBuiltInTypeRules(fieldType FieldType, rawRules json.RawMessage) (Rules, error) {
	switch fieldType {
	case FieldTypeString:
		var rules StringRules
		if err := json.Unmarshal(rawRules, &rules); err != nil {
			return nil, fmt.Errorf("error parsing string rules: %w", err)
		}
		return rules, nil

	case FieldTypeInt:
		var rules IntRules
		if err := json.Unmarshal(rawRules, &rules); err != nil {
			return nil, fmt.Errorf("error parsing int rules: %w", err)
		}
		return rules, nil

	case FieldTypeFloat:
		var rules FloatRules
		if err := json.Unmarshal(rawRules, &rules); err != nil {
			return nil, fmt.Errorf("error parsing float rules: %w", err)
		}
		return rules, nil

	case FieldTypeBoolean:
		var rules BooleanRules
		if err := json.Unmarshal(rawRules, &rules); err != nil {
			return nil, fmt.Errorf("error parsing boolean rules: %w", err)
		}
		return rules, nil

	case FieldTypeObject:
		var rules ObjectRules
		if err := json.Unmarshal(rawRules, &rules); err != nil {
			return nil, fmt.Errorf("error parsing object rules: %w", err)
		}
		return rules, nil

	case FieldTypeArray:
		var rules ArrayRules
		if err := json.Unmarshal(rawRules, &rules); err != nil {
			return nil, fmt.Errorf("error parsing array rules: %w", err)
		}
		return rules, nil

	default:
		return nil, fmt.Errorf("unsupported field type: %s", fieldType.Value)
	}
}

// parseCustomTypeRules parses rules for custom type fields
func parseCustomTypeRules(rawRules json.RawMessage) (Rules, error) {
	var rules CustomTypeRules
	if err := json.Unmarshal(rawRules, &rules); err != nil {
		return nil, fmt.Errorf("error parsing custom type rules: %w", err)
	}
	return rules, nil
}

// UnmarshalJSON ensures proper initialization of a Field when unmarshaling from JSON
func (f *Field) UnmarshalJSON(data []byte) error {
	type FieldTemp Field
	var temp FieldTemp

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*f = Field(temp)

	if f.Type == FieldTypeObject.Value && f.Fields == nil {
		f.Fields = make(map[string]Field)
	}

	return nil
}

// UnmarshalJSON ensures proper initialization of a Schema when unmarshaling from JSON
func (s *Schema) UnmarshalJSON(data []byte) error {
	type SchemaTemp Schema
	var temp SchemaTemp

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*s = Schema(temp)

	if s.Types == nil {
		s.Types = make(map[string]Field)
	}
	if s.Procedures == nil {
		s.Procedures = make(map[string]Procedure)
	}

	return nil
}
