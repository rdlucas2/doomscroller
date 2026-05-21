package game

import (
	"fmt"
	"math/rand"
)

type Enemy struct {
	Name   string
	HP     int
	MaxHP  int
	STR    int
	INT    int
	AGI    int
	DEF    int
	XPVal  int
	GoldVal int
	IsBoss bool
	Skills []EnemySkill
}

type EnemySkill struct {
	Name    string
	MPCost  int
	Damage  int
	Healing int
}

func (e *Enemy) IsAlive() bool { return e.HP > 0 }

func (e *Enemy) StatusLine() string {
	return fmt.Sprintf("%s  HP:%d/%d", e.Name, e.HP, e.MaxHP)
}

// EnemyTable returns a list of enemies scaled to floor and run state.
func RandomEncounterEnemies(floor int, rng *rand.Rand) []*Enemy {
	count := 1 + rng.Intn(3)
	if floor > 3 {
		count = min(count+1, 4)
	}
	enemies := make([]*Enemy, count)
	for i := range enemies {
		enemies[i] = randomEnemy(floor, rng)
	}
	return enemies
}

func BossEnemy(name string, floor int) *Enemy {
	e := &Enemy{
		Name:    name,
		IsBoss:  true,
		XPVal:   200 + floor*100,
		GoldVal: 80 + floor*40,
	}
	scale := 1 + floor
	e.MaxHP = 80 * scale
	e.HP = e.MaxHP
	e.STR = 10 + floor*3
	e.INT = 8 + floor*2
	e.AGI = 6 + floor
	e.DEF = 8 + floor*2
	e.Skills = []EnemySkill{
		{Name: "Power Strike", Damage: 20 + floor*5},
		{Name: "Dark Bolt", Damage: 15 + floor*4},
	}
	return e
}

func randomEnemy(floor int, rng *rand.Rand) *Enemy {
	type template struct {
		name    string
		hpBase  int
		strBase int
		intBase int
		agiBase int
		defBase int
		xpBase  int
		gold    int
	}
	templates := []template{
		{"Goblin", 20, 6, 3, 8, 3, 15, 5},
		{"Skeleton", 25, 8, 2, 5, 6, 18, 6},
		{"Slime", 15, 4, 1, 3, 2, 10, 4},
		{"Shadow Bat", 18, 5, 4, 10, 2, 14, 5},
		{"Stone Golem", 40, 10, 2, 3, 12, 30, 12},
		{"Cursed Knight", 35, 12, 5, 6, 10, 35, 15},
		{"Void Wraith", 28, 7, 12, 7, 4, 32, 14},
		{"Tower Mage", 22, 5, 14, 6, 3, 38, 16},
	}
	t := templates[rng.Intn(len(templates))]
	scale := 1 + (floor-1)/2
	return &Enemy{
		Name:    t.name,
		MaxHP:   t.hpBase * scale,
		HP:      t.hpBase * scale,
		STR:     t.strBase + floor,
		INT:     t.intBase + floor,
		AGI:     t.agiBase + floor/2,
		DEF:     t.defBase + floor/2,
		XPVal:   t.xpBase + floor*5,
		GoldVal: t.gold + floor*2,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
