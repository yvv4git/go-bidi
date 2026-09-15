package protocol

// Browsing context types accepted by browsingContext.create.
const (
	// ContextTypeTab is a new tab.
	ContextTypeTab = "tab"
	// ContextTypeWindow is a new window.
	ContextTypeWindow = "window"
)

// Readiness states used when waiting for a navigation to settle.
const (
	// ReadinessNone does not wait.
	ReadinessNone = "none"
	// ReadinessInteractive waits for the document to become interactive.
	ReadinessInteractive = "interactive"
	// ReadinessComplete waits for the document to finish loading.
	ReadinessComplete = "complete"
)

// Screenshot origins accepted by browsingContext.captureScreenshot.
const (
	// OriginViewport captures only the viewport.
	OriginViewport = "viewport"
	// OriginDocument captures the full document.
	OriginDocument = "document"
)

// Input source types accepted by input.performActions.
const (
	// InputSourceTypeKey is a keyboard input source.
	InputSourceTypeKey = "key"
	// InputSourceTypePointer is a pointer input source.
	InputSourceTypePointer = "pointer"
	// InputSourceTypeWheel is a wheel input source.
	InputSourceTypeWheel = "wheel"
	// InputSourceTypeNone introduces a pause tick to synchronize sources.
	InputSourceTypeNone = "none"
)

// Action types of a key input source.
const (
	// KeyActionKeyDown presses a key.
	KeyActionKeyDown = "keyDown"
	// KeyActionKeyUp releases a key.
	KeyActionKeyUp = "keyUp"
	// KeyActionPause waits for the given duration.
	KeyActionPause = "pause"
)

// Action types of a pointer input source.
const (
	// PointerActionPointerDown presses a pointer button.
	PointerActionPointerDown = "pointerDown"
	// PointerActionPointerUp releases a pointer button.
	PointerActionPointerUp = "pointerUp"
	// PointerActionPointerMove moves the pointer to the given coordinates.
	PointerActionPointerMove = "pointerMove"
	// PointerActionPointerCancel cancels the pointer input.
	PointerActionPointerCancel = "pointerCancel"
	// PointerActionPause waits for the given duration.
	PointerActionPause = "pause"
)

// Action types of a wheel input source.
const (
	// WheelActionScroll scrolls by the given deltas.
	WheelActionScroll = "scroll"
	// WheelActionPause waits for the given duration.
	WheelActionPause = "pause"
)

// Pointer origins accepted by pointer and wheel actions.
const (
	// OriginPointer uses the current pointer position.
	OriginPointer = "pointer"
	// OriginElement uses the referenced element.
	OriginElement = "element"
)

// Result ownership modes accepted by script evaluation commands.
const (
	// ResultOwnershipNone does not retain the result.
	ResultOwnershipNone = "none"
	// ResultOwnershipRoot retains the root of the result.
	ResultOwnershipRoot = "root"
)

// Kinds of a script evaluation result.
const (
	// EvaluateSuccess marks a successful evaluation.
	EvaluateSuccess = "success"
	// EvaluateException marks an evaluation that threw.
	EvaluateException = "exception"
)

// Remote value types carried in the "type" field of a script.RemoteValue.
const (
	// RemoteTypeUndefined is the undefined value.
	RemoteTypeUndefined = "undefined"
	// RemoteTypeNull is the null value.
	RemoteTypeNull = "null"
	// RemoteTypeString is a string value.
	RemoteTypeString = "string"
	// RemoteTypeBoolean is a boolean value.
	RemoteTypeBoolean = "boolean"
	// RemoteTypeNumber is a number value.
	RemoteTypeNumber = "number"
	// RemoteTypeBigInt is a bigint value.
	RemoteTypeBigInt = "bigint"
	// RemoteTypeArray is an array value.
	RemoteTypeArray = "array"
	// RemoteTypeSet is a set value.
	RemoteTypeSet = "set"
	// RemoteTypeObject is an object value.
	RemoteTypeObject = "object"
	// RemoteTypeMap is a map value.
	RemoteTypeMap = "map"
)

// Special numeric representations used by script.RemoteValue.
const (
	// NumberNaN is the not-a-number value.
	NumberNaN = "NaN"
	// NumberInfinity is positive infinity.
	NumberInfinity = "Infinity"
	// NumberNegativeInfinity is negative infinity.
	NumberNegativeInfinity = "-Infinity"
)
