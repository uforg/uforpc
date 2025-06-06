//nolint:unused
package pieces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
)

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Server Types
// -----------------------------------------------------------------------------

// ServerRequestResponseProvider provides the required methods for UFO RPC server
// to handle a request and write a response to the client.
type ServerRequestResponseProvider[T any] interface {
	// RequestGetInitialContext returns the initial context for the request.
	RequestGetInitialContext() T
	// RequestGetBodyReader returns the body reader for the request.
	RequestGetBodyReader() io.Reader
	// ResponseSetHeader sets a header in the response.
	ResponseSetHeader(key, value string)
	// ResponseWrite writes data to the response.
	ResponseWrite(data []byte) (int, error)
	// ResponseFlush flushes the response to the client.
	ResponseFlush()
}

// ServerNetHTTPRequestResponseProvider implements the ServerRequestResponseProvider interface for net/http.
type ServerNetHTTPRequestResponseProvider[T any] struct {
	initialContext T
	responseWriter http.ResponseWriter
	request        *http.Request
}

// NewServerNetHTTPRequestResponseProvider creates a new ServerNetHTTPRequestResponseProvider.
func NewServerNetHTTPRequestResponseProvider[T any](initialContext T, w http.ResponseWriter, r *http.Request) ServerRequestResponseProvider[T] {
	return &ServerNetHTTPRequestResponseProvider[T]{
		initialContext: initialContext,
		responseWriter: w,
		request:        r,
	}
}

func (r *ServerNetHTTPRequestResponseProvider[T]) RequestGetInitialContext() T {
	return r.initialContext
}

func (r *ServerNetHTTPRequestResponseProvider[T]) RequestGetBodyReader() io.Reader {
	return r.request.Body
}

func (r *ServerNetHTTPRequestResponseProvider[T]) ResponseSetHeader(key, value string) {
	r.responseWriter.Header().Set(key, value)
}

func (r *ServerNetHTTPRequestResponseProvider[T]) ResponseWrite(data []byte) (int, error) {
	return r.responseWriter.Write(data)
}

func (r *ServerNetHTTPRequestResponseProvider[T]) ResponseFlush() {
	if f, ok := r.responseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// MiddlewareBefore runs before request processing. Both for procedures and streams.
type MiddlewareBefore[T any] func(requestType string, requestName string, context T) (T, error)

// MiddlewareAfter runs after request processing. Only supports procedures, not streams.
type MiddlewareAfter[T any] func(requestType string, requestName string, context T, response Response[any]) Response[any]

// -----------------------------------------------------------------------------
// Server Internal Implementation
// -----------------------------------------------------------------------------

// internalServer handles RPC requests.
type internalServer[T any] struct {
	procNames             []string
	streamNames           []string
	procHandlers          map[string]func(context T, input json.RawMessage) (any, error)
	procInputProcessors   map[string]func(context T, input any) (any, error)
	streamHandlers        map[string]func(context T, input json.RawMessage, emit func(any) error) error
	streamInputProcessors map[string]func(context T, input any) (any, error)
	beforeMiddlewares     []MiddlewareBefore[T]
	afterMiddlewares      []MiddlewareAfter[T]
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
		procNames:             procNames,
		streamNames:           streamNames,
		procHandlers:          map[string]func(T, json.RawMessage) (any, error){},
		procInputProcessors:   map[string]func(T, any) (any, error){},
		streamHandlers:        map[string]func(T, json.RawMessage, func(any) error) error{},
		streamInputProcessors: map[string]func(T, any) (any, error){},
		beforeMiddlewares:     []MiddlewareBefore[T]{},
		afterMiddlewares:      []MiddlewareAfter[T]{},
	}
}

// addMiddlewareBefore adds a middleware function that runs before the handler.
//
// It modifies the request context before it reaches the main procedure.
//
// Multiple MiddlewareBefore can be added and are processed in order.
func (s *internalServer[T]) addMiddlewareBefore(fn MiddlewareBefore[T]) *internalServer[T] {
	s.beforeMiddlewares = append(s.beforeMiddlewares, fn)
	return s
}

// addMiddlewareAfter adds a middleware function that runs after the handler.
//
// It modifies the response before it is sent back to the client.
//
// Multiple MiddlewareAfter can be added and are processed in order.
func (s *internalServer[T]) addMiddlewareAfter(fn MiddlewareAfter[T]) *internalServer[T] {
	s.afterMiddlewares = append(s.afterMiddlewares, fn)
	return s
}

// setProcHandler registers the handler for the provided procedure name
func (s *internalServer[T]) setProcHandler(
	procName string,
	handler func(context T, input json.RawMessage) (any, error),
) *internalServer[T] {
	s.procHandlers[procName] = handler
	return s
}

// setProcInputProcessor registers the input processor for the provided procedure name
func (s *internalServer[T]) setProcInputProcessor(
	procName string,
	processor func(context T, input any) (any, error),
) *internalServer[T] {
	s.procInputProcessors[procName] = processor
	return s
}

// setStreamHandler registers the handler for the provided stream name
func (s *internalServer[T]) setStreamHandler(
	streamName string,
	handler func(context T, input json.RawMessage, emit func(any) error) error,
) *internalServer[T] {
	s.streamHandlers[streamName] = handler
	return s
}

// setStreamInputProcessor registers the input processor for the provided stream name
func (s *internalServer[T]) setStreamInputProcessor(
	streamName string,
	processor func(context T, input any) (any, error),
) *internalServer[T] {
	s.streamInputProcessors[streamName] = processor
	return s
}

// handleRequest processes an incoming RPC request
func (s *internalServer[T]) handleRequest(reqResProvider ServerRequestResponseProvider[T]) error {
	if reqResProvider == nil {
		res := Response[any]{
			Ok:    false,
			Error: Error{Message: "ServerRequestResponseProvider is nil, please provide a valid provider"},
		}
		return s.writeProcResponse(reqResProvider, res)
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

	if isStream {
		return s.handleStreamRequest(jsonBody.Name, jsonBody.Input, reqResProvider)
	}

	return s.handleProcRequest(jsonBody.Name, jsonBody.Input, reqResProvider)
}

// writeProcResponse writes a procedure response to the client
func (s *internalServer[T]) writeProcResponse(
	reqResProvider ServerRequestResponseProvider[T],
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
	procName string,
	rawInput json.RawMessage,
	reqResProvider ServerRequestResponseProvider[T],
) error {
	currentContext := reqResProvider.RequestGetInitialContext()
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
	if _, ok := s.procHandlers[procName]; response.Ok && !ok {
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
		for _, fn := range s.beforeMiddlewares {
			var err error
			if currentContext, err = fn("proc", procName, currentContext); err != nil {
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
		if output, err := s.procHandlers[procName](currentContext, rawInput); err != nil {
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
	for _, fn := range s.afterMiddlewares {
		response = fn("proc", procName, currentContext, response)
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
	streamName string,
	rawInput json.RawMessage,
	reqResProvider ServerRequestResponseProvider[T],
) error {
	currentContext := reqResProvider.RequestGetInitialContext()

	// emit sends a response event to the client.
	// If the client disconnects, it returns an error.
	emit := func(data any) error {
		response := Response[any]{
			Ok:     true,
			Output: data,
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			return err
		}

		resPayload := fmt.Sprintf("data: %s\n\n", jsonData)
		_, err = reqResProvider.ResponseWrite([]byte(resPayload))
		if err != nil {
			return err
		}

		reqResProvider.ResponseFlush()
		return nil
	}

	// Upgrade the request to a stream
	reqResProvider.ResponseSetHeader("Content-Type", "text/event-stream")
	reqResProvider.ResponseSetHeader("Cache-Control", "no-cache")
	reqResProvider.ResponseSetHeader("Connection", "keep-alive")

	// Validate stream name
	if !slices.Contains(s.streamNames, streamName) {
		response := Response[any]{
			Ok: false,
			Error: Error{
				Message: streamName + " stream not found",
				Details: map[string]any{"stream": streamName},
			},
		}
		return emit(response)
	}

	// Validate stream implementation
	if _, ok := s.streamHandlers[streamName]; !ok {
		response := Response[any]{
			Ok: false,
			Error: Error{
				Message: streamName + " stream not implemented",
				Details: map[string]any{"stream": streamName},
			},
		}
		return emit(response)
	}

	// Execute Before middlewares
	for _, fn := range s.beforeMiddlewares {
		var err error
		if currentContext, err = fn("stream", streamName, currentContext); err != nil {
			response := Response[any]{
				Ok:    false,
				Error: asError(err),
			}
			return emit(response)
		}
	}

	// Run the stream handler and wait for it to finish
	err := s.streamHandlers[streamName](currentContext, rawInput, emit)
	if err != nil {
		response := Response[any]{
			Ok:    false,
			Error: asError(err),
		}
		return emit(response)
	}

	return nil
}
