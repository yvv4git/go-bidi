// Package bidi is a Go client for the WebDriver BiDi browser automation
// protocol. It drives BiDi-enabled browsers (Firefox Nightly and other
// WebDriver BiDi implementations) over a single full-duplex WebSocket
// channel using JSON-RPC 2.0 messages.
//
// The package is split into layers:
//
//	transport/  low-level byte channels: WebSocket and pipe (fd3/fd4)
//	protocol/   BiDi message types and JSON encoding/decoding
//	browser/    high-level browser and page control
//
// Dependencies flow one way, from high level to low level: browser/ uses
// transport/ and protocol/, while protocol/ stays independent so that
// transport moves raw frames and protocol gives them meaning.
//
// This root package re-exports the browser API so that most callers only
// need a single import.
package bidi

import (
	"github.com/yvv4git/go-bidi/browser"
	"github.com/yvv4git/go-bidi/protocol"
	"github.com/yvv4git/go-bidi/transport"
)

// Browser is a high-level handle over an established BiDi session.
type Browser = browser.Browser

// Page is a handle to a single browsing context.
type Page = browser.Page

// Session issues commands over a Caller and tracks subscriptions.
type Session = browser.Session

// Subscription delivers events for one session.subscribe call.
type Subscription = browser.Subscription

// Requests records network events for a browsing context.
type Requests = browser.Requests

// Logs records log entries for a browsing context.
type Logs = browser.Logs

// Entry summarizes one log.entryAdded event.
type Entry = browser.Entry

// HandshakeResult is the outcome of the WebDriver classic handshake.
type HandshakeResult = browser.HandshakeResult

// Option customizes Handshake, Connect, ConnectEndpoint and NewClient.
type Option = browser.Option

// Transport is the bidirectional channel between the client and the browser.
type Transport = transport.Transport

// RemoteValue is a serialized JavaScript value.
type RemoteValue = protocol.RemoteValue

// Constructors exposed through the root package.
var (
	// Connect establishes a BiDi session over an existing transport.
	Connect = browser.Connect
	// ConnectEndpoint performs the handshake and connects over HTTP.
	ConnectEndpoint = browser.ConnectEndpoint
	// Handshake creates a session over the WebDriver classic endpoint.
	Handshake = browser.Handshake
	// NewClient starts a routed client over a transport.
	NewClient = browser.NewClient
	// NewSession returns a handle for an existing session id.
	NewSession = browser.NewSession
	// NewPage wraps an existing browsing context id.
	NewPage = browser.NewPage
	// NewRequests records network events for a browsing context.
	NewRequests = browser.NewRequests
	// NewLogs records log entries for a browsing context.
	NewLogs = browser.NewLogs
)

// Options exposed through the root package.
var (
	// ErrClosed is returned when operating on a closed Client.
	ErrClosed = browser.ErrClosed
	// WithTimeout sets the default deadline for commands.
	WithTimeout = browser.WithTimeout
	// WithHTTPClient sets the client used for session creation.
	WithHTTPClient = browser.WithHTTPClient
	// WithSubprotocol overrides the negotiated WebSocket subprotocols.
	WithSubprotocol = browser.WithSubprotocol
	// WithCapabilities merges extra alwaysMatch capabilities.
	WithCapabilities = browser.WithCapabilities
)
