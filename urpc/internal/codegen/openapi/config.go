package openapi

import (
	"fmt"
	"strings"
)

// Config is the configuration for the OpenAPI generator.
type Config struct {
	// OutputFile is the file to output the generated code to.
	OutputFile string `toml:"output_file"`
	// Title is the title of the OpenAPI spec.
	Title string `toml:"title"`
	// Description is the description of the OpenAPI spec.
	Description string `toml:"description"`
	// BaseURL is the base URL to use for the OpenAPI spec.
	BaseURL string `toml:"base_url"`
}

func (c Config) Validate() error {
	if c.OutputFile == "" {
		return fmt.Errorf(`"output_file" is required`)
	}
	if !strings.HasSuffix(c.OutputFile, ".json") && !strings.HasSuffix(c.OutputFile, ".yaml") && !strings.HasSuffix(c.OutputFile, ".yml") {
		return fmt.Errorf(`"output_file" must end with ".json", ".yaml" or ".yml"`)
	}
	return nil
}
