package golang

import (
	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

func generateCoreTypes(g *genkit.GenKit, _ schema.Schema, _ Config) error {
	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Core Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Line()

	g.Line("// UFOHTTPMethod represents an HTTP method.")
	g.Line("type UFOHTTPMethod string")
	g.Line()

	g.Line("const (")
	g.Block(func() {
		g.Line("// GET represents the HTTP GET method.")
		g.Line("GET UFOHTTPMethod = \"GET\"")
		g.Line("// POST represents the HTTP POST method.")
		g.Line("POST UFOHTTPMethod = \"POST\"")
	})
	g.Line(")")
	g.Line()

	g.Line("// UFOResponse represents the response of a UFO RPC call.")
	g.Line("type UFOResponse[T any] struct {")
	g.Block(func() {
		g.Line("Ok     bool             `json:\"ok\"`")
		g.Line("Output T                `json:\"output,omitempty,omitzero\"`")
		g.Line("Error  UFOError         `json:\"error,omitempty,omitzero\"`")
	})
	g.Line("}")
	g.Line()

	g.Line("// UFOError represents a standardized error in the UFO RPC system.")
	g.Line("//")
	g.Line("// It provides structured information about errors that occur within the system,")
	g.Line("// enabling consistent error handling across servers and clients.")
	g.Line("//")
	g.Line("// Fields:")
	g.Line("//   - Message: A human-readable description of the error.")
	g.Line("//   - Category: Optional. Categorizes the error by its nature or source (e.g., \"ValidationError\", \"DatabaseError\").")
	g.Line("//   - Code: Optional. A machine-readable identifier for the specific error condition (e.g., \"INVALID_EMAIL\").")
	g.Line("//   - Details: Optional. Additional information about the error.")
	g.Line("//")
	g.Line("// The struct implements the error interface.")
	g.Line("type UFOError struct {")
	g.Block(func() {
		g.Line("// Message provides a human-readable description of the error.")
		g.Line("//")
		g.Line("// This message can be displayed to end-users or used for logging and debugging purposes.")
		g.Line("//")
		g.Line("// Use Cases:")
		g.Line("//   1. If localization is not implemented, Message can be directly shown to the user to inform them of the issue.")
		g.Line("//   2. Developers can use Message in logs to diagnose problems during development or in production.")
		g.Line("Message string `json:\"message\"`")
		g.Line()

		g.Line("// Category categorizes the error by its nature or source.")
		g.Line("//")
		g.Line("// Examples:")
		g.Line("//   - \"ValidationError\" for input validation errors.")
		g.Line("//   - \"DatabaseError\" for errors originating from database operations.")
		g.Line("//   - \"AuthenticationError\" for authentication-related issues.")
		g.Line("//")
		g.Line("// Use Cases:")
		g.Line("//   1. In middleware, you can use Category to determine how to handle the error.")
		g.Line("//      For instance, you might log \"InternalError\" types and return a generic message to the client.")
		g.Line("//   2. Clients can inspect the Category to decide whether to prompt the user for action,")
		g.Line("//      such as re-authentication if the Category is \"AuthenticationError\".")
		g.Line("Category string `json:\"category,omitempty,omitzero\"`")
		g.Line()

		g.Line("// Code is a machine-readable identifier for the specific error condition.")
		g.Line("//")
		g.Line("// Examples:")
		g.Line("//   - \"INVALID_EMAIL\" when an email address fails validation.")
		g.Line("//   - \"USER_NOT_FOUND\" when a requested user does not exist.")
		g.Line("//   - \"RATE_LIMIT_EXCEEDED\" when a client has made too many requests.")
		g.Line("//")
		g.Line("// Use Cases:")
		g.Line("//   1. Clients can map Codes to localized error messages for internationalization (i18n),")
		g.Line("//      displaying appropriate messages based on the user's language settings.")
		g.Line("//   2. Clients or middleware can implement specific logic based on the Code,")
		g.Line("//      such as retry mechanisms for \"TEMPORARY_FAILURE\" or showing captcha for \"RATE_LIMIT_EXCEEDED\".")
		g.Line("Code string `json:\"code,omitempty,omitzero\"`")
		g.Line()

		g.Line("// Details contains optional additional information about the error.")
		g.Line("//")
		g.Line("// This field can include any relevant data that provides more context about the error.")
		g.Line("// The contents should be serializable to JSON.")
		g.Line("//")
		g.Line("// Use Cases:")
		g.Line("//   1. Providing field-level validation errors, e.g., Details could be:")
		g.Line("//      {\"fields\": {\"email\": \"Email is invalid\", \"password\": \"Password is too short\"}}")
		g.Line("//   2. Including diagnostic information such as timestamps, request IDs, or stack traces")
		g.Line("//      (ensure sensitive information is not exposed to clients).")
		g.Line("Details map[string]any `json:\"details,omitempty,omitzero\"`")
	})
	g.Line("}")
	g.Line()

	// Add Error method
	g.Line("// Error implements the error interface, returning the error message.")
	g.Line("func (e UFOError) Error() string {")
	g.Block(func() {
		g.Line("return e.Message")
	})
	g.Line("}")
	g.Line()

	// Add String method
	g.Line("// String implements the fmt.Stringer interface, returning the error message.")
	g.Line("func (e UFOError) String() string {")
	g.Block(func() {
		g.Line("return e.Message")
	})
	g.Line("}")
	g.Line()

	// Add ToJSON method
	g.Line("// ToJSON returns the UFOError as a JSON-formatted string including all its fields.")
	g.Line("// This is useful for logging and debugging purposes.")
	g.Line("//")
	g.Line("// Example usage:")
	g.Line("//   err := UFOError{")
	g.Line("//     Category: \"ValidationError\",")
	g.Line("//     Code:     \"INVALID_EMAIL\",")
	g.Line("//     Message:  \"The email address provided is invalid.\",")
	g.Line("//     Details:  map[string]any{")
	g.Line("//       \"field\": \"email\",")
	g.Line("//     },")
	g.Line("//   }")
	g.Line("//   log.Println(err.ToJSON())")
	g.Line("func (e UFOError) ToJSON() string {")
	g.Block(func() {
		g.Line("b, err := json.Marshal(e)")
		g.Line("if err != nil {")
		g.Block(func() {
			g.Line("return fmt.Sprintf(")
			g.Line("  `{\"message\":%q,\"error\":\"Failed to marshal UFOError: %s\"}`,")
			g.Line("  e.Message, err.Error(),")
			g.Line(")")
		})
		g.Line("}")
		g.Line("return string(b)")
	})
	g.Line("}")
	g.Line()

	// Add asUFOError function
	g.Line("// asUFOError converts any error into a UFOError.")
	g.Line("// If the provided error is already a UFOError, it returns it as is.")
	g.Line("// Otherwise, it wraps the error message into a new UFOError.")
	g.Line("//")
	g.Line("// This function ensures that all errors conform to the UFOError structure,")
	g.Line("// facilitating consistent error handling across the system.")
	g.Line("func asUFOError(err error) UFOError {")
	g.Block(func() {
		g.Line("switch e := err.(type) {")
		g.Line("case UFOError:")
		g.Block(func() {
			g.Line("return e")
		})
		g.Line("case *UFOError:")
		g.Block(func() {
			g.Line("return *e")
		})
		g.Line("default:")
		g.Block(func() {
			g.Line("return UFOError{")
			g.Block(func() {
				g.Line("Message: err.Error(),")
			})
			g.Line("}")
		})
		g.Line("}")
	})
	g.Line("}")

	return nil
}
