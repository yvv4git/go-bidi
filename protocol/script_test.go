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

func TestScriptCallFunctionParams(t *testing.T) {
	params := ScriptCallFunctionParams{
		FunctionDeclaration: "() => 42",
		Arguments: []ScriptArgument{
			{Type: RemoteTypeNumber, Value: json.RawMessage(`7`)},
		},
		Target:       ScriptTarget{Context: "c1"},
		AwaitPromise: true,
	}

	raw, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ScriptCallFunctionParams
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if decoded.FunctionDeclaration != "() => 42" {
		t.Errorf("FunctionDeclaration = %q", decoded.FunctionDeclaration)
	}

	if len(decoded.Arguments) != 1 || decoded.Arguments[0].Type != RemoteTypeNumber {
		t.Errorf("Arguments = %+v", decoded.Arguments)
	}
}

func TestScriptGetRealmsResult(t *testing.T) {
	frame := `{"type":"success","id":1,"result":{"realms":[{` +
		`"realm":"r1","origin":"https://example.com","type":"window"}]}}`

	msg, err := Decode([]byte(frame))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	resp, ok := msg.(*Response)
	if !ok {
		t.Fatalf("got %T, want *Response", msg)
	}

	var result ScriptGetRealmsResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if len(result.Realms) != 1 {
		t.Fatalf("realms: got %d, want 1", len(result.Realms))
	}

	if got := result.Realms[0]; got.Realm != "r1" ||
		got.Origin != "https://example.com" || got.Type != "window" {
		t.Errorf("realm = %+v", got)
	}
}

func TestScriptAddPreloadScriptResult(t *testing.T) {
	frame := `{"type":"success","id":1,"result":{"script":"s1"}}`

	var resp Response
	if err := json.Unmarshal([]byte(frame), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	var result ScriptAddPreloadScriptResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if result.Script != "s1" {
		t.Errorf("Script = %q, want %q", result.Script, "s1")
	}
}

func TestScriptDisownParams(t *testing.T) {
	params := ScriptDisownParams{
		Handles: []string{"h1"},
		Target:  ScriptTarget{Realm: "r1"},
	}

	raw, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ScriptDisownParams
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(decoded.Handles) != 1 || decoded.Handles[0] != "h1" {
		t.Errorf("Handles = %+v", decoded.Handles)
	}

	if decoded.Target.Realm != "r1" {
		t.Errorf("Target = %+v", decoded.Target)
	}
}

func TestScriptMessageParams(t *testing.T) {
	frame := `{"type":"event","method":"script.message","params":{` +
		`"channel":"chan","data":{"type":"string","value":"hi"},` +
		`"source":{"realm":"r1","origin":"https://example.com","type":"window"}}}`

	msg, err := Decode([]byte(frame))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	event, ok := msg.(*Event)
	if !ok {
		t.Fatalf("got %T, want *Event", msg)
	}

	var params ScriptMessageParams
	if err := json.Unmarshal(event.Params, &params); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if params.Channel != "chan" {
		t.Errorf("Channel = %q, want %q", params.Channel, "chan")
	}

	if params.Source.Realm != "r1" {
		t.Errorf("Source = %+v", params.Source)
	}
}
