// Command basic opens a page, navigates to a URL and saves a screenshot of
// the result.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/yvv4git/go-bidi/examples/internal/support"
)

const (
	viewportWidth  = 1024
	viewportHeight = 768
)

func main() {
	opts := support.Flag(flag.CommandLine)
	out := flag.String("out", "basic.png", "path to write the screenshot to")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	firefox, b, err := support.Connect(ctx, opts)
	support.Must(err)

	defer func() {
		_ = b.Close(ctx)
		if firefox != nil {
			_ = firefox.Close()
		}
	}()

	page, err := b.NewPage(ctx)
	support.Must(err)

	defer func() { _ = page.Close(ctx) }()

	support.Must(page.SetViewport(ctx, viewportWidth, viewportHeight))

	nav, err := page.Navigate(ctx, "https://example.com/")
	support.Must(err)
	fmt.Printf("loaded %s (navigation %s)\n", nav.URL, nav.Navigation)

	png, err := page.Screenshot(ctx, false)
	support.Must(err)
	support.Must(os.WriteFile(*out, png, 0o600))
	fmt.Printf("wrote %s (%d bytes)\n", *out, len(png))
}
