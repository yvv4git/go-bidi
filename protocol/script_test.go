package protocol

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRemoteValueInterface(t *testing.T) {
	tests := []struct {
		give string
		want any
	}{
		{give: `{"type":"undefined"}`, want: nil},
		{give: `{"type":"null"}`, want: nil},
		{give: `{"type":"string","value":"hello"}`, want: "hello"},
		{give: `{"type":"boolean","value":true}`, want: true},
		{give: `{"type":"number","value":42}`, want: float64(42)},
		{give: `{"type":"number","value":"NaN"}`, want: NumberNaN},
		{give: `{"type":"number","value":"Infinity"}`, want: NumberInfinity},
		{give: `{"type":"object","handle":"h1"}`, want: "h1"},
		{
			give: `{"type":"array","value":[{"type":"number","value":1},{"type":"string","value":"a"}]}`,
			want: []any{float64(1), "a"},
		},
		{
			give: `{"type":"object","value":[["key",{"type":"string","value":"v"}]]}`,
			want: map[string]any{"key": "v"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			var value RemoteValue
			if err := json.Unmarshal([]byte(tt.give), &value); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			got, err := value.Interface()
			if err != nil {
				t.Fatalf("interface: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestScriptEvaluateResult(t *testing.T) {
	frame := `{"type":"success","id":1,"result":{` +
		`"type":"success","realm":"r1",` +
		`"result":{"type":"string","value":"title"}}}`

	msg, err := Decode([]byte(frame))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	resp, ok := msg.(*Response)
	if !ok {
		t.Fatalf("got %T, want *Response", msg)
	}

	var result ScriptEvaluateResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if result.Type != EvaluateSuccess {
		t.Fatalf("type: got %q, want %q", result.Type, EvaluateSuccess)
	}

	value, err := result.Result.Interface()
	if err != nil {
		t.Fatalf("interface: %v", err)
	}

	if value != "title" {
		t.Fatalf("value: got %#v, want %q", value, "title")
	}
}
