// Command network subscribes to network events and prints them while
// navigating to a website.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yvv4git/go-bidi/examples/internal/support"
	"github.com/yvv4git/go-bidi/protocol"
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

		sub, err := b.Client().Subscribe(ctx, []string{
			protocol.NetworkBeforeRequestSent,
			protocol.NetworkResponseStarted,
			protocol.NetworkResponseCompleted,
		})
		support.Must(err)

		defer func() { _ = sub.Close(ctx) }()

		page, err := b.NewPage(ctx)
		support.Must(err)

		nav, err := page.Navigate(ctx, "https://www.wikipedia.org/")
		support.Must(err)
		fmt.Printf("loaded %s\n\n", nav.URL)

		deadline := time.After(5 * time.Second)

		for {
			select {
			case event, ok := <-sub.Receive():
				if !ok {
					return
				}

				printEvent(event)

			case <-deadline:
				return

			case <-ctx.Done():
				return
			}
		}
	})
}

func printEvent(event protocol.Event) {
	switch event.Method {
	case protocol.NetworkBeforeRequestSent:
		var p protocol.NetworkBeforeRequestSentParams
		if err := json.Unmarshal(event.Params, &p); err != nil {
			log.Printf("decode %s: %v", event.Method, err)
			return
		}

		fmt.Printf("%-28s %s %s\n", event.Method, p.Request.Method, p.Request.URL)

	case protocol.NetworkResponseStarted:
		var p protocol.NetworkResponseStartedParams
		if err := json.Unmarshal(event.Params, &p); err != nil {
			log.Printf("decode %s: %v", event.Method, err)
			return
		}

		fmt.Printf("%-28s status=%d %s\n", event.Method, p.Response.Status, p.Request.URL)

	case protocol.NetworkResponseCompleted:
		var p protocol.NetworkResponseCompletedParams
		if err := json.Unmarshal(event.Params, &p); err != nil {
			log.Printf("decode %s: %v", event.Method, err)
			return
		}

		fmt.Printf("%-28s status=%d %s\n", event.Method, p.Response.Status, p.Request.URL)

	default:
		fmt.Printf("%-28s %s\n", event.Method, event.Params)
	}
}
