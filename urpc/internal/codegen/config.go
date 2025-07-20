package codegen

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/uforg/uforpc/urpc/internal/codegen/golang"
	"github.com/uforg/uforpc/urpc/internal/codegen/openapi"
	"github.com/uforg/uforpc/urpc/internal/codegen/playground"
	"github.com/uforg/uforpc/urpc/internal/codegen/typescript"
)

// Config is the configuration for the code generator.
type Config struct {
	Version    int                `toml:"version"`
	Schema     string             `toml:"schema"`
	OpenAPI    *openapi.Config    `toml:"openapi"`
	Playground *playground.Config `toml:"playground"`
	Golang     *golang.Config     `toml:"golang"`
	Typescript *typescript.Config `toml:"typescript"`
}

func (c *Config) HasOpenAPI() bool {
	return c.OpenAPI != nil && c.OpenAPI.OutputFile != ""
}

func (c *Config) HasPlayground() bool {
	return c.Playground != nil
}

func (c *Config) HasGolang() bool {
	return c.Golang != nil
}

func (c *Config) HasTypescript() bool {
	return c.Typescript != nil
}

func (c *Config) Unmarshal(data []byte) error {
	if err := toml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to unmarshal TOML config: %w", err)
	}
	return nil
}

func (c *Config) Validate() error {
	if c.Version == 0 {
		return fmt.Errorf(`"version" is required`)
	}

	if c.Version != 1 {
		return fmt.Errorf("unsupported version: %d", c.Version)
	}

	if c.Schema == "" {
		return fmt.Errorf(`"schema" is required`)
	}

	if c.OpenAPI != nil {
		if err := c.OpenAPI.Validate(); err != nil {
			return fmt.Errorf("openapi config is invalid: %w", err)
		}
	}

	if c.Playground != nil {
		if err := c.Playground.Validate(); err != nil {
			return fmt.Errorf("playground config is invalid: %w", err)
		}
	}

	if c.Golang != nil {
		if err := c.Golang.Validate(); err != nil {
			return fmt.Errorf("golang config is invalid: %w", err)
		}
	}

	if c.Typescript != nil {
		if err := c.Typescript.Validate(); err != nil {
			return fmt.Errorf("typescript config is invalid: %w", err)
		}
	}

	return nil
}

// UnmarshalAndValidate unmarshals and validates a TOML config
func (c *Config) UnmarshalAndValidate(configBytes []byte) error {
	if err := c.Unmarshal(configBytes); err != nil {
		return err
	}

	if err := c.Validate(); err != nil {
		return err
	}

	return nil
}
