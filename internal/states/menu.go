package states

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type MenuChoice int

const (
	MenuNewGame MenuChoice = iota
	MenuContinue
	MenuSettings
	MenuQuit
)

type MenuState struct {
	Cursor  int
	HasSave bool
}

func NewMenuState() *MenuState {
	return &MenuState{HasSave: game.HasSave()}
}

func (s *MenuState) Items() []string {
	items := []string{"New Game"}
	if s.HasSave {
		items = append(items, "Continue")
	}
	items = append(items, "Settings", "Quit")
	return items
}

func (s *MenuState) Selected() MenuChoice {
	items := s.Items()
	if s.Cursor >= len(items) {
		return MenuQuit
	}
	switch items[s.Cursor] {
	case "New Game":
		return MenuNewGame
	case "Continue":
		return MenuContinue
	case "Settings":
		return MenuSettings
	default:
		return MenuQuit
	}
}

func (s *MenuState) MoveUp() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}

func (s *MenuState) MoveDown() {
	items := s.Items()
	if s.Cursor < len(items)-1 {
		s.Cursor++
	}
}

func (s *MenuState) View(width, height int) string {
	lines := make([]string, 0, height)

	banner := []string{
		`  ╔══════════════════════════════════════════╗`,
		`  ║   DOOM SCROLLER  ─  Tower of the Void   ║`,
		`  ╚══════════════════════════════════════════╝`,
	}

	topPad := (height - 20) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	for _, l := range banner {
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(ui.TitleStyle.Render(l)))
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.SubtitleStyle.Render("Scale the tower. Break the Void."),
	))
	lines = append(lines, "")
	lines = append(lines, "")

	items := s.Items()
	for i, item := range items {
		var rendered string
		if i == s.Cursor {
			rendered = ui.SelectedStyle.Render("▶  " + item)
		} else {
			rendered = ui.NormalStyle.Render("   " + item)
		}
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(rendered))
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.DimStyle.Render("↑↓ Navigate   Enter Select   Esc Back"),
	))

	return strings.Join(lines, "\n")
}
