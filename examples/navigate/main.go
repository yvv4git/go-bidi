// Command navigate opens a tab and navigates to Wikipedia.
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/yvv4git/go-bidi/examples/internal/support"
)

func main() {
	opts := support.Flag(flag.CommandLine)

	flag.Parse()

	support.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		firefox, b, err := support.Connect(ctx, opts)
		support.Must(err)

		defer support.Cleanup(b, firefox)

		page, err := b.NewPage(ctx)
		support.Must(err)

		nav, err := page.Navigate(ctx, "https://www.wikipedia.org/")
		support.Must(err)
		fmt.Printf("loaded %s\n", nav.URL)
		fmt.Println("waiting 3 seconds...")
		time.Sleep(3 * time.Second)

		_ = page.Close(ctx)
	})
}
