package stories

import "github.com/rdlucas2/doomscroller/internal/game"

// WeeklyStories are optional story beats that can be toggled in/out for weekly/monthly content drops.
// To add a new weekly story: create a new StoryBeat and call game.RegisterStory in an init() or
// a dedicated RegisterWeekly() function called from main.

func init() {
	game.RegisterStory(&game.StoryBeat{
		ID:       "weekly_lost_merchant",
		Name:     "The Lost Merchant",
		FloorMin: 1, FloorMax: 3,
		Trigger: game.TriggerEnterRoom,
		Dialogues: []game.Dialogue{
			{Speaker: "Merchant", Text: "Oh thank the stars! A living soul! I've been trapped here a week!"},
			{Speaker: "Hero", Text: "How did you end up inside the tower?"},
			{Speaker: "Merchant", Text: "Chasing a debt. Not my finest hour. Here, take these supplies as thanks."},
		},
		Reward: game.Reward{Gold: 120, XP: 75, ItemName: "Merchant's Pack"},
	})

	game.RegisterStory(&game.StoryBeat{
		ID:       "monthly_shadow_twins",
		Name:     "The Shadow Twins",
		FloorMin: 2, FloorMax: 4,
		Trigger: game.TriggerBossDefeated,
		BossName: "Shadow Twins",
		Dialogues: []game.Dialogue{
			{Speaker: "Shadow A", Text: "We were born in this tower. It is our home."},
			{Speaker: "Shadow B", Text: "And you will NOT take it from us."},
			{Speaker: "Hero", Text: "The tower is corrupted. Come with me—there's a world outside these walls."},
			{Speaker: "Shadow A", Text: "...Outside. We have... never seen outside."},
		},
		Reward: game.Reward{Gold: 200, XP: 350, ItemName: "Twin Shadow Essence"},
	})
}
