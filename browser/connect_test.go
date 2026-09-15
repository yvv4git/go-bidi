package browser

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yvv4git/go-bidi/transport"
)

func TestConnect(t *testing.T) {
	tr := newFakeTransport()

	go respondSuccess(t, tr, `{"sessionId":"session-42","capabilities":{}}`)

	b, err := Connect(context.Background(), tr)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	if got, want := b.Session().ID(), "session-42"; got != want {
		t.Errorf("Session().ID() = %q, want %q", got, want)
	}

	go respondSuccess(t, tr, `{}`)

	if err := b.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestConnectSessionError(t *testing.T) {
	tr := newFakeTransport()

	go respondError(t, tr, "session not created", "boom")

	if _, err := Connect(context.Background(), tr); err == nil {
		t.Fatal("Connect: expected error")
	}

	select {
	case <-tr.closed:
	default:
		t.Error("transport not closed after failed Connect")
	}
}

func TestBrowserNewPageAndTracking(t *testing.T) {
	tr := newFakeTransport()

	go func() {
		respondSuccess(t, tr, `{"sessionId":"s1","capabilities":{}}`)
		respondSuccess(t, tr, `{"context":"ctx-2"}`)
		respondSuccess(t, tr, `{"context":"ctx-1"}`)
	}()

	b, err := Connect(context.Background(), tr)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	p1, err := b.NewPage(context.Background())
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	if p1.ID() != "ctx-2" {
		t.Errorf("page 1 ID = %q, want %q", p1.ID(), "ctx-2")
	}

	p2, err := b.NewPage(context.Background())
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	if p2.ID() != "ctx-1" {
		t.Errorf("page 2 ID = %q, want %q", p2.ID(), "ctx-1")
	}

	if b.Page("ctx-2") != p1 {
		t.Error("Page(\"ctx-2\") does not return the tracked page")
	}

	unknown := b.Page("ctx-untracked")
	if unknown.ID() != "ctx-untracked" {
		t.Errorf("Page ID = %q, want %q", unknown.ID(), "ctx-untracked")
	}

	pages := b.Pages()
	if len(pages) != 3 {
		t.Fatalf("Pages() length = %d, want 3", len(pages))
	}
	for i, want := range []string{"ctx-1", "ctx-2", "ctx-untracked"} {
		if pages[i].ID() != want {
			t.Errorf("Pages()[%d] ID = %q, want %q", i, pages[i].ID(), want)
		}
	}

	go respondSuccess(t, tr, `{}`)

	if err := b.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestBrowserCloseStopsCommands(t *testing.T) {
	tr := newFakeTransport()

	go respondSuccess(t, tr, `{"sessionId":"s1","capabilities":{}}`)

	b, err := Connect(context.Background(), tr)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	go respondSuccess(t, tr, `{}`)

	if err := b.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if _, err := b.Session().Status(context.Background()); !errors.Is(err, ErrClosed) {
		t.Fatalf("Status error = %v, want ErrClosed", err)
	}
}

func TestClientConfiguredTimeout(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr, WithTimeout(25*time.Millisecond))
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	err := client.Call(context.Background(), "test.slow", nil, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Call error = %v, want DeadlineExceeded", err)
	}
}

func TestClientKeepsExplicitDeadline(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr, WithTimeout(time.Hour))
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := client.Call(ctx, "test.slow", nil, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Call error = %v, want DeadlineExceeded", err)
	}
}

func TestConnectEndpoint(t *testing.T) {
	var srv *httptest.Server

	srv = newHandshakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
		wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
		payload := map[string]any{
			"value": map[string]any{
				"sessionId":    "s1",
				"capabilities": map[string]any{_capWebSocketURL: wsURL},
			},
		}

		w.Header().Set(_headerContentType, _contentTypeJSON)

		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("encode response: %v", err)
		}
	})

	b, err := ConnectEndpoint(
		context.Background(),
		srv.URL,
		WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("ConnectEndpoint: %v", err)
	}

	t.Cleanup(func() {
		if err := b.client.Close(); err != nil {
			t.Errorf("close client: %v", err)
		}
	})

	if b.Session().ID() != "s1" {
		t.Errorf("Session().ID() = %q, want %q", b.Session().ID(), "s1")
	}

	if b.client.timeout != time.Second {
		t.Errorf("client timeout = %v, want %v", b.client.timeout, time.Second)
	}
}

func TestConnectEndpointSubprotocolOption(t *testing.T) {
	var srv *httptest.Server

	srv = newHandshakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
		var body struct {
			Value struct {
				SessionID    string         `json:"sessionId"`
				Capabilities map[string]any `json:"capabilities"`
			} `json:"value"`
		}
		body.Value.SessionID = "s1"
		body.Value.Capabilities = map[string]any{
			_capWebSocketURL: "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws",
		}

		w.Header().Set(_headerContentType, _contentTypeJSON)

		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Errorf("encode response: %v", err)
		}
	})

	b, err := ConnectEndpoint(
		context.Background(),
		srv.URL,
		WithSubprotocol(transport.BidiSubprotocol),
	)
	if err != nil {
		t.Fatalf("ConnectEndpoint: %v", err)
	}

	t.Cleanup(func() {
		if err := b.client.Close(); err != nil {
			t.Errorf("close client: %v", err)
		}
	})

	if b.Session().ID() != "s1" {
		t.Errorf("Session().ID() = %q, want %q", b.Session().ID(), "s1")
	}
}

func TestConnectEndpointStatusError(t *testing.T) {
	srv := newHandshakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
		if _, err := io.WriteString(w, "nope"); err != nil {
			t.Errorf("write response: %v", err)
		}
	})

	if _, err := ConnectEndpoint(context.Background(), srv.URL); err == nil {
		t.Fatal("ConnectEndpoint: expected error")
	}
}
