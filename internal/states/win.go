package states

import "github.com/rdlucas2/doomscroller/internal/game"

var winItems = []string{"New Game", "Main Menu", "Quit"}

type WinState struct {
	Character *game.Character
	Cursor    int
}

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
func (s *WinState) Selected() string { return winItems[s.Cursor] }
