package browser

import (
	"context"
	"reflect"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestSessionEvaluate(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1",` +
			`"result":{"type":"string","value":"hi"}}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.Evaluate(context.Background(), protocol.ScriptEvaluateParams{
		Expression: "document.title",
		Target:     protocol.ScriptTarget{Context: "c1"},
	})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}

	value, err := got.Result.Interface()
	if err != nil {
		t.Fatalf("Interface: %v", err)
	}

	if value != "hi" {
		t.Errorf("value = %#v, want %q", value, "hi")
	}

	if m := caller.lastCall(t).method; m != protocol.ScriptEvaluate {
		t.Errorf("method = %q, want %q", m, protocol.ScriptEvaluate)
	}
}

func TestSessionCallFunction(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptCallFunction: `{"type":"success","realm":"r1",` +
			`"result":{"type":"number","value":42}}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.CallFunction(context.Background(), protocol.ScriptCallFunctionParams{
		FunctionDeclaration: "() => 42",
		Target:              protocol.ScriptTarget{Context: "c1"},
	})
	if err != nil {
		t.Fatalf("CallFunction: %v", err)
	}

	value, err := got.Result.Interface()
	if err != nil {
		t.Fatalf("Interface: %v", err)
	}

	if value != float64(42) {
		t.Errorf("value = %#v, want %v", value, float64(42))
	}
}

func TestSessionGetRealms(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptGetRealms: `{"realms":[{"realm":"r1","origin":"o","type":"window"}]}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.GetRealms(context.Background(), protocol.ScriptGetRealmsParams{})
	if err != nil {
		t.Fatalf("GetRealms: %v", err)
	}

	if len(got.Realms) != 1 || got.Realms[0].Realm != "r1" {
		t.Errorf("Realms = %+v, want one r1", got.Realms)
	}
}

func TestSessionAddPreloadScript(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptAddPreloadScript: `{"script":"s1"}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.AddPreloadScript(context.Background(), protocol.ScriptAddPreloadScriptParams{
		FunctionDeclaration: "() => window.__preloaded = true",
	})
	if err != nil {
		t.Fatalf("AddPreloadScript: %v", err)
	}

	if got.Script != "s1" {
		t.Errorf("Script = %q, want %q", got.Script, "s1")
	}
}

func TestSessionRemovePreloadScript(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.RemovePreloadScript(
		context.Background(),
		protocol.ScriptRemovePreloadScriptParams{Script: "s1"},
	); err != nil {
		t.Fatalf("RemovePreloadScript: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.ScriptRemovePreloadScript {
		t.Errorf("method = %q, want %q", m, protocol.ScriptRemovePreloadScript)
	}
}

func TestSessionDisown(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.Disown(context.Background(), protocol.ScriptDisownParams{
		Handles: []string{"h1"},
		Target:  protocol.ScriptTarget{Realm: "r1"},
	}); err != nil {
		t.Fatalf("Disown: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.ScriptDisown {
		t.Errorf("method = %q, want %q", m, protocol.ScriptDisown)
	}
}

func TestPageEvaluate(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1","result":{"type":"number","value":2}}`,
	}}
	page := NewPage(caller, "c1")

	got, err := page.Evaluate(context.Background(), "1 + 1")
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}

	value, err := got.Result.Interface()
	if err != nil {
		t.Fatalf("Interface: %v", err)
	}

	if value != float64(2) {
		t.Errorf("value = %#v, want 2", value)
	}

	params, ok := caller.lastCall(t).params.(protocol.ScriptEvaluateParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Target.Context != "c1" {
		t.Errorf("Target = %+v, want context c1", params.Target)
	}

	if !params.AwaitPromise {
		t.Error("AwaitPromise is false, want true")
	}

	if params.ResultOwnership != protocol.ResultOwnershipNone {
		t.Errorf("ResultOwnership = %q, want %q", params.ResultOwnership, protocol.ResultOwnershipNone)
	}
}

func TestPageCallFunction(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptCallFunction: `{"type":"success","realm":"r1",` +
			`"result":{"type":"number","value":42}}`,
	}}
	page := NewPage(caller, "c1")

	if _, err := page.CallFunction(context.Background(), "(a, b) => a + b"); err != nil {
		t.Fatalf("CallFunction no args: %v", err)
	}

	params := caller.lastCall(t).params.(protocol.ScriptCallFunctionParams)
	if params.Arguments == nil || len(params.Arguments) != 0 {
		t.Errorf("Arguments = %#v, want empty non-nil slice", params.Arguments)
	}

	_, err := page.CallFunction(
		context.Background(),
		"(a, b) => a + b",
		protocol.ScriptArgument{Type: protocol.RemoteTypeNumber, Value: []byte(`1`)},
		protocol.ScriptArgument{Type: protocol.RemoteTypeNumber, Value: []byte(`2`)},
	)
	if err != nil {
		t.Fatalf("CallFunction with args: %v", err)
	}

	params = caller.lastCall(t).params.(protocol.ScriptCallFunctionParams)
	if !reflect.DeepEqual(params.Arguments, []protocol.ScriptArgument{
		{Type: protocol.RemoteTypeNumber, Value: []byte(`1`)},
		{Type: protocol.RemoteTypeNumber, Value: []byte(`2`)},
	}) {
		t.Errorf("Arguments = %+v", params.Arguments)
	}
}
