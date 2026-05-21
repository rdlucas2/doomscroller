package states

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type PauseChoice int

const (
	PauseResume PauseChoice = iota
	PauseSave
	PauseSettings
	PauseQuitToMenu
)

type PauseState struct {
	Cursor int
}

var pauseItems = []string{
	"Resume",
	"Save & Quit Run",
	"Settings",
	"Quit to Main Menu",
}

func (ps *PauseState) MoveUp() {
	if ps.Cursor > 0 {
		ps.Cursor--
	}
}

func (ps *PauseState) MoveDown() {
	if ps.Cursor < len(pauseItems)-1 {
		ps.Cursor++
	}
}

func (ps *PauseState) Selected() PauseChoice {
	return PauseChoice(ps.Cursor)
}

func (ps *PauseState) View(width, height int) string {
	lines := make([]string, 0, height)

	topPad := (height - 16) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.TitleStyle.Render("⏸  PAUSED"),
	))
	lines = append(lines, "")
	lines = append(lines, "")

	for i, item := range pauseItems {
		var rendered string
		if i == ps.Cursor {
			rendered = ui.SelectedStyle.Render("▶  " + item)
		} else {
			rendered = ui.NormalStyle.Render("   " + item)
		}
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(rendered))
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.DimStyle.Render("↑↓ Navigate   Enter Select   P Resume"),
	))

	return strings.Join(lines, "\n")
}
