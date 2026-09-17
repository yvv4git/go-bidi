package browser

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestSessionCreateBrowsingContext(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextCreate: `{"context":"c1"}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.CreateBrowsingContext(
		context.Background(),
		protocol.BrowsingContextCreateParams{Type: protocol.ContextTypeTab},
	)
	if err != nil {
		t.Fatalf("CreateBrowsingContext: %v", err)
	}

	if got.Context != "c1" {
		t.Errorf("Context = %q, want %q", got.Context, "c1")
	}

	if m := caller.lastCall(t).method; m != protocol.BrowsingContextCreate {
		t.Errorf("method = %q, want %q", m, protocol.BrowsingContextCreate)
	}
}

func TestSessionNavigate(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextNavigate: `{"navigation":"n1","url":"https://example.com"}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.Navigate(context.Background(), protocol.BrowsingContextNavigateParams{
		Context: "c1",
		URL:     "https://example.com",
		Wait:    protocol.ReadinessComplete,
	})
	if err != nil {
		t.Fatalf("Navigate: %v", err)
	}

	if got.Navigation != "n1" {
		t.Errorf("Navigation = %q, want %q", got.Navigation, "n1")
	}

	if got.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", got.URL, "https://example.com")
	}

	params, ok := caller.lastCall(t).params.(protocol.BrowsingContextNavigateParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Wait != protocol.ReadinessComplete {
		t.Errorf("Wait = %q, want %q", params.Wait, protocol.ReadinessComplete)
	}
}

func TestSessionGetTree(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextGetTree: `{"contexts":[` +
			`{"context":"c1","url":"about:blank","children":[]}` +
			`]}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.GetTree(context.Background(), protocol.BrowsingContextGetTreeParams{})
	if err != nil {
		t.Fatalf("GetTree: %v", err)
	}

	if len(got.Contexts) != 1 || got.Contexts[0].Context != "c1" {
		t.Errorf("Contexts = %+v, want one c1", got.Contexts)
	}
}

func TestSessionCloseBrowsingContext(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.CloseBrowsingContext(context.Background(), protocol.BrowsingContextCloseParams{
		Context: "c1",
	}); err != nil {
		t.Fatalf("CloseBrowsingContext: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.BrowsingContextClose {
		t.Errorf("method = %q, want %q", m, protocol.BrowsingContextClose)
	}
}

func TestSessionReload(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextReload: `{"navigation":"n2","url":"https://example.com"}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.Reload(context.Background(), protocol.BrowsingContextReloadParams{
		Context:     "c1",
		IgnoreCache: true,
	})
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if got.Navigation != "n2" {
		t.Errorf("Navigation = %q, want %q", got.Navigation, "n2")
	}
}

func TestSessionTraverseHistory(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.TraverseHistory(
		context.Background(),
		protocol.BrowsingContextTraverseHistoryParams{Context: "c1", Delta: -1},
	); err != nil {
		t.Fatalf("TraverseHistory: %v", err)
	}

	params, ok := caller.lastCall(t).params.(protocol.BrowsingContextTraverseHistoryParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Delta != -1 {
		t.Errorf("Delta = %d, want -1", params.Delta)
	}
}

func TestSessionCaptureScreenshot(t *testing.T) {
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	encoded := base64.StdEncoding.EncodeToString(png)
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextCaptureScreenshot: fmt.Sprintf(`{"data":"%s"}`, encoded),
	}}
	session := NewSession(caller, "s1")

	got, err := session.CaptureScreenshot(
		context.Background(),
		protocol.BrowsingContextCaptureScreenshotParams{Context: "c1", Origin: protocol.OriginViewport},
	)
	if err != nil {
		t.Fatalf("CaptureScreenshot: %v", err)
	}

	if !bytes.Equal(got, png) {
		t.Errorf("data = %x, want %x", got, png)
	}
}

func TestSessionSetViewport(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	err := session.SetViewport(context.Background(), protocol.BrowsingContextSetViewportParams{
		Context:  "c1",
		Viewport: &protocol.Viewport{Width: 800, Height: 600},
	})
	if err != nil {
		t.Fatalf("SetViewport: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.BrowsingContextSetViewport {
		t.Errorf("method = %q, want %q", m, protocol.BrowsingContextSetViewport)
	}
}
