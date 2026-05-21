package engine

import "github.com/hajimehoshi/ebiten/v2"

var watchKeys = []ebiten.Key{
	ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
	ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
	ebiten.KeyH, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL,
	ebiten.KeyEnter, ebiten.KeySpace,
	ebiten.KeyEscape, ebiten.KeyP,
	ebiten.KeyBackspace,
}

// keyMap maps Ebiten keys to short string names (same strings as the old bubbletea model).
var keyMap = map[ebiten.Key]string{
	ebiten.KeyArrowUp:    "up",
	ebiten.KeyArrowDown:  "down",
	ebiten.KeyArrowLeft:  "left",
	ebiten.KeyArrowRight: "right",
	ebiten.KeyW: "w", ebiten.KeyA: "a", ebiten.KeyS: "s", ebiten.KeyD: "d",
	ebiten.KeyH: "h", ebiten.KeyJ: "j", ebiten.KeyK: "k", ebiten.KeyL: "l",
	ebiten.KeyEnter:     "enter",
	ebiten.KeySpace:     " ",
	ebiten.KeyEscape:    "esc",
	ebiten.KeyP:         "p",
	ebiten.KeyBackspace: "backspace",
}

type InputState struct {
	prev map[ebiten.Key]bool
	curr map[ebiten.Key]bool
}

func NewInputState() *InputState {
	return &InputState{
		prev: make(map[ebiten.Key]bool),
		curr: make(map[ebiten.Key]bool),
	}
}

func (s *InputState) Update() {
	s.prev = s.curr
	s.curr = make(map[ebiten.Key]bool, len(watchKeys))
	for _, k := range watchKeys {
		s.curr[k] = ebiten.IsKeyPressed(k)
	}
}

// JustPressed returns true only on the first frame the key is held.
func (s *InputState) JustPressed(k ebiten.Key) bool {
	return s.curr[k] && !s.prev[k]
}

// JustPressedKeys returns string names of all just-pressed keys.
func (s *InputState) JustPressedKeys() []string {
	var keys []string
	for k, name := range keyMap {
		if s.JustPressed(k) {
			keys = append(keys, name)
		}
	}
	return keys
}

// PressedRunes returns rune characters typed this frame (for text input).
func PressedRunes() []rune {
	return ebiten.AppendInputChars(nil)
}
