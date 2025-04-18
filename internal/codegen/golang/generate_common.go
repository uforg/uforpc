package golang

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

type generateCommonRenderFieldParams struct {
	// The field to render
	field schema.FieldDefinition
	// Whether to only return the type or the full field
	typeOnly bool
	// Whether to omit the JSON tags
	omitTag bool
}

// generateCommonRenderField generates the code for a field
func generateCommonRenderField(params generateCommonRenderFieldParams) string {
	field := params.field
	name := field.Name
	typeOnly := params.typeOnly
	omitTag := params.omitTag

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
		og.Inline("struct {")
		og.Block(func() {
			for _, fieldDef := range field.TypeInline.Fields {
				og.Line(generateCommonRenderField(generateCommonRenderFieldParams{
					field:    fieldDef,
					typeOnly: false,
					omitTag:  false,
				}))
			}
		})
		og.Line("}")
		typeLiteral = og.String()
	}

	if field.Depth > 0 {
		typeLiteral = strings.Repeat("[]", field.Depth) + typeLiteral
	}

	if isOptional {
		switch typeLiteral {
		case "string":
			typeLiteral = "NullString"
		case "int":
			typeLiteral = "NullInt"
		case "float64":
			typeLiteral = "NullFloat64"
		case "bool":
			typeLiteral = "NullBool"
		case "time.Time":
			typeLiteral = "NullTime"
		default:
			typeLiteral = fmt.Sprintf("Null[%s]", typeLiteral)
		}
	}

	if typeOnly {
		return typeLiteral
	}

	jsonTag := ""
	if !omitTag {
		jsonTag = fmt.Sprintf(" `json:\"%s,omitempty,omitzero\"`", nameCamel)
	}

	result := fmt.Sprintf("%s %s", namePascal, typeLiteral)
	return result + jsonTag
}

// generateCommonRenderStructFromFieldSlice generates the code for a slice of fields
func generateCommonRenderStructFromFieldSlice(fieldSlice []schema.FieldDefinition) string {
	if len(fieldSlice) == 0 {
		return "struct{}"
	}

	og := genkit.NewGenKit().WithTabs()
	og.Inline("struct {")
	og.Block(func() {
		for _, fieldItem := range fieldSlice {
			og.Line(generateCommonRenderField(generateCommonRenderFieldParams{
				field:    fieldItem,
				typeOnly: false,
				omitTag:  false,
			}))
		}
	})
	og.Line("}")
	return og.String()
}
