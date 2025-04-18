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
	// OmitServerRequestValidation disables server request validation in the generated server code.
	OmitServerRequestValidation bool `toml:"omit_server_request_validation"`
	// OmitClientRequestValidation disables client request validation in the generated client code.
	OmitClientRequestValidation bool `toml:"omit_client_request_validation"`
	// OmitClientDefaultFetch disables the default fetch implementation in the generated client code.
	OmitClientDefaultFetch bool `toml:"omit_client_default_fetch"`
}

func (c Config) Validate() error {
	if c.OutputDir == "" {
		return fmt.Errorf("output_dir is required")
	}
	return nil
}
