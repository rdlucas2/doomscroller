package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type StoryState struct {
	Beat    *game.StoryBeat
	Line    int
	Done    bool
}

func NewStoryState(beat *game.StoryBeat) *StoryState {
	return &StoryState{Beat: beat}
}

func (ss *StoryState) Advance() {
	if ss.Line < len(ss.Beat.Dialogues)-1 {
		ss.Line++
	} else {
		ss.Done = true
		game.MarkCompleted(ss.Beat.ID)
	}
}

func (ss *StoryState) View(width, height int) string {
	beat := ss.Beat
	lines := make([]string, 0, height)

	topPad := (height - 14) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	// Story title
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
		ui.TitleStyle.Render(fmt.Sprintf("✦  %s", beat.Name)),
	))
	lines = append(lines, "")

	if len(beat.Dialogues) == 0 {
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("(No dialogue)"),
		))
	} else {
		d := beat.Dialogues[ss.Line]
		speaker := ui.SubtitleStyle.Render(d.Speaker + ":")
		text := wrapText(d.Text, 60)
		dialogue := fmt.Sprintf("%s\n\n%s", speaker, text)
		box := ui.HighlightBox.Width(66).Render(dialogue)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(box))
		lines = append(lines, "")

		// Progress indicator
		prog := fmt.Sprintf("(%d / %d)", ss.Line+1, len(beat.Dialogues))
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render(prog),
		))
	}

	lines = append(lines, "")

	if ss.Line == len(beat.Dialogues)-1 {
		// Show reward
		r := beat.Reward
		var rewardParts []string
		if r.Gold > 0 {
			rewardParts = append(rewardParts, ui.GoldStyle.Render(fmt.Sprintf("+%d Gold", r.Gold)))
		}
		if r.XP > 0 {
			rewardParts = append(rewardParts, ui.XPStyle.Render(fmt.Sprintf("+%d XP", r.XP)))
		}
		if r.ItemName != "" {
			rewardParts = append(rewardParts, ui.SuccessStyle.Render("Item: "+r.ItemName))
		}
		if len(rewardParts) > 0 {
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
				"Reward: "+strings.Join(rewardParts, "  "),
			))
			lines = append(lines, "")
		}
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.SuccessStyle.Render("Press Enter to continue"),
		))
	} else {
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("Press Enter or Space to advance"),
		))
	}

	return strings.Join(lines, "\n")
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	var lines []string
	line := ""
	for _, w := range words {
		if len(line)+len(w)+1 > width {
			lines = append(lines, line)
			line = w
		} else {
			if line != "" {
				line += " "
			}
			line += w
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
