package protocol

import (
	"encoding/json"
	"testing"
)

func TestLogEntryDecode(t *testing.T) {
	params := LogEntryParams{
		Type:      LogTypeConsole,
		Level:     LogLevelInfo,
		Source:    LogSource{Realm: "r1", Context: "c1"},
		Text:      "hello",
		Timestamp: 5,
		Method:    "log",
	}

	frame := eventFrame(t, LogEntryAdded, params)
	event := decodeEvent(t, frame)

	var got LogEntryParams
	if err := json.Unmarshal(event.Params, &got); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if got.Type != LogTypeConsole || got.Level != LogLevelInfo || got.Method != "log" {
		t.Errorf("entry = %+v", got)
	}

	if got.Source.Context != "c1" || got.Source.Realm != "r1" {
		t.Errorf("Source = %+v, want c1/r1", got.Source)
	}
}

func TestLogJavascriptEntryDecode(t *testing.T) {
	params := LogEntryParams{
		Type:      LogTypeJavascript,
		Level:     LogLevelError,
		Source:    LogSource{Realm: "r1", Context: "c1"},
		Text:      "boom",
		Timestamp: 9,
		StackTrace: &StackTrace{
			CallFrames: []StackFrame{
				{FunctionName: "f", URL: "https://example.com/app.js", LineNumber: 1, ColumnNumber: 2},
			},
		},
	}

	frame := eventFrame(t, LogEntryAdded, params)
	event := decodeEvent(t, frame)

	var got LogEntryParams
	if err := json.Unmarshal(event.Params, &got); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if got.Type != LogTypeJavascript || got.StackTrace == nil {
		t.Fatalf("entry = %+v", got)
	}

	frame0 := got.StackTrace.CallFrames[0]
	if frame0.URL != "https://example.com/app.js" || frame0.LineNumber != 1 {
		t.Errorf("frame = %+v", frame0)
	}
}
