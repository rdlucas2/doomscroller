package game

import "fmt"

type Class int

const (
	Warrior Class = iota
	Mage
	Rogue
)

func (c Class) String() string {
	switch c {
	case Warrior:
		return "Warrior"
	case Mage:
		return "Mage"
	case Rogue:
		return "Rogue"
	}
	return "Unknown"
}

func (c Class) Description() string {
	switch c {
	case Warrior:
		return "High HP & STR. Excels in physical combat."
	case Mage:
		return "High MP & INT. Commands devastating spells."
	case Rogue:
		return "High AGI. Fast, evasive, deadly criticals."
	}
	return ""
}

type Character struct {
	Name      string
	Class     Class
	Level     int
	XP        int
	XPToNext  int
	HP        int
	MaxHP     int
	MP        int
	MaxMP     int
	STR       int
	INT       int
	AGI       int
	DEF       int
	Gold      int
	Perks     []Perk
	Equipment Equipment
}

type Equipment struct {
	Weapon string
	Armor  string
}

func NewCharacter(name string, class Class) *Character {
	c := &Character{
		Name:     name,
		Class:    class,
		Level:    1,
		XP:       0,
		XPToNext: 100,
		Gold:     50,
	}
	switch class {
	case Warrior:
		c.MaxHP, c.HP = 120, 120
		c.MaxMP, c.MP = 20, 20
		c.STR, c.INT, c.AGI, c.DEF = 14, 6, 8, 12
		c.Equipment = Equipment{Weapon: "Iron Sword", Armor: "Chain Mail"}
	case Mage:
		c.MaxHP, c.HP = 70, 70
		c.MaxMP, c.MP = 80, 80
		c.STR, c.INT, c.AGI, c.DEF = 6, 16, 10, 5
		c.Equipment = Equipment{Weapon: "Oak Staff", Armor: "Robe"}
	case Rogue:
		c.MaxHP, c.HP = 90, 90
		c.MaxMP, c.MP = 40, 40
		c.STR, c.INT, c.AGI, c.DEF = 10, 8, 16, 7
		c.Equipment = Equipment{Weapon: "Twin Daggers", Armor: "Leather Vest"}
	}
	return c
}

func (c *Character) AddXP(amount int) (levelUp bool) {
	c.XP += amount
	if c.XP >= c.XPToNext {
		c.XP -= c.XPToNext
		c.Level++
		c.XPToNext = 100 + c.Level*50
		c.STR += 2
		c.INT += 2
		c.AGI += 1
		c.DEF += 1
		hpGain := 10 + c.Level
		c.MaxHP += hpGain
		c.HP += hpGain
		mpGain := 5 + c.Level
		c.MaxMP += mpGain
		c.MP += mpGain
		return true
	}
	return false
}

func (c *Character) PhysicalDamage() int {
	base := c.STR + 5
	switch c.Class {
	case Warrior:
		base += 4
	case Rogue:
		base += 2
	}
	return base
}

func (c *Character) MagicDamage() int {
	base := c.INT + 8
	if c.Class == Mage {
		base += 6
	}
	return base
}

func (c *Character) HasPerk(id string) bool {
	for _, p := range c.Perks {
		if p.ID == id {
			return true
		}
	}
	return false
}

func (c *Character) StatusLine() string {
	return fmt.Sprintf("HP:%d/%d  MP:%d/%d  LV:%d  XP:%d/%d",
		c.HP, c.MaxHP, c.MP, c.MaxMP, c.Level, c.XP, c.XPToNext)
}
