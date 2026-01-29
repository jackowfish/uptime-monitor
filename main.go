package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"uptime_monitor/internal/monitor"
)

var (
	cyan    = lipgloss.Color("#00d7d7")
	muted   = lipgloss.Color("#6c6c6c")
	white   = lipgloss.Color("#e4e4e4")
	red     = lipgloss.Color("#ff5f5f")
	green   = lipgloss.Color("#00d787")
	headerBg = lipgloss.Color("#262626")

	logoStyle = lipgloss.NewStyle().
			Background(cyan).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(muted)

	promptStyle = lipgloss.NewStyle().
			Foreground(white)

	inputStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(muted).
			Italic(true)
)

func printUsage() {
	fmt.Println()
	fmt.Println(logoStyle.Render("UPTIME") + " " + titleStyle.Render("monitor"))
	fmt.Println(subtitleStyle.Render("Real-time HTTP endpoint monitoring"))
	fmt.Println()
	fmt.Println(promptStyle.Render("Usage:"))
	fmt.Println(subtitleStyle.Render("  uptime_monitor [url]"))
	fmt.Println(subtitleStyle.Render("  uptime_monitor"))
	fmt.Println()
	fmt.Println(promptStyle.Render("Examples:"))
	fmt.Println(subtitleStyle.Render("  uptime_monitor https://example.com"))
	fmt.Println(subtitleStyle.Render("  uptime_monitor example.com/health"))
	fmt.Println(subtitleStyle.Render("  TARGET_URL=https://example.com uptime_monitor"))
	fmt.Println()
	fmt.Println(promptStyle.Render("Keyboard shortcuts:"))
	fmt.Println(subtitleStyle.Render("  q        Quit"))
	fmt.Println(subtitleStyle.Render("  ?        Show help"))
	fmt.Println(subtitleStyle.Render("  t        Toggle TLS verification"))
	fmt.Println(subtitleStyle.Render("  f        Toggle follow redirects"))
	fmt.Println(subtitleStyle.Render("  +/-      Adjust check interval"))
	fmt.Println(subtitleStyle.Render("  [/]      Adjust timeout"))
	fmt.Println()
}

func promptURL(reader *bufio.Reader) string {
	fmt.Println()
	fmt.Println(logoStyle.Render("UPTIME") + " " + titleStyle.Render("monitor"))
	fmt.Println(subtitleStyle.Render("Real-time HTTP endpoint monitoring"))
	fmt.Println()
	fmt.Print(promptStyle.Render("Enter URL to monitor: "))

	input, _ := reader.ReadString('\n')
	url := strings.TrimSpace(input)

	if url == "" {
		return ""
	}

	return normalizeURL(url)
}

func normalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}

	// Add https:// if no scheme provided
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return url
}

func main() {
	args := os.Args[1:]

	// Handle help flag
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			printUsage()
			os.Exit(0)
		}
	}

	var targetURL string

	// Priority: 1. CLI argument, 2. Environment variable, 3. Interactive prompt
	if len(args) > 0 {
		targetURL = normalizeURL(args[0])
	} else if envURL := os.Getenv("TARGET_URL"); envURL != "" {
		targetURL = normalizeURL(envURL)
	} else {
		reader := bufio.NewReader(os.Stdin)
		targetURL = promptURL(reader)
		if targetURL == "" {
			fmt.Println()
			fmt.Println(errorStyle.Render("Error: URL is required"))
			fmt.Println(hintStyle.Render("Run with --help for usage information"))
			fmt.Println()
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println(subtitleStyle.Render("Starting monitor for ") + inputStyle.Render(targetURL))
		fmt.Println(hintStyle.Render("Press ? for help, q to quit"))
		fmt.Println()
	}

	p := tea.NewProgram(monitor.NewModel(targetURL), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
