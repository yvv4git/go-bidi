package transport

import (
	"context"
	"net/url"

	"github.com/coder/websocket"
	"go.uber.org/atomic"
)

// BidiSubprotocol is the WebSocket subprotocol negotiated for WebDriver BiDi.
const BidiSubprotocol = "webdriver.bidi"

// WebSocketTransport carries messages over a single WebSocket connection.
// The underlying connection is safe for concurrent writes, so the transport
// only tracks its own closed state.
type WebSocketTransport struct {
	conn   *websocket.Conn
	closed atomic.Bool
}

var _ Transport = (*WebSocketTransport)(nil)

// DialWebSocket establishes a WebSocket transport to the browser's
// endpoint, optionally negotiating the given subprotocols.
func DialWebSocket(
	ctx context.Context, endpoint *url.URL, subprotocols ...string,
) (*WebSocketTransport, error) {
	var opts *websocket.DialOptions
	if len(subprotocols) > 0 {
		opts = &websocket.DialOptions{Subprotocols: subprotocols}
	}

	conn, _, err := websocket.Dial(ctx, endpoint.String(), opts)
	if err != nil {
		return nil, err
	}

	return &WebSocketTransport{conn: conn}, nil
}

// DialBidi establishes a WebDriver BiDi transport, negotiating the
// "webdriver.bidi" subprotocol.
func DialBidi(ctx context.Context, endpoint *url.URL) (*WebSocketTransport, error) {
	return DialWebSocket(ctx, endpoint, BidiSubprotocol)
}

// Subprotocol reports the WebSocket subprotocol negotiated with the server.
func (t *WebSocketTransport) Subprotocol() string {
	return t.conn.Subprotocol()
}

// Send writes a single JSON frame.
func (t *WebSocketTransport) Send(ctx context.Context, data []byte) error {
	if t.closed.Load() {
		return ErrClosed
	}

	return t.conn.Write(ctx, websocket.MessageText, data)
}

// Receive reads the next JSON frame.
func (t *WebSocketTransport) Receive(ctx context.Context) ([]byte, error) {
	_, data, err := t.conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Close tears down the WebSocket connection. It is safe to call repeatedly.
func (t *WebSocketTransport) Close() error {
	if t.closed.Swap(true) {
		return nil
	}

	return t.conn.CloseNow()
}
