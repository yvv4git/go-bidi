package browser

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestPageID(t *testing.T) {
	page := NewPage(&fakeCaller{}, "c1")

	if got := page.ID(); got != "c1" {
		t.Errorf("ID = %q, want %q", got, "c1")
	}
}

func TestPageNavigate(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextNavigate: `{"navigation":"n1","url":"https://example.com"}`,
	}}
	page := NewPage(caller, "c1")

	got, err := page.Navigate(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("Navigate: %v", err)
	}

	if got.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", got.URL, "https://example.com")
	}

	params, ok := caller.lastCall(t).params.(protocol.BrowsingContextNavigateParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Context != "c1" {
		t.Errorf("Context = %q, want %q", params.Context, "c1")
	}

	if params.Wait != protocol.ReadinessComplete {
		t.Errorf("Wait = %q, want %q", params.Wait, protocol.ReadinessComplete)
	}
}

func TestPageReload(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextReload: `{"navigation":"n2","url":"https://example.com"}`,
	}}
	page := NewPage(caller, "c1")

	got, err := page.Reload(context.Background())
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if got.Navigation != "n2" {
		t.Errorf("Navigation = %q, want %q", got.Navigation, "n2")
	}

	params, ok := caller.lastCall(t).params.(protocol.BrowsingContextReloadParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Context != "c1" || params.Wait != protocol.ReadinessComplete {
		t.Errorf("params = %+v, want context c1 wait %s", params, protocol.ReadinessComplete)
	}
}

func TestPageBackAndForward(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.Back(context.Background()); err != nil {
		t.Fatalf("Back: %v", err)
	}

	back := caller.lastCall(t).params.(protocol.BrowsingContextTraverseHistoryParams)
	if back.Delta != -1 {
		t.Errorf("Back Delta = %d, want -1", back.Delta)
	}

	if err := page.Forward(context.Background()); err != nil {
		t.Fatalf("Forward: %v", err)
	}

	forward := caller.lastCall(t).params.(protocol.BrowsingContextTraverseHistoryParams)
	if forward.Delta != 1 {
		t.Errorf("Forward Delta = %d, want 1", forward.Delta)
	}
}

func TestPageScreenshot(t *testing.T) {
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	encoded := base64.StdEncoding.EncodeToString(png)
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowsingContextCaptureScreenshot: fmt.Sprintf(`{"data":"%s"}`, encoded),
	}}
	page := NewPage(caller, "c1")

	got, err := page.Screenshot(context.Background(), true)
	if err != nil {
		t.Fatalf("Screenshot: %v", err)
	}

	if !bytes.Equal(got, png) {
		t.Errorf("data = %x, want %x", got, png)
	}

	params, ok := caller.lastCall(t).params.(protocol.BrowsingContextCaptureScreenshotParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if !params.CaptureBeyondViewport {
		t.Error("CaptureBeyondViewport is false, want true")
	}

	if params.Origin != protocol.OriginViewport {
		t.Errorf("Origin = %q, want %q", params.Origin, protocol.OriginViewport)
	}
}

func TestPageSetViewport(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.SetViewport(context.Background(), 800, 600); err != nil {
		t.Fatalf("SetViewport: %v", err)
	}

	params, ok := caller.lastCall(t).params.(protocol.BrowsingContextSetViewportParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Viewport == nil || params.Viewport.Width != 800 || params.Viewport.Height != 600 {
		t.Errorf("Viewport = %+v, want 800x600", params.Viewport)
	}
}

func TestPageClose(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.BrowsingContextClose {
		t.Errorf("method = %q, want %q", m, protocol.BrowsingContextClose)
	}
}
