package golang

import (
	_ "embed"
	"fmt"
	"strings"

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

	g.Raw(piece)
	g.Break()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Server Generated Implementation")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	g.Line("// hookRegistry provides a type-safe interface for registering server hooks")
	g.Line("// that execute at different stages of request processing. Hooks allow for")
	g.Line("// cross-cutting concerns like authentication, logging, and response transformation.")
	g.Line("type hookRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBeforeHandler registers a hook that executes before any procedure or stream")
	g.Line("// handler is invoked. This hook is useful for cross-cutting concerns like")
	g.Line("// authentication, authorization, logging, or request validation that apply")
	g.Line("// to all handlers.")
	g.Line("//")
	g.Line("// Hooks are executed in the order they were registered.")
	g.Line("//")
	g.Line("// The hook can modify the context and UFO context, or return an error to")
	g.Line("// abort the request with an error response.")
	g.Line("func (m *hookRegistry[T]) AddBeforeHandler(hook ServerHookBeforeHandler[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookBeforeHandler(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBeforeProcRespond registers a hook that executes before the procedure")
	g.Line("// response is sent to the client. This hook can modify the response data,")
	g.Line("// allowing for response transformation, filtering, or adding metadata.")
	g.Line("//")
	g.Line("// Hooks are executed in the order they were registered.")
	g.Line("//")
	g.Line("// This hook is useful for standardizing response formats, adding data,")
	g.Line("// or implementing response logic.")
	g.Line("func (m *hookRegistry[T]) AddBeforeProcRespond(hook ServerHookBeforeProcRespond[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookBeforeProcRespond(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBeforeStreamEmit registers a hook that executes before each stream event")
	g.Line("// is emitted to the client. This hook can modify the response data for each")
	g.Line("// event, allowing for real-time response transformation or filtering.")
	g.Line("//")
	g.Line("// Hooks are executed in the order they were registered.")
	g.Line("//")
	g.Line("// This hook is useful for adding metadata to stream events, implementing")
	g.Line("// event filtering logic, or transforming event data before transmission.")
	g.Line("func (m *hookRegistry[T]) AddBeforeStreamEmit(hook ServerHookBeforeStreamEmit[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookBeforeStreamEmit(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfterProc registers a hook that executes after procedure request")
	g.Line("// processing completes, regardless of success or failure. This hook cannot")
	g.Line("// modify the response as it has already been sent to the client.")
	g.Line("//")
	g.Line("// Hooks are executed in the order they were registered.")
	g.Line("//")
	g.Line("// This hook is useful for logging, metrics collection, cleanup operations,")
	g.Line("// or audit trails specific to procedure calls.")
	g.Line("func (m *hookRegistry[T]) AddAfterProc(hook ServerHookAfterProc[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookAfterProc(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfterStream registers a hook that executes after stream request")
	g.Line("// processing ends, either successfully or with an error. This hook cannot")
	g.Line("// modify anything as the stream has already completed.")
	g.Line("//")
	g.Line("// Hooks are executed in the order they were registered.")
	g.Line("//")
	g.Line("// This hook is useful for cleanup operations, logging, metrics collection,")
	g.Line("// or resource deallocation specific to stream handling.")
	g.Line("func (m *hookRegistry[T]) AddAfterStream(hook ServerHookAfterStream[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookAfterStream(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfterStreamEmit registers a hook that executes after each stream event")
	g.Line("// has been successfully emitted to the client. This hook cannot modify")
	g.Line("// anything as the event has already been sent.")
	g.Line("//")
	g.Line("// Hooks are executed in the order they were registered.")
	g.Line("//")
	g.Line("// This hook is useful for logging individual stream events, metrics")
	g.Line("// collection, or post-emission cleanup operations.")
	g.Line("func (m *hookRegistry[T]) AddAfterStreamEmit(hook ServerHookAfterStreamEmit[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookAfterStreamEmit(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// procRegistry provides a type-safe interface for registering procedure")
	g.Line("// handlers and input processors. It handles the marshaling and validation")
	g.Line("// of procedure inputs and manages the procedure execution lifecycle.")
	g.Line("type procRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		name := strutil.ToPascalCase(procNode.Name)

		g.Linef("// Set%sInputHandler registers an input processor for the %s procedure.", name, name)
		g.Line("//")
		g.Line("// This handler is responsible for validating and transforming the input data")
		g.Line("// before it reaches the main procedure handler. It receives the typed input")
		g.Line("// and must return the validated/transformed input or an error.")
		g.Line("//")
		g.Line("// This is useful for implementing custom validation logic, input sanitization,")
		g.Line("// or data transformation that goes beyond the default required fields validation.")
		renderDeprecated(g, procNode.Deprecated)
		g.Linef("func (p *procRegistry[T]) Set%sInputHandler(", name)
		g.Block(func() {
			g.Linef("inputHandler func(ctx context.Context, ufoCtx T, input %sInput) (%sInput, error),", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef("typedInputHandler := func(ctx context.Context, ufoCtx T, input any) (any, error) {")
			g.Block(func() {
				g.Linef("typedInput, ok := input.(%sInput)", name)
				g.Line("if !ok {")
				g.Block(func() {
					g.Linef("return nil, fmt.Errorf(\"invalid input type for %s procedure, expected %sInput, got %%T\", input)", name, name)
				})
				g.Line("}")
				g.Line("return inputHandler(ctx, ufoCtx, typedInput)")
			})
			g.Line("}")
			g.Linef("p.intServer.setProcInputHandler(\"%s\", typedInputHandler)", name)
		})
		g.Line("}")
		g.Break()

		g.Linef("// Set%sHandler registers the main implementation for the %s procedure.", name, name)
		g.Line("//")
		g.Line("// The handler function will be called when a client invokes this procedure")
		g.Line("// via RPC. It receives the typed input data and must return the typed output")
		g.Line("// or an error.")
		g.Line("//")
		g.Line("// The handler is executed after all before-handler hooks and input processing.")
		g.Line("// If an input handler is registered, the input will be processed through that")
		g.Line("// handler before reaching this main handler.")
		renderDeprecated(g, procNode.Deprecated)
		g.Linef("func (p *procRegistry[T]) Set%sHandler(", name)
		g.Block(func() {
			g.Linef("handler func(ctx context.Context, ufoCtx T, input %sInput) (%sOutput, error),", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef(" p.intServer.setProcHandler(\"%s\", func(ctx context.Context, ufoCtx T, rawInput json.RawMessage) (any, error) {", name)
			g.Block(func() {
				g.Line("var preTypedInput pre" + name + "Input")
				g.Line("if err := json.Unmarshal(rawInput, &preTypedInput); err != nil {")
				g.Block(func() {
					g.Linef(`return nil, fmt.Errorf("failed to unmarshal %s input: %%w", err)`, name)
				})
				g.Line("}")
				g.Break()

				g.Line("if err := preTypedInput.validate(); err != nil {")
				g.Block(func() {
					g.Line("return nil, err")
				})
				g.Line("}")
				g.Break()

				g.Line("typedInput := preTypedInput.transform()")

				g.Linef("if inputHandler, ok := p.intServer.procInputHandlers[\"%s\"]; ok {", name)
				g.Block(func() {
					g.Line("processedInput, err := inputHandler(ctx, ufoCtx, typedInput)")
					g.Line("if err != nil {")
					g.Block(func() {
						g.Line("return nil, err")
					})
					g.Line("}")
					g.Linef("typedInput = processedInput.(%sInput)", name)
				})
				g.Line("}")
				g.Break()

				g.Line("return handler(ctx, ufoCtx, typedInput)")
			})
			g.Line("})")
		})
		g.Line("}")
		g.Break()
	}

	g.Line("// streamRegistry provides a type-safe interface for registering stream")
	g.Line("// handlers and input processors. It handles the marshaling and validation")
	g.Line("// of stream inputs and manages the stream execution lifecycle with")
	g.Line("// Server-Sent Events delivery.")
	g.Line("type streamRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	for _, streamNode := range sch.GetStreamNodes() {
		name := strutil.ToPascalCase(streamNode.Name)

		g.Linef("// Set%sInputHandler registers an input processor for the %s stream.", name, name)
		g.Line("//")
		g.Line("// This handler is responsible for validating and transforming the input data")
		g.Line("// before the stream begins. It receives the typed input and must return")
		g.Line("// the validated/transformed input or an error.")
		g.Line("//")
		g.Line("// This is useful for implementing custom validation logic, input sanitization,")
		g.Line("// or data transformation specific to stream initialization that goes beyond")
		g.Line("// the default required fields validation.")
		renderDeprecated(g, streamNode.Deprecated)
		g.Linef("func (s *streamRegistry[T]) Set%sInputHandler(", name)
		g.Block(func() {
			g.Linef("inputHandler func(ctx context.Context, ufoCtx T, input %sInput) (%sInput, error),", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef("typedInputHandler := func(ctx context.Context, ufoCtx T, input any) (any, error) {")
			g.Block(func() {
				g.Linef("typedInput, ok := input.(%sInput)", name)
				g.Line("if !ok {")
				g.Block(func() {
					g.Linef("return nil, fmt.Errorf(\"invalid input type for %s stream, expected %sInput, got %%T\", input)", name, name)
				})
				g.Line("}")
				g.Line("return inputHandler(ctx, ufoCtx, typedInput)")
			})
			g.Line("}")
			g.Linef("s.intServer.setStreamInputHandler(\"%s\", typedInputHandler)", name)
		})
		g.Line("}")
		g.Break()

		g.Linef("// Set%sHandler registers the main implementation for the %s stream.", name, name)
		g.Line("//")
		g.Line("// The handler function will be called when a client initiates this stream")
		g.Line("// via RPC. It receives the typed input data and an emit function for sending")
		g.Line("// events to the client. The handler should call emit for each event and")
		g.Line("// return when the stream is complete or an error occurs.")
		g.Line("//")
		g.Line("// The handler is executed after all before-handler hooks and input processing.")
		g.Line("//")
		g.Line("// Each emitted event goes through before-emit and after-emit hooks.")
		renderDeprecated(g, streamNode.Deprecated)
		g.Linef("func (s *streamRegistry[T]) Set%sHandler(", name)
		g.Block(func() {
			g.Linef("handler func(ctx context.Context, ufoCtx T, input %sInput, emit func(%sOutput) error) error,", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef("s.intServer.setStreamHandler(\"%s\", func(ctx context.Context, ufoCtx T, rawInput json.RawMessage, rawEmit func(any) error) error {", name)
			g.Block(func() {
				g.Line("var preTypedInput pre" + name + "Input")
				g.Line("if err := json.Unmarshal(rawInput, &preTypedInput); err != nil {")
				g.Block(func() {
					g.Linef(`return fmt.Errorf("failed to unmarshal %s input: %%w", err)`, name)
				})
				g.Line("}")
				g.Break()

				g.Line("if err := preTypedInput.validate(); err != nil {")
				g.Block(func() {
					g.Line("return err")
				})
				g.Line("}")
				g.Break()

				g.Line("typedInput := preTypedInput.transform()")
				g.Line("typedEmit := func(output " + name + "Output) error {")
				g.Block(func() {
					g.Line("return rawEmit(output)")
				})
				g.Line("}")
				g.Break()

				g.Linef("if inputHandler, ok := s.intServer.streamInputHandlers[\"%s\"]; ok {", name)
				g.Block(func() {
					g.Line("processedInput, err := inputHandler(ctx, ufoCtx, typedInput)")
					g.Line("if err != nil {")
					g.Block(func() {
						g.Line("return err")
					})
					g.Line("}")
					g.Linef("typedInput = processedInput.(%sInput)", name)
				})
				g.Line("}")
				g.Break()

				g.Line("return handler(ctx, ufoCtx, typedInput, typedEmit)")
			})
			g.Line("})")
		})
		g.Line("}")
		g.Break()
	}

	g.Line("// Server provides a high-level, type-safe interface for UFO RPC server")
	g.Line("// functionality. It wraps the internal server implementation and organizes")
	g.Line("// methods into logical groups for hooks, procedures, and streams.")
	g.Line("//")
	g.Line("// The generic type T represents the UFO context type, allowing you to pass")
	g.Line("// custom data (authentication info, user sessions, etc.) through the entire")
	g.Line("// request processing pipeline.")
	g.Line("type Server[T any] struct {")
	g.Block(func() {
		g.Line("intServer  *internalServer[T]")
		g.Line("Hooks      *hookRegistry[T]")
		g.Line("Procs      *procRegistry[T]")
		g.Line("Streams    *streamRegistry[T]")
	})
	g.Line("}")
	g.Break()

	g.Line("// NewServer creates a new UFO RPC server instance configured to handle")
	g.Line("// all defined procedures and streams. The server is initialized with empty")
	g.Line("// handler registrations and hook chains, ready for configuration.")
	g.Line("//")
	g.Line("// The generic type T represents the UFO context type, used to pass additional")
	g.Line("// data to handlers throughout the request lifecycle. This can include:")
	g.Line("//   - Authentication information")
	g.Line("//   - User session data")
	g.Line("//   - Database connections")
	g.Line("//   - Any other request-scoped data")
	g.Line("//")
	g.Line("// Example usage:")
	g.Line("//   type AppContext struct {")
	g.Line("//       DB     *sql.DB")
	g.Line("//       UserID string")
	g.Line("//   }")
	g.Line("//   server := NewServer[AppContext]()")
	g.Line("func NewServer[T any]() *Server[T] {")
	g.Block(func() {
		g.Line("intServer := newInternalServer[T](ufoProcedureNames, ufoStreamNames)")
		g.Line("return &Server[T]{")
		g.Block(func() {
			g.Line("intServer:  intServer,")
			g.Line("Hooks:      &hookRegistry[T]{intServer: intServer},")
			g.Line("Procs:      &procRegistry[T]{intServer: intServer},")
			g.Line("Streams:    &streamRegistry[T]{intServer: intServer},")
		})
		g.Line("}")
	})
	g.Line("}")
	g.Break()

	g.Line("// HandleRequest processes an incoming RPC request and manages the complete")
	g.Line("// request-response lifecycle. This is the main entry point for handling")
	g.Line("// client requests in your HTTP server.")
	g.Line("//")
	g.Line("// The method handles all aspects of request processing including:")
	g.Line("//   - Request parsing and validation")
	g.Line("//   - Context management and hook execution")
	g.Line("//   - Handler dispatch (procedures and streams)")
	g.Line("//   - Response formatting and delivery")
	g.Line("//   - Error handling and reporting")
	g.Line("//")
	g.Line("// The operationName parameter specifies the name of the procedure or stream")
	g.Line("// to invoke. It must be one of the procedure or stream names in the schema,")
	g.Line("// otherwise the request will be rejected with an error.")
	g.Line("//")
	g.Line("// The operationName parameter must be extracted from the last request url path")
	g.Line("// segment. For example, if the request url is \"/api/v1/urpc/GetUser\", the")
	g.Line("// operationName is \"GetUser\".")
	g.Line("//")
	g.Line("// The httpAdapter parameter provides the bridge between UFO RPC and your")
	g.Line("// HTTP server implementation:")
	g.Line("//   - For standard net/http: use NewServerNetHTTPAdapter(w, r)")
	g.Line("//   - For other frameworks (gin, echo, etc.): implement ServerHTTPAdapter")
	g.Line("//")
	g.Line("// Example with net/http:")
	g.Line("//   http.HandleFunc(\"POST /api/v1/urpc/{operationName}\", func(w http.ResponseWriter, r *http.Request) {")
	g.Line("//       ctx := r.Context()")
	g.Line("//       ufoCtx := AppContext{DB: db, Foo: \"bar\"}")
	g.Line("//       operationName := r.PathValue(\"operationName\")")
	g.Line("//       httpAdapter := NewServerNetHTTPAdapter(w, r)")
	g.Line("//       server.HandleRequest(ctx, ufoCtx, operationName, httpAdapter)")
	g.Line("//   })")
	g.Line("func (s *Server[T]) HandleRequest(ctx context.Context, ufoCtx T, operationName string, httpAdapter ServerHTTPAdapter) error {")
	g.Block(func() {
		g.Line("return s.intServer.handleRequest(ctx, ufoCtx, operationName, httpAdapter)")
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}

// renderDeprecated receives a pointer to a string and if it is not nil, it will
// render a comment with the deprecated message to the given genkit.GenKit.
func renderDeprecated(g *genkit.GenKit, deprecated *string) {
	if deprecated == nil {
		return
	}

	desc := "Deprecated: "
	if *deprecated == "" {
		desc += "This is deprecated and should not be used in new code."
	} else {
		desc += *deprecated
	}

	g.Line("//")
	for line := range strings.SplitSeq(desc, "\n") {
		g.Linef("// %s", line)
	}
}
