package browser

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/yvv4git/go-bidi/protocol"
)

// CreateBrowsingContext creates a new browsing context.
func (s *Session) CreateBrowsingContext(
	ctx context.Context,
	params protocol.BrowsingContextCreateParams,
) (*protocol.BrowsingContextCreateResult, error) {
	var result protocol.BrowsingContextCreateResult
	if err := s.caller.Call(ctx, protocol.BrowsingContextCreate, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Navigate navigates a browsing context to a URL.
func (s *Session) Navigate(
	ctx context.Context,
	params protocol.BrowsingContextNavigateParams,
) (*protocol.BrowsingContextNavigateResult, error) {
	var result protocol.BrowsingContextNavigateResult
	if err := s.caller.Call(ctx, protocol.BrowsingContextNavigate, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTree returns the tree of browsing contexts.
func (s *Session) GetTree(
	ctx context.Context,
	params protocol.BrowsingContextGetTreeParams,
) (*protocol.BrowsingContextGetTreeResult, error) {
	var result protocol.BrowsingContextGetTreeResult
	if err := s.caller.Call(ctx, protocol.BrowsingContextGetTree, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CloseBrowsingContext closes a browsing context.
func (s *Session) CloseBrowsingContext(
	ctx context.Context,
	params protocol.BrowsingContextCloseParams,
) error {
	return s.caller.Call(ctx, protocol.BrowsingContextClose, params, nil)
}

// Reload reloads a browsing context.
func (s *Session) Reload(
	ctx context.Context,
	params protocol.BrowsingContextReloadParams,
) (*protocol.BrowsingContextReloadResult, error) {
	var result protocol.BrowsingContextReloadResult
	if err := s.caller.Call(ctx, protocol.BrowsingContextReload, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// TraverseHistory moves a browsing context through the session history.
func (s *Session) TraverseHistory(
	ctx context.Context,
	params protocol.BrowsingContextTraverseHistoryParams,
) error {
	return s.caller.Call(ctx, protocol.BrowsingContextTraverseHistory, params, nil)
}

// CaptureScreenshot captures a screenshot of a browsing context. Data is
// returned as decoded PNG bytes.
func (s *Session) CaptureScreenshot(
	ctx context.Context,
	params protocol.BrowsingContextCaptureScreenshotParams,
) ([]byte, error) {
	var result protocol.BrowsingContextCaptureScreenshotResult
	if err := s.caller.Call(
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

// SetViewport sets the viewport size of a browsing context.
func (s *Session) SetViewport(
	ctx context.Context,
	params protocol.BrowsingContextSetViewportParams,
) error {
	return s.caller.Call(ctx, protocol.BrowsingContextSetViewport, params, nil)
}
