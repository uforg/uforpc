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

	g.Line("// middlewareRegistry handles middleware registration")
	g.Line("type middlewareRegistry[T any] struct {")
	g.Block(func() {
		g.Line("intServer *internalServer[T]")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddBefore adds a middleware function that runs before the handler. Both for procedures and streams.")
	g.Line("func (m *middlewareRegistry[T]) AddBefore(fn MiddlewareBefore[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addMiddlewareBefore(fn)")
	})
	g.Line("}")
	g.Break()

	g.Line("// AddAfter adds a middleware function that runs after the handler. Only for procedures, not for streams.")
	g.Line("func (m *middlewareRegistry[T]) AddAfter(fn MiddlewareAfter[T]) {")
	g.Block(func() {
		g.Line("m.intServer.addMiddlewareAfter(fn)")
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

		g.Linef("// Set%sHandler registers the handler for the %s procedure", namePascal, name)
		g.Linef("func (p *procRegistry[T]) Set%sHandler(", namePascal)
		g.Block(func() {
			g.Linef("handler func(context T, input %sInput) (%sOutput, error),", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef(" p.intServer.setProcHandler(\"%s\", func(context T, rawInput json.RawMessage) (any, error) {", namePascal)
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
				g.Line("return handler(context, typedInput)")
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

		g.Linef("// Set%sHandler registers the handler for the %s stream", namePascal, name)
		g.Linef("func (s *streamRegistry[T]) Set%sHandler(", namePascal)
		g.Block(func() {
			g.Linef("handler func(context T, input %sInput, emit func(%sOutput) error) error,", name, name)
		})
		g.Linef(") {")
		g.Block(func() {
			g.Linef("s.intServer.setStreamHandler(\"%s\", func(context T, rawInput json.RawMessage, rawEmit func(any) error) error {", namePascal)
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

				g.Line("return handler(context, typedInput, typedEmit)")
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
		g.Line("Middlewares     *middlewareRegistry[T]")
		g.Line("Procs           *procRegistry[T]")
		g.Line("Streams         *streamRegistry[T]")
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
			g.Line("intServer:    intServer,")
			g.Line("Middlewares:  &middlewareRegistry[T]{intServer: intServer},")
			g.Line("Procs:        &procRegistry[T]{intServer: intServer},")
			g.Line("Streams:      &streamRegistry[T]{intServer: intServer},")
		})
		g.Line("}")
	})
	g.Line("}")
	g.Break()

	g.Line("// HandleRequest processes an incoming RPC request")
	g.Line("//")
	g.Line("// It handles all the processing of the request and response, the only thing")
	g.Line("// you need to do is implement the ServerRequestResponseProvider interface")
	g.Line("// that allows UFO RPC to access the request and response resources it needs")
	g.Line("func (s *Server[T]) HandleRequest(requestResponseProvider ServerRequestResponseProvider[T]) error {")
	g.Block(func() {
		g.Line("return s.intServer.handleRequest(requestResponseProvider)")
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
