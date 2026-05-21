package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type BattleUIState struct {
	Battle          *game.Battle
	StoryAfterBattle *game.StoryBeat
}

func NewBattleUIState(b *game.Battle, story *game.StoryBeat) *BattleUIState {
	return &BattleUIState{Battle: b, StoryAfterBattle: story}
}

func (bs *BattleUIState) HandleInput(key string) {
	b := bs.Battle
	if b.Phase != game.PhasePlayerMenu {
		return
	}
	switch key {
	case "up", "k", "w":
		if b.SelectedAction > 0 {
			b.SelectedAction--
		}
	case "down", "j", "s":
		if b.SelectedAction < len(game.ActionNames)-1 {
			b.SelectedAction++
		}
	case "left", "h", "a":
		lives := b.LiveEnemies()
		if len(lives) > 0 && b.SelectedTarget > 0 {
			b.SelectedTarget--
		}
	case "right", "l", "d":
		lives := b.LiveEnemies()
		if len(lives) > 0 && b.SelectedTarget < len(lives)-1 {
			b.SelectedTarget++
		}
	case "enter", " ":
		bs.executeAction()
	}
}

func (bs *BattleUIState) HandlePerkChoice(key string) {
	b := bs.Battle
	if b.Phase != game.PhasePerkChoice {
		return
	}
	switch key {
	case "up", "k", "w":
		if b.SelectedPerk > 0 {
			b.SelectedPerk--
		}
	case "down", "j", "s":
		if b.SelectedPerk < len(b.PerkChoices)-1 {
			b.SelectedPerk++
		}
	case "enter", " ":
		if b.SelectedPerk < len(b.PerkChoices) {
			b.Player.GrantPerk(b.PerkChoices[b.SelectedPerk])
		}
		b.Phase = game.PhaseVictory
	}
}

func (bs *BattleUIState) executeAction() {
	b := bs.Battle
	switch game.BattleAction(b.SelectedAction) {
	case game.ActionAttack:
		b.Phase = game.PhasePlayerAction
		b.ExecutePlayerAttack(b.SelectedTarget)
		bs.afterPlayerAction()
	case game.ActionMagic:
		b.Phase = game.PhasePlayerAction
		ok := b.ExecutePlayerMagic(b.SelectedTarget)
		if ok {
			bs.afterPlayerAction()
		} else {
			b.Phase = game.PhasePlayerMenu
		}
	case game.ActionItem:
		b.AddLog("No items available.")
	case game.ActionFlee:
		b.Phase = game.PhasePlayerAction
		fled := b.ExecutePlayerFlee()
		if !fled {
			bs.afterPlayerAction()
		}
	}
}

func (bs *BattleUIState) afterPlayerAction() {
	b := bs.Battle
	if b.CheckVictory() {
		b.CollectRewards()
		// Offer a perk if eligible (every 2nd win)
		choices := game.RandomPerkChoices(b.Player, 3)
		if len(choices) > 0 {
			b.PerkChoices = choices
			b.Phase = game.PhasePerkChoice
		} else {
			b.Phase = game.PhaseVictory
		}
		return
	}
	if b.CheckDefeat() {
		b.Phase = game.PhaseDefeat
		return
	}
	b.Phase = game.PhaseEnemyAction
	b.ExecuteEnemyTurn()
	if b.CheckDefeat() {
		b.Phase = game.PhaseDefeat
		return
	}
	// Clean up dead enemies
	b.SelectedTarget = 0
	b.Phase = game.PhasePlayerMenu
}

func (bs *BattleUIState) View(width, height int) string {
	b := bs.Battle
	lines := make([]string, 0, height)

	// Title bar
	titleStr := "⚔  BATTLE"
	if b.IsBossEncounter {
		titleStr = "💀  BOSS BATTLE"
	}
	lines = append(lines,
		lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(ui.TitleStyle.Render(titleStr)),
	)
	lines = append(lines, "")

	// --- Enemies ---
	liveEnemies := b.LiveEnemies()
	enemyCards := make([]string, len(b.Enemies))
	for i, e := range b.Enemies {
		var body string
		if !e.IsAlive() {
			body = ui.DimStyle.Render(fmt.Sprintf("%s\n[Defeated]", e.Name))
			enemyCards[i] = ui.BoxStyle.Width(16).Render(body)
		} else {
			hpBar := ui.HPBar(e.HP, e.MaxHP, 10)
			body = fmt.Sprintf("%s\n%s\n%s",
				ui.HPStyle.Render(e.Name),
				hpBar,
				fmt.Sprintf("HP:%d/%d", e.HP, e.MaxHP))
			// Highlight selected target
			targetIdx := -1
			for ti, le := range liveEnemies {
				if le == e {
					targetIdx = ti
					break
				}
			}
			if targetIdx == b.SelectedTarget && b.Phase == game.PhasePlayerMenu {
				enemyCards[i] = ui.HighlightBox.Width(16).Render(body)
			} else {
				enemyCards[i] = ui.BoxStyle.Width(16).Render(body)
			}
		}
	}
	enemyRow := lipgloss.JoinHorizontal(lipgloss.Top, enemyCards...)
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(enemyRow))
	lines = append(lines, "")

	// --- Player status ---
	c := b.Player
	playerStatus := fmt.Sprintf("%s  [%s Lv.%d]  |  HP %s %d/%d  |  MP %s %d/%d",
		ui.TitleStyle.Render(c.Name),
		ui.GoldStyle.Render(c.Class.String()),
		c.Level,
		ui.HPBar(c.HP, c.MaxHP, 10),
		c.HP, c.MaxHP,
		ui.MPBar(c.MP, c.MaxMP, 8),
		c.MP, c.MaxMP,
	)
	lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(playerStatus))
	lines = append(lines, "")

	// --- Battle log ---
	logContent := strings.Join(b.Log, "\n")
	logBox := ui.BattleLogStyle.Width(width - 4).Render(logContent)
	lines = append(lines, logBox)
	lines = append(lines, "")

	// --- Action menu / perk choice / outcome ---
	switch b.Phase {
	case game.PhasePlayerMenu:
		actionItems := make([]string, len(game.ActionNames))
		for i, name := range game.ActionNames {
			if i == b.SelectedAction {
				actionItems[i] = ui.SelectedStyle.Render("▶ " + name)
			} else {
				actionItems[i] = ui.NormalStyle.Render("  " + name)
			}
		}
		actionRow := lipgloss.JoinHorizontal(lipgloss.Top,
			actionItems[0]+"   ", actionItems[1]+"   ",
			actionItems[2]+"   ", actionItems[3])
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.BoxStyle.Render(actionRow),
		))
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("↑↓ Action   ◄► Target   Enter Confirm"),
		))

	case game.PhasePerkChoice:
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.TitleStyle.Render("Choose a Perk!"),
		))
		lines = append(lines, "")
		perkCards := make([]string, len(b.PerkChoices))
		for i, p := range b.PerkChoices {
			body := fmt.Sprintf("%s\n[%s]\n\n%s",
				perkTypeStyle(p.Type).Render(p.Name),
				ui.DimStyle.Render(p.Type.String()),
				p.Desc)
			if i == b.SelectedPerk {
				perkCards[i] = ui.HighlightBox.Width(20).Render(body)
			} else {
				perkCards[i] = ui.BoxStyle.Width(20).Render(body)
			}
		}
		perkRow := lipgloss.JoinHorizontal(lipgloss.Top, perkCards...)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(perkRow))
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("↑↓ Navigate   Enter Select"),
		))

	case game.PhaseVictory:
		msg := fmt.Sprintf("Victory!  +%d XP  +%d Gold", b.TotalXP, b.TotalGold)
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.SuccessStyle.Render(msg),
		))
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("Press Enter to continue"),
		))

	case game.PhaseDefeat:
		if b.Player.HP <= 0 {
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
				ui.ErrorStyle.Render("You have been defeated..."),
			))
		} else {
			lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
				ui.DimStyle.Render("Escaped from battle."),
			))
		}
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			ui.DimStyle.Render("Press Enter to continue"),
		))
	}

	return strings.Join(lines, "\n")
}
