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

- [x] Initialize `go.mod` (module `github.com/yvv4git/go-bidi`)
- [x] Add the `github.com/coder/websocket` dependency
- [x] Add `.golangci.yml`, `.markdownlint.json`, `NOTICE`
- [x] Set up the package layout (`transport/`, `protocol/`, `browser/`, `examples/`)
- [x] Describe the architecture and layers in the `bidi.go` root package doc comment

## Phase 1. Transport (transport/)

- [x] Define the `Transport` interface (`Send`, `Receive`, `Close`)
- [x] Implement `WebSocketTransport` (Dial, Send, Receive, Close)
- [x] Implement `PipeTransport` for fd 3/fd 4 (pipe-based launch)
- [x] Ensure thread safety and correct channel shutdown
- [x] Unit tests for the transport

## Phase 2. Protocol (protocol/)

- [x] BiDi message types: `Command`, `Response`, `Error`, `Event`
- [x] JSON-RPC 2.0 encoding/decoding (marshal/unmarshal)
- [x] Distinguish responses from events when parsing a frame
- [x] `Error` type with error code and message
- [x] Method-name constants per module (session, browsingContext, script, ...)
- [x] `RemoteValue` type and result wrappers for `script.evaluate`
- [x] Parameter/result types for commands and events
- [x] Unit tests for encoding/decoding

## Phase 3. Session module

- [x] Handshake: HTTP `POST /session` with capabilities (`webSocketUrl: true`)
- [x] Parse `sessionId` and `webSocketUrl` from the response
- [x] Establish the WebSocket with the `webdriver.bidi` subprotocol
- [x] `session.status`
- [x] `session.new` / `session.end`
- [x] `session.subscribe` / `session.unsubscribe` (subscribe to events by module)
- [x] Track active subscriptions

## Phase 4. Client core and message routing

- [x] Reader loop: continuously read frames from the transport
- [x] Match responses to requests by `id` (pending map -> chan)
- [x] Generate monotonic `id`s (atomic)
- [x] Demultiplex events by subscription
- [x] Publish events to subscribers (fan-out / callback)
- [x] Graceful shutdown: cancel contexts, close `pending` channels
- [x] Handle timeouts and `context` while waiting for a response

## Phase 5. BrowsingContext module

- [x] `browsingContext.create` (new tab/window/frame)
- [x] `browsingContext.navigate` (with `wait` and `readinessState`)
- [x] `browsingContext.getTree` (context tree, top-level frames)
- [x] `browsingContext.close`
- [x] `browsingContext.reload`, `traverseHistory` (back/forward)
- [x] `browsingContext.captureScreenshot` (viewport / full page)
- [x] `browsingContext.setViewport`
- [x] Events: `contextCreated`, `contextDestroyed`, `navigationStarted`, `load`, `domContentLoaded`, `fragmentNavigated`
- [x] High-level `Page` type (wrapper over a context id)

## Phase 6. Script module

- [x] `script.evaluate` (with `awaitPromise`, `resultOwnership`)
- [x] `script.callFunction` (argument and this passing)
- [x] Decode `RemoteValue` into Go values (primitives, arrays, objects)
- [x] `script.getRealms`
- [x] `script.addPreloadScript` / `script.removePreloadScript`
- [x] `script.disown`
- [x] Events: `script.message`, `script.realmCreated`, `script.realmDestroyed`
- [x] High-level `Evaluate` and `CallFunction` helpers on `Page`

## Phase 7. Input module

- [x] `input.performActions` (keyboard, mouse, wheel)
- [x] `input.releaseActions`
- [x] `input.setFiles`
- [x] `Click`, `Type`, `Press`, `Scroll` helpers on `Page`

## Phase 8. Network module

- [x] Subscribe to network-module events
- [x] Events: `beforeRequestSent`, `responseStarted`, `responseCompleted`,
      `fetchError`, `authRequired`
- [x] `network.continueRequest`, `continueResponse`
- [x] `network.failRequest`
- [x] `network.provideResponse`
- [x] `network.addIntercept` / `network.removeIntercept`
- [x] `network.setCacheBehavior`
- [x] Collect and store page network requests (list + poll helper)

## Phase 9. Log module and other modules

- [x] `log` — `entryAdded` event, log-reading helpers
- [x] `storage` — cookies, localStorage/sessionStorage
- [x] `permissions` — `setPermission`
- [x] `browser` — `browser.close`, `browser.createUserContext`
- [x] `emulation` — `setGeolocationOverride`, `setTimezoneOverride`

## Phase 10. High-level API (browser/ and bidi.go)

- [x] `Browser` type (wrapper over transport + router)
- [x] `Connect(ctx, tr)` — connect over an existing transport
- [x] `ConnectEndpoint(ctx, addr)` — handshake over an HTTP address
- [x] `Page`/`Tab` type with control methods
- [x] `NewPage`, `Page`, `Pages`, `Close` methods
- [x] Re-export the public API from the root `bidi.go` package
- [x] Client options (`WithTimeout`, `WithHTTPClient`, `WithSubprotocol`)

## Phase 11. Launching Firefox (launcher)

- [x] Launch options (`WithExecPath`, `WithHeadless`, `WithTimeout`,
      `WithArgs`, `WithProfile`)
- [x] Spawn the Firefox process with `--remote-debugging-port` (Remote Agent)
- [x] Wait for endpoint readiness (poll `session.status` / HTTP)
- [x] Pipe-based connection mode (fd 3/fd 4) where supported
- [x] Graceful process shutdown (`Close`, kill on context cancel)
- [x] `Launch(ctx, opts...)` at the root package level

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
