package protocol

import (
	"encoding/json"
	"testing"
)

func TestStorageGetCookies(t *testing.T) {
	filter := StorageCookieFilter{Name: strPtr("sid"), SameSite: CookieSameSiteLax}

	params := StorageGetCookiesParams{Filter: &filter}
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"filter":{"name":"sid","sameSite":"Lax"}}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	var result StorageGetCookiesResult
	raw := []byte(`{"cookies":[{"name":"sid","value":{"type":"string","value":"1"},` +
		`"domain":"example.com","path":"/","size":3,"expires":0,` +
		`"httpOnly":true,"secure":true,"sameSite":"Lax"}]}`)
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if len(result.Cookies) != 1 || !result.Cookies[0].HTTPOnly ||
		result.Cookies[0].Value.Value != "1" {
		t.Errorf("result = %+v", result.Cookies)
	}
}

func TestStorageSetCookieParams(t *testing.T) {
	params := StorageSetCookieParams{
		Cookie: StorageCookie{
			Name:     "a",
			Value:    NetworkBytesValue{Type: BytesValueString, Value: "b"},
			Domain:   "example.com",
			Path:     "/",
			SameSite: CookieSameSiteStrict,
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"cookie":{"name":"a","value":{"type":"string","value":"b"},` +
		`"domain":"example.com","path":"/","size":0,"expires":0,` +
		`"httpOnly":false,"secure":false,"sameSite":"Strict"}}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestStorageDeleteCookiesParams(t *testing.T) {
	params := StorageDeleteCookiesParams{
		Filter: &StorageCookieFilter{Domain: strPtr("example.com")},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(data); got != `{"filter":{"domain":"example.com"}}` {
		t.Errorf("got %s", got)
	}
}
