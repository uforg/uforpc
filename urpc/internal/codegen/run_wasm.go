package codegen

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/codegen/golang"
	"github.com/uforg/uforpc/urpc/internal/codegen/typescript"
	"github.com/uforg/uforpc/urpc/internal/transpile"
	"github.com/uforg/uforpc/urpc/internal/urpc/parser"
)

// RunWasmOptions contains options for running code generators in WASM mode
// without writing to files.
type RunWasmOptions struct {
	// Generator must be one of: "golang-server", "golang-client", "typescript-client".
	Generator string `json:"generator"`
	// SchemaInput is the schema content as a string (URPC schema only).
	SchemaInput string `json:"schemaInput"`
	// GolangPackageName is required when Generator is golang-server or golang-client.
	GolangPackageName string `json:"golangPackageName"`
}

// RunWasm executes a single generator and returns the generated code as a string.
func RunWasm(opts RunWasmOptions) (string, error) {
	if opts.Generator == "" {
		return "", fmt.Errorf("missing generator")
	}
	if opts.SchemaInput == "" {
		return "", fmt.Errorf("missing schema input")
	}

	// Parse input into JSON schema
	astSchema, err := parser.ParserInstance.ParseString("schema.urpc", opts.SchemaInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse URPC schema: %s", err)
	}
	jsonSchema, err := transpile.ToJSON(*astSchema)
	if err != nil {
		return "", fmt.Errorf("failed to transpile URPC to JSON: %s", err)
	}

	if opts.Generator == "golang-server" {
		if opts.GolangPackageName == "" {
			return "", fmt.Errorf("golang-server requires 'GolangPackageName'")
		}
		cfg := golang.Config{PackageName: opts.GolangPackageName, IncludeServer: true, IncludeClient: false}
		return golang.Generate(jsonSchema, cfg)
	}

	if opts.Generator == "golang-client" {
		if opts.GolangPackageName == "" {
			return "", fmt.Errorf("golang-client requires 'GolangPackageName'")
		}
		cfg := golang.Config{PackageName: opts.GolangPackageName, IncludeServer: false, IncludeClient: true}
		return golang.Generate(jsonSchema, cfg)
	}

	if opts.Generator == "typescript-client" {
		cfg := typescript.Config{IncludeServer: false, IncludeClient: true}
		return typescript.Generate(jsonSchema, cfg)
	}

	return "", fmt.Errorf("unsupported generator: %s", opts.Generator)
}
