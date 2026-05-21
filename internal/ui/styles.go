package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorGold    = lipgloss.Color("#FFD700")
	ColorRed     = lipgloss.Color("#FF4444")
	ColorGreen   = lipgloss.Color("#44FF88")
	ColorBlue    = lipgloss.Color("#4488FF")
	ColorPurple  = lipgloss.Color("#AA44FF")
	ColorCyan    = lipgloss.Color("#44FFFF")
	ColorWhite   = lipgloss.Color("#FFFFFF")
	ColorGray    = lipgloss.Color("#888888")
	ColorDark    = lipgloss.Color("#1A1A2E")
	ColorMid     = lipgloss.Color("#16213E")
	ColorAccent  = lipgloss.Color("#0F3460")

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorGold).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Italic(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorGold).
			Bold(true).
			PaddingLeft(2)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			PaddingLeft(2)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			PaddingLeft(2)

	HPStyle = lipgloss.NewStyle().
		Foreground(ColorRed).
		Bold(true)

	MPStyle = lipgloss.NewStyle().
		Foreground(ColorBlue).
		Bold(true)

	XPStyle = lipgloss.NewStyle().
		Foreground(ColorCyan)

	GoldStyle = lipgloss.NewStyle().
			Foreground(ColorGold)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(0, 1)

	HighlightBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGold).
			Padding(0, 1)

	BattleLogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	PerkCombatStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	PerkMagicStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	PerkUtilStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	PerkPassiveStyle = lipgloss.NewStyle().
				Foreground(ColorGold).
				Bold(true)
)

func Bar(current, max, width int, filledChar, emptyChar string, filledStyle, emptyStyle lipgloss.Style) string {
	if max <= 0 {
		return ""
	}
	filled := current * width / max
	if filled > width {
		filled = width
	}
	out := ""
	for i := 0; i < width; i++ {
		if i < filled {
			out += filledStyle.Render(filledChar)
		} else {
			out += emptyStyle.Render(emptyChar)
		}
	}
	return out
}

func HPBar(hp, maxHP, width int) string {
	return Bar(hp, maxHP, width, "█", "░",
		lipgloss.NewStyle().Foreground(ColorRed),
		lipgloss.NewStyle().Foreground(ColorGray))
}

func MPBar(mp, maxMP, width int) string {
	return Bar(mp, maxMP, width, "█", "░",
		lipgloss.NewStyle().Foreground(ColorBlue),
		lipgloss.NewStyle().Foreground(ColorGray))
}

func XPBar(xp, xpMax, width int) string {
	return Bar(xp, xpMax, width, "▪", "·",
		lipgloss.NewStyle().Foreground(ColorCyan),
		lipgloss.NewStyle().Foreground(ColorGray))
}

func Center(s string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
}
