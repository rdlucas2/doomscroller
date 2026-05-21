package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type GameOverState struct {
	Character *game.Character
	Floor     int
	Cursor    int
}

var gameOverItems = []string{"New Game", "Main Menu", "Quit"}

func (s *GameOverState) MoveUp() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}

func (s *GameOverState) MoveDown() {
	if s.Cursor < len(gameOverItems)-1 {
		s.Cursor++
	}
}

func (s *GameOverState) Selected() string {
	if s.Cursor < len(gameOverItems) {
		return gameOverItems[s.Cursor]
	}
	return "Quit"
}

func (s *GameOverState) View(width, height int) string {
	lines := make([]string, 0, height)

	topPad := (height - 20) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	deathArt := []string{
		`   ██████╗  █████╗ ███╗   ███╗███████╗`,
		`  ██╔════╝ ██╔══██╗████╗ ████║██╔════╝`,
		`  ██║  ███╗███████║██╔████╔██║█████╗  `,
		`  ██║   ██║██╔══██║██║╚██╔╝██║██╔══╝  `,
		`  ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗`,
		`   ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝`,
		``,
		`        ██████╗ ██╗   ██╗███████╗██████╗ `,
		`       ██╔═══██╗██║   ██║██╔════╝██╔══██╗`,
		`       ██║   ██║██║   ██║█████╗  ██████╔╝`,
		`       ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗`,
		`       ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║`,
		`        ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝`,
	}
	for _, l := range deathArt {
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.ErrorStyle.Render(l),
		))
	}

	lines = append(lines, "")

	c := s.Character
	if c != nil {
		summary := fmt.Sprintf("%-16s reached Floor %d  |  Level %d  |  Gold: %d",
			c.Name, s.Floor, c.Level, c.Gold)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.SubtitleStyle.Render(summary),
		))
	}
	lines = append(lines, "")
	lines = append(lines, "")

	for i, item := range gameOverItems {
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
