// Command screenshot captures a page twice: the visible viewport and the
// whole document.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/yvv4git/go-bidi"
	"github.com/yvv4git/go-bidi/examples/internal/support"
)

const (
	viewportWidth  = 800
	viewportHeight = 400
)

func main() {
	opts := support.Flag(flag.CommandLine)
	dir := flag.String("out", ".", "directory to write the captures to")

	flag.Parse()

	support.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		firefox, b, err := support.Connect(ctx, opts)
		support.Must(err)

		defer support.Cleanup(b, firefox)

		page, err := b.NewPage(ctx)
		support.Must(err)

		defer func() { _ = page.Close(ctx) }()

		support.Must(page.SetViewport(ctx, viewportWidth, viewportHeight))

		_, err = page.Navigate(ctx, "https://example.com/")
		support.Must(err)

		support.Must(capture(ctx, page, filepath.Join(*dir, "viewport.png"), false))
		support.Must(capture(ctx, page, filepath.Join(*dir, "full.png"), true))
	})
}

// capture screenshots the page and writes the PNG to path, reporting the
// decoded image dimensions.
func capture(ctx context.Context, page *bidi.Page, path string, full bool) error {
	data, err := page.Screenshot(ctx, full)
	if err != nil {
		return err
	}

	config, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}

	kind := "viewport"
	if full {
		kind = "full document"
	}

	fmt.Printf("%-13s %dx%d -> %s\n", kind, config.Width, config.Height, path)

	return nil
}
