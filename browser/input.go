package browser

import (
	"context"

	"github.com/yvv4git/go-bidi/protocol"
)

// PerformActions dispatches the given input actions in the context.
func (s *Session) PerformActions(
	ctx context.Context,
	params protocol.InputPerformActionsParams,
) error {
	return s.caller.Call(
		ctx,
		protocol.InputPerformActions,
		params,
		nil,
	)
}

// ReleaseActions releases the input state and cancels queued input for the
// given context.
func (s *Session) ReleaseActions(
	ctx context.Context,
	params protocol.InputReleaseActionsParams,
) error {
	return s.caller.Call(
		ctx,
		protocol.InputReleaseActions,
		params,
		nil,
	)
}

// SetFiles sets the files of the file input in the given context. The files
// must be paths readable to the browser.
func (s *Session) SetFiles(ctx context.Context, params protocol.InputSetFilesParams) error {
	return s.caller.Call(ctx, protocol.InputSetFiles, params, nil)
}

// Click clicks at the viewport coordinates x, y.
func (p *Page) Click(ctx context.Context, x, y int64) error {
	source := protocol.InputSourceActions{
		Type: protocol.InputSourceTypePointer,
		ID:   "mouse",
		Actions: []protocol.InputAction{
			{
				Type:   protocol.PointerActionPointerMove,
				X:      x,
				Y:      y,
				Origin: ptrOrigin(protocol.PointerOrigin(protocol.OriginViewport)),
			},
			{Type: protocol.PointerActionPointerDown, Button: 0},
			{Type: protocol.PointerActionPointerUp, Button: 0},
		},
	}

	return p.perform(ctx, source)
}

// Press presses and releases a keyboard key, for example "Enter", "a" or
// "ArrowUp".
func (p *Page) Press(ctx context.Context, key string) error {
	value := keyValue(key)

	source := protocol.InputSourceActions{
		Type: protocol.InputSourceTypeKey,
		ID:   "keyboard",
		Actions: []protocol.InputAction{
			{Type: protocol.KeyActionKeyDown, Value: value},
			{Type: protocol.KeyActionKeyUp, Value: value},
		},
	}

	return p.perform(ctx, source)
}

// keyValue maps a key name to the value expected by the BiDi protocol.
func keyValue(key string) string {
	switch key {
	case "Enter":
		return "\r"
	case "Backspace":
		return "\b"
	case "Tab":
		return "\t"
	case "Escape":
		return "\x1b"
	case "Delete":
		return "\x7f"
	default:
		return key
	}
}

// Type types each rune of text into the focused element.
func (p *Page) Type(ctx context.Context, text string) error {
	actions := make([]protocol.InputAction, 0, 2*len(text))
	for _, r := range text {
		actions = append(actions,
			protocol.InputAction{Type: protocol.KeyActionKeyDown, Value: string(r)},
			protocol.InputAction{Type: protocol.KeyActionKeyUp, Value: string(r)},
		)
	}

	source := protocol.InputSourceActions{
		Type:    protocol.InputSourceTypeKey,
		ID:      "keyboard",
		Actions: actions,
	}

	return p.perform(ctx, source)
}

// Scroll scrolls the page by the given deltas.
func (p *Page) Scroll(ctx context.Context, deltaX, deltaY int64) error {
	source := protocol.InputSourceActions{
		Type: protocol.InputSourceTypeWheel,
		ID:   "wheel",
		Actions: []protocol.InputAction{
			{
				Type:   protocol.WheelActionScroll,
				DX:     deltaX,
				DY:     deltaY,
				Origin: ptrOrigin(protocol.PointerOrigin(protocol.OriginViewport)),
			},
		},
	}

	return p.perform(ctx, source)
}

func (p *Page) perform(ctx context.Context, sources ...protocol.InputSourceActions) error {
	params := protocol.InputPerformActionsParams{
		Context: p.id,
		Actions: sources,
	}

	return p.caller.Call(ctx, protocol.InputPerformActions, params, nil)
}

func ptrOrigin(o protocol.PointerOrigin) *protocol.PointerOrigin {
	return &o
}
