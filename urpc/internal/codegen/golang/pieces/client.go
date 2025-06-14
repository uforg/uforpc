//nolint:unused
package pieces

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Client Types
// -----------------------------------------------------------------------------

// internalClient manages HTTP communication with a UFO RPC server. It supports
// invoking procedure endpoints (standard request / response) as well as
// subscribing to stream endpoints delivered via Server-Sent Events (SSE).
//
// This implementation purposefully keeps all identifiers unexported because the
// code generator will copy this file verbatim into the generated client code
// and wrap it with a public, type-safe manner. Hiding the implementation
// details ensures that only the generated wrapper is exposed to end-users.
//
// Design goals:
//   1. Minimal dependencies – it only relies on Go's standard library.
//   2. Idiomatic, well-documented code that is trivial to wrap.
//   3. Flexibility – callers can override the underlying *http.Client or the
//      base URL used to contact the server.
//   4. Robustness – context cancellation is honoured and every goroutine is
//      cleaned up.
//
// The transport contract is the mirror image of the server implementation:
//   • Procedures  – JSON request body, JSON response body.
//   • Streams     – JSON request body, SSE response with each event encoded as
//                   a "data: { … }\n\n" payload identical to Response[T].

// internalClient is the core engine used by the generated client abstraction. It is
// thread-safe and can be reused across concurrent requests.
//
// The client validates the requested operation name against the schema to fail
// fast when a typo occurs in the generated wrapper (a programmer error).
//
// The zero value is not usable – use newInternalClient to construct one.
type internalClient struct {
	// Immutable after construction.
	baseURL        string
	httpClient     *http.Client
	procNames      []string
	procNamesMap   map[string]bool
	streamNames    []string
	streamNamesMap map[string]bool

	// header configuration (global on every request)
	globalHeaders map[string]string

	// SSE data size configuration in bytes
	maxStreamEventDataSize int
}

// internalClientOption represents a configuration option for internalClient.
type internalClientOption func(*internalClient)

// withHTTPClient supplies a custom *http.Client. If nil, the default client is
// used. Callers can leverage this to inject time-outs, proxies, or a transport
// with advanced TLS settings.
func withHTTPClient(hc *http.Client) internalClientOption {
	return func(c *internalClient) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// withGlobalHeader sets a header that will be attached to every request.
func withGlobalHeader(key string, value string) internalClientOption {
	return func(c *internalClient) {
		c.globalHeaders[key] = value
	}
}

// withMaxStreamEventDataSize sets a global default maximum size for SSE data
// events when the per-stream size is not configured.
func withMaxStreamEventDataSize(size int) internalClientOption {
	return func(c *internalClient) {
		if size > 0 {
			c.maxStreamEventDataSize = size
		}
	}
}

// newInternalClient creates a new internalClient capable of talking to the UFO
// RPC server described by procNames and streamNames.
//
// The caller can optionally pass functional options to tweak the configuration
// (base URL, custom *http.Client, …).
func newInternalClient(
	baseURL string,
	procNames []string,
	streamNames []string,
	opts ...internalClientOption,
) *internalClient {
	procMap := make(map[string]bool, len(procNames))
	for _, n := range procNames {
		procMap[n] = true
	}
	streamMap := make(map[string]bool, len(streamNames))
	for _, n := range streamNames {
		streamMap[n] = true
	}

	// Sensible defaults – baseURL will be "" (relative) and the default
	// http.Client will follow redirects and has no timeout.
	cli := &internalClient{
		baseURL:                strings.TrimRight(baseURL, "/"),
		httpClient:             http.DefaultClient,
		procNames:              procNames,
		procNamesMap:           procMap,
		streamNames:            streamNames,
		streamNamesMap:         streamMap,
		maxStreamEventDataSize: 1 << 20, // 1 MiB per event by default
	}

	// Apply functional options.
	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

/* Internal client builder */

// internalClientBuilder helps constructing an internalClient using chained
// configuration methods before calling Build().
type internalClientBuilder struct {
	baseURL     string
	procNames   []string
	streamNames []string
	opts        []internalClientOption
}

// newClientBuilder creates a builder with the schema information (procedure and
// stream names). Generated code will pass the automatically produced slices.
func newClientBuilder(baseURL string, procNames, streamNames []string) *internalClientBuilder {
	return &internalClientBuilder{
		baseURL:     baseURL,
		procNames:   procNames,
		streamNames: streamNames,
		opts:        []internalClientOption{},
	}
}

// withHTTPClient sets a custom *http.Client.
func (b *internalClientBuilder) withHTTPClient(hc *http.Client) *internalClientBuilder {
	b.opts = append(b.opts, withHTTPClient(hc))
	return b
}

// withGlobalHeader adds a global header that will be sent with every request.
func (b *internalClientBuilder) withGlobalHeader(key, value string) *internalClientBuilder {
	b.opts = append(b.opts, withGlobalHeader(key, value))
	return b
}

// withMaxStreamEventDataSize overrides the default maximum SSE payload size in
// bytes.
func (b *internalClientBuilder) withMaxStreamEventDataSize(size int) *internalClientBuilder {
	b.opts = append(b.opts, withMaxStreamEventDataSize(size))
	return b
}

// Build creates the internalClient applying all accumulated options.
func (b *internalClientBuilder) Build() *internalClient {
	return newInternalClient(b.baseURL, b.procNames, b.streamNames, b.opts...)
}

/* Internal client procedure handling */

// callProc invokes the given procedure with the provided input and returns the
// raw JSON response from the server wrapped in a Response object.
//
// Any transport or decoding error is converted into a Response with Ok set to
// false and the Error field describing the failure.
func (c *internalClient) callProc(
	ctx context.Context,
	procName string,
	input any,
	extraHeaders map[string]string,
) Response[json.RawMessage] {
	if !c.procNamesMap[procName] {
		return Response[json.RawMessage]{
			Ok: false,
			Error: Error{
				Category: "ClientError",
				Code:     "INVALID_PROC",
				Message:  fmt.Sprintf("%s procedure not found in schema", procName),
				Details:  map[string]any{"procedure": procName},
			},
		}
	}

	// Encode the input.
	var payload []byte
	var err error
	if input == nil {
		payload = []byte("{}")
	} else {
		payload, err = json.Marshal(input)
		if err != nil {
			return Response[json.RawMessage]{
				Ok: false,
				Error: Error{
					Category: "ClientError",
					Code:     "ENCODE_INPUT",
					Message:  fmt.Sprintf("failed to marshal input for %s: %v", procName, err),
				},
			}
		}
	}

	// Build URL – <baseURL>/<procName> . Leading slash added if missing.
	url := c.baseURL + "/" + procName

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return Response[json.RawMessage]{
			Ok:    false,
			Error: asError(fmt.Errorf("failed to create HTTP request: %w", err)),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Apply headers: global + per-call extras.
	for key, value := range c.globalHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Response[json.RawMessage]{
			Ok:    false,
			Error: asError(fmt.Errorf("http request failed: %w", err)),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Response[json.RawMessage]{
			Ok: false,
			Error: Error{
				Category: "HTTPError",
				Code:     "BAD_STATUS",
				Message:  fmt.Sprintf("unexpected HTTP status: %s", resp.Status),
				Details:  map[string]any{"status": resp.StatusCode},
			},
		}
	}

	// Decode the generic response first so that we can decide what to do next.
	var raw struct {
		Ok     bool            `json:"ok"`
		Output json.RawMessage `json:"output"`
		Error  Error           `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return Response[json.RawMessage]{
			Ok:    false,
			Error: asError(fmt.Errorf("failed to decode UFO RPC response: %w", err)),
		}
	}

	if !raw.Ok {
		return Response[json.RawMessage]{
			Ok:    false,
			Error: raw.Error,
		}
	}

	return Response[json.RawMessage]{
		Ok:     true,
		Output: raw.Output,
	}
}

// procCallBuilder is a fluent builder for invoking a procedure.
type procCallBuilder struct {
	client  *internalClient
	name    string
	input   any
	headers map[string]string
}

// withHeader adds a header to this procedure invocation.
func (p *procCallBuilder) withHeader(key, value string) *procCallBuilder {
	p.headers[key] = value
	return p
}

// execute performs the RPC call and returns the Response.
func (p *procCallBuilder) execute(ctx context.Context) Response[json.RawMessage] {
	return p.client.callProc(ctx, p.name, p.input, p.headers)
}

// newProcCallBuilder creates a builder for calling the given procedure.
func (c *internalClient) newProcCallBuilder(name string, input any) *procCallBuilder {
	return &procCallBuilder{
		client:  c,
		name:    name,
		input:   input,
		headers: map[string]string{},
	}
}

/* Internal client stream handling */

// stream establishes a Server-Sent Events subscription for the given stream
// name. Each received event is forwarded on the returned channel until ctx is
// cancelled or the server closes the connection.
//
// The channel is closed on termination and MUST be fully drained by the caller
// to avoid goroutine leaks.
func (c *internalClient) stream(
	ctx context.Context,
	streamName string,
	input any,
	extraHeaders map[string]string,
	maxEventDataSize int,
) <-chan Response[json.RawMessage] {
	if !c.streamNamesMap[streamName] {
		ch := make(chan Response[json.RawMessage], 1)
		ch <- Response[json.RawMessage]{
			Ok: false,
			Error: Error{
				Category: "ClientError",
				Code:     "INVALID_STREAM",
				Message:  fmt.Sprintf("%s stream not found in schema", streamName),
				Details:  map[string]any{"stream": streamName},
			},
		}
		close(ch)
		return ch
	}

	// Encode input.
	var payload []byte
	var err error
	if input == nil {
		payload = []byte("{}")
	} else {
		payload, err = json.Marshal(input)
		if err != nil {
			ch := make(chan Response[json.RawMessage], 1)
			ch <- Response[json.RawMessage]{
				Ok:    false,
				Error: asError(fmt.Errorf("failed to marshal input for %s: %w", streamName, err)),
			}
			close(ch)
			return ch
		}
	}

	// Build URL – <baseURL>/<streamName>
	url := c.baseURL + "/" + streamName

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		ch := make(chan Response[json.RawMessage], 1)
		ch <- Response[json.RawMessage]{
			Ok:    false,
			Error: asError(fmt.Errorf("failed to create HTTP request: %w", err)),
		}
		close(ch)
		return ch
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	// Apply headers: global + per-call extras.
	for key, value := range c.globalHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		ch := make(chan Response[json.RawMessage], 1)
		ch <- Response[json.RawMessage]{
			Ok:    false,
			Error: asError(fmt.Errorf("stream request failed: %w", err)),
		}
		close(ch)
		return ch
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		ch := make(chan Response[json.RawMessage], 1)
		ch <- Response[json.RawMessage]{
			Ok: false,
			Error: Error{
				Category: "HTTPError",
				Code:     "BAD_STATUS",
				Message:  fmt.Sprintf("unexpected HTTP status: %s", resp.Status),
				Details:  map[string]any{"status": resp.StatusCode},
			},
		}
		close(ch)
		return ch
	}

	// Channel for events.
	events := make(chan Response[json.RawMessage])

	go func() {
		defer close(events)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		// Determine the maximum token size.
		maxSize := c.maxStreamEventDataSize
		if maxEventDataSize > 0 {
			maxSize = maxEventDataSize
		}

		scanner.Buffer(make([]byte, 0, 64*1024), maxSize)

		var dataBuf bytes.Buffer

		flush := func() {
			if dataBuf.Len() == 0 {
				return
			}
			var evt Response[json.RawMessage]
			if err := json.Unmarshal(dataBuf.Bytes(), &evt); err != nil {
				// Protocol violation – stop the stream.
				events <- Response[json.RawMessage]{
					Ok:    false,
					Error: asError(fmt.Errorf("received invalid SSE payload: %v", err)),
				}
				return
			}
			select {
			case events <- evt:
			case <-ctx.Done():
			}
			dataBuf.Reset()
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if !scanner.Scan() {
				// EOF or error – exit.
				return
			}
			line := scanner.Text()
			if line == "" { // Blank line marks end of event.
				flush()
				continue
			}
			if strings.HasPrefix(line, "data:") {
				// Strip the "data:" prefix and optional leading space.
				chunk := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				dataBuf.WriteString(chunk)
			}
			// Everything else is ignored (e.g., id:, retry:, …).
		}
	}()

	return events
}

// streamCall is a fluent builder for SSE subscriptions.
type streamCall struct {
	client           *internalClient
	name             string
	input            any
	headers          map[string]string
	maxEventDataSize int
}

// withHeader adds a header to this stream invocation.
func (s *streamCall) withHeader(key, value string) *streamCall {
	s.headers[key] = value
	return s
}

// withMaxEventDataSize sets the maximum event data size.
func (s *streamCall) withMaxEventDataSize(size int) *streamCall {
	s.maxEventDataSize = size
	return s
}

// execute starts the stream and returns the channel of events.
func (s *streamCall) execute(ctx context.Context) <-chan Response[json.RawMessage] {
	return s.client.stream(ctx, s.name, s.input, s.headers, s.maxEventDataSize)
}

// newStreamCallBuilder creates a builder for the given stream.
func (c *internalClient) newStreamCallBuilder(name string, input any) *streamCall {
	return &streamCall{
		client:           c,
		name:             name,
		input:            input,
		headers:          map[string]string{},
		maxEventDataSize: 0,
	}
}
