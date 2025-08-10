package typescript

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
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
			typeLiteral = "number"
		case "float":
			typeLiteral = "number"
		case "bool":
			typeLiteral = "boolean"
		case "datetime":
			typeLiteral = "Date"
		}
	}

	if isInline {
		typeLiteral = parentTypeName + namePascal
	}

	if field.IsArray {
		typeLiteral = fmt.Sprintf("%s[]", typeLiteral)
	}

	finalName := nameCamel
	if isOptional {
		finalName += "?"
	}

	return fmt.Sprintf("%s: %s", finalName, typeLiteral)
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
		og.Linef("/**")
		renderPartialMultilineComment(og, fmt.Sprintf("%s %s", name, desc))
		og.Linef(" */")
	}
	og.Linef("export type %s = {", name)
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

// renderHydrateField generates the code for a field in a hydrate type
func renderHydrateField(parentTypeName string, field schema.FieldDefinition) string {
	name := field.Name
	isNamed := field.IsNamed()
	isInline := field.IsInline()

	// Protect against empty fields
	if !isNamed && !isInline {
		return ""
	}

	nameHydrated := "hydrated" + strutil.ToPascalCase(name)
	namePascal := strutil.ToPascalCase(name)
	isOptional := field.Optional
	isCustomType := field.IsCustomType()
	isBuiltInType := field.IsBuiltInType()

	valueLiteral := "input." + name

	if isNamed && isCustomType {
		valueLiteral = fmt.Sprintf("hydrate%s(input.%s)", namePascal, name)
	}

	if isNamed && isBuiltInType {
		switch *field.TypeName {
		case "string":
			valueLiteral = "input." + name
		case "int":
			valueLiteral = "input." + name
		case "float":
			valueLiteral = "input." + name
		case "bool":
			valueLiteral = "input." + name
		case "datetime":
			valueLiteral = "new Date(input." + name + ")"
		}
	}

	if isInline {
		valueLiteral = fmt.Sprintf("hydrate%s%s(input.%s)", parentTypeName, namePascal, name)
	}

	if field.IsArray {
		valueLiteral = fmt.Sprintf("input.%s.map(%s => {const input = {%s}; return %s})", name, name, name, valueLiteral)
	}

	if isOptional {
		valueLiteral = fmt.Sprintf("input.%s ? %s : input.%s", name, valueLiteral, name)
	}

	return fmt.Sprintf("const %s = %s", nameHydrated, valueLiteral)
}

// renderHydrateType renders a function used to transform a type returned from JSON.parse to it's
// final type.
func renderHydrateType(parentName string, name string, fields []schema.FieldDefinition) string {
	name = parentName + name

	og := genkit.NewGenKit().WithTabs()
	og.Linef("function hydrate%s(input: %s): %s {", name, name, name)
	og.Block(func() {
		for _, fieldDef := range fields {
			og.Line(renderHydrateField(name, fieldDef))
		}
		og.Linef("return {")
		og.Block(func() {
			for _, fieldDef := range fields {
				name := fieldDef.Name
				hydratedName := "hydrated" + strutil.ToPascalCase(name)
				og.Linef("%s: %s,", name, hydratedName)
			}
		})
		og.Linef("}")
	})
	og.Line("}")
	og.Break()

	// Render children inline types
	for _, fieldDef := range fields {
		if !fieldDef.IsInline() {
			continue
		}

		og.Line(renderHydrateType(name, strutil.ToPascalCase(fieldDef.Name), fieldDef.TypeInline.Fields))
	}

	return og.String()
}

// renderPartialMultilineComment receives a text and renders it to the given genkit.GenKit
// as a partial multiline comment.
func renderPartialMultilineComment(g *genkit.GenKit, text string) {
	for line := range strings.SplitSeq(text, "\n") {
		g.Linef(" * %s", line)
	}
}

// renderDeprecated receives a pointer to a string and if it is not nil, it will
// render a comment with the deprecated message to the given genkit.GenKit.
func renderDeprecated(g *genkit.GenKit, deprecated *string) {
	if deprecated == nil {
		return
	}

	desc := "@deprecated "
	if *deprecated == "" {
		desc += "This is deprecated and should not be used in new code."
	} else {
		desc += *deprecated
	}

	g.Line(" *")
	renderPartialMultilineComment(g, desc)
}
