package protocol

// LogSource identifies where a log entry was produced. Realm is a script
// realm id; Context is the owning browsing context when known.
type LogSource struct {
	Realm   string `json:"realm"`
	Context string `json:"context,omitempty"`
}

// LogEntryParams is the payload of log.entryAdded. Type is one of the log
// entry type constants; Method and StackTrace are populated for console and
// JavaScript entries respectively.
type LogEntryParams struct {
	Type       string      `json:"type"`
	Level      string      `json:"level"`
	Source     LogSource   `json:"source"`
	Text       string      `json:"text"`
	Timestamp  int64       `json:"timestamp"`
	Method     string      `json:"method,omitempty"`
	StackTrace *StackTrace `json:"stackTrace,omitempty"`
}
