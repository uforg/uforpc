package golang

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

// renderField generates the code for a field
func renderField(parentTypeName string, field schema.FieldDefinition) string {
	name := field.Name
	isNamed := field.IsNamed()
	isInline := field.IsInline()

	// Protect against empty fields
	if !isNamed && !isInline {
		return ""
	}

	namePascal := strutil.ToPascalCase(name)
	nameCamel := strutil.ToCamelCase(name)
	isOptional := field.Optional
	isCustomType := field.IsCustomType()
	isBuiltInType := field.IsBuiltInType()

	typeLiteral := "any"

	if isNamed && isCustomType {
		typeLiteral = *field.TypeName
	}

	if isNamed && isBuiltInType {
		switch *field.TypeName {
		case "string":
			typeLiteral = "string"
		case "int":
			typeLiteral = "int"
		case "float":
			typeLiteral = "float64"
		case "boolean":
			typeLiteral = "bool"
		case "datetime":
			typeLiteral = "time.Time"
		}
	}

	if isInline {
		typeLiteral = parentTypeName + namePascal
	}

	if field.Depth > 0 {
		typeLiteral = strings.Repeat("[]", field.Depth) + typeLiteral
	}

	if isOptional {
		typeLiteral = fmt.Sprintf("Optional[%s]", typeLiteral)
	}

	jsonTag := fmt.Sprintf(" `json:\"%s,omitempty,omitzero\"`", nameCamel)
	result := fmt.Sprintf("%s %s", namePascal, typeLiteral)
	return result + jsonTag
}

// renderType renders a type definition with all its fields
func renderType(
	parentName string,
	name string,
	desc string,
	fields []schema.FieldDefinition,
) string {
	name = parentName + name

	og := genkit.NewGenKit().WithTabs()
	if desc != "" {
		og.Linef("/* %s %s */", name, desc)
	}
	og.Linef("type %s struct {", name)
	og.Block(func() {
		for _, fieldDef := range fields {
			og.Line(renderField(name, fieldDef))
		}
	})
	og.Line("}")
	og.Break()

	// Render children inline types
	for _, fieldDef := range fields {
		if !fieldDef.IsInline() {
			continue
		}

		og.Line(renderType(name, strutil.ToPascalCase(fieldDef.Name), "", fieldDef.TypeInline.Fields))
	}

	return og.String()
}

// renderPreField generates the code for a field in a pre type
func renderPreField(parentTypeName string, field schema.FieldDefinition) string {
	name := field.Name
	isNamed := field.IsNamed()
	isInline := field.IsInline()

	// Protect against empty fields
	if !isNamed && !isInline {
		return ""
	}

	namePascal := strutil.ToPascalCase(name)
	nameCamel := strutil.ToCamelCase(name)
	isCustomType := field.IsCustomType()
	isBuiltInType := field.IsBuiltInType()

	typeLiteral := "any"

	if isNamed && isCustomType {
		typeLiteral = "pre" + *field.TypeName
	}

	if isNamed && isBuiltInType {
		switch *field.TypeName {
		case "string":
			typeLiteral = "string"
		case "int":
			typeLiteral = "int"
		case "float":
			typeLiteral = "float64"
		case "boolean":
			typeLiteral = "bool"
		case "datetime":
			typeLiteral = "time.Time"
		}
	}

	if isInline {
		typeLiteral = "pre" + parentTypeName + namePascal
	}

	if field.Depth > 0 {
		typeLiteral = strings.Repeat("[]", field.Depth) + typeLiteral
	}

	typeLiteral = fmt.Sprintf("Optional[%s]", typeLiteral)

	jsonTag := fmt.Sprintf(" `json:\"%s,omitempty,omitzero\"`", nameCamel)
	result := fmt.Sprintf("%s %s", namePascal, typeLiteral)
	return result + jsonTag
}

// renderPreType renders a type definition with all its fields marked as optional
// and helpers to validate the required fields and transform to the final type
func renderPreType(
	parentName string,
	name string,
	fields []schema.FieldDefinition,
) string {
	name = parentName + name

	og := genkit.NewGenKit().WithTabs()
	og.Linef("// pre%s is the version of %s previous to the required field validation", name, name)
	og.Linef("type pre%s struct {", name)
	og.Block(func() {
		for _, fieldDef := range fields {
			og.Line(renderPreField(name, fieldDef))
		}
	})
	og.Line("}")
	og.Break()

	// Render children inline types
	for _, fieldDef := range fields {
		if !fieldDef.IsInline() {
			continue
		}

		og.Line(renderPreType(name, strutil.ToPascalCase(fieldDef.Name), fieldDef.TypeInline.Fields))
	}

	// Render validate function
	og.Linef("// validate validates the required fields of %s", name)
	og.Linef("func (p *pre%s) validate() error {", name)
	og.Block(func() {
		og.Line("if p == nil {")
		og.Block(func() {
			og.Linef("return errorMissingRequiredField(\"pre%s is nil\")", name)
		})
		og.Line("}")
		og.Break()

		// Required fields
		for _, fieldDef := range fields {
			fieldName := strutil.ToPascalCase(fieldDef.Name)
			isRequired := !fieldDef.Optional
			isCustomType := fieldDef.IsCustomType()
			isInline := fieldDef.IsInline()
			isArray := fieldDef.Depth > 0
			arrayDepth := fieldDef.Depth

			og.Linef(`// Required validations for field "%s"`, fieldDef.Name)

			if isRequired {
				og.Linef("if !p.%s.Present {", fieldName)
				og.Block(func() {
					og.Linef("return errorMissingRequiredField(\"field %s is required\")", fieldDef.Name)
				})
				og.Line("}")
			}

			if (isCustomType || isInline) && !isArray {
				og.Linef("if p.%s.Present {", fieldName)
				og.Block(func() {
					og.Linef("if err := p.%s.Value.validate(); err != nil {", fieldName)
					og.Block(func() {
						og.Linef("return errorMissingRequiredField(\"field %s: \" + err.Error())", fieldDef.Name)
					})
					og.Line("}")
				})
				og.Line("}")
			}

			if (isCustomType || isInline) && isArray {
				og.Linef("if p.%s.Present {", fieldName)
				og.Block(func() {
					og.Linef("item%d := p.%s.Value", arrayDepth, fieldName)

					var renderLevel func(arrayDepth int)
					renderLevel = func(arrayDepth int) {
						nextArrayDepth := arrayDepth - 1

						if arrayDepth > 1 {
							og.Linef("for _, item%d := range item%d {", nextArrayDepth, arrayDepth)
							og.Block(func() {
								renderLevel(nextArrayDepth)
							})
							og.Line("}")
							return
						}

						og.Linef("for _, item := range item%d {", arrayDepth)
						og.Block(func() {
							og.Linef("if err := item.validate(); err != nil {")
							og.Block(func() {
								og.Linef("return errorMissingRequiredField(\"field %s: \" + err.Error())", fieldDef.Name)
							})
							og.Line("}")
						})
						og.Line("}")
					}

					renderLevel(arrayDepth)
				})
				og.Line("}")
			}

			og.Break()
		}

		og.Line("return nil")
	})
	og.Line("}")
	og.Break()

	return og.String()
}
