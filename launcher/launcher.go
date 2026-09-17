// Package launcher spawns BiDi-enabled browser instances.
//
// It currently targets Firefox but any browser exposing a BiDi WebSocket
// endpoint or pipe transport can be driven through the same Option set.
package launcher

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/yvv4git/go-bidi/transport"
)

// Option customizes the Firefox launcher.
type Option func(*config)

type config struct {
	execPath string
	headless bool
	timeout  time.Duration
	args     []string
	profile  string
	pipe     bool
	readyFn  func(ctx context.Context, endpoint *url.URL) error
	execFn   func(ctx context.Context, name string, args ...string) *exec.Cmd
}

const (
	_defaultTimeout = 30 * time.Second
	_defaultHost    = "127.0.0.1"
	_pollInterval   = 50 * time.Millisecond
	_dialTimeout    = 500 * time.Millisecond
)

// WithExecPath sets the path to the browser binary. The default resolves
// "firefox" via exec.LookPath.
func WithExecPath(path string) Option {
	return func(c *config) {
		c.execPath = path
	}
}

// WithHeadless starts the browser in headless mode.
func WithHeadless(headless bool) Option {
	return func(c *config) {
		c.headless = headless
	}
}

// WithTimeout sets the maximum wait for browser readiness. The default is
// 30 seconds.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}

// WithArgs adds extra command-line arguments passed to the browser.
func WithArgs(args ...string) Option {
	return func(c *config) {
		c.args = append(c.args, args...)
	}
}

// WithProfile sets the Firefox profile directory (-profile).
func WithProfile(path string) Option {
	return func(c *config) {
		c.profile = path
	}
}

// WithPipe enables pipe-based connection mode (fd 3/fd 4) instead of a
// WebSocket. Supported by browsers implementing the WebDriver BiDi pipe
// transport.
func WithPipe(pipe bool) Option {
	return func(c *config) {
		c.pipe = pipe
	}
}

// Firefox is a handle to a launched browser process. The process is killed
// by Close or by cancellation of the context passed to Launch.
type Firefox struct {
	cmd       *exec.Cmd
	endpoint  *url.URL
	port      int
	pipe      bool
	tr        transport.Transport
	readyFn   func(ctx context.Context, endpoint *url.URL) error
	cleanupFn func()

	mu   sync.Mutex
	done chan struct{}
}

// Launch starts the browser process and waits until it is ready. In pipe
// mode readiness is immediate. The returned handle owns the process.
func Launch(ctx context.Context, opts ...Option) (*Firefox, error) {
	cfg := config{
		timeout: _defaultTimeout,
		execFn:  exec.CommandContext,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	f, extraFiles, err := cfg.init(ctx)
	if err != nil {
		return nil, err
	}

	if err := cfg.startCmd(ctx, f, extraFiles); err != nil {
		return nil, err
	}

	if cfg.pipe {
		return f, nil
	}

	if err := f.waitForReadiness(ctx, cfg.timeout); err != nil {
		f.kill()
		<-f.done

		return nil, err
	}

	return f, nil
}

// Endpoint returns the WebSocket URL of the browser BiDi endpoint. It is
// nil in pipe mode.
func (f *Firefox) Endpoint() *url.URL {
	return f.endpoint
}

// Port returns the debugging port. It is 0 in pipe mode.
func (f *Firefox) Port() int {
	return f.port
}

// Transport returns the browser transport. In pipe mode it is a
// transport.PipeTransport; in WebSocket mode it is nil and the caller
// should use ConnectEndpoint with Endpoint.
func (f *Firefox) Transport() transport.Transport {
	return f.tr
}

// Close kills the browser process and waits for it to exit.
func (f *Firefox) Close() error {
	f.kill()

	if f.cmd != nil {
		<-f.done
	}

	if f.cleanupFn != nil {
		f.cleanupFn()
	}

	return nil
}

func (f *Firefox) kill() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.cmd != nil && f.cmd.Process != nil {
		_ = f.cmd.Process.Kill()
	}
}

func (f *Firefox) watchProcess() {
	go func() {
		_ = f.cmd.Wait()
		close(f.done)
	}()
}

func (f *Firefox) waitForReadiness(ctx context.Context, timeout time.Duration) error {
	readyFn := f.readyFn
	if readyFn == nil {
		readyFn = f.httpReady
	}

	return f.waitReady(ctx, readyFn, timeout)
}

func (f *Firefox) waitReady(
	ctx context.Context,
	fn func(context.Context, *url.URL) error,
	timeout time.Duration,
) error {
	ticker := time.NewTicker(_pollInterval)
	defer ticker.Stop()

	deadline := time.After(timeout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline:
			return fmt.Errorf("timeout waiting for browser readiness")
		case <-f.done:
			return fmt.Errorf("browser process exited unexpectedly")
		case <-ticker.C:
			if err := fn(ctx, f.endpoint); err == nil {
				return nil
			}
		}
	}
}

func (f *Firefox) httpReady(_ context.Context, endpoint *url.URL) error {
	conn, err := net.DialTimeout("tcp", endpoint.Host, _dialTimeout)
	if err != nil {
		return err
	}

	_ = conn.Close()

	return nil
}

func (f *Firefox) setupPipe() (transport.Transport, []*os.File, error) {
	childW, parentR, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("create read pipe: %w", err)
	}

	parentW, childR, err := os.Pipe()
	if err != nil {
		_ = childW.Close()
		_ = parentR.Close()

		return nil, nil, fmt.Errorf("create write pipe: %w", err)
	}

	tr := transport.NewPipe(parentR, parentW)
	extraFiles := []*os.File{childW, childR}

	return tr, extraFiles, nil
}

func (cfg *config) resolveExecPath() error {
	if cfg.execPath != "" {
		return nil
	}

	// Try exec.LookPath first (works on Linux and if Firefox is in PATH)
	if path, err := exec.LookPath("firefox"); err == nil {
		cfg.execPath = path
		return nil
	}

	// Try common macOS application paths
	if path := findFile(
		"/Applications/Firefox.app/Contents/MacOS/firefox",
		os.Getenv("HOME")+"/Applications/Firefox.app/Contents/MacOS/firefox",
	); path != "" {
		cfg.execPath = path
		return nil
	}

	// Try common Linux binary paths
	if path := findFile(
		"/usr/bin/firefox",
		"/usr/local/bin/firefox",
		"/opt/firefox/firefox",
	); path != "" {
		cfg.execPath = path
		return nil
	}

	return fmt.Errorf("find firefox: firefox executable not found in PATH or standard application directories")
}

// findFile returns the first existing file path from the given list.
func findFile(paths ...string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func (cfg *config) init(ctx context.Context) (*Firefox, []*os.File, error) {
	if err := cfg.resolveExecPath(); err != nil {
		return nil, nil, err
	}

	// Create a temporary profile if none was specified. This isolates
	// the launched browser from any already-running Firefox instance.
	var cleanupFn func()

	if cfg.profile == "" {
		tmpDir, err := os.MkdirTemp("", "go-bidi-firefox-*")
		if err != nil {
			return nil, nil, fmt.Errorf("create temp profile: %w", err)
		}

		cfg.profile = tmpDir
		cleanupFn = func() { _ = os.RemoveAll(tmpDir) }
	}

	f := &Firefox{
		pipe:      cfg.pipe,
		done:      make(chan struct{}),
		readyFn:   cfg.readyFn,
		cleanupFn: cleanupFn,
	}

	extraFiles, err := cfg.setupTransport(ctx, f)
	if err != nil {
		return nil, nil, err
	}

	return f, extraFiles, nil
}

func (cfg *config) setupTransport(_ context.Context, f *Firefox) ([]*os.File, error) {
	if cfg.pipe {
		tr, ef, err := f.setupPipe()
		if err != nil {
			return nil, err
		}

		f.tr = tr

		return ef, nil
	}

	if err := cfg.allocPort(f); err != nil {
		return nil, err
	}

	return nil, nil
}

func (cfg *config) allocPort(f *Firefox) error {
	ln, err := net.Listen("tcp", _defaultHost+":0")
	if err != nil {
		return fmt.Errorf("allocate port: %w", err)
	}

	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	f.port = port
	f.endpoint = &url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", _defaultHost, port),
		Path:   "/session",
	}

	return nil
}

func (cfg *config) buildArgs(f *Firefox) []string {
	cmdArgs := make([]string, 0, 4+len(cfg.args))

	if cfg.headless {
		cmdArgs = append(cmdArgs, "--headless")
	}

	cmdArgs = append(cmdArgs, "--no-remote")

	if cfg.profile != "" {
		cmdArgs = append(cmdArgs, "-profile", cfg.profile)
	}

	if f.pipe {
		cmdArgs = append(cmdArgs, "--remote-debugging-pipe")
	} else {
		cmdArgs = append(cmdArgs,
			"--remote-debugging-port", strconv.Itoa(f.port))
	}

	cmdArgs = append(cmdArgs, cfg.args...)

	return cmdArgs
}

func (cfg *config) startCmd(ctx context.Context, f *Firefox, extraFiles []*os.File) error {
	args := cfg.buildArgs(f)
	f.cmd = cfg.execFn(ctx, cfg.execPath, args...)
	f.cmd.ExtraFiles = extraFiles

	err := f.cmd.Start()

	for _, fd := range extraFiles {
		_ = fd.Close()
	}

	if err != nil {
		return fmt.Errorf("start firefox: %w", err)
	}

	f.watchProcess()

	return nil
}
