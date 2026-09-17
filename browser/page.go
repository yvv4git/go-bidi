package browser

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/yvv4git/go-bidi/protocol"
)

// Page is a handle to a single browsing context. Its methods target the
// context id it was created for.
type Page struct {
	caller Caller
	id     string
}

// NewPage returns a page handle for the given browsing context id.
func NewPage(caller Caller, id string) *Page {
	return &Page{caller: caller, id: id}
}

// ID returns the underlying browsing context id.
func (p *Page) ID() string {
	return p.id
}

// Navigate navigates the page to url and waits for the document to finish
// loading.
func (p *Page) Navigate(
	ctx context.Context,
	url string,
) (*protocol.BrowsingContextNavigateResult, error) {
	params := protocol.BrowsingContextNavigateParams{
		Context: p.id,
		URL:     url,
		Wait:    protocol.ReadinessComplete,
	}

	var result protocol.BrowsingContextNavigateResult
	if err := p.caller.Call(ctx, protocol.BrowsingContextNavigate, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Reload reloads the page and waits for the document to finish loading.
func (p *Page) Reload(ctx context.Context) (*protocol.BrowsingContextReloadResult, error) {
	params := protocol.BrowsingContextReloadParams{
		Context: p.id,
		Wait:    protocol.ReadinessComplete,
	}

	var result protocol.BrowsingContextReloadResult
	if err := p.caller.Call(ctx, protocol.BrowsingContextReload, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Back moves the page one entry back in the session history.
func (p *Page) Back(ctx context.Context) error {
	return p.traverseHistory(ctx, -1)
}

// Forward moves the page one entry forward in the session history.
func (p *Page) Forward(ctx context.Context) error {
	return p.traverseHistory(ctx, 1)
}

// Screenshot captures the page. When full is set the whole document is
// captured; otherwise only the viewport.
func (p *Page) Screenshot(ctx context.Context, full bool) ([]byte, error) {
	params := protocol.BrowsingContextCaptureScreenshotParams{
		Context:               p.id,
		Origin:                protocol.OriginViewport,
		CaptureBeyondViewport: full,
	}

	var result protocol.BrowsingContextCaptureScreenshotResult
	if err := p.caller.Call(
		ctx, protocol.BrowsingContextCaptureScreenshot, params, &result,
	); err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(result.Data)
	if err != nil {
		return nil, fmt.Errorf("decode screenshot: %w", err)
	}

	return data, nil
}

// SetViewport resizes the page viewport.
func (p *Page) SetViewport(ctx context.Context, width, height int) error {
	params := protocol.BrowsingContextSetViewportParams{
		Context:  p.id,
		Viewport: &protocol.Viewport{Width: width, Height: height},
	}

	return p.caller.Call(ctx, protocol.BrowsingContextSetViewport, params, nil)
}

// Evaluate runs a JavaScript expression in the page context. The promise is
// awaited when the expression returns one.
func (p *Page) Evaluate(
	ctx context.Context,
	expression string,
) (*protocol.ScriptEvaluateResult, error) {
	params := protocol.ScriptEvaluateParams{
		Expression:      expression,
		Target:          protocol.ScriptTarget{Context: p.id},
		ResultOwnership: protocol.ResultOwnershipNone,
		AwaitPromise:    true,
	}

	var result protocol.ScriptEvaluateResult
	if err := p.caller.Call(ctx, protocol.ScriptEvaluate, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CallFunction calls a JavaScript function in the page context, passing the
// given arguments.
func (p *Page) CallFunction(
	ctx context.Context,
	function string,
	args ...protocol.ScriptArgument,
) (*protocol.ScriptEvaluateResult, error) {
	arguments := append([]protocol.ScriptArgument(nil), args...)
	if arguments == nil {
		arguments = []protocol.ScriptArgument{}
	}

	params := protocol.ScriptCallFunctionParams{
		FunctionDeclaration: function,
		Arguments:           arguments,
		Target:              protocol.ScriptTarget{Context: p.id},
		ResultOwnership:     protocol.ResultOwnershipNone,
		AwaitPromise:        true,
	}

	var result protocol.ScriptEvaluateResult
	if err := p.caller.Call(ctx, protocol.ScriptCallFunction, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Close closes the browsing context.
func (p *Page) Close(ctx context.Context) error {
	params := protocol.BrowsingContextCloseParams{Context: p.id}

	return p.caller.Call(ctx, protocol.BrowsingContextClose, params, nil)
}

func (p *Page) traverseHistory(ctx context.Context, delta int) error {
	params := protocol.BrowsingContextTraverseHistoryParams{Context: p.id, Delta: delta}

	return p.caller.Call(ctx, protocol.BrowsingContextTraverseHistory, params, nil)
}
