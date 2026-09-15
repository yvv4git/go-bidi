package browser

import (
	"context"

	"github.com/yvv4git/go-bidi/protocol"
)

// Evaluate evaluates a JavaScript expression in the target realm or context.
func (s *Session) Evaluate(
	ctx context.Context,
	params protocol.ScriptEvaluateParams,
) (*protocol.ScriptEvaluateResult, error) {
	var result protocol.ScriptEvaluateResult
	if err := s.caller.Call(ctx, protocol.ScriptEvaluate, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CallFunction calls a JavaScript function in the target realm or context.
func (s *Session) CallFunction(
	ctx context.Context,
	params protocol.ScriptCallFunctionParams,
) (*protocol.ScriptEvaluateResult, error) {
	var result protocol.ScriptEvaluateResult
	if err := s.caller.Call(ctx, protocol.ScriptCallFunction, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRealms returns the realms of the session, optionally filtered by
// context.
func (s *Session) GetRealms(
	ctx context.Context,
	params protocol.ScriptGetRealmsParams,
) (*protocol.ScriptGetRealmsResult, error) {
	var result protocol.ScriptGetRealmsResult
	if err := s.caller.Call(ctx, protocol.ScriptGetRealms, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddPreloadScript registers a script that runs before page scripts.
func (s *Session) AddPreloadScript(
	ctx context.Context,
	params protocol.ScriptAddPreloadScriptParams,
) (*protocol.ScriptAddPreloadScriptResult, error) {
	var result protocol.ScriptAddPreloadScriptResult
	if err := s.caller.Call(ctx, protocol.ScriptAddPreloadScript, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemovePreloadScript removes a previously added preload script.
func (s *Session) RemovePreloadScript(
	ctx context.Context,
	params protocol.ScriptRemovePreloadScriptParams,
) error {
	return s.caller.Call(ctx, protocol.ScriptRemovePreloadScript, params, nil)
}

// Disown releases remote object handles in the target realm or context.
func (s *Session) Disown(
	ctx context.Context,
	params protocol.ScriptDisownParams,
) error {
	return s.caller.Call(ctx, protocol.ScriptDisown, params, nil)
}
