package analyzer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

var primitiveTypes = map[string]bool{
	"string":   true,
	"int":      true,
	"float":    true,
	"bool":     true,
	"datetime": true,
}

// semanalyzer is the semantic alyzer phase for the URPC schema analyzer.
//
// It performs the following checks:
//   - Custom validation rule names are unique and valid.
//   - Custom type names are unique and valid.
//   - Custom procedure names are unique and valid.
//   - All referenced validation rules exist.
//   - All referenced types exist.
type semanalyzer struct {
	astSchema   *ast.Schema
	diagnostics []Diagnostic
}

// newSemanalyzer creates a new Semanalyzer for the given URPC schema. See more
// details in the Semanalyzer struct.
func newSemanalyzer(astSchema *ast.Schema) *semanalyzer {
	return &semanalyzer{
		astSchema:   astSchema,
		diagnostics: []Diagnostic{},
	}
}

// Analyze analyzes the provided URPC schema.
//
// Returns:
//   - A list of diagnostics that occurred during the analysis.
//   - The first diagnostic converted to Error interface if any.
func (a *semanalyzer) analyze() ([]Diagnostic, error) {
	a.validateCustomValidationRuleNames()
	a.validateCustomTypeNames()
	a.validateProcNames()
	a.validateStreamNames()
	a.validateCustomTypeReferences()
	a.validateCustomRuleReferences()
	a.validateTypeFieldUniqueness()
	a.validateTypeCircularDependencies()
	a.validateRuleStructure()
	a.validateProcStructure()
	a.validateStreamStructure()

	if len(a.diagnostics) > 0 {
		return a.diagnostics, a.diagnostics[0]
	}
	return nil, nil
}

// validateCustomValidationRuleNames validates the custom validation rule names and detects duplicates.
func (a *semanalyzer) validateCustomValidationRuleNames() {
	visited := map[string]Positions{}

	for _, ruleDecl := range a.astSchema.GetRules() {
		positions := Positions(ruleDecl.Positions)
		ruleName := ruleDecl.Name

		if decl, isDecl := visited[ruleName]; isDecl {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("custom validation rule \"%s\" is already declared at %s", ruleName, decl.Pos.String()),
			})
			continue
		}
		visited[ruleName] = positions

		if !strutil.IsCamelCase(ruleName) {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("custom validation rule name \"%s\" must be in camelCase", ruleName),
			})
			continue
		}
	}
}

// validateCustomTypeNames validates the custom type names and detects duplicates.
func (a *semanalyzer) validateCustomTypeNames() {
	visited := map[string]Positions{}

	for _, typeDecl := range a.astSchema.GetTypes() {
		positions := Positions(typeDecl.Positions)
		typeName := typeDecl.Name

		if decl, isDecl := visited[typeName]; isDecl {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("custom type \"%s\" is already declared at %s", typeName, decl.Pos.String()),
			})
			continue
		}
		visited[typeName] = positions

		if !strutil.IsPascalCase(typeName) {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("custom type name \"%s\" must be in PascalCase", typeName),
			})
			continue
		}
	}
}

// validateProcNames validates the procedure names and detects duplicates.
func (a *semanalyzer) validateProcNames() {
	visited := map[string]Positions{}

	for _, procDecl := range a.astSchema.GetProcs() {
		positions := Positions(procDecl.Positions)
		procName := procDecl.Name

		if decl, isDecl := visited[procName]; isDecl {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("procedure \"%s\" is already declared at %s", procName, decl.Pos.String()),
			})
			continue
		}
		visited[procName] = positions

		if !strutil.IsPascalCase(procName) {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("procedure name \"%s\" must be in PascalCase", procName),
			})
			continue
		}
	}
}

// validateStreamNames validates the stream names and detects duplicates.
func (a *semanalyzer) validateStreamNames() {
	visited := map[string]Positions{}

	for _, streamDecl := range a.astSchema.GetStreams() {
		positions := Positions(streamDecl.Positions)
		streamName := streamDecl.Name

		if decl, isDecl := visited[streamName]; isDecl {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("stream \"%s\" is already declared at %s", streamName, decl.Pos.String()),
			})
			continue
		}
		visited[streamName] = positions

		if !strutil.IsPascalCase(streamName) {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: positions,
				Message:   fmt.Sprintf("stream name \"%s\" must be in PascalCase", streamName),
			})
			continue
		}
	}
}

// validateCustomTypeReferences validates that all referenced custom types exist.
func (a *semanalyzer) validateCustomTypeReferences() {
	isValidType := func(typeName string) bool {
		if primitiveTypes[typeName] {
			return true
		}

		for _, typeDecl := range a.astSchema.GetTypes() {
			if typeDecl.Name == typeName {
				return true
			}
		}

		return false
	}

	var checkFieldTypeReferences func([]*ast.Field, string)
	checkFieldTypeReferences = func(fields []*ast.Field, context string) {
		for _, field := range fields {
			if field.Type.Base.Named != nil {
				typeName := *field.Type.Base.Named

				if !isValidType(typeName) {
					a.diagnostics = append(a.diagnostics, Diagnostic{
						Positions: Positions{
							Pos:    field.Type.Pos,
							EndPos: field.Type.EndPos,
						},
						Message: fmt.Sprintf("type \"%s\" referenced %s is not declared", typeName, context),
					})
				}
			} else if field.Type.Base.Object != nil {
				// Extract fields from inline object and recursively check them
				inlineFields := extractFields(field.Type.Base.Object.Children)
				checkFieldTypeReferences(inlineFields, fmt.Sprintf("at inline object at field \"%s\"", field.Name))
			}
		}
	}

	// Check custom validation rule declarations
	for _, ruleDecl := range a.astSchema.GetRules() {
		for _, child := range ruleDecl.Children {
			if child.For == nil {
				continue
			}

			if !isValidType(child.For.Type) {
				a.diagnostics = append(a.diagnostics, Diagnostic{
					Positions: Positions{
						Pos:    child.For.Pos,
						EndPos: child.For.EndPos,
					},
					Message: fmt.Sprintf("type \"%s\" referenced at for clause of rule \"%s\" is not declared", child.For.Type, ruleDecl.Name),
				})
			}
		}
	}

	// Check type declarations
	for _, typeDecl := range a.astSchema.GetTypes() {
		// Check fields
		typeFields := extractFields(typeDecl.Children)
		checkFieldTypeReferences(typeFields, fmt.Sprintf("at type \"%s\"", typeDecl.Name))
	}

	// Check procedure declarations
	for _, proc := range a.astSchema.GetProcs() {
		for _, child := range proc.Children {
			// Check input fields
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				checkFieldTypeReferences(inputFields, fmt.Sprintf("at input of procedure \"%s\"", proc.Name))
			}

			// Check output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				checkFieldTypeReferences(outputFields, fmt.Sprintf("at output of procedure \"%s\"", proc.Name))
			}
		}
	}

	// Check stream declarations
	for _, stream := range a.astSchema.GetStreams() {
		for _, child := range stream.Children {
			// Check input fields
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				checkFieldTypeReferences(inputFields, fmt.Sprintf("at input of stream \"%s\"", stream.Name))
			}

			// Check output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				checkFieldTypeReferences(outputFields, fmt.Sprintf("at output of stream \"%s\"", stream.Name))
			}
		}
	}
}

// Helper function to extract fields from FieldOrComment array
func extractFields(fieldOrComments []*ast.FieldOrComment) []*ast.Field {
	var fields []*ast.Field
	for _, foc := range fieldOrComments {
		if foc.Field != nil {
			fields = append(fields, foc.Field)
		}
	}
	return fields
}

// validateCustomRuleReferences validates that all custom rule references are valid.
func (a *semanalyzer) validateCustomRuleReferences() {
	// Map of primitive types and their supported built-in rules
	typesAndRulesMap := map[string][]string{
		"string":   {"equals", "contains", "minlen", "maxlen", "enum", "lowercase", "uppercase"},
		"int":      {"equals", "min", "max", "enum"},
		"float":    {"min", "max"},
		"bool":     {"equals"},
		"array":    {"minlen", "maxlen"},
		"datetime": {"min", "max"},
	}

	// Is the reverse of typesAndRulesMap and maps rules to their supported types
	rulesAndTypesMap := func() map[string][]string {
		rulesMap := map[string][]string{}
		for typeName, rules := range typesAndRulesMap {
			for _, rule := range rules {
				rulesMap[rule] = append(rulesMap[rule], typeName)
			}
		}
		return rulesMap
	}()

	// Helper function to get the base type of a field
	getFieldBaseType := func(field *ast.Field) (string, bool) {
		var baseType string
		isArray := field.Type.IsArray

		if field.Type.Base.Named != nil {
			baseType = *field.Type.Base.Named
		} else if field.Type.Base.Object != nil {
			baseType = "object"
		} else {
			return "", false
		}

		return baseType, isArray
	}

	// Helper function to check if a rule can be applied to a type
	canRuleApplyToType := func(ruleName, fieldType string, isArray bool) bool {
		isCustomRule := false
		for _, ruleDecl := range a.astSchema.GetRules() {
			if ruleDecl.Name == ruleName {
				isCustomRule = true
				break
			}
		}

		if isCustomRule {
			rule := &ast.RuleDecl{}
			for _, ruleItem := range a.astSchema.GetRules() {
				if ruleItem.Name == ruleName {
					rule = ruleItem
					break
				}
			}

			// Find the 'for' clause in the rule declaration
			var forType string
			for _, child := range rule.Children {
				if child.For != nil {
					forType = child.For.Type
					break
				}
			}

			// Check if rule applies to arrays
			if forType == "array" && isArray {
				return true
			}

			// Check if rule applies to the field type
			return forType == fieldType
		}

		builtinRules, hasBuiltinRules := typesAndRulesMap[fieldType]
		if !isArray && !hasBuiltinRules {
			return false
		}
		if isArray {
			builtinRules = typesAndRulesMap["array"]
		}

		return slices.Contains(builtinRules, ruleName)
	}

	// Helper function to check rules in a field
	var checkFieldRules func([]*ast.Field, string)
	checkFieldRules = func(fields []*ast.Field, context string) {
		for _, field := range fields {
			baseType, isArray := getFieldBaseType(field)

			// Check rules in field.Children
			for _, child := range field.Children {
				if child.Rule == nil {
					continue
				}

				// Check for rule existence
				rule := child.Rule
				if _, isBuiltIn := rulesAndTypesMap[rule.Name]; !isBuiltIn {
					isCustomRule := false
					for _, ruleDecl := range a.astSchema.GetRules() {
						if ruleDecl.Name == rule.Name {
							isCustomRule = true
							break
						}
					}

					if !isCustomRule {
						a.diagnostics = append(a.diagnostics, Diagnostic{
							Positions: Positions{
								Pos:    rule.Pos,
								EndPos: rule.EndPos,
							},
							Message: fmt.Sprintf("referenced rule \"%s\" %s is not declared", rule.Name, context),
						})
						continue
					}
				}

				// Check if the rule can be applied to the field type
				if !canRuleApplyToType(rule.Name, baseType, isArray) {
					var ruleAppliesTo string

					if customRule, isCustomRule := a.astSchema.GetRulesMap()[rule.Name]; isCustomRule {
						// Find the 'for' clause in the rule declaration
						var forType string
						for _, ruleChild := range customRule.Children {
							if ruleChild.For != nil {
								forType = ruleChild.For.Type
								break
							}
						}
						ruleAppliesTo = forType
					} else if supportedTypes, isBuiltIn := typesAndRulesMap[rule.Name]; isBuiltIn {
						ruleAppliesTo = strings.Join(supportedTypes, " or ")
					} else {
						ruleAppliesTo = "unknown"
					}

					if isArray {
						a.diagnostics = append(a.diagnostics, Diagnostic{
							Positions: Positions{
								Pos:    rule.Pos,
								EndPos: rule.EndPos,
							},
							Message: fmt.Sprintf("rule \"%s\" %s cannot be applied to array type \"%s[]\", it can only be applied to %s",
								rule.Name, context, baseType, ruleAppliesTo),
						})
					}

					if !isArray {
						a.diagnostics = append(a.diagnostics, Diagnostic{
							Positions: Positions{
								Pos:    rule.Pos,
								EndPos: rule.EndPos,
							},
							Message: fmt.Sprintf("rule \"%s\" %s cannot be applied to type \"%s\", it can only be applied to %s",
								rule.Name, context, baseType, ruleAppliesTo),
						})
					}
				}
			}

			// Check inline object fields
			if field.Type.Base.Object != nil {
				inlineFields := extractFields(field.Type.Base.Object.Children)
				checkFieldRules(inlineFields, fmt.Sprintf("at inline object at field \"%s\"", field.Name))
			}
		}
	}

	// Check type declarations
	for _, typeDecl := range a.astSchema.GetTypes() {
		typeFields := extractFields(typeDecl.Children)
		checkFieldRules(typeFields, fmt.Sprintf("at type \"%s\"", typeDecl.Name))
	}

	// Check procedure declarations
	for _, proc := range a.astSchema.GetProcs() {
		// Check input fields
		for _, child := range proc.Children {
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				checkFieldRules(inputFields, fmt.Sprintf("at input of procedure \"%s\"", proc.Name))
			}

			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				checkFieldRules(outputFields, fmt.Sprintf("at output of procedure \"%s\"", proc.Name))
			}
		}
	}

	// Check stream declarations
	for _, stream := range a.astSchema.GetStreams() {
		// Check input fields
		for _, child := range stream.Children {
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				checkFieldRules(inputFields, fmt.Sprintf("at input of stream \"%s\"", stream.Name))
			}

			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				checkFieldRules(outputFields, fmt.Sprintf("at output of stream \"%s\"", stream.Name))
			}
		}
	}
}

// validateTypeFieldUniqueness validates that fields in a type (including extended types) are unique.
func (a *semanalyzer) validateTypeFieldUniqueness() {
	for _, typeDecl := range a.astSchema.GetTypes() {
		// Collect all fields from this type and its extensions
		allFields := make(map[string]Positions)

		// First collect fields from the type itself
		fields := extractFields(typeDecl.Children)
		for _, field := range fields {
			// Check if this field already exists
			if existingPos, exists := allFields[field.Name]; exists {
				a.diagnostics = append(a.diagnostics, Diagnostic{
					Positions: Positions{
						Pos:    typeDecl.Pos,
						EndPos: typeDecl.EndPos,
					},
					Message: fmt.Sprintf(
						"field \"%s\" in type \"%s\" is already defined at %s",
						field.Name, typeDecl.Name, existingPos.Pos.String(),
					),
				})
			} else {
				// Add the field to the map
				allFields[field.Name] = Positions{
					Pos:    field.Pos,
					EndPos: field.EndPos,
				}
			}
		}
	}
}

// validateTypeCircularDependencies validates that there are no circular dependencies between types.
func (a *semanalyzer) validateTypeCircularDependencies() {
	types := a.astSchema.GetTypesMap()
	for name, typeDecl := range types {
		if err := validateTypeCircularDependenciesCheckType(name, types, []string{}); err != nil {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    typeDecl.Pos,
					EndPos: typeDecl.EndPos,
				},
				Message: err.Error(),
			})
		}
	}
}

// validateTypeCircularDependenciesCheckType checks if a type has a circular dependency.
func validateTypeCircularDependenciesCheckType(name string, types map[string]*ast.TypeDecl, stack []string) error {
	// Is it already in the stack (cycle)?
	for _, stackItem := range stack {
		if stackItem == name {
			return fmt.Errorf("circular dependency detected between types: %s", strings.Join(stack, " -> "))
		}
	}

	// Add it to the stack
	stack = append(stack, name)

	// Check every field in the type (including nested types)
	typ := types[name]
	for _, field := range extractFields(typ.Children) {
		if err := validateTypeCircularDependenciesCheckField(field.Type, types, stack); err != nil {
			return err
		}
	}

	return nil
}

// validateTypeCircularDependenciesCheckField checks if a field has a circular dependency.
func validateTypeCircularDependenciesCheckField(fieldType ast.FieldType, types map[string]*ast.TypeDecl, stack []string) error {
	// If it's a custom named type, check it
	if fieldType.Base.Named != nil {
		typeName := *fieldType.Base.Named
		if !primitiveTypes[typeName] {
			return validateTypeCircularDependenciesCheckType(typeName, types, stack)
		}
	}

	// If it's an inline object, check all its fields
	if fieldType.Base.Object != nil {
		objectFields := extractFields(fieldType.Base.Object.Children)
		for _, field := range objectFields {
			if err := validateTypeCircularDependenciesCheckField(field.Type, types, stack); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateRuleStructure validates that rule declarations have the correct structure:
// - Exactly one 'for' clause
// - At most one 'param' clause
// - At most one 'error' clause
func (a *semanalyzer) validateRuleStructure() {
	for _, ruleDecl := range a.astSchema.GetRules() {
		forCount := 0
		paramCount := 0
		errorCount := 0

		// Count the number of each clause
		for _, child := range ruleDecl.Children {
			if child.For != nil {
				forCount++
			}
			if child.Param != nil {
				paramCount++
			}
			if child.Error != nil {
				errorCount++
			}
		}

		// Validate 'for' clause
		if forCount == 0 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    ruleDecl.Pos,
					EndPos: ruleDecl.EndPos,
				},
				Message: fmt.Sprintf("rule \"%s\" must have exactly one 'for' clause", ruleDecl.Name),
			})
		} else if forCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    ruleDecl.Pos,
					EndPos: ruleDecl.EndPos,
				},
				Message: fmt.Sprintf("rule \"%s\" cannot have more than one 'for' clause", ruleDecl.Name),
			})
		}

		// Validate 'param' clause
		if paramCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    ruleDecl.Pos,
					EndPos: ruleDecl.EndPos,
				},
				Message: fmt.Sprintf("rule \"%s\" cannot have more than one 'param' clause", ruleDecl.Name),
			})
		}

		// Validate 'error' clause
		if errorCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    ruleDecl.Pos,
					EndPos: ruleDecl.EndPos,
				},
				Message: fmt.Sprintf("rule \"%s\" cannot have more than one 'error' clause", ruleDecl.Name),
			})
		}
	}
}

// validateProcStructure validates that procedure declarations have the correct structure:
// - At most one 'input' section
// - At most one 'output' section
// - At most one 'meta' section
func (a *semanalyzer) validateProcStructure() {
	for _, procDecl := range a.astSchema.GetProcs() {
		inputCount := 0
		outputCount := 0
		metaCount := 0

		// Count the number of each section
		for _, child := range procDecl.Children {
			if child.Input != nil {
				inputCount++
			}
			if child.Output != nil {
				outputCount++
			}
			if child.Meta != nil {
				metaCount++
			}
		}

		// Validate 'input' section
		if inputCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    procDecl.Pos,
					EndPos: procDecl.EndPos,
				},
				Message: fmt.Sprintf("procedure \"%s\" cannot have more than one 'input' section", procDecl.Name),
			})
		}

		// Validate 'output' section
		if outputCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    procDecl.Pos,
					EndPos: procDecl.EndPos,
				},
				Message: fmt.Sprintf("procedure \"%s\" cannot have more than one 'output' section", procDecl.Name),
			})
		}

		// Validate 'meta' section
		if metaCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    procDecl.Pos,
					EndPos: procDecl.EndPos,
				},
				Message: fmt.Sprintf("procedure \"%s\" cannot have more than one 'meta' section", procDecl.Name),
			})
		}
	}
}

// validateStreamStructure validates that stream declarations have the correct structure:
// - At most one 'input' section
// - At most one 'output' section
// - At most one 'meta' section
func (a *semanalyzer) validateStreamStructure() {
	for _, streamDecl := range a.astSchema.GetStreams() {
		inputCount := 0
		outputCount := 0
		metaCount := 0

		// Count the number of each section
		for _, child := range streamDecl.Children {
			if child.Input != nil {
				inputCount++
			}
			if child.Output != nil {
				outputCount++
			}
			if child.Meta != nil {
				metaCount++
			}
		}

		// Validate 'input' section
		if inputCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    streamDecl.Pos,
					EndPos: streamDecl.EndPos,
				},
				Message: fmt.Sprintf("stream \"%s\" cannot have more than one 'input' section", streamDecl.Name),
			})
		}

		// Validate 'output' section
		if outputCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    streamDecl.Pos,
					EndPos: streamDecl.EndPos,
				},
				Message: fmt.Sprintf("stream \"%s\" cannot have more than one 'output' section", streamDecl.Name),
			})
		}

		// Validate 'meta' section
		if metaCount > 1 {
			a.diagnostics = append(a.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    streamDecl.Pos,
					EndPos: streamDecl.EndPos,
				},
				Message: fmt.Sprintf("stream \"%s\" cannot have more than one 'meta' section", streamDecl.Name),
			})
		}
	}
}
