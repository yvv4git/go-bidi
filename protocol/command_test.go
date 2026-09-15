package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewCommandWithParams(t *testing.T) {
	cmd, err := NewCommand(3, BrowsingContextNavigate, BrowsingContextNavigateParams{
		Context: "abc",
		URL:     "https://example.com",
		Wait:    ReadinessComplete,
	})
	if err != nil {
		t.Fatalf("new command: %v", err)
	}

	if cmd.ID != 3 {
		t.Fatalf("id: got %d, want 3", cmd.ID)
	}

	var got BrowsingContextNavigateParams
	if err := json.Unmarshal(cmd.Params, &got); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if got.Context != "abc" || got.URL != "https://example.com" || got.Wait != ReadinessComplete {
		t.Fatalf("params: got %+v", got)
	}
}

func TestNewCommandWithoutParams(t *testing.T) {
	cmd, err := NewCommand(1, SessionStatus, nil)
	if err != nil {
		t.Fatalf("new command: %v", err)
	}

	if cmd.Params != nil {
		t.Fatalf("params: got %s, want nil", cmd.Params)
	}
}

func TestNewCommandRawParams(t *testing.T) {
	raw := json.RawMessage(`{"type":"tab"}`)

	cmd, err := NewCommand(2, BrowsingContextCreate, raw)
	if err != nil {
		t.Fatalf("new command: %v", err)
	}

	if string(cmd.Params) != string(raw) {
		t.Fatalf("params: got %s, want %s", cmd.Params, raw)
	}
}

func TestCommandEncode(t *testing.T) {
	cmd, err := NewCommand(7, SessionStatus, nil)
	if err != nil {
		t.Fatalf("new command: %v", err)
	}

	data, err := cmd.Encode()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if got, want := string(data), `{"id":7,"method":"session.status"}`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
