package stories

import "github.com/rdlucas2/doomscroller/internal/game"

func init() {
	game.RegisterStory(&game.StoryBeat{
		ID:       "prologue_village",
		Name:     "The Burning Village",
		FloorMin: 1, FloorMax: 1,
		Trigger: game.TriggerEnterRoom,
		Dialogues: []game.Dialogue{
			{Speaker: "Elder Maris", Text: "You there! Adventurer! The dark tower erupted three nights ago."},
			{Speaker: "Elder Maris", Text: "Monsters pour from its depths. Our village cannot survive another attack."},
			{Speaker: "Hero", Text: "I'll enter the tower and put an end to it."},
			{Speaker: "Elder Maris", Text: "The Void Sovereign waits at the summit. May the ancients guide you."},
		},
		Reward: game.Reward{Gold: 30, XP: 50},
	})

	game.RegisterStory(&game.StoryBeat{
		ID:       "mid_revelation",
		Name:     "The Forgotten Archive",
		FloorMin: 3, FloorMax: 4,
		Trigger: game.TriggerEnterRoom,
		Dialogues: []game.Dialogue{
			{Speaker: "Spirit", Text: "So you found this place. The archive holds the tower's darkest secret."},
			{Speaker: "Spirit", Text: "The Void Sovereign was once a hero like you—consumed by power."},
			{Speaker: "Hero", Text: "Then there must be a way to reach what remains of their humanity."},
			{Speaker: "Spirit", Text: "The Sovereign's Heart is shattered across three seals. Destroy them."},
		},
		Reward: game.Reward{Gold: 60, XP: 120, ItemName: "Archive Key"},
	})

	game.RegisterStory(&game.StoryBeat{
		ID:       "seal_guardian",
		Name:     "Guardian of the First Seal",
		FloorMin: 2, FloorMax: 3,
		Trigger: game.TriggerBossDefeated,
		BossName: "Stone Colossus",
		Dialogues: []game.Dialogue{
			{Speaker: "Colossus", Text: "INTRUDER... THE SEAL... MUST HOLD..."},
			{Speaker: "Hero", Text: "I'm sorry. The seal must be broken."},
			{Speaker: "Colossus", Text: "...Then prove... your worth..."},
		},
		Reward: game.Reward{Gold: 80, XP: 200, ItemName: "First Seal Fragment"},
	})

	game.RegisterStory(&game.StoryBeat{
		ID:       "final_confrontation",
		Name:     "Summit of the Void",
		FloorMin: 5, FloorMax: 5,
		Trigger: game.TriggerEnterRoom,
		Dialogues: []game.Dialogue{
			{Speaker: "Void Sovereign", Text: "You dare ascend my tower? No one has reached this height in centuries."},
			{Speaker: "Hero", Text: "It ends here. The village—the people—they deserve peace."},
			{Speaker: "Void Sovereign", Text: "Peace? I sought peace once. Power was the answer I found."},
			{Speaker: "Hero", Text: "Then let me show you a different answer."},
		},
		BossName: "Void Sovereign",
		Reward:   game.Reward{Gold: 500, XP: 1000, ItemName: "Sovereign's Crown"},
	})
}
