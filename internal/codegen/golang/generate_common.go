package golang

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

// renderField generates the code for a field
func renderField(field schema.FieldDefinition) string {
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
		og := genkit.NewGenKit().WithTabs()
		og.Line("struct {")
		og.Block(func() {
			for _, fieldDef := range field.TypeInline.Fields {
				og.Line(renderField(fieldDef))
			}
		})
		og.Inline("}")
		typeLiteral = og.String()
	}

	if field.Depth > 0 {
		typeLiteral = strings.Repeat("[]", field.Depth) + typeLiteral
	}

	if isOptional {
		if isCustomType {
			typeLiteral = *field.TypeName + "Optional"
		} else {
			switch typeLiteral {
			case "string":
				typeLiteral = "StringOptional"
			case "int":
				typeLiteral = "IntOptional"
			case "float64":
				typeLiteral = "Float64Optional"
			case "bool":
				typeLiteral = "BoolOptional"
			case "time.Time":
				typeLiteral = "TimeOptional"
			default:
				typeLiteral = fmt.Sprintf("Optional[%s]", strings.TrimSpace(typeLiteral))
			}
		}
	}

	jsonTag := fmt.Sprintf(" `json:\"%s,omitempty,omitzero\"`", nameCamel)
	result := fmt.Sprintf("%s %s", namePascal, typeLiteral)
	return result + jsonTag
}

// renderType renders a type definition with all its fields
func renderType(name string, desc string, fields []schema.FieldDefinition) string {
	og := genkit.NewGenKit().WithTabs()
	if desc != "" {
		og.Linef("/* %s %s */", name, desc)
	}
	og.Linef("type %s struct {", name)
	og.Block(func() {
		for _, fieldDef := range fields {
			og.Line(renderField(fieldDef))
		}
	})
	og.Line("}")

	return og.String()
}

// renderPreField generates the code for a field in a pre type
func renderPreField(field schema.FieldDefinition) string {
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
		typeLiteral = "pre" + *field.TypeName + "Optional"
	}

	if isNamed && isBuiltInType {
		switch *field.TypeName {
		case "string":
			typeLiteral = "StringOptional"
		case "int":
			typeLiteral = "IntOptional"
		case "float":
			typeLiteral = "Float64Optional"
		case "boolean":
			typeLiteral = "BoolOptional"
		case "datetime":
			typeLiteral = "TimeOptional"
		}
	}

	if isInline {
		og := genkit.NewGenKit().WithTabs()
		og.Line("struct {")
		og.Block(func() {
			for _, fieldDef := range field.TypeInline.Fields {
				og.Line(renderPreField(fieldDef))
			}
		})
		og.Inline("}")
		typeLiteral = og.String()
	}

	if field.Depth > 0 {
		typeLiteral = strings.Repeat("[]", field.Depth) + typeLiteral
	}

	jsonTag := fmt.Sprintf(" `json:\"%s,omitempty,omitzero\"`", nameCamel)
	result := fmt.Sprintf("%s %s", namePascal, typeLiteral)
	return result + jsonTag
}

// renderPreType renders a type definition with all its fields marked as optional
// and helpers to validate the required fields and transform to the final type
func renderPreType(name string, fields []schema.FieldDefinition) string {
	og := genkit.NewGenKit().WithTabs()
	og.Linef("// pre%s is the version of %s previous to the required field validation", name, name)
	og.Linef("type pre%s struct {", name)
	og.Block(func() {
		for _, fieldDef := range fields {
			og.Line(renderPreField(fieldDef))
		}
	})
	og.Line("}")
	og.Break()

	// Pre optional
	og.Linef("// pre%sOptional is the optional version of pre%s", name, name)
	og.Linef("type pre%sOptional = Optional[pre%s]", name, name)
	og.Break()

	// Validate function
	og.Linef("// validate validates the required fields of %s", name)
	og.Linef("func (p *pre%s) validate() error {", name)
	og.Block(func() {
		og.Line("if p == nil {")
		og.Block(func() {
			og.Linef("return fmt.Errorf(\"pre%s is nil\")", name)
		})
		og.Line("}")
		og.Break()

		// Required fields
		for _, fieldDef := range fields {
			if fieldDef.Optional {
				continue
			}

			fieldName := strutil.ToPascalCase(fieldDef.Name)
			og.Linef("if !p.%s.Present {", fieldName)
			og.Block(func() {
				og.Linef("return fmt.Errorf(\"%s is required\")", fieldDef.Name)
			})
			og.Line("}")
		}
		og.Break()

		// Deep validation for custom types
		for _, fieldDef := range fields {
			if !fieldDef.IsCustomType() {
				continue
			}

			fieldName := strutil.ToPascalCase(fieldDef.Name)
			og.Linef("if p.%s.Present {", fieldName)
			og.Block(func() {
				og.Linef("if err := p.%s.Value.validate(); err != nil {", fieldName)
				og.Block(func() {
					og.Linef("return fmt.Errorf(\"%s: %%w\", err)", fieldDef.Name)
				})
				og.Line("}")
			})
			og.Line("}")

		}
		og.Break()

		og.Line("return nil")
	})
	og.Line("}")
	og.Break()

	return og.String()
}
