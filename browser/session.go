package browser

import (
	"context"
	"sync"

	"github.com/yvv4git/go-bidi/protocol"
)

// Session is a handle to an established BiDi session. It issues session-module
// commands through a Caller and keeps track of active subscriptions.
type Session struct {
	caller Caller
	id     string

	mu   sync.Mutex
	subs map[string][]string
}

// NewSession returns a handle for an already-created session id. The caller
// must be able to route commands to that session's transport.
func NewSession(caller Caller, id string) *Session {
	return &Session{
		caller: caller,
		id:     id,
		subs:   make(map[string][]string),
	}
}

// CreateSession runs session.new over caller and returns the created session.
func CreateSession(
	ctx context.Context,
	caller Caller,
	params protocol.SessionNewParams,
) (*Session, error) {
	var result protocol.SessionNewResult
	if err := caller.Call(ctx, protocol.SessionNew, params, &result); err != nil {
		return nil, err
	}

	return NewSession(caller, result.SessionID), nil
}

// ID returns the session id assigned by the protocol.
func (s *Session) ID() string {
	return s.id
}

// Status reports whether the remote end is ready to accept commands.
func (s *Session) Status(ctx context.Context) (*protocol.SessionStatusResult, error) {
	var result protocol.SessionStatusResult
	if err := s.caller.Call(ctx, protocol.SessionStatus, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// End terminates the session.
func (s *Session) End(ctx context.Context) error {
	return s.caller.Call(ctx, protocol.SessionEnd, nil, nil)
}

// Subscribe subscribes to the requested events and records the subscription.
func (s *Session) Subscribe(
	ctx context.Context,
	params protocol.SessionSubscribeParams,
) (*protocol.SessionSubscribeResult, error) {
	var result protocol.SessionSubscribeResult
	if err := s.caller.Call(ctx, protocol.SessionSubscribe, params, &result); err != nil {
		return nil, err
	}

	s.track(result.Subscription, params.Events)

	return &result, nil
}

// Unsubscribe removes a subscription and forgets its recorded events.
func (s *Session) Unsubscribe(ctx context.Context, subscription string) error {
	params := protocol.SessionUnsubscribeParams{Subscriptions: []string{subscription}}
	if err := s.caller.Call(ctx, protocol.SessionUnsubscribe, params, nil); err != nil {
		return err
	}

	s.untrack(subscription)

	return nil
}

// Subscriptions returns a copy of the active subscriptions keyed by id.
func (s *Session) Subscriptions() map[string][]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs := make(map[string][]string, len(s.subs))
	for id, events := range s.subs {
		subs[id] = append([]string(nil), events...)
	}

	return subs
}

func (s *Session) track(id string, events []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subs[id] = append([]string(nil), events...)
}

func (s *Session) untrack(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subs, id)
}
