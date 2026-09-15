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
