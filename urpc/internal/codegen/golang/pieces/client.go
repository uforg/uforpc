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
	globalHeaders http.Header

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

// withGlobalHeaders sets headers that will be attached to every request.
// The provided header map is copied, so further mutations will not affect the
// client after construction.
func withGlobalHeaders(h http.Header) internalClientOption {
	return func(c *internalClient) {
		if h == nil {
			return
		}
		c.globalHeaders = clientHelperCloneHeader(h)
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

// WithHTTPClient sets a custom *http.Client.
func (b *internalClientBuilder) WithHTTPClient(hc *http.Client) *internalClientBuilder {
	b.opts = append(b.opts, withHTTPClient(hc))
	return b
}

// WithGlobalHeaders adds global headers that will be sent with every request.
func (b *internalClientBuilder) WithGlobalHeaders(h http.Header) *internalClientBuilder {
	b.opts = append(b.opts, withGlobalHeaders(h))
	return b
}

// WithMaxStreamEventDataSize overrides the default maximum SSE payload size in
// bytes.
func (b *internalClientBuilder) WithMaxStreamEventDataSize(size int) *internalClientBuilder {
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
	extraHeaders http.Header,
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

	// Apply configured headers (global headers + per-call extras overriding).
	clientHelperAddHeaders(req.Header, c.globalHeaders)
	clientHelperAddHeaders(req.Header, extraHeaders)

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
	headers http.Header
}

// WithHeader adds a header to this procedure invocation.
func (p *procCallBuilder) WithHeader(key, value string) *procCallBuilder {
	p.headers.Add(key, value)
	return p
}

// WithHeaders adds all headers h to this invocation (merged).
func (p *procCallBuilder) WithHeaders(h http.Header) *procCallBuilder {
	clientHelperAddHeaders(p.headers, h)
	return p
}

// Execute performs the RPC call and returns the Response.
func (p *procCallBuilder) Execute(ctx context.Context) Response[json.RawMessage] {
	return p.client.callProc(ctx, p.name, p.input, p.headers)
}

// newProcCallBuilder creates a builder for calling the given procedure.
func (c *internalClient) newProcCallBuilder(name string, input any) *procCallBuilder {
	return &procCallBuilder{
		client:  c,
		name:    name,
		input:   input,
		headers: http.Header{},
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
	extraHeaders http.Header,
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
	clientHelperAddHeaders(req.Header, c.globalHeaders)
	clientHelperAddHeaders(req.Header, extraHeaders)

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
	headers          http.Header
	maxEventDataSize int
}

// WithHeader adds a header to this stream invocation.
func (s *streamCall) WithHeader(key, value string) *streamCall {
	s.headers.Add(key, value)
	return s
}

// WithHeaders adds multiple headers.
func (s *streamCall) WithHeaders(h http.Header) *streamCall {
	clientHelperAddHeaders(s.headers, h)
	return s
}

// WithMaxEventDataSize sets the maximum event data size.
func (s *streamCall) WithMaxEventDataSize(size int) *streamCall {
	s.maxEventDataSize = size
	return s
}

// Execute starts the stream and returns the channel of events.
func (s *streamCall) Execute(ctx context.Context) <-chan Response[json.RawMessage] {
	return s.client.stream(ctx, s.name, s.input, s.headers, s.maxEventDataSize)
}

// newStreamCallBuilder creates a builder for the given stream.
func (c *internalClient) newStreamCallBuilder(name string, input any) *streamCall {
	return &streamCall{
		client:           c,
		name:             name,
		input:            input,
		headers:          http.Header{},
		maxEventDataSize: 0,
	}
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

// clientHelperAddHeaders copies all headers from src into dst. Duplicate keys are appended
// (preserving existing entries).
func clientHelperAddHeaders(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

// clientHelperCloneHeader returns a deep copy of h.
func clientHelperCloneHeader(h http.Header) http.Header {
	out := make(http.Header, len(h))
	for k, vs := range h {
		cpy := make([]string, len(vs))
		copy(cpy, vs)
		out[k] = cpy
	}
	return out
}
