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

// Connect establishes a session over o.Endpoint, or over a freshly
// launched headless Firefox when that is empty.
func Connect(ctx context.Context, o *Options) (*bidi.Firefox, *bidi.Browser, error) {
	addr, firefox, err := target(ctx, o)
	if err != nil {
		return nil, nil, err
	}

	b, err := bidi.ConnectEndpoint(ctx, addr, bidi.WithTimeout(o.Timeout))
	if err != nil {
		if firefox != nil {
			_ = firefox.Close()
		}

		return nil, nil, fmt.Errorf("connect: %w", err)
	}

	return firefox, b, nil
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

// Must logs the error and exits when err is non-nil.
func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
