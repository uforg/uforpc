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

		if typeNode.Deprecated != nil {
			desc += "\n\nDeprecated: "
			if *typeNode.Deprecated == "" {
				desc += "This type is deprecated and should not be used in new code."
			} else {
				desc += *typeNode.Deprecated
			}
		}

		g.Line(renderType(typeNode.Name, desc, typeNode.Fields))
		g.Break()

		g.Line(renderPreType(typeNode.Name, typeNode.Fields))
		g.Break()

		g.Linef("// %sOptional is the optional version of %s", typeNode.Name, typeNode.Name)
		g.Linef("type %sOptional = Optional[%s]", typeNode.Name, typeNode.Name)
		g.Break()

	}

	return g.String(), nil
}
