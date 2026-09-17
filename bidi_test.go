package bidi

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Exported API surface referenced so the re-exports stay in sync.
var (
	_ = Connect
	_ = ConnectEndpoint
	_ = Handshake
	_ = NewClient
	_ = NewSession
	_ = NewPage
	_ = NewRequests
	_ = NewLogs
	_ = WithTimeout
	_ = WithHTTPClient
	_ = WithSubprotocol
	_ = WithCapabilities
	_ = ErrClosed
)

var (
	_ RemoteValue
	_ Transport = (*scriptedTransport)(nil)
)

type scriptedTransport struct {
	recv chan []byte

	mu sync.Mutex
	n  int
}

func newScriptedTransport() *scriptedTransport {
	return &scriptedTransport{recv: make(chan []byte, 16)}
}

func (s *scriptedTransport) Send(_ context.Context, data []byte) error {
	s.mu.Lock()
	s.n++
	nth := s.n
	s.mu.Unlock()

	go func() {
		var cmd struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(data, &cmd); err != nil {
			return
		}

		result := `{}`
		if nth == 1 {
			result = `{"sessionId":"root-s1","capabilities":{}}`
		}

		s.recv <- fmt.Appendf(nil,
			`{"type":"success","id":%d,"result":%s}`,
			cmd.ID, result,
		)
	}()

	return nil
}

func (s *scriptedTransport) Receive(ctx context.Context) ([]byte, error) {
	select {
	case data := <-s.recv:
		return data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *scriptedTransport) Close() error {
	return nil
}

func TestRootConnect(t *testing.T) {
	tr := newScriptedTransport()

	b, err := Connect(context.Background(), tr, WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	if got, want := b.Session().ID(), "root-s1"; got != want {
		t.Errorf("Session().ID() = %q, want %q", got, want)
	}

	if err := b.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
