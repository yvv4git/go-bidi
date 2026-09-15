package protocol

// PointerOrigin selects the coordinate origin of a pointer or wheel action.
// Element is only used when Type is OriginElement.
type PointerOrigin struct {
	Type    string            `json:"type"`
	Element *ElementReference `json:"element,omitempty"`
}

// ElementReference points to an element by its shared id.
type ElementReference struct {
	SharedID string `json:"sharedId"`
}

// InputAction is a single input action. Type selects the action kind and is
// one of the input action constants; the remaining fields are populated
// according to the input source and action.
type InputAction struct {
	Type     string         `json:"type"`
	Value    string         `json:"value,omitempty"`
	Button   int            `json:"button,omitempty"`
	X        int64          `json:"x,omitempty"`
	Y        int64          `json:"y,omitempty"`
	DX       int64          `json:"dx,omitempty"`
	DY       int64          `json:"dy,omitempty"`
	Duration int64          `json:"duration,omitempty"`
	Origin   *PointerOrigin `json:"origin,omitempty"`
}

// InputSourceActions lists the actions performed by a single input source.
// Type is one of the input source type constants.
type InputSourceActions struct {
	Type    string        `json:"type"`
	ID      string        `json:"id"`
	Actions []InputAction `json:"actions"`
}

// InputPerformActionsParams are the parameters of input.performActions.
type InputPerformActionsParams struct {
	Context string               `json:"context"`
	Actions []InputSourceActions `json:"actions"`
}

// InputReleaseActionsParams are the parameters of input.releaseActions.
type InputReleaseActionsParams struct {
	Context string `json:"context"`
}

// InputSetFilesParams are the parameters of input.setFiles. Files holds
// paths to files the browser can read.
type InputSetFilesParams struct {
	Context string   `json:"context"`
	Files   []string `json:"files"`
}
