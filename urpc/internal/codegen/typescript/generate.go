package typescript

import (
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

// Generate takes a schema and a config and generates the TypeScript code for the schema.
func Generate(sch schema.Schema, config Config) (string, error) {
	subGenerators := []func(schema.Schema, Config) (string, error){
		generateCoreTypes,
		generateDomainTypes,
		generateProcedureTypes,
		generateStreamTypes,
		generateClient,
	}

	g := genkit.NewGenKit().WithTabs()
	for _, generator := range subGenerators {
		codeChunk, err := generator(sch, config)
		if err != nil {
			return "", err
		}

		codeChunk = strings.TrimSpace(codeChunk)
		g.Raw(codeChunk)
		g.Break()
		g.Break()
	}

	generatedCode := g.String()
	generatedCode = strutil.LimitConsecutiveNewlines(generatedCode, 2)
	return generatedCode, nil
}
