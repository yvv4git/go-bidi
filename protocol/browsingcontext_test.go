package protocol

import (
	"encoding/json"
	"testing"
)

func TestBrowsingContextGetTreeResult(t *testing.T) {
	frame := `{"type":"success","id":1,"result":{"contexts":[` +
		`{"context":"c1","url":"https://example.com","children":[]}` +
		`]}}`

	msg, err := Decode([]byte(frame))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	resp, ok := msg.(*Response)
	if !ok {
		t.Fatalf("got %T, want *Response", msg)
	}

	var result BrowsingContextGetTreeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if len(result.Contexts) != 1 {
		t.Fatalf("contexts: got %d, want 1", len(result.Contexts))
	}

	if result.Contexts[0].Context != "c1" {
		t.Fatalf("context: got %q, want %q", result.Contexts[0].Context, "c1")
	}
}
