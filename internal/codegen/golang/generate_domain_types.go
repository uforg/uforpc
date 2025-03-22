package golang

import (
	"fmt"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

func generateDomainTypesRenderField(name string, content schema.Field) string {
	namePascal := strutil.ToPascalCase(name)
	nameCamel := strutil.ToCamelCase(name)
	isNameless := name == ""
	isOptional := content.ProcessedRules.IsOptional()
	isCustomType := content.IsCustomType()
	isBuiltInType := content.IsBuiltInType()

	typeLiteral := "any"
	if isCustomType {
		typeLiteral = content.Type
	}
	if isBuiltInType {
		switch content.Type {
		case "string":
			typeLiteral = "string"
		case "int":
			typeLiteral = "int"
		case "float":
			typeLiteral = "float64"
		case "boolean":
			typeLiteral = "bool"
		case "object":
			og := genkit.NewGenKit().WithTabs()
			og.Inline("struct {")
			og.Block(func() {
				for fieldName, fieldContent := range content.Fields {
					og.Line(generateDomainTypesRenderField(fieldName, fieldContent))
				}
			})
			og.Line("}")
			typeLiteral = og.String()
		case "array":
			if content.ArrayType != nil {
				underlyingType := generateDomainTypesRenderField("", *content.ArrayType)
				typeLiteral = fmt.Sprintf("[]%s", underlyingType)
			}
		}
	}

	if isOptional {
		switch content.Type {
		case "string":
			typeLiteral = "NullString"
		case "int":
			typeLiteral = "NullInt"
		case "float":
			typeLiteral = "NullFloat64"
		case "boolean":
			typeLiteral = "NullBool"
		default:
			typeLiteral = fmt.Sprintf("Null[%s]", typeLiteral)
		}
	}

	if isNameless {
		return typeLiteral
	}

	result := ""

	if content.Description != "" {
		result += fmt.Sprintf("// %s\n", content.Description)
	}

	result += fmt.Sprintf("%s %s `json:\"%s,omitempty,omitzero\"`", namePascal, typeLiteral, nameCamel)
	return result
}

func generateDomainTypes(g *genkit.GenKit, sch schema.Schema, config Config) error {
	g.Inline("// -----------------------------------------------------------------------------")
	g.Line("// Domain Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	for typeName, typeContent := range sch.Types {
		desc := typeContent.Description
		if desc == "" {
			desc = "is a domain type defined in UFO RPC with no description."
		}

		g.Linef("// %s %s", typeName, desc)
		g.Linef("type %s struct {", typeName)

		g.Block(func() {
			for fieldName, fieldContent := range typeContent.Fields {
				g.Line(generateDomainTypesRenderField(fieldName, fieldContent))
			}
		})

		g.Line("}")
		g.Break()

		g.Linef("// Null%s is the nullable version of %s", typeName, typeName)
		g.Linef("type Null%s Null[%s]", typeName, typeName)
		g.Break()
	}

	return nil
}
