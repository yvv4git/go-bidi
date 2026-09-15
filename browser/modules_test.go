package browser

import (
	"context"
	"testing"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestSessionSetPermission(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	err := session.SetPermission(context.Background(), protocol.PermissionsSetPermissionParams{
		Descriptor: protocol.PermissionsDescriptor{Name: "notifications"},
		State:      protocol.PermissionStateGranted,
		Origin:     "https://example.com",
	})
	if err != nil {
		t.Fatalf("SetPermission: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.PermissionsSetPermission {
		t.Errorf("method = %q, want %q", call.method, protocol.PermissionsSetPermission)
	}

	params := call.params.(protocol.PermissionsSetPermissionParams)
	if params.State != protocol.PermissionStateGranted || params.Origin != "https://example.com" {
		t.Errorf("params = %+v", params)
	}
}

func TestSessionSetGeolocationOverride(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	err := session.SetGeolocationOverride(
		context.Background(),
		protocol.EmulationSetGeolocationOverrideParams{
			Context:     "c1",
			Coordinates: &protocol.EmulationCoordinates{Latitude: 1, Longitude: 2},
		},
	)
	if err != nil {
		t.Fatalf("SetGeolocationOverride: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.EmulationSetGeolocationOverride {
		t.Errorf("method = %q, want %q", m, protocol.EmulationSetGeolocationOverride)
	}
}

func TestSessionSetTimezoneOverride(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	err := session.SetTimezoneOverride(
		context.Background(),
		protocol.EmulationSetTimezoneOverrideParams{
			Context:  "c1",
			Timezone: "Europe/Berlin",
		},
	)
	if err != nil {
		t.Fatalf("SetTimezoneOverride: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.EmulationSetTimezoneOverride {
		t.Errorf("method = %q, want %q", call.method, protocol.EmulationSetTimezoneOverride)
	}

	params := call.params.(protocol.EmulationSetTimezoneOverrideParams)
	if params.Timezone != "Europe/Berlin" {
		t.Errorf("params = %+v", params)
	}
}

func TestSessionCreateUserContext(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.BrowserCreateUserContext: `{"userContext":"uc1"}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.CreateUserContext(context.Background())
	if err != nil {
		t.Fatalf("CreateUserContext: %v", err)
	}

	if got != "uc1" {
		t.Errorf("UserContext = %q, want uc1", got)
	}
}

func TestSessionCloseBrowser(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.CloseBrowser(context.Background()); err != nil {
		t.Fatalf("CloseBrowser: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.BrowserClose {
		t.Errorf("method = %q, want %q", m, protocol.BrowserClose)
	}
}
