package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

type call struct {
	method string
	params any
}

type fakeCaller struct {
	mu      sync.Mutex
	calls   []call
	results map[string]string
	err     error
}

func (f *fakeCaller) Call(_ context.Context, method string, params, result any) error {
	f.mu.Lock()
	f.calls = append(f.calls, call{method: method, params: params})
	f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	if result == nil {
		return nil
	}

	raw, ok := f.results[method]
	if !ok {
		return fmt.Errorf("fakeCaller: no result for %s", method)
	}

	return json.Unmarshal([]byte(raw), result)
}

func (f *fakeCaller) lastCall(t *testing.T) call {
	t.Helper()

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.calls) == 0 {
		t.Fatal("fakeCaller: no calls recorded")
	}

	return f.calls[len(f.calls)-1]
}

func TestSessionStatus(t *testing.T) {
	tests := []struct {
		name    string
		ready   bool
		message string
	}{
		{name: "ready", ready: true, message: "ok"},
		{name: "not ready", ready: false, message: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(protocol.SessionStatusResult{Ready: tt.ready, Message: tt.message})
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			caller := &fakeCaller{results: map[string]string{protocol.SessionStatus: string(payload)}}
			session := NewSession(caller, "s1")

			got, err := session.Status(context.Background())
			if err != nil {
				t.Fatalf("Status: %v", err)
			}

			if got.Ready != tt.ready || got.Message != tt.message {
				t.Errorf("Status = %+v, want ready=%v message=%q", got, tt.ready, tt.message)
			}

			if method := caller.lastCall(t).method; method != protocol.SessionStatus {
				t.Errorf("method = %q, want %q", method, protocol.SessionStatus)
			}
		})
	}
}

func TestSessionEnd(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.End(context.Background()); err != nil {
		t.Fatalf("End: %v", err)
	}

	if method := caller.lastCall(t).method; method != protocol.SessionEnd {
		t.Errorf("method = %q, want %q", method, protocol.SessionEnd)
	}
}

func TestSessionSubscribeAndUnsubscribe(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.SessionSubscribe:   `{"subscription":"sub-1"}`,
		protocol.SessionUnsubscribe: `{}`,
	}}
	session := NewSession(caller, "s1")

	events := []string{"log.entryAdded"}

	got, err := session.Subscribe(
		context.Background(),
		protocol.SessionSubscribeParams{Events: events},
	)
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	if got.Subscription != "sub-1" {
		t.Fatalf("Subscription = %q, want %q", got.Subscription, "sub-1")
	}

	subs := session.Subscriptions()
	if len(subs["sub-1"]) != 1 || subs["sub-1"][0] != "log.entryAdded" {
		t.Fatalf("Subscriptions = %v, want sub-1 -> [log.entryAdded]", subs)
	}

	subs["sub-1"][0] = "changed"

	if session.Subscriptions()["sub-1"][0] != "log.entryAdded" {
		t.Error("Subscriptions returned an internal slice")
	}

	if err := session.Unsubscribe(context.Background(), "sub-1"); err != nil {
		t.Fatalf("Unsubscribe: %v", err)
	}

	if len(session.Subscriptions()) != 0 {
		t.Errorf("Subscriptions = %v, want empty", session.Subscriptions())
	}

	params, ok := caller.lastCall(t).params.(protocol.SessionUnsubscribeParams)
	if !ok {
		t.Fatalf("params = %T, want protocol.SessionUnsubscribeParams", caller.lastCall(t).params)
	}

	if len(params.Subscriptions) != 1 || params.Subscriptions[0] != "sub-1" {
		t.Errorf("Subscriptions = %v, want [sub-1]", params.Subscriptions)
	}
}

func TestCreateSession(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.SessionNew: `{"sessionId":"s-42","capabilities":{}}`,
	}}

	session, err := CreateSession(context.Background(), caller, protocol.SessionNewParams{})
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if session.ID() != "s-42" {
		t.Errorf("ID = %q, want %q", session.ID(), "s-42")
	}

	if method := caller.lastCall(t).method; method != protocol.SessionNew {
		t.Errorf("method = %q, want %q", method, protocol.SessionNew)
	}
}

func TestSessionCallError(t *testing.T) {
	caller := &fakeCaller{err: fmt.Errorf("boom")}
	session := NewSession(caller, "s1")

	if _, err := session.Status(context.Background()); err == nil {
		t.Fatal("Status: expected error")
	}

	if err := session.End(context.Background()); err == nil {
		t.Fatal("End: expected error")
	}
}
