// Package analyzer provides a semantic analyzer for the URPC AST.
//
// It performs the following checks:
//   - Custom validation rule names are unique and valid.
//   - Custom type names are unique and valid.
//   - Custom procedure names are unique and valid.
//   - All referenced validation rules exist.
//   - All referenced types exist.

package analyzer

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/util/strutil"
)

// AnalyzerError is an error that occurs during the analysis of the URPC schema.
type AnalyzerError struct {
	Message string
	Pos     ast.Position
	EndPos  ast.Position
}

func (e AnalyzerError) String() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Message)
}

func (e *AnalyzerError) Error() string {
	return e.String()
}

// Analyzer is a semantic analyzer for the URPC schema.
type Analyzer struct {
	sch             *ast.Schema
	customRuleNames map[string]ast.RuleDecl
	customTypeNames map[string]ast.TypeDecl
	procNames       map[string]ast.ProcDecl
	errors          []AnalyzerError
}

// NewAnalyzer creates a new Analyzer for the given URPC schema.
func NewAnalyzer(sch *ast.Schema) *Analyzer {
	return &Analyzer{
		sch:             sch,
		customRuleNames: map[string]ast.RuleDecl{},
		customTypeNames: map[string]ast.TypeDecl{},
		procNames:       map[string]ast.ProcDecl{},
		errors:          []AnalyzerError{},
	}
}

// Analyze analyzes the provided URPC schema.
//
// Returns:
//   - A list of errors that occurred during the analysis.
//   - The first error converted to Error interface if there are any errors.
func (a *Analyzer) Analyze() ([]AnalyzerError, error) {
	if err := a.validateVersion(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if err := a.validateImports(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if err := a.collectAndValidateCustomRuleNames(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if err := a.collectAndValidateCustomTypeNames(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if err := a.collectAndValidateProcNames(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if err := a.validateCustomRuleReferences(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if err := a.validateCustomTypeReferences(); err != nil {
		a.errors = append(a.errors, *err)
	}

	if len(a.errors) > 0 {
		return a.errors, &a.errors[0]
	}
	return nil, nil
}

// validateVersion validates the version of the URPC schema.
func (a *Analyzer) validateVersion() *AnalyzerError {
	versionsLen := len(a.sch.GetVersions())

	if versionsLen == 0 {
		return &AnalyzerError{
			Message: "version is required",
			Pos:     a.sch.Pos,
			EndPos:  a.sch.EndPos,
		}
	}

	if versionsLen > 1 {
		return &AnalyzerError{
			Message: "only one version is allowed",
			Pos:     a.sch.GetVersions()[1].Pos,
			EndPos:  a.sch.GetVersions()[1].EndPos,
		}
	}

	if a.sch.GetVersions()[0].Number != 1 {
		return &AnalyzerError{
			Message: "version must be 1",
			Pos:     a.sch.GetVersions()[0].Pos,
			EndPos:  a.sch.GetVersions()[0].EndPos,
		}
	}

	return nil
}

// validateImports validates the imports of the URPC schema.
func (a *Analyzer) validateImports() *AnalyzerError {
	importedNames := make(map[string]bool)

	for _, importStmt := range a.sch.GetImports() {
		if importStmt.Path == "" {
			return &AnalyzerError{
				Message: "import path is required",
				Pos:     importStmt.Pos,
				EndPos:  importStmt.EndPos,
			}
		}

		if !strings.HasSuffix(importStmt.Path, ".urpc") {
			return &AnalyzerError{
				Message: "import path must end with .urpc",
				Pos:     importStmt.Pos,
				EndPos:  importStmt.EndPos,
			}
		}

		if importedNames[importStmt.Path] {
			return &AnalyzerError{
				Message: "import path must be unique",
				Pos:     importStmt.Pos,
				EndPos:  importStmt.EndPos,
			}
		}

		importedNames[importStmt.Path] = true
	}

	return nil
}

// collectAndValidateCustomRuleNames collects and validates the custom rule names.
func (a *Analyzer) collectAndValidateCustomRuleNames() *AnalyzerError {
	for _, rule := range a.sch.GetRules() {
		if existing, ok := a.customRuleNames[rule.Name]; ok {
			return &AnalyzerError{
				Message: fmt.Sprintf(
					"\"%s\" custom rule name is already defined at %s:%d:%d",
					rule.Name, existing.Pos.Filename, existing.Pos.Line, existing.Pos.Column,
				),
				Pos:    rule.Pos,
				EndPos: rule.EndPos,
			}
		}

		if !strutil.IsCamelCase(rule.Name) {
			return &AnalyzerError{
				Message: fmt.Sprintf("\"%s\" custom rule name must be in camelCase", rule.Name),
				Pos:     rule.Pos,
				EndPos:  rule.EndPos,
			}
		}

		a.customRuleNames[rule.Name] = *rule
	}

	return nil
}

// collectAndValidateCustomTypeNames collects and validates the custom type names.
func (a *Analyzer) collectAndValidateCustomTypeNames() *AnalyzerError {
	for _, typeDecl := range a.sch.GetTypes() {
		if existing, ok := a.customTypeNames[typeDecl.Name]; ok {
			return &AnalyzerError{
				Message: fmt.Sprintf(
					"\"%s\" custom type name is already defined at %s:%d:%d",
					typeDecl.Name, existing.Pos.Filename, existing.Pos.Line, existing.Pos.Column,
				),
				Pos:    typeDecl.Pos,
				EndPos: typeDecl.EndPos,
			}
		}

		if !strutil.IsPascalCase(typeDecl.Name) {
			return &AnalyzerError{
				Message: fmt.Sprintf("\"%s\" custom type name must be in PascalCase", typeDecl.Name),
				Pos:     typeDecl.Pos,
				EndPos:  typeDecl.EndPos,
			}
		}

		a.customTypeNames[typeDecl.Name] = *typeDecl
	}

	return nil
}

// collectAndValidateProcNames collects and validates the procedure names.
func (a *Analyzer) collectAndValidateProcNames() *AnalyzerError {
	for _, proc := range a.sch.GetProcs() {
		if existing, ok := a.procNames[proc.Name]; ok {
			return &AnalyzerError{
				Message: fmt.Sprintf(
					"\"%s\" procedure is already defined at %s:%d:%d",
					proc.Name, existing.Pos.Filename, existing.Pos.Line, existing.Pos.Column,
				),
				Pos:    proc.Pos,
				EndPos: proc.EndPos,
			}
		}

		if !strutil.IsPascalCase(proc.Name) {
			return &AnalyzerError{
				Message: fmt.Sprintf("\"%s\" custom procedure name must be in PascalCase", proc.Name),
				Pos:     proc.Pos,
				EndPos:  proc.EndPos,
			}
		}

		a.procNames[proc.Name] = *proc
	}

	return nil
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
func (a *Analyzer) validateCustomRuleReferences() *AnalyzerError {
	// Map of built-in rules and the types they can be applied to
	builtInRulesMap := map[string][]string{
		// String rules
		"equals":    {"string", "boolean", "int", "float"},
		"contains":  {"string"},
		"minlen":    {"string", "array"},
		"maxlen":    {"string", "array"},
		"enum":      {"string", "int", "float"},
		"lowercase": {"string"},
		"uppercase": {"string"},
		"email":     {"string"},
		"uuid":      {"string"},
		"iso8601":   {"string"},
		"json":      {"string"},

		// Number rules
		"min":   {"int", "float", "datetime"},
		"max":   {"int", "float", "datetime"},
		"range": {"int", "float"},
	}

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
		if rule, isCustomRule := a.customRuleNames[ruleName]; isCustomRule {
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

			// Check if rule applies to custom types
			if _, isCustomType := a.customTypeNames[fieldType]; isCustomType {
				return forType == fieldType
			}

			// Check if rule applies to the field type
			return forType == fieldType
		}

		if supportedTypes, exists := builtInRulesMap[ruleName]; exists {
			for _, supportedType := range supportedTypes {
				if supportedType == "array" && isArray {
					return true
				}

				if supportedType == fieldType && (!isArray || supportedType == "array") {
					return true
				}
			}
			return false
		}

		return true
	}

	var checkFieldRules func([]*ast.Field, string) *AnalyzerError

	checkFieldRules = func(fields []*ast.Field, context string) *AnalyzerError {
		for _, field := range fields {
			baseType, isArray := getFieldBaseType(field)

			// Check rules in field.Children
			for _, child := range field.Children {
				if child.Rule == nil {
					continue
				}

				rule := child.Rule
				if _, isBuiltIn := builtInRulesMap[rule.Name]; !isBuiltIn {
					if _, isCustomRule := a.customRuleNames[rule.Name]; !isCustomRule {
						return &AnalyzerError{
							Message: fmt.Sprintf("referenced rule \"%s\" in %s is not defined", rule.Name, context),
							Pos:     rule.Pos,
							EndPos:  rule.EndPos,
						}
					}
				}

				if !canRuleApplyToType(rule.Name, baseType, isArray) {
					var ruleAppliesTo string

					if customRule, isCustomRule := a.customRuleNames[rule.Name]; isCustomRule {
						// Find the 'for' clause in the rule declaration
						var forType string
						for _, ruleChild := range customRule.Children {
							if ruleChild.For != nil {
								forType = ruleChild.For.For
								break
							}
						}
						ruleAppliesTo = forType
					} else if supportedTypes, isBuiltIn := builtInRulesMap[rule.Name]; isBuiltIn {
						ruleAppliesTo = strings.Join(supportedTypes, " or ")
					} else {
						ruleAppliesTo = "unknown"
					}

					if isArray {
						return &AnalyzerError{
							Message: fmt.Sprintf("rule \"%s\" in %s cannot be applied to array type \"%s[]\", it can only be applied to %s",
								rule.Name, context, baseType, ruleAppliesTo),
							Pos:    rule.Pos,
							EndPos: rule.EndPos,
						}
					}

					return &AnalyzerError{
						Message: fmt.Sprintf("rule \"%s\" in %s cannot be applied to type \"%s\", it can only be applied to %s",
							rule.Name, context, baseType, ruleAppliesTo),
						Pos:    rule.Pos,
						EndPos: rule.EndPos,
					}
				}
			}

			// Check inline object fields
			if field.Type.Base.Object != nil {
				inlineFields := extractFields(field.Type.Base.Object.Children)
				if err := checkFieldRules(inlineFields, fmt.Sprintf("inline object in field \"%s\"", field.Name)); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Check type declarations
	for _, typeDecl := range a.sch.GetTypes() {
		typeFields := extractFields(typeDecl.Children)
		if err := checkFieldRules(typeFields, fmt.Sprintf("type \"%s\"", typeDecl.Name)); err != nil {
			return err
		}
	}

	// Check procedure declarations
	for _, proc := range a.sch.GetProcs() {
		// Check input fields
		for _, child := range proc.Children {
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				if err := checkFieldRules(inputFields, fmt.Sprintf("input of procedure \"%s\"", proc.Name)); err != nil {
					return err
				}
			}

			// Check output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				if err := checkFieldRules(outputFields, fmt.Sprintf("output of procedure \"%s\"", proc.Name)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateCustomTypeReferences validates that all referenced custom types exist.
func (a *Analyzer) validateCustomTypeReferences() *AnalyzerError {
	primitiveTypes := map[string]bool{
		"string":   true,
		"int":      true,
		"float":    true,
		"boolean":  true,
		"datetime": true,
	}

	isValidType := func(typeName string) bool {
		return primitiveTypes[typeName] || a.customTypeNames[typeName].Name != ""
	}

	var checkFieldTypeReferences func([]*ast.Field, string) *AnalyzerError

	checkFieldTypeReferences = func(fields []*ast.Field, context string) *AnalyzerError {
		for _, field := range fields {
			if field.Type.Base.Named != nil {
				typeName := *field.Type.Base.Named

				if !primitiveTypes[typeName] {
					if !isValidType(typeName) {
						return &AnalyzerError{
							Message: fmt.Sprintf("referenced type \"%s\" in %s is not defined", typeName, context),
							Pos:     field.Type.Pos,
							EndPos:  field.Type.EndPos,
						}
					}
				}
			} else if field.Type.Base.Object != nil {
				// Extract fields from inline object
				inlineFields := extractFields(field.Type.Base.Object.Children)
				if err := checkFieldTypeReferences(inlineFields, fmt.Sprintf("inline object in field \"%s\"", field.Name)); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Check type declarations
	for _, typeDecl := range a.sch.GetTypes() {
		// Check extends clauses
		for _, extendTypeName := range typeDecl.Extends {
			if !isValidType(extendTypeName) {
				return &AnalyzerError{
					Message: fmt.Sprintf("type \"%s\" extends undefined type \"%s\"", typeDecl.Name, extendTypeName),
					Pos:     typeDecl.Pos,
					EndPos:  typeDecl.EndPos,
				}
			}
		}

		// Check fields
		typeFields := extractFields(typeDecl.Children)
		if err := checkFieldTypeReferences(typeFields, fmt.Sprintf("type \"%s\"", typeDecl.Name)); err != nil {
			return err
		}
	}

	// Check procedure declarations
	for _, proc := range a.sch.GetProcs() {
		for _, child := range proc.Children {
			// Check input fields
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				if err := checkFieldTypeReferences(inputFields, fmt.Sprintf("input of procedure \"%s\"", proc.Name)); err != nil {
					return err
				}
			}

			// Check output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				if err := checkFieldTypeReferences(outputFields, fmt.Sprintf("output of procedure \"%s\"", proc.Name)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
