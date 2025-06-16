package typescript

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

//go:embed pieces/client.ts
var clientRawPiece string

func generateClient(sch schema.Schema, config Config) (string, error) {
	if !config.IncludeClient {
		return "", nil
	}

	piece := strutil.GetStrAfter(clientRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("client.ts: could not find start delimiter")
	}

	g := genkit.NewGenKit().WithTabs()

	g.Raw(piece)
	g.Break()

	g.Line("// =============================================================================")
	g.Line("// Generated Client Implementation")
	g.Line("// =============================================================================")
	g.Break()

	// Generate main client builder function
	generateClientBuilder(g)
	g.Break()

	// Generate main Client class
	generateClientClass(g)
	g.Break()

	// Generate procedure registry and builders
	generateProcedureImplementation(g, sch)
	g.Break()

	// Generate stream registry and builders
	generateStreamImplementation(g, sch)

	return g.String(), nil
}

// generateClientBuilder creates the main NewClient function
func generateClientBuilder(g *genkit.GenKit) {
	g.Line("/**")
	g.Line(" * Creates a new UFO RPC client builder.")
	g.Line(" *")
	g.Line(" * @param baseURL - The base URL for the RPC endpoint")
	g.Line(" * @returns A fluent builder for configuring the client")
	g.Line(" *")
	g.Line(" * @example")
	g.Line(" * ```typescript")
	g.Line(" * const client = NewClient(\"https://api.example.com/v1/urpc\")")
	g.Line(" *   .withGlobalHeader(\"Authorization\", \"Bearer token\")")
	g.Line(" *   .build();")
	g.Line(" * ```")
	g.Line(" */")
	g.Line("export function NewClient(baseURL: string): ClientBuilder {")
	g.Block(func() {
		g.Line("return new ClientBuilder(baseURL);")
	})
	g.Line("}")
	g.Break()

	g.Line("/**")
	g.Line(" * Fluent builder for configuring UFO RPC client options.")
	g.Line(" */")
	g.Line("class ClientBuilder {")
	g.Block(func() {
		g.Line("private builder: clientBuilder;")
		g.Break()

		g.Line("constructor(baseURL: string) {")
		g.Block(func() {
			g.Line("this.builder = new clientBuilder(baseURL);")
		})
		g.Line("}")
		g.Break()

		g.Line("/**")
		g.Line(" * Sets a custom fetch function for HTTP requests.")
		g.Line(" * Useful for environments without global fetch or for custom configurations.")
		g.Line(" */")
		g.Line("withCustomFetch(fetchFn: typeof fetch): ClientBuilder {")
		g.Block(func() {
			g.Line("this.builder.withFetch(fetchFn);")
			g.Line("return this;")
		})
		g.Line("}")
		g.Break()

		g.Line("/**")
		g.Line(" * Adds a global header that will be sent with every request.")
		g.Line(" * Can be called multiple times to set different headers.")
		g.Line(" */")
		g.Line("withGlobalHeader(key: string, value: string): ClientBuilder {")
		g.Block(func() {
			g.Line("this.builder.withGlobalHeader(key, value);")
			g.Line("return this;")
		})
		g.Line("}")
		g.Break()

		g.Line("/**")
		g.Line(" * Builds the configured client instance.")
		g.Line(" * @returns A fully configured Client ready for use")
		g.Line(" */")
		g.Line("build(): Client {")
		g.Block(func() {
			g.Line("const intClient = this.builder.build(ufoProcedureNames, ufoStreamNames);")
			g.Line("return new Client(intClient);")
		})
		g.Line("}")
	})
	g.Line("}")
}

// generateClientClass creates the main Client class
func generateClientClass(g *genkit.GenKit) {
	g.Line("/**")
	g.Line(" * Main UFO RPC client providing type-safe access to procedures and streams.")
	g.Line(" */")
	g.Line("export class Client {")
	g.Block(func() {
		g.Line("/** Registry for accessing RPC procedures */")
		g.Line("public readonly procs: ProcRegistry;")
		g.Break()
		g.Line("/** Registry for accessing RPC streams */")
		g.Line("public readonly streams: StreamRegistry;")
		g.Break()

		g.Line("constructor(private intClient: internalClient) {")
		g.Block(func() {
			g.Line("this.procs = new ProcRegistry(intClient);")
			g.Line("this.streams = new StreamRegistry(intClient);")
		})
		g.Line("}")
	})
	g.Line("}")
}

// generateProcedureImplementation generates all procedure-related code
func generateProcedureImplementation(g *genkit.GenKit, sch schema.Schema) {
	g.Line("// =============================================================================")
	g.Line("// Procedure Implementation")
	g.Line("// =============================================================================")
	g.Break()

	// Generate procedure registry
	g.Line("/**")
	g.Line(" * Registry providing access to all RPC procedures.")
	g.Line(" */")
	g.Line("class ProcRegistry {")
	g.Block(func() {
		g.Line("constructor(private intClient: internalClient) {}")
		g.Break()

		// Generate method for each procedure
		for _, procNode := range sch.GetProcNodes() {
			name := strutil.ToPascalCase(procNode.Name)
			builderName := fmt.Sprintf("builder%s", name)

			g.Linef("/**")
			g.Linef(" * Creates a call builder for the %s procedure.", name)
			renderDeprecated(g, procNode.Deprecated)
			g.Linef(" */")
			g.Linef("%s(): %s {", strutil.ToCamelCase(procNode.Name), builderName)
			g.Block(func() {
				g.Linef("return new %s(this.intClient, \"%s\");", builderName, procNode.Name)
			})
			g.Line("}")
			g.Break()
		}
	})
	g.Line("}")
	g.Break()

	// Generate individual procedure builders
	for _, procNode := range sch.GetProcNodes() {
		name := strutil.ToPascalCase(procNode.Name)
		builderName := fmt.Sprintf("builder%s", name)
		inputType := fmt.Sprintf("%sInput", name)
		outputType := fmt.Sprintf("%sOutput", name)

		g.Linef("/**")
		g.Linef(" * Fluent builder for the %s procedure.", name)
		if procNode.Deprecated != nil && *procNode.Deprecated != "" {
			g.Linef(" * @deprecated %s", *procNode.Deprecated)
		}
		g.Linef(" */")
		g.Linef("class %s {", builderName)
		g.Block(func() {
			g.Line("private headers: Record<string, string> = {};")
			g.Break()

			g.Line("constructor(")
			g.Block(func() {
				g.Line("private intClient: internalClient,")
				g.Line("private procName: string")
			})
			g.Line(") {}")
			g.Break()

			g.Line("/**")
			g.Linef(" * Adds a custom header to the %s request.", name)
			g.Line(" * Can be called multiple times to set different headers.")
			g.Line(" */")
			g.Linef("withHeader(key: string, value: string): %s {", builderName)
			g.Block(func() {
				g.Line("this.headers[key] = value;")
				g.Line("return this;")
			})
			g.Line("}")
			g.Break()

			g.Line("/**")
			g.Linef(" * Executes the %s procedure.", name)
			g.Line(" *")
			g.Linef(" * @param input - The %s input parameters", name)
			g.Linef(" * @returns Promise resolving to %s or throws UfoError if something went wrong", outputType)
			g.Line(" */")
			g.Linef("async execute(input: %s): Promise<%s> {", inputType, outputType)
			g.Block(func() {
				g.Line("const rawResponse = await this.intClient.callProc(")
				g.Block(func() {
					g.Line("this.procName,")
					g.Line("input,")
					g.Line("this.headers")
				})
				g.Line(");")

				g.Line("if (!rawResponse.ok) throw rawResponse.error;")
				g.Linef("return rawResponse.output as %s;", outputType)
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()
	}
}

// generateStreamImplementation generates all stream-related code
func generateStreamImplementation(g *genkit.GenKit, sch schema.Schema) {
	g.Line("// =============================================================================")
	g.Line("// Stream Implementation")
	g.Line("// =============================================================================")
	g.Break()

	// Generate stream registry
	g.Line("/**")
	g.Line(" * Registry providing access to all RPC streams.")
	g.Line(" */")
	g.Line("class StreamRegistry {")
	g.Block(func() {
		g.Line("constructor(private intClient: internalClient) {}")
		g.Break()

		// Generate method for each stream
		for _, streamNode := range sch.GetStreamNodes() {
			name := strutil.ToPascalCase(streamNode.Name)
			builderName := fmt.Sprintf("builder%sStream", name)

			g.Linef("/**")
			g.Linef(" * Creates a stream builder for the %s stream.", name)
			renderDeprecated(g, streamNode.Deprecated)
			g.Linef(" */")
			g.Linef("%s(): %s {", strutil.ToCamelCase(streamNode.Name), builderName)
			g.Block(func() {
				g.Linef("return new %s(this.intClient, \"%s\");", builderName, streamNode.Name)
			})
			g.Line("}")
			g.Break()
		}
	})
	g.Line("}")
	g.Break()

	// Generate individual stream builders
	for _, streamNode := range sch.GetStreamNodes() {
		name := strutil.ToPascalCase(streamNode.Name)
		builderName := fmt.Sprintf("builder%sStream", name)
		inputType := fmt.Sprintf("%sInput", name)
		outputType := fmt.Sprintf("%sOutput", name)

		g.Linef("/**")
		g.Linef(" * Fluent builder for the %s stream.", name)
		if streamNode.Deprecated != nil && *streamNode.Deprecated != "" {
			g.Linef(" * @deprecated %s", *streamNode.Deprecated)
		}
		g.Linef(" */")
		g.Linef("class %s {", builderName)
		g.Block(func() {
			g.Line("private headers: Record<string, string> = {};")
			g.Break()

			g.Line("constructor(")
			g.Block(func() {
				g.Line("private intClient: internalClient,")
				g.Line("private streamName: string")
			})
			g.Line(") {}")
			g.Break()

			g.Line("/**")
			g.Linef(" * Adds a custom header to the %s stream request.", name)
			g.Line(" * Can be called multiple times to set different headers.")
			g.Line(" */")
			g.Linef("withHeader(key: string, value: string): %s {", builderName)
			g.Block(func() {
				g.Line("this.headers[key] = value;")
				g.Line("return this;")
			})
			g.Line("}")
			g.Break()

			g.Line("/**")
			g.Linef(" * Opens the %s Server-Sent Events stream.", name)
			g.Line(" *")
			g.Linef(" * @param input - The %s input parameters", name)
			g.Line(" * @returns Object containing:")
			g.Linef(" *   - stream: AsyncGenerator yielding Response<%s> events", outputType)
			g.Line(" *   - cancel: Function for cancelling the stream")
			g.Line(" *")
			g.Line(" * @example")
			g.Line(" * ```typescript")
			g.Linef(" * const { stream, cancel } = client.streams.%s().execute(input);", strutil.ToCamelCase(streamNode.Name))
			g.Line(" * ")
			g.Line(" * // All stream events are received here")
			g.Line(" * for await (const event of stream) {")
			g.Line(" *   if (event.ok) {")
			g.Line(" *     console.log('Received:', event.output);")
			g.Line(" *   } else {")
			g.Line(" *     console.error('Error:', event.error);")
			g.Line(" *   }")
			g.Line(" * }")
			g.Line(" * ")
			g.Line(" * // Cancel the stream when needed")
			g.Line(" * cancel();")
			g.Line(" * ```")
			g.Line(" */")
			g.Linef("execute(input: %s): {", inputType)
			g.Block(func() {
				g.Linef("stream: AsyncGenerator<Response<%s>, void, unknown>;", outputType)
				g.Line("cancel: () => void;")
			})
			g.Line("} {")
			g.Block(func() {
				g.Line("const { stream, cancel } = this.intClient.callStream(")
				g.Block(func() {
					g.Line("this.streamName,")
					g.Line("input,")
					g.Line("this.headers")
				})
				g.Line(");")

				g.Linef("const typedStream = async function* (): AsyncGenerator<Response<%s>, void, unknown> {", outputType)
				g.Block(func() {
					g.Line("for await (const event of stream) {")
					g.Block(func() {
						g.Linef("yield event as Response<%s>;", outputType)
					})
					g.Line("}")
				})
				g.Line("};")
				g.Break()

				g.Line("return {")
				g.Block(func() {
					g.Line("stream: typedStream(),")
					g.Line("cancel: cancel")
				})
				g.Line("};")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()
	}
}
