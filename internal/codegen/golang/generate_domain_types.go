package golang

import (
	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

func generateDomainTypes(sch schema.Schema, config Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

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
				g.Line(generateCommonRenderField(generateCommonRenderFieldParams{
					name:     fieldName,
					field:    fieldContent,
					typeOnly: false,
					omitTag:  false,
				}))
			}
		})

		g.Line("}")
		g.Break()

		g.Linef("// Null%s is the nullable version of %s", typeName, typeName)
		g.Linef("type Null%s = Null[%s]", typeName, typeName)
		g.Break()
	}

	return g.String(), nil
}
