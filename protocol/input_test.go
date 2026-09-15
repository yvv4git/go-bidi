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
		`"actions":[{"type":"keyDown","value":"a"},{"type":"keyUp","value":"a"},` +
		`{"type":"pause","duration":100}]}]}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestInputPointerOrigin(t *testing.T) {
	action := InputAction{
		Type: PointerActionPointerMove,
		X:    10,
		Y:    20,
		Origin: &PointerOrigin{
			Type:    OriginElement,
			Element: &ElementReference{SharedID: "e1"},
		},
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"type":"pointerMove","x":10,"y":20,` +
		`"origin":{"type":"element","element":{"sharedId":"e1"}}}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestInputWheelAction(t *testing.T) {
	action := InputAction{
		Type:   WheelActionScroll,
		DX:     -50,
		DY:     100,
		Origin: &PointerOrigin{Type: OriginViewport},
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"type":"scroll","dx":-50,"dy":100,"origin":{"type":"viewport"}}`
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
