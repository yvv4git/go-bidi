package launcher

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"
)

const _testBrowserFlag = "--go-bidi-test-browser"

// TestMain enables the fake browser process: when the test binary is
// re-invoked with _testBrowserFlag it blocks forever, standing in for a
// real browser during process-lifecycle tests.
func TestMain(m *testing.M) {
	for _, arg := range os.Args[1:] {
		if arg == _testBrowserFlag {
			for {
				time.Sleep(time.Hour)
			}
		}
	}

	os.Exit(m.Run())
}

func alwaysReady(_ context.Context, _ *url.URL) error {
	return nil
}

func neverReady(_ context.Context, _ *url.URL) error {
	return fmt.Errorf("not ready")
}

func withReadyFn(fn func(ctx context.Context, endpoint *url.URL) error) Option {
	return func(c *config) {
		c.readyFn = fn
	}
}

func newFakeLaunch(opts ...Option) (*Firefox, context.CancelFunc) {
	base := []Option{
		WithExecPath(os.Args[0]),
		WithArgs(_testBrowserFlag),
		withReadyFn(alwaysReady),
	}

	ctx, cancel := context.WithCancel(context.Background())
	f, err := Launch(ctx, append(base, opts...)...)
	if err != nil {
		cancel()

		return nil, cancel
	}

	return f, cancel
}

func TestOptions(t *testing.T) {
	tests := []struct {
		name  string
		opts  []Option
		check func(t *testing.T, cfg *config)
	}{
		{
			name: "exec path",
			opts: []Option{WithExecPath("/custom/firefox")},
			check: func(t *testing.T, cfg *config) {
				if cfg.execPath != "/custom/firefox" {
					t.Errorf("execPath = %q, want %q", cfg.execPath, "/custom/firefox")
				}
			},
		},
		{
			name: "headless",
			opts: []Option{WithHeadless(true)},
			check: func(t *testing.T, cfg *config) {
				if !cfg.headless {
					t.Error("headless is false, want true")
				}
			},
		},
		{
			name: "timeout",
			opts: []Option{WithTimeout(5 * time.Second)},
			check: func(t *testing.T, cfg *config) {
				if cfg.timeout != 5*time.Second {
					t.Errorf("timeout = %v, want %v", cfg.timeout, 5*time.Second)
				}
			},
		},
		{
			name: "args",
			opts: []Option{WithArgs("--one", "--two")},
			check: func(t *testing.T, cfg *config) {
				if len(cfg.args) != 2 {
					t.Fatalf("args = %v, want 2 entries", cfg.args)
				}
			},
		},
		{
			name: "profile",
			opts: []Option{WithProfile("/tmp/prof")},
			check: func(t *testing.T, cfg *config) {
				if cfg.profile != "/tmp/prof" {
					t.Errorf("profile = %q, want %q", cfg.profile, "/tmp/prof")
				}
			},
		},
		{
			name: "pipe",
			opts: []Option{WithPipe(true)},
			check: func(t *testing.T, cfg *config) {
				if !cfg.pipe {
					t.Error("pipe is false, want true")
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg config
			for _, opt := range tt.opts {
				opt(&cfg)
			}

			tt.check(t, &cfg)
		})
	}
}

func TestLaunchWebSocket(t *testing.T) {
	f, cancel := newFakeLaunch(WithHeadless(true))
	if f == nil {
		t.Fatal("newFakeLaunch returned nil after Launch error")
	}
	defer cancel()

	if f.Endpoint() == nil {
		t.Fatal("Endpoint is nil in WebSocket mode")
	}

	if f.Endpoint().Scheme != "ws" {
		t.Errorf("Endpoint scheme = %q, want %q", f.Endpoint().Scheme, "ws")
	}

	if f.Endpoint().Path != "/session" {
		t.Errorf("Endpoint path = %q, want %q", f.Endpoint().Path, "/session")
	}

	if f.Port() == 0 {
		t.Error("Port is 0 in WebSocket mode")
	}

	if f.Transport() != nil {
		t.Error("Transport is not nil in WebSocket mode")
	}

	hasPort := false
	for _, arg := range f.cmd.Args[1:] {
		if arg == "--remote-debugging-port" {
			hasPort = true
		}
	}

	if !hasPort {
		t.Error("args missing --remote-debugging-port")
	}

	hasHeadless := false
	for _, arg := range f.cmd.Args[1:] {
		if arg == "--headless" {
			hasHeadless = true
		}
	}

	if !hasHeadless {
		t.Error("args missing --headless")
	}

	if err := f.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestLaunchPipe(t *testing.T) {
	f, cancel := newFakeLaunch(WithPipe(true))
	if f == nil {
		t.Fatal("newFakeLaunch returned nil after Launch error")
	}
	defer cancel()

	if f.Transport() == nil {
		t.Fatal("Transport is nil in pipe mode")
	}

	if f.Endpoint() != nil {
		t.Error("Endpoint is not nil in pipe mode")
	}

	hasPipe := false
	for _, arg := range f.cmd.Args[1:] {
		if arg == "--remote-debugging-pipe" {
			hasPipe = true
		}
	}

	if !hasPipe {
		t.Error("args missing --remote-debugging-pipe")
	}

	if err := f.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestHTTPReady(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	endpoint := &url.URL{Host: ln.Addr().String()}
	f := &Firefox{}

	if err := f.httpReady(context.Background(), endpoint); err != nil {
		t.Fatalf("httpReady on open port: %v", err)
	}

	if err := ln.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	if err := f.httpReady(context.Background(), endpoint); err == nil {
		t.Fatal("httpReady on closed port: expected error")
	}
}

func TestLaunchReadinessTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Launch(
		ctx,
		WithExecPath(os.Args[0]),
		WithArgs(_testBrowserFlag),
		WithTimeout(100*time.Millisecond),
		withReadyFn(neverReady),
	)
	if err == nil {
		t.Fatal("Launch: expected error")
	}
}

func TestLaunchExecNotFound(t *testing.T) {
	_, err := Launch(
		context.Background(),
		WithExecPath("/nonexistent/firefox-binary"),
		withReadyFn(alwaysReady),
	)
	if err == nil {
		t.Fatal("Launch: expected error")
	}
}

func TestLaunchProfileArg(t *testing.T) {
	f, cancel := newFakeLaunch(WithProfile("/tmp/prof"))
	if f == nil {
		t.Fatal("newFakeLaunch returned nil after Launch error")
	}
	defer cancel()

	hasProfile := false
	for _, arg := range f.cmd.Args[1:] {
		if arg == "-profile" {
			hasProfile = true
		}
	}

	if !hasProfile {
		t.Error("args missing -profile")
	}

	if err := f.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestLaunchContextCancel(t *testing.T) {
	f, cancel := newFakeLaunch()
	if f == nil {
		t.Fatal("newFakeLaunch returned nil after Launch error")
	}

	cancel()

	select {
	case <-f.done:
	case <-time.After(5 * time.Second):
		t.Fatal("process not killed after context cancel")
	}

	if err := f.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestCloseKillsProcess(t *testing.T) {
	f, cancel := newFakeLaunch()
	if f == nil {
		t.Fatal("newFakeLaunch returned nil after Launch error")
	}
	defer cancel()

	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	select {
	case <-f.done:
	default:
		t.Error("done not closed after Close")
	}
}
