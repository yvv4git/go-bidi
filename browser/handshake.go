package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yvv4git/go-bidi/protocol"
	"github.com/yvv4git/go-bidi/transport"
)

const (
	_sessionPath       = "/session"
	_headerContentType = "Content-Type"
	_contentTypeJSON   = "application/json"
	_capWebSocketURL   = "webSocketUrl"
)

// HandshakeResult is the outcome of the WebDriver classic handshake. The
// transport is connected and speaks the webdriver.bidi subprotocol; the
// caller owns it and must close it when done.
type HandshakeResult struct {
	SessionID    string
	WebSocketURL string
	Capabilities map[string]any
	Transport    transport.Transport
}

// Option customizes Handshake, Connect, ConnectEndpoint and NewClient.
type Option func(*config)

type config struct {
	httpClient   *http.Client
	capabilities map[string]any
	timeout      time.Duration
	subprotocols []string
}

type sessionResponse struct {
	SessionID    string         `json:"sessionId"`
	Capabilities map[string]any `json:"capabilities"`
}

// WithHTTPClient sets the client used for the session creation request.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		c.httpClient = client
	}
}

// WithCapabilities merges extra alwaysMatch capabilities into the request.
func WithCapabilities(caps map[string]any) Option {
	return func(c *config) {
		c.capabilities = make(map[string]any, len(caps))
		for key, value := range caps {
			c.capabilities[key] = value
		}
	}
}

// WithTimeout sets the default deadline applied to commands that do not
// carry their own context deadline.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}

// WithSubprotocol overrides the WebSocket subprotocols negotiated by
// Handshake and ConnectEndpoint. The default is transport.BidiSubprotocol.
func WithSubprotocol(subprotocols ...string) Option {
	return func(c *config) {
		c.subprotocols = append([]string(nil), subprotocols...)
	}
}

// Handshake creates a session over the WebDriver classic endpoint at addr,
// then connects the returned WebSocket and negotiates webdriver.bidi.
func Handshake(
	ctx context.Context,
	addr string,
	opts ...Option,
) (*HandshakeResult, error) {
	cfg := config{httpClient: http.DefaultClient}
	for _, opt := range opts {
		opt(&cfg)
	}

	response, err := cfg.createSession(ctx, addr)
	if err != nil {
		return nil, err
	}

	wsURL, err := response.webSocketURL()
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("parse websocket url: %w", err)
	}

	tr, err := cfg.dial(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("dial websocket: %w", err)
	}

	return response.result(wsURL, tr), nil
}

func (c *config) dial(ctx context.Context, endpoint *url.URL) (transport.Transport, error) {
	if len(c.subprotocols) > 0 {
		return transport.DialWebSocket(ctx, endpoint, c.subprotocols...)
	}

	return transport.DialBidi(ctx, endpoint)
}

func (c *config) newSessionRequest(
	ctx context.Context,
	addr string,
) (*http.Request, error) {
	body, err := c.requestBody()
	if err != nil {
		return nil, err
	}

	endpoint := strings.TrimRight(addr, "/") + _sessionPath

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new session request: %w", err)
	}

	req.Header.Set(_headerContentType, _contentTypeJSON)

	return req, nil
}

func (c *config) requestBody() ([]byte, error) {
	alwaysMatch := make(map[string]any, len(c.capabilities)+1)

	for key, value := range c.capabilities {
		alwaysMatch[key] = value
	}

	if _, ok := alwaysMatch[_capWebSocketURL]; !ok {
		alwaysMatch[_capWebSocketURL] = true
	}

	request := struct {
		Capabilities struct {
			AlwaysMatch map[string]any `json:"alwaysMatch"`
		} `json:"capabilities"`
	}{}
	request.Capabilities.AlwaysMatch = alwaysMatch

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode capabilities: %w", err)
	}

	return body, nil
}

func (c *config) createSession(
	ctx context.Context,
	addr string,
) (*sessionResponse, error) {
	req, err := c.newSessionRequest(ctx, addr)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create session: unexpected status %d", resp.StatusCode)
	}

	var decoded struct {
		Value sessionResponse `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode session response: %w", err)
	}

	return &decoded.Value, nil
}

func (r *sessionResponse) webSocketURL() (string, error) {
	raw, ok := r.Capabilities[_capWebSocketURL]
	if !ok {
		return "", errors.New("handshake: response has no webSocketUrl")
	}

	wsURL, ok := raw.(string)
	if !ok || wsURL == "" {
		return "", errors.New("handshake: webSocketUrl is not a non-empty string")
	}

	return wsURL, nil
}

func (r *sessionResponse) result(wsURL string, tr transport.Transport) *HandshakeResult {
	return &HandshakeResult{
		SessionID:    r.SessionID,
		WebSocketURL: wsURL,
		Capabilities: r.Capabilities,
		Transport:    tr,
	}
}

// buildConfig extracts a config from options, filling in defaults.
func buildConfig(opts []Option) config {
	cfg := config{httpClient: http.DefaultClient}
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// httpToWSURL converts an HTTP or WebSocket base address to the BiDi
// WebSocket session endpoint ws://host[:port]/session.
func httpToWSURL(addr string) string {
	addr = strings.TrimRight(addr, "/")

	u, err := url.Parse(addr)
	if err != nil {
		return "ws://" + addr + "/session"
	}

	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		u.Scheme = "ws"
	}

	u.Path = _sessionPath

	return u.String()
}

// ConnectBiDi connects directly over a BiDi WebSocket, creating a session
// with session.new.  addr may be an HTTP URL (http://host:port), a
// WebSocket URL (ws://host:port/session), or a bare host:port pair.
// This is the connection flow used by Firefox 158+ and newer BiDi
// implementations that no longer expose the classic POST /session
// endpoint.
func ConnectBiDi(
	ctx context.Context,
	addr string,
	opts ...Option,
) (*Browser, error) {
	cfg := buildConfig(opts)

	wsURL := httpToWSURL(addr)

	endpoint, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("parse websocket url: %w", err)
	}

	tr, err := cfg.dial(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("dial biDi: %w", err)
	}

	alwaysMatch := make(map[string]any, len(cfg.capabilities)+1)
	for key, value := range cfg.capabilities {
		alwaysMatch[key] = value
	}

	alwaysMatch[_capWebSocketURL] = true

	client := NewClient(tr, opts...)

	params := protocol.SessionNewParams{
		Capabilities: protocol.SessionCapabilities{
			AlwaysMatch: alwaysMatch,
		},
	}

	session, err := CreateSession(ctx, client, params)
	if err != nil {
		_ = client.Close()

		return nil, err
	}

	return newBrowser(client, session), nil
}
