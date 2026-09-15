# go-bidi — implementation plan

High-level plan for implementing a WebDriver BiDi client in pure Go to
drive Firefox (and other BiDi-compatible browsers).

Protocol overview: `man/it_protocol_net_web_browser_playwright_bidi.md`.
W3C specification: <https://w3c.github.io/webdriver-bidi/>.

The architecture mirrors `go-juggler` (the analogous Juggler-protocol
project): the same layers and one-way dependencies between them.

```text
go-bidi/
├── bidi.go          # root package: public API facade
├── transport/       # byte channels: WebSocket and pipe (fd 3/fd 4)
├── protocol/        # BiDi message types, JSON-RPC 2.0, commands and events
├── browser/         # high-level browser and page control
└── examples/        # runnable example programs
```

Dependencies: `browser/` -> `transport/` + `protocol/`, and `protocol/`
is independent of `transport/`.

BiDi specifics that matter for the implementation:

- JSON-RPC 2.0 framing: command `{id, method, params}`,
  response `{id, result|error}`, event `{method, params}` without `id`;
- connection via HTTP `POST /session` (WebDriver classic) with
  `webSocketUrl: true`; the server returns `sessionId` and `webSocketUrl`;
- everything then flows over a single WebSocket with the `webdriver.bidi`
  subprotocol;
- functionality is split into modules (session, browsingContext, script,
  network, input, log, storage, permissions, etc.);
- events are pushed to the same channel.

---

## Phase 0. Project bootstrap

- [ ] Initialize `go.mod` (module `github.com/yvv4git/go-bidi`)
- [ ] Add the `github.com/coder/websocket` dependency
- [ ] Add `.golangci.yml`, `.markdownlint.json`, `NOTICE`
- [ ] Set up the package layout (`transport/`, `protocol/`, `browser/`, `examples/`)
- [ ] Describe the architecture and layers in the `bidi.go` root package doc comment

## Phase 1. Transport (transport/)

- [ ] Define the `Transport` interface (`Send`, `Receive`, `Close`)
- [ ] Implement `WebSocketTransport` (Dial, Send, Receive, Close)
- [ ] Implement `PipeTransport` for fd 3/fd 4 (pipe-based launch)
- [ ] Ensure thread safety and correct channel shutdown
- [ ] Unit tests for the transport

## Phase 2. Protocol (protocol/)

- [ ] BiDi message types: `Command`, `Response`, `Error`, `Event`
- [ ] JSON-RPC 2.0 encoding/decoding (marshal/unmarshal)
- [ ] Distinguish responses from events when parsing a frame
- [ ] `Error` type with error code and message
- [ ] Method-name constants per module (session, browsingContext, script, ...)
- [ ] `RemoteValue` type and result wrappers for `script.evaluate`
- [ ] Parameter/result types for commands and events
- [ ] Unit tests for encoding/decoding

## Phase 3. session module

- [ ] Handshake: HTTP `POST /session` with capabilities (`webSocketUrl: true`)
- [ ] Parse `sessionId` and `webSocketUrl` from the response
- [ ] Establish the WebSocket with the `webdriver.bidi` subprotocol
- [ ] `session.status`
- [ ] `session.new` / `session.end`
- [ ] `session.subscribe` / `session.unsubscribe` (subscribe to events by module)
- [ ] Track active subscriptions

## Phase 4. Client core and message routing

- [ ] Reader loop: continuously read frames from the transport
- [ ] Match responses to requests by `id` (pending map -> chan)
- [ ] Generate monotonic `id`s (atomic)
- [ ] Demultiplex events by subscription
- [ ] Publish events to subscribers (fan-out / callback)
- [ ] Graceful shutdown: cancel contexts, close `pending` channels
- [ ] Handle timeouts and `context` while waiting for a response

## Phase 5. browsingContext module

- [ ] `browsingContext.create` (new tab/window/frame)
- [ ] `browsingContext.navigate` (with `wait` and `readinessState`)
- [ ] `browsingContext.getTree` (context tree, top-level frames)
- [ ] `browsingContext.close`
- [ ] `browsingContext.reload`, `traverseHistory` (back/forward)
- [ ] `browsingContext.captureScreenshot` (viewport / full page)
- [ ] `browsingContext.setViewport`
- [ ] Events: `contextCreated`, `contextDestroyed`, `navigationStarted`,
      `load`, `domContentLoaded`, `fragmentNavigated`
- [ ] High-level `Page` type (wrapper over a context id)

## Phase 6. script module

- [ ] `script.evaluate` (with `awaitPromise`, `resultOwnership`)
- [ ] `script.callFunction` (argument and this passing)
- [ ] Decode `RemoteValue` into Go values (primitives, arrays, objects)
- [ ] `script.getRealms`
- [ ] `script.addPreloadScript` / `script.removePreloadScript`
- [ ] `script.disown`
- [ ] Events: `script.message`, `script.realmCreated`, `script.realmDestroyed`
- [ ] High-level `Evaluate` and `CallFunction` helpers on `Page`

## Phase 7. input module

- [ ] `input.performActions` (keyboard, mouse, wheel)
- [ ] `input.releaseActions`
- [ ] `input.setFiles`
- [ ] `Click`, `Type`, `Press`, `Scroll` helpers on `Page`

## Phase 8. network module

- [ ] Subscribe to network-module events
- [ ] Events: `beforeRequestSent`, `responseStarted`, `responseCompleted`,
      `fetchError`, `authRequired`
- [ ] `network.continueRequest`, `continueResponse`
- [ ] `network.failRequest`
- [ ] `network.provideResponse`
- [ ] `network.addIntercept` / `network.removeIntercept`
- [ ] `network.setCacheBehavior`
- [ ] Collect and store page network requests (list + poll helper)

## Phase 9. log module and other modules

- [ ] `log` — `entryAdded` event, log-reading helpers
- [ ] `storage` — cookies, localStorage/sessionStorage
- [ ] `permissions` — `setPermission`
- [ ] `browser` — `browser.close`, `browser.createUserContext`
- [ ] `emulation` — `setGeolocationOverride`, `setTimezoneOverride`

## Phase 10. High-level API (browser/ and bidi.go)

- [ ] `Browser` type (wrapper over transport + router)
- [ ] `Connect(ctx, tr)` — connect over an existing transport
- [ ] `ConnectEndpoint(ctx, addr)` — handshake over an HTTP address
- [ ] `Page`/`Tab` type with control methods
- [ ] `NewPage`, `Page`, `Pages`, `Close` methods
- [ ] Re-export the public API from the root `bidi.go` package
- [ ] Client options (`WithTimeout`, `WithHTTPClient`, `WithSubprotocol`)

## Phase 11. Launching Firefox (launcher)

- [ ] Launch options (`WithExecPath`, `WithHeadless`, `WithTimeout`,
      `WithArgs`, `WithProfile`)
- [ ] Spawn the Firefox process with `--remote-debugging-port` (Remote Agent)
- [ ] Wait for endpoint readiness (poll `session.status` / HTTP)
- [ ] Pipe-based connection mode (fd 3/fd 4) where supported
- [ ] Graceful process shutdown (`Close`, kill on context cancel)
- [ ] `Launch(ctx, opts...)` at the root package level

## Phase 12. Examples (examples/)

- [ ] `connect/` — handshake and `session.status`
- [ ] `basic/` — open a page, navigate, screenshot
- [ ] `evaluate/` — run JS and decode `RemoteValue`
- [ ] `input/` — click, type text, key presses
- [ ] `events/` — subscribe to browsingContext/script events
- [ ] `network/` — intercept and modify requests
- [ ] `screenshot/` — viewport and full-page capture

## Phase 13. Quality, tests, documentation

- [ ] Integration tests against Firefox Nightly with BiDi
- [ ] Test fixture pages for navigation/network/DOM
- [ ] Run `go test ./...`, `go vet ./...`, `golangci-lint run`
- [ ] Update `README.md` (install, example, capability table, layout)
- [ ] CI (GitHub Actions): build, lint, test
- [ ] Note the analogy with `go-juggler` and link to `go-juggler-mcp`
