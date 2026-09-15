// Command connect performs the WebDriver classic handshake against a
// running Firefox and reports the created session and its status.
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/yvv4git/go-bidi"
	"github.com/yvv4git/go-bidi/examples/internal/support"
)

func main() {
	endpoint := flag.String(
		"endpoint",
		"http://127.0.0.1:9222",
		"WebDriver BiDi endpoint of a running Firefox",
	)

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Printf("connecting to %s\n", *endpoint)

	result, err := bidi.Handshake(ctx, *endpoint)
	support.Must(err)

	defer func() { _ = result.Transport.Close() }()

	fmt.Printf("session:      %s\n", result.SessionID)
	fmt.Printf("websocket:    %s\n", result.WebSocketURL)
	fmt.Printf("capabilities: %v\n", result.Capabilities)

	client := bidi.NewClient(result.Transport)
	session := bidi.NewSession(client, result.SessionID)

	status, err := session.Status(ctx)
	support.Must(err)
	fmt.Printf("status: ready=%t message=%q\n", status.Ready, status.Message)
}
