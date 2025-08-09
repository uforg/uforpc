package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

//go:embed pieces/server.go
var serverRawPiece string

func generateServer(sch schema.Schema, config Config) (string, error) {
	if !config.IncludeServer {
		return "", nil
	}

	piece := strutil.GetStrAfter(serverRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("server.go: could not find start delimiter")
	}

	g := genkit.NewGenKit().WithTabs()

	// Core server piece (types + internal implementation)
	g.Raw(piece)
	g.Break()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Server generated implementation")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	// Server facade
	g.Line("// Server provides a high-level, type-safe interface for UFO RPC server.")
	g.Line("// It exposes groups for procedures and streams and a global middleware API.")
	g.Line("type Server[P any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[P]")
		g.Line("Procs     *serverProcRegistry[P]")
		g.Line("Streams   *serverStreamRegistry[P]")
	})
	g.Line("}")
	g.Break()

	g.Line("// NewServer creates a new UFO RPC server instance configured to handle all")
	g.Line("// defined procedures and streams.")
	g.Line("func NewServer[P any]() *Server[P] {")
	g.Block(func() {
		g.Line("intServer := newInternalServer[P](ufoProcedureNames, ufoStreamNames)")
		g.Line("return &Server[P]{")
		g.Block(func() {
			g.Line("intServer: intServer,")
			g.Line("Procs:     newServerProcRegistry(intServer),")
			g.Line("Streams:   newServerStreamRegistry(intServer),")
		})
		g.Line("}")
	})
	g.Line("}")
	g.Break()

	g.Line("// Use registers a global middleware that executes for every request (proc and stream).")
	g.Line("func (s *Server[P]) Use(mw GlobalMiddleware[P]) { s.intServer.addGlobalMiddleware(mw) }")
	g.Break()

	g.Line("// HandleRequest is the main entry point to process incoming requests.")
	g.Line("func (s *Server[P]) HandleRequest(ctx context.Context, props P, operationName string, httpAdapter ServerHTTPAdapter) error {")
	g.Block(func() {
		g.Line("return s.intServer.handleRequest(ctx, props, operationName, httpAdapter)")
	})
	g.Line("}")
	g.Break()

	// -----------------------------------------------------------------------------
	// Procedures registry and entries
	// -----------------------------------------------------------------------------
	g.Line("type serverProcRegistry[P any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[P]")
		for _, procNode := range sch.GetProcNodes() {
			name := strutil.ToPascalCase(procNode.Name)
			g.Linef("%s proc%sEntry[P]", name, name)
		}
	})
	g.Line("}")
	g.Break()

	g.Line("func newServerProcRegistry[P any](intServer *internalServer[P]) *serverProcRegistry[P] {")
	g.Block(func() {
		g.Line("r := &serverProcRegistry[P]{intServer: intServer}")
		for _, procNode := range sch.GetProcNodes() {
			name := strutil.ToPascalCase(procNode.Name)
			g.Linef("r.%s = proc%sEntry[P]{intServer: intServer}", name, name)
		}
		g.Line("return r")
	})
	g.Line("}")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		name := strutil.ToPascalCase(procNode.Name)

		g.Linef("type proc%sEntry[P any] struct {", name)
		g.Block(func() {
			g.Line("intServer *internalServer[P]")
		})
		g.Line("}")
		g.Break()

		// Generate type aliases
		g.Linef("// Type aliases for %s procedure", name)
		g.Linef("type %sHandlerContext[P any] = HandlerContext[P, %sInput]", name, name)
		g.Linef("type %sHandlerFunc[P any] func(c *%sHandlerContext[P]) (%sOutput, error)", name, name, name)
		g.Linef("type %sMiddlewareFunc[P any] func(next %sHandlerFunc[P]) %sHandlerFunc[P]", name, name, name)
		g.Break()

		// Use (procedure middleware)
		g.Linef("// Use registers a typed middleware for the %s procedure.", name)
		g.Linef("func (e proc%sEntry[P]) Use(mw %sMiddlewareFunc[P]) {", name, name)
		g.Block(func() {
			g.Linef("adapted := func(next ProcHandlerFunc[P, any, any]) ProcHandlerFunc[P, any, any] {")
			g.Block(func() {
				g.Linef("return func(cGeneric *HandlerContext[P, any]) (any, error) {")
				g.Block(func() {
					g.Line("// 1. The \"final link\" in the specific middleware chain. When called,")
					g.Line("//    it invokes the original generic 'next' handler.")
					g.Linef("finalLink := func(c *%sHandlerContext[P]) (%sOutput, error) {", name, name)
					g.Block(func() {
						g.Line("// Call the next generic handler in the chain.")
						g.Line("genericOutput, err := next(cGeneric)")
						g.Line("if err != nil {")
						g.Block(func() {
							g.Linef("// On error, return the zero value for the specific output type.")
							g.Linef("var zero %sOutput", name)
							g.Line("return zero, err")
						})
						g.Line("}")

						g.Line("// On success, assert the 'any' output to the specific output type.")
						g.Line("// It's assumed a higher layer guarantees the type is correct.")
						g.Linef("specificOutput, _ := genericOutput.(%sOutput)", name)
						g.Line("return specificOutput, nil")
					})
					g.Line("}")

					g.Line("// 2. Apply the user-defined typed middleware to the final link.")
					g.Line("handlerChain := mw(finalLink)")

					g.Line("// 3. Construct the typed context from the generic one.")
					g.Line("//    It's assumed a higher layer guarantees the type is correct.")
					g.Linef("input, _ := cGeneric.Input.(%sInput)", name)
					g.Linef("cSpecific := &%sHandlerContext[P]{", name)
					g.Block(func() {
						g.Line("Input:         input,")
						g.Line("Props:         cGeneric.Props,")
						g.Line("Context:       cGeneric.Context,")
						g.Line("operationName: cGeneric.operationName,")
						g.Line("operationType: cGeneric.operationType,")
					})
					g.Line("}")

					g.Line("// 4. Execute the complete middleware chain with the typed context.")
					g.Line("return handlerChain(cSpecific)")
				})
				g.Line("}")
			})
			g.Line("}")
			g.Linef("e.intServer.addProcMiddleware(\"%s\", adapted)", name)
		})
		g.Line("}")
		g.Break()

		// Handle (procedure handler)
		g.Linef("// Handle registers the business handler for the %s procedure.", name)
		g.Linef("func (e proc%sEntry[P]) Handle(handler %sHandlerFunc[P]) {", name, name)
		g.Block(func() {
			g.Linef("adaptedHandler := func(cGeneric *HandlerContext[P, any]) (any, error) {")
			g.Block(func() {
				g.Line("// 1. Create the specific context from the generic one provided by the server.")
				g.Line("//    This requires a type assertion on the input field.")
				g.Line("//    It's assumed a higher layer guarantees the type is correct.")
				g.Linef("input, _ := cGeneric.Input.(%sInput)", name)
				g.Linef("cSpecific := &%sHandlerContext[P]{", name)
				g.Block(func() {
					g.Line("Input:         input,")
					g.Line("Props:         cGeneric.Props,")
					g.Line("Context:       cGeneric.Context,")
					g.Line("operationName: cGeneric.operationName,")
					g.Line("operationType: cGeneric.operationType,")
				})
				g.Line("}")

				g.Line("// 2. Call the user-provided, type-safe handler with the adapted context.")
				g.Line("//    The return values are compatible with (any, error).")
				g.Line("return handler(cSpecific)")
			})
			g.Line("}")

			g.Linef("deserializer := func(raw json.RawMessage) (any, error) {")
			g.Block(func() {
				g.Linef("var pre pre%sInput", name)
				g.Line("if err := json.Unmarshal(raw, &pre); err != nil {")
				g.Block(func() { g.Linef("return nil, fmt.Errorf(\"failed to unmarshal %s input: %%w\", err)", name) })
				g.Line("}")
				g.Line("if err := pre.validate(); err != nil { return nil, err }")
				g.Line("typed := pre.transform()")
				g.Line("return typed, nil")
			})
			g.Line("}")

			g.Linef("e.intServer.setProcHandler(\"%s\", adaptedHandler, deserializer)", name)
		})
		g.Line("}")
		g.Break()
	}

	// -----------------------------------------------------------------------------
	// Streams registry and entries
	// -----------------------------------------------------------------------------
	g.Line("type serverStreamRegistry[P any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[P]")
		for _, streamNode := range sch.GetStreamNodes() {
			name := strutil.ToPascalCase(streamNode.Name)
			g.Linef("%s stream%sEntry[P]", name, name)
		}
	})
	g.Line("}")
	g.Break()

	g.Line("func newServerStreamRegistry[P any](intServer *internalServer[P]) *serverStreamRegistry[P] {")
	g.Block(func() {
		g.Line("r := &serverStreamRegistry[P]{intServer: intServer}")
		for _, streamNode := range sch.GetStreamNodes() {
			name := strutil.ToPascalCase(streamNode.Name)
			g.Linef("r.%s = stream%sEntry[P]{intServer: intServer}", name, name)
		}
		g.Line("return r")
	})
	g.Line("}")
	g.Break()

	for _, streamNode := range sch.GetStreamNodes() {
		name := strutil.ToPascalCase(streamNode.Name)
		g.Linef("type stream%sEntry[P any] struct {", name)
		g.Block(func() {
			g.Line("intServer *internalServer[P]")
		})
		g.Line("}")
		g.Break()

		// Generate type aliases
		g.Linef("// Type aliases for %s stream", name)
		g.Linef("type %sHandlerContext[P any] = HandlerContext[P, %sInput]", name, name)
		g.Linef("type %sEmitFunc[P any] = EmitFunc[P, %sInput, %sOutput]", name, name, name)
		g.Linef("type %sHandlerFunc[P any] func(c *%sHandlerContext[P], emit %sEmitFunc[P]) error", name, name, name)
		g.Linef("type %sMiddlewareFunc[P any] func(next %sHandlerFunc[P]) %sHandlerFunc[P]", name, name, name)
		g.Linef("type %sEmitMiddlewareFunc[P any] func(next %sEmitFunc[P]) %sEmitFunc[P]", name, name, name)
		g.Break()

		// Generate Use (stream middleware)
		g.Linef("// Use registers a typed middleware for the %s stream.", name)
		g.Linef("//")
		g.Linef("// This function allows you to add a middleware specific to the %s stream.", name)
		g.Linef("// The middleware is applied to the stream's handler chain, enabling you to intercept,")
		g.Linef("// modify, or augment the handling of incoming stream requests for %s.", name)
		g.Linef("func (e stream%sEntry[P]) Use(mw %sMiddlewareFunc[P]) {", name, name)
		g.Block(func() {
			g.Linef("adapted := func(next StreamHandlerFunc[P, any, any]) StreamHandlerFunc[P, any, any] {")
			g.Block(func() {
				// Returns the final handler that the system will execute.
				g.Line("// Returns the final handler that the system will execute.")
				g.Linef("return func(cGeneric *HandlerContext[P, any], emitGeneric EmitFunc[P, any, any]) error {")
				g.Block(func() {
					g.Line("// 1. The \"final link\" in the specific middleware chain. When called,")
					g.Line("//    it invokes the original generic 'next' handler, using the generic")
					g.Line("//    arguments captured from this closure.")
					g.Linef("finalLink := func(c *%sHandlerContext[P], emit %sEmitFunc[P]) error {", name, name)
					g.Block(func() {
						g.Line("return next(cGeneric, emitGeneric)")
					})
					g.Line("}")

					g.Line("// 2. Apply the user-defined typed middleware to the final link.")
					g.Line("handlerChain := mw(finalLink)")

					g.Line("// 3. Create a typed 'emit' function that delegates to the generic emit.")
					g.Line("//    Uses 'cGeneric' from the outer scope, which is the correct context.")
					g.Linef("emitSpecific := func(c *%sHandlerContext[P], output %sOutput) error {", name, name)
					g.Block(func() {
						g.Line("return emitGeneric(cGeneric, output)")
					})
					g.Line("}")

					g.Line("// 4. Construct the typed context from the generic context.")
					g.Linef("input, _ := cGeneric.Input.(%sInput)", name)
					g.Linef("cSpecific := &%sHandlerContext[P]{", name)
					g.Block(func() {
						g.Line("Input:         input,")
						g.Line("Props:         cGeneric.Props,")
						g.Line("Context:       cGeneric.Context,")
						g.Line("operationName: cGeneric.operationName,")
						g.Line("operationType: cGeneric.operationType,")
					})
					g.Line("}")

					g.Line("// 5. Execute the complete middleware chain with the typed arguments.")
					g.Line("return handlerChain(cSpecific, emitSpecific)")
				})
				g.Line("}")
			})
			g.Line("}")
			g.Linef("e.intServer.addStreamMiddleware(\"%s\", adapted)", name)
		})
		g.Line("}")
		g.Break()

		// UseEmit (emit middleware)
		g.Linef("// UseEmit registers a typed emit middleware for the %s stream.", name)
		g.Linef("func (e stream%sEntry[P]) UseEmit(mw %sEmitMiddlewareFunc[P]) {", name, name)
		g.Block(func() {
			g.Linef("adapted := func(next EmitFunc[P, any, any]) EmitFunc[P, any, any] {")
			g.Block(func() {
				g.Line("// Return a new generic 'emit' function that wraps the logic.")
				g.Line("return func(cGeneric *HandlerContext[P, any], outputGeneric any) error {")
				g.Block(func() {
					g.Line("// 1. The \"next link\" in the emit chain. It is a wrapper that")
					g.Line("//    calls the original generic 'next' provided by the system.")
					g.Linef("nextSpecific := func(c *%sHandlerContext[P], output %sOutput) error {", name, name)
					g.Block(func() {
						g.Line("return next(cGeneric, output)")
					})
					g.Line("}")

					g.Line("// 2. Apply the specific middleware to obtain the final 'emit' function.")
					g.Line("//    The middleware now only takes the next function in the chain.")
					g.Line("emitChain := mw(nextSpecific)")

					g.Line("// 3. Create the specific context from the generic one. This is still")
					g.Line("//    needed to call the final emitChain function.")
					g.Linef("input, _ := cGeneric.Input.(%sInput)", name)
					g.Linef("cSpecific := &%sHandlerContext[P]{", name)
					g.Block(func() {
						g.Line("Input:         input,")
						g.Line("Props:         cGeneric.Props,")
						g.Line("Context:       cGeneric.Context,")
						g.Line("operationName: cGeneric.operationName,")
						g.Line("operationType: cGeneric.operationType,")
					})
					g.Line("}")

					g.Line("// 4. Perform the type assertion for the output object.")
					g.Linef("outputSpecific, _ := outputGeneric.(%sOutput)", name)

					g.Line("// 5. Execute the middleware chain with the specific arguments.")
					g.Line("return emitChain(cSpecific, outputSpecific)")
				})
				g.Line("}")
			})
			g.Line("}")
			g.Linef("e.intServer.addStreamEmitMiddleware(\"%s\", adapted)", name)
		})
		g.Line("}")
		g.Break()

		// Handle (stream handler)
		g.Linef("// Handle registers the business handler for the %s stream.", name)
		g.Linef("func (e stream%sEntry[P]) Handle(handler %sHandlerFunc[P]) {", name, name)
		g.Block(func() {
			g.Linef("adaptedHandler := func(cGeneric *HandlerContext[P, any], emitGeneric EmitFunc[P, any, any]) error {")
			g.Block(func() {
				g.Line("// 1. Create the specific, type-safe emit function by wrapping the generic one.")
				g.Line("//    It uses 'cGeneric' from the outer scope, which has the correct type for the generic call.")
				g.Linef("emitSpecific := func(c *%sHandlerContext[P], output %sOutput) error {", name, name)
				g.Block(func() {
					g.Line("return emitGeneric(cGeneric, output)")
				})
				g.Line("}")

				g.Line("// 2. Create the specific context from the generic one provided by the server.")
				g.Line("//    This requires a type assertion on the input field.")
				g.Linef("input, _ := cGeneric.Input.(%sInput)", name)
				g.Linef("cSpecific := &%sHandlerContext[P]{", name)
				g.Block(func() {
					g.Line("Input:         input,")
					g.Line("Props:         cGeneric.Props,")
					g.Line("Context:       cGeneric.Context,")
					g.Line("operationName: cGeneric.operationName,")
					g.Line("operationType: cGeneric.operationType,")
				})
				g.Line("}")

				g.Line("// 3. Call the user-provided, type-safe handler with the adapted, specific arguments.")
				g.Line("return handler(cSpecific, emitSpecific)")
			})
			g.Line("}")

			g.Linef("deserializer := func(raw json.RawMessage) (any, error) {")
			g.Block(func() {
				g.Linef("var pre pre%sInput", name)
				g.Line("if err := json.Unmarshal(raw, &pre); err != nil {")
				g.Block(func() { g.Linef("return nil, fmt.Errorf(\"failed to unmarshal %s input: %%w\", err)", name) })
				g.Line("}")
				g.Line("if err := pre.validate(); err != nil { return nil, err }")
				g.Line("typed := pre.transform()")
				g.Line("return typed, nil")
			})
			g.Line("}")

			g.Linef("e.intServer.setStreamHandler(\"%s\", adaptedHandler, deserializer)", name)
		})
		g.Line("}")
		g.Break()
	}

	return g.String(), nil
}
