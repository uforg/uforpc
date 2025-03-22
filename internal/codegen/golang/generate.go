package golang

import (
	"strings"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

// Generate takes a schema and a config and generates the Go code for the schema.
func Generate(sch schema.Schema, config Config) (string, error) {
	subGenerators := []func(schema.Schema, Config) (string, error){
		generatePackage,
		generateCoreTypes,
		generateDomainTypes,
		generateNullUtils,
		generateValidator,
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

	return g.String(), nil
}
