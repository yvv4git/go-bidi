# screenshot

Captures viewport and full-page screenshots.

## Run

```bash
go run examples/screenshot/main.go
```

## Examples

```bash
# Default: launch headless Firefox, capture screenshots
go run examples/screenshot/main.go

# Connect to running Firefox
go run examples/screenshot/main.go -endpoint http://127.0.0.1:9222

# Launch visible Firefox
go run examples/screenshot/main.go -headless=false

# Custom output directory
go run examples/screenshot/main.go -out /tmp/screenshots

# With custom timeout
go run examples/screenshot/main.go -timeout 30s

# Combined flags
go run examples/screenshot/main.go -endpoint http://127.0.0.1:9222 -out ./captures
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-endpoint` | | WebDriver BiDi endpoint of a running Firefox |
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-headless` | `true` | Run Firefox in headless mode |
| `-timeout` | `2m` | Overall command timeout |
| `-out` | `.` | Output directory |

## What it demonstrates

- Setting viewport dimensions
- Capturing viewport screenshot
- Capturing full-page screenshot
- Decoding PNG dimensions
