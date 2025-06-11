//nolint:unused
package pieces

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
)

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Server Types
// -----------------------------------------------------------------------------

const (
	ServerOperationTypeProc   = "proc"
	ServerOperationTypeStream = "stream"
)

// ServerHTTPAdapter defines the interface required by UFO RPC server to handle
// incoming HTTP requests and write responses to clients. This abstraction allows
// the server to work with different HTTP frameworks while maintaining the same
// core functionality.
//
// Implementations must provide methods to read request bodies, set response headers,
// write response data, and flush the response buffer to ensure immediate delivery
// to the client.
type ServerHTTPAdapter interface {
	// RequestBody returns the body reader for the incoming HTTP request.
	// The returned io.Reader allows the server to read the request payload
	// containing RPC call data.
	RequestBody() io.Reader

	// SetHeader sets a response header with the specified key-value pair.
	// This is used to configure response headers like Content-Type and
	// caching directives for both procedure and stream responses.
	SetHeader(key, value string)

	// Write writes the provided data to the response body.
	// Returns the number of bytes written and any error encountered.
	// For procedures, this writes the complete JSON response. For streams,
	// this writes individual Server-Sent Events data chunks.
	Write(data []byte) (int, error)

	// Flush immediately sends any buffered response data to the client.
	// This is crucial for streaming responses to ensure real-time delivery
	// of events. Returns an error if the flush operation fails.
	Flush() error
}

// ServerNetHTTPAdapter implements ServerHTTPAdapter for Go's standard net/http package.
// This adapter bridges the UFO RPC server with the standard HTTP library, allowing
// seamless integration with existing HTTP servers and middleware.
type ServerNetHTTPAdapter struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

// NewServerNetHTTPAdapter creates a new ServerNetHTTPAdapter that implements the
// ServerHTTPAdapter interface for net/http.
//
// Parameters:
//   - w: The http.ResponseWriter to write responses to
//   - r: The http.Request containing the incoming request data
//
// Returns a ServerHTTPAdapter implementation ready for use with UFO RPC server.
func NewServerNetHTTPAdapter(w http.ResponseWriter, r *http.Request) ServerHTTPAdapter {
	return &ServerNetHTTPAdapter{
		responseWriter: w,
		request:        r,
	}
}

// RequestBody returns the body reader for the HTTP request.
// This provides access to the request payload containing the RPC call data.
func (r *ServerNetHTTPAdapter) RequestBody() io.Reader {
	return r.request.Body
}

// SetHeader sets a response header with the specified key-value pair.
// This configures headers for the HTTP response, such as Content-Type
// for JSON responses or streaming-specific headers.
func (r *ServerNetHTTPAdapter) SetHeader(key, value string) {
	r.responseWriter.Header().Set(key, value)
}

// Write writes the provided data to the HTTP response body.
// Returns the number of bytes written and any error encountered during writing.
func (r *ServerNetHTTPAdapter) Write(data []byte) (int, error) {
	return r.responseWriter.Write(data)
}

// Flush immediately sends any buffered response data to the client.
// For streaming responses, this ensures real-time delivery of events.
// If the underlying ResponseWriter doesn't support flushing, this is a no-op.
func (r *ServerNetHTTPAdapter) Flush() error {
	if f, ok := r.responseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// ServerHookBeforeHandler defines a hook function that executes before
// any procedure or stream handler is invoked. This allows for cross-cutting
// concerns like authentication, logging, rate limiting, etc.
//
// The hook receives the current context, UFO context, operation type
// (ServerOperationTypeProc or ServerOperationTypeStream), and operation name.
// It must return the potentially modified context and UFO context, or an
// error to abort the request.
//
// If an error is returned, the request is terminated and an error response
// is sent to the client.
type ServerHookBeforeHandler[T any] func(
	ctx context.Context,
	ufoCtx T,
	operationType string,
	operationName string,
) (context.Context, T, error)

// ServerHookBeforeProcRespond defines a hook function that executes before
// the procedure handler responds to the client. This hook can modify the response
// that will be sent to the client, allowing for response transformation, filtering,
// or adding metadata.
//
// The hook receives the current context, UFO context, procedure name, and
// the response about to be sent. It must return the potentially modified
// response that will be delivered to the client.
type ServerHookBeforeProcRespond[T any] func(
	ctx context.Context,
	ufoCtx T,
	procName string,
	response Response[any],
) Response[any]

// ServerHookBeforeStreamEmit defines a hook function that executes before
// each stream event is emitted to the client. This hook can modify the response
// that will be sent as part of the Server-Sent Events stream.
//
// The hook receives the current context, UFO context, stream name, and
// the response about to be emitted. It must return the potentially modified
// response that will be sent to the client.
type ServerHookBeforeStreamEmit[T any] func(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	response Response[any],
) Response[any]

// ServerHookAfterProc defines a hook function that executes after
// a procedure request has been processed, regardless of success or failure.
// This allows for logging, metrics collection, or cleanup operations.
// This hook cannot modify the response as it has already been sent to the client.
//
// The hook receives the current context, UFO context, procedure name,
// and the response that was generated by the handler.
type ServerHookAfterProc[T any] func(
	ctx context.Context,
	ufoCtx T,
	procName string,
	response Response[any],
)

// ServerHookAfterStream defines a hook function that executes after
// a stream request processing ends, either successfully or with an error.
// This allows for cleanup operations, logging, metrics collection, or
// resource deallocation. This hook cannot modify anything as the stream
// has already completed.
//
// The hook receives the current context, UFO context, stream name,
// and any error that occurred during stream processing (nil if successful).
type ServerHookAfterStream[T any] func(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	err error,
)

// ServerHookAfterStreamEmit defines a hook function that executes after
// each stream event has been successfully emitted to the client.
// This allows for logging, metrics collection, or post-emission cleanup operations.
// This hook cannot modify anything as the event has already been sent.
//
// The hook receives the current context, UFO context, stream name,
// and the response that was sent to the client.
type ServerHookAfterStreamEmit[T any] func(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	response Response[any],
)

// -----------------------------------------------------------------------------
// Server Internal Implementation
// -----------------------------------------------------------------------------

// internalServer manages RPC request handling and hook execution for
// both procedures and streams. It maintains handler registrations, hook
// chains, and coordinates the complete request lifecycle.
//
// The generic type T represents the UFO context type, allowing users to pass
// custom data (authentication info, user sessions, etc.) through the entire
// request processing pipeline.
type internalServer[T any] struct {
	// procNames contains the list of all registered procedure names
	procNames []string
	// procNamesMap contains the list of all registered procedure names
	procNamesMap map[string]bool
	// streamNames contains the list of all registered stream names
	streamNames []string
	// streamNamesMap contains the list of all registered stream names
	streamNamesMap map[string]bool
	// operationNamesMap contains the list of all registered operation names
	// and its corresponding type
	operationNamesMap map[string]string
	// handlersMu protects all handler maps and hook slices from concurrent access
	handlersMu sync.RWMutex
	// procHandlers maps procedure names to their implementation functions
	procHandlers map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage) (any, error)
	// procInputHandlers maps procedure names to their input processing functions
	procInputHandlers map[string]func(ctx context.Context, ufoCtx T, input any) (any, error)
	// streamHandlers maps stream names to their implementation functions
	streamHandlers map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage, emit func(any) error) error
	// streamInputHandlers maps stream names to their input processing functions
	streamInputHandlers map[string]func(ctx context.Context, ufoCtx T, input any) (any, error)
	// hooksBeforeHandler contains hooks that run before any handler execution
	hooksBeforeHandler []ServerHookBeforeHandler[T]
	// hooksBeforeProcRespond contains hooks that run before procedure response is sent
	hooksBeforeProcRespond []ServerHookBeforeProcRespond[T]
	// hooksBeforeStreamEmit contains hooks that run before each stream event emission
	hooksBeforeStreamEmit []ServerHookBeforeStreamEmit[T]
	// hooksAfterProc contains hooks that run after procedure completion
	hooksAfterProc []ServerHookAfterProc[T]
	// hooksAfterStream contains hooks that run after stream completion
	hooksAfterStream []ServerHookAfterStream[T]
	// hooksAfterStreamEmit contains hooks that run after each stream event emission
	hooksAfterStreamEmit []ServerHookAfterStreamEmit[T]
}

// newInternalServer creates a new UFO RPC server instance with the specified
// procedure and stream names. The server is initialized with empty handler
// maps and hook slices, ready for registration.
//
// The generic type T represents the UFO context type, used to pass additional
// data to handlers, such as authentication information, user sessions, or any
// other request-scoped data.
//
// Parameters:
//   - procNames: List of procedure names that this server will handle
//   - streamNames: List of stream names that this server will handle
//
// Returns a new internalServer instance ready for handler and hook registration.
func newInternalServer[T any](
	procNames []string,
	streamNames []string,
) *internalServer[T] {
	procNamesMap := make(map[string]bool)
	streamNamesMap := make(map[string]bool)
	operationNames := make(map[string]string)
	for _, procName := range procNames {
		procNamesMap[procName] = true
		operationNames[procName] = ServerOperationTypeProc
	}
	for _, streamName := range streamNames {
		streamNamesMap[streamName] = true
		operationNames[streamName] = ServerOperationTypeStream
	}

	return &internalServer[T]{
		procNames:              procNames,
		procNamesMap:           procNamesMap,
		streamNames:            streamNames,
		streamNamesMap:         streamNamesMap,
		operationNamesMap:      operationNames,
		handlersMu:             sync.RWMutex{},
		procHandlers:           map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage) (any, error){},
		procInputHandlers:      map[string]func(ctx context.Context, ufoCtx T, input any) (any, error){},
		streamHandlers:         map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage, emit func(any) error) error{},
		streamInputHandlers:    map[string]func(ctx context.Context, ufoCtx T, input any) (any, error){},
		hooksBeforeHandler:     []ServerHookBeforeHandler[T]{},
		hooksBeforeProcRespond: []ServerHookBeforeProcRespond[T]{},
		hooksBeforeStreamEmit:  []ServerHookBeforeStreamEmit[T]{},
		hooksAfterProc:         []ServerHookAfterProc[T]{},
		hooksAfterStream:       []ServerHookAfterStream[T]{},
		hooksAfterStreamEmit:   []ServerHookAfterStreamEmit[T]{},
	}
}

// addHookBeforeHandler registers a hook function that executes before
// any handler (both procedures and streams) is invoked. Hook functions are
// executed in the order they were registered.
//
// This hook is useful for cross-cutting concerns like authentication,
// authorization, logging, or request validation that apply to all handlers.
//
// Parameters:
//   - hook: The hook function to register
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) addHookBeforeHandler(hook ServerHookBeforeHandler[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.hooksBeforeHandler = append(s.hooksBeforeHandler, hook)
	return s
}

// addHookBeforeProcRespond registers a hook function that executes before
// the procedure response is sent to the client. Hook functions are executed
// in the order they were registered.
//
// This hook is useful for response transformation, filtering, adding metadata,
// or modifying the response data before it reaches the client.
//
// Parameters:
//   - hook: The hook function to register
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) addHookBeforeProcRespond(hook ServerHookBeforeProcRespond[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.hooksBeforeProcRespond = append(s.hooksBeforeProcRespond, hook)
	return s
}

// addHookBeforeStreamEmit registers a hook function that executes before
// each stream event is emitted to the client. Hook functions are executed
// in the order they were registered.
//
// This hook is useful for response transformation, filtering, adding metadata,
// or logging stream events before they reach the client.
//
// Parameters:
//   - hook: The hook function to register
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) addHookBeforeStreamEmit(hook ServerHookBeforeStreamEmit[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.hooksBeforeStreamEmit = append(s.hooksBeforeStreamEmit, hook)
	return s
}

// addHookAfterProc registers a hook function that executes after
// procedure request processing, regardless of success or failure. Hook
// functions are executed in the order they were registered.
//
// This hook is useful for logging, metrics collection, or cleanup
// operations specific to procedure calls. This hook cannot modify
// the response as it has already been sent to the client.
//
// Parameters:
//   - hook: The hook function to register
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) addHookAfterProc(hook ServerHookAfterProc[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.hooksAfterProc = append(s.hooksAfterProc, hook)
	return s
}

// addHookAfterStream registers a hook function that executes after
// stream request processing ends, either successfully or with an error.
// Hook functions are executed in the order they were registered.
//
// This hook is useful for cleanup operations, logging, metrics collection,
// or resource deallocation specific to stream handling. This hook cannot
// modify anything as the stream has already completed.
//
// Parameters:
//   - hook: The hook function to register
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) addHookAfterStream(hook ServerHookAfterStream[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.hooksAfterStream = append(s.hooksAfterStream, hook)
	return s
}

// addHookAfterStreamEmit registers a hook function that executes after
// each stream event has been successfully emitted to the client. Hook
// functions are executed in the order they were registered.
//
// This hook is useful for logging, metrics collection, or post-emission
// cleanup operations for individual stream events. This hook cannot modify
// anything as the event has already been sent to the client.
//
// Parameters:
//   - hook: The hook function to register
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) addHookAfterStreamEmit(hook ServerHookAfterStreamEmit[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.hooksAfterStreamEmit = append(s.hooksAfterStreamEmit, hook)
	return s
}

// setProcHandler registers the implementation function for the specified procedure name.
// The handler function will be called when a client invokes the procedure via RPC.
//
// The handler receives the current context, UFO context, and raw JSON input,
// and must return the response data or an error. The input is provided as
// json.RawMessage to allow for flexible input processing.
//
// Parameters:
//   - procName: The name of the procedure to register
//   - handler: The function that implements the procedure logic
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) setProcHandler(
	procName string,
	handler func(
		ctx context.Context,
		ufoCtx T,
		input json.RawMessage,
	) (any, error),
) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.procHandlers[procName] = handler
	return s
}

// setProcInputHandler registers the input processing function for the specified
// procedure name. This handler is responsible for validating and transforming
// the input data before it reaches the main procedure handler.
//
// The input handler receives the current context, UFO context, and processed input,
// and must return the validated/transformed input or an error.
//
// Parameters:
//   - procName: The name of the procedure to register the input handler for
//   - inputHandler: The function that processes and validates procedure input
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) setProcInputHandler(
	procName string,
	inputHandler func(
		ctx context.Context,
		ufoCtx T,
		input any,
	) (any, error),
) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.procInputHandlers[procName] = inputHandler
	return s
}

// setStreamHandler registers the implementation function for the specified stream name.
// The handler function will be called when a client initiates a stream via RPC.
//
// The handler receives the current context, UFO context, raw JSON input, and an
// emit function for sending events to the client. The handler should call emit
// for each event and return when the stream is complete or an error occurs.
//
// Parameters:
//   - streamName: The name of the stream to register
//   - handler: The function that implements the stream logic
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) setStreamHandler(
	streamName string,
	handler func(
		ctx context.Context,
		ufoCtx T,
		input json.RawMessage,
		emit func(any) error,
	) error,
) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.streamHandlers[streamName] = handler
	return s
}

// setStreamInputHandler registers the input processing function for the specified
// stream name. This handler is responsible for validating and transforming
// the input data before the stream begins.
//
// The input handler receives the current context, UFO context, and processed input,
// and must return the validated/transformed input or an error.
//
// Parameters:
//   - streamName: The name of the stream to register the input handler for
//   - inputHandler: The function that processes and validates stream input
//
// Returns the server instance for method chaining.
func (s *internalServer[T]) setStreamInputHandler(
	streamName string,
	inputHandler func(
		ctx context.Context,
		ufoCtx T,
		input any,
	) (any, error),
) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.streamInputHandlers[streamName] = inputHandler
	return s
}

// handleRequest processes an incoming RPC request by parsing the request body,
// determining the request type (procedure or stream), and dispatching to the
// appropriate handler. This is the main entry point for all RPC requests.
//
// The request body can contain JSON with the input data for the handler.
//
// Parameters:
//   - ctx: The request context
//   - ufoCtx: The UFO context containing user-defined data
//   - operationName: The name of the procedure or stream to invoke
//   - httpAdapter: The HTTP adapter for reading requests and writing responses
//
// Returns an error if request processing fails at the transport level.
func (s *internalServer[T]) handleRequest(
	ctx context.Context,
	ufoCtx T,
	operationName string,
	httpAdapter ServerHTTPAdapter,
) error {
	if httpAdapter == nil {
		return fmt.Errorf("ServerRequestResponseProvider is nil, please provide a valid provider")
	}

	var input json.RawMessage
	if err := json.NewDecoder(httpAdapter.RequestBody()).Decode(&input); err != nil {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "Invalid request body"},
		}
		return s.writeProcResponse(httpAdapter, res)
	}

	operationType, operationExists := s.operationNamesMap[operationName]
	isStream := operationType == ServerOperationTypeStream
	if !operationExists {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "Invalid operation name"},
		}
		return s.writeProcResponse(httpAdapter, res)
	}

	// Lock the handlers map for reading
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	if isStream {
		return s.handleStreamRequest(ctx, ufoCtx, operationName, input, httpAdapter)
	}

	return s.handleProcRequest(ctx, ufoCtx, operationName, input, httpAdapter)
}

// writeProcResponse writes a procedure response to the client as JSON.
// This helper method sets the appropriate Content-Type header and marshals
// the response data before sending it to the client.
//
// Parameters:
//   - httpAdapter: The HTTP adapter for writing the response
//   - response: The response data to send to the client
//
// Returns an error if writing the response fails.
func (s *internalServer[T]) writeProcResponse(
	httpAdapter ServerHTTPAdapter,
	response Response[any],
) error {
	httpAdapter.SetHeader("Content-Type", "application/json")
	_, err := httpAdapter.Write(response.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// handleProcRequest processes a procedure RPC request through the complete
// hook chain and handler execution. This includes validation, hook
// execution, handler invocation, and response generation.
//
// The processing flow is:
// 1. Validate procedure name and implementation
// 2. Execute before-handler hooks
// 3. Invoke the procedure handler
// 4. Execute before-response hooks (can modify response)
// 5. Send the response to the client
// 6. Execute after-procedure hooks (for cleanup/logging)
//
// Parameters:
//   - ctx: The request context
//   - ufoCtx: The UFO context containing user-defined data
//   - procName: The name of the procedure to invoke
//   - rawInput: The raw JSON input for the procedure
//   - httpAdapter: The HTTP adapter for writing the response
//
// Returns an error if the response cannot be written to the client.
func (s *internalServer[T]) handleProcRequest(
	ctx context.Context,
	ufoCtx T,
	procName string,
	rawInput json.RawMessage,
	httpAdapter ServerHTTPAdapter,
) error {
	response := Response[any]{Ok: true}

	// Validate procedure name
	if !slices.Contains(s.procNames, procName) {
		response = Response[any]{
			Ok: false,
			Error: Error{
				Message: procName + " procedure not found",
				Details: map[string]any{"procedure": procName},
			},
		}
	}

	// Validate procedure implementation
	proc, found := s.procHandlers[procName]
	if response.Ok && !found {
		response = Response[any]{
			Ok: false,
			Error: Error{
				Message: procName + " procedure not implemented",
				Details: map[string]any{"procedure": procName},
			},
		}
	}

	// Execute Before middlewares if we haven't encountered an error yet
	if response.Ok {
		for _, hook := range s.hooksBeforeHandler {
			var err error
			if ctx, ufoCtx, err = hook(ctx, ufoCtx, ServerOperationTypeProc, procName); err != nil {
				response = Response[any]{
					Ok:    false,
					Error: asError(err),
				}
				break
			}
		}
	}

	// Run handler if no errors have occurred
	if response.Ok {
		if output, err := proc(ctx, ufoCtx, rawInput); err != nil {
			response = Response[any]{
				Ok:    false,
				Error: asError(err),
			}
		} else {
			response = Response[any]{
				Ok:     true,
				Output: output,
			}
		}
	}

	// Execute before proc respond hooks
	for _, hook := range s.hooksBeforeProcRespond {
		response = hook(ctx, ufoCtx, procName, response)
	}

	// Write the response to the client
	httpAdapter.SetHeader("Content-Type", "application/json")
	_, err := httpAdapter.Write(response.Bytes())
	if err != nil {
		return err
	}

	// Execute after proc hooks
	for _, hook := range s.hooksAfterProc {
		hook(ctx, ufoCtx, procName, response)
	}

	return nil
}

// handleStreamRequest processes a stream RPC request by setting up Server-Sent Events,
// executing the hook chain, and managing the stream lifecycle. This includes
// validation, hook execution, stream handler invocation, and event emission.
//
// The processing flow is:
// 1. Set up SSE headers and emit functions
// 2. Validate stream name and implementation
// 3. Execute before-handler hooks
// 4. Invoke the stream handler with emit capability
// 5. Execute after-stream hooks
//
// Stream handlers receive an emit function that allows them to send events to the client
// in real-time. Each emitted event goes through before-emit and after-emit hooks.
//
// Parameters:
//   - ctx: The request context
//   - ufoCtx: The UFO context containing user-defined data
//   - streamName: The name of the stream to invoke
//   - rawInput: The raw JSON input for the stream
//   - httpAdapter: The HTTP adapter for writing SSE responses
//
// Returns an error if the stream setup or emission fails.
func (s *internalServer[T]) handleStreamRequest(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	rawInput json.RawMessage,
	httpAdapter ServerHTTPAdapter,
) error {
	// emit sends a response event to the client.
	// If the client disconnects, it returns an error.
	emit := func(data Response[any]) error {
		// execute before emit middlewares
		for _, hook := range s.hooksBeforeStreamEmit {
			data = hook(ctx, ufoCtx, streamName, data)
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal stream data: %w", err)
		}

		resPayload := fmt.Sprintf("data: %s\n\n", jsonData)
		_, err = httpAdapter.Write([]byte(resPayload))
		if err != nil {
			return err
		}

		if err := httpAdapter.Flush(); err != nil {
			return err
		}

		// execute after emit middlewares
		for _, hook := range s.hooksAfterStreamEmit {
			hook(ctx, ufoCtx, streamName, data)
		}

		return nil
	}

	emitSuccess := func(data any) error {
		response := Response[any]{
			Ok:     true,
			Output: data,
		}
		return emit(response)
	}

	emitError := func(err error) error {
		response := Response[any]{
			Ok:    false,
			Error: asError(err),
		}
		return emit(response)
	}

	// Upgrade the request to a stream
	httpAdapter.SetHeader("Content-Type", "text/event-stream")
	httpAdapter.SetHeader("Cache-Control", "no-cache")
	httpAdapter.SetHeader("Connection", "keep-alive")

	// Validate stream name
	if !slices.Contains(s.streamNames, streamName) {
		return emitError(Error{
			Message: streamName + " stream not found",
			Details: map[string]any{"stream": streamName},
		})
	}

	// Validate stream implementation
	stream, found := s.streamHandlers[streamName]
	if !found {
		return emitError(Error{
			Message: streamName + " stream not implemented",
			Details: map[string]any{"stream": streamName},
		})
	}

	// Execute Before middlewares
	for _, hook := range s.hooksBeforeHandler {
		var err error
		if ctx, ufoCtx, err = hook(ctx, ufoCtx, ServerOperationTypeStream, streamName); err != nil {
			return emitError(err)
		}
	}

	// Run the stream handler and wait for it to finish
	err := stream(ctx, ufoCtx, rawInput, emitSuccess)
	if err != nil {
		// execute after stream middlewares
		for _, hook := range s.hooksAfterStream {
			hook(ctx, ufoCtx, streamName, err)
		}

		return emitError(err)
	}

	// execute after stream middlewares
	for _, hook := range s.hooksAfterStream {
		hook(ctx, ufoCtx, streamName, nil)
	}

	return nil
}
