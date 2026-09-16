// Command evaluate runs JavaScript in the page and decodes the returned
// RemoteValue into natural Go values.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"

	"github.com/yvv4git/go-bidi/examples/internal/support"
	"github.com/yvv4git/go-bidi/protocol"
)

func main() {
	opts := support.Flag(flag.CommandLine)

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

	const doc = `<!doctype html><p>a small page</p>`

	_, err = page.Navigate(ctx, "data:text/html,"+url.PathEscape(doc))
	support.Must(err)

	for _, expression := range []string{
		"1 + 2",
		"'answer: ' + 42",
		"({ answer: 42, items: [1, 2, 3] })",
		"Promise.resolve({ ready: document.readyState, url: location.href })",
	} {
		result, err := page.Evaluate(ctx, expression)
		support.Must(err)
		fmt.Printf("%-45q => %v\n", expression, decoded(result))
	}

	result, err := page.CallFunction(ctx, "function (a, b) { return a * b; }",
		protocol.ScriptArgument{Type: protocol.RemoteTypeNumber, Value: json.RawMessage("6")},
		protocol.ScriptArgument{Type: protocol.RemoteTypeNumber, Value: json.RawMessage("7")},
	)
	support.Must(err)
	fmt.Printf("call function            => %v\n", decoded(result))
}

// decoded converts an evaluation result into a printable Go value.
func decoded(result *protocol.ScriptEvaluateResult) any {
	if result.Type == protocol.EvaluateException {
		return "exception: " + result.ExceptionDetails.Text
	}

	value, err := result.Result.Interface()
	support.Must(err)

	return value
}
