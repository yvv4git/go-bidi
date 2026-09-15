package browser

import (
	"context"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestSessionGetCookies(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.StorageGetCookies: `{"cookies":[{"name":"sid","value":{"type":"string","value":"1"}}]}`,
	}}
	session := NewSession(caller, "s1")

	cookies, err := session.GetCookies(context.Background(), protocol.StorageGetCookiesParams{})
	if err != nil {
		t.Fatalf("GetCookies: %v", err)
	}

	if len(cookies) != 1 || cookies[0].Name != "sid" {
		t.Errorf("cookies = %+v", cookies)
	}

	if m := caller.lastCall(t).method; m != protocol.StorageGetCookies {
		t.Errorf("method = %q, want %q", m, protocol.StorageGetCookies)
	}
}

func TestSessionSetCookie(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	err := session.SetCookie(context.Background(), protocol.StorageSetCookieParams{
		Cookie: protocol.StorageCookie{
			Name:   "a",
			Value:  protocol.NetworkBytesValue{Type: protocol.BytesValueString, Value: "b"},
			Domain: "example.com",
		},
	})
	if err != nil {
		t.Fatalf("SetCookie: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.StorageSetCookie {
		t.Errorf("method = %q, want %q", call.method, protocol.StorageSetCookie)
	}

	params := call.params.(protocol.StorageSetCookieParams)
	if params.Cookie.Name != "a" || params.Cookie.Domain != "example.com" {
		t.Errorf("params = %+v", params)
	}
}

func TestSessionDeleteCookies(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	err := session.DeleteCookies(context.Background(), protocol.StorageDeleteCookiesParams{})
	if err != nil {
		t.Fatalf("DeleteCookies: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.StorageDeleteCookies {
		t.Errorf("method = %q, want %q", m, protocol.StorageDeleteCookies)
	}
}

func TestPageGetLocalStorage(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1",` +
			`"result":{"type":"string","value":"val1"}}`,
	}}
	page := NewPage(caller, "c1")

	got, err := page.GetLocalStorage(context.Background(), "k")
	if err != nil {
		t.Fatalf("GetLocalStorage: %v", err)
	}

	if got != "val1" {
		t.Errorf("value = %q, want val1", got)
	}

	assertExpression(t, caller, `localStorage.getItem("k")`)
}

func TestPageGetLocalStorageMissing(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1","result":{"type":"undefined"}}`,
	}}
	page := NewPage(caller, "c1")

	got, err := page.GetLocalStorage(context.Background(), "k")
	if err != nil {
		t.Fatalf("GetLocalStorage: %v", err)
	}

	if got != "" {
		t.Errorf("value = %q, want empty", got)
	}
}

func TestPageSetLocalStorage(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1","result":{"type":"undefined"}}`,
	}}
	page := NewPage(caller, "c1")

	if err := page.SetLocalStorage(context.Background(), "k", "val1"); err != nil {
		t.Fatalf("SetLocalStorage: %v", err)
	}

	assertExpression(t, caller, `localStorage.setItem("k", "val1")`)
}

func TestPageRemoveLocalStorage(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1","result":{"type":"undefined"}}`,
	}}
	page := NewPage(caller, "c1")

	if err := page.RemoveLocalStorage(context.Background(), "k"); err != nil {
		t.Fatalf("RemoveLocalStorage: %v", err)
	}

	assertExpression(t, caller, `localStorage.removeItem("k")`)
}

func TestPageSessionStorage(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.ScriptEvaluate: `{"type":"success","realm":"r1",` +
			`"result":{"type":"string","value":"s"}}`,
	}}
	page := NewPage(caller, "c1")

	got, err := page.GetSessionStorage(context.Background(), "k")
	if err != nil {
		t.Fatalf("GetSessionStorage: %v", err)
	}

	if got != "s" {
		t.Errorf("value = %q, want s", got)
	}

	assertExpression(t, caller, `sessionStorage.getItem("k")`)
}

func assertExpression(t *testing.T, caller *fakeCaller, want string) {
	t.Helper()

	params, ok := caller.lastCall(t).params.(protocol.ScriptEvaluateParams)
	if !ok {
		t.Fatalf("params = %T", caller.lastCall(t).params)
	}

	if params.Expression != want {
		t.Errorf("Expression = %q, want %q", params.Expression, want)
	}
}
