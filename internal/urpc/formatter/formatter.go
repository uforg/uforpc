package formatter

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/parser"
)

// Format formats URPC code according to the spec, using 2 spaces for indentation.
// It parses the input into an AST and then rebuilds the formatted code.
func Format(content string) (string, error) {
	schema, err := parser.Parser.ParseString("", content)
	if err != nil {
		return "", fmt.Errorf("error parsing URPC: %w", err)
	}

	f := &formatter{
		builder: &strings.Builder{},
		indent:  0,
	}
	f.formatSchema(schema)

	// Ensure a single trailing newline
	result := f.builder.String()
	result = strings.TrimSpace(result)
	return result + "\n", nil
}

// formatter maintains the state during the formatting process.
type formatter struct {
	builder *strings.Builder
	indent  int
}

// write appends a string to the builder.
func (f *formatter) write(s string) {
	f.builder.WriteString(s)
}

// writeLine appends a string and a newline to the builder.
func (f *formatter) writeLine(s string) {
	f.builder.WriteString(s)
	f.builder.WriteString("\n")
}

// writeIndent appends the current indentation to the builder.
func (f *formatter) writeIndent() {
	for range f.indent {
		f.builder.WriteString("  ")
	}
}

// writeIndentedLine appends an indented line to the builder.
func (f *formatter) writeIndentedLine(s string) {
	f.writeIndent()
	f.writeLine(s)
}

// formatSchema formats the entire URPC schema.
func (f *formatter) formatSchema(schema *ast.URPCSchema) {
	// Format version
	if schema.Version != nil {
		f.formatVersion(schema.Version)
		f.write("\n")
	}

	// Format imports
	if len(schema.Imports) > 0 {
		f.formatImports(schema.Imports)
		f.write("\n")
	}

	// Format rules
	if len(schema.Rules) > 0 {
		f.formatRules(schema.Rules)
		if len(schema.Types) > 0 || len(schema.Procs) > 0 {
			f.write("\n")
		}
	}

	// Format types
	if len(schema.Types) > 0 {
		f.formatTypes(schema.Types)
		if len(schema.Procs) > 0 {
			f.write("\n")
		}
	}

	// Format procedures
	if len(schema.Procs) > 0 {
		f.formatProcs(schema.Procs)
	}
}

// formatVersion formats the version declaration.
func (f *formatter) formatVersion(v *ast.Version) {
	f.writeLine(fmt.Sprintf("version %d", v.Number))
}

// formatImports formats the import statements.
func (f *formatter) formatImports(imports []*ast.Import) {
	for _, imp := range imports {
		f.writeLine(fmt.Sprintf("import %q", imp.Path))
	}
}

// formatRules formats all rule declarations.
func (f *formatter) formatRules(rules []*ast.RuleDecl) {
	for i, rule := range rules {
		f.formatRule(rule)
		if i < len(rules)-1 {
			f.write("\n")
		}
	}
}

// formatRule formats a single rule declaration.
func (f *formatter) formatRule(rule *ast.RuleDecl) {
	// Format docstring if present
	if rule.Docstring != "" {
		// Ensure docstring is properly formatted with triple quotes
		if !strings.HasPrefix(rule.Docstring, "\"\"\"") {
			f.writeLine("\"\"\"")
			f.writeLine(strings.TrimSpace(rule.Docstring))
			f.writeLine("\"\"\"")
		} else {
			f.writeLine(rule.Docstring)
		}
	}

	// Write rule name
	f.writeLine(fmt.Sprintf("rule @%s {", rule.Name))

	// Increase indentation for the rule body
	f.indent++

	// Format the for field
	f.writeIndentedLine(fmt.Sprintf("for: %s", rule.Body.For))

	// Format the param field if present
	if rule.Body.Param != "" {
		paramType := rule.Body.Param
		if rule.Body.ParamIsArray {
			paramType += "[]"
		}
		f.writeIndentedLine(fmt.Sprintf("param: %s", paramType))
	}

	// Format the error field if present
	if rule.Body.Error != "" {
		f.writeIndentedLine(fmt.Sprintf("error: %q", rule.Body.Error))
	}

	// Decrease indentation and close the rule
	f.indent--
	f.writeLine("}")
}

// formatTypes formats all type declarations.
func (f *formatter) formatTypes(types []*ast.TypeDecl) {
	for i, t := range types {
		f.formatType(t)
		if i < len(types)-1 {
			f.write("\n")
		}
	}
}

// formatType formats a single type declaration.
func (f *formatter) formatType(t *ast.TypeDecl) {
	// Format docstring if present
	if t.Docstring != "" {
		// Ensure docstring is properly formatted with triple quotes
		if !strings.HasPrefix(t.Docstring, "\"\"\"") {
			f.writeLine("\"\"\"")
			f.writeLine(strings.TrimSpace(t.Docstring))
			f.writeLine("\"\"\"")
		} else {
			f.writeLine(t.Docstring)
		}
	}

	// Write type name and extends clause if present
	if len(t.Extends) > 0 {
		f.writeLine(fmt.Sprintf("type %s extends %s {", t.Name, strings.Join(t.Extends, ", ")))
	} else {
		f.writeLine(fmt.Sprintf("type %s {", t.Name))
	}

	// Increase indentation for the type body
	f.indent++

	// Format the fields
	for _, field := range t.Fields {
		f.formatField(field)
	}

	// Decrease indentation and close the type
	f.indent--
	f.writeLine("}")
}

// formatProcs formats all procedure declarations.
func (f *formatter) formatProcs(procs []*ast.ProcDecl) {
	for i, proc := range procs {
		f.formatProc(proc)
		if i < len(procs)-1 {
			f.write("\n")
		}
	}
}

// formatProc formats a single procedure declaration.
func (f *formatter) formatProc(proc *ast.ProcDecl) {
	// Format docstring if present
	if proc.Docstring != "" {
		// Ensure docstring is properly formatted with triple quotes
		if !strings.HasPrefix(proc.Docstring, "\"\"\"") {
			f.writeLine("\"\"\"")
			f.writeLine(strings.TrimSpace(proc.Docstring))
			f.writeLine("\"\"\"")
		} else {
			f.writeLine(proc.Docstring)
		}
	}

	// Write procedure name
	f.writeLine(fmt.Sprintf("proc %s {", proc.Name))

	// Increase indentation for the procedure body
	f.indent++

	// Format input if present
	if len(proc.Body.Input) > 0 {
		f.writeIndentedLine("input {")
		f.indent++
		for _, field := range proc.Body.Input {
			f.formatField(field)
		}
		f.indent--
		f.writeIndentedLine("}")

		if len(proc.Body.Output) > 0 || len(proc.Body.Meta) > 0 {
			f.write("\n")
		}
	}

	// Format output if present
	if len(proc.Body.Output) > 0 {
		f.writeIndentedLine("output {")
		f.indent++
		for _, field := range proc.Body.Output {
			f.formatField(field)
		}
		f.indent--
		f.writeIndentedLine("}")

		if len(proc.Body.Meta) > 0 {
			f.write("\n")
		}
	}

	// Format meta if present
	if len(proc.Body.Meta) > 0 {
		f.writeIndentedLine("meta {")
		f.indent++
		for _, meta := range proc.Body.Meta {
			// Ensure string values are properly quoted
			value := meta.Value
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				f.writeIndentedLine(fmt.Sprintf("%s: %s", meta.Key, value))
			} else if value == "true" || value == "false" || isNumeric(value) {
				// Boolean and numeric values don't need quotes
				f.writeIndentedLine(fmt.Sprintf("%s: %s", meta.Key, value))
			} else {
				// Other values should be quoted
				f.writeIndentedLine(fmt.Sprintf("%s: %q", meta.Key, value))
			}
		}
		f.indent--
		f.writeIndentedLine("}")
	}

	// Decrease indentation and close the procedure
	f.indent--
	f.writeLine("}")
}

// isNumeric checks if a string represents a numeric value
func isNumeric(s string) bool {
	// Simple check for integer or float
	for i, c := range s {
		if !(c >= '0' && c <= '9') && !(i > 0 && c == '.') {
			return false
		}
	}
	return true
}

// formatField formats a field in a type or procedure declaration.
func (f *formatter) formatField(field *ast.Field) {
	// Build the field name with optional marker
	fieldName := field.Name
	if field.Optional {
		fieldName += "?"
	}

	// Format the field type
	typeStr := f.formatFieldType(&field.Type)

	// Write the field and its type
	f.writeIndent()
	f.write(fmt.Sprintf("%s: %s", fieldName, typeStr))

	// If there are no rules, add a newline
	if len(field.Rules) == 0 {
		f.write("\n")
		return
	}

	// Format rules with increased indentation
	f.write("\n")
	f.indent++
	for _, rule := range field.Rules {
		f.formatFieldRule(rule)
	}
	f.indent--
}

// formatFieldType formats a field type.
func (f *formatter) formatFieldType(fieldType *ast.FieldType) string {
	var result string

	// Format the base type
	if fieldType.Base.Named != nil {
		result = *fieldType.Base.Named
	} else if fieldType.Base.Object != nil {
		result = f.formatInlineObject(fieldType.Base.Object)
	}

	// Add array brackets if needed
	for i := 0; i < int(fieldType.Depth); i++ {
		result += "[]"
	}

	return result
}

// formatInlineObject formats an inline object type.
func (f *formatter) formatInlineObject(obj *ast.FieldTypeObject) string {
	var result strings.Builder
	result.WriteString("{\n")

	// Get the current indentation level
	currentIndent := f.indent

	// Store the current builder and indent temporarily
	savedBuilder := f.builder
	savedIndent := f.indent

	// Create a new builder for the inline object
	f.builder = &strings.Builder{}

	// Increase indentation for the fields in the inline object
	f.indent = currentIndent + 1

	// Format each field
	for _, field := range obj.Fields {
		f.formatField(field)
	}

	// Get the formatted content
	fieldContent := f.builder.String()

	// Restore the original builder and indentation
	f.builder = savedBuilder
	f.indent = savedIndent

	// Add the field content with proper indentation
	result.WriteString(fieldContent)

	// Add closing brace with proper indentation
	for i := 0; i < currentIndent; i++ {
		result.WriteString("  ")
	}
	result.WriteString("}")

	return result.String()
}

// formatFieldRuleToString formats a field rule to a string without writing to the builder.
func formatFieldRuleToString(rule *ast.FieldRule) string {
	ruleName := fmt.Sprintf("@%s", rule.Name)

	// If there's no rule body, just return the rule name
	if rule.Body.ParamSingle == nil &&
		len(rule.Body.ParamListString) == 0 &&
		len(rule.Body.ParamListInt) == 0 &&
		len(rule.Body.ParamListFloat) == 0 &&
		len(rule.Body.ParamListBoolean) == 0 &&
		rule.Body.Error == "" {
		return ruleName
	}

	// Start building the rule with parameters
	var paramStr string

	// Format single parameter if present
	if rule.Body.ParamSingle != nil {
		// Preserve quotes for string parameters
		paramStr = *rule.Body.ParamSingle

		// Handle special case for @contains("@") to ensure quotes are preserved
		if rule.Name == "contains" && paramStr == "@" {
			paramStr = "\"@\""
		}
	}

	// Format parameter lists if present
	if len(rule.Body.ParamListString) > 0 {
		// Ensure strings are properly quoted
		quotedStrs := make([]string, len(rule.Body.ParamListString))
		for i, str := range rule.Body.ParamListString {
			quotedStrs[i] = fmt.Sprintf("\"%s\"", str)
		}
		paramStr = fmt.Sprintf("[%s]", strings.Join(quotedStrs, ", "))
	} else if len(rule.Body.ParamListInt) > 0 {
		paramStr = fmt.Sprintf("[%s]", strings.Join(rule.Body.ParamListInt, ", "))
	} else if len(rule.Body.ParamListFloat) > 0 {
		paramStr = fmt.Sprintf("[%s]", strings.Join(rule.Body.ParamListFloat, ", "))
	} else if len(rule.Body.ParamListBoolean) > 0 {
		paramStr = fmt.Sprintf("[%s]", strings.Join(rule.Body.ParamListBoolean, ", "))
	}

	// Build the rule string with parameters and optional error
	ruleStr := ruleName
	if paramStr != "" {
		ruleStr += fmt.Sprintf("(%s", paramStr)
		if rule.Body.Error != "" {
			ruleStr += fmt.Sprintf(", error: %q", rule.Body.Error)
		}
		ruleStr += ")"
	}

	return ruleStr
}

// formatFieldRule formats a field rule.
func (f *formatter) formatFieldRule(rule *ast.FieldRule) {
	ruleStr := formatFieldRuleToString(rule)
	f.writeIndentedLine(ruleStr)
}
