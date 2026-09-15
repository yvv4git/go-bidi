package browser

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/yvv4git/go-bidi/protocol"
)

// AddIntercept registers a network intercept and returns its id.
func (s *Session) AddIntercept(
	ctx context.Context,
	params protocol.NetworkAddInterceptParams,
) (string, error) {
	var result protocol.NetworkAddInterceptResult
	if err := s.caller.Call(ctx, protocol.NetworkAddIntercept, params, &result); err != nil {
		return "", err
	}

	return result.Intercept, nil
}

// RemoveIntercept unregisters a previously added network intercept.
func (s *Session) RemoveIntercept(
	ctx context.Context,
	params protocol.NetworkRemoveInterceptParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkRemoveIntercept, params, nil)
}

// ContinueRequest continues an intercepted request, optionally modifying its
// body, cookies, headers, method or URL.
func (s *Session) ContinueRequest(
	ctx context.Context,
	params protocol.NetworkContinueRequestParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkContinueRequest, params, nil)
}

// ContinueResponse continues an intercepted response, optionally
// modifying its cookies, credentials, headers, reason phrase or status.
func (s *Session) ContinueResponse(
	ctx context.Context,
	params protocol.NetworkContinueResponseParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkContinueResponse, params, nil)
}

// ContinueWithAuth resolves an authentication challenge with the given
// action and optional credentials.
func (s *Session) ContinueWithAuth(
	ctx context.Context,
	params protocol.NetworkContinueWithAuthParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkContinueWithAuth, params, nil)
}

// FailRequest fails an intercepted request with a network error.
func (s *Session) FailRequest(
	ctx context.Context,
	params protocol.NetworkFailRequestParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkFailRequest, params, nil)
}

// ProvideResponse supplies a synthetic response for an intercepted request.
func (s *Session) ProvideResponse(
	ctx context.Context,
	params protocol.NetworkProvideResponseParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkProvideResponse, params, nil)
}

// SetCacheBehavior changes how the given contexts use the network cache.
func (s *Session) SetCacheBehavior(
	ctx context.Context,
	params protocol.NetworkSetCacheBehaviorParams,
) error {
	return s.caller.Call(ctx, protocol.NetworkSetCacheBehavior, params, nil)
}

// Request summarizes a single request tracked by Requests. Status and
// MimeType are only meaningful when HasStatus is true.
type Request struct {
	ID        string
	URL       string
	Method    string
	Headers   []protocol.NetworkHeader
	Timestamp int64
	Status    int
	MimeType  string
	ErrorText string
	HasStatus bool
	Complete  bool
	Failed    bool
}

// Requests records the network events delivered on a subscription for a
// single browsing context. Create one with NewRequests.
type Requests struct {
	sub       *Subscription
	contextID string
	mu        sync.Mutex
	reqs      []*Request
	byID      map[string]*Request
	wake      chan struct{}
	done      chan struct{}
	closed    bool
}

// NewRequests starts recording requests for contextID from the events on
// sub. The subscription must cover the network module events.
func NewRequests(sub *Subscription, contextID string) *Requests {
	return &Requests{
		sub:       sub,
		contextID: contextID,
		byID:      make(map[string]*Request),
		wake:      make(chan struct{}, 1),
		done:      make(chan struct{}),
	}
}

// List returns a snapshot of the requests recorded so far.
func (r *Requests) List() []*Request {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*Request, 0, len(r.reqs))
	for _, req := range r.reqs {
		cloned := *req
		out = append(out, &cloned)
	}

	return out
}

// Count returns the number of recorded requests.
func (r *Requests) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.reqs)
}

// Poll waits until cond reports true for the snapshot of recorded requests,
// or until ctx is done. It returns the snapshot that satisfied cond.
func (r *Requests) Poll(ctx context.Context, cond func([]*Request) bool) ([]*Request, error) {
	for {
		if out := r.List(); cond(out) {
			return out, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-r.wake:
		case <-r.done:
			return nil, context.Canceled
		}
	}
}

// Close stops recording and unsubscribes the underlying subscription.
func (r *Requests) Close(ctx context.Context) error {
	err := r.sub.Close(ctx)

	r.mu.Lock()
	if !r.closed {
		close(r.done)
		r.closed = true
	}
	r.mu.Unlock()

	return err
}

// run drains the subscription and records events until it is closed.
func (r *Requests) run() {
	for event := range r.sub.Receive() {
		r.record(event.Method, event.Params)
	}
}

func (r *Requests) record(method string, params json.RawMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-r.done:
		return
	default:
	}

	var frame struct {
		Context string `json:"context"`
	}
	if err := json.Unmarshal(params, &frame); err != nil {
		return
	}

	if frame.Context != r.contextID {
		return
	}

	switch method {
	case protocol.NetworkBeforeRequestSent:
		var p protocol.NetworkBeforeRequestSentParams
		if err := json.Unmarshal(params, &p); err != nil {
			return
		}

		req := &Request{
			ID:        p.Request.Request,
			URL:       p.Request.URL,
			Method:    p.Request.Method,
			Headers:   p.Request.Headers,
			Timestamp: p.Timestamp,
		}
		r.reqs = append(r.reqs, req)
		r.byID[req.ID] = req

	case protocol.NetworkResponseStarted, protocol.NetworkResponseCompleted:
		var p protocol.NetworkResponseStartedParams
		if err := json.Unmarshal(params, &p); err != nil {
			return
		}

		if req := r.byID[p.Request.Request]; req != nil {
			req.Status = p.Response.Status
			req.MimeType = p.Response.MimeType
			req.HasStatus = true
			req.Complete = method == protocol.NetworkResponseCompleted
		}

	case protocol.NetworkFetchError:
		var p protocol.NetworkFetchErrorParams
		if err := json.Unmarshal(params, &p); err != nil {
			return
		}

		if req := r.byID[p.Request.Request]; req != nil {
			req.ErrorText = p.ErrorText
			req.Failed = true
		}
	}

	select {
	case r.wake <- struct{}{}:
	default:
	}
}
