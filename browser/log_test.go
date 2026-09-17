package browser

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestLogsCollectsEntries(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	logs := NewLogs(sub, "c1")

	logs.record(protocol.LogEntryAdded, entryParams(t, logEntryParams("c1", "hello")))
	logs.record(protocol.LogEntryAdded, entryParams(t, logEntryParams("c1", "warn")))
	logs.record(protocol.LogEntryAdded, entryParams(t, logEntryParams("c2", "other")))
	logs.record("some.module.event", entryParams(t, logEntryParams("c1", "ignored")))

	got := logs.List()
	if len(got) != 2 {
		t.Fatalf("entries = %d, want 2", len(got))
	}

	if got[0].Text != "hello" || got[0].Level != protocol.LogLevelInfo {
		t.Errorf("entry = %+v", got[0])
	}

	if got[0].Type != protocol.LogTypeConsole || got[0].Method != "log" {
		t.Errorf("entry = %+v", got[0])
	}

	if got[0].Context != "c1" || got[0].Realm != "r1" {
		t.Errorf("entry source = %+v", got[0])
	}
}

func TestLogsPoll(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	logs := NewLogs(sub, "c1")

	done := make(chan struct{})
	go func() {
		defer close(done)
		got, err := logs.Poll(context.Background(), func(entries []*Entry) bool {
			return len(entries) == 1 && entries[0].Level == protocol.LogLevelError
		})
		if err != nil {
			t.Errorf("Poll: %v", err)
		}

		if len(got) != 1 || got[0].Text != "boom" {
			t.Errorf("Poll result = %+v", got)
		}
	}()

	logs.record(protocol.LogEntryAdded, entryParams(t, logEntryParams("c1", "boom")))

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Poll did not return")
	}
}

func TestLogsRun(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	logs := NewLogs(sub, "c1")
	go logs.run()

	sub.ch <- protocol.Event{
		Method: protocol.LogEntryAdded,
		Params: entryParams(t, logEntryParams("c1", "run")),
	}

	waitLogs(t, logs, func(entries []*Entry) bool { return len(entries) == 1 })
}

func TestLogsCount(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	logs := NewLogs(sub, "c1")

	if got := logs.Count(); got != 0 {
		t.Errorf("Count = %d, want 0", got)
	}

	logs.record(protocol.LogEntryAdded, entryParams(t, logEntryParams("c1", "a")))
	logs.record(protocol.LogEntryAdded, entryParams(t, logEntryParams("c1", "b")))

	if got := logs.Count(); got != 2 {
		t.Errorf("Count = %d, want 2", got)
	}
}

func TestLogsClose(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer closeClient(t, client)

	go respondSuccess(t, tr, `{"subscription":"sub-1"}`)

	sub, err := client.Subscribe(context.Background(), []string{"log"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	logs := NewLogs(sub, "c1")
	go logs.run()

	go respondSuccess(t, tr, `{}`)

	if err := logs.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}

	select {
	case _, open := <-sub.Receive():
		if open {
			t.Error("subscription channel still open after Close")
		}
	default:
	}

	if _, err := logs.Poll(context.Background(), func(_ []*Entry) bool {
		return false
	}); err != context.Canceled {
		t.Errorf("Poll error = %v, want context.Canceled", err)
	}
}

func waitLogs(t *testing.T, logs *Logs, cond func([]*Entry) bool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := logs.Poll(ctx, cond); err != nil {
		t.Fatalf("Poll: %v", err)
	}
}

func logEntryParams(contextID, text string) protocol.LogEntryParams {
	params := protocol.LogEntryParams{
		Type:      protocol.LogTypeConsole,
		Level:     protocol.LogLevelInfo,
		Source:    protocol.LogSource{Realm: "r1", Context: contextID},
		Text:      text,
		Timestamp: 1,
		Method:    "log",
	}

	if text == "boom" {
		params.Level = protocol.LogLevelError
	}

	return params
}

func entryParams(t *testing.T, params protocol.LogEntryParams) json.RawMessage {
	t.Helper()

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return data
}
