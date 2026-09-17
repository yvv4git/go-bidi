package browser

import "context"

// Caller sends a BiDi command and decodes its result. It is the port the
// session module depends on, so the message router can be plugged in later
// without coupling the two.
type Caller interface {
	Call(ctx context.Context, method string, params, result any) error
}
