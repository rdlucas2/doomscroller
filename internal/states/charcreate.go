package states

import (
	"strings"
	"unicode"

	"github.com/rdlucas2/doomscroller/internal/game"
)

type CharCreatePhase int

const (
	PhasePickName CharCreatePhase = iota
	PhasePickClass
	PhaseConfirm
)

type CharCreateState struct {
	Phase         CharCreatePhase
	NameInput     string
	ClassCursor   int
	ConfirmCursor int
}

var classOptions = []game.Class{game.Warrior, game.Mage, game.Rogue}

func (s *CharCreateState) CurrentClass() game.Class { return classOptions[s.ClassCursor] }

func (s *CharCreateState) HandleRune(r rune) {
	if len([]rune(s.NameInput)) >= 14 {
		return
	}
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == ' ' {
		s.NameInput += string(r)
	}
}

func (s *CharCreateState) HandleBackspace() {
	runes := []rune(s.NameInput)
	if len(runes) > 0 {
		s.NameInput = string(runes[:len(runes)-1])
	}
}

func (s *CharCreateState) IsNameReady() bool { return strings.TrimSpace(s.NameInput) != "" }

func (s *CharCreateState) MoveClassLeft() {
	if s.ClassCursor > 0 {
		s.ClassCursor--
	}
}
func (s *CharCreateState) MoveClassRight() {
	if s.ClassCursor < len(classOptions)-1 {
		s.ClassCursor++
	}
}

func (s *CharCreateState) BuildCharacter() *game.Character {
	name := strings.TrimSpace(s.NameInput)
	if name == "" {
		name = "Hero"
	}
	return game.NewCharacter(name, s.CurrentClass())
}
