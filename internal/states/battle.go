package states

import "github.com/rdlucas2/doomscroller/internal/game"

type BattleUIState struct {
	Battle           *game.Battle
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
		if len(b.LiveEnemies()) > 0 && b.SelectedTarget > 0 {
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
	case "left", "h", "a":
		if b.SelectedPerk > 0 {
			b.SelectedPerk--
		}
	case "right", "l", "d":
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
		if b.ExecutePlayerMagic(b.SelectedTarget) {
			bs.afterPlayerAction()
		} else {
			b.Phase = game.PhasePlayerMenu
		}
	case game.ActionItem:
		b.AddLog("No items available.")
	case game.ActionFlee:
		b.Phase = game.PhasePlayerAction
		if !b.ExecutePlayerFlee() {
			bs.afterPlayerAction()
		}
	}
}

func (bs *BattleUIState) afterPlayerAction() {
	b := bs.Battle
	if b.CheckVictory() {
		b.CollectRewards()
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
	b.SelectedTarget = 0
	b.Phase = game.PhasePlayerMenu
}
