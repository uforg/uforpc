package golang

import (
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

func generateDomainTypes(sch schema.Schema, config Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Domain Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	for _, typeNode := range sch.GetTypeNodes() {
		desc := "is a domain type defined in UFO RPC with no documentation."
		if typeNode.Doc != nil {
			desc = strings.TrimSpace(*typeNode.Doc)
		}

		g.Linef("/* %s %s */", typeNode.Name, desc)
		g.Linef("type %s struct {", typeNode.Name)

		g.Block(func() {
			for _, fieldDef := range typeNode.Fields {
				g.Line(generateCommonRenderField(generateCommonRenderFieldParams{
					field:    fieldDef,
					typeOnly: false,
					omitTag:  false,
				}))
			}
		})

		g.Line("}")
		g.Break()

		g.Linef("// Null%s is the nullable version of %s", typeNode.Name, typeNode.Name)
		g.Linef("type Null%s = Null[%s]", typeNode.Name, typeNode.Name)
		g.Break()
	}

	return g.String(), nil
}
