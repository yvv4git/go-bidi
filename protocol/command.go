package protocol

import (
	"encoding/json"
	"fmt"
)

// Command is a client-to-browser request. The ID is chosen by the client and
// used to match the eventual Response. Params holds JSON-encoded command
// parameters.
type Command struct {
	ID     int64           `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// NewCommand builds a Command, encoding params as JSON. A nil params is
// omitted from the wire frame, as for commands that take no parameters.
func NewCommand(id int64, method string, params any) (*Command, error) {
	cmd := &Command{ID: id, Method: method}

	if params == nil {
		return cmd, nil
	}

	if raw, ok := params.(json.RawMessage); ok {
		cmd.Params = raw

		return cmd, nil
	}

	raw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("encode params: %w", err)
	}

	cmd.Params = raw

	return cmd, nil
}

// Encode serialises the command to a wire JSON frame.
func (c *Command) Encode() ([]byte, error) {
	return json.Marshal(c)
}
