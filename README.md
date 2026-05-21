# Doom Scroller — Tower of the Void

A terminal JRPG roguelike built in Go with a game loop and finite state machine.

## Features

- **Intro / Loading screen** → animated logo with progress bar
- **Main Menu** — New Game, Continue, Settings, Quit
- **Character Creation** — Name input + class selection (Warrior / Mage / Rogue)
- **Procedurally generated dungeon maps** — BSP room placement, fog of war
- **5-floor tower** with a boss on each floor and escalating difficulty
- **Story beats system** — modular JRPG story scenes injected into procedural maps
  - Core story: prologue → mid-revelation → seal guardians → final boss
  - Weekly / monthly beats: add new `StoryBeat` structs to `content/stories/` and they appear automatically in the relevant floor range
- **Random + triggered encounters** — move around the map to trigger battles
- **Turn-based battle system** — Attack / Magic / Item / Flee, enemy AI, critical hits
- **Roguelike perk system** — 4 perk types (Combat, Magic, Utility, Passive), choose 1-of-3 after each victory
- **Pause menu** — Resume, Save & Quit, Settings, Quit to Menu
- **Game Over** and **Victory** screens
- **Save / Load** — JSON save file in `~/.doomscroller/save.json`

## Running

```bash
go run ./cmd/doomscroller
```

Or build first:

```bash
go build -o doomscroller ./cmd/doomscroller
./doomscroller
```

## Controls

| Context | Key | Action |
|---------|-----|--------|
| All menus | ↑↓ / `j``k` | Navigate |
| All menus | Enter / Space | Select |
| Character creation | ◄► / `h``l` | Cycle class |
| World map | WASD / Arrows | Move player |
| World map | `P` / Esc | Pause |
| Battle | ↑↓ | Choose action |
| Battle | ◄► | Choose target |
| Battle | Enter | Confirm |
| Settings | ◄► | Adjust value |
| Story | Enter / Space | Advance dialogue |
| Story | Esc | Skip scene |

## Adding New Story Beats

Create a new file in `content/stories/` or add to `weekly.go`:

```go
func init() {
    game.RegisterStory(&game.StoryBeat{
        ID:       "my_story_week_3",
        Name:     "The Cursed Merchant",
        FloorMin: 2, FloorMax: 4,
        Trigger: game.TriggerEnterRoom,
        Dialogues: []game.Dialogue{
            {Speaker: "Merchant", Text: "You look like someone who can help..."},
        },
        Reward: game.Reward{Gold: 100, XP: 80},
    })
}
```

Story beats with `Trigger: TriggerEnterRoom` are placed on `!` tiles in generated maps.
Story beats with `Trigger: TriggerBossDefeated` fire after defeating the named boss.

## Architecture

```
cmd/doomscroller/     Entry point
internal/engine/      Game loop (bubbletea model) + FSM
internal/game/        Core game logic (character, battle, map, story, perks, save)
internal/states/      Per-screen rendering and input handling
internal/ui/          Styles and rendering helpers
content/stories/      Story beat definitions (core + weekly/monthly)
```

### State Machine

```
Intro → MainMenu ↔ Settings
       ↓
  CharCreate → World ↔ Battle
                     ↔ Story
                     ↔ Pause → Settings
               ↓         ↓
             Win      GameOver
               ↓         ↓
           MainMenu   MainMenu
```
