package monitor

import (
	"fmt"
	"strings"
	"time"
)

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func SafePadding(total, used, minPad int) int {
	pad := total - used
	if pad < minPad {
		return minPad
	}
	return pad
}

func Truncate(s string, maxLen int) string {
	if maxLen < 4 {
		return s[:maxLen]
	}
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func RenderSection(title string, width int) string {
	return SectionStyle.Render("  "+title) + "\n" +
		DividerStyle.Render("  "+strings.Repeat("─", Clamp(width-4, 1, 200))) + "\n"
}

func ValueOrDash(hasData bool, value string) string {
	if !hasData {
		return "—"
	}
	return value
}

func RenderSparkline(latencies []time.Duration, width int) string {
	if len(latencies) == 0 {
		return strings.Repeat("─", width)
	}

	sparks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	var minLat, maxLat time.Duration = latencies[0], latencies[0]
	for _, l := range latencies {
		if l < minLat {
			minLat = l
		}
		if l > maxLat {
			maxLat = l
		}
	}

	latRange := maxLat - minLat
	if latRange == 0 {
		latRange = 1
	}

	var result strings.Builder
	data := latencies
	if len(data) > width {
		data = data[len(data)-width:]
	}

	for _, l := range data {
		scaled := float64(l-minLat) / float64(latRange) * 7
		idx := int(scaled)
		if idx > 7 {
			idx = 7
		}

		spark := string(sparks[idx])
		if l < 100*time.Millisecond {
			result.WriteString(SparkLowStyle.Render(spark))
		} else if l < 300*time.Millisecond {
			result.WriteString(SparkMedStyle.Render(spark))
		} else {
			result.WriteString(SparkHighStyle.Render(spark))
		}
	}

	remaining := width - len(data)
	if remaining > 0 {
		result.WriteString(GraphAxisStyle.Render(strings.Repeat("─", remaining)))
	}

	return result.String()
}

func RenderProgressBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	var bar strings.Builder
	bar.WriteString(ProgressFullStyle.Render(strings.Repeat("█", filled)))
	bar.WriteString(ProgressEmptyStyle.Render(strings.Repeat("░", width-filled)))
	return bar.String()
}

func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0ms"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func FormatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func RenderGraph(data []float64, width, height int) string {
	if len(data) == 0 {
		var lines []string
		for i := 0; i < height; i++ {
			lines = append(lines, strings.Repeat(" ", width))
		}
		return strings.Join(lines, "\n")
	}

	blocks := []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	dataLen := len(data)
	startCol := width - dataLen
	if startCol < 0 {
		startCol = 0
		data = data[dataLen-width:]
		dataLen = width
	}

	for i, val := range data {
		col := startCol + i
		if col >= width {
			break
		}

		scaledVal := val / 100.0 * float64(height)
		filledRows := int(scaledVal)
		partial := scaledVal - float64(filledRows)

		for row := 0; row < filledRows && row < height; row++ {
			grid[row][col] = GraphBarStyle.Render("█")
		}

		if filledRows < height && partial > 0 {
			blockIdx := int(partial * 8)
			if blockIdx > 0 && blockIdx < len(blocks) {
				grid[filledRows][col] = GraphBarStyle.Render(blocks[blockIdx])
			}
		}
	}

	var lines []string
	for row := height - 1; row >= 0; row-- {
		lines = append(lines, strings.Join(grid[row], ""))
	}

	return strings.Join(lines, "\n")
}
