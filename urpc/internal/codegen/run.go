package codegen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uforg/uforpc/urpc/internal/codegen/golang"
	"github.com/uforg/uforpc/urpc/internal/codegen/openapi"
	"github.com/uforg/uforpc/urpc/internal/codegen/playground"
	"github.com/uforg/uforpc/urpc/internal/codegen/typescript"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/transpile"
	"github.com/uforg/uforpc/urpc/internal/urpc/analyzer"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/urpc/docstore"
	"github.com/uforg/uforpc/urpc/internal/util/filepathutil"
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

	///////////////////////////////////////
	// PARSE AND ANALYZE THE URPC SCHEMA //
	///////////////////////////////////////

	absConfigPath, err := filepathutil.NormalizeFromWD(configPath)
	if err != nil {
		return fmt.Errorf("failed to normalize config path: %w", err)
	}

	absConfigDir := filepath.Dir(absConfigPath)
	absSchemaPath := filepath.Join(absConfigDir, config.Schema)

	an, err := analyzer.NewAnalyzer(docstore.NewDocstore())
	if err != nil {
		return fmt.Errorf("failed to create URPC analyzer: %w", err)
	}

	astSchema, _, err := an.Analyze(absSchemaPath)
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	///////////////////////
	// TRANSPILE TO JSON //
	///////////////////////

	jsonSchema, err := transpile.ToJSON(*astSchema)
	if err != nil {
		return fmt.Errorf("failed to transpile schema to its JSON representation: %w", err)
	}

	/////////////////////////
	// RUN CODE GENERATORS //
	/////////////////////////

	if config.HasOpenAPI() {
		if err := runOpenAPI(absConfigDir, config.OpenAPI, astSchema); err != nil {
			return fmt.Errorf("failed to run openapi code generator: %w", err)
		}
	}

	if config.HasPlayground() {
		if err := runPlayground(absConfigDir, config.Playground, astSchema); err != nil {
			return fmt.Errorf("failed to run playground code generator: %w", err)
		}
	}

	if config.HasGolang() {
		if err := runGolang(absConfigDir, config.Golang, jsonSchema); err != nil {
			return fmt.Errorf("failed to run golang code generator: %w", err)
		}
	}

	if config.HasTypescript() {
		if err := runTypescript(absConfigDir, config.Typescript, jsonSchema); err != nil {
			return fmt.Errorf("failed to run typescript code generator: %w", err)
		}
	}

	return nil
}

func runOpenAPI(absConfigDir string, config *openapi.Config, astSchema *ast.Schema) error {
	outputFile := filepath.Join(absConfigDir, config.OutputFile)
	outputDir := filepath.Dir(outputFile)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the code
	code, err := openapi.Generate(astSchema, *config)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	if err := os.WriteFile(outputFile, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	return nil
}

func runPlayground(absConfigDir string, config *playground.Config, astSchema *ast.Schema) error {
	outputDir := filepath.Join(absConfigDir, config.OutputDir)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the playground
	err := playground.Generate(absConfigDir, astSchema, *config)
	if err != nil {
		return fmt.Errorf("failed to generate playground: %w", err)
	}

	return nil
}

func runGolang(absConfigDir string, config *golang.Config, schema schema.Schema) error {
	outputFile := filepath.Join(absConfigDir, config.OutputFile)
	outputDir := filepath.Dir(outputFile)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the code
	code, err := golang.Generate(schema, *config)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Write the code to the output file
	if err := os.WriteFile(outputFile, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	return nil
}

func runTypescript(absConfigDir string, config *typescript.Config, schema schema.Schema) error {
	outputFile := filepath.Join(absConfigDir, config.OutputFile)
	outputDir := filepath.Dir(outputFile)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the code
	code, err := typescript.Generate(schema, *config)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Write the code to the output file
	if err := os.WriteFile(outputFile, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	return nil
}
