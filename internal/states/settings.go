package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type SettingsState struct {
	Settings *game.Settings
	Cursor   int
}

var settingsItems = []string{
	"Music Volume",
	"Sound Volume",
	"Difficulty",
	"Show Full Map",
	"Back",
}

func NewSettingsState(s *game.Settings) *SettingsState {
	return &SettingsState{Settings: s}
}

func (s *SettingsState) MoveUp() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}

func (s *SettingsState) MoveDown() {
	if s.Cursor < len(settingsItems)-1 {
		s.Cursor++
	}
}

func (s *SettingsState) AdjustLeft() {
	switch s.Cursor {
	case 0:
		if s.Settings.MusicVolume > 0 {
			s.Settings.MusicVolume -= 10
		}
	case 1:
		if s.Settings.SoundVolume > 0 {
			s.Settings.SoundVolume -= 10
		}
	case 2:
		if s.Settings.Difficulty > game.DifficultyEasy {
			s.Settings.Difficulty--
		}
	case 3:
		s.Settings.ShowFullMap = !s.Settings.ShowFullMap
	}
}

func (s *SettingsState) AdjustRight() {
	switch s.Cursor {
	case 0:
		if s.Settings.MusicVolume < 100 {
			s.Settings.MusicVolume += 10
		}
	case 1:
		if s.Settings.SoundVolume < 100 {
			s.Settings.SoundVolume += 10
		}
	case 2:
		if s.Settings.Difficulty < game.DifficultyHard {
			s.Settings.Difficulty++
		}
	case 3:
		s.Settings.ShowFullMap = !s.Settings.ShowFullMap
	}
}

func (s *SettingsState) IsBack() bool {
	return s.Cursor == len(settingsItems)-1
}

func (s *SettingsState) valueFor(idx int) string {
	switch idx {
	case 0:
		return fmt.Sprintf("◄ %3d%% ►  %s", s.Settings.MusicVolume, ui.HPBar(s.Settings.MusicVolume, 100, 10))
	case 1:
		return fmt.Sprintf("◄ %3d%% ►  %s", s.Settings.SoundVolume, ui.HPBar(s.Settings.SoundVolume, 100, 10))
	case 2:
		return fmt.Sprintf("◄ %-6s ►", s.Settings.Difficulty.String())
	case 3:
		if s.Settings.ShowFullMap {
			return "◄  ON  ►"
		}
		return "◄  OFF ►"
	}
	return ""
}

func (s *SettingsState) View(width, height int) string {
	lines := make([]string, 0, height)

	topPad := (height - 20) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	title := ui.TitleStyle.Render("⚙  SETTINGS")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(title))
	lines = append(lines, "")
	lines = append(lines, "")

	for i, item := range settingsItems {
		var row string
		if i == len(settingsItems)-1 {
			// Back button
			if i == s.Cursor {
				row = ui.SelectedStyle.Render("▶  " + item)
			} else {
				row = ui.NormalStyle.Render("   " + item)
			}
		} else {
			val := s.valueFor(i)
			label := fmt.Sprintf("%-16s  %s", item, val)
			if i == s.Cursor {
				row = ui.SelectedStyle.Render("▶  " + label)
			} else {
				row = ui.NormalStyle.Render("   " + label)
			}
		}
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(row))
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.DimStyle.Render("↑↓ Navigate   ◄► Adjust   Enter/Esc Back"),
	))

	return strings.Join(lines, "\n")
}
