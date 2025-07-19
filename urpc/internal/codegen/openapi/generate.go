package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
)

func Generate(schema *ast.Schema, config Config) (string, error) {
	if config.Title == "" {
		config.Title = "UFO RPC API"
	}

	spec := Spec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       config.Title,
			Description: config.Description,
			Version:     "1.0.0",
		},
	}

	if config.BaseURL != "" {
		spec.Servers = []Server{
			{
				URL: config.BaseURL,
			},
		}
	}

	code, err := encodeSpec(spec, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate spec file: %w", err)
	}

	return code, nil
}

func encodeSpec(spec Spec, config Config) (string, error) {
	isYAML := strings.HasSuffix(config.OutputFile, ".yaml") || strings.HasSuffix(config.OutputFile, ".yml")
	var buf bytes.Buffer

	if isYAML {
		enc := yaml.NewEncoder(&buf)
		if err := enc.Encode(spec); err != nil {
			return "", fmt.Errorf("failed to encode yaml spec: %w", err)
		}
		return buf.String(), nil
	}

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(spec); err != nil {
		return "", fmt.Errorf("failed to encode json spec: %w", err)
	}
	return buf.String(), nil
}
