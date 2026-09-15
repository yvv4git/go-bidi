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

// HandshakeOption customizes Handshake.
type HandshakeOption func(*handshakeConfig)

type handshakeConfig struct {
	httpClient   *http.Client
	capabilities map[string]any
}

type sessionResponse struct {
	SessionID    string         `json:"sessionId"`
	Capabilities map[string]any `json:"capabilities"`
}

// WithHTTPClient sets the client used for the session creation request.
func WithHTTPClient(client *http.Client) HandshakeOption {
	return func(c *handshakeConfig) {
		c.httpClient = client
	}
}

// WithCapabilities merges extra alwaysMatch capabilities into the request.
func WithCapabilities(caps map[string]any) HandshakeOption {
	return func(c *handshakeConfig) {
		c.capabilities = make(map[string]any, len(caps))
		for key, value := range caps {
			c.capabilities[key] = value
		}
	}
}

// Handshake creates a session over the WebDriver classic endpoint at addr,
// then connects the returned WebSocket and negotiates webdriver.bidi.
func Handshake(
	ctx context.Context,
	addr string,
	opts ...HandshakeOption,
) (*HandshakeResult, error) {
	cfg := handshakeConfig{httpClient: http.DefaultClient}
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

	tr, err := transport.DialBidi(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("dial websocket: %w", err)
	}

	return response.result(wsURL, tr), nil
}

func (c *handshakeConfig) newSessionRequest(
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

func (c *handshakeConfig) requestBody() ([]byte, error) {
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

func (c *handshakeConfig) createSession(
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
