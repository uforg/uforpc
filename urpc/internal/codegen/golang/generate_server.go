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
				g.Linef("return func(c *HandlerContext[P, any]) (any, error) {")
				g.Block(func() {
					g.Linef("var pre pre%sInput", name)
					g.Line("raw, _ := c.Input.(json.RawMessage)")
					g.Line("if err := json.Unmarshal(raw, &pre); err != nil {")
					g.Block(func() { g.Linef("return nil, fmt.Errorf(\"failed to unmarshal %s input: %%w\", err)", name) })
					g.Line("}")
					g.Line("if err := pre.validate(); err != nil { return nil, err }")
					g.Linef("typedInput := pre.transform()")
					g.Linef("typedCtx := &HandlerContext[P, %sInput]{Props: c.Props, Input: typedInput, Context: c.Context, operationName: c.operationName, operationType: c.operationType}", name)
					g.Linef("nextTyped := func(tc *HandlerContext[P, %sInput]) (%sOutput, error) {", name, name)
					g.Block(func() {
						g.Line("outAny, err := next(&HandlerContext[P, any]{Props: tc.Props, Input: tc.Input, Context: tc.Context, operationName: tc.operationName, operationType: tc.operationType})")
						g.Linef("if err != nil { var zero %sOutput; return zero, err }", name)
						g.Linef("out, _ := outAny.(%sOutput)", name)
						g.Line("return out, nil")
					})
					g.Line("}")
					g.Line("typedComposed := mw(nextTyped)")
					g.Line("return typedComposed(typedCtx)")
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
			g.Linef("untyped := func(c *HandlerContext[P, any]) (any, error) {")
			g.Block(func() {
				g.Linef("var pre pre%sInput", name)
				g.Line("raw, _ := c.Input.(json.RawMessage)")
				g.Line("if err := json.Unmarshal(raw, &pre); err != nil {")
				g.Block(func() { g.Linef("return nil, fmt.Errorf(\"failed to unmarshal %s input: %%w\", err)", name) })
				g.Line("}")
				g.Line("if err := pre.validate(); err != nil { return nil, err }")
				g.Line("typedInput := pre.transform()")
				g.Linef("typedCtx := &HandlerContext[P, %sInput]{Props: c.Props, Input: typedInput, Context: c.Context, operationName: c.operationName, operationType: c.operationType}", name)
				g.Line("out, err := handler(typedCtx)")
				g.Line("if err != nil { return nil, err }")
				g.Line("return any(out), nil")
			})
			g.Line("}")
			g.Linef("e.intServer.setProcHandler(\"%s\", untyped)", name)
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
		g.Linef("type %sEmitFunc = EmitFunc[%sOutput]", name, name)
		g.Linef("type %sHandlerFunc[P any] func(c *%sHandlerContext[P], emit %sEmitFunc) error", name, name, name)
		g.Linef("type %sMiddlewareFunc[P any] func(next %sHandlerFunc[P]) %sHandlerFunc[P]", name, name, name)
		g.Linef("type %sEmitMiddlewareFunc[P any] func(c *%sHandlerContext[P], next %sEmitFunc) %sEmitFunc", name, name, name, name)
		g.Break()

		// Use (stream middleware)
		g.Linef("// Use registers a typed middleware for the %s stream.", name)
		g.Linef("func (e stream%sEntry[P]) Use(mw %sMiddlewareFunc[P]) {", name, name)
		g.Block(func() {
			g.Linef("adapted := func(next StreamHandlerFunc[P, any, any]) StreamHandlerFunc[P, any, any] {")
			g.Block(func() {
				g.Linef("return func(c *HandlerContext[P, any], emit EmitFunc[any]) error {")
				g.Block(func() {
					g.Linef("var pre pre%sInput", name)
					g.Line("raw, _ := c.Input.(json.RawMessage)")
					g.Line("if err := json.Unmarshal(raw, &pre); err != nil {")
					g.Block(func() { g.Linef("return fmt.Errorf(\"failed to unmarshal %s input: %%w\", err)", name) })
					g.Line("}")
					g.Line("if err := pre.validate(); err != nil { return err }")
					g.Linef("typedInput := pre.transform()")
					g.Linef("typedCtx := &HandlerContext[P, %sInput]{Props: c.Props, Input: typedInput, Context: c.Context, operationName: c.operationName, operationType: c.operationType}", name)
					g.Linef("emitTyped := func(o %sOutput) error { return emit(any(o)) }", name)
					g.Linef("nextTyped := func(tc *HandlerContext[P, %sInput], et EmitFunc[%sOutput]) error {", name, name)
					g.Block(func() {
						g.Line("return next(&HandlerContext[P, any]{Props: tc.Props, Input: tc.Input, Context: tc.Context, operationName: tc.operationName, operationType: tc.operationType}, func(o any) error { return et(o.(" + name + "Output)) })")
					})
					g.Line("}")
					g.Line("typedComposed := mw(nextTyped)")
					g.Line("return typedComposed(typedCtx, emitTyped)")
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
			g.Linef("adapted := func(c *HandlerContext[P, any], next EmitFunc[any]) EmitFunc[any] {")
			g.Block(func() {
				g.Linef("typedInput, _ := c.Input.(%sInput)", name)
				g.Linef("typedCtx := &%sHandlerContext[P]{Props: c.Props, Input: typedInput, Context: c.Context, operationName: c.operationName, operationType: c.operationType}", name)
				g.Linef("typedNext := func(o %sOutput) error { return next(any(o)) }", name)
				g.Line("typedWrapped := mw(typedCtx, typedNext)")
				g.Linef("return func(o any) error { return typedWrapped(o.(%sOutput)) }", name)
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
			g.Linef("untyped := func(c *HandlerContext[P, any], emit EmitFunc[any]) error {")
			g.Block(func() {
				g.Linef("var pre pre%sInput", name)
				g.Line("raw, _ := c.Input.(json.RawMessage)")
				g.Line("if err := json.Unmarshal(raw, &pre); err != nil {")
				g.Block(func() { g.Linef("return fmt.Errorf(\"failed to unmarshal %s input: %%w\", err)", name) })
				g.Line("}")
				g.Line("if err := pre.validate(); err != nil { return err }")
				g.Line("typedInput := pre.transform()")
				g.Linef("typedCtx := &HandlerContext[P, %sInput]{Props: c.Props, Input: typedInput, Context: c.Context, operationName: c.operationName, operationType: c.operationType}", name)
				g.Linef("emitTyped := func(o %sOutput) error { return emit(any(o)) }", name)
				g.Line("return handler(typedCtx, emitTyped)")
			})
			g.Line("}")
			g.Linef("e.intServer.setStreamHandler(\"%s\", untyped)", name)
		})
		g.Line("}")
		g.Break()
	}

	return g.String(), nil
}
