package monitor

import "github.com/charmbracelet/lipgloss"

var (
	ColorCyan     = lipgloss.Color("#00d7d7")
	ColorYellow   = lipgloss.Color("#d7af00")
	ColorGreen    = lipgloss.Color("#00d787")
	ColorRed      = lipgloss.Color("#ff5f5f")
	ColorBlue     = lipgloss.Color("#5f87d7")
	ColorMagenta  = lipgloss.Color("#d787d7")
	ColorOrange   = lipgloss.Color("#ff8700")
	ColorMuted    = lipgloss.Color("#6c6c6c")
	ColorDimmed   = lipgloss.Color("#4e4e4e")
	ColorWhite    = lipgloss.Color("#e4e4e4")
	ColorHeaderBg = lipgloss.Color("#262626")
	ColorPanelBg  = lipgloss.Color("#1a1a1a")
)

var (
	HeaderStyle = lipgloss.NewStyle().
			Background(ColorHeaderBg).
			Foreground(ColorCyan).
			Bold(true)

	LogoStyle = lipgloss.NewStyle().
			Background(ColorCyan).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	HeaderInfoStyle = lipgloss.NewStyle().
			Background(ColorHeaderBg).
			Foreground(ColorMuted)

	SectionStyle = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Width(10)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Bold(true)

	ValueGreenStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	ValueRedStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	StatusOnlineStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Bold(true)

	StatusOfflineStyle = lipgloss.NewStyle().
				Foreground(ColorRed).
				Bold(true)

	URLStyle = lipgloss.NewStyle().
			Foreground(ColorBlue)

	GraphBarStyle = lipgloss.NewStyle().
			Foreground(ColorCyan)

	GraphAxisStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed)

	LogTimeStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	LogSuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	LogErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	FooterStyle = lipgloss.NewStyle().
			Background(ColorHeaderBg).
			Foreground(ColorMuted)

	FooterKeyStyle = lipgloss.NewStyle().
			Background(ColorHeaderBg).
			Foreground(ColorCyan).
			Bold(true)

	FooterDescStyle = lipgloss.NewStyle().
			Background(ColorHeaderBg).
			Foreground(ColorMuted)

	DividerStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed)

	LogoTextStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)

	LatencyStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta)

	ProgressFullStyle = lipgloss.NewStyle().
				Foreground(ColorGreen)

	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorDimmed)

	HelpStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorCyan).
			Padding(1, 2).
			Background(ColorPanelBg)

	HelpTitleStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true).
			MarginBottom(1)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Bold(true).
			Width(12)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ColorWhite)

	SparkLowStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	SparkMedStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	SparkHighStyle = lipgloss.NewStyle().
			Foreground(ColorOrange)
)
