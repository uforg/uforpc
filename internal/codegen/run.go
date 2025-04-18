package codegen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uforg/uforpc/internal/codegen/golang"
	"github.com/uforg/uforpc/internal/codegen/typescript"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/transpile"
	"github.com/uforg/uforpc/internal/urpc/analyzer"
	"github.com/uforg/uforpc/internal/urpc/docstore"
	"github.com/uforg/uforpc/internal/urpc/typeflattener"
	"github.com/uforg/uforpc/internal/util/filepathutil"
)

// Run runs the code generator and returns an error if one occurred.
func Run(configPath string) error {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read %s config file: %s", configPath, err)
	}

	config := Config{}
	if err := config.UnmarshalAndValidate(configBytes); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	/////////////////////////////
	// ANALYZE THE URPC SCHEMA //
	/////////////////////////////

	absSchemaPath, err := filepathutil.NormalizeFromWD(config.Schema)
	if err != nil {
		return fmt.Errorf("failed to normalize schema path: %w", err)
	}

	an, err := analyzer.NewAnalyzer(docstore.NewDocstore())
	if err != nil {
		return fmt.Errorf("failed to create URPC analyzer: %w", err)
	}

	combinedSchema, _, err := an.Analyze(absSchemaPath)
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	flattenedSchema := typeflattener.Flatten(combinedSchema.Schema)

	///////////////////////
	// TRANSPILE TO JSON //
	///////////////////////

	jsonSchema, err := transpile.ToJSON(*flattenedSchema)
	if err != nil {
		return fmt.Errorf("failed to transpile schema to its JSON representation: %w", err)
	}

	/////////////////////////
	// RUN CODE GENERATORS //
	/////////////////////////

	if config.HasGolang() {
		if err := runGolang(config.Golang, jsonSchema); err != nil {
			return fmt.Errorf("failed to run golang code generator: %w", err)
		}
	}

	if config.HasTypescript() {
		if err := runTypescript(config.Typescript, jsonSchema); err != nil {
			return fmt.Errorf("failed to run typescript code generator: %w", err)
		}
	}

	return nil
}

func runGolang(config *golang.Config, schema schema.Schema) error {
	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the code
	code, err := golang.Generate(schema, *config)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Write the code to the output directory
	filePath := filepath.Join(config.OutputDir, "uforpc_gen.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	return nil
}

func runTypescript(_ *typescript.Config, _ schema.Schema) error {
	return nil
}
