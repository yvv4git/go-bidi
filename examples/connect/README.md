# connect

Connects to an already-running Firefox via remote debugging.

## Run

```bash
go run examples/connect/main.go
```

Firefox must be running with `--remote-debugging-port 9222`.

## What it demonstrates

- Connecting to an existing browser session
- Performing the WebDriver BiDi handshake
- Querying session status
