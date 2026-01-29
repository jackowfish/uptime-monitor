# Uptime Monitor

A terminal-based HTTP endpoint monitor built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Screenshots

<img width="835" height="499" alt="Screenshot 2026-01-29 at 9 43 09 AM" src="https://github.com/user-attachments/assets/0c74c1cc-e256-4854-8853-24d2582cc646" />
<img width="834" height="499" alt="Screenshot 2026-01-29 at 9 43 20 AM" src="https://github.com/user-attachments/assets/cb7dd294-660d-43c5-8f37-c002e5794cbd" />

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
