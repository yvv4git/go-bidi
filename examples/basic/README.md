# basic

Launches a headless Firefox, navigates to a page and saves a screenshot.

## Run

```bash
go run examples/basic/main.go
```

## Examples

```bash
# Default: launch headless Firefox, save screenshot
go run examples/basic/main.go

# Connect to running Firefox
go run examples/basic/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox (see the window)
go run examples/basic/main.go -headless=false

# Custom Firefox binary and timeout
go run examples/basic/main.go -exec /usr/local/bin/firefox -timeout 5m

# Custom output file
go run examples/basic/main.go -out screenshot.png

# Combined flags
go run examples/basic/main.go -endpoint http://127.0.0.1:9222 -out wiki.png -timeout 10s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |
| `-out` | `basic.png` | Output file path |

## What it demonstrates

- Launching a headless Firefox instance
- Performing the WebDriver BiDi handshake
- Creating a new browsing context (tab)
- Setting viewport dimensions
- Navigating to a URL
- Taking a viewport screenshot
