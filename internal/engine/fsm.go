package engine

type StateID int

const (
	StateIntro StateID = iota
	StateMainMenu
	StateSettings
	StateCharCreate
	StateWorld
	StateBattle
	StatePause
	StateGameOver
	StateWin
	StateStory
)

func (s StateID) String() string {
	switch s {
	case StateIntro:
		return "intro"
	case StateMainMenu:
		return "main_menu"
	case StateSettings:
		return "settings"
	case StateCharCreate:
		return "char_create"
	case StateWorld:
		return "world"
	case StateBattle:
		return "battle"
	case StatePause:
		return "pause"
	case StateGameOver:
		return "game_over"
	case StateWin:
		return "win"
	case StateStory:
		return "story"
	}
	return "unknown"
}

type FSM struct {
	current StateID
	stack   []StateID
}

func NewFSM() *FSM {
	return &FSM{current: StateIntro}
}

func (f *FSM) Current() StateID {
	return f.current
}

func (f *FSM) Transition(to StateID) {
	f.stack = nil
	f.current = to
}

// Push overlays a state (e.g. pause) keeping prior state on the stack.
func (f *FSM) Push(s StateID) {
	f.stack = append(f.stack, f.current)
	f.current = s
}

// Pop returns to the previous stacked state.
func (f *FSM) Pop() bool {
	if len(f.stack) == 0 {
		return false
	}
	f.current = f.stack[len(f.stack)-1]
	f.stack = f.stack[:len(f.stack)-1]
	return true
}
