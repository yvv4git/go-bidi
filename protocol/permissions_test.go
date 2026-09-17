package protocol

import (
	"encoding/json"
	"testing"
)

func TestPermissionsSetPermissionParams(t *testing.T) {
	unrestricted := false
	params := PermissionsSetPermissionParams{
		Descriptor: PermissionsDescriptor{
			Name:            "notifications",
			UserVisibleOnly: &unrestricted,
		},
		State:  PermissionStateGranted,
		Origin: "https://example.com",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"descriptor":{"name":"notifications","userVisibleOnly":false},` +
		`"state":"granted","origin":"https://example.com"}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
