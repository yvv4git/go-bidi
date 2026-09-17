package browser

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/coder/websocket"

	"github.com/yvv4git/go-bidi/transport"
)

// roundTripFunc adapts a function into an http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newHandshakeServer(t *testing.T, session http.HandlerFunc) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc(_sessionPath, session)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			Subprotocols: []string{transport.BidiSubprotocol},
		})
		if err != nil {
			return
		}

		defer func() { _ = conn.CloseNow() }()

		if _, _, err := conn.Read(r.Context()); err != nil {
			return
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

func TestHandshake(t *testing.T) {
	var (
		mu   sync.Mutex
		body []byte
	)

	var srv *httptest.Server

	srv = newHandshakeServer(t, func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		body = data
		mu.Unlock()

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

	result, err := Handshake(
		context.Background(),
		srv.URL,
		WithCapabilities(map[string]any{"browserName": "firefox"}),
	)
	if err != nil {
		t.Fatalf("Handshake: %v", err)
	}

	t.Cleanup(func() {
		if err := result.Transport.Close(); err != nil {
			t.Errorf("close transport: %v", err)
		}
	})

	if result.SessionID != "s1" {
		t.Errorf("SessionID = %q, want %q", result.SessionID, "s1")
	}

	if !strings.HasSuffix(result.WebSocketURL, "/ws") {
		t.Errorf("WebSocketURL = %q, want suffix /ws", result.WebSocketURL)
	}

	mu.Lock()
	got := string(body)
	mu.Unlock()

	if !strings.Contains(got, `"browserName":"firefox"`) {
		t.Errorf("request body %q missing capability", got)
	}

	if !strings.Contains(got, `"webSocketUrl":true`) {
		t.Errorf("request body %q missing webSocketUrl", got)
	}
}

func TestHandshakeNoWebSocketURL(t *testing.T) {
	srv := newHandshakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
		payload := map[string]any{
			"value": map[string]any{"sessionId": "s1", "capabilities": map[string]any{}},
		}

		w.Header().Set(_headerContentType, _contentTypeJSON)

		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("encode response: %v", err)
		}
	})

	if _, err := Handshake(context.Background(), srv.URL); err == nil {
		t.Fatal("Handshake: expected error")
	}
}

func TestHandshakeStatusError(t *testing.T) {
	srv := newHandshakeServer(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	})

	if _, err := Handshake(context.Background(), srv.URL); err == nil {
		t.Fatal("Handshake: expected error")
	}
}

func TestHandshakeWithHTTPClient(t *testing.T) {
	sentinel := errors.New("custom client used")

	_, err := Handshake(
		context.Background(),
		"http://127.0.0.1:1",
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, sentinel
			}),
		}),
	)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Handshake error = %v, want %v", err, sentinel)
	}
}

func TestHTTPToWSURL(t *testing.T) {
	for _, tc := range []struct {
		in, want string
	}{
		{in: "http://127.0.0.1:9222", want: "ws://127.0.0.1:9222/session"},
		{in: "https://127.0.0.1:9222", want: "wss://127.0.0.1:9222/session"},
		{in: "ws://127.0.0.1:9222/session", want: "ws://127.0.0.1:9222/session"},
		{in: "wss://127.0.0.1:9222/session", want: "wss://127.0.0.1:9222/session"},
		{in: "http://host:8080/", want: "ws://host:8080/session"},
		{in: "http://127.0.0.1:9222/anything", want: "ws://127.0.0.1:9222/session"},
	} {
		got := httpToWSURL(tc.in)
		if got != tc.want {
			t.Errorf("httpToWSURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestConnectBiDiCreatesSession(t *testing.T) {
	srv := newDirectBiDiServer(t)

	browser, err := ConnectBiDi(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("ConnectBiDi: %v", err)
	}
	defer func() { _ = browser.Close(context.Background()) }()

	if browser.Session().ID() == "" {
		t.Error("session ID is empty")
	}
}

func TestConnectEndpointFallsBackToDirectBiDi(t *testing.T) {
	srv := newDirectBiDiServer(t)

	browser, err := ConnectEndpoint(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("ConnectEndpoint: %v", err)
	}
	defer func() { _ = browser.Close(context.Background()) }()

	if browser.Session().ID() == "" {
		t.Error("session ID is empty")
	}
}

// newDirectBiDiServer returns a test server that only speaks the direct
// BiDi WebSocket flow: POST /session returns 400 (classic unavailable),
// while GET /session upgrades to a BiDi WebSocket that accepts session.new.
func newDirectBiDiServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc(_sessionPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Header.Get("Upgrade") != "" {
			conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
				Subprotocols: []string{transport.BidiSubprotocol},
			})
			if err != nil {
				return
			}

			defer func() { _ = conn.CloseNow() }()

			for {
				_, data, err := conn.Read(r.Context())
				if err != nil {
					return
				}

				var cmd struct {
					ID     int64  `json:"id"`
					Method string `json:"method"`
				}
				if err := json.Unmarshal(data, &cmd); err != nil {
					continue
				}

				var result any

				switch cmd.Method {
				case "session.status":
					result = map[string]any{"ready": true, "message": ""}
				case "session.new":
					result = map[string]any{
						"sessionId":    "direct-s1",
						"capabilities": map[string]any{"browserName": "firefox"},
					}
				case "session.end":
					result = map[string]any{}
				}

				resp, _ := json.Marshal(map[string]any{
					"type":   "success",
					"id":     cmd.ID,
					"result": result,
				})

				_ = conn.Write(r.Context(), websocket.MessageText, resp)
			}
		}

		http.Error(w, "The handshake request must use GET method", http.StatusBadRequest)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}
