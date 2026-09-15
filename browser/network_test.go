package browser

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/yvv4git/go-bidi/protocol"
)

func TestSessionAddIntercept(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{
		protocol.NetworkAddIntercept: `{"intercept":"i1"}`,
	}}
	session := NewSession(caller, "s1")

	got, err := session.AddIntercept(context.Background(), protocol.NetworkAddInterceptParams{
		Phases: []string{protocol.InterceptPhaseBeforeRequestSent},
	})
	if err != nil {
		t.Fatalf("AddIntercept: %v", err)
	}

	if got != "i1" {
		t.Errorf("Intercept = %q, want i1", got)
	}
}

func TestSessionRemoveIntercept(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.RemoveIntercept(context.Background(), protocol.NetworkRemoveInterceptParams{
		Intercept: "i1",
	}); err != nil {
		t.Fatalf("RemoveIntercept: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.NetworkRemoveIntercept {
		t.Errorf("method = %q, want %q", m, protocol.NetworkRemoveIntercept)
	}
}

func TestSessionContinueRequest(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.ContinueRequest(context.Background(), protocol.NetworkContinueRequestParams{
		Request: "r1",
		Method:  "POST",
	}); err != nil {
		t.Fatalf("ContinueRequest: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.NetworkContinueRequest {
		t.Errorf("method = %q, want %q", call.method, protocol.NetworkContinueRequest)
	}

	params := call.params.(protocol.NetworkContinueRequestParams)
	if params.Request != "r1" || params.Method != "POST" {
		t.Errorf("params = %+v", params)
	}
}

func TestSessionContinueResponse(t *testing.T) {
	status := int(200)
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.ContinueResponse(context.Background(), protocol.NetworkContinueResponseParams{
		Request:    "r1",
		StatusCode: &status,
	}); err != nil {
		t.Fatalf("ContinueResponse: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.NetworkContinueResponse {
		t.Errorf("method = %q, want %q", m, protocol.NetworkContinueResponse)
	}
}

func TestSessionContinueWithAuth(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.ContinueWithAuth(context.Background(), protocol.NetworkContinueWithAuthParams{
		Request: "r1",
		Action:  protocol.AuthActionProvideCredentials,
		Credentials: &protocol.NetworkAuthCredentials{
			Type:     protocol.AuthCredentialsTypePassword,
			Username: "u",
			Password: "p",
		},
	}); err != nil {
		t.Fatalf("ContinueWithAuth: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.NetworkContinueWithAuth {
		t.Errorf("method = %q, want %q", call.method, protocol.NetworkContinueWithAuth)
	}

	params := call.params.(protocol.NetworkContinueWithAuthParams)
	if params.Action != protocol.AuthActionProvideCredentials ||
		params.Credentials == nil ||
		params.Credentials.Username != "u" {
		t.Errorf("params = %+v", params)
	}
}

func TestSessionFailRequest(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.FailRequest(context.Background(), protocol.NetworkFailRequestParams{
		Request: "r1",
	}); err != nil {
		t.Fatalf("FailRequest: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.NetworkFailRequest {
		t.Errorf("method = %q, want %q", m, protocol.NetworkFailRequest)
	}
}

func TestSessionProvideResponse(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.ProvideResponse(context.Background(), protocol.NetworkProvideResponseParams{
		Request: "r1",
		Body:    &protocol.NetworkBytesValue{Type: protocol.BytesValueString, Value: "ok"},
	}); err != nil {
		t.Fatalf("ProvideResponse: %v", err)
	}

	if m := caller.lastCall(t).method; m != protocol.NetworkProvideResponse {
		t.Errorf("method = %q, want %q", m, protocol.NetworkProvideResponse)
	}
}

func TestSessionSetCacheBehavior(t *testing.T) {
	caller := &fakeCaller{results: map[string]string{}}
	session := NewSession(caller, "s1")

	if err := session.SetCacheBehavior(context.Background(), protocol.NetworkSetCacheBehaviorParams{
		CacheBehavior: protocol.CacheBehaviorBypass,
		Contexts:      []string{"c1"},
	}); err != nil {
		t.Fatalf("SetCacheBehavior: %v", err)
	}

	call := caller.lastCall(t)
	if call.method != protocol.NetworkSetCacheBehavior {
		t.Errorf("method = %q, want %q", call.method, protocol.NetworkSetCacheBehavior)
	}

	params := call.params.(protocol.NetworkSetCacheBehaviorParams)
	if params.CacheBehavior != protocol.CacheBehaviorBypass {
		t.Errorf("params = %+v", params)
	}
}

func TestRequestsCollectsEvents(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	r := NewRequests(sub, "c1")

	r.record(protocol.NetworkBeforeRequestSent,
		eventParams(t, beforeRequestParams("r1", "c1")))
	r.record(protocol.NetworkResponseCompleted,
		eventParams(t, responseCompletedParams("r1", "c1", 200)))
	r.record(protocol.NetworkBeforeRequestSent,
		eventParams(t, beforeRequestParams("r2", "c2")))

	got := r.List()
	if len(got) != 1 {
		t.Fatalf("requests = %d, want 1", len(got))
	}

	req := got[0]
	if req.ID != "r1" || req.URL != "https://example.com" || req.Method != "GET" {
		t.Errorf("request = %+v", req)
	}

	if !req.Complete || req.Status != 200 {
		t.Errorf("request state = %+v, want complete 200", req)
	}

	if len(req.Headers) != 1 || req.Headers[0].Name != "Host" {
		t.Errorf("headers = %+v", req.Headers)
	}
}

func TestRequestsRecordsFetchError(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	r := NewRequests(sub, "c1")

	r.record(protocol.NetworkBeforeRequestSent,
		eventParams(t, beforeRequestParams("r1", "c1")))
	r.record(protocol.NetworkFetchError,
		eventParams(t, fetchErrorParams("r1", "c1", "connection refused")))

	reqs := r.List()
	if len(reqs) != 1 {
		t.Fatalf("requests = %d, want 1", len(reqs))
	}

	if !reqs[0].Failed || reqs[0].ErrorText != "connection refused" {
		t.Errorf("request = %+v, want failed", reqs[0])
	}
}

func TestRequestsPoll(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	r := NewRequests(sub, "c1")

	r.record(protocol.NetworkBeforeRequestSent,
		eventParams(t, beforeRequestParams("r1", "c1")))

	done := make(chan struct{})
	go func() {
		defer close(done)
		got, err := r.Poll(context.Background(), func(reqs []*Request) bool {
			return len(reqs) == 1 && reqs[0].Complete
		})
		if err != nil {
			t.Errorf("Poll: %v", err)
		}

		if len(got) != 1 || !got[0].Complete {
			t.Errorf("Poll result = %+v", got)
		}
	}()

	r.record(protocol.NetworkResponseCompleted,
		eventParams(t, responseCompletedParams("r1", "c1", 200)))

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Poll did not return")
	}
}

func TestRequestsRun(t *testing.T) {
	sub := &Subscription{ch: make(chan protocol.Event, _eventBuffer)}
	r := NewRequests(sub, "c1")
	go r.run()

	before := beforeRequestParams("r1", "c1")
	sub.ch <- protocol.Event{
		Method: protocol.NetworkBeforeRequestSent,
		Params: eventParams(t, before),
	}

	complete := responseCompletedParams("r1", "c1", 200)
	sub.ch <- protocol.Event{
		Method: protocol.NetworkResponseCompleted,
		Params: eventParams(t, complete),
	}

	waitPoll(t, r, func(reqs []*Request) bool {
		return len(reqs) == 1 && reqs[0].Complete
	})
}

func TestRequestsClose(t *testing.T) {
	tr := newFakeTransport()
	client := NewClient(tr)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	go respondSuccess(t, tr, `{"subscription":"sub1"}`)
	sub, err := client.Subscribe(context.Background(), []string{"network"})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	r := NewRequests(sub, "c1")
	go r.run()

	sub.ch <- protocol.Event{
		Method: protocol.NetworkBeforeRequestSent,
		Params: eventParams(t, beforeRequestParams("r1", "c1")),
	}
	waitPoll(t, r, func(reqs []*Request) bool { return len(reqs) == 1 })

	go respondSuccess(t, tr, `{}`)
	if err := r.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}

	r.record(protocol.NetworkBeforeRequestSent,
		eventParams(t, beforeRequestParams("r2", "c1")))

	if got := r.Count(); got != 1 {
		t.Errorf("Count = %d, want 1", got)
	}
}

func waitPoll(t *testing.T, r *Requests, cond func([]*Request) bool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := r.Poll(ctx, cond); err != nil {
		t.Fatalf("Poll: %v", err)
	}
}

func eventParams(t *testing.T, params any) json.RawMessage {
	t.Helper()

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return data
}

func beforeRequestParams(id, contextID string) protocol.NetworkBeforeRequestSentParams {
	return protocol.NetworkBeforeRequestSentParams{
		Context:   contextID,
		Timestamp: 100,
		Request: protocol.NetworkRequestData{
			Request: id,
			URL:     "https://example.com",
			Method:  "GET",
			Headers: []protocol.NetworkHeader{
				{
					Name:  "Host",
					Value: protocol.NetworkBytesValue{Type: protocol.BytesValueString, Value: "example.com"},
				},
			},
		},
	}
}

func responseCompletedParams(
	id, contextID string,
	status int,
) protocol.NetworkResponseCompletedParams {
	return protocol.NetworkResponseCompletedParams{
		Context:   contextID,
		Timestamp: 200,
		Request:   protocol.NetworkRequestData{Request: id},
		Response:  protocol.NetworkResponseData{Status: status, MimeType: "text/html"},
	}
}

func fetchErrorParams(id, contextID, text string) protocol.NetworkFetchErrorParams {
	return protocol.NetworkFetchErrorParams{
		Context:   contextID,
		Timestamp: 300,
		Request:   protocol.NetworkRequestData{Request: id},
		ErrorText: text,
	}
}
