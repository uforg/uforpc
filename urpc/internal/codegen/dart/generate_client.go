package dart

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

func generateClient(sch schema.Schema, _ Config) (string, error) {
	g := genkit.NewGenKit().WithSpaces(2)

	g.Line("// =============================================================================")
	g.Line("// Generated Client Implementation")
	g.Line("// =============================================================================")
	g.Break()

	generateClientBuilder(g)
	g.Break()

	generateClientClass(g)
	g.Break()

	generateProcedureImplementation(g, sch)
	g.Break()

	generateStreamImplementation(g, sch)
	g.Break()

	return g.String(), nil
}

func generateClientBuilder(g *genkit.GenKit) {
	g.Line("/// Creates a new UFO RPC client builder.")
	g.Line("_ClientBuilder NewClient(String baseURL) => _ClientBuilder(baseURL);")
	g.Break()

	g.Line("/// Fluent builder for configuring UFO RPC client options.")
	g.Line("class _ClientBuilder {")
	g.Block(func() {
		g.Line("final _InternalClientBuilder _builder;")
		g.Break()
		g.Line("_ClientBuilder(String baseURL) : _builder = _InternalClientBuilder(baseURL);")
		g.Break()
		g.Line("/// Sets a custom fetch-like function for HTTP requests.")
		g.Line("_ClientBuilder withCustomFetch(_FetchLike fetchFn) { _builder.withFetch(fetchFn); return this; }")
		g.Break()
		g.Line("/// Adds a global header that will be sent with every request.")
		g.Line("_ClientBuilder withGlobalHeader(String key, String value) { _builder.withGlobalHeader(key, value); return this; }")
		g.Break()
		g.Line("/// Builds the configured client instance.")
		g.Line("Client build() { final intClient = _builder.build(__ufoProcedureNames, __ufoStreamNames); return Client._internal(intClient); }")
	})
	g.Line("}")
}

func generateClientClass(g *genkit.GenKit) {
	g.Line("/// Main UFO RPC client providing type-safe access to procedures and streams.")
	g.Line("class Client {")
	g.Block(func() {
		g.Line("final _ProcRegistry procs;")
		g.Line("final _StreamRegistry streams;")
		g.Break()
		g.Line("Client._internal(_InternalClient intClient) : procs = _ProcRegistry(intClient), streams = _StreamRegistry(intClient);")
	})
	g.Line("}")
}

func generateProcedureImplementation(g *genkit.GenKit, sch schema.Schema) {
	g.Line("// =============================================================================")
	g.Line("// Procedure Implementation")
	g.Line("// =============================================================================")
	g.Break()

	g.Line("/// Registry providing access to all RPC procedures.")
	g.Line("class _ProcRegistry {")
	g.Block(func() {
		g.Line("final _InternalClient _intClient;")
		g.Line("_ProcRegistry(this._intClient);")
		g.Break()
		for _, procNode := range sch.GetProcNodes() {
			name := strutil.ToPascalCase(procNode.Name)
			builderName := fmt.Sprintf("_Builder%s", name)
			g.Linef("/// Creates a call builder for the %s procedure.", name)
			renderDeprecatedDart(g, procNode.Deprecated)
			g.Linef("%s %s() => %s(_intClient, '%s');", builderName, strutil.ToCamelCase(procNode.Name), builderName, procNode.Name)
			g.Break()
		}
	})
	g.Line("}")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		name := strutil.ToPascalCase(procNode.Name)
		builderName := fmt.Sprintf("_Builder%s", name)
		hydrateFuncName := fmt.Sprintf("%sOutput.fromJson", name)
		inputType := fmt.Sprintf("%sInput", name)
		outputType := fmt.Sprintf("%sOutput", name)

		g.Linef("/// Fluent builder for the %s procedure.", name)
		if procNode.Deprecated != nil && *procNode.Deprecated != "" {
			g.Linef("/// @deprecated %s", *procNode.Deprecated)
		}
		g.Linef("class %s {", builderName)
		g.Block(func() {
			g.Line("final _InternalClient _intClient;")
			g.Line("final String _procName;")
			g.Line("final Map<String, String> _headers = {};")
			g.Line("_RetryConfig? _retryConfig;")
			g.Line("_TimeoutConfig? _timeoutConfig;")
			g.Break()
			g.Linef("%s(this._intClient, this._procName);", builderName)
			g.Break()
			g.Linef("%s withHeader(String key, String value) { _headers[key] = value; return this; }", builderName)
			g.Break()
			g.Linef("%s withRetries(_RetryConfig config) { _retryConfig = _RetryConfig.sanitised(config); return this; }", builderName)
			g.Break()
			g.Linef("%s withTimeout(_TimeoutConfig config) { _timeoutConfig = _TimeoutConfig.sanitised(config); return this; }", builderName)
			g.Break()
			g.Linef("Future<%s> execute(%s input) async {", outputType, inputType)
			g.Block(func() {
				g.Line("final rawResponse = await _intClient.callProc(_procName, input.toJson(), _headers, _retryConfig, _timeoutConfig);")
				g.Line("if (!rawResponse.ok) { throw rawResponse.error!; }")
				g.Linef("final out = %s((rawResponse.output as Map).cast<String, dynamic>());", hydrateFuncName)
				g.Line("return out;")
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()
	}
}

func generateStreamImplementation(g *genkit.GenKit, sch schema.Schema) {
	g.Line("// =============================================================================")
	g.Line("// Stream Implementation")
	g.Line("// =============================================================================")
	g.Break()

	g.Line("/// Registry providing access to all RPC streams.")
	g.Line("class _StreamRegistry {")
	g.Block(func() {
		g.Line("final _InternalClient _intClient;")
		g.Line("_StreamRegistry(this._intClient);")
		g.Break()
		for _, streamNode := range sch.GetStreamNodes() {
			name := strutil.ToPascalCase(streamNode.Name)
			builderName := fmt.Sprintf("_Builder%sStream", name)
			g.Linef("/// Creates a stream builder for the %s stream.", name)
			renderDeprecatedDart(g, streamNode.Deprecated)
			g.Linef("%s %s() => %s(_intClient, '%s');", builderName, strutil.ToCamelCase(streamNode.Name), builderName, streamNode.Name)
			g.Break()
		}
	})
	g.Line("}")
	g.Break()

	for _, streamNode := range sch.GetStreamNodes() {
		name := strutil.ToPascalCase(streamNode.Name)
		builderName := fmt.Sprintf("_Builder%sStream", name)
		hydrateFuncName := fmt.Sprintf("%sOutput.fromJson", name)
		inputType := fmt.Sprintf("%sInput", name)
		outputType := fmt.Sprintf("%sOutput", name)

		g.Linef("/// Fluent builder for the %s stream.", name)
		if streamNode.Deprecated != nil && *streamNode.Deprecated != "" {
			g.Linef("/// @deprecated %s", *streamNode.Deprecated)
		}
		g.Linef("class %s {", builderName)
		g.Block(func() {
			g.Line("final _InternalClient _intClient;")
			g.Line("final String _streamName;")
			g.Line("final Map<String, String> _headers = {};")
			g.Line("_ReconnectConfig? _reconnectConfig;")
			g.Break()
			g.Linef("%s(this._intClient, this._streamName);", builderName)
			g.Break()
			g.Linef("%s withHeader(String key, String value) { _headers[key] = value; return this; }", builderName)
			g.Break()
			g.Linef("%s withReconnect(_ReconnectConfig config) { _reconnectConfig = _ReconnectConfig.sanitised(config); return this; }", builderName)
			g.Break()
			g.Linef("_StreamHandle<%s> execute(%s input) {", outputType, inputType)
			g.Block(func() {
				g.Line("final handle = _intClient.callStream(_streamName, input.toJson(), _headers, _reconnectConfig);")
				g.Linef("final typed = handle.stream.map((event) { if (event.ok) { final out = %s((event.output as Map).cast<String, dynamic>()); return Response<%s>.ok(out); } else { return Response<%s>.error(event.error!); } });", hydrateFuncName, outputType, outputType)
				g.Linef("return _StreamHandle<%s>(stream: typed, cancel: handle.cancel);", outputType)
			})
			g.Line("}")
		})
		g.Line("}")
		g.Break()
	}
}
