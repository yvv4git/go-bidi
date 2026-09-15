package protocol

import (
	"encoding/json"
	"testing"
)

func TestDecodeSuccessResponse(t *testing.T) {
	msg, err := Decode([]byte(`{"type":"success","id":1,"result":{"ready":true,"message":"ok"}}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	resp, ok := msg.(*Response)
	if !ok {
		t.Fatalf("got %T, want *Response", msg)
	}

	if resp.ID != 1 {
		t.Fatalf("id: got %d, want 1", resp.ID)
	}

	if resp.Err != nil {
		t.Fatalf("unexpected error: %v", resp.Err)
	}

	var status SessionStatusResult
	if err := json.Unmarshal(resp.Result, &status); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if !status.Ready {
		t.Fatalf("ready: got false, want true")
	}
}

func TestDecodeErrorResponse(t *testing.T) {
	msg, err := Decode([]byte(
		`{"type":"error","id":2,"error":"no such element",` +
			`"message":"element not found","stacktrace":"at foo"}`,
	))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	resp, ok := msg.(*Response)
	if !ok {
		t.Fatalf("got %T, want *Response", msg)
	}

	if resp.Err == nil {
		t.Fatal("expected an error")
	}

	if resp.Err.Code != "no such element" {
		t.Fatalf("code: got %q", resp.Err.Code)
	}

	if resp.Err.Message != "element not found" {
		t.Fatalf("message: got %q", resp.Err.Message)
	}

	if resp.Err.Stacktrace != "at foo" {
		t.Fatalf("stacktrace: got %q", resp.Err.Stacktrace)
	}
}

func TestDecodeEvent(t *testing.T) {
	msg, err := Decode([]byte(`{"type":"event","method":"log.entryAdded","params":{"text":"hi"}}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	event, ok := msg.(*Event)
	if !ok {
		t.Fatalf("got %T, want *Event", msg)
	}

	if event.Method != LogEntryAdded {
		t.Fatalf("method: got %q, want %q", event.Method, LogEntryAdded)
	}

	if string(event.Params) != `{"text":"hi"}` {
		t.Fatalf("params: got %s", event.Params)
	}
}

func TestDecodeWithoutType(t *testing.T) {
	tests := []struct {
		name  string
		give  string
		event bool
	}{
		{name: "response has id", give: `{"id":1,"result":{}}`, event: false},
		{name: "event has method", give: `{"method":"log.entryAdded","params":{}}`, event: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := Decode([]byte(tt.give))
			if err != nil {
				t.Fatalf("decode: %v", err)
			}

			_, isEvent := msg.(*Event)
			if isEvent != tt.event {
				t.Fatalf("got %T, want event=%v", msg, tt.event)
			}
		})
	}
}

func TestErrorString(t *testing.T) {
	err := &Error{Code: "unknown command", Message: "bad method"}

	if got, want := err.Error(), "bidi: unknown command: bad method"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
