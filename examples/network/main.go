// Command network intercepts the requests made by the page and rewrites the
// document request to a different URL, proving the modification by
// evaluating the loaded body.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/yvv4git/go-bidi"
	"github.com/yvv4git/go-bidi/examples/internal/support"
	"github.com/yvv4git/go-bidi/protocol"
)

func main() {
	opts := support.Flag(flag.CommandLine)

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	server := httptest.NewServer(documentHandlers())
	defer server.Close()

	firefox, b, err := support.Connect(ctx, opts)
	support.Must(err)

	defer func() {
		_ = b.Close(ctx)
		if firefox != nil {
			_ = firefox.Close()
		}
	}()

	sub, err := b.Client().Subscribe(ctx, []string{
		protocol.NetworkBeforeRequestSent,
		protocol.NetworkResponseStarted,
		protocol.NetworkResponseCompleted,
	})
	support.Must(err)

	intercept, err := b.Session().AddIntercept(ctx, protocol.NetworkAddInterceptParams{
		Phases: []string{protocol.InterceptPhaseBeforeRequestSent},
	})
	support.Must(err)

	target := server.URL + "/page"
	redirect := server.URL + "/moved"

	var wg sync.WaitGroup
	wg.Go(func() {
		redirectRequests(ctx, b.Session(), sub, target, redirect)
	})

	page, err := b.NewPage(ctx)
	support.Must(err)

	nav, err := page.Navigate(ctx, target)
	support.Must(err)
	fmt.Printf("requested %s, loaded %s\n", target, nav.URL)

	wg.Wait()

	result, err := page.Evaluate(ctx, "document.body.textContent")
	support.Must(err)
	body, err := result.Result.Interface()
	support.Must(err)
	fmt.Printf("document body: %q\n", body)

	support.Must(b.Session().RemoveIntercept(
		ctx, protocol.NetworkRemoveInterceptParams{Intercept: intercept},
	))
	support.Must(sub.Close(ctx))
}

// documentHandlers serves the original and the redirect target documents.
func documentHandlers() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/page", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `<!doctype html><p>the original document</p>`)
	})
	mux.HandleFunc("/moved", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `<!doctype html><p>the redirected document</p>`)
	})

	return mux
}

// redirectRequests drains sub and continues every blocked request,
// rewriting the URL of the document request to redirect.
func redirectRequests(
	ctx context.Context,
	session *bidi.Session,
	sub *bidi.Subscription,
	target,
	redirect string,
) {
	for {
		select {
		case event, ok := <-sub.Receive():
			if !ok {
				return
			}

			continueBlocked(ctx, session, event, target, redirect)

		case <-ctx.Done():
			return
		}
	}
}

// continueBlocked continues an intercepted request, printing a summary and
// rewriting the URL of the document request to redirect.
func continueBlocked(
	ctx context.Context,
	session *bidi.Session,
	event protocol.Event,
	target,
	redirect string,
) {
	switch event.Method {
	case protocol.NetworkBeforeRequestSent:
		var p protocol.NetworkBeforeRequestSentParams
		if err := json.Unmarshal(event.Params, &p); err != nil {
			log.Printf("decode %s: %v", event.Method, err)

			return
		}

		fmt.Printf("%-28s %s %s\n", event.Method, p.Request.Method, p.Request.URL)

		params := protocol.NetworkContinueRequestParams{Request: p.Request.Request}
		if p.Request.URL == target {
			params.URL = redirect
		}

		if err := session.ContinueRequest(ctx, params); err != nil {
			log.Printf("continue: %v", err)
		}

	case protocol.NetworkResponseStarted, protocol.NetworkResponseCompleted:
		var p protocol.NetworkResponseStartedParams
		if err := json.Unmarshal(event.Params, &p); err != nil {
			return
		}

		fmt.Printf("%-28s status=%d %s\n", event.Method, p.Response.Status, p.Request.URL)
	}
}
