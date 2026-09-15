// Package browser provides high-level control of a BiDi-enabled browser:
// connecting to a session and managing pages.
//
// A Client routes commands and events over a transport.Transport: commands
// are matched to responses by id and events are demultiplexed to
// subscriptions, so callers work with Session and Page handles instead of
// raw JSON.
package browser
