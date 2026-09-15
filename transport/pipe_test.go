package transport

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestPipeRoundTrip(t *testing.T) {
	peerR, clientW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	clientR, peerW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	client := NewPipe(clientR, clientW)
	peer := NewPipe(peerR, peerW)

	t.Cleanup(func() { _ = client.Close() })
	t.Cleanup(func() { _ = peer.Close() })

	ctx := context.Background()

	frame := `{"id":1,"method":"session.new","params":{"capabilities":{}}}`
	if err := client.Send(ctx, []byte(frame)); err != nil {
		t.Fatalf("send: %v", err)
	}

	got, err := peer.Receive(ctx)
	if err != nil {
		t.Fatalf("receive: %v", err)
	}

	if string(got) != frame {
		t.Fatalf("got %q, want %q", got, frame)
	}

	reply := `{"id":1,"result":{"sessionId":"abc"}}`
	if err := peer.Send(ctx, []byte(reply)); err != nil {
		t.Fatalf("send reply: %v", err)
	}

	got, err = client.Receive(ctx)
	if err != nil {
		t.Fatalf("receive reply: %v", err)
	}

	if string(got) != reply {
		t.Fatalf("got %q, want %q", got, reply)
	}
}

func TestPipeReceiveContextCancel(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	client := NewPipe(r, w)

	t.Cleanup(func() { _ = client.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := client.Receive(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
}

func TestPipeClose(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	client := NewPipe(r, w)

	if err := client.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}

	ctx := context.Background()

	if err := client.Send(ctx, []byte(`{"id":1}`)); !errors.Is(err, ErrClosed) {
		t.Fatalf("send after close: got %v, want ErrClosed", err)
	}

	if _, err := client.Receive(ctx); !errors.Is(err, ErrClosed) {
		t.Fatalf("receive after close: got %v, want ErrClosed", err)
	}
}

func TestPipeReceiveTrimsLineEnding(t *testing.T) {
	tests := []struct {
		give string
		want string
	}{
		{give: "hello\n", want: "hello"},
		{give: "hello\r\n", want: "hello"},
		{give: "{}\n", want: "{}"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			client := NewPipe(r, w)

			t.Cleanup(func() { _ = client.Close() })

			if _, err := w.WriteString(tt.give); err != nil {
				t.Fatal(err)
			}

			got, err := client.Receive(context.Background())
			if err != nil {
				t.Fatalf("receive: %v", err)
			}

			if string(got) != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
