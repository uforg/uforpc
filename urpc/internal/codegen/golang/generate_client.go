package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

//go:embed pieces/client.go
var clientRawPiece string

func generateClient(sch schema.Schema, config Config) (string, error) {
	if !config.IncludeClient {
		return "", nil
	}

	piece := strutil.GetStrAfter(clientRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("client.go: could not find start delimiter")
	}

	g := genkit.NewGenKit().WithTabs()

	g.Raw(piece)
	g.Break()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Client generated implementation")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	g.Line("// clientBuilder provides a fluent API for configuring the UFO RPC client before instantiation.")
	g.Line("//")
	g.Line("// A builder is obtained by calling NewClient(baseURL), then optional")
	g.Line("// configuration methods can be chained before calling Build() to obtain a *Client ready for use.")
	g.Line("type clientBuilder struct {")
	g.Block(func() {
		g.Line("baseURL string")
		g.Line("opts    []internalClientOption")
	})
	g.Line("}")
	g.Break()

	g.Line("// NewClient creates a new clientBuilder instance. The baseURL parameter is mandatory and")
	g.Line("// should contain the root URL of your UFO RPC server (e.g., \"https://api.example.com/urpc\").")
	g.Line("//")
	g.Line("// Example usage:")
	g.Line("//   cli := NewClient(\"https://api.example.com\").WithHTTPClient(myHTTP).Build()")
	g.Line("func NewClient(baseURL string) *clientBuilder {")
	g.Block(func() {
		g.Line("return &clientBuilder{baseURL: baseURL, opts: []internalClientOption{}}")
	})
	g.Line("}")
	g.Break()

	g.Line("// WithBaseURL sets the server base URL used for every request.")
	g.Line("func (b *clientBuilder) WithBaseURL(u string) *clientBuilder {")
	g.Block(func() {
		g.Line("b.baseURL = u")
		g.Line("return b")
	})
	g.Line("}")
	g.Break()

	g.Line("// WithHTTPClient supplies a custom *http.Client (e.g., with timeouts or custom transport).")
	g.Line("func (b *clientBuilder) WithHTTPClient(hc *http.Client) *clientBuilder {")
	g.Block(func() {
		g.Line("b.opts = append(b.opts, withHTTPClient(hc))")
		g.Line("return b")
	})
	g.Line("}")
	g.Break()

	g.Line("// WithGlobalHeaders sets HTTP headers that will be sent with every request.")
	g.Line("func (b *clientBuilder) WithGlobalHeaders(h http.Header) *clientBuilder {")
	g.Block(func() {
		g.Line("b.opts = append(b.opts, withGlobalHeaders(h))")
		g.Line("return b")
	})
	g.Line("}")
	g.Break()

	g.Line("// WithMaxStreamEventDataSize overrides the default maximum size (bytes) for SSE event payloads.")
	g.Line("func (b *clientBuilder) WithMaxStreamEventDataSize(size int) *clientBuilder {")
	g.Block(func() {
		g.Line("b.opts = append(b.opts, withMaxStreamEventDataSize(size))")
		g.Line("return b")
	})
	g.Line("}")
	g.Break()

	g.Line("// Build constructs the *Client using the configured options.")
	g.Line("func (b *clientBuilder) Build() *Client {")
	g.Block(func() {
		g.Line("intClient := newInternalClient(b.baseURL, ufoProcedureNames, ufoStreamNames, b.opts...)")
		g.Line("return &Client{Procs: &clientProcRegistry{intClient: intClient}, Streams: &clientStreamRegistry{intClient: intClient}}")
	})
	g.Line("}")
	g.Break()

	g.Line("// Client provides a high-level, type-safe interface for invoking RPC procedures and streams.")
	g.Line("type Client struct {")
	g.Block(func() {
		g.Line("Procs     *clientProcRegistry")
		g.Line("Streams   *clientStreamRegistry")
	})
	g.Line("}")
	g.Break()

	// -----------------------------------------------------------------------------
	// Generate procedure wrappers
	// -----------------------------------------------------------------------------

	g.Line("type clientProcRegistry struct {")
	g.Block(func() {
		g.Line("intClient *internalClient")
	})
	g.Line("}")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		name := strutil.ToPascalCase(procNode.Name)
		builderName := "clientBuilder" + name

		// Client method to create builder
		g.Linef("// %s creates a call builder for the %s procedure.", name, name)
		g.Linef("func (registry *clientProcRegistry) %s() *%s {", name, builderName)
		g.Block(func() {
			g.Linef("return &%s{client: registry.intClient, headers: http.Header{}, name: \"%s\"}", builderName, name)
		})
		g.Line("}")
		g.Break()

		// Builder struct
		g.Linef("// %s represents a fluent call builder for the %s procedure.", builderName, name)
		g.Linef("type %s struct {", builderName)
		g.Block(func() {
			g.Line("client  *internalClient")
			g.Line("headers http.Header")
			g.Line("name    string")
		})
		g.Line("}")
		g.Break()

		// WithHeader method
		g.Linef("// WithHeader adds a single HTTP header to the %s invocation.", name)
		g.Linef("func (b *%s) WithHeader(key, value string) *%s {", builderName, builderName)
		g.Block(func() {
			g.Line("b.headers.Add(key, value)")
			g.Line("return b")
		})
		g.Line("}")
		g.Break()

		// WithHeaders method
		g.Linef("// WithHeaders merges the provided headers into the %s invocation.", name)
		g.Linef("func (b *%s) WithHeaders(h http.Header) *%s {", builderName, builderName)
		g.Block(func() {
			g.Line("clientHelperAddHeaders(b.headers, h)")
			g.Line("return b")
		})
		g.Line("}")
		g.Break()

		// Execute method
		g.Linef("// Execute performs the %s RPC call.", name)
		g.Linef("func (b *%s) Execute(ctx context.Context, input %sInput) (%sOutput, error) {", builderName, name, name)
		g.Block(func() {
			g.Line("raw := b.client.callProc(ctx, b.name, input, b.headers)")

			g.Line("if !raw.Ok {")
			g.Block(func() {
				g.Linef("return %sOutput{}, raw.Error", name)
			})
			g.Line("}")

			g.Linef("var out %sOutput", name)
			g.Line("if err := json.Unmarshal(raw.Output, &out); err != nil {")
			g.Block(func() {
				g.Linef("return %sOutput{}, Error{Message: fmt.Sprintf(\"failed to decode %s output: %%v\", err)}", name, name)
			})
			g.Line("}")

			g.Line("return out, nil")
		})
		g.Line("}")
		g.Break()
	}

	// -----------------------------------------------------------------------------
	// Generate stream wrappers
	// -----------------------------------------------------------------------------

	g.Line("type clientStreamRegistry struct {")
	g.Block(func() {
		g.Line("intClient *internalClient")
	})
	g.Line("}")
	g.Break()

	for _, streamNode := range sch.GetStreamNodes() {
		name := strutil.ToPascalCase(streamNode.Name)
		builderStream := "clientBuilder" + name + "Stream"

		// Client method to create stream builder
		g.Linef("// %s creates a stream builder for the %s stream.", name, name)
		g.Linef("func (registry *clientStreamRegistry) %s() *%s {", name, builderStream)
		g.Block(func() {
			g.Linef("return &%s{client: registry.intClient, headers: http.Header{}, name: \"%s\"}", builderStream, name)
		})
		g.Line("}")
		g.Break()

		g.Linef("// %s represents a fluent call builder for the %s stream.", builderStream, name)
		g.Linef("type %s struct {", builderStream)
		g.Block(func() {
			g.Line("client  *internalClient")
			g.Line("headers http.Header")
			g.Line("name    string")
			g.Line("maxEvt  int")
		})
		g.Line("}")
		g.Break()

		// WithHeader
		g.Linef("// WithHeader adds a single HTTP header to the %s stream subscription.", name)
		g.Linef("func (b *%s) WithHeader(key, value string) *%s {", builderStream, builderStream)
		g.Block(func() {
			g.Line("b.headers.Add(key, value)")
			g.Line("return b")
		})
		g.Line("}")
		g.Break()

		// WithHeaders
		g.Linef("// WithHeaders merges the provided headers into the %s stream subscription.", name)
		g.Linef("func (b *%s) WithHeaders(h http.Header) *%s {", builderStream, builderStream)
		g.Block(func() {
			g.Line("clientHelperAddHeaders(b.headers, h)")
			g.Line("return b")
		})
		g.Line("}")
		g.Break()

		// WithMaxEventDataSize
		g.Linef("// WithMaxEventDataSize sets the max allowed SSE payload size for this subscription.")
		g.Linef("func (b *%s) WithMaxEventDataSize(size int) *%s {", builderStream, builderStream)
		g.Block(func() {
			g.Line("b.maxEvt = size")
			g.Line("return b")
		})
		g.Line("}")
		g.Break()

		// Execute
		g.Linef("// Execute opens the %s stream and returns a typed event channel.", name)
		g.Linef("func (b *%s) Execute(ctx context.Context, input %sInput) <-chan Response[%sOutput] {", builderStream, name, name)
		g.Block(func() {
			g.Line("rawCh := b.client.stream(ctx, b.name, input, b.headers, b.maxEvt)")
			g.Linef("outCh := make(chan Response[%sOutput])", name)
			g.Line("go func() {")
			g.Block(func() {
				g.Line("for evt := range rawCh {")
				g.Block(func() {
					g.Line("if !evt.Ok {")
					g.Block(func() {
						g.Linef("outCh <- Response[%sOutput]{Ok: false, Error: evt.Error}", name)
					})
					g.Line("continue")
					g.Line("}")
					g.Linef("var out %sOutput", name)
					g.Line("if err := json.Unmarshal(evt.Output, &out); err != nil {")
					g.Block(func() {
						g.Linef("outCh <- Response[%sOutput]{Ok: false, Error: Error{Message: fmt.Sprintf(\"failed to decode %s output: %%v\", err)}}", name, name)
					})
					g.Line("continue")
					g.Line("}")
					g.Linef("outCh <- Response[%sOutput]{Ok: true, Output: out}", name)
				})
				g.Line("}")
				g.Line("close(outCh)")
			})
			g.Line("}()")
			g.Linef("return outCh")
		})
		g.Line("}")
		g.Break()
	}

	return g.String(), nil
}
