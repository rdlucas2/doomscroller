package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type WinState struct {
	Character *game.Character
	Cursor    int
}

var winItems = []string{"New Game", "Main Menu", "Quit"}

func (s *WinState) MoveUp() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}

func (s *WinState) MoveDown() {
	if s.Cursor < len(winItems)-1 {
		s.Cursor++
	}
}

func (s *WinState) Selected() string {
	if s.Cursor < len(winItems) {
		return winItems[s.Cursor]
	}
	return "Quit"
}

func (s *WinState) View(width, height int) string {
	lines := make([]string, 0, height)

	topPad := (height - 22) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	victoryArt := []string{
		`  ██╗   ██╗██╗ ██████╗████████╗ ██████╗ ██████╗ ██╗   ██╗`,
		`  ██║   ██║██║██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗╚██╗ ██╔╝`,
		`  ██║   ██║██║██║        ██║   ██║   ██║██████╔╝ ╚████╔╝ `,
		`  ╚██╗ ██╔╝██║██║        ██║   ██║   ██║██╔══██╗  ╚██╔╝  `,
		`   ╚████╔╝ ██║╚██████╗   ██║   ╚██████╔╝██║  ██║   ██║   `,
		`    ╚═══╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   `,
	}
	for _, l := range victoryArt {
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.GoldStyle.Render(l),
		))
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.SubtitleStyle.Render("The Void Sovereign has fallen. The tower is at peace."),
	))
	lines = append(lines, "")

	c := s.Character
	if c != nil {
		perksStr := fmt.Sprintf("%d", len(c.Perks))
		summary := fmt.Sprintf(
			"%s  |  Level %d  |  Gold: %d  |  Perks: %s",
			ui.TitleStyle.Render(c.Name), c.Level, c.Gold, ui.GoldStyle.Render(perksStr),
		)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(summary))
		lines = append(lines, "")

		if len(c.Perks) > 0 {
			perkNames := make([]string, len(c.Perks))
			for i, p := range c.Perks {
				perkNames[i] = perkTypeStyle(p.Type).Render(p.Name)
			}
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
				"Perks: "+strings.Join(perkNames, "  "),
			))
			lines = append(lines, "")
		}
	}

	lines = append(lines, "")
	for i, item := range winItems {
		var rendered string
		if i == s.Cursor {
			rendered = ui.SelectedStyle.Render("▶  " + item)
		} else {
			rendered = ui.NormalStyle.Render("   " + item)
		}
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(rendered))
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}
