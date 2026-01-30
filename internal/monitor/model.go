package monitor

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	TargetURL       string
	Spinner         spinner.Model
	TotalReqs       int
	SuccessReqs     int
	FailedReqs      int
	LastStatus      bool
	LastLatency     time.Duration
	AvgLatency      time.Duration
	MinLatency      time.Duration
	MaxLatency      time.Duration
	LatencySum      time.Duration
	UptimeHistory   []float64
	LatencyHistory  []time.Duration
	RecentResults   []bool // Rolling window of recent success/failure
	Logs            []LogEntry
	Quitting        bool
	ShowHelp        bool
	Width           int
	Height          int
	StartTime       time.Time
	ConsecutiveUp   int
	ConsecutiveDown int
	InsecureTLS     bool
	FollowRedirects bool
	Timeout         time.Duration
	Interval        time.Duration
	Workers         int
}

func NewModel(targetURL string) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(ColorCyan)

	return Model{
		TargetURL:       targetURL,
		Spinner:         s,
		UptimeHistory:   make([]float64, 0, MaxHistory),
		LatencyHistory:  make([]time.Duration, 0, MaxLatencies),
		RecentResults:   make([]bool, 0, UptimeWindow),
		Logs:            make([]LogEntry, 0, MaxLogs),
		Width:           80,
		Height:          24,
		StartTime:       time.Now(),
		MinLatency:      time.Hour,
		InsecureTLS:     true,
		FollowRedirects: true,
		Timeout:         5 * time.Second,
		Interval:        CheckInterval,
		Workers:         1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, m.doCheck(), m.tickCmd())
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(m.Interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) doCheck() tea.Cmd {
	return func() tea.Msg {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: m.InsecureTLS},
		}

		client := &http.Client{
			Timeout:   m.Timeout,
			Transport: transport,
		}

		if !m.FollowRedirects {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}

		start := time.Now()
		resp, err := client.Get(m.TargetURL)
		latency := time.Since(start)

		if err != nil {
			return CheckResult{
				Success:   false,
				Err:       err,
				Latency:   latency,
				Timestamp: time.Now(),
			}
		}
		defer resp.Body.Close()

		success := resp.StatusCode >= 200 && resp.StatusCode < 300

		result := CheckResult{
			Success:    success,
			StatusCode: resp.StatusCode,
			Latency:    latency,
			Timestamp:  time.Now(),
		}

		if !success {
			result.Err = fmt.Errorf("HTTP %d", resp.StatusCode)
		}

		return result
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "?":
			m.ShowHelp = !m.ShowHelp
			return m, nil
		case "t":
			m.InsecureTLS = !m.InsecureTLS
			return m, nil
		case "f":
			m.FollowRedirects = !m.FollowRedirects
			return m, nil
		case "+", "=":
			if m.Interval < 10*time.Millisecond {
				m.Interval += 1 * time.Millisecond
			} else if m.Interval < 100*time.Millisecond {
				m.Interval += 10 * time.Millisecond
			} else if m.Interval < 1*time.Second {
				m.Interval += 100 * time.Millisecond
			} else if m.Interval < 10*time.Second {
				m.Interval += 500 * time.Millisecond
			}
			return m, nil
		case "-", "_":
			if m.Interval > 100*time.Millisecond {
				m.Interval -= 100 * time.Millisecond
			} else if m.Interval > 10*time.Millisecond {
				m.Interval -= 10 * time.Millisecond
			} else if m.Interval > 1*time.Millisecond {
				m.Interval -= 1 * time.Millisecond
			}
			return m, nil
		case ">", ".":
			if m.Workers < 1000 {
				if m.Workers < 10 {
					m.Workers++
				} else if m.Workers < 100 {
					m.Workers += 10
				} else {
					m.Workers += 100
				}
			}
			return m, nil
		case "<", ",":
			if m.Workers > 1 {
				if m.Workers <= 10 {
					m.Workers--
				} else if m.Workers <= 100 {
					m.Workers -= 10
				} else {
					m.Workers -= 100
				}
			}
			return m, nil
		case "[":
			if m.Timeout > 1*time.Second {
				m.Timeout -= 1 * time.Second
			}
			return m, nil
		case "]":
			if m.Timeout < 30*time.Second {
				m.Timeout += 1 * time.Second
			}
			return m, nil
		case "esc":
			if m.ShowHelp {
				m.ShowHelp = false
			}
			return m, nil
		}

	case TickMsg:
		cmds := make([]tea.Cmd, m.Workers+1)
		for i := 0; i < m.Workers; i++ {
			cmds[i] = m.doCheck()
		}
		cmds[m.Workers] = m.tickCmd()
		return m, tea.Batch(cmds...)

	case CheckResult:
		m.TotalReqs++
		m.LastLatency = msg.Latency

		if msg.Success {
			m.LatencySum += msg.Latency
			if msg.Latency < m.MinLatency {
				m.MinLatency = msg.Latency
			}
			if msg.Latency > m.MaxLatency {
				m.MaxLatency = msg.Latency
			}
		}

		m.LatencyHistory = append(m.LatencyHistory, msg.Latency)
		if len(m.LatencyHistory) > MaxLatencies {
			m.LatencyHistory = m.LatencyHistory[1:]
		}

		if msg.Success {
			m.SuccessReqs++
			m.AvgLatency = m.LatencySum / time.Duration(m.SuccessReqs)
			m.LastStatus = true
			m.ConsecutiveUp++
			m.ConsecutiveDown = 0
		} else {
			m.FailedReqs++
			m.LastStatus = false
			m.ConsecutiveDown++
			m.ConsecutiveUp = 0
		}

		// Track rolling window of recent results
		m.RecentResults = append(m.RecentResults, msg.Success)
		if len(m.RecentResults) > UptimeWindow {
			m.RecentResults = m.RecentResults[1:]
		}

		// Calculate uptime from rolling window (last 1000 requests)
		recentSuccesses := 0
		for _, success := range m.RecentResults {
			if success {
				recentSuccesses++
			}
		}
		uptime := float64(recentSuccesses) / float64(len(m.RecentResults)) * 100

		m.UptimeHistory = append(m.UptimeHistory, uptime)
		if len(m.UptimeHistory) > MaxHistory {
			m.UptimeHistory = m.UptimeHistory[1:]
		}

		entry := LogEntry{
			Timestamp:  msg.Timestamp,
			Success:    msg.Success,
			StatusCode: msg.StatusCode,
			Latency:    msg.Latency,
		}
		if msg.Err != nil {
			entry.ErrMsg = msg.Err.Error()
		}
		m.Logs = append(m.Logs, entry)
		if len(m.Logs) > MaxLogs {
			m.Logs = m.Logs[1:]
		}

		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	width := m.Width
	if width < 70 {
		width = 70
	}

	var b strings.Builder

	b.WriteString(m.renderHeader(width))
	b.WriteString("\n")
	b.WriteString(m.renderSecondaryHeader(width))
	b.WriteString("\n\n")
	b.WriteString(m.renderStatistics(width))
	b.WriteString(m.renderLatency(width))
	b.WriteString(m.renderUptimeGraph(width))
	b.WriteString(m.renderActivityLog(width))
	b.WriteString(m.renderFooter(width))

	if m.ShowHelp {
		return m.overlayHelp(b.String(), width)
	}

	return b.String()
}

func (m Model) renderHeader(width int) string {
	// Status badge
	statusBadge := "INIT"
	statusStyle := HeaderInfoStyle
	if m.TotalReqs > 0 {
		if m.LastStatus {
			if m.ConsecutiveUp > 5 {
				statusBadge = "HEALTHY"
			} else {
				statusBadge = "ONLINE"
			}
			statusStyle = StatusOnlineStyle
		} else {
			if m.ConsecutiveDown > 3 {
				statusBadge = "CRITICAL"
			} else {
				statusBadge = "OFFLINE"
			}
			statusStyle = StatusOfflineStyle
		}
	}

	// Fixed elements with their visual widths
	logoBadge := LogoStyle.Render("UPTIME")
	logoWidth := lipgloss.Width(logoBadge)
	statusRendered := statusStyle.Render(statusBadge)
	statusWidth := lipgloss.Width(statusRendered)

	// Reserve space: logo + 2 spaces + url + 2 spaces + status + 2 margins
	reserved := logoWidth + statusWidth + 6
	maxURLWidth := width - reserved
	if maxURLWidth < 10 {
		maxURLWidth = 10
	}

	// Truncate URL to fit
	displayURL := m.TargetURL
	if len(displayURL) > maxURLWidth {
		displayURL = Truncate(displayURL, maxURLWidth)
	}
	urlRendered := URLStyle.Render(displayURL)
	urlWidth := lipgloss.Width(urlRendered)

	// Calculate actual padding needed
	usedWidth := logoWidth + 2 + urlWidth + 2 + statusWidth
	padding := width - usedWidth - 2 // -2 for outer margins
	if padding < 1 {
		padding = 1
	}

	content := " " + logoBadge + "  " + urlRendered + strings.Repeat(" ", padding) + statusRendered + " "
	return HeaderStyle.Width(width).Render(content)
}

func (m Model) renderSecondaryHeader(width int) string {
	timeStr := time.Now().Format("15:04:05")
	dateStr := time.Now().Format("2006-01-02")
	runningFor := FormatUptime(time.Since(m.StartTime))

	tlsIndicator := ValueGreenStyle.Render("tls:skip")
	if !m.InsecureTLS {
		tlsIndicator = ValueRedStyle.Render("tls:verify")
	}
	redirIndicator := HeaderInfoStyle.Render("redir:on")
	if !m.FollowRedirects {
		redirIndicator = HeaderInfoStyle.Render("redir:off")
	}

	// Calculate requests per second: workers * (1000ms / interval)
	rps := float64(m.Workers) * (1000.0 / float64(m.Interval.Milliseconds()))
	rpsStr := fmt.Sprintf("%.0f", rps)
	if rps >= 1000 {
		rpsStr = fmt.Sprintf("%.1fk", rps/1000)
	}

	rightInfo := fmt.Sprintf("%s │ %s │ %s │ %s │ %s %s %dms %ds %dw %s",
		HeaderInfoStyle.Render(Version),
		HeaderInfoStyle.Render(dateStr),
		ValueStyle.Render(timeStr),
		HeaderInfoStyle.Render("↑"+runningFor),
		tlsIndicator,
		redirIndicator,
		m.Interval.Milliseconds(),
		int(m.Timeout.Seconds()),
		m.Workers,
		ValueGreenStyle.Render(rpsStr+"rps"),
	)

	rightInfoPadding := SafePadding(width, lipgloss.Width(rightInfo)+2, 0)
	return FooterStyle.Width(width).Render(" " + rightInfo + strings.Repeat(" ", rightInfoPadding))
}

func (m Model) renderStatistics(width int) string {
	var b strings.Builder

	uptime := 0.0
	if m.TotalReqs > 0 {
		uptime = float64(m.SuccessReqs) / float64(m.TotalReqs) * 100
	}

	b.WriteString(RenderSection("STATISTICS", width))

	progressWidth := 25
	b.WriteString(fmt.Sprintf("  %s %s %s\n",
		LabelStyle.Width(12).Render("Uptime:"),
		RenderProgressBar(uptime, progressWidth),
		ValueGreenStyle.Render(fmt.Sprintf(" %.2f%%", uptime)),
	))

	b.WriteString(fmt.Sprintf("  %s %s          %s %s          %s %s\n",
		LabelStyle.Render("Total:"),
		ValueStyle.Render(fmt.Sprintf("%-6d", m.TotalReqs)),
		LabelStyle.Render("Success:"),
		ValueGreenStyle.Render(fmt.Sprintf("%-6d", m.SuccessReqs)),
		LabelStyle.Render("Failed:"),
		ValueRedStyle.Render(fmt.Sprintf("%-6d", m.FailedReqs)),
	))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderLatency(width int) string {
	var b strings.Builder

	b.WriteString(RenderSection("LATENCY", width))

	hasData := m.TotalReqs > 0
	minLat := ValueOrDash(hasData, FormatDuration(m.MinLatency))
	maxLat := ValueOrDash(hasData, FormatDuration(m.MaxLatency))
	avgLat := ValueOrDash(hasData, FormatDuration(m.AvgLatency))
	lastLat := ValueOrDash(hasData, FormatDuration(m.LastLatency))

	b.WriteString(fmt.Sprintf("  %s %s    %s %s    %s %s    %s %s\n",
		LabelStyle.Render("Last:"),
		LatencyStyle.Render(fmt.Sprintf("%-8s", lastLat)),
		LabelStyle.Render("Avg:"),
		LatencyStyle.Render(fmt.Sprintf("%-8s", avgLat)),
		LabelStyle.Render("Min:"),
		ValueGreenStyle.Render(fmt.Sprintf("%-8s", minLat)),
		LabelStyle.Render("Max:"),
		ValueRedStyle.Render(fmt.Sprintf("%-8s", maxLat)),
	))

	sparkWidth := Clamp(width-16, 10, 60)
	b.WriteString(fmt.Sprintf("  %s %s\n",
		LabelStyle.Render("Trend:"),
		RenderSparkline(m.LatencyHistory, sparkWidth),
	))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderUptimeGraph(width int) string {
	var b strings.Builder

	b.WriteString(SectionStyle.Render("  UPTIME HISTORY"))
	b.WriteString(GraphAxisStyle.Render(fmt.Sprintf(" (last %d requests)", len(m.RecentResults))))
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("  " + strings.Repeat("─", Clamp(width-4, 1, 200))))
	b.WriteString("\n")

	graphWidth := Clamp(width-14, 30, 55)
	graph := RenderGraph(m.UptimeHistory, graphWidth, GraphHeight)
	graphLines := strings.Split(graph, "\n")

	yLabels := []string{"100", "", "50", "", "0"}
	if GraphHeight != 5 {
		yLabels = make([]string, GraphHeight)
		yLabels[0] = "100"
		yLabels[GraphHeight-1] = "0"
	}

	for i, line := range graphLines {
		label := ""
		if i < len(yLabels) {
			label = yLabels[i]
		}
		b.WriteString(GraphAxisStyle.Render(fmt.Sprintf("  %4s │", label)))
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString(GraphAxisStyle.Render("       └" + strings.Repeat("─", graphWidth)))
	b.WriteString("\n\n")

	return b.String()
}

func (m Model) renderActivityLog(width int) string {
	var b strings.Builder

	availableLogLines := Clamp(m.Height-27, 1, MaxLogs)
	b.WriteString(RenderSection("RECENT ACTIVITY", width))

	if len(m.Logs) == 0 {
		b.WriteString(GraphAxisStyle.Render("  Waiting for first check..."))
		b.WriteString("\n")
	} else {
		logsToShow := m.Logs
		if len(logsToShow) > availableLogLines {
			logsToShow = logsToShow[len(logsToShow)-availableLogLines:]
		}
		for _, entry := range logsToShow {
			ts := LogTimeStyle.Render(entry.Timestamp.Format("15:04:05.000"))
			lat := LatencyStyle.Render(fmt.Sprintf("[%s]", FormatDuration(entry.Latency)))
			if entry.Success {
				b.WriteString(fmt.Sprintf("  %s %s %s %s\n",
					ts,
					LogSuccessStyle.Render("✔"),
					lat,
					LogSuccessStyle.Render(fmt.Sprintf("HTTP %d OK", entry.StatusCode)),
				))
			} else {
				var errDisplay string
				if entry.StatusCode > 0 {
					errDisplay = fmt.Sprintf("HTTP %d", entry.StatusCode)
				} else {
					errDisplay = Truncate(entry.ErrMsg, 40)
				}
				b.WriteString(fmt.Sprintf("  %s %s %s %s\n",
					ts,
					LogErrorStyle.Render("✖"),
					lat,
					LogErrorStyle.Render(errDisplay),
				))
			}
		}
	}

	return b.String()
}

func (m Model) renderFooter(width int) string {
	footerKeys := []struct{ key, desc string }{
		{"q", "Quit"},
		{"?", "Help"},
		{"t", "TLS"},
		{"f", "Redir"},
		{"+/-", "Int"},
		{"[/]", "Tout"},
		{"</>", "Wrk"},
	}

	var footerParts []string
	for _, fk := range footerKeys {
		footerParts = append(footerParts, fmt.Sprintf("%s %s",
			FooterKeyStyle.Render("<"+fk.key+">"),
			FooterDescStyle.Render(fk.desc),
		))
	}
	footerContent := " " + strings.Join(footerParts, "   ")

	footerRight := LogoTextStyle.Render("UPTIME") + FooterDescStyle.Render(" monitor")
	footerPadding := SafePadding(width, lipgloss.Width(footerContent)+lipgloss.Width(footerRight)+2, 1)

	return FooterStyle.Width(width).Render(footerContent + strings.Repeat(" ", footerPadding) + footerRight + " ")
}

func (m Model) renderHelp(width int) string {
	var b strings.Builder

	b.WriteString(HelpTitleStyle.Render("KEYBOARD SHORTCUTS"))
	b.WriteString("\n\n")

	tlsStatus := "ON"
	if !m.InsecureTLS {
		tlsStatus = "OFF"
	}
	redirectStatus := "ON"
	if !m.FollowRedirects {
		redirectStatus = "OFF"
	}

	shortcuts := []struct{ key, desc string }{
		{"q / Ctrl+C", "Quit application"},
		{"?", "Toggle this help"},
		{"", ""},
		{"t", fmt.Sprintf("Toggle skip TLS verify [%s]", tlsStatus)},
		{"f", fmt.Sprintf("Toggle follow redirects [%s]", redirectStatus)},
		{"+/-", fmt.Sprintf("Adjust check interval [%dms]", m.Interval.Milliseconds())},
		{"[/]", fmt.Sprintf("Adjust timeout [%ds]", int(m.Timeout.Seconds()))},
		{"</>", fmt.Sprintf("Adjust workers [%d]", m.Workers)},
		{"", ""},
		{"Esc", "Close help overlay"},
	}

	for _, s := range shortcuts {
		if s.key == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString(HelpKeyStyle.Render(s.key))
		b.WriteString(HelpDescStyle.Render(s.desc))
		b.WriteString("\n")
	}

	return HelpStyle.Width(width - 10).Render(b.String())
}

func (m Model) overlayHelp(output string, width int) string {
	helpBox := m.renderHelp(width)
	helpLines := strings.Split(helpBox, "\n")
	helpHeight := len(helpLines)
	helpWidth := lipgloss.Width(helpLines[0])

	outputLines := strings.Split(output, "\n")

	startRow := (len(outputLines) - helpHeight) / 2
	if startRow < 2 {
		startRow = 2
	}
	startCol := (width - helpWidth) / 2
	if startCol < 0 {
		startCol = 0
	}

	for i, helpLine := range helpLines {
		targetRow := startRow + i
		if targetRow < len(outputLines) {
			outputLines[targetRow] = strings.Repeat(" ", startCol) + helpLine
		}
	}

	return strings.Join(outputLines, "\n")
}
