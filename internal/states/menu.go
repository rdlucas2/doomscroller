package states

import "github.com/rdlucas2/doomscroller/internal/game"

type MenuChoice int

const (
	MenuNewGame MenuChoice = iota
	MenuContinue
	MenuSettings
	MenuQuit
)

type MenuState struct {
	Cursor  int
	hasSave bool
}

func NewMenuState() *MenuState {
	return &MenuState{hasSave: game.HasSave()}
}

func (s *MenuState) Items() []string {
	items := []string{"New Game"}
	if s.hasSave {
		items = append(items, "Continue")
	}
	return append(items, "Settings", "Quit")
}

func (s *MenuState) Selected() MenuChoice {
	item := s.Items()[s.Cursor]
	switch item {
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
	if s.Cursor < len(s.Items())-1 {
		s.Cursor++
	}
}
