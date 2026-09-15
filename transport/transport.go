// Package transport provides the low-level channels between the client and
// a BiDi-enabled browser: a WebSocket or an inherited OS pipe pair
// (fd 3/fd 4).
//
// A Transport carries opaque JSON frames in both directions. Message
// semantics live in the protocol package; transports only move bytes.
package transport

import (
	"context"
	"errors"
)

// ErrClosed is returned when operating on a closed transport.
var ErrClosed = errors.New("transport: closed")

// Transport is the bidirectional channel between the client and the
// browser. Implementations must be safe for concurrent use.
type Transport interface {
	// Send writes a single JSON frame to the browser.
	Send(ctx context.Context, data []byte) error
	// Receive reads the next JSON frame from the browser.
	Receive(ctx context.Context) ([]byte, error)
	// Close tears down the transport and releases associated resources.
	Close() error
}
