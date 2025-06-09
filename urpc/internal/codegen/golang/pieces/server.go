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

// ServerRequestResponseProvider provides the required methods for UFO RPC server
// to handle a request and write a response to the client.
type ServerRequestResponseProvider interface {
	// RequestGetBodyReader returns the body reader for the request.
	RequestGetBodyReader() io.Reader
	// ResponseSetHeader sets a header in the response.
	ResponseSetHeader(key, value string)
	// ResponseWrite writes data to the response.
	ResponseWrite(data []byte) (int, error)
	// ResponseFlush flushes the response to the client.
	ResponseFlush() error
}

// ServerNetHTTPRequestResponseProvider implements the ServerRequestResponseProvider interface for net/http.
type ServerNetHTTPRequestResponseProvider struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

// NewServerNetHTTPRequestResponseProvider creates a new ServerNetHTTPRequestResponseProvider.
func NewServerNetHTTPRequestResponseProvider[T any](initialUFOCtx T, w http.ResponseWriter, r *http.Request) ServerRequestResponseProvider {
	return &ServerNetHTTPRequestResponseProvider{
		responseWriter: w,
		request:        r,
	}
}

func (r *ServerNetHTTPRequestResponseProvider) RequestGetBodyReader() io.Reader {
	return r.request.Body
}

func (r *ServerNetHTTPRequestResponseProvider) ResponseSetHeader(key, value string) {
	r.responseWriter.Header().Set(key, value)
}

func (r *ServerNetHTTPRequestResponseProvider) ResponseWrite(data []byte) (int, error) {
	return r.responseWriter.Write(data)
}

func (r *ServerNetHTTPRequestResponseProvider) ResponseFlush() error {
	if f, ok := r.responseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// MiddlewareBeforeHandler runs before procedure and stream handler is executed.
type MiddlewareBeforeHandler[T any] func(
	ctx context.Context,
	ufoCtx T,
	handlerType string,
	handlerName string,
) (context.Context, T, error)

// MiddlewareBeforeStreamEmit runs before stream event is emitted.
type MiddlewareBeforeStreamEmit[T any] func(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	response Response[any],
) Response[any]

// MiddlewareAfterProc runs after procedure request processing.
type MiddlewareAfterProc[T any] func(
	ctx context.Context,
	ufoCtx T,
	procName string,
	response Response[any],
) (context.Context, T, Response[any])

// MiddlewareAfterStream runs after stream request processing ends.
type MiddlewareAfterStream[T any] func(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	err error,
)

// MiddlewareAfterStreamEmit runs after stream event is emitted.
type MiddlewareAfterStreamEmit[T any] func(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	response Response[any],
)

// -----------------------------------------------------------------------------
// Server Internal Implementation
// -----------------------------------------------------------------------------

// internalServer handles RPC requests.
type internalServer[T any] struct {
	procNames                   []string
	streamNames                 []string
	handlersMu                  sync.RWMutex
	procHandlers                map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage) (any, error)
	procInputProcessors         map[string]func(ctx context.Context, ufoCtx T, input any) (any, error)
	streamHandlers              map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage, emit func(any) error) error
	streamInputProcessors       map[string]func(ctx context.Context, ufoCtx T, input any) (any, error)
	beforeHandlerMiddlewares    []MiddlewareBeforeHandler[T]
	beforeStreamEmitMiddlewares []MiddlewareBeforeStreamEmit[T]
	afterProcMiddlewares        []MiddlewareAfterProc[T]
	afterStreamMiddlewares      []MiddlewareAfterStream[T]
	afterStreamEmitMiddlewares  []MiddlewareAfterStreamEmit[T]
}

// newInternalServer creates a new UFO RPC server
//
// The generic type T represents the context type, used to pass additional data
// to procedures, such as authentication information, user session or any
// other data you want to pass to procedures before they are executed.
func newInternalServer[T any](
	procNames []string,
	streamNames []string,
) *internalServer[T] {
	return &internalServer[T]{
		procNames:                   procNames,
		streamNames:                 streamNames,
		handlersMu:                  sync.RWMutex{},
		procHandlers:                map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage) (any, error){},
		procInputProcessors:         map[string]func(ctx context.Context, ufoCtx T, input any) (any, error){},
		streamHandlers:              map[string]func(ctx context.Context, ufoCtx T, input json.RawMessage, emit func(any) error) error{},
		streamInputProcessors:       map[string]func(ctx context.Context, ufoCtx T, input any) (any, error){},
		beforeHandlerMiddlewares:    []MiddlewareBeforeHandler[T]{},
		beforeStreamEmitMiddlewares: []MiddlewareBeforeStreamEmit[T]{},
		afterProcMiddlewares:        []MiddlewareAfterProc[T]{},
		afterStreamMiddlewares:      []MiddlewareAfterStream[T]{},
		afterStreamEmitMiddlewares:  []MiddlewareAfterStreamEmit[T]{},
	}
}

// addMiddlewareBeforeHandler adds a middleware function that runs before the handler. Both for procedures and streams.
func (s *internalServer[T]) addMiddlewareBeforeHandler(fn MiddlewareBeforeHandler[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.beforeHandlerMiddlewares = append(s.beforeHandlerMiddlewares, fn)
	return s
}

// addMiddlewareBeforeStreamEmit adds a middleware function that runs before the stream event is emitted.
func (s *internalServer[T]) addMiddlewareBeforeStreamEmit(fn MiddlewareBeforeStreamEmit[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.beforeStreamEmitMiddlewares = append(s.beforeStreamEmitMiddlewares, fn)
	return s
}

// addMiddlewareAfterProc adds a middleware function that runs after the handler. Only for procedures, not for streams.
func (s *internalServer[T]) addMiddlewareAfterProc(fn MiddlewareAfterProc[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.afterProcMiddlewares = append(s.afterProcMiddlewares, fn)
	return s
}

// addMiddlewareAfterStream adds a middleware function that runs after the stream handler.
func (s *internalServer[T]) addMiddlewareAfterStream(fn MiddlewareAfterStream[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.afterStreamMiddlewares = append(s.afterStreamMiddlewares, fn)
	return s
}

// addMiddlewareAfterStreamEmit adds a middleware function that runs after the stream event is emitted.
func (s *internalServer[T]) addMiddlewareAfterStreamEmit(fn MiddlewareAfterStreamEmit[T]) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.afterStreamEmitMiddlewares = append(s.afterStreamEmitMiddlewares, fn)
	return s
}

// setProcHandler registers the handler for the provided procedure name
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

// setProcInputProcessor registers the input processor for the provided procedure name
func (s *internalServer[T]) setProcInputProcessor(
	procName string,
	processor func(
		ctx context.Context,
		ufoCtx T,
		input any,
	) (any, error),
) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.procInputProcessors[procName] = processor
	return s
}

// setStreamHandler registers the handler for the provided stream name
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

// setStreamInputProcessor registers the input processor for the provided stream name
func (s *internalServer[T]) setStreamInputProcessor(
	streamName string,
	processor func(
		ctx context.Context,
		ufoCtx T,
		input any,
	) (any, error),
) *internalServer[T] {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.streamInputProcessors[streamName] = processor
	return s
}

// handleRequest processes an incoming RPC request
func (s *internalServer[T]) handleRequest(
	ctx context.Context,
	ufoCtx T,
	reqResProvider ServerRequestResponseProvider,
) error {
	if reqResProvider == nil {
		return fmt.Errorf("ServerRequestResponseProvider is nil, please provide a valid provider")
	}

	var jsonBody struct {
		Type  string          `json:"type"` // "proc" or "stream"
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}
	if err := json.NewDecoder(reqResProvider.RequestGetBodyReader()).Decode(&jsonBody); err != nil {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "Invalid request body"},
		}
		return s.writeProcResponse(reqResProvider, res)
	}

	isProc := jsonBody.Type == "proc"
	isStream := jsonBody.Type == "stream"
	if !isProc && !isStream {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "Invalid request body, type must be 'proc' or 'stream'"},
		}
		return s.writeProcResponse(reqResProvider, res)
	}

	// Lock the handlers map for reading
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	if isStream {
		return s.handleStreamRequest(ctx, ufoCtx, jsonBody.Name, jsonBody.Input, reqResProvider)
	}

	return s.handleProcRequest(ctx, ufoCtx, jsonBody.Name, jsonBody.Input, reqResProvider)
}

// writeProcResponse writes a procedure response to the client
func (s *internalServer[T]) writeProcResponse(
	reqResProvider ServerRequestResponseProvider,
	response Response[any],
) error {
	reqResProvider.ResponseSetHeader("Content-Type", "application/json")
	_, err := reqResProvider.ResponseWrite(response.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// handleProcRequest processes a procedure request
func (s *internalServer[T]) handleProcRequest(
	ctx context.Context,
	ufoCtx T,
	procName string,
	rawInput json.RawMessage,
	reqResProvider ServerRequestResponseProvider,
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
		for _, fn := range s.beforeHandlerMiddlewares {
			var err error
			if ctx, ufoCtx, err = fn(ctx, ufoCtx, "proc", procName); err != nil {
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

	// Always execute After middlewares, regardless of any previous errors
	for _, fn := range s.afterProcMiddlewares {
		ctx, ufoCtx, response = fn(ctx, ufoCtx, procName, response)
	}

	// Write the response to the client
	reqResProvider.ResponseSetHeader("Content-Type", "application/json")
	_, err := reqResProvider.ResponseWrite(response.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// handleStreamRequest processes a stream request
func (s *internalServer[T]) handleStreamRequest(
	ctx context.Context,
	ufoCtx T,
	streamName string,
	rawInput json.RawMessage,
	reqResProvider ServerRequestResponseProvider,
) error {
	// emit sends a response event to the client.
	// If the client disconnects, it returns an error.
	emit := func(data Response[any]) error {
		// execute before emit middlewares
		for _, fn := range s.beforeStreamEmitMiddlewares {
			data = fn(ctx, ufoCtx, streamName, data)
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal stream data: %w", err)
		}

		resPayload := fmt.Sprintf("data: %s\n\n", jsonData)
		_, err = reqResProvider.ResponseWrite([]byte(resPayload))
		if err != nil {
			return err
		}

		if err := reqResProvider.ResponseFlush(); err != nil {
			return err
		}

		// execute after emit middlewares
		for _, fn := range s.afterStreamEmitMiddlewares {
			fn(ctx, ufoCtx, streamName, data)
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
	reqResProvider.ResponseSetHeader("Content-Type", "text/event-stream")
	reqResProvider.ResponseSetHeader("Cache-Control", "no-cache")
	reqResProvider.ResponseSetHeader("Connection", "keep-alive")

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
	for _, fn := range s.beforeHandlerMiddlewares {
		var err error
		if ctx, ufoCtx, err = fn(ctx, ufoCtx, "stream", streamName); err != nil {
			return emitError(err)
		}
	}

	// Run the stream handler and wait for it to finish
	err := stream(ctx, ufoCtx, rawInput, emitSuccess)
	if err != nil {
		// execute after stream middlewares
		for _, fn := range s.afterStreamMiddlewares {
			fn(ctx, ufoCtx, streamName, err)
		}

		return emitError(err)
	}

	// execute after stream middlewares
	for _, fn := range s.afterStreamMiddlewares {
		fn(ctx, ufoCtx, streamName, nil)
	}

	return nil
}
