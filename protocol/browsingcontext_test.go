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

func TestBrowsingContextReloadParams(t *testing.T) {
	params := BrowsingContextReloadParams{
		Context:     "c1",
		IgnoreCache: true,
		Wait:        ReadinessNone,
	}

	raw, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded BrowsingContextReloadParams
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if decoded != params {
		t.Errorf("round-trip = %+v, want %+v", decoded, params)
	}
}

func TestBrowsingContextSetViewportParams(t *testing.T) {
	dpr := 2.0
	params := BrowsingContextSetViewportParams{
		Context:          "c1",
		Viewport:         &Viewport{Width: 800, Height: 600},
		DevicePixelRatio: &dpr,
	}

	raw, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded BrowsingContextSetViewportParams
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if decoded.Context != "c1" {
		t.Errorf("Context = %q, want %q", decoded.Context, "c1")
	}

	if decoded.Viewport == nil || decoded.Viewport.Width != 800 || decoded.Viewport.Height != 600 {
		t.Errorf("Viewport = %+v, want 800x600", decoded.Viewport)
	}

	if decoded.DevicePixelRatio == nil || *decoded.DevicePixelRatio != 2.0 {
		t.Errorf("DevicePixelRatio = %v, want 2.0", decoded.DevicePixelRatio)
	}
}

func TestNavigationInfoDecode(t *testing.T) {
	frame := `{"type":"event","method":"browsingContext.navigationStarted",` +
		`"params":{"context":"c1","navigation":"n1","timestamp":1234,"url":"https://example.com"}}`

	msg, err := Decode([]byte(frame))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	event, ok := msg.(*Event)
	if !ok {
		t.Fatalf("got %T, want *Event", msg)
	}

	var info NavigationInfo
	if err := json.Unmarshal(event.Params, &info); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if info.Context != "c1" || info.Navigation != "n1" || info.URL != "https://example.com" {
		t.Errorf("info = %+v", info)
	}

	if info.Timestamp != 1234 {
		t.Errorf("Timestamp = %d, want 1234", info.Timestamp)
	}
}
