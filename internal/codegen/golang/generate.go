package golang

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

// Generate takes a schema and a config and generates the Go code for the schema.
func Generate(sch schema.Schema, config Config) (string, error) {
	subGenerators := []func(schema.Schema, Config) (string, error){
		generatePackage,
		generateCoreTypes,
		generateDomainTypes,
		generateProcedureTypes,
		generateNullUtils,
		generateValidator,
		generateServer,
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
	formattedCode, err := format.Source([]byte(generatedCode))
	if err != nil {
		return "", fmt.Errorf("failed to format generated code: %w", err)
	}

	return string(formattedCode), nil
}
