// Package support contains the connection and error-handling helpers shared
// by the runnable example programs in the sibling directories below.
package support

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yvv4git/go-bidi"
	"github.com/yvv4git/go-bidi/protocol"
)

// Options holds the command-line flags shared by every example.
type Options struct {
	Endpoint string
	ExecPath string
	Timeout  time.Duration
}

// Flag registers the shared -endpoint, -exec and -timeout flags on fs and
// returns the Options they fill. Callers register topic-specific flags first
// and then call flag.Parse.
func Flag(fs *flag.FlagSet) *Options {
	o := &Options{}
	fs.StringVar(&o.Endpoint, "endpoint", "", "connect to this WebDriver BiDi endpoint")
	fs.StringVar(&o.ExecPath, "exec", "", "path to the Firefox binary to launch")
	fs.DurationVar(&o.Timeout, "timeout", 2*time.Minute, "overall command timeout")

	return o
}

// Browser is a connected BiDi session together with the client that owns
// its transport. It is created by Connect and shut down with Close.
type Browser struct {
	client  *bidi.Client
	session *bidi.Session
	firefox *bidi.Firefox
	ctx     context.Context
}

// Connect establishes a session over o.Endpoint, or over a freshly
// launched headless Firefox when that is empty.
func Connect(ctx context.Context, o *Options) (*Browser, error) {
	addr, firefox, err := target(ctx, o)
	if err != nil {
		return nil, err
	}

	result, err := bidi.Handshake(ctx, addr, bidi.WithTimeout(o.Timeout))
	if err != nil {
		if firefox != nil {
			_ = firefox.Close()
		}

		return nil, fmt.Errorf("handshake: %w", err)
	}

	client := bidi.NewClient(result.Transport, bidi.WithTimeout(o.Timeout))

	return &Browser{
		client:  client,
		session: bidi.NewSession(client, result.SessionID),
		firefox: firefox,
		ctx:     ctx,
	}, nil
}

// target returns the endpoint address of o.Endpoint, launching a headless
// Firefox with the -exec binary when that is empty.
func target(ctx context.Context, o *Options) (string, *bidi.Firefox, error) {
	if o.Endpoint != "" {
		return o.Endpoint, nil, nil
	}

	firefox, err := bidi.Launch(ctx,
		bidi.WithLaunchExecPath(o.ExecPath),
		bidi.WithLaunchHeadless(true),
		bidi.WithLaunchTimeout(o.Timeout),
	)
	if err != nil {
		return "", nil, fmt.Errorf("launch firefox: %w", err)
	}

	return firefox.Endpoint().String(), firefox, nil
}

// Client returns the client that owns the transport.
func (b *Browser) Client() *bidi.Client {
	return b.client
}

// Session returns the session handle.
func (b *Browser) Session() *bidi.Session {
	return b.session
}

// NewPage opens a new tab and returns its page handle.
func (b *Browser) NewPage(ctx context.Context) (*bidi.Page, error) {
	result, err := b.session.CreateBrowsingContext(
		ctx,
		protocol.BrowsingContextCreateParams{Type: protocol.ContextTypeTab},
	)
	if err != nil {
		return nil, fmt.Errorf("create browsing context: %w", err)
	}

	return bidi.NewPage(b.client, result.Context), nil
}

// Close ends the session, closes the transport and stops a launched
// Firefox. It is safe to call more than once.
func (b *Browser) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = b.session.End(ctx)
	_ = b.client.Close()

	if b.firefox != nil {
		_ = b.firefox.Close()
	}
}

// Must logs the error and exits when err is non-nil.
func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
