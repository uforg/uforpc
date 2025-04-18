package golang

import "fmt"

// Config is the configuration for the Go code generator.
type Config struct {
	// OutputDir is the directory to output the generated code to.
	OutputDir string `toml:"output_dir"`
	// PackageName is the name of the package to generate the code in.
	PackageName string `toml:"package_name"`
	// IncludeServer enables server code generation.
	IncludeServer bool `toml:"include_server"`
	// IncludeClient enables client code generation.
	IncludeClient bool `toml:"include_client"`
	// OmitServerRequestValidation disables server request validation in the generated server code.
	OmitServerRequestValidation bool `toml:"omit_server_request_validation"`
	// OmitClientRequestValidation disables client request validation in the generated client code.
	OmitClientRequestValidation bool `toml:"omit_client_request_validation"`
}

func (c Config) Validate() error {
	if c.OutputDir == "" {
		return fmt.Errorf(`"output_dir" is required`)
	}
	if c.PackageName == "" {
		return fmt.Errorf(`"package_name" is required`)
	}
	return nil
}
