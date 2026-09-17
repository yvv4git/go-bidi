package browser

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/yvv4git/go-bidi/protocol"
)

const _eventBuffer = 256

// Subscription delivers events for a session.subscribe call.
type Subscription struct {
	client *Client
	id     string
	events []string
	ch     chan protocol.Event

	dropped uint64
}

// Subscribe subscribes to the given event names or modules and returns the
// subscription handle that receives the matching events.
func (c *Client) Subscribe(ctx context.Context, events []string) (*Subscription, error) {
	if len(events) == 0 {
		return nil, errors.New("bidi: subscribe: no events")
	}

	params := protocol.SessionSubscribeParams{Events: events}

	var result protocol.SessionSubscribeResult
	if err := c.Call(ctx, protocol.SessionSubscribe, params, &result); err != nil {
		return nil, err
	}

	if result.Subscription == "" {
		return nil, errors.New("bidi: subscribe: empty subscription id")
	}

	sub := &Subscription{
		client: c,
		id:     result.Subscription,
		events: append([]string(nil), events...),
		ch:     make(chan protocol.Event, c.eventBuffer),
	}

	c.subMu.Lock()
	c.subs[sub.id] = sub
	c.subMu.Unlock()

	return sub, nil
}

// ID returns the subscription id.
func (s *Subscription) ID() string {
	return s.id
}

// Events returns the subscribed event names.
func (s *Subscription) Events() []string {
	return append([]string(nil), s.events...)
}

// Receive returns the channel that the matched events are published to. The
// channel is closed when the subscription ends.
func (s *Subscription) Receive() <-chan protocol.Event {
	return s.ch
}

// Dropped reports how many events were dropped because the channel was full.
func (s *Subscription) Dropped() uint64 {
	s.client.subMu.Lock()
	defer s.client.subMu.Unlock()

	return s.dropped
}

// Close ends the subscription and closes the event channel.
func (s *Subscription) Close(ctx context.Context) error {
	params := protocol.SessionUnsubscribeParams{Subscriptions: []string{s.id}}
	if err := s.client.Call(ctx, protocol.SessionUnsubscribe, params, nil); err != nil {
		return err
	}

	s.client.removeSubscription(s.id)

	return nil
}

func (c *Client) removeSubscription(id string) {
	c.subMu.Lock()
	defer c.subMu.Unlock()

	if sub, ok := c.subs[id]; ok {
		delete(c.subs, id)
		close(sub.ch)
	}
}

func (c *Client) dispatchEvent(event *protocol.Event) {
	c.subMu.Lock()
	defer c.subMu.Unlock()

	for _, sub := range c.subs {
		if !sub.matches(event.Method) {
			continue
		}

		delivery := *event
		delivery.Params = append(json.RawMessage(nil), event.Params...)

		select {
		case sub.ch <- delivery:
		default:
			sub.dropped++
		}
	}
}

func (s *Subscription) matches(method string) bool {
	for _, event := range s.events {
		if event == "*" || event == method || strings.HasPrefix(method, event+".") {
			return true
		}
	}

	return false
}
