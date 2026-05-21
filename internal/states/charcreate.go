package states

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type CharCreatePhase int

const (
	PhasePickName CharCreatePhase = iota
	PhasePickClass
	PhaseConfirm
)

type CharCreateState struct {
	Phase      CharCreatePhase
	NameInput  string
	ClassCursor int
	ConfirmCursor int
}

var classOptions = []game.Class{game.Warrior, game.Mage, game.Rogue}

func (s *CharCreateState) CurrentClass() game.Class {
	return classOptions[s.ClassCursor]
}

func (s *CharCreateState) HandleNameKey(key string) {
	switch key {
	case "backspace":
		if len(s.NameInput) > 0 {
			runes := []rune(s.NameInput)
			s.NameInput = string(runes[:len(runes)-1])
		}
	default:
		if len(key) == 1 && len([]rune(s.NameInput)) < 14 {
			r := rune(key[0])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == ' ' {
				s.NameInput += key
			}
		}
	}
}

func (s *CharCreateState) IsNameReady() bool {
	return strings.TrimSpace(s.NameInput) != ""
}

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

func (s *CharCreateState) View(width, height int) string {
	lines := make([]string, 0, height)

	topPad := (height - 24) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	title := ui.TitleStyle.Render("⚔  CHARACTER CREATION")
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(title))
	lines = append(lines, "")

	switch s.Phase {
	case PhasePickName:
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.SubtitleStyle.Render("Enter your hero's name:"),
		))
		lines = append(lines, "")
		cursor := "█"
		nameDisplay := fmt.Sprintf("[ %s%s ]", s.NameInput, cursor)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.TitleStyle.Render(nameDisplay),
		))
		lines = append(lines, "")
		lines = append(lines, "")
		if s.IsNameReady() {
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
				ui.SuccessStyle.Render("Press Enter to continue"),
			))
		} else {
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
				ui.DimStyle.Render("Type a name and press Enter"),
			))
		}

	case PhasePickClass:
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.SubtitleStyle.Render("Choose your class:"),
		))
		lines = append(lines, "")
		// Class cards
		cards := make([]string, len(classOptions))
		for i, cls := range classOptions {
			dummy := game.NewCharacter("x", cls)
			statLine := fmt.Sprintf("HP:%-3d  MP:%-3d\nSTR:%-3d INT:%-3d\nAGI:%-3d DEF:%-3d",
				dummy.MaxHP, dummy.MaxMP,
				dummy.STR, dummy.INT,
				dummy.AGI, dummy.DEF)
			desc := cls.Description()
			body := fmt.Sprintf("%s\n\n%s\n\n%s", ui.TitleStyle.Render(cls.String()), desc, statLine)
			if i == s.ClassCursor {
				cards[i] = ui.HighlightBox.Width(24).Render(body)
			} else {
				cards[i] = ui.BoxStyle.Width(24).Render(body)
			}
		}
		cardRow := lipgloss.JoinHorizontal(lipgloss.Top, cards[0], "  ", cards[1], "  ", cards[2])
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(cardRow))
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("◄► Change Class   Enter Confirm"),
		))

	case PhaseConfirm:
		c := s.BuildCharacter()
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.SubtitleStyle.Render("Your hero is ready!"),
		))
		lines = append(lines, "")
		summary := fmt.Sprintf(
			"%s  the  %s\n\nHP: %d   MP: %d\nSTR: %d  INT: %d  AGI: %d  DEF: %d\n\nWeapon: %s\nArmor:  %s",
			ui.TitleStyle.Render(c.Name),
			ui.GoldStyle.Render(c.Class.String()),
			c.MaxHP, c.MaxMP,
			c.STR, c.INT, c.AGI, c.DEF,
			c.Equipment.Weapon,
			c.Equipment.Armor,
		)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.HighlightBox.Width(40).Render(summary),
		))
		lines = append(lines, "")

		confirmItems := []string{"Begin Adventure", "Back"}
		for i, item := range confirmItems {
			var rendered string
			if i == s.ConfirmCursor {
				rendered = ui.SelectedStyle.Render("▶  " + item)
			} else {
				rendered = ui.NormalStyle.Render("   " + item)
			}
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(rendered))
		}
	}

	return strings.Join(lines, "\n")
}
