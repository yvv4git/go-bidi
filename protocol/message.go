package protocol

import (
	"encoding/json"
	"fmt"
)

// Message kinds carried in the "type" field of BiDi wire frames.
const (
	TypeSuccess = "success"
	TypeError   = "error"
	TypeEvent   = "event"
)

// Message is a decoded BiDi message received from the browser. It is either
// a Response that answers a Command or a browser-initiated Event.
type Message interface {
	isMessage()
}

// Error is a protocol-level failure reported by the browser in an error
// response. Code is one of the WebDriver error codes.
type Error struct {
	Code       string `json:"error"`
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return "bidi: " + e.Code + ": " + e.Message
}

// Response answers a previously sent Command. Result is populated on success
// and Err on failure.
type Response struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Err    *Error          `json:"-"`
}

func (*Response) isMessage() {}

// Event is a browser-initiated notification. Events carry no id and are
// delivered only for the subscriptions enabled on the session.
type Event struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

func (*Event) isMessage() {}

// Decode parses a wire frame into a Response or an Event. The "type" field
// selects the kind; when it is absent the frame is classified by the
// presence of a method.
func Decode(data []byte) (Message, error) {
	var raw struct {
		Type       string          `json:"type"`
		ID         int64           `json:"id"`
		Method     string          `json:"method"`
		Params     json.RawMessage `json:"params"`
		Result     json.RawMessage `json:"result"`
		ErrorCode  string          `json:"error"`
		ErrorText  string          `json:"message"`
		Stacktrace string          `json:"stacktrace"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode message: %w", err)
	}

	if raw.Type == TypeEvent || (raw.Type == "" && raw.Method != "") {
		event := Event{Method: raw.Method, Params: raw.Params}

		return &event, nil
	}

	response := Response{ID: raw.ID, Result: raw.Result}
	if raw.Type == TypeError || raw.ErrorCode != "" {
		response.Err = &Error{
			Code:       raw.ErrorCode,
			Message:    raw.ErrorText,
			Stacktrace: raw.Stacktrace,
		}
	}

	return &response, nil
}
