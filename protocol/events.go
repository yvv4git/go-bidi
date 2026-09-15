package protocol

// Event names defined by the WebDriver BiDi protocol, grouped by module.

// browsingContext module events.
const (
	// BrowsingContextContextCreated is emitted when a context is created.
	BrowsingContextContextCreated = "browsingContext.contextCreated"
	// BrowsingContextContextDestroyed is emitted when a context is destroyed.
	BrowsingContextContextDestroyed = "browsingContext.contextDestroyed"
	// BrowsingContextDomContentLoaded is emitted on DOMContentLoaded.
	BrowsingContextDomContentLoaded = "browsingContext.domContentLoaded"
	// BrowsingContextDownloadEnd is emitted when a download ends.
	BrowsingContextDownloadEnd = "browsingContext.downloadEnd"
	// BrowsingContextDownloadWillBegin is emitted when a download begins.
	BrowsingContextDownloadWillBegin = "browsingContext.downloadWillBegin"
	// BrowsingContextFragmentNavigated is emitted on fragment navigation.
	BrowsingContextFragmentNavigated = "browsingContext.fragmentNavigated"
	// BrowsingContextHistoryUpdated is emitted when history is updated.
	BrowsingContextHistoryUpdated = "browsingContext.historyUpdated"
	// BrowsingContextLoad is emitted on page load.
	BrowsingContextLoad = "browsingContext.load"
	// BrowsingContextNavigationAborted is emitted when navigation aborts.
	BrowsingContextNavigationAborted = "browsingContext.navigationAborted"
	// BrowsingContextNavigationCommitted is emitted when navigation commits.
	BrowsingContextNavigationCommitted = "browsingContext.navigationCommitted"
	// BrowsingContextNavigationFailed is emitted when navigation fails.
	BrowsingContextNavigationFailed = "browsingContext.navigationFailed"
	// BrowsingContextNavigationStarted is emitted when navigation starts.
	BrowsingContextNavigationStarted = "browsingContext.navigationStarted"
	// BrowsingContextUserPromptClosed is emitted when a prompt is closed.
	BrowsingContextUserPromptClosed = "browsingContext.userPromptClosed"
	// BrowsingContextUserPromptOpened is emitted when a prompt is opened.
	BrowsingContextUserPromptOpened = "browsingContext.userPromptOpened"
)

// script module events.
const (
	// ScriptMessage is emitted for a channel message.
	ScriptMessage = "script.message"
	// ScriptRealmCreated is emitted when a realm is created.
	ScriptRealmCreated = "script.realmCreated"
	// ScriptRealmDestroyed is emitted when a realm is destroyed.
	ScriptRealmDestroyed = "script.realmDestroyed"
)

// network module events.
const (
	// NetworkAuthRequired is emitted when authentication is required.
	NetworkAuthRequired = "network.authRequired"
	// NetworkBeforeRequestSent is emitted before a request is sent.
	NetworkBeforeRequestSent = "network.beforeRequestSent"
	// NetworkFetchError is emitted when a fetch fails.
	NetworkFetchError = "network.fetchError"
	// NetworkResponseCompleted is emitted when a response completes.
	NetworkResponseCompleted = "network.responseCompleted"
	// NetworkResponseStarted is emitted when a response starts.
	NetworkResponseStarted = "network.responseStarted"
)

// log module events.
const (
	// LogEntryAdded is emitted when a log entry is added.
	LogEntryAdded = "log.entryAdded"
)

// input module events.
const (
	// InputFileDialogOpened is emitted when a file dialog opens.
	InputFileDialogOpened = "input.fileDialogOpened"
)
