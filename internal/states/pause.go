package states

type PauseChoice int

const (
	PauseResume PauseChoice = iota
	PauseSave
	PauseSettings
	PauseQuitToMenu
)

type PauseState struct{ Cursor int }

func (ps *PauseState) MoveUp() {
	if ps.Cursor > 0 {
		ps.Cursor--
	}
}
func (ps *PauseState) MoveDown() {
	if ps.Cursor < 3 {
		ps.Cursor++
	}
}
