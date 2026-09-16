# screenshot

Captures viewport and full-page screenshots.

## Run

```bash
go run examples/screenshot/main.go
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-exec` | (auto-detect) | Path to Firefox binary |
| `-timeout` | `2m` | Overall command timeout |
| `-out` | `.` | Output directory |

## What it demonstrates

- Setting viewport dimensions
- Capturing viewport screenshot
- Capturing full-page screenshot
- Decoding PNG dimensions
