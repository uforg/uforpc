//nolint:unused
package pieces

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// -----------------------------------------------------------------------------
// Middleware-based Server Architecture
// -----------------------------------------------------------------------------

// HandlerContext is the unified container for all request information and state
// that flows through the entire request processing pipeline.
//
// The generic type P represents the user-defined container for application
// dependencies and request data (e.g., UserID, DB connection, etc.).
//
// The generic type I represents the input type, which can be any type depending
// on the operation.
type HandlerContext[P any, I any] struct {
	// Props is the user-defined container, created per request,
	// for application dependencies and request data (e.g., UserID).
	Props P

	// Input contains the request body, already deserialized and typed.
	// For global middlewares, the type I will be any.
	Input I

	// Context is the standard Go context.Context for cancellations and deadlines.
	Context context.Context

	// operationName is the name of the invoked proc or stream (e.g., "CreateUser").
	operationName string

	// operationType is the type of operation ("proc" or "stream").
	operationType string
}

// OperationName returns the name of the operation (e.g. "CreateUser", "GetPost", etc.)
func (h *HandlerContext[P, I]) OperationName() string { return h.operationName }

// OperationType returns the type of operation (e.g. "proc" or "stream")
func (h *HandlerContext[P, I]) OperationType() string { return h.operationType }

// GlobalHandlerFunc is the signature for a global handler function.
// Both for procedures and streams
type GlobalHandlerFunc[P any] func(
	c *HandlerContext[P, any],
) (any, error)

// GlobalMiddleware is the signature for a middleware applied to all requests.
type GlobalMiddleware[P any] func(
	next GlobalHandlerFunc[P],
) GlobalHandlerFunc[P]

// ProcHandlerFunc is the signature of the final business handler for a proc.
type ProcHandlerFunc[P any, I any, O any] func(
	c *HandlerContext[P, I],
) (O, error)

// ProcMiddlewareFunc is the signature for a proc-specific typed middleware.
// It uses a wrapper pattern for a clean composition.
//
// This is the same as [GlobalMiddleware] but for specific procedures and with types.
type ProcMiddlewareFunc[P any, I any, O any] func(
	next ProcHandlerFunc[P, I, O],
) ProcHandlerFunc[P, I, O]

// StreamHandlerFunc is the signature of the main handler that initializes a stream.
type StreamHandlerFunc[P any, I any, O any] func(
	c *HandlerContext[P, I],
	emit EmitFunc[P, I, O],
) error

// StreamMiddlewareFunc is the signature for a middleware that wraps the main stream handler.
type StreamMiddlewareFunc[P any, I any, O any] func(
	next StreamHandlerFunc[P, I, O],
) StreamHandlerFunc[P, I, O]

// EmitFunc is the signature for emitting events from a stream.
type EmitFunc[P any, I any, O any] func(
	c *HandlerContext[P, I],
	output O,
) error

// EmitMiddlewareFunc is the signature for a middleware that wraps each call to emit.
type EmitMiddlewareFunc[P any, I any, O any] func(
	next EmitFunc[P, I, O],
) EmitFunc[P, I, O]

// -----------------------------------------------------------------------------
// Server Internal Implementation
// -----------------------------------------------------------------------------

// internalServer manages RPC request handling and middleware execution for
// both procedures and streams. It maintains handler registrations, middleware
// chains, and coordinates the complete request lifecycle.
//
// The generic type P represents the user context type, allowing users to pass
// custom data (authentication info, user sessions, etc.) through the entire
// request processing pipeline.
type internalServer[P any] struct {
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
	// handlersMu protects all handler maps and middleware slices from concurrent access
	handlersMu sync.RWMutex
	// procHandlers maps procedure names to their final implementation functions
	// The function receives a generic HandlerContext with any input and returns (any, error)
	procHandlers map[string]ProcHandlerFunc[P, any, any]
	// streamHandlers maps stream names to their final implementation functions
	// The function receives a generic HandlerContext with any input and a generic emit function.
	streamHandlers map[string]StreamHandlerFunc[P, any, any]
	// globalMiddlewares contains middlewares that run for every request (both procs and streams)
	globalMiddlewares []GlobalMiddleware[P]
	// procMiddlewares contains per-procedure middlewares (already adapted to untyped layer)
	procMiddlewares map[string][]ProcMiddlewareFunc[P, any, any]
	// streamMiddlewares contains per-stream middlewares (already adapted to untyped layer)
	streamMiddlewares map[string][]StreamMiddlewareFunc[P, any, any]
	// streamEmitMiddlewares contains per-stream emit middlewares (already adapted to untyped layer)
	streamEmitMiddlewares map[string][]EmitMiddlewareFunc[P, any, any]
}

// newInternalServer creates a new UFO RPC server instance with the specified
// procedure and stream names. The server is initialized with empty handler
// maps and middleware slices, ready for registration.
//
// The generic type T represents the user context type, used to pass additional
// data to handlers, such as authentication information, user sessions, or any
// other request-scoped data.
//
// Parameters:
//   - procNames: List of procedure names that this server will handle
//   - streamNames: List of stream names that this server will handle
//
// Returns a new internalServer instance ready for handler and middleware registration.
func newInternalServer[P any](
	procNames []string,
	streamNames []string,
) *internalServer[P] {
	procNamesMap := make(map[string]bool)
	streamNamesMap := make(map[string]bool)
	operationNamesMap := make(map[string]string)
	for _, procName := range procNames {
		procNamesMap[procName] = true
		operationNamesMap[procName] = ServerOperationTypeProc
	}
	for _, streamName := range streamNames {
		streamNamesMap[streamName] = true
		operationNamesMap[streamName] = ServerOperationTypeStream
	}

	return &internalServer[P]{
		procNames:             procNames,
		procNamesMap:          procNamesMap,
		streamNames:           streamNames,
		streamNamesMap:        streamNamesMap,
		operationNamesMap:     operationNamesMap,
		handlersMu:            sync.RWMutex{},
		procHandlers:          map[string]ProcHandlerFunc[P, any, any]{},
		streamHandlers:        map[string]StreamHandlerFunc[P, any, any]{},
		globalMiddlewares:     []GlobalMiddleware[P]{},
		procMiddlewares:       map[string][]ProcMiddlewareFunc[P, any, any]{},
		streamMiddlewares:     map[string][]StreamMiddlewareFunc[P, any, any]{},
		streamEmitMiddlewares: map[string][]EmitMiddlewareFunc[P, any, any]{},
	}
}

// addGlobalMiddleware registers a global middleware that executes for every request (proc and stream).
// Middlewares are executed in the order they were registered.
func (s *internalServer[P]) addGlobalMiddleware(
	mw GlobalMiddleware[P],
) *internalServer[P] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.globalMiddlewares = append(s.globalMiddlewares, mw)
	return s
}

// addProcMiddleware registers an untyped wrapper middleware for a specific procedure.
// Middlewares are executed in the order they were registered.
func (s *internalServer[P]) addProcMiddleware(
	procName string,
	mw ProcMiddlewareFunc[P, any, any],
) *internalServer[P] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.procMiddlewares[procName] = append(s.procMiddlewares[procName], mw)
	return s
}

// addStreamMiddleware registers an untyped wrapper middleware for a specific stream.
// Middlewares are executed in the order they were registered.
func (s *internalServer[P]) addStreamMiddleware(
	streamName string,
	mw StreamMiddlewareFunc[P, any, any],
) *internalServer[P] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.streamMiddlewares[streamName] = append(s.streamMiddlewares[streamName], mw)
	return s
}

// addStreamEmitMiddleware registers an untyped emit wrapper middleware for a specific stream.
// Middlewares are executed in the order they were registered.
func (s *internalServer[P]) addStreamEmitMiddleware(
	streamName string,
	mw EmitMiddlewareFunc[P, any, any],
) *internalServer[P] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.streamEmitMiddlewares[streamName] = append(s.streamEmitMiddlewares[streamName], mw)
	return s
}

// setProcHandler registers the final implementation function for the specified procedure name.
// The provided function is stored as-is. Middlewares are composed at request time.
//
// Panics if a handler is already registered for the given procedure name.
func (s *internalServer[P]) setProcHandler(
	procName string,
	handler ProcHandlerFunc[P, any, any],
) *internalServer[P] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	if _, exists := s.procHandlers[procName]; exists {
		panic(fmt.Sprintf("the procedure handler for %s is already registered", procName))
	}
	s.procHandlers[procName] = handler
	return s
}

// setStreamHandler registers the final implementation function for the specified stream name.
// The provided function is stored as-is. Middlewares are composed at request time.
//
// Panics if a handler is already registered for the given stream name.
func (s *internalServer[P]) setStreamHandler(
	streamName string,
	handler StreamHandlerFunc[P, any, any],
) *internalServer[P] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	if _, exists := s.streamHandlers[streamName]; exists {
		panic(fmt.Sprintf("the stream handler for %s is already registered", streamName))
	}
	s.streamHandlers[streamName] = handler
	return s
}

// handleRequest processes an incoming RPC request by parsing the request body,
// building the global middleware chain, and dispatching to the appropriate
// adapter (procedure or stream).
//
// The request body must contain a JSON object with the input data for the handler.
//
// Parameters:
//   - ctx: The request context
//   - props: The UFO context containing user-defined data
//   - operationName: The name of the procedure or stream to invoke
//   - httpAdapter: The HTTP adapter for reading requests and writing responses
//
// Returns an error if request processing fails at the transport level.
func (s *internalServer[P]) handleRequest(
	ctx context.Context,
	props P,
	operationName string,
	httpAdapter ServerHTTPAdapter,
) error {
	if httpAdapter == nil {
		return fmt.Errorf("the HTTP adapter is nil, please provide a valid adapter")
	}

	// Decode the request body into a json.RawMessage as the initial input container
	var input json.RawMessage
	if err := json.NewDecoder(httpAdapter.RequestBody()).Decode(&input); err != nil {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "Invalid request body"},
		}
		return s.writeProcResponse(httpAdapter, res)
	}

	operationType, operationExists := s.operationNamesMap[operationName]
	if !operationExists {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "Invalid operation name"},
		}
		return s.writeProcResponse(httpAdapter, res)
	}

	// Build the unified handler context for the global middleware chain
	c := &HandlerContext[P, any]{
		Input:         input,
		Props:         props,
		Context:       ctx,
		operationName: operationName,
		operationType: operationType,
	}

	// Track whether the stream connection has been started (headers sent)
	startedStream := false

	// Adapter bridges the global chain with the specific proc/stream function
	adapter := func(c *HandlerContext[P, any]) (any, error) {
		switch operationType {
		case ServerOperationTypeProc:
			return s.handleProcRequest(c, operationName)
		case ServerOperationTypeStream:
			startedStream = true
			return nil, s.handleStreamRequest(c, operationName, httpAdapter)
		default:
			return nil, fmt.Errorf("unsupported operation type: %s", operationType)
		}
	}

	// Build the global middleware chain (in reverse registration order)
	exec := adapter
	if len(s.globalMiddlewares) > 0 {
		mwChain := append([]GlobalMiddleware[P](nil), s.globalMiddlewares...)
		for i := len(mwChain) - 1; i >= 0; i-- {
			exec = mwChain[i](exec)
		}
	}

	// Execute the chain
	output, err := exec(c)

	// Stream response path
	if operationType == ServerOperationTypeStream {
		if err != nil {
			if startedStream {
				// Emit a final error event for stream failures after connection started
				response := Response[any]{
					Ok:    false,
					Error: asError(err),
				}
				jsonData, marshalErr := json.Marshal(response)
				if marshalErr != nil {
					return fmt.Errorf("failed to marshal stream error: %w", marshalErr)
				}
				resPayload := fmt.Sprintf("data: %s\n\n", jsonData)
				if _, writeErr := httpAdapter.Write([]byte(resPayload)); writeErr != nil {
					return writeErr
				}
				if flushErr := httpAdapter.Flush(); flushErr != nil {
					return flushErr
				}
			} else {
				// Before establishing the stream, return a single JSON error response
				res := Response[any]{
					Ok:    false,
					Error: asError(err),
				}
				return s.writeProcResponse(httpAdapter, res)
			}
		}
		return nil
	}

	// Procedure response path
	response := Response[any]{}
	if err != nil {
		response.Ok = false
		response.Error = asError(err)
	} else {
		response.Ok = true
		response.Output = output
	}

	return s.writeProcResponse(httpAdapter, response)
}

// handleProcRequest builds the per-request middleware chain for a procedure and executes it.
// It returns the procedure output (as any) and an error if the handler failed.
func (s *internalServer[P]) handleProcRequest(
	c *HandlerContext[P, any],
	procName string,
) (any, error) {
	// Snapshot handler and middlewares under read lock
	s.handlersMu.RLock()
	baseHandler, ok := s.procHandlers[procName]
	mws := s.procMiddlewares[procName]
	s.handlersMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%s procedure not implemented", procName)
	}

	// Compose middlewares around the base handler (reverse registration order)
	final := baseHandler
	if len(mws) > 0 {
		mwChain := append([]ProcMiddlewareFunc[P, any, any](nil), mws...)
		for i := len(mwChain) - 1; i >= 0; i-- {
			final = mwChain[i](final)
		}
	}

	return final(c)
}

// handleStreamRequest builds the per-request middleware chain for a stream, sets up SSE,
// composes emit middlewares, and executes the stream handler.
func (s *internalServer[P]) handleStreamRequest(
	c *HandlerContext[P, any],
	streamName string,
	httpAdapter ServerHTTPAdapter,
) error {
	// Set SSE headers
	httpAdapter.SetHeader("Content-Type", "text/event-stream")
	httpAdapter.SetHeader("Cache-Control", "no-cache")
	httpAdapter.SetHeader("Connection", "keep-alive")

	// Snapshot handler and middlewares under read lock
	s.handlersMu.RLock()
	baseHandler, ok := s.streamHandlers[streamName]
	streamMws := s.streamMiddlewares[streamName]
	emitMws := s.streamEmitMiddlewares[streamName]
	s.handlersMu.RUnlock()

	if !ok {
		return fmt.Errorf("%s stream not implemented", streamName)
	}

	// Base emit writes SSE envelope with {ok:true, output}
	baseEmit := func(_ *HandlerContext[P, any], data any) error {
		response := Response[any]{
			Ok:     true,
			Output: data,
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			return fmt.Errorf("failed to marshal stream data: %w", err)
		}
		resPayload := fmt.Sprintf("data: %s\n\n", jsonData)
		if _, err = httpAdapter.Write([]byte(resPayload)); err != nil {
			return err
		}
		if err := httpAdapter.Flush(); err != nil {
			return err
		}
		return nil
	}

	// Compose emit middlewares (reverse registration order)
	emitFinal := baseEmit
	if len(emitMws) > 0 {
		mwChain := append([]EmitMiddlewareFunc[P, any, any](nil), emitMws...)
		for i := len(mwChain) - 1; i >= 0; i-- {
			emitFinal = func(c *HandlerContext[P, any], data any) error {
				return mwChain[i](emitFinal)(c, data)
			}
		}
	}

	// Compose stream middlewares around the base handler (reverse order)
	final := baseHandler
	if len(streamMws) > 0 {
		mwChain := append([]StreamMiddlewareFunc[P, any, any](nil), streamMws...)
		for i := len(mwChain) - 1; i >= 0; i-- {
			final = mwChain[i](final)
		}
	}

	// Execute chain
	return final(c, emitFinal)
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
func (s *internalServer[P]) writeProcResponse(
	httpAdapter ServerHTTPAdapter,
	response Response[any],
) error {
	httpAdapter.SetHeader("Content-Type", "application/json")
	_, err := httpAdapter.Write(response.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}
