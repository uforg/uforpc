package golang

import (
	"fmt"

	"github.com/uforg/uforpc/internal/codegen/genkit"
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
		g.Line("handlers          map[ProcedureName]func(context T, input any) (any, error)")
		g.Line("beforeMiddlewares []MiddlewareBefore[T]")
		g.Line("afterMiddlewares  []MiddlewareAfter[T]")
		g.Line("methodMap         map[ProcedureName]HTTPMethod")
	})
	g.Line("}")
	g.Break()

	g.Line("// ServerRequest represents an incoming RPC request")
	g.Line("type ServerRequest[T any] struct {")
	g.Block(func() {
		g.Line("Method     HTTPMethod")
		g.Line("Context    T")
		g.Line("Procedure  string")
		g.Line("Input      any")
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
			g.Line("handlers:         	map[ProcedureName]func(T, any) (any, error){},")
			g.Line("beforeMiddlewares: 	[]MiddlewareBefore[T]{},")
			g.Line("afterMiddlewares:  	[]MiddlewareAfter[T]{},")
			g.Line("methodMap: map[ProcedureName]HTTPMethod{")
			g.Block(func() {
				for name, procedure := range sch.Procedures {
					namePascal := strutil.ToPascalCase(name)
					method := "GET"
					if procedure.Type == schema.ProcedureTypeMutation {
						method = "POST"
					}

					g.Line(fmt.Sprintf("ProcedureNames.%s: \"%s\",", namePascal, method))
				}
			})
			g.Line("},")
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

	for name := range sch.Procedures {
		namePascal := strutil.ToPascalCase(name)

		g.Linef("// Set%sHandler registers the handler for the %s procedure", namePascal, name)
		g.Linef("func (s *Server[T]) Set%sHandler(", namePascal)
		g.Block(func() {
			g.Linef("handler func(context T, input P%sInput) (P%sOutput, error),", name, name)
		})
		g.Linef(") *Server[T] {")
		g.Block(func() {
			g.Linef("s.handlers[ProcedureNames.%s] = func(context T, input any) (any, error) {", namePascal)
			g.Block(func() {
				g.Linef("typedInput, ok := input.(P%sInput)", name)
				g.Linef("if !ok {")
				g.Block(func() {
					g.Linef("return nil, &Error{Message: \"Invalid input type for %s\"}", name)
				})
				g.Line("}")
				g.Line("return handler(context, typedInput)")
			})
			g.Line("}")
			g.Line("return s")
		})
		g.Line("}")
		g.Break()
	}

	g.Line("// HandleRequest processes an incoming RPC request")
	g.Line("func (s *Server[T]) HandleRequest(request ServerRequest[T]) (Response[any], error) {")
	g.Block(func() {
		g.Line("procedureName := ProcedureName(request.Procedure)")
		g.Line("currentContext := request.Context")
		g.Line("response := Response[any]{Ok: true}")
		g.Line("shouldSkipHandler := false")
		g.Break()

		g.Line("// Initial validation for procedure and method")
		g.Line("if _, exists := s.handlers[procedureName]; !exists {")
		g.Block(func() {
			g.Line("response = Response[any]{")
			g.Block(func() {
				g.Line("Ok: false,")
				g.Line("Error: Error{")
				g.Block(func() {
					g.Line("Message: fmt.Sprintf(\"Handler not defined for procedure %s\", request.Procedure),")
				})
				g.Line("},")
			})
			g.Line("}")
			g.Line("shouldSkipHandler = true")
		})
		g.Line("} else if expectedMethod := s.methodMap[procedureName]; expectedMethod != request.Method {")
		g.Block(func() {
			g.Line("response = Response[any]{")
			g.Block(func() {
				g.Line("Ok: false,")
				g.Line("Error: Error{")
				g.Block(func() {
					g.Line("Message: fmt.Sprintf(\"Method %s not allowed for %s procedure\", request.Method, request.Procedure),")
				})
				g.Line("},")
			})
			g.Line("}")
			g.Line("shouldSkipHandler = true")
		})
		g.Line("}")
		g.Break()

		g.Line("// Execute Before middleware if we haven't encountered an error yet")
		g.Line("if !shouldSkipHandler {")
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
					g.Line("shouldSkipHandler = true")
					g.Line("break")
				})
				g.Line("}")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()

		g.Line("// Run handler if no errors have occurred")
		g.Line("if !shouldSkipHandler {")
		g.Block(func() {
			g.Line("//TODO: Implement validation logic here")
			g.Break()
			g.Line("if output, err := s.handlers[procedureName](currentContext, request.Input); err != nil {")
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

		g.Line("return response, nil")
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
