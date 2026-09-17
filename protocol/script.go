package protocol

import (
	"encoding/json"
	"fmt"
)

// ScriptTarget identifies the realm or browsing context a script runs in.
// Exactly one of Realm, Context or Sandbox is expected.
type ScriptTarget struct {
	Realm   string `json:"realm,omitempty"`
	Context string `json:"context,omitempty"`
	Sandbox string `json:"sandbox,omitempty"`
}

// SerializationOptions control how script values are serialized.
type SerializationOptions struct {
	MaxDepth          int64  `json:"maxDepth,omitempty"`
	MaxDomDepth       int64  `json:"maxDomDepth,omitempty"`
	IncludeShadowTree string `json:"includeShadowTree,omitempty"`
}

// ScriptEvaluateParams are the parameters of script.evaluate.
type ScriptEvaluateParams struct {
	SerializationOptions *SerializationOptions `json:"serializationOptions,omitempty"`
	Expression           string                `json:"expression"`
	Target               ScriptTarget          `json:"target"`
	ResultOwnership      string                `json:"resultOwnership,omitempty"`
	AwaitPromise         bool                  `json:"awaitPromise,omitempty"`
	UserActivation       bool                  `json:"userActivation,omitempty"`
}

// ScriptEvaluateResult is the result of script.evaluate. Type is either
// EvaluateSuccess or EvaluateException.
type ScriptEvaluateResult struct {
	ExceptionDetails *ExceptionDetails `json:"exceptionDetails,omitempty"`
	Result           *RemoteValue      `json:"result,omitempty"`
	Type             string            `json:"type"`
	Realm            string            `json:"realm,omitempty"`
}

// ExceptionDetails describes an exception raised by script evaluation.
type ExceptionDetails struct {
	Exception    *RemoteValue `json:"exception,omitempty"`
	StackTrace   *StackTrace  `json:"stackTrace,omitempty"`
	Text         string       `json:"text"`
	LineNumber   int64        `json:"lineNumber"`
	ColumnNumber int64        `json:"columnNumber"`
}

// StackTrace is an ordered list of stack frames.
type StackTrace struct {
	CallFrames []StackFrame `json:"callFrames"`
}

// StackFrame is a single frame of a stack trace.
type StackFrame struct {
	FunctionName string `json:"functionName"`
	URL          string `json:"url"`
	LineNumber   int64  `json:"lineNumber"`
	ColumnNumber int64  `json:"columnNumber"`
}

// RemoteValue is a JavaScript value produced by script evaluation. Type
// selects the representation, Value holds the serialized contents and
// Handle is a remote object reference.
type RemoteValue struct {
	Type       string          `json:"type"`
	Value      json.RawMessage `json:"value,omitempty"`
	Handle     string          `json:"handle,omitempty"`
	SharedID   string          `json:"sharedId,omitempty"`
	InternalID string          `json:"internalId,omitempty"`
}

// Interface decodes the remote value into a natural Go value. Remote object
// references are returned as their handle string, arrays as []any and
// objects as map[string]any.
func (v *RemoteValue) Interface() (any, error) {
	if v.Handle != "" {
		return v.Handle, nil
	}

	switch v.Type {
	case RemoteTypeUndefined, RemoteTypeNull:
		return nil, nil
	case RemoteTypeString:
		return v.decodeString()
	case RemoteTypeBoolean:
		return v.decodeBoolean()
	case RemoteTypeNumber, RemoteTypeBigInt:
		return v.decodeNumber()
	case RemoteTypeArray, RemoteTypeSet:
		return v.decodeArray()
	case RemoteTypeObject, RemoteTypeMap:
		return v.decodeObject()
	default:
		return v.Value, nil
	}
}

// decodeString decodes a primitive string value.
func (v *RemoteValue) decodeString() (string, error) {
	var out string
	if err := json.Unmarshal(v.Value, &out); err != nil {
		return "", fmt.Errorf("decode string: %w", err)
	}

	return out, nil
}

// decodeBoolean decodes a primitive boolean value.
func (v *RemoteValue) decodeBoolean() (bool, error) {
	var out bool
	if err := json.Unmarshal(v.Value, &out); err != nil {
		return false, fmt.Errorf("decode boolean: %w", err)
	}

	return out, nil
}

// decodeNumber decodes a number, including the special values NaN, Infinity
// and -Infinity.
func (v *RemoteValue) decodeNumber() (any, error) {
	var out float64
	if err := json.Unmarshal(v.Value, &out); err == nil {
		return out, nil
	}

	var special string
	if err := json.Unmarshal(v.Value, &special); err != nil {
		return nil, fmt.Errorf("decode number: %w", err)
	}

	switch special {
	case NumberNaN, NumberInfinity, NumberNegativeInfinity:
		return special, nil
	default:
		return nil, fmt.Errorf("decode number: %q", special)
	}
}

// decodeArray decodes an array or set value into a slice.
func (v *RemoteValue) decodeArray() ([]any, error) {
	var items []RemoteValue
	if err := json.Unmarshal(v.Value, &items); err != nil {
		return nil, fmt.Errorf("decode array: %w", err)
	}

	out := make([]any, 0, len(items))

	for i := range items {
		item, err := items[i].Interface()
		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}

// decodeObject decodes an object or map value into a map.
func (v *RemoteValue) decodeObject() (map[string]any, error) {
	var entries [][2]json.RawMessage
	if err := json.Unmarshal(v.Value, &entries); err != nil {
		return nil, fmt.Errorf("decode object: %w", err)
	}

	out := make(map[string]any, len(entries))

	for _, entry := range entries {
		var key string
		if err := json.Unmarshal(entry[0], &key); err != nil {
			return nil, fmt.Errorf("decode object key: %w", err)
		}

		var item RemoteValue
		if err := json.Unmarshal(entry[1], &item); err != nil {
			return nil, fmt.Errorf("decode object value: %w", err)
		}

		value, err := item.Interface()
		if err != nil {
			return nil, err
		}

		out[key] = value
	}

	return out, nil
}

// ScriptArgument is a value passed to script.callFunction or a preload
// script. It is either a local value ({type, value}) or a reference to a
// remote value ({handle} / {sharedId} / {internalId}).
type ScriptArgument struct {
	Type       string          `json:"type,omitempty"`
	Value      json.RawMessage `json:"value,omitempty"`
	Handle     string          `json:"handle,omitempty"`
	SharedID   string          `json:"sharedId,omitempty"`
	InternalID string          `json:"internalId,omitempty"`
}

// ScriptCallFunctionParams are the parameters of script.callFunction.
// Arguments holds the values passed to the function and This sets the this
// value.
type ScriptCallFunctionParams struct {
	FunctionDeclaration  string                `json:"functionDeclaration"`
	Arguments            []ScriptArgument      `json:"arguments"`
	SerializationOptions *SerializationOptions `json:"serializationOptions,omitempty"`
	Target               ScriptTarget          `json:"target"`
	ResultOwnership      string                `json:"resultOwnership,omitempty"`
	AwaitPromise         bool                  `json:"awaitPromise,omitempty"`
	This                 *RemoteValue          `json:"this,omitempty"`
	UserActivation       bool                  `json:"userActivation,omitempty"`
}

// ScriptRealmInfo describes a realm of an execution context.
type ScriptRealmInfo struct {
	Realm  string `json:"realm"`
	Origin string `json:"origin"`
	Type   string `json:"type"`
}

// ScriptGetRealmsParams are the parameters of script.getRealms. Contexts
// filters the returned realms to specific browsing contexts.
type ScriptGetRealmsParams struct {
	Contexts []string `json:"contexts,omitempty"`
}

// ScriptGetRealmsResult is the result of script.getRealms.
type ScriptGetRealmsResult struct {
	Realms []ScriptRealmInfo `json:"realms"`
}

// ScriptAddPreloadScriptParams are the parameters of
// script.addPreloadScript. Contexts limits the script to specific browsing
// contexts.
type ScriptAddPreloadScriptParams struct {
	FunctionDeclaration string           `json:"functionDeclaration"`
	Arguments           []ScriptArgument `json:"arguments,omitempty"`
	Contexts            []string         `json:"contexts,omitempty"`
	Sandbox             string           `json:"sandbox,omitempty"`
}

// ScriptAddPreloadScriptResult is the result of script.addPreloadScript.
type ScriptAddPreloadScriptResult struct {
	Script string `json:"script"`
}

// ScriptRemovePreloadScriptParams are the parameters of
// script.removePreloadScript.
type ScriptRemovePreloadScriptParams struct {
	Script string `json:"script"`
}

// ScriptDisownParams are the parameters of script.disown. Handles holds the
// remote object handles to release.
type ScriptDisownParams struct {
	Handles []string     `json:"handles"`
	Target  ScriptTarget `json:"target"`
}

// ScriptMessageParams are the parameters of the script.message event.
type ScriptMessageParams struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
	Source  ScriptRealmInfo `json:"source"`
}
