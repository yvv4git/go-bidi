# basic

Launches a headless Firefox, navigates to a page and saves a screenshot.

## Run

```bash
go run examples/basic/main.go
```

Firefox will be launched automatically in headless mode. No manual setup required.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-timeout` | `2m` | Overall command timeout |
| `-out` | `basic.png` | Output file path |

## What it demonstrates

- Launching a headless Firefox instance
- Performing the WebDriver BiDi handshake
- Creating a new browsing context (tab)
- Setting viewport dimensions
- Navigating to a URL
- Taking a viewport screenshot

## Connecting to an existing Firefox

If Firefox is already running with remote debugging enabled (`--remote-debugging-port`), connect directly:

```bash
go run examples/basic/main.go -endpoint http://127.0.0.1:9222
```
