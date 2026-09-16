package browser

import (
	"context"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestSessionPerformActions(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	params := protocol.InputPerformActionsParams{
		Context: "c1",
		Actions: []protocol.InputSourceActions{{Type: protocol.InputSourceTypeNone}},
	}

	if err := session.PerformActions(context.Background(), params); err != nil {
		t.Fatalf("PerformActions: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.InputPerformActions {
		t.Errorf("method = %q, want %q", call.method, protocol.InputPerformActions)
	}

	got := call.params.(protocol.InputPerformActionsParams)
	if got.Context != "c1" || len(got.Actions) != 1 {
		t.Errorf("params = %+v, want context c1 with one source", got)
	}
}

func TestSessionReleaseActions(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.ReleaseActions(context.Background(), protocol.InputReleaseActionsParams{
		Context: "c1",
	}); err != nil {
		t.Fatalf("ReleaseActions: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.InputReleaseActions {
		t.Errorf("method = %q, want %q", m, protocol.InputReleaseActions)
	}
}

func TestSessionSetFiles(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.SetFiles(context.Background(), protocol.InputSetFilesParams{
		Context: "c1",
		Files:   []string{"/tmp/a.txt"},
	}); err != nil {
		t.Fatalf("SetFiles: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.InputSetFiles {
		t.Errorf("method = %q, want %q", call.method, protocol.InputSetFiles)
	}

	params := call.params.(protocol.InputSetFilesParams)
	if params.Context != "c1" || len(params.Files) != 1 {
		t.Errorf("params = %+v, want context c1 with one file", params)
	}
}

func TestPageClick(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.Click(context.Background(), 10, 20); err != nil {
		t.Fatalf("Click: %v", err)
	}

	action := sourceAction(t, caller, protocol.InputSourceTypePointer, 0)

	if m := caller.lastCall(t).method; m != protocol.InputPerformActions {
		t.Errorf("method = %q, want %q", m, protocol.InputPerformActions)
	}

	if action.X != 10 || action.Y != 20 {
		t.Errorf("move = %+v, want 10, 20", action)
	}

	if action.Origin == nil || *action.Origin != protocol.PointerOrigin(protocol.OriginViewport) {
		t.Errorf("origin = %+v, want viewport", action.Origin)
	}
}

func TestPagePress(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.Press(context.Background(), "Enter"); err != nil {
		t.Fatalf("Press: %v", err)
	}

	source := source(t, caller, protocol.InputSourceTypeKey)

	if len(source.Actions) != 2 {
		t.Fatalf("actions = %+v, want keyDown and keyUp", source.Actions)
	}

	if source.Actions[0].Type != protocol.KeyActionKeyDown || source.Actions[0].Value != "\r" {
		t.Errorf("actions[0] = %+v, want keyDown Enter", source.Actions[0])
	}

	if source.Actions[1].Type != protocol.KeyActionKeyUp {
		t.Errorf("actions[1] = %+v, want keyUp", source.Actions[1])
	}
}

func TestPageType(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.Type(context.Background(), "ab"); err != nil {
		t.Fatalf("Type: %v", err)
	}

	source := source(t, caller, protocol.InputSourceTypeKey)

	if len(source.Actions) != 4 {
		t.Fatalf("actions = %+v, want four key events", source.Actions)
	}

	if source.Actions[0].Value != "a" || source.Actions[2].Value != "b" {
		t.Errorf("values = %q, %q, want a and b", source.Actions[0].Value, source.Actions[2].Value)
	}
}

func TestPageScroll(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	page := NewPage(caller, "c1")

	if err := page.Scroll(context.Background(), -50, 100); err != nil {
		t.Fatalf("Scroll: %v", err)
	}

	source := source(t, caller, protocol.InputSourceTypeWheel)

	if len(source.Actions) != 1 {
		t.Fatalf("actions = %+v, want one scroll", source.Actions)
	}

	action := source.Actions[0]
	if action.Type != protocol.WheelActionScroll || action.DX != -50 || action.DY != 100 {
		t.Errorf("scroll = %+v, want scroll -50, 100", action)
	}
}

func source(t *testing.T, caller *fakeCaller, sourceType string) protocol.InputSourceActions {
	t.Helper()

	params := caller.lastCall(t).params.(protocol.InputPerformActionsParams)
	if params.Context != "c1" {
		t.Errorf("context = %q, want c1", params.Context)
	}

	if len(params.Actions) != 1 {
		t.Fatalf("sources = %+v, want one source", params.Actions)
	}

	source := params.Actions[0]
	if source.Type != sourceType {
		t.Fatalf("source type = %q, want %q", source.Type, sourceType)
	}

	return source
}

func sourceAction(
	t *testing.T,
	caller *fakeCaller,
	sourceType string,
	index int,
) protocol.InputAction {
	t.Helper()

	actions := source(t, caller, sourceType).Actions
	if index >= len(actions) {
		t.Fatalf("actions = %+v, want index %d", actions, index)
	}

	return actions[index]
}
