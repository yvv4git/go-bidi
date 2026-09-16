# navigate

Opens a tab and navigates to Wikipedia.

## Run

```bash
go run examples/navigate/main.go
```

## Examples

```bash
# Default: launch headless Firefox, navigate to Wikipedia
go run examples/navigate/main.go

# Connect to running Firefox
go run examples/navigate/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox (see the window)
go run examples/navigate/main.go -headless=false

# With custom timeout
go run examples/navigate/main.go -timeout 30s
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |

## What it demonstrates

- Creating a new browsing context (tab)
- Navigating to an external URL
- Printing the loaded URL
