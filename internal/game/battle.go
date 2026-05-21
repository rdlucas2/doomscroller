package game

import (
	"fmt"
	"math/rand"
)

type BattlePhase int

const (
	PhasePlayerMenu BattlePhase = iota
	PhasePlayerAction
	PhaseEnemyAction
	PhaseVictory
	PhaseDefeat
	PhasePerkChoice
)

type BattleAction int

const (
	ActionAttack BattleAction = iota
	ActionMagic
	ActionItem
	ActionFlee
)

var ActionNames = []string{"Attack", "Magic", "Item", "Flee"}

type Battle struct {
	Player         *Character
	Enemies        []*Enemy
	Phase          BattlePhase
	Log            []string
	SelectedAction int
	SelectedTarget int
	SelectedPerk   int
	PerkChoices    []Perk
	TotalXP        int
	TotalGold      int
	IsBossEncounter bool
	Rng            *rand.Rand
}

func NewBattle(player *Character, enemies []*Enemy, rng *rand.Rand) *Battle {
	return &Battle{
		Player:  player,
		Enemies: enemies,
		Phase:   PhasePlayerMenu,
		Rng:     rng,
		Log:     []string{"A battle begins!"},
	}
}

func (b *Battle) LiveEnemies() []*Enemy {
	var out []*Enemy
	for _, e := range b.Enemies {
		if e.IsAlive() {
			out = append(out, e)
		}
	}
	return out
}

func (b *Battle) ExecutePlayerAttack(targetIdx int) {
	enemies := b.LiveEnemies()
	if targetIdx >= len(enemies) {
		return
	}
	target := enemies[targetIdx]
	dmg := b.Player.PhysicalDamage() + b.Rng.Intn(8) - b.Rng.Intn(4)
	dmg -= target.DEF / 2
	// Last Stand perk
	if b.Player.HasPerk("last_stand") && b.Player.HP <= b.Player.MaxHP/4 {
		dmg = dmg * 3 / 2
	}
	// Crit check based on AGI
	crit := b.Rng.Intn(100) < b.Player.AGI
	if crit {
		dmg = dmg * 2
	}
	if dmg < 1 {
		dmg = 1
	}
	target.HP -= dmg
	msg := fmt.Sprintf("%s attacks %s for %d damage!", b.Player.Name, target.Name, dmg)
	if crit {
		msg = "CRITICAL! " + msg
	}
	b.AddLog(msg)
}

func (b *Battle) ExecutePlayerMagic(targetIdx int) bool {
	enemies := b.LiveEnemies()
	if targetIdx >= len(enemies) {
		return false
	}
	target := enemies[targetIdx]
	mpCost := 12
	// mana_flow perk reduces cost
	if b.Player.HasPerk("mana_flow") {
		mpCost -= 2
	}
	if b.Player.MP < mpCost {
		b.AddLog("Not enough MP!")
		return false
	}
	b.Player.MP -= mpCost
	dmg := b.Player.MagicDamage() + b.Rng.Intn(10)
	target.HP -= dmg
	b.AddLog(fmt.Sprintf("%s casts a spell on %s for %d magic damage!", b.Player.Name, target.Name, dmg))
	return true
}

func (b *Battle) ExecutePlayerFlee() bool {
	if b.Player.HasPerk("escape_artist") {
		b.AddLog("You slip away effortlessly!")
		b.Phase = PhaseDefeat
		return true
	}
	chance := 40 + b.Player.AGI*2
	if b.Rng.Intn(100) < chance {
		b.AddLog("You successfully fled!")
		b.Phase = PhaseDefeat
		return true
	}
	b.AddLog("Couldn't escape!")
	return false
}

func (b *Battle) ExecuteEnemyTurn() {
	for _, e := range b.LiveEnemies() {
		if !e.IsAlive() {
			continue
		}
		action := b.Rng.Intn(3)
		var dmg int
		var msg string
		if len(e.Skills) > 0 && action == 2 {
			skill := e.Skills[b.Rng.Intn(len(e.Skills))]
			dmg = skill.Damage + b.Rng.Intn(6)
			dmg -= b.Player.DEF / 2
			if dmg < 1 {
				dmg = 1
			}
			msg = fmt.Sprintf("%s uses %s on you for %d damage!", e.Name, skill.Name, dmg)
		} else {
			dmg = e.STR + b.Rng.Intn(6)
			dmg -= b.Player.DEF / 2
			if dmg < 1 {
				dmg = 1
			}
			msg = fmt.Sprintf("%s attacks you for %d damage!", e.Name, dmg)
		}
		b.Player.HP -= dmg
		b.AddLog(msg)
	}
}

func (b *Battle) CheckVictory() bool {
	for _, e := range b.Enemies {
		if e.IsAlive() {
			return false
		}
	}
	return true
}

func (b *Battle) CheckDefeat() bool {
	return b.Player.HP <= 0
}

func (b *Battle) CollectRewards() {
	xpMult := 1
	if b.Player.HasPerk("quick_learner") {
		xpMult = 6 // 20% more via integer: 5+1/5 → use 120%
	}
	for _, e := range b.Enemies {
		b.TotalXP += e.XPVal
		b.TotalGold += e.GoldVal
	}
	if b.Player.HasPerk("quick_learner") {
		b.TotalXP = b.TotalXP * 6 / 5
	}
	_ = xpMult
	if b.Player.HasPerk("treasure_sense") {
		b.TotalGold = b.TotalGold * 5 / 4
	}
	b.Player.AddXP(b.TotalXP)
	b.Player.Gold += b.TotalGold
	// Regeneration perk
	if b.Player.HasPerk("regeneration") {
		heal := 5
		b.Player.HP = min(b.Player.HP+heal, b.Player.MaxHP)
	}
}

func (b *Battle) AddLog(msg string) {
	b.Log = append(b.Log, msg)
	if len(b.Log) > 8 {
		b.Log = b.Log[len(b.Log)-8:]
	}
}
