package analyzer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/util/strutil"
)

// semanalyzer is the semantic alyzer phase for the URPC schema analyzer.
//
// It performs the following checks:
//   - Custom validation rule names are unique and valid.
//   - Custom type names are unique and valid.
//   - Custom procedure names are unique and valid.
//   - All referenced validation rules exist.
//   - All referenced types exist.
type semanalyzer struct {
	combinedSchema *CombinedSchema
	diagnostics    []Diagnostic
}

// newSemanalyzer creates a new Semanalyzer for the given URPC schema. See more
// details in the Semanalyzer struct.
func newSemanalyzer(combinedSchema CombinedSchema) *semanalyzer {
	return &semanalyzer{
		combinedSchema: &combinedSchema,
		diagnostics:    []Diagnostic{},
	}
}

// Analyze analyzes the provided URPC schema.
//
// Returns:
//   - A list of diagnostics that occurred during the analysis.
//   - The first diagnostic converted to Error interface if any.
func (a *semanalyzer) analyze() ([]Diagnostic, error) {
	a.validateCustomRuleNames()
	a.validateCustomTypeNames()
	a.validateProcNames()
	a.validateCustomTypeReferences()
	a.validateCustomRuleReferences()
	a.validateTypeExtendRecursion()
	a.validateTypeFieldUniqueness()
	a.validateTypeCircularDependencies()
	a.validateRuleStructure()
	a.validateProcStructure()

	if len(a.diagnostics) > 0 {
		return a.diagnostics, a.diagnostics[0]
	}
	return nil, nil
}

// validateCustomRuleNames validates the custom rule names and detects duplicates.
func (a *semanalyzer) validateCustomRuleNames() {
	visited := map[string]Positions{}

	for _, ruleDecl := range a.combinedSchema.Schema.GetRules() {
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

	for _, typeDecl := range a.combinedSchema.Schema.GetTypes() {
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

	for _, procDecl := range a.combinedSchema.Schema.GetProcs() {
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

// validateCustomTypeReferences validates that all referenced custom types exist.
func (a *semanalyzer) validateCustomTypeReferences() {
	primitiveTypes := map[string]bool{
		"string":   true,
		"int":      true,
		"float":    true,
		"boolean":  true,
		"datetime": true,
	}

	isValidType := func(typeName string) bool {
		if primitiveTypes[typeName] {
			return true
		}
		_, isDeclared := a.combinedSchema.TypeDecls[typeName]
		return isDeclared
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

	// Check type declarations
	for _, typeDecl := range a.combinedSchema.Schema.GetTypes() {
		// Check extends clauses
		for _, extendTypeName := range typeDecl.Extends {
			if !isValidType(extendTypeName) {
				a.diagnostics = append(a.diagnostics, Diagnostic{
					Positions: Positions{
						Pos:    typeDecl.Pos,
						EndPos: typeDecl.EndPos,
					},
					Message: fmt.Sprintf("type \"%s\" extends non-declared type \"%s\"", typeDecl.Name, extendTypeName),
				})
			}

			// Cannot reference itself
			if extendTypeName == typeDecl.Name {
				a.diagnostics = append(a.diagnostics, Diagnostic{
					Positions: Positions{
						Pos:    typeDecl.Pos,
						EndPos: typeDecl.EndPos,
					},
					Message: fmt.Sprintf("type \"%s\" cannot extend itself", typeDecl.Name),
				})
			}
		}

		// Check fields
		typeFields := extractFields(typeDecl.Children)
		checkFieldTypeReferences(typeFields, fmt.Sprintf("at type \"%s\"", typeDecl.Name))
	}

	// Check procedure declarations
	for _, proc := range a.combinedSchema.Schema.GetProcs() {
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

// validateTypeExtendRecursion validates that type extensions are not recursive.
func (a *semanalyzer) validateTypeExtendRecursion() {
	// For each type, check if it extends itself directly or indirectly
	for _, typeDecl := range a.combinedSchema.Schema.GetTypes() {
		visited := make(map[string]bool)
		visited[typeDecl.Name] = true

		for _, extendTypeName := range typeDecl.Extends {
			a.checkTypeExtendRecursion(typeDecl, extendTypeName, visited)
		}
	}
}

// checkTypeExtendRecursion recursively checks if a type extension creates a cycle.
func (a *semanalyzer) checkTypeExtendRecursion(originalType *ast.TypeDecl, currentExtendName string, visited map[string]bool) {
	// If we've already visited this type in this path, we have a cycle
	if visited[currentExtendName] {
		a.diagnostics = append(a.diagnostics, Diagnostic{
			Positions: Positions{
				Pos:    originalType.Pos,
				EndPos: originalType.EndPos,
			},
			Message: fmt.Sprintf(
				"recursive type extension detected: type \"%s\" extends \"%s\" which creates a cycle",
				originalType.Name, currentExtendName,
			),
		})
		return
	}

	// Mark this type as visited in the current path
	visited[currentExtendName] = true

	// Get the extended type
	extendedType, exists := a.combinedSchema.TypeDecls[currentExtendName]
	if !exists {
		// Skip if the extended type doesn't exist (this is caught by another validation)
		return
	}

	// Check all extensions of the extended type
	for _, nextExtendName := range extendedType.Extends {
		a.checkTypeExtendRecursion(originalType, nextExtendName, visited)
	}

	// Remove this type from the visited set when backtracking
	visited[currentExtendName] = false
}

// validateTypeFieldUniqueness validates that fields in a type (including extended types) are unique.
func (a *semanalyzer) validateTypeFieldUniqueness() {
	for _, typeDecl := range a.combinedSchema.Schema.GetTypes() {
		// Collect all fields from this type and its extensions
		allFields := make(map[string]Positions)

		// First collect fields from the type itself
		fields := extractFields(typeDecl.Children)
		for _, field := range fields {
			allFields[field.Name] = Positions{
				Pos:    field.Pos,
				EndPos: field.EndPos,
			}
		}

		// Then collect fields from all extended types
		for _, extendTypeName := range typeDecl.Extends {
			extendedType, exists := a.combinedSchema.TypeDecls[extendTypeName]
			if !exists {
				// Skip if the extended type doesn't exist (this is caught by another validation)
				continue
			}

			extendedFields := extractFields(extendedType.Children)
			for _, extendedField := range extendedFields {
				// Check if this field already exists
				if existingPos, exists := allFields[extendedField.Name]; exists {
					a.diagnostics = append(a.diagnostics, Diagnostic{
						Positions: Positions{
							Pos:    typeDecl.Pos,
							EndPos: typeDecl.EndPos,
						},
						Message: fmt.Sprintf("field \"%s\" in type \"%s\" is already defined in extended type \"%s\" at %s",
							extendedField.Name, typeDecl.Name, extendTypeName, existingPos.Pos.String()),
					})
				} else {
					// Add the field to the map
					allFields[extendedField.Name] = Positions{
						Pos:    extendedField.Pos,
						EndPos: extendedField.EndPos,
					}
				}
			}
		}
	}
}

// validateTypeCircularDependencies validates that there are no circular dependencies between types,
// unless one of the fields in the cycle is optional.
func (a *semanalyzer) validateTypeCircularDependencies() {
	// Build a dependency graph
	dependencyGraph := make(map[string]map[string]bool) // type -> dependencies
	optionalFields := make(map[string]map[string]bool)  // type -> optional field types

	// Initialize the graph
	for typeName := range a.combinedSchema.TypeDecls {
		dependencyGraph[typeName] = make(map[string]bool)
		optionalFields[typeName] = make(map[string]bool)
	}

	// Build the dependency graph
	for typeName, typeDecl := range a.combinedSchema.TypeDecls {
		fields := extractFields(typeDecl.Children)
		for _, field := range fields {
			if field.Type.Base.Named != nil {
				fieldTypeName := *field.Type.Base.Named

				// Check if this is a custom type (not a primitive)
				if _, isCustomType := a.combinedSchema.TypeDecls[fieldTypeName]; isCustomType {
					// Add the dependency
					dependencyGraph[typeName][fieldTypeName] = true

					// If the field is optional, mark it
					if field.Optional {
						optionalFields[typeName][fieldTypeName] = true
					}
				}
			}
		}
	}

	// Check for circular dependencies
	for typeName := range dependencyGraph {
		visited := make(map[string]bool)
		recStack := make(map[string]bool)
		path := []string{}

		a.checkCircularDependency(typeName, dependencyGraph, optionalFields, visited, recStack, path)
	}
}

// checkCircularDependency performs a DFS to detect circular dependencies.
func (a *semanalyzer) checkCircularDependency(typeName string, graph map[string]map[string]bool,
	optionalFields map[string]map[string]bool, visited, recStack map[string]bool, path []string) {

	// Mark the current node as visited and part of recursion stack
	visited[typeName] = true
	recStack[typeName] = true
	path = append(path, typeName)

	// Visit all dependencies
	for dependency := range graph[typeName] {
		// If not visited, then recurse
		if !visited[dependency] {
			a.checkCircularDependency(dependency, graph, optionalFields, visited, recStack, path)
		} else if recStack[dependency] {
			// If the dependency is already in the recursion stack, we have a cycle

			// Find the start of the cycle in the path
			cycleStart := -1
			for i, t := range path {
				if t == dependency {
					cycleStart = i
					break
				}
			}

			// Extract the cycle
			cycle := append(path[cycleStart:], dependency)

			// Check if any field in the cycle is optional
			hasOptionalField := false
			for i := 0; i < len(cycle)-1; i++ {
				current := cycle[i]
				next := cycle[i+1]
				if optionalFields[current][next] {
					hasOptionalField = true
					break
				}
			}

			// If no optional field is found, report the circular dependency
			if !hasOptionalField {
				// Get the type declaration for the current type
				typeDecl := a.combinedSchema.TypeDecls[typeName]

				a.diagnostics = append(a.diagnostics, Diagnostic{
					Positions: Positions{
						Pos:    typeDecl.Pos,
						EndPos: typeDecl.EndPos,
					},
					Message: fmt.Sprintf(
						"circular dependency detected between types: %s. At least one field in the cycle must be optional",
						strings.Join(cycle, " -> "),
					),
				})
			}
		}
	}

	// Remove the current node from recursion stack
	recStack[typeName] = false
}

// validateRuleStructure validates that rule declarations have the correct structure:
// - Exactly one 'for' clause
// - At most one 'param' clause
// - At most one 'error' clause
func (a *semanalyzer) validateRuleStructure() {
	for _, ruleDecl := range a.combinedSchema.Schema.GetRules() {
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
	for _, procDecl := range a.combinedSchema.Schema.GetProcs() {
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

// validateCustomRuleReferences validates that all custom rule references are valid.
func (a *semanalyzer) validateCustomRuleReferences() {
	// Map of primitive types and their supported built-in rules
	typesAndRulesMap := map[string][]string{
		"string":   {"equals", "contains", "minlen", "maxlen", "enum", "lowercase", "uppercase"},
		"int":      {"equals", "min", "max", "enum"},
		"float":    {"min", "max"},
		"boolean":  {"equals"},
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

	getFieldBaseType := func(field *ast.Field) (string, bool) {
		var baseType string
		isArray := field.Type.Depth > 0

		if field.Type.Base.Named != nil {
			baseType = *field.Type.Base.Named
		} else if field.Type.Base.Object != nil {
			baseType = "object"
		} else {
			return "", false
		}

		return baseType, isArray
	}

	canRuleApplyToType := func(ruleName, fieldType string, isArray bool) bool {
		_, isCustomRule := a.combinedSchema.RuleDecls[ruleName]

		if isCustomRule {
			rule := &ast.RuleDecl{}
			for _, ruleItem := range a.combinedSchema.Schema.GetRules() {
				if ruleItem.Name == ruleName {
					rule = ruleItem
					break
				}
			}

			// Find the 'for' clause in the rule declaration
			var forType string
			for _, child := range rule.Children {
				if child.For != nil {
					forType = child.For.For
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

	var checkFieldRules func([]*ast.Field, string)
	checkFieldRules = func(fields []*ast.Field, context string) {
		for _, field := range fields {
			baseType, isArray := getFieldBaseType(field)

			// Check rules in field.Children
			for _, child := range field.Children {
				if child.Rule == nil {
					continue
				}

				rule := child.Rule
				if _, isBuiltIn := rulesAndTypesMap[rule.Name]; !isBuiltIn {
					if _, isCustomRule := a.combinedSchema.RuleDecls[rule.Name]; !isCustomRule {
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

				if !canRuleApplyToType(rule.Name, baseType, isArray) {
					var ruleAppliesTo string

					if customRule, isCustomRule := a.combinedSchema.RuleDecls[rule.Name]; isCustomRule {
						// Find the 'for' clause in the rule declaration
						var forType string
						for _, ruleChild := range customRule.Children {
							if ruleChild.For != nil {
								forType = ruleChild.For.For
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
	for _, typeDecl := range a.combinedSchema.Schema.GetTypes() {
		typeFields := extractFields(typeDecl.Children)
		checkFieldRules(typeFields, fmt.Sprintf("at type \"%s\"", typeDecl.Name))
	}

	// Check procedure declarations
	for _, proc := range a.combinedSchema.Schema.GetProcs() {
		// Check input fields
		for _, child := range proc.Children {
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				checkFieldRules(inputFields, fmt.Sprintf("at input of procedure \"%s\"", proc.Name))
			}

			// Check output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				checkFieldRules(outputFields, fmt.Sprintf("at output of procedure \"%s\"", proc.Name))
			}
		}
	}
}
