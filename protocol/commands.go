package protocol

// Command names defined by the WebDriver BiDi protocol, grouped by module.

// session module commands.
const (
	// SessionStatus reports whether the remote end is ready.
	SessionStatus = "session.status"
	// SessionNew creates a new session with the requested capabilities.
	SessionNew = "session.new"
	// SessionEnd ends the current session.
	SessionEnd = "session.end"
	// SessionSubscribe enables a subscription to the given events.
	SessionSubscribe = "session.subscribe"
	// SessionUnsubscribe removes a subscription.
	SessionUnsubscribe = "session.unsubscribe"
)

// browser module commands.
const (
	// BrowserClose closes the browser.
	BrowserClose = "browser.close"
	// BrowserCreateUserContext creates a new user context.
	BrowserCreateUserContext = "browser.createUserContext"
)

// browsingContext module commands.
const (
	// BrowsingContextActivate brings a browsing context to the foreground.
	BrowsingContextActivate = "browsingContext.activate"
	// BrowsingContextCaptureScreenshot captures a PNG screenshot.
	BrowsingContextCaptureScreenshot = "browsingContext.captureScreenshot"
	// BrowsingContextClose closes a browsing context.
	BrowsingContextClose = "browsingContext.close"
	// BrowsingContextCreate creates a new browsing context.
	BrowsingContextCreate = "browsingContext.create"
	// BrowsingContextGetTree returns the tree of browsing contexts.
	BrowsingContextGetTree = "browsingContext.getTree"
	// BrowsingContextHandleUserPrompt handles an open user prompt.
	BrowsingContextHandleUserPrompt = "browsingContext.handleUserPrompt"
	// BrowsingContextLocateNodes finds nodes matching a locator.
	BrowsingContextLocateNodes = "browsingContext.locateNodes"
	// BrowsingContextNavigate navigates a browsing context to a URL.
	BrowsingContextNavigate = "browsingContext.navigate"
	// BrowsingContextPrint prints a browsing context to PDF.
	BrowsingContextPrint = "browsingContext.print"
	// BrowsingContextReload reloads a browsing context.
	BrowsingContextReload = "browsingContext.reload"
	// BrowsingContextSetViewport sets the viewport of a browsing context.
	BrowsingContextSetViewport = "browsingContext.setViewport"
	// BrowsingContextTraverseHistory moves through the session history.
	BrowsingContextTraverseHistory = "browsingContext.traverseHistory"
)

// script module commands.
const (
	// ScriptAddPreloadScript adds a script run before page scripts.
	ScriptAddPreloadScript = "script.addPreloadScript"
	// ScriptCallFunction calls a JavaScript function.
	ScriptCallFunction = "script.callFunction"
	// ScriptDisown removes a remote object handle.
	ScriptDisown = "script.disown"
	// ScriptEvaluate evaluates a JavaScript expression.
	ScriptEvaluate = "script.evaluate"
	// ScriptGetRealms lists the realms of a browsing context.
	ScriptGetRealms = "script.getRealms"
	// ScriptRemovePreloadScript removes a preload script.
	ScriptRemovePreloadScript = "script.removePreloadScript"
)

// input module commands.
const (
	// InputPerformActions dispatches a sequence of input actions.
	InputPerformActions = "input.performActions"
	// InputReleaseActions releases all held input state.
	InputReleaseActions = "input.releaseActions"
	// InputSetFiles sets files on a file input element.
	InputSetFiles = "input.setFiles"
)

// network module commands.
const (
	// NetworkAddIntercept registers a network intercept.
	NetworkAddIntercept = "network.addIntercept"
	// NetworkContinueRequest continues an intercepted request.
	NetworkContinueRequest = "network.continueRequest"
	// NetworkContinueResponse continues an intercepted response.
	NetworkContinueResponse = "network.continueResponse"
	// NetworkContinueWithAuth responds to an authentication challenge.
	NetworkContinueWithAuth = "network.continueWithAuth"
	// NetworkFailRequest fails an intercepted request.
	NetworkFailRequest = "network.failRequest"
	// NetworkProvideResponse provides a synthetic response.
	NetworkProvideResponse = "network.provideResponse"
	// NetworkRemoveIntercept removes a network intercept.
	NetworkRemoveIntercept = "network.removeIntercept"
	// NetworkSetCacheBehavior sets the cache behavior of the context.
	NetworkSetCacheBehavior = "network.setCacheBehavior"
)

// storage module commands.
const (
	// StorageDeleteCookies deletes cookies.
	StorageDeleteCookies = "storage.deleteCookies"
	// StorageGetCookies returns cookies.
	StorageGetCookies = "storage.getCookies"
	// StorageSetCookie sets a cookie.
	StorageSetCookie = "storage.setCookie"
)

// emulation module commands.
const (
	// EmulationSetGeolocationOverride overrides the geolocation.
	EmulationSetGeolocationOverride = "emulation.setGeolocationOverride"
	// EmulationSetLocaleOverride overrides the locale.
	EmulationSetLocaleOverride = "emulation.setLocaleOverride"
	// EmulationSetTimezoneOverride overrides the timezone.
	EmulationSetTimezoneOverride = "emulation.setTimezoneOverride"
	// EmulationSetUserAgentOverride overrides the user agent.
	EmulationSetUserAgentOverride = "emulation.setUserAgentOverride"
)

// permissions module commands.
const (
	// PermissionsSetPermission grants or denies a permission.
	PermissionsSetPermission = "permissions.setPermission"
)
