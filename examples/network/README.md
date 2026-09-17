# network

Subscribes to network events and prints them while navigating to a website.

## Run

```bash
go run examples/network/main.go
```

## Examples

```bash
# Default: launch headless Firefox, print network events for wikipedia.org
go run examples/network/main.go

# Connect to running Firefox
go run examples/network/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox
go run examples/network/main.go -headless=false

# With custom timeout
go run examples/network/main.go -timeout 30s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |

## What it demonstrates

- Subscribing to network events (beforeRequestSent, responseStarted, responseCompleted)
- Observing real HTTP requests to wikipedia.org
- Printing request methods, URLs and response status codes
