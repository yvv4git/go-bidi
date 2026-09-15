// Package browser provides high-level control of a BiDi-enabled browser:
// connecting to a session and managing pages.
//
// A Browser wraps a transport.Transport with a message router that turns
// protocol messages into typed calls, so callers work with Browser and Page
// handles instead of raw JSON.
package browser
