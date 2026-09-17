package browser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yvv4git/go-bidi/protocol"
	"github.com/yvv4git/go-bidi/transport"
)

type fakeTransport struct {
	sendMu sync.Mutex
	sent   [][]byte
	sentCh chan []byte

	recvCh    chan []byte
	closed    chan struct{}
	closeOnce sync.Once
}

func newFakeTransport() *fakeTransport {
	return &fakeTransport{
		sentCh: make(chan []byte, 64),
		recvCh: make(chan []byte, 64),
		closed: make(chan struct{}),
	}
}

func (f *fakeTransport) Send(_ context.Context, data []byte) error {
	f.sendMu.Lock()
	defer f.sendMu.Unlock()

	select {
	case <-f.closed:
		return transport.ErrClosed
	default:
	}

	f.sent = append(f.sent, data)
	f.sentCh <- data

	return nil
}

func (f *fakeTransport) Receive(ctx context.Context) ([]byte, error) {
	select {
	case data := <-f.recvCh:
		return data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-f.closed:
		return nil, transport.ErrClosed
	}
}

func (f *fakeTransport) Close() error {
	f.closeOnce.Do(func() {
		close(f.closed)
	})

	return nil
}

func (f *fakeTransport) frames(t *testing.T) []string {
	t.Helper()

	f.sendMu.Lock()
	defer f.sendMu.Unlock()

	frames := make([]string, 0, len(f.sent))
	for _, data := range f.sent {
		frames = append(frames, string(data))
	}

	return frames
}

func respondSuccess(t *testing.T, tr *fakeTransport, result string) {
	t.Helper()

	frame := <-tr.sentCh

	var cmd struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(frame, &cmd); err != nil {
		t.Errorf("decode command: %v", err)

		return
	}

	tr.recvCh <- []byte(fmt.Sprintf(`{"type":"success","id":%d,"result":%s}`, cmd.ID, result))
}

func respondError(t *testing.T, tr *fakeTransport, code, message string) {
	t.Helper()

	frame := <-tr.sentCh

	var cmd struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(frame, &cmd); err != nil {
		t.Errorf("decode command: %v", err)

		return
	}

	tr.recvCh <- []byte(fmt.Sprintf(
		`{"type":"error","id":%d,"error":%q,"message":%q,"stacktrace":""}`,
		cmd.ID, code, message,
	))
}

func TestCallRoundTrip(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	go respondSuccess(t, tr, `{"ok":true}`)

	var out struct {
		OK bool `json:"ok"`
	}
	if err := client.Call(context.Background(), "test.method", nil, &out); err != nil {
		t.Fatalf("Call: %v", err)
	}

	if !out.OK {
		t.Error("out.OK is false, want true")
	}
}

func TestCallNoResult(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	go respondSuccess(t, tr, `{}`)
	if err := client.Call(context.Background(), "test.noResult", nil, nil); err != nil {
		t.Fatalf("Call: %v", err)
	}
}

func TestCallError(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	go respondError(t, tr, "unknown command", "no such command")

	err := client.Call(context.Background(), "test.bad", nil, nil)
	if err == nil {
		t.Fatal("Call: expected error")
	}

	var perr *protocol.Error
	if !errors.As(err, &perr) {
		t.Fatalf("error type = %T, want *protocol.Error", err)
	}

	if perr.Code != "unknown command" {
		t.Errorf("Code = %q, want %q", perr.Code, "unknown command")
	}
}

func TestCallTimeout(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
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

func TestCallMonotonicIDs(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		respondSuccess(t, tr, `{}`)
		respondSuccess(t, tr, `{}`)
	}()

	for i := 0; i < 2; i++ {
		if err := client.Call(context.Background(), "test.method", nil, nil); err != nil {
			t.Fatalf("Call %d: %v", i, err)
		}
	}
	<-done

	frames := tr.frames(t)
	if len(frames) != 2 {
		t.Fatalf("sent frames = %d, want 2", len(frames))
	}

	for i, frame := range frames {
		var cmd struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal([]byte(frame), &cmd); err != nil {
			t.Fatalf("decode frame %d: %v", i, err)
		}

		if want := int64(i + 1); cmd.ID != want {
			t.Errorf("frame %d ID = %d, want %d", i, cmd.ID, want)
		}
	}
}

func TestCallAfterClose(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)

	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}

	if err := client.Call(context.Background(), "test.method", nil, nil); !errors.Is(err, ErrClosed) {
		t.Fatalf("Call error = %v, want ErrClosed", err)
	}
}

func TestCallCloseFailsPending(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Call(context.Background(), "test.pending", nil, nil)
	}()

	<-tr.sentCh

	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if err := <-errCh; !errors.Is(err, ErrClosed) {
		t.Fatalf("Call error = %v, want ErrClosed", err)
	}
}
