package protocol

// BrowsingContextCreateParams are the parameters of browsingContext.create.
type BrowsingContextCreateParams struct {
	Type             string `json:"type,omitempty"`
	ReferenceContext string `json:"referenceContext,omitempty"`
	UserContext      string `json:"userContext,omitempty"`
	Background       bool   `json:"background,omitempty"`
}

// BrowsingContextCreateResult is the result of browsingContext.create. The
// Context id identifies the new browsing context.
type BrowsingContextCreateResult struct {
	Context string `json:"context"`
}

// BrowsingContextNavigateParams are the parameters of
// browsingContext.navigate. Wait selects the readiness state to await.
type BrowsingContextNavigateParams struct {
	Context string `json:"context"`
	URL     string `json:"url"`
	Wait    string `json:"wait,omitempty"`
}

// BrowsingContextNavigateResult is the result of browsingContext.navigate.
type BrowsingContextNavigateResult struct {
	Navigation string `json:"navigation"`
	URL        string `json:"url"`
}

// BrowsingContextReloadParams are the parameters of browsingContext.reload.
type BrowsingContextReloadParams struct {
	Context     string `json:"context"`
	IgnoreCache bool   `json:"ignoreCache,omitempty"`
	Wait        string `json:"wait,omitempty"`
}

// BrowsingContextReloadResult is the result of browsingContext.reload.
type BrowsingContextReloadResult struct {
	Navigation string `json:"navigation"`
	URL        string `json:"url"`
}

// Viewport describes the size of a browsing context viewport.
type Viewport struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// BrowsingContextSetViewportParams are the parameters of
// browsingContext.setViewport. Viewport and DevicePixelRatio are optional;
// when both are nil the override is cleared.
type BrowsingContextSetViewportParams struct {
	Context          string    `json:"context"`
	Viewport         *Viewport `json:"viewport,omitempty"`
	DevicePixelRatio *float64  `json:"devicePixelRatio,omitempty"`
}

// BrowsingContextInfo describes a browsing context and its descendants.
type BrowsingContextInfo struct {
	Children       []BrowsingContextInfo `json:"children"`
	Context        string                `json:"context"`
	URL            string                `json:"url"`
	Parent         *string               `json:"parent,omitempty"`
	UserContext    string                `json:"userContext,omitempty"`
	OriginalOpener *string               `json:"originalOpener,omitempty"`
}

// BrowsingContextGetTreeParams are the parameters of
// browsingContext.getTree. Root scopes the returned tree to one context.
type BrowsingContextGetTreeParams struct {
	MaxDepth int    `json:"maxDepth,omitempty"`
	Root     string `json:"root,omitempty"`
}

// BrowsingContextGetTreeResult is the result of browsingContext.getTree.
type BrowsingContextGetTreeResult struct {
	Contexts []BrowsingContextInfo `json:"contexts"`
}

// BrowsingContextCloseParams are the parameters of browsingContext.close.
type BrowsingContextCloseParams struct {
	Context      string `json:"context"`
	PromptUnload bool   `json:"promptUnload,omitempty"`
}

// BrowsingContextTraverseHistoryParams are the parameters of
// browsingContext.traverseHistory. Delta is the number of history entries to
// move, positive forward and negative backward.
type BrowsingContextTraverseHistoryParams struct {
	Context string `json:"context"`
	Delta   int    `json:"delta"`
}

// BrowsingContextCaptureScreenshotParams are the parameters of
// browsingContext.captureScreenshot. CaptureBeyondViewport captures the whole
// document when set.
type BrowsingContextCaptureScreenshotParams struct {
	Context               string `json:"context"`
	Origin                string `json:"origin,omitempty"`
	CaptureBeyondViewport bool   `json:"captureBeyondViewport,omitempty"`
}

// BrowsingContextCaptureScreenshotResult is the result of
// browsingContext.captureScreenshot. Data is a base64-encoded PNG image.
type BrowsingContextCaptureScreenshotResult struct {
	Data string `json:"data"`
}

// NavigationInfo describes the navigation state reported by the
// navigation-related browsingContext events.
type NavigationInfo struct {
	Context    string `json:"context"`
	Navigation string `json:"navigation,omitempty"`
	Timestamp  int64  `json:"timestamp"`
	URL        string `json:"url"`
}
