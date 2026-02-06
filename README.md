# Uptime Monitor

A terminal-based HTTP endpoint monitor built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Screenshots
<img width="1306" height="790" alt="uptime_monitor (Window) 2026-02-06 10:06 AM" src="https://github.com/user-attachments/assets/1057f953-6253-494a-93a1-5f39d7463522" />
<img width="1306" height="790" alt="uptime_monitor (Window) 2026-02-06 10:07 AM" src="https://github.com/user-attachments/assets/b9f23204-eddb-46e2-ad25-69f0c7fcf1be" />



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
