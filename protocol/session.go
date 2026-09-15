package protocol

// SessionStatusResult is the result of session.status. It reports whether
// the remote end is ready to accept commands.
type SessionStatusResult struct {
	Ready   bool   `json:"ready"`
	Message string `json:"message"`
}

// SessionCapabilities are the capabilities requested by session.new.
type SessionCapabilities struct {
	AlwaysMatch map[string]any   `json:"alwaysMatch,omitempty"`
	FirstMatch  []map[string]any `json:"firstMatch,omitempty"`
}

// SessionNewParams are the parameters of session.new.
type SessionNewParams struct {
	Capabilities SessionCapabilities `json:"capabilities"`
}

// SessionNewResult is the result of session.new.
type SessionNewResult struct {
	SessionID    string         `json:"sessionId"`
	Capabilities map[string]any `json:"capabilities"`
}

// SessionSubscribeParams are the parameters of session.subscribe. Events
// holds event names or module names; Contexts and UserContexts scope the
// subscription to specific browsing contexts.
type SessionSubscribeParams struct {
	Events       []string `json:"events,omitempty"`
	Contexts     []string `json:"contexts,omitempty"`
	UserContexts []string `json:"userContexts,omitempty"`
}

// SessionSubscribeResult is the result of session.subscribe. The
// Subscription id is used to unsubscribe later.
type SessionSubscribeResult struct {
	Subscription string `json:"subscription"`
}
