package typescript

// Config is the configuration for the TypeScript code generator.
type Config struct {
	// IncludeServer enables server code generation.
	IncludeServer bool `json:"includeServer"`
	// IncludeClient enables client code generation.
	IncludeClient bool `json:"includeClient"`
	// OmitServerRequestValidation disables server request validation in the generated server code.
	OmitServerRequestValidation bool `json:"omitServerRequestValidation"`
	// OmitClientRequestValidation disables client request validation in the generated client code.
	OmitClientRequestValidation bool `json:"omitClientRequestValidation"`
	// OmitClientDefaultFetch disables the default fetch implementation in the generated client code.
	OmitClientDefaultFetch bool `json:"omitClientDefaultFetch"`
}
