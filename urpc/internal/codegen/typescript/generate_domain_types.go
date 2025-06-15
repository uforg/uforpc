package typescript

import (
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
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

		if typeNode.Deprecated != nil {
			desc += "\n\n@deprecated "
			if *typeNode.Deprecated == "" {
				desc += "This type is deprecated and should not be used in new code."
			} else {
				desc += *typeNode.Deprecated
			}
		}

		g.Line(renderType("", typeNode.Name, desc, typeNode.Fields))
		g.Break()
	}

	return g.String(), nil
}
