package golang

import (
	"fmt"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

type generateCommonRenderFieldParams struct {
	// The name of the field
	name string
	// The field to render
	field schema.Field
	// Whether to only return the type or the full field
	typeOnly bool
	// Whether to omit the JSON tags
	omitTag bool
}

// generateCommonRenderField generates the code for a field
func generateCommonRenderField(params generateCommonRenderFieldParams) string {
	name := params.name
	typeOnly := params.typeOnly
	field := params.field
	omitTag := params.omitTag

	// Protect against empty fields
	if field.Type == "" {
		return ""
	}

	namePascal := strutil.ToPascalCase(name)
	nameCamel := strutil.ToCamelCase(name)
	isOptional := field.Optional
	isCustomType := field.IsCustomType()
	isBuiltInType := field.IsBuiltInType()

	typeLiteral := "any"
	if isCustomType {
		typeLiteral = field.Type
	}
	if isBuiltInType {
		switch field.Type {
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
				for fieldName, fieldContent := range field.Fields {
					og.Line(generateCommonRenderField(generateCommonRenderFieldParams{
						name:     fieldName,
						field:    fieldContent,
						typeOnly: false,
						omitTag:  false,
					}))
				}
			})
			og.Line("}")
			typeLiteral = og.String()
		case "array":
			if field.ArrayType != nil {
				underlyingType := generateCommonRenderField(generateCommonRenderFieldParams{
					name:     "",
					field:    *field.ArrayType,
					typeOnly: true,
					omitTag:  true,
				})
				typeLiteral = fmt.Sprintf("[]%s", underlyingType)
			}
		}
	}

	if isOptional {
		switch field.Type {
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

	if typeOnly {
		return typeLiteral
	}

	description := ""
	if field.Description != "" {
		description = fmt.Sprintf("// %s\n", field.Description)
	}

	jsonTag := ""
	if !omitTag {
		jsonTag = fmt.Sprintf(" `json:\"%s,omitempty,omitzero\"`", nameCamel)
	}

	result := fmt.Sprintf("%s %s", namePascal, typeLiteral)
	return description + result + jsonTag
}
