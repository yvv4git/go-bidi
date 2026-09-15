package protocol

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNetworkBeforeRequestSentDecode(t *testing.T) {
	params := NetworkBeforeRequestSentParams{
		Context:    "c1",
		Navigation: strPtr("n1"),
		Timestamp:  100,
		Request: NetworkRequestData{
			Request: "r1",
			URL:     "https://example.com/a",
			Method:  "GET",
			Headers: []NetworkHeader{
				{Name: "Host", Value: NetworkBytesValue{Type: BytesValueString, Value: "example.com"}},
			},
			Cookies:     []NetworkCookie{},
			HeadersSize: 100,
			Destination: "document",
			Initiator:   NetworkInitiator{Type: "parser"},
			Timing: NetworkTimingInfo{
				OriginTime:    0,
				RequestTime:   0,
				FetchStart:    int64Ptr(1),
				DNSStart:      int64Ptr(1),
				DNSEnd:        int64Ptr(1),
				ConnectStart:  int64Ptr(2),
				ConnectEnd:    int64Ptr(3),
				TLSStart:      int64Ptr(3),
				RequestStart:  int64Ptr(4),
				ResponseStart: int64Ptr(5),
			},
		},
	}

	frame := eventFrame(t, NetworkBeforeRequestSent, params)
	event := decodeEvent(t, frame)

	var got NetworkBeforeRequestSentParams
	if err := json.Unmarshal(event.Params, &got); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if got.Context != "c1" {
		t.Errorf("Context = %q, want c1", got.Context)
	}

	req := got.Request
	if req.Request != "r1" || req.URL != "https://example.com/a" || req.Method != "GET" {
		t.Errorf("Request = %+v", req)
	}

	if req.BodySize != nil {
		t.Errorf("BodySize = %v, want nil", *req.BodySize)
	}

	if got := req.Timing.DNSStart; got == nil || *got != 1 {
		t.Errorf("DNSStart = %v, want 1", got)
	}

	if got := req.Timing.ConnectStart; got == nil || *got != 2 {
		t.Errorf("ConnectStart = %v, want 2", got)
	}

	if got.Navigation == nil || *got.Navigation != "n1" {
		t.Errorf("Navigation = %v, want n1", got.Navigation)
	}
}

func TestNetworkResponseCompletedDecode(t *testing.T) {
	frame := `{"type":"event","method":"network.responseCompleted","params":{` +
		`"context":"c1","isBlocked":false,"redirectCount":0,"timestamp":100,` +
		`"request":{"request":"r1","url":"https://example.com/a","method":"GET",` +
		`"headers":[],"cookies":[],"headersSize":0,"bodySize":1,` +
		`"destination":"document","initiator":{"type":"other"},` +
		`"timing":{"originTime":0,"requestTime":0,"fetchStart":1,` +
		`"redirectStart":1,"redirectEnd":2,"dnsStart":null,"dnsEnd":null,` +
		`"connectStart":null,"connectEnd":null,"tlsStart":null,"requestStart":4,` +
		`"responseStart":5,"responseEnd":6}},` +
		`"response":{"url":"https://example.com/a","protocol":"http/1.1","status":200,` +
		`"statusText":"OK","fromCache":false,` +
		`"headers":[{"name":"Content-Type",` +
		`"value":{"type":"string","value":"text/html"}}],` +
		`"mimeType":"text/html","knownLength":300,"bodySize":300,"cookies":[]}}}`

	event := decodeEvent(t, frame)
	var params NetworkResponseCompletedParams
	if err := json.Unmarshal(event.Params, &params); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if params.Response.Status != 200 || params.Response.Protocol != "http/1.1" {
		t.Errorf("Response = %+v", params.Response)
	}

	if got := params.Response.KnownLength; got == nil || *got != 300 {
		t.Errorf("KnownLength = %v, want 300", got)
	}

	if got := params.Request.Timing.RedirectStart; got == nil || *got != 1 {
		t.Errorf("RedirectStart = %v, want 1", got)
	}

	if len(params.Response.Headers) != 1 ||
		params.Response.Headers[0].Value.Value != "text/html" {
		t.Errorf("Headers = %+v", params.Response.Headers)
	}

	if len(params.Response.Cookies) != 0 {
		t.Errorf("Cookies = %+v, want none", params.Response.Cookies)
	}
}

func TestNetworkFetchErrorDecode(t *testing.T) {
	params := NetworkFetchErrorParams{
		Context:   "c1",
		Timestamp: 100,
		Request: NetworkRequestData{
			Request: "r1",
			URL:     "https://example.com/x",
			Method:  "GET",
		},
		ErrorText: "connection refused",
	}

	frame := eventFrame(t, NetworkFetchError, params)
	event := decodeEvent(t, frame)

	var got NetworkFetchErrorParams
	if err := json.Unmarshal(event.Params, &got); err != nil {
		t.Fatalf("decode params: %v", err)
	}

	if got.ErrorText != "connection refused" {
		t.Errorf("ErrorText = %q, want connection refused", got.ErrorText)
	}
}

func TestNetworkContinueRequestMarshal(t *testing.T) {
	headers := []NetworkHeader{
		{Name: "X-Trace", Value: NetworkBytesValue{Type: BytesValueString, Value: "b"}},
	}
	params := NetworkContinueRequestParams{
		Request: "r1",
		Body:    &NetworkBytesValue{Type: BytesValueBase64, Value: "aGk="},
		Headers: headers,
		Method:  "POST",
		URL:     "https://example.com/x",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"request":"r1","body":{"type":"base64","value":"aGk="},` +
		`"headers":[{"name":"X-Trace","value":{"type":"string","value":"b"}}],` +
		`"method":"POST","url":"https://example.com/x"}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestNetworkContinueResponseMarshal(t *testing.T) {
	status := 204
	params := NetworkContinueResponseParams{
		Request:      "r1",
		ReasonPhrase: "No Content",
		StatusCode:   &status,
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"request":"r1","reasonPhrase":"No Content","statusCode":204}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestNetworkContinueWithAuthMarshal(t *testing.T) {
	params := NetworkContinueWithAuthParams{
		Request: "r1",
		Action:  AuthActionProvideCredentials,
		Credentials: &NetworkAuthCredentials{
			Type:     AuthCredentialsTypePassword,
			Username: "u",
			Password: "p",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"request":"r1","action":"provideCredentials",` +
		`"credentials":{"type":"password","username":"u","password":"p"}}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestNetworkProvideResponseMarshal(t *testing.T) {
	params := NetworkProvideResponseParams{
		Request: "r1",
		Body:    &NetworkBytesValue{Type: BytesValueString, Value: "ok"},
		Headers: []NetworkHeader{
			{Name: "X-A", Value: NetworkBytesValue{Type: BytesValueString, Value: "1"}},
		},
		ReasonPhrase: "Tea",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"request":"r1","body":{"type":"string","value":"ok"},` +
		`"headers":[{"name":"X-A","value":{"type":"string","value":"1"}}],` +
		`"reasonPhrase":"Tea"}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestNetworkAddInterceptResultDecode(t *testing.T) {
	var result NetworkAddInterceptResult
	if err := json.Unmarshal([]byte(`{"intercept":"i1"}`), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	if result.Intercept != "i1" {
		t.Errorf("Intercept = %q, want i1", result.Intercept)
	}
}

func TestNetworkAddInterceptMarshal(t *testing.T) {
	params := NetworkAddInterceptParams{
		Phases: []string{
			InterceptPhaseBeforeRequestSent,
			InterceptPhaseResponseStarted,
		},
		Contexts: []string{"c1"},
		URLPatterns: []NetworkURLPattern{
			{
				Type:     "pattern",
				Protocol: "https",
				Hostname: "example.com",
			},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(data)
	want := `{"phases":["beforeRequestSent","responseStarted"],"contexts":["c1"],` +
		`"urlPatterns":[{"type":"pattern","protocol":"https","hostname":"example.com"}]}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestNetworkSetCacheBehaviorMarshal(t *testing.T) {
	params := NetworkSetCacheBehaviorParams{
		CacheBehavior: CacheBehaviorBypass,
		Contexts:      []string{"c1"},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(data); got != `{"cacheBehavior":"bypass","contexts":["c1"]}` {
		t.Errorf("got %s", got)
	}
}

func decodeEvent(t *testing.T, frame string) *Event {
	t.Helper()

	msg, err := Decode([]byte(frame))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	event, ok := msg.(*Event)
	if !ok {
		t.Fatalf("got %T, want *Event", msg)
	}

	return event
}

func eventFrame(t *testing.T, method string, params any) string {
	t.Helper()

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return fmt.Sprintf(`{"type":"event","method":%q,"params":%s}`, method, data)
}

func strPtr(s string) *string {
	return &s
}

func int64Ptr(n int64) *int64 {
	return &n
}
