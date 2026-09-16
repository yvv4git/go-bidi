# input

Demonstrates typing text, pressing keys and clicking on a page.

## Run

```bash
go run examples/input/main.go
```

## Examples

```bash
# Default: launch headless Firefox, demonstrate input
go run examples/input/main.go

# Connect to running Firefox
go run examples/input/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox
go run examples/input/main.go -headless=false

# With custom timeout
go run examples/input/main.go -timeout 30s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |

## What it demonstrates

- Typing text into an input field
- Pressing keyboard keys (Enter)
- Clicking at specific coordinates
- Evaluating page state after interactions
