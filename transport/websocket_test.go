package transport

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/coder/websocket"
)

func newEchoServer(t *testing.T) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			Subprotocols: []string{BidiSubprotocol},
		})
		if err != nil {
			return
		}

		defer func() { _ = conn.CloseNow() }()

		ctx := r.Context()
		for {
			typ, data, err := conn.Read(ctx)
			if err != nil {
				return
			}

			if err := conn.Write(ctx, typ, data); err != nil {
				return
			}
		}
	}))

	t.Cleanup(srv.Close)

	return srv
}

func dialTestServer(t *testing.T, srv *httptest.Server) *WebSocketTransport {
	t.Helper()

	endpoint, err := url.Parse("ws" + strings.TrimPrefix(srv.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}

	tr, err := DialBidi(context.Background(), endpoint)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	return tr
}

func TestWebSocketRoundTrip(t *testing.T) {
	srv := newEchoServer(t)
	tr := dialTestServer(t, srv)

	t.Cleanup(func() { _ = tr.Close() })

	if got := tr.Subprotocol(); got != BidiSubprotocol {
		t.Fatalf("subprotocol: got %q, want %q", got, BidiSubprotocol)
	}

	ctx := context.Background()

	frame := `{"id":1,"method":"session.status","params":{}}`
	if err := tr.Send(ctx, []byte(frame)); err != nil {
		t.Fatalf("send: %v", err)
	}

	got, err := tr.Receive(ctx)
	if err != nil {
		t.Fatalf("receive: %v", err)
	}

	if string(got) != frame {
		t.Fatalf("got %q, want %q", got, frame)
	}
}

func TestWebSocketClose(t *testing.T) {
	srv := newEchoServer(t)
	tr := dialTestServer(t, srv)

	if err := tr.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if err := tr.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}

	if err := tr.Send(context.Background(), []byte(`{"id":1}`)); !errors.Is(err, ErrClosed) {
		t.Fatalf("send after close: got %v, want ErrClosed", err)
	}
}
