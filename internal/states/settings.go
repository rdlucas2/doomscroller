package states

import (
	"fmt"
	"github.com/rdlucas2/doomscroller/internal/game"
)

type SettingsState struct {
	Settings *game.Settings
	Cursor   int
}

var settingLabels = []string{"Music Volume", "Sound Volume", "Difficulty", "Show Full Map", "Back"}

func NewSettingsState(s *game.Settings) *SettingsState { return &SettingsState{Settings: s} }

func (s *SettingsState) MoveUp() {
	if s.Cursor > 0 {
		s.Cursor--
	}
}
func (s *SettingsState) MoveDown() {
	if s.Cursor < len(settingLabels)-1 {
		s.Cursor++
	}
}
func (s *SettingsState) IsBack() bool { return s.Cursor == len(settingLabels)-1 }

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

// ItemsAndValues returns (labels, values, cursor) for the renderer.
func (s *SettingsState) ItemsAndValues() ([]string, []string, int) {
	vals := make([]string, len(settingLabels))
	vals[0] = fmt.Sprintf("◄ %3d%% ►", s.Settings.MusicVolume)
	vals[1] = fmt.Sprintf("◄ %3d%% ►", s.Settings.SoundVolume)
	vals[2] = fmt.Sprintf("◄ %-6s ►", s.Settings.Difficulty.String())
	if s.Settings.ShowFullMap {
		vals[3] = "◄  ON  ►"
	} else {
		vals[3] = "◄  OFF ►"
	}
	vals[4] = ""
	return settingLabels, vals, s.Cursor
}
