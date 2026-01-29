# Uptime Monitor

A terminal-based HTTP endpoint monitor built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Usage

```bash
# Pass URL directly
uptime_monitor https://example.com

# Or use environment variable
TARGET_URL=https://example.com uptime_monitor

# Or run interactively
uptime_monitor
```

## Build

```bash
make build
```

## Keyboard Shortcuts

- `q` - Quit
- `?` - Show help
- `t` - Toggle TLS verification
- `f` - Toggle follow redirects
- `+/-` - Adjust check interval
- `[/]` - Adjust timeout
