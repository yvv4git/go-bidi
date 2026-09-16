# evaluate

Runs JavaScript in the page and decodes returned values into Go types.

## Run

```bash
go run examples/evaluate/main.go
```

## Examples

```bash
# Default: launch headless Firefox, run JS
go run examples/evaluate/main.go

# Connect to running Firefox
go run examples/evaluate/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox
go run examples/evaluate/main.go -headless=false

# With custom timeout
go run examples/evaluate/main.go -timeout 30s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |

## What it demonstrates

- Evaluating JavaScript expressions
- Decoding RemoteValue into native Go values
- Calling functions with arguments
- Handling Promises
