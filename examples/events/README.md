# events

Subscribes to browser events and prints them while a page loads.

## Run

```bash
go run examples/events/main.go
```

## Examples

```bash
# Default: launch headless Firefox, subscribe to events
go run examples/events/main.go

# Connect to running Firefox
go run examples/events/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox
go run examples/events/main.go -headless=false

# With custom timeout
go run examples/events/main.go -timeout 30s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |

## What it demonstrates

- Subscribing to browsingContext and script events
- Receiving real-time event notifications
- Inspecting browsingContext created, DOM content loaded, page load, and realm created events
