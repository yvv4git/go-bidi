package browser

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/yvv4git/go-bidi/protocol"
)

// Entry summarizes a single log.entryAdded event collected by Logs.
type Entry struct {
	Timestamp int64
	Level     string
	Type      string
	Text      string
	Method    string
	Context   string
	Realm     string
}

// Logs records log entries delivered on a subscription for a single
// browsing context. Create one with NewLogs.
type Logs struct {
	sub       *Subscription
	contextID string
	mu        sync.Mutex
	entries   []*Entry
	wake      chan struct{}
	done      chan struct{}
	closed    bool
}

// NewLogs starts recording entryAdded events for contextID from the events
// on sub. The subscription must cover the log module events.
func NewLogs(sub *Subscription, contextID string) *Logs {
	return &Logs{
		sub:       sub,
		contextID: contextID,
		wake:      make(chan struct{}, 1),
		done:      make(chan struct{}),
	}
}

// List returns a snapshot of the entries recorded so far.
func (l *Logs) List() []*Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]*Entry, 0, len(l.entries))
	for _, entry := range l.entries {
		cloned := *entry
		out = append(out, &cloned)
	}

	return out
}

// Count returns the number of recorded entries.
func (l *Logs) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.entries)
}

// Poll waits until cond reports true for the snapshot of recorded entries,
// or until ctx is done. It returns the snapshot that satisfied cond.
func (l *Logs) Poll(ctx context.Context, cond func([]*Entry) bool) ([]*Entry, error) {
	for {
		if out := l.List(); cond(out) {
			return out, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-l.wake:
		case <-l.done:
			return nil, context.Canceled
		}
	}
}

// Close stops recording and unsubscribes the underlying subscription.
func (l *Logs) Close(ctx context.Context) error {
	err := l.sub.Close(ctx)

	l.mu.Lock()
	if !l.closed {
		close(l.done)
		l.closed = true
	}
	l.mu.Unlock()

	return err
}

// run drains the subscription and records events until it is closed.
func (l *Logs) run() {
	for event := range l.sub.Receive() {
		l.record(event.Method, event.Params)
	}
}

func (l *Logs) record(method string, params json.RawMessage) {
	if method != protocol.LogEntryAdded {
		return
	}

	var p protocol.LogEntryParams
	if err := json.Unmarshal(params, &p); err != nil {
		return
	}

	if p.Source.Context != l.contextID {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	select {
	case <-l.done:
		return
	default:
	}

	l.entries = append(l.entries, &Entry{
		Timestamp: p.Timestamp,
		Level:     p.Level,
		Type:      p.Type,
		Text:      p.Text,
		Method:    p.Method,
		Context:   p.Source.Context,
		Realm:     p.Source.Realm,
	})

	select {
	case l.wake <- struct{}{}:
	default:
	}
}
