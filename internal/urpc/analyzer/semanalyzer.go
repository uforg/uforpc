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
