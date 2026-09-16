package browser

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/yvv4git/go-bidi/protocol"
	"github.com/yvv4git/go-bidi/transport"
)

// Browser is a high-level handle over an established BiDi session. It owns
// the client transport and tracks the pages created through it. Browser is
// safe for concurrent use.
type Browser struct {
	client  *Client
	session *Session

	mu    sync.Mutex
	pages map[string]*Page
}

// Connect establishes a BiDi session over tr and returns a Browser that
// owns the transport. Truncated handshakes that already carry a session
// should use ConnectEndpoint or NewClient with NewSession instead.
func Connect(ctx context.Context, tr transport.Transport, opts ...Option) (*Browser, error) {
	client := NewClient(tr, opts...)

	session, err := CreateSession(ctx, client, protocol.SessionNewParams{})
	if err != nil {
		_ = client.Close()

		return nil, err
	}

	return newBrowser(client, session), nil
}

// ConnectEndpoint performs the WebDriver classic handshake at addr and
// returns a Browser for the created session.  If the classic POST /session
// endpoint is unavailable (e.g. Firefox 158+), ConnectEndpoint
// transparently falls back to the direct BiDi WebSocket flow.
func ConnectEndpoint(ctx context.Context, addr string, opts ...Option) (*Browser, error) {
	result, err := Handshake(ctx, addr, opts...)
	if err == nil {
		client := NewClient(result.Transport, opts...)
		session := NewSession(client, result.SessionID)

		return newBrowser(client, session), nil
	}

	// Classic handshake failed — try the direct BiDi WebSocket session
	// flow used by newer browser implementations.
	browser, fallbackErr := ConnectBiDi(ctx, addr, opts...)
	if fallbackErr != nil {
		return nil, fmt.Errorf("handshake: %w; direct BiDi: %w", err, fallbackErr)
	}

	return browser, nil
}

func newBrowser(client *Client, session *Session) *Browser {
	return &Browser{
		client:  client,
		session: session,
		pages:   make(map[string]*Page),
	}
}

// Session returns the underlying session handle.
func (b *Browser) Session() *Session {
	return b.session
}

// Client returns the underlying client that owns the transport.
func (b *Browser) Client() *Client {
	return b.client
}

// NewPage creates a new tab and returns its Page handle.
func (b *Browser) NewPage(ctx context.Context) (*Page, error) {
	result, err := b.session.CreateBrowsingContext(
		ctx,
		protocol.BrowsingContextCreateParams{Type: protocol.ContextTypeTab},
	)
	if err != nil {
		return nil, err
	}

	return b.track(result.Context), nil
}

// Page returns the tracked page for the given context id. Unknown ids are
// added as untracked handles.
func (b *Browser) Page(id string) *Page {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p, ok := b.pages[id]; ok {
		return p
	}

	p := NewPage(b.client, id)
	b.pages[id] = p

	return p
}

// Pages returns the tracked pages ordered by context id.
func (b *Browser) Pages() []*Page {
	b.mu.Lock()
	defer b.mu.Unlock()

	ids := make([]string, 0, len(b.pages))
	for id := range b.pages {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	pages := make([]*Page, 0, len(ids))
	for _, id := range ids {
		pages = append(pages, b.pages[id])
	}

	return pages
}

// Close ends the session and closes the transport.
func (b *Browser) Close(ctx context.Context) error {
	err := b.session.End(ctx)
	_ = b.client.Close()

	return err
}

func (b *Browser) track(id string) *Page {
	b.mu.Lock()
	defer b.mu.Unlock()

	p := NewPage(b.client, id)
	b.pages[id] = p

	return p
}
