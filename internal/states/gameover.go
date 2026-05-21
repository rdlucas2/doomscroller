package states

import "github.com/rdlucas2/doomscroller/internal/game"

var gameOverItems = []string{"New Game", "Main Menu", "Quit"}

type GameOverState struct {
	Character *game.Character
	Floor     int
	Cursor    int
}

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
func (s *GameOverState) Selected() string { return gameOverItems[s.Cursor] }
