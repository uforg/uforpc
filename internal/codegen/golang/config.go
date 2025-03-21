package golang

// Config is the configuration for the Go code generator.
type Config struct {
	// PackageName is the name of the package to generate the code in.
	PackageName string `json:"packageName"`
	// IncludeServer enables server code generation.
	IncludeServer bool `json:"includeServer"`
	// IncludeClient enables client code generation.
	IncludeClient bool `json:"includeClient"`
	// OmitServerRequestValidation disables server request validation in the generated server code.
	OmitServerRequestValidation bool `json:"omitServerRequestValidation"`
	// OmitClientRequestValidation disables client request validation in the generated client code.
	OmitClientRequestValidation bool `json:"omitClientRequestValidation"`
}
