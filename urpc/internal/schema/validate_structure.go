package schema

import (
	"fmt"
)

// TODO: Transpile to URPC and execute the semantic analyzer to better validate the schema.

// validateStructure validates the structure of an already parsed
// schema, it verifies:
//
//   - All custom rules are defined
//   - No duplicate custom rule names
//   - All custom types are defined
//   - No duplicate custom type names
//   - All rules applied to fields are defined
//
// This is a basic validation, to perform a more complex and complete validation
// first transpile the JSON schema to URPC schema and run the semantic analysis.
func validateStructure(sch Schema) error {
	// Check for duplicate rule names
	if err := validateRuleNames(sch); err != nil {
		return err
	}

	// Check for duplicate type names
	if err := validateTypeNames(sch); err != nil {
		return err
	}

	// Check that all custom types referenced in rules are defined
	if err := validateRuleTypeReferences(sch); err != nil {
		return err
	}

	// Check that all custom types referenced in fields are defined
	if err := validateFieldTypeReferences(sch); err != nil {
		return err
	}

	// Check that all rules applied to fields are defined
	if err := validateAppliedRules(sch); err != nil {
		return err
	}

	return nil
}

// validateRuleNames checks for duplicate rule names
func validateRuleNames(sch Schema) error {
	ruleNodes := sch.GetRuleNodes()
	ruleNames := make(map[string]bool)

	for _, rule := range ruleNodes {
		if _, exists := ruleNames[rule.Name]; exists {
			return fmt.Errorf("duplicate rule name: %s", rule.Name)
		}
		ruleNames[rule.Name] = true
	}

	return nil
}

// validateTypeNames checks for duplicate type names
func validateTypeNames(sch Schema) error {
	typeNodes := sch.GetTypeNodes()
	typeNames := make(map[string]bool)

	for _, typeNode := range typeNodes {
		if _, exists := typeNames[typeNode.Name]; exists {
			return fmt.Errorf("duplicate type name: %s", typeNode.Name)
		}
		typeNames[typeNode.Name] = true
	}

	return nil
}

// validateRuleTypeReferences checks that all custom types referenced in rules are defined
func validateRuleTypeReferences(sch Schema) error {
	ruleNodes := sch.GetRuleNodes()
	typeNodes := sch.GetTypeNodes()

	// Create a map of defined type names for quick lookup
	definedTypes := make(map[string]bool)
	for _, typeNode := range typeNodes {
		definedTypes[typeNode.Name] = true
	}

	// Check primitive types
	primitiveTypes := map[string]bool{
		"string":   true,
		"int":      true,
		"float":    true,
		"bool":     true,
		"datetime": true,
	}

	for _, rule := range ruleNodes {
		if rule.For == nil {
			return fmt.Errorf("rule '%s' has no 'for' clause", rule.Name)
		}

		// If the rule is for a custom type (starts with uppercase letter)
		if len(rule.For.Type) > 0 && rule.For.Type[0] >= 'A' && rule.For.Type[0] <= 'Z' {
			if !definedTypes[rule.For.Type] {
				return fmt.Errorf("rule '%s' references undefined type: %s", rule.Name, rule.For.Type)
			}
		} else if !primitiveTypes[rule.For.Type] {
			// If it's not a custom type, it must be a valid primitive type
			return fmt.Errorf("rule '%s' references invalid primitive type: %s", rule.Name, rule.For.Type)
		}
	}

	return nil
}

// validateFieldTypeReferences checks that all custom types referenced in fields are defined
func validateFieldTypeReferences(sch Schema) error {
	typeNodes := sch.GetTypeNodes()
	procNodes := sch.GetProcNodes()

	// Create a map of defined type names for quick lookup
	definedTypes := make(map[string]bool)
	for _, typeNode := range typeNodes {
		definedTypes[typeNode.Name] = true
	}

	// Check primitive types
	primitiveTypes := map[string]bool{
		"string":   true,
		"int":      true,
		"float":    true,
		"bool":     true,
		"datetime": true,
	}

	// Check type fields
	for _, typeNode := range typeNodes {
		for _, field := range typeNode.Fields {
			if field.IsNamed() && *field.TypeName != "" {
				typeName := *field.TypeName
				// If it's a custom type (starts with uppercase letter)
				if typeName[0] >= 'A' && typeName[0] <= 'Z' {
					if !definedTypes[typeName] {
						return fmt.Errorf("field '%s' in type '%s' references undefined type: %s",
							field.Name, typeNode.Name, typeName)
					}
				} else if !primitiveTypes[typeName] {
					// If it's not a custom type, it must be a valid primitive type
					return fmt.Errorf("field '%s' in type '%s' references invalid primitive type: %s",
						field.Name, typeNode.Name, typeName)
				}
			}

			// Check inline types recursively
			if field.IsInline() {
				if err := validateInlineTypeFields(field.TypeInline.Fields, definedTypes, primitiveTypes,
					fmt.Sprintf("inline type in field '%s' of type '%s'", field.Name, typeNode.Name)); err != nil {
					return err
				}
			}
		}
	}

	// Check procedure input and output fields
	for _, procNode := range procNodes {
		// Check input fields
		for _, field := range procNode.Input {
			if field.IsNamed() && *field.TypeName != "" {
				typeName := *field.TypeName
				// If it's a custom type (starts with uppercase letter)
				if typeName[0] >= 'A' && typeName[0] <= 'Z' {
					if !definedTypes[typeName] {
						return fmt.Errorf("input field '%s' in procedure '%s' references undefined type: %s",
							field.Name, procNode.Name, typeName)
					}
				} else if !primitiveTypes[typeName] {
					// If it's not a custom type, it must be a valid primitive type
					return fmt.Errorf("input field '%s' in procedure '%s' references invalid primitive type: %s",
						field.Name, procNode.Name, typeName)
				}
			}

			// Check inline types recursively
			if field.IsInline() {
				if err := validateInlineTypeFields(field.TypeInline.Fields, definedTypes, primitiveTypes,
					fmt.Sprintf("inline type in input field '%s' of procedure '%s'", field.Name, procNode.Name)); err != nil {
					return err
				}
			}
		}

		// Check output fields
		for _, field := range procNode.Output {
			if field.IsNamed() && *field.TypeName != "" {
				typeName := *field.TypeName
				// If it's a custom type (starts with uppercase letter)
				if typeName[0] >= 'A' && typeName[0] <= 'Z' {
					if !definedTypes[typeName] {
						return fmt.Errorf("output field '%s' in procedure '%s' references undefined type: %s",
							field.Name, procNode.Name, typeName)
					}
				} else if !primitiveTypes[typeName] {
					// If it's not a custom type, it must be a valid primitive type
					return fmt.Errorf("output field '%s' in procedure '%s' references invalid primitive type: %s",
						field.Name, procNode.Name, typeName)
				}
			}

			// Check inline types recursively
			if field.IsInline() {
				if err := validateInlineTypeFields(field.TypeInline.Fields, definedTypes, primitiveTypes,
					fmt.Sprintf("inline type in output field '%s' of procedure '%s'", field.Name, procNode.Name)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateInlineTypeFields recursively checks fields in inline type definitions
func validateInlineTypeFields(fields []FieldDefinition, definedTypes, primitiveTypes map[string]bool, context string) error {
	for _, field := range fields {
		if field.IsNamed() && *field.TypeName != "" {
			typeName := *field.TypeName
			// If it's a custom type (starts with uppercase letter)
			if typeName[0] >= 'A' && typeName[0] <= 'Z' {
				if !definedTypes[typeName] {
					return fmt.Errorf("field '%s' in %s references undefined type: %s",
						field.Name, context, typeName)
				}
			} else if !primitiveTypes[typeName] {
				// If it's not a custom type, it must be a valid primitive type
				return fmt.Errorf("field '%s' in %s references invalid primitive type: %s",
					field.Name, context, typeName)
			}
		}

		// Recursively check nested inline types
		if field.IsInline() {
			nestedContext := fmt.Sprintf("inline type in field '%s' of %s", field.Name, context)
			if err := validateInlineTypeFields(field.TypeInline.Fields, definedTypes, primitiveTypes, nestedContext); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAppliedRules checks that all rules applied to fields are defined
func validateAppliedRules(sch Schema) error {
	// Get all defined rule names
	ruleNodes := sch.GetRuleNodes()
	definedRules := make(map[string]bool)
	for _, rule := range ruleNodes {
		definedRules[rule.Name] = true
	}

	// Check rules in type fields
	typeNodes := sch.GetTypeNodes()
	for _, typeNode := range typeNodes {
		for _, field := range typeNode.Fields {
			if err := validateFieldRules(field, definedRules,
				fmt.Sprintf("field '%s' in type '%s'", field.Name, typeNode.Name)); err != nil {
				return err
			}

			// Check inline types recursively
			if field.IsInline() {
				if err := validateInlineTypeRules(field.TypeInline.Fields, definedRules,
					fmt.Sprintf("inline type in field '%s' of type '%s'", field.Name, typeNode.Name)); err != nil {
					return err
				}
			}
		}
	}

	// Check rules in procedure input and output fields
	procNodes := sch.GetProcNodes()
	for _, procNode := range procNodes {
		// Check input fields
		for _, field := range procNode.Input {
			if err := validateFieldRules(field, definedRules,
				fmt.Sprintf("input field '%s' in procedure '%s'", field.Name, procNode.Name)); err != nil {
				return err
			}

			// Check inline types recursively
			if field.IsInline() {
				if err := validateInlineTypeRules(field.TypeInline.Fields, definedRules,
					fmt.Sprintf("inline type in input field '%s' of procedure '%s'", field.Name, procNode.Name)); err != nil {
					return err
				}
			}
		}

		// Check output fields
		for _, field := range procNode.Output {
			if err := validateFieldRules(field, definedRules,
				fmt.Sprintf("output field '%s' in procedure '%s'", field.Name, procNode.Name)); err != nil {
				return err
			}

			// Check inline types recursively
			if field.IsInline() {
				if err := validateInlineTypeRules(field.TypeInline.Fields, definedRules,
					fmt.Sprintf("inline type in output field '%s' of procedure '%s'", field.Name, procNode.Name)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateFieldRules checks that all rules applied to a field are defined
func validateFieldRules(field FieldDefinition, definedRules map[string]bool, context string) error {
	for _, rule := range field.Rules {
		if !definedRules[rule.Rule] {
			return fmt.Errorf("%s uses undefined rule: %s", context, rule.Rule)
		}
	}
	return nil
}

// validateInlineTypeRules recursively checks rules in inline type definitions
func validateInlineTypeRules(fields []FieldDefinition, definedRules map[string]bool, context string) error {
	for _, field := range fields {
		if err := validateFieldRules(field, definedRules,
			fmt.Sprintf("field '%s' in %s", field.Name, context)); err != nil {
			return err
		}

		// Recursively check nested inline types
		if field.IsInline() {
			nestedContext := fmt.Sprintf("inline type in field '%s' of %s", field.Name, context)
			if err := validateInlineTypeRules(field.TypeInline.Fields, definedRules, nestedContext); err != nil {
				return err
			}
		}
	}
	return nil
}
