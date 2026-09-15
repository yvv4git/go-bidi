package protocol

import (
	"encoding/json"
	"testing"
)

func TestBrowserCreateUserContextResult(t *testing.T) {
	var result BrowserCreateUserContextResult
	if err := json.Unmarshal([]byte(`{"userContext":"uc1"}`), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if result.UserContext != "uc1" {
		t.Errorf("UserContext = %q, want uc1", result.UserContext)
	}
}
