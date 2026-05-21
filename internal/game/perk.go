package game

import "math/rand"

type PerkType int

const (
	CombatPerk PerkType = iota
	MagicPerk
	UtilityPerk
	PassivePerk
)

func (pt PerkType) String() string {
	switch pt {
	case CombatPerk:
		return "Combat"
	case MagicPerk:
		return "Magic"
	case UtilityPerk:
		return "Utility"
	case PassivePerk:
		return "Passive"
	}
	return "Unknown"
}

type Perk struct {
	ID       string
	Name     string
	Desc     string
	Type     PerkType
	Rarity   int // 1=common 2=uncommon 3=rare
	ApplyFn  func(c *Character) `json:"-"`
}

var allPerks = []Perk{
	// Combat
	{
		ID: "warrior_blood", Name: "Warrior's Blood", Type: CombatPerk, Rarity: 1,
		Desc:    "+15 Max HP, +3 STR",
		ApplyFn: func(c *Character) { c.MaxHP += 15; c.HP += 15; c.STR += 3 },
	},
	{
		ID: "iron_skin", Name: "Iron Skin", Type: CombatPerk, Rarity: 2,
		Desc:    "+8 DEF, +10 Max HP",
		ApplyFn: func(c *Character) { c.DEF += 8; c.MaxHP += 10; c.HP += 10 },
	},
	{
		ID: "berserker", Name: "Berserker", Type: CombatPerk, Rarity: 3,
		Desc:    "+10 STR, -5 DEF",
		ApplyFn: func(c *Character) { c.STR += 10; c.DEF -= 5 },
	},
	{
		ID: "critical_eye", Name: "Critical Eye", Type: CombatPerk, Rarity: 2,
		Desc:    "+5 AGI (increases crit chance)",
		ApplyFn: func(c *Character) { c.AGI += 5 },
	},
	// Magic
	{
		ID: "arcane_surge", Name: "Arcane Surge", Type: MagicPerk, Rarity: 1,
		Desc:    "+20 Max MP, +3 INT",
		ApplyFn: func(c *Character) { c.MaxMP += 20; c.MP += 20; c.INT += 3 },
	},
	{
		ID: "mana_flow", Name: "Mana Flow", Type: MagicPerk, Rarity: 2,
		Desc:    "+8 INT, spells cost 2 less MP",
		ApplyFn: func(c *Character) { c.INT += 8 },
	},
	{
		ID: "spell_echo", Name: "Spell Echo", Type: MagicPerk, Rarity: 3,
		Desc:    "+12 INT, +15 Max MP",
		ApplyFn: func(c *Character) { c.INT += 12; c.MaxMP += 15; c.MP += 15 },
	},
	// Utility
	{
		ID: "cartographer", Name: "Cartographer", Type: UtilityPerk, Rarity: 1,
		Desc:    "Reveals the full floor map on entry",
		ApplyFn: func(c *Character) {},
	},
	{
		ID: "treasure_sense", Name: "Treasure Sense", Type: UtilityPerk, Rarity: 2,
		Desc:    "+25% gold from all sources",
		ApplyFn: func(c *Character) { c.Gold += c.Gold / 4 },
	},
	{
		ID: "escape_artist", Name: "Escape Artist", Type: UtilityPerk, Rarity: 1,
		Desc:    "Always succeed at fleeing battle",
		ApplyFn: func(c *Character) {},
	},
	// Passive
	{
		ID: "regeneration", Name: "Regeneration", Type: PassivePerk, Rarity: 2,
		Desc:    "Recover 5 HP after each battle",
		ApplyFn: func(c *Character) {},
	},
	{
		ID: "quick_learner", Name: "Quick Learner", Type: PassivePerk, Rarity: 1,
		Desc:    "+20% XP from all battles",
		ApplyFn: func(c *Character) {},
	},
	{
		ID: "last_stand", Name: "Last Stand", Type: PassivePerk, Rarity: 3,
		Desc:    "At low HP, deal 50% more damage",
		ApplyFn: func(c *Character) {},
	},
}

// RandomPerkChoices returns n unique perks the player doesn't already have.
func RandomPerkChoices(c *Character, n int) []Perk {
	owned := map[string]bool{}
	for _, p := range c.Perks {
		owned[p.ID] = true
	}
	pool := make([]Perk, 0, len(allPerks))
	for _, p := range allPerks {
		if !owned[p.ID] {
			pool = append(pool, p)
		}
	}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	if n > len(pool) {
		n = len(pool)
	}
	return pool[:n]
}

func (c *Character) GrantPerk(p Perk) {
	p.ApplyFn(c)
	c.Perks = append(c.Perks, p)
}
