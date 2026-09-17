// Package browser provides a client-side API for the WebDriver BiDi
// protocol: sessions, pages, input, network and log handling.
package browser

import (
	"context"

	"github.com/yvv4git/go-bidi/protocol"
)

// CloseBrowser closes the browser itself.
func (s *Session) CloseBrowser(ctx context.Context) error {
	return s.caller.Call(ctx, protocol.BrowserClose, nil, nil)
}

// CreateUserContext creates a new user context and returns its id.
func (s *Session) CreateUserContext(ctx context.Context) (string, error) {
	var result protocol.BrowserCreateUserContextResult
	if err := s.caller.Call(ctx, protocol.BrowserCreateUserContext, nil, &result); err != nil {
		return "", err
	}

	return result.UserContext, nil
}
