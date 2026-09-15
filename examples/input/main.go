// Command input demonstrates typing text, pressing keys and clicking on a
// page served from a data: URL.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"

	"github.com/yvv4git/go-bidi"
	"github.com/yvv4git/go-bidi/examples/internal/support"
	"github.com/yvv4git/go-bidi/protocol"
)

const (
	clickButtonX = 150
	clickButtonY = 30
)

const demoPage = `<!doctype html>
<style>
  button { position: fixed; left: 0; top: 0; width: 300px; height: 60px; }
</style>
<input id="name">
<button id="go">Go</button>
<p id="output">none</p>
<script>
  const input = document.getElementById('name');
  const output = document.getElementById('output');
  input.addEventListener('keydown', function (event) {
    if (event.key === 'Enter') {
      output.textContent = 'entered: ' + input.value;
    }
  });
  document.getElementById('go').addEventListener('click', function () {
    output.textContent = 'clicked with: ' + input.value;
  });
</script>`

func main() {
	opts := support.Flag(flag.CommandLine)

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	b, err := support.Connect(ctx, opts)
	support.Must(err)

	defer b.Close()

	page, err := b.NewPage(ctx)
	support.Must(err)

	defer func() { _ = page.Close(ctx) }()

	_, err = page.Navigate(ctx, "data:text/html,"+url.PathEscape(demoPage))
	support.Must(err)

	_, err = page.Evaluate(ctx, "document.getElementById('name').focus()")
	support.Must(err)
	support.Must(page.Type(ctx, "hello bidi"))
	support.Must(page.Press(ctx, "Enter"))
	fmt.Printf("after typing and Enter: %s\n", evalString(ctx, page, "output"))

	support.Must(page.Click(ctx, clickButtonX, clickButtonY))
	fmt.Printf("after click:            %s\n", evalString(ctx, page, "output"))
}

// evalString evaluates an expression in the page and returns its value as a
// string.
func evalString(ctx context.Context, page *bidi.Page, expression string) string {
	result, err := page.Evaluate(ctx, expression)
	support.Must(err)

	if result.Type == protocol.EvaluateException {
		return "exception: " + result.ExceptionDetails.Text
	}

	value, err := result.Result.Interface()
	support.Must(err)

	return fmt.Sprint(value)
}
