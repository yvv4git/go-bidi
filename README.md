# go-bidi

![go-bidi logo](./assets/logo.jpeg)

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white&style=flat-square)][Go] [![License](https://img.shields.io/badge/license-custom-blue.svg?style=flat-square)][LICENSE] [![CI](https://github.com/yvv4git/go-bidi/actions/workflows/ci.yml/badge.svg)][CI]

Go client for the [WebDriver BiDi] protocol. It drives BiDi-enabled
browsers — Firefox Nightly and other WebDriver BiDi implementations — over a
single full-duplex WebSocket channel using JSON-RPC 2.0 messages, with no
Selenium, Playwright or browser drivers required.

The connection follows the WebDriver BiDi handshake: the client POSTs to the
WebDriver classic `HTTP /session` endpoint with `webSocketUrl: true`, receives
`sessionId` and `webSocketUrl`, and then speaks BiDi over the returned
WebSocket with the `webdriver.bidi` subprotocol. Everything — commands,
responses and pushed events — flows over that one channel.

## Features

- Pure Go, zero browser-driver binaries
- WebSocket (`webSocketUrl`) and pipe (fd 3/fd 4) transports
- Full session lifecycle: `session.status`, `session.new`, `session.end`
- Tab/window and frame control: create, navigate, history, screenshots,
  viewport
- JavaScript via `script.evaluate` / `script.callFunction` with decoded
  `RemoteValue` results
- Keyboard, mouse and wheel input, file upload
- Network interception and request/response modification
- Cookies and localStorage/sessionStorage access
- Event subscriptions for browsingContext, script, network, log and input
  modules; `Requests` / `Logs` collectors with polling
- High-level `Browser` / `Page` API plus the raw `Client` / `Session` handles
- Firefox launcher with headless mode and readiness waiting

## Requirements

- Go 1.26+
- A browser that speaks WebDriver BiDi: **Firefox Nightly** with the Remote
  Agent enabled, or any compliant WebDriver BiDi endpoint

## Install

```sh
go get github.com/yvv4git/go-bidi@latest
```

## Quick start

Launch a headless Firefox, perform the handshake, open a page, and save a
screenshot — all from the root package:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/yvv4git/go-bidi"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Launch headless Firefox and wait until its BiDi endpoint is ready.
    firefox, err := bidi.Launch(ctx,
        bidi.WithLaunchHeadless(true),
        bidi.WithLaunchTimeout(time.Minute),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer firefox.Close()

    // WebDriver classic handshake: POST /session, then dial the WebSocket.
    browser, err := bidi.ConnectEndpoint(ctx, "http://"+firefox.Endpoint().Host)
    if err != nil {
        log.Fatal(err)
    }
    defer func() { _ = browser.Close(ctx) }()

    page, err := browser.NewPage(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer func() { _ = page.Close(ctx) }()

    if err := page.SetViewport(ctx, 1280, 720); err != nil {
        log.Fatal(err)
    }

    nav, err := page.Navigate(ctx, "https://example.com/")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("loaded %s (navigation %s)\n", nav.URL, nav.Navigation)

    // Evaluate an expression in the page and decode the result.
    value, err := page.Evaluate(ctx, "document.title")
    if err != nil {
        log.Fatal(err)
    }
    title, err := value.Result.Interface()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("title:", title)

    png, err := page.Screenshot(ctx, false)
    if err != nil {
        log.Fatal(err)
    }
    if err := os.WriteFile("page.png", png, 0o600); err != nil {
        log.Fatal(err)
    }
}
```

Connect to an already-running Firefox instead with `ConnectEndpoint` pointing
at its `http` endpoint, or take full control of the raw channel:

```go
result, err := bidi.Handshake(ctx, "http://127.0.0.1:9222")
if err != nil {
	log.Fatal(err)
}
defer result.Transport.Close()

client := bidi.NewClient(result.Transport, bidi.WithTimeout(30*time.Second))
session := bidi.NewSession(client, result.SessionID)
```

Firefox 158+ no longer supports the classic `POST /session` handshake.
`ConnectEndpoint` falls back to the direct BiDi WebSocket flow
automatically; you can also use `ConnectBiDi` explicitly:

```go
browser, err := bidi.ConnectBiDi(ctx, "http://127.0.0.1:9222")
if err != nil {
	log.Fatal(err)
}
defer func() { _ = browser.Close(ctx) }()
```

## Session management

Firefox 158+ supports only **one** BiDi session at a time. If a client
disconnects without calling `session.end`, the zombie session blocks all
subsequent `session.new` calls with "Maximum number of active sessions".

`ConnectEndpoint` and `ConnectBiDi` handle this automatically:

1. On success, the library stores the session ID in the OS temp directory
   (`os.TempDir()`).
2. On the next connect, if `session.new` fails with "Maximum number of active
   sessions", the library reads the stored ID, sends `session.end` for the
   zombie session, and retries `session.new`.

If you call `browser.Close(ctx)` or `Session.End(ctx)` on every code path
(including defers), zombie sessions are avoided entirely. The recovery
mechanism is a safety net for abnormal exits.

## Capabilities

Implemented protocol modules and the commands they expose:

| Module             | Commands |
| ------------------ | -------- |
| `session`          | `status` · `new` · `end` · `subscribe` · `unsubscribe` |
| `browser`          | `close` · `createUserContext` |
| `browsingContext`  | `create` · `navigate` · `getTree` · `close` · `reload` · `traverseHistory` · `captureScreenshot` · `setViewport` |
| `script`           | `evaluate` · `callFunction` · `getRealms` · `addPreloadScript` · `removePreloadScript` · `disown` |
| `input`            | `performActions` · `releaseActions` · `setFiles` |
| `network`          | `addIntercept` · `removeIntercept` · `continueRequest` · `continueResponse` · `continueWithAuth` · `failRequest` · `provideResponse` · `setCacheBehavior` |
| `storage`          | `getCookies` · `setCookie` · `deleteCookies` |
| `permissions`      | `setPermission` |
| `emulation`        | `setGeolocationOverride` · `setTimezoneOverride` |

Subscribable events:

| Module             | Events |
| ------------------ | ------ |
| `browsingContext`  | `contextCreated` · `contextDestroyed` · `navigationStarted` · `navigationCommitted` · `navigationAborted` · `navigationFailed` · `domContentLoaded` · `load` · `fragmentNavigated` · `historyUpdated` · `userPromptOpened` · `userPromptClosed` · `downloadBegin` · `downloadEnd` |
| `script`           | `message` · `realmCreated` · `realmDestroyed` |
| `network`          | `beforeRequestSent` · `responseStarted` · `responseCompleted` · `fetchError` · `authRequired` |
| `log`              | `entryAdded` |
| `input`            | `fileDialogOpened` |

High-level helpers on top of the raw commands:

- `Browser` — `NewPage` / `Page` / `Pages` / `Session` / `Close`
- `Page` — `Navigate` / `Reload` / `Back` / `Forward` / `Screenshot` /
  `SetViewport` / `Evaluate` / `CallFunction` / `Click` / `Type` / `Press` /
  `Scroll` / storage accessors
- `Client` — `Call` / `Subscribe` / `Close`
- `Session` — `Status` / `End` / `CreateBrowsingContext` /
  `Subscribe` / `Unsubscribe` / network & storage commands
- `Requests` / `Logs` — collectors with `Count` / `List` / `Poll` / `Close`

## Project layout

```text
go-bidi/
├── bidi.go          # root package: public API facade
├── transport/       # byte channels: WebSocket and pipe (fd 3/fd 4)
├── protocol/        # BiDi message types, JSON-RPC 2.0, commands and events
├── browser/         # high-level browser and page control
├── launcher/        # Firefox process launch and readiness
└── examples/        # runnable example programs
```

Dependencies flow one way, from high level to low level: `browser/` uses
`transport/` and `protocol/`, while `protocol/` stays independent so that
`transport/` moves raw frames and `protocol/` gives them meaning. The root
package re-exports the public API so most callers need a single import.

## Examples

The `examples/` directory contains runnable programs: `connect`, `basic`,
`evaluate`, `input`, `events`, `network` and `screenshot`. Each takes
`-endpoint` (connect to a running Firefox) or launches its own headless
instance:

```sh
go run ./examples/basic -endpoint http://127.0.0.1:9222
go run ./examples/network -out network.png
```

## Development

All automation is driven through the `Makefile` — CI invokes the same
targets:

```sh
make check       # fmt-check + vet + lint + test (the full quality gate)
make build       # compile all packages
make test        # run unit tests
make test-race   # run unit tests with the race detector
make cover       # run tests and print the coverage report
make vet         # go vet ./...
make fmt         # format the code with gofumpt
make lint        # golangci-lint
make examples    # build the runnable examples
```

## Relationship to go-juggler

`go-bidi` shares its layer architecture with
[`go-juggler`](https://github.com/yvv4git/go-juggler), the analogous Go
client for the Juggler protocol: the same `transport/`, `protocol/`,
`browser/` and `examples/` split with one-way dependencies between them.
Where `go-juggler` targets the Juggler protocol, `go-bidi` targets the
WebDriver BiDi protocol.

[`go-juggler-mcp`](https://github.com/yvv4git/go-juggler-mcp) is an
MCP-compatible server on top of the Juggler interface that exposes browser
automation as MCP tools — the same idea could be built on top of this
library for a WebDriver BiDi-based MCP server.

## License

Custom license — see [LICENSE] for full text.

**Summary:**
- Free to use, including in commercial products
- No selling or redistribution as a standalone product
- No claiming authorship or ownership
- Contributions via pull requests are welcome
- Commercial redistribution available by agreement with YVV

[WebDriver BiDi]: https://w3c.github.io/webdriver-bidi/
[Go]: https://go.dev/
[LICENSE]: LICENSE
[CI]: https://github.com/yvv4git/go-bidi/actions/workflows/ci.yml
