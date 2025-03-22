package golang

import (
	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

// Generate takes a schema and a config and generates the Go code for the schema.
func Generate(sch schema.Schema, config Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	subGenerators := []func(*genkit.GenKit, schema.Schema, Config) error{
		generatePackage,
		generateCoreTypes,
		generateDomainTypes,
		generateNullUtils,
		generateValidator,
	}

	for _, generator := range subGenerators {
		if err := generator(g, sch, config); err != nil {
			return "", err
		}
		g.Break()
	}

	return g.String(), nil
}
