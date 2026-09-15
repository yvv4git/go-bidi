// Package transport provides the low-level channels between the client and
// a BiDi-enabled browser: a WebSocket or an inherited OS pipe pair
// (fd 3/fd 4).
//
// A Transport carries opaque JSON frames in both directions. Message
// semantics live in the protocol package; transports only move bytes.
package transport
