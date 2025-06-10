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

	g.Raw(piece)
	g.Break()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Server Generated Implementation")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	g.Line("// hookRegistry handles hook registration")
	g.Line("type hookRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBeforeHandler adds a hook function that runs before the handler. Both for procedures and streams.")
	g.Line("func (m *hookRegistry[T]) AddBeforeHandler(hook ServerHookBeforeHandler[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookBeforeHandler(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBeforeProcRespond adds a hook function that runs before the procedure handler responds.")
	g.Line("func (m *hookRegistry[T]) AddBeforeProcRespond(hook ServerHookBeforeProcRespond[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookBeforeProcRespond(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBeforeStreamEmit adds a hook function that runs before the stream event is emitted.")
	g.Line("func (m *hookRegistry[T]) AddBeforeStreamEmit(hook ServerHookBeforeStreamEmit[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookBeforeStreamEmit(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfterProc adds a hook function that runs after the procedure handler.")
	g.Line("func (m *hookRegistry[T]) AddAfterProc(hook ServerHookAfterProc[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookAfterProc(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfterStream adds a hook function that runs after the stream handler.")
	g.Line("func (m *hookRegistry[T]) AddAfterStream(hook ServerHookAfterStream[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookAfterStream(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfterStreamEmit adds a hook function that runs after the stream event is emitted.")
	g.Line("func (m *hookRegistry[T]) AddAfterStreamEmit(hook ServerHookAfterStreamEmit[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addHookAfterStreamEmit(hook)")
	})
	g.Line("}")
	g.Break()

	g.Line("// procRegistry handles procedure registration")
	g.Line("type procRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		name := procNode.Name
		namePascal := strutil.ToPascalCase(name)

		g.Linef("// Set%sInputHandler registers the input handler for the %s procedure", namePascal, name)
		g.Linef("func (p *procRegistry[T]) Set%sInputHandler(", namePascal)
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
			g.Linef("p.intServer.setProcInputHandler(\"%s\", typedInputHandler)", namePascal)
		})
		g.Line("}")
		g.Break()

		g.Linef("// Set%sHandler registers the handler for the %s procedure", namePascal, name)
		g.Linef("func (p *procRegistry[T]) Set%sHandler(", namePascal)
		g.Block(func() {
			g.Linef("handler func(ctx context.Context, ufoCtx T, input %sInput) (%sOutput, error),", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef(" p.intServer.setProcHandler(\"%s\", func(ctx context.Context, ufoCtx T, rawInput json.RawMessage) (any, error) {", namePascal)
			g.Block(func() {
				g.Line("var preTypedInput pre" + namePascal + "Input")
				g.Line("if err := json.Unmarshal(rawInput, &preTypedInput); err != nil {")
				g.Block(func() {
					g.Linef(`return nil, fmt.Errorf("failed to unmarshal %s input: %%w", err)`, namePascal)
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

				g.Linef("if inputHandler, ok := p.intServer.procInputHandlers[\"%s\"]; ok {", namePascal)
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

	g.Line("// streamRegistry handles stream registration")
	g.Line("type streamRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	for _, streamNode := range sch.GetStreamNodes() {
		name := streamNode.Name
		namePascal := strutil.ToPascalCase(name)

		g.Linef("// Set%sInputHandler registers the input handler for the %s stream", namePascal, name)
		g.Linef("func (s *streamRegistry[T]) Set%sInputHandler(", namePascal)
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
			g.Linef("s.intServer.setStreamInputHandler(\"%s\", typedInputHandler)", namePascal)
		})
		g.Line("}")
		g.Break()

		g.Linef("// Set%sHandler registers the handler for the %s stream", namePascal, name)
		g.Linef("func (s *streamRegistry[T]) Set%sHandler(", namePascal)
		g.Block(func() {
			g.Linef("handler func(ctx context.Context, ufoCtx T, input %sInput, emit func(%sOutput) error) error,", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef("s.intServer.setStreamHandler(\"%s\", func(ctx context.Context, ufoCtx T, rawInput json.RawMessage, rawEmit func(any) error) error {", namePascal)
			g.Block(func() {
				g.Line("var preTypedInput pre" + namePascal + "Input")
				g.Line("if err := json.Unmarshal(rawInput, &preTypedInput); err != nil {")
				g.Block(func() {
					g.Linef(`return fmt.Errorf("failed to unmarshal %s input: %%w", err)`, namePascal)
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
				g.Line("typedEmit := func(output " + namePascal + "Output) error {")
				g.Block(func() {
					g.Line("return rawEmit(output)")
				})
				g.Line("}")
				g.Break()

				g.Linef("if inputHandler, ok := s.intServer.streamInputHandlers[\"%s\"]; ok {", namePascal)
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

	g.Line("// Server wraps the UFO RPC internal generic server with organized, type-safe methods")
	g.Line("type Server[T any] struct {")
	g.Block(func() {
		g.Line("intServer  *internalServer[T]")
		g.Line("Hooks      *hookRegistry[T]")
		g.Line("Procs      *procRegistry[T]")
		g.Line("Streams    *streamRegistry[T]")
	})
	g.Line("}")
	g.Break()

	g.Line("// NewServer creates a new UFO RPC server that handles all procedures and streams")
	g.Line("//")
	g.Line("// The generic type T represents the context type, used to pass additional data")
	g.Line("// to procedures, such as authentication information, user session or any")
	g.Line("// other data you want to pass to procedures before they are executed")
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

	g.Line("// HandleRequest processes an incoming RPC request and manages the complete request-response lifecycle")
	g.Line("//")
	g.Line("// This method handles all aspects of request processing including:")
	g.Line("//	- Request parsing and validation")
	g.Line("//	- Context management")
	g.Line("//	- Response formatting and delivery")
	g.Line("//")
	g.Line("// To use this method, you need to implement the ServerHTTPAdapter interface which provides")
	g.Line("// the necessary methods for UFO RPC to interact with your HTTP server implementation")
	g.Line("//")
	g.Line("//	- For standard net/http implementations, use NewServerNetHTTPAdapter")
	g.Line("//	- For custom HTTP servers (gin, echo, etc.), implement the ServerHTTPAdapter interface with your specific logic")
	g.Line("func (s *Server[T]) HandleRequest(ctx context.Context, ufoCtx T, httpAdapter ServerHTTPAdapter) error {")
	g.Block(func() {
		g.Line("return s.intServer.handleRequest(ctx, ufoCtx, httpAdapter)")
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
