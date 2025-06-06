package typescript

import "fmt"

// Config is the configuration for the TypeScript code generator.
type Config struct {
	// OutputDir is the directory to output the generated code to.
	OutputDir string `toml:"output_dir"`
	// IncludeServer enables server code generation.
	IncludeServer bool `toml:"include_server"`
	// IncludeClient enables client code generation.
	IncludeClient bool `toml:"include_client"`
	// OmitClientDefaultFetch disables the default fetch implementation in the generated client code.
	OmitClientDefaultFetch bool `toml:"omit_client_default_fetch"`
}

func (c Config) Validate() error {
	if c.OutputDir == "" {
		return fmt.Errorf("output_dir is required")
	}
	return nil
}
