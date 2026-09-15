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
