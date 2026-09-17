package protocol

import (
	"encoding/json"
	"testing"
)

func TestInputPerformActionsParams(t *testing.T) {
	params := InputPerformActionsParams{
		Context: "c1",
		Actions: []InputSourceActions{
			{
				Type: InputSourceTypeKey,
				ID:   "keyboard",
				Actions: []InputAction{
					{Type: KeyActionKeyDown, Value: "a"},
					{Type: KeyActionKeyUp, Value: "a"},
					{Type: KeyActionPause, Duration: 100},
				},
			},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"context":"c1","actions":[{"type":"key","id":"keyboard",` +
		`"actions":[{"type":"keyDown","value":"a","button":0},{"type":"keyUp","value":"a","button":0},` +
		`{"type":"pause","button":0,"duration":100}]}]}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestInputPointerOrigin(t *testing.T) {
	action := InputAction{
		Type:   PointerActionPointerMove,
		X:      10,
		Y:      20,
		Origin: ptrOrigin(OriginPointer),
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"type":"pointerMove","button":0,"x":10,"y":20,"origin":"pointer"}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestInputWheelAction(t *testing.T) {
	action := InputAction{
		Type:   WheelActionScroll,
		DX:     -50,
		DY:     100,
		Origin: ptrOrigin(PointerOrigin(OriginViewport)),
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"type":"scroll","button":0,"dx":-50,"dy":100,"origin":"viewport"}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestInputReleaseActionsParams(t *testing.T) {
	params := InputReleaseActionsParams{Context: "c1"}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(data); got != `{"context":"c1"}` {
		t.Errorf("got %s, want %s", got, `{"context":"c1"}`)
	}
}

func ptrOrigin(o PointerOrigin) *PointerOrigin {
	return &o
}

func TestInputSetFilesParams(t *testing.T) {
	params := InputSetFilesParams{Context: "c1", Files: []string{"a.txt"}}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"context":"c1","files":["a.txt"]}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
