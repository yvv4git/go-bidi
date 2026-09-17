# connect

Connects to an already-running Firefox and reports the session status.

## Run

```bash
go run examples/connect/main.go
```

## Examples

```bash
# Default: connect to localhost:9222
go run examples/connect/main.go

# Connect to custom endpoint
go run examples/connect/main.go -endpoint http://127.0.0.1:9333

# With timeout
go run examples/connect/main.go -endpoint http://127.0.0.1:9222 -timeout 10s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | `http://127.0.0.1:9222` | WebDriver BiDi endpoint of a running Firefox |
| `-timeout` | `2m` | Overall command timeout |

## What it demonstrates

- Connecting to an existing browser session
- Performing the WebDriver BiDi handshake
- Querying session status
