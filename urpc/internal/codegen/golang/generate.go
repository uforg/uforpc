package golang

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"golang.org/x/tools/imports"
)

// Generate takes a schema and a config and generates the Go code for the schema.
func Generate(sch schema.Schema, config Config) (string, error) {
	subGenerators := []func(schema.Schema, Config) (string, error){
		generatePackage,
		generateCoreTypes,
		generateDomainTypes,
		generateProcedureTypes,
		generateOptional,
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
	formattedCode, err := imports.Process("", []byte(generatedCode), nil)
	if err != nil {
		return "", fmt.Errorf("failed to format generated code: %w", err)
	}

	return string(formattedCode), nil
}
