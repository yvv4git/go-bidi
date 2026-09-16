// Command events subscribes to browsingContext and script events and prints
// them while a page loads in the background.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
			protocol.BrowsingContextContextCreated,
			protocol.BrowsingContextDomContentLoaded,
			protocol.BrowsingContextLoad,
			protocol.ScriptRealmCreated,
		})
		support.Must(err)

		defer func() { _ = sub.Close(ctx) }()

		page, err := b.NewPage(ctx)
		support.Must(err)

		_, err = page.Navigate(ctx, "https://example.com/")
		support.Must(err)

		deadline := time.After(3 * time.Second)

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

// printEvent decodes the params of a subscribed event into its documented
// shape and prints a one-line summary.
func printEvent(event protocol.Event) {
	switch event.Method {
	case protocol.ScriptRealmCreated:
		var realm protocol.ScriptRealmInfo
		if err := json.Unmarshal(event.Params, &realm); err != nil {
			break
		}

		fmt.Printf("%-32s realm=%s type=%s\n", event.Method, realm.Realm, realm.Type)

		return

	case protocol.BrowsingContextContextCreated:
		var info protocol.BrowsingContextInfo
		if err := json.Unmarshal(event.Params, &info); err != nil {
			break
		}

		fmt.Printf("%-32s context=%s url=%s\n", event.Method, info.Context, info.URL)

		return

	case protocol.BrowsingContextDomContentLoaded, protocol.BrowsingContextLoad:
		var nav protocol.NavigationInfo
		if err := json.Unmarshal(event.Params, &nav); err != nil {
			break
		}

		fmt.Printf("%-32s context=%s url=%s\n", event.Method, nav.Context, nav.URL)

		return
	}

	fmt.Printf("%-32s %s\n", event.Method, event.Params)
}
