package browser

import (
	"context"
	"strings"
	"testing"
	"time"
)

func closeClient(t *testing.T, client *Client) {
	t.Helper()

	if err := client.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func nextEvent(t *testing.T, sub *Subscription) (method string) {
	t.Helper()

	select {
	case event, ok := <-sub.Receive():
		if !ok {
			t.Fatal("subscription channel closed")
		}

		return event.Method
	case <-time.After(time.Second):
		t.Fatal("no event received")
	}

	return ""
}

func TestSubscribeReceivesMatchingEvent(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	sub, err := client.Subscribe(context.Background(), []string{"log.entryAdded"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	if got := sub.ID(); got != "sub-1" {
		t.Errorf("ID = %q, want %q", got, "sub-1")
	}

	if events := strings.Join(sub.Events(), ","); events != "log.entryAdded" {
		t.Errorf("Events = %q, want %q", events, "log.entryAdded")
	}

	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{"level":"info"}}`)

	if method := nextEvent(t, sub); method != "log.entryAdded" {
		t.Errorf("method = %q, want %q", method, "log.entryAdded")
	}
}

func TestSubscribeModuleMatch(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	sub, err := client.Subscribe(context.Background(), []string{"log"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{}}`)
	if method := nextEvent(t, sub); method != "log.entryAdded" {
		t.Errorf("method = %q, want %q", method, "log.entryAdded")
	}
}

func TestSubscribeNonMatchingEvent(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	sub, err := client.Subscribe(context.Background(), []string{"network"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{}}`)

	select {
	case event, ok := <-sub.Receive():
		if ok {
			t.Fatalf("unexpected event %s", event.Method)
		}

		t.Fatal("subscription channel closed")
	case <-time.After(20 * time.Millisecond):
	}
}

func TestSubscribeFanout(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	first, err := client.Subscribe(context.Background(), []string{"log.entryAdded"})
	if err != nil {
		t.Fatalf("first Subscribe: %v", err)
	}

	go respondSuccess(t, tr, `{"subscription":"sub-2"}`)

	second, err := client.Subscribe(context.Background(), []string{"log"})
	if err != nil {
		t.Fatalf("second Subscribe: %v", err)
	}

	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{}}`)

	if method := nextEvent(t, first); method != "log.entryAdded" {
		t.Errorf("first method = %q", method)
	}

	if method := nextEvent(t, second); method != "log.entryAdded" {
		t.Errorf("second method = %q", method)
	}
}

func TestSubscribeDropsWhenFull(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	client.eventBuffer = 1

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	sub, err := client.Subscribe(context.Background(), []string{"log"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{}}`)
	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{}}`)

	deadline := time.Now().Add(time.Second)
	for sub.Dropped() != 1 {
		if time.Now().After(deadline) {
			t.Fatalf("Dropped = %d, want 1", sub.Dropped())
		}

		time.Sleep(time.Millisecond)
	}

	if method := nextEvent(t, sub); method != "log.entryAdded" {
		t.Errorf("method = %q, want %q", method, "log.entryAdded")
	}
}

func TestSubscriptionClose(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	sub, err := client.Subscribe(context.Background(), []string{"log"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	go respondSuccess(t, tr, `{}`)
	if err := sub.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}

	var sentUnsubscribe bool
	for _, frame := range tr.frames(t) {
		if strings.Contains(frame, `"method":"session.unsubscribe"`) &&
			strings.Contains(frame, `"subscriptions":["sub-1"]`) {
			sentUnsubscribe = true
		}
	}
	if !sentUnsubscribe {
		t.Errorf("frames = %v, want session.unsubscribe for sub-1", tr.frames(t))
	}

	if event, ok := <-sub.Receive(); ok {
		t.Fatalf("channel not closed, got event %s", event.Method)
	}

	tr.recvCh <- []byte(`{"type":"event","method":"log.entryAdded","params":{}}`)
	time.Sleep(10 * time.Millisecond)
}

func TestSubscribeNoEvents(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	if _, err := client.Subscribe(context.Background(), nil); err == nil ||
		!strings.Contains(err.Error(), "no events") {
		t.Fatalf("Subscribe error = %v, want no-events error", err)
	}
}

func TestSubscribeEmptyID(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{}`)

	if _, err := client.Subscribe(context.Background(), []string{"log"}); err == nil ||
		!strings.Contains(err.Error(), "empty subscription id") {
		t.Fatalf("Subscribe error = %v, want empty-id error", err)
	}
}
