package browser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/atomic"

	"github.com/yvv4git/go-bidi/protocol"
	"github.com/yvv4git/go-bidi/transport"
)

// ErrClosed is returned when operating on a closed Client.
var ErrClosed = errors.New("bidi: connection closed")

// pending is the response slot delivered to a waiting Call.
type pending struct {
	result json.RawMessage
	err    error
}

// Client routes commands and events over a single transport. Commands are
// matched to responses by id; events are demultiplexed to the subscriptions
// that registered for them. Client is safe for concurrent use.
type Client struct {
	tr transport.Transport

	mu      sync.Mutex
	pending map[int64]chan pending
	nextID  atomic.Int64

	loopCtx    context.Context
	loopCancel context.CancelFunc
	closed     chan struct{}
	closeOnce  sync.Once
	closeErr   error
	done       chan struct{}

	eventBuffer int

	subMu sync.Mutex
	subs  map[string]*Subscription
}

var _ Caller = (*Client)(nil)

// NewClient starts a reader loop over tr and returns a routed client.
func NewClient(tr transport.Transport) *Client {
	loopCtx, loopCancel := context.WithCancel(context.Background())

	c := &Client{
		tr:          tr,
		pending:     make(map[int64]chan pending),
		loopCtx:     loopCtx,
		loopCancel:  loopCancel,
		closed:      make(chan struct{}),
		done:        make(chan struct{}),
		eventBuffer: _eventBuffer,
		subs:        make(map[string]*Subscription),
	}

	go c.readerLoop()

	return c
}

// Call sends a command and waits for its response, decoding the result into
// result when non-nil.
func (c *Client) Call(ctx context.Context, method string, params, result any) error {
	select {
	case <-c.closed:
		return ErrClosed
	default:
	}

	id := c.nextID.Add(1)

	command, err := protocol.NewCommand(id, method, params)
	if err != nil {
		return err
	}

	frame, err := command.Encode()
	if err != nil {
		return err
	}

	ch := make(chan pending, 1)

	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	if err := c.tr.Send(ctx, frame); err != nil {
		return fmt.Errorf("send command: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return ErrClosed
	case resp := <-ch:
		if resp.err != nil {
			return resp.err
		}

		if result != nil && len(resp.result) > 0 {
			if err := json.Unmarshal(resp.result, result); err != nil {
				return fmt.Errorf("decode result: %w", err)
			}
		}

		return nil
	}
}

// Close shuts the reader loop down and closes the transport. Pending calls
// fail with ErrClosed.
func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		c.loopCancel()
		c.closeErr = c.tr.Close()
		close(c.closed)
		<-c.done

		c.subMu.Lock()
		for id, sub := range c.subs {
			delete(c.subs, id)
			close(sub.ch)
		}
		c.subMu.Unlock()
	})

	return c.closeErr
}

func (c *Client) readerLoop() {
	defer close(c.done)

	for {
		data, err := c.tr.Receive(c.loopCtx)
		if err != nil {
			c.failPending()

			return
		}

		msg, err := protocol.Decode(data)
		if err != nil {
			continue
		}

		switch m := msg.(type) {
		case *protocol.Response:
			c.dispatchResponse(m)
		case *protocol.Event:
			c.dispatchEvent(m)
		}
	}
}

func (c *Client) dispatchResponse(resp *protocol.Response) {
	c.mu.Lock()
	ch, ok := c.pending[resp.ID]
	delete(c.pending, resp.ID)
	c.mu.Unlock()

	if !ok {
		return
	}

	p := pending{result: resp.Result}
	if resp.Err != nil {
		p.err = resp.Err
	}

	ch <- p
}

func (c *Client) failPending() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, ch := range c.pending {
		select {
		case ch <- pending{err: ErrClosed}:
		default:
		}

		delete(c.pending, id)
	}
}
