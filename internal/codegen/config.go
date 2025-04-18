package codegen

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/uforg/uforpc/internal/codegen/golang"
	"github.com/uforg/uforpc/internal/codegen/typescript"
)

type Config struct {
	Version    int                `toml:"version"`
	Schema     string             `toml:"schema"`
	Golang     *golang.Config     `toml:"golang"`
	Typescript *typescript.Config `toml:"typescript"`
}

func (c Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version: %d", c.Version)
	}

	if c.Schema == "" {
		return fmt.Errorf("schema path is required")
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

// UnmarshalTOMLConfig unmarshals and validates a TOML config
// string into a Config struct.
func UnmarshalTOMLConfig(configStr string) (Config, error) {
	var config Config
	if err := toml.Unmarshal([]byte(configStr), &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal TOML config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}
