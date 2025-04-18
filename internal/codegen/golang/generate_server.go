package golang

import (
	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

func generateServer(sch schema.Schema, config Config) (string, error) {
	if !config.IncludeServer {
		return "", nil
	}

	g := genkit.NewGenKit().WithTabs()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Server Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	g.Line("// Server handles RPC requests.")
	g.Line("type Server[T any] struct {")
	g.Block(func() {
		g.Line("handlers          map[ProcedureName]func(context T, input json.RawMessage) (any, error)")
		g.Line("beforeMiddlewares []MiddlewareBefore[T]")
		g.Line("afterMiddlewares  []MiddlewareAfter[T]")
	})
	g.Line("}")
	g.Break()

	g.Line("// ServerRequest represents an incoming RPC request")
	g.Line("type ServerRequest[T any] struct {")
	g.Block(func() {
		g.Line("InitialContext	T")
		g.Line("RequestBody			io.Reader")
	})
	g.Line("}")
	g.Break()

	g.Line("// MiddlewareBefore runs before request processing.")
	g.Line("type MiddlewareBefore[T any] func(context T) (T, error)")
	g.Break()

	g.Line("// MiddlewareAfter runs after request processing.")
	g.Line("type MiddlewareAfter[T any] func(context T, response Response[any]) Response[any]")
	g.Break()

	g.Line("// NewServer creates a new UFO RPC server")
	g.Line("//")
	g.Line("// The generic type T represents the context type, used to pass additional data")
	g.Line("// to procedures, such as authentication information, user session or any")
	g.Line("// other data you want to pass to procedures before they are executed.")
	g.Line("func NewServer[T any]() *Server[T] {")
	g.Block(func() {
		g.Line("return &Server[T]{")
		g.Block(func() {
			g.Line("handlers:         	map[ProcedureName]func(T, json.RawMessage) (any, error){},")
			g.Line("beforeMiddlewares: 	[]MiddlewareBefore[T]{},")
			g.Line("afterMiddlewares:  	[]MiddlewareAfter[T]{},")
		})
		g.Line("}")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddMiddlewareBefore adds a middleware function that runs before the handler.")
	g.Line("//")
	g.Line("// It modifies the request context before it reaches the main procedure.")
	g.Line("//")
	g.Line("// Multiple MiddlewareBefore can be added and are processed in order.")
	g.Line("func (s *Server[T]) AddMiddlewareBefore(fn MiddlewareBefore[T]) *Server[T] {")
	g.Block(func() {
		g.Line("s.beforeMiddlewares = append(s.beforeMiddlewares, fn)")
		g.Line("return s")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddMiddlewareAfter adds a middleware function that runs after the handler.")
	g.Line("//")
	g.Line("// It modifies the response before it is sent back to the client.")
	g.Line("//")
	g.Line("// Multiple MiddlewareAfter can be added and are processed in order.")
	g.Line("func (s *Server[T]) AddMiddlewareAfter(fn MiddlewareAfter[T]) *Server[T] {")
	g.Block(func() {
		g.Line("s.afterMiddlewares = append(s.afterMiddlewares, fn)")
		g.Line("return s")
	})
	g.Line("}")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		name := procNode.Name
		namePascal := strutil.ToPascalCase(name)

		g.Linef("// Set%sHandler registers the handler for the %s procedure", namePascal, name)
		g.Linef("func (s *Server[T]) Set%sHandler(", namePascal)
		g.Block(func() {
			g.Linef("handler func(context T, input P%sInput) (P%sOutput, error),", name, name)
		})
		g.Linef(") *Server[T] {")
		g.Block(func() {
			g.Linef("s.handlers[ProcedureNames.%s] = func(context T, input json.RawMessage) (any, error) {", namePascal)
			g.Block(func() {
				g.Line("var typedInput P" + namePascal + "Input")
				g.Line("if err := json.Unmarshal(input, &typedInput); err != nil {")
				g.Block(func() {
					g.Line("return nil, &Error{")
					g.Block(func() {
						g.Line("Message: \"Invalid input for " + name + " procedure\",")
						g.Line("Details: map[string]any{\"procedure\": \"" + name + "\", \"error\": err.Error()},")
					})
					g.Line("}")
				})
				g.Line("}")
				g.Break()
				g.Line("//TODO: Use rule validation utilities here")
				g.Break()
				g.Line("return handler(context, typedInput)")
			})
			g.Line("}")
			g.Line("return s")
		})
		g.Line("}")
		g.Break()
	}

	g.Line("// HandleRequest processes an incoming RPC request")
	g.Line("func (s *Server[T]) HandleRequest(request ServerRequest[T]) Response[any] {")
	g.Block(func() {
		g.Line("var jsonBody struct {")
		g.Block(func() {
			g.Line("Procedure ProcedureName		`json:\"procedure\"`")
			g.Line("Input 		json.RawMessage	`json:\"input\"`")
		})
		g.Line("}")
		g.Line("if err := json.NewDecoder(request.RequestBody).Decode(&jsonBody); err != nil {")
		g.Block(func() {
			g.Line("return Response[any]{")
			g.Block(func() {
				g.Line("Ok: false,")
				g.Line("Error: Error{Message: \"Invalid request body\"},")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()

		g.Line("currentContext := request.InitialContext")
		g.Line("response := Response[any]{Ok: true}")
		g.Break()

		g.Line("// Validate procedure name")
		g.Line("if !slices.Contains(ProcedureNamesList, jsonBody.Procedure) {")
		g.Block(func() {
			g.Line("response = Response[any]{")
			g.Block(func() {
				g.Line("Ok: false,")
				g.Line("Error: Error{")
				g.Block(func() {
					g.Line("Message: string(jsonBody.Procedure) + \" procedure not found\",")
					g.Line("Details: map[string]any{\"procedure\": jsonBody.Procedure},")
				})
				g.Line("},")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()

		g.Line("// Validate procedure implementation")
		g.Line("if _, ok := s.handlers[jsonBody.Procedure]; response.Ok && !ok {")
		g.Block(func() {
			g.Line("response = Response[any]{")
			g.Block(func() {
				g.Line("Ok: false,")
				g.Line("Error: Error{")
				g.Block(func() {
					g.Line("Message: string(jsonBody.Procedure) + \" procedure not implemented\",")
					g.Line("Details: map[string]any{\"procedure\": jsonBody.Procedure},")
				})
				g.Line("},")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()

		g.Line("// Execute Before middleware if we haven't encountered an error yet")
		g.Line("if response.Ok {")
		g.Block(func() {
			g.Line("for _, fn := range s.beforeMiddlewares {")
			g.Block(func() {
				g.Line("var err error")
				g.Line("if currentContext, err = fn(currentContext); err != nil {")
				g.Block(func() {
					g.Line("response = Response[any]{")
					g.Block(func() {
						g.Line("Ok: false,")
						g.Line("Error: asError(err),")
					})
					g.Line("}")
					g.Line("break")
				})
				g.Line("}")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()

		g.Line("// Run handler if no errors have occurred")
		g.Line("if response.Ok {")
		g.Block(func() {
			g.Line("if output, err := s.handlers[jsonBody.Procedure](currentContext, jsonBody.Input); err != nil {")
			g.Block(func() {
				g.Line("response = Response[any]{")
				g.Block(func() {
					g.Line("Ok: false,")
					g.Line("Error: asError(err),")
				})
				g.Line("}")
			})
			g.Line("} else {")
			g.Block(func() {
				g.Line("response = Response[any]{")
				g.Block(func() {
					g.Line("Ok: true,")
					g.Line("Output: output,")
				})
				g.Line("}")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()

		g.Line("// Always execute After middleware, regardless of any previous errors")
		g.Line("for _, fn := range s.afterMiddlewares {")
		g.Block(func() {
			g.Line("response = fn(currentContext, response)")
		})
		g.Line("}")
		g.Break()

		g.Line("return response")
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
