package render

import "github.com/rdlucas2/doomscroller/internal/game"

// StateID mirrors engine.StateID (defined here to avoid import cycles).
type StateID int

const (
	StateIntro StateID = iota
	StateMainMenu
	StateSettings
	StateCharCreate
	StateWorld
	StateBattle
	StatePause
	StateGameOver
	StateWin
	StateStory
)

// EngineState is built by the engine each frame and consumed by the renderer.
// The engine imports render; the renderer only uses this struct (no engine import).
type EngineState struct {
	State   StateID
	TimeSec float64 // seconds since game start

	// Intro
	IntroProgress float64 // 0–1
	IntroDone     bool

	// Menu
	MenuItems  []string
	MenuCursor int

	// Settings
	SettingsLabels []string
	SettingsValues []string
	SettingsCursor int

	// Character creation
	CharPhase     int // 0=name 1=class 2=confirm
	CharName      string
	ClassCursor   int
	ClassName     string           // chosen class name
	ClassCards    []ClassCardData  // one per class
	ConfirmCursor int
	ConfirmStats  ConfirmStatsData

	// World
	Session      *game.Session
	WorldMsg     string
	PlayerSprite int // sprite sheet ID
	WalkFrame    int // 0 or 1

	// Battle
	Battle *game.Battle

	// Story
	Story     *game.StoryBeat
	StoryLine int

	// Pause
	PauseCursor int
	PauseItems  []string

	// Game Over
	GameOverChar   *game.Character
	GameOverFloor  int
	GameOverCursor int
	GameOverItems  []string

	// Win
	WinChar   *game.Character
	WinCursor int
	WinItems  []string
}

// ClassCardData holds display data for one class card.
type ClassCardData struct {
	Name string
	Desc string
	HP, MP, STR, INT, AGI, DEF int
}

// ConfirmStatsData holds the final character preview stats.
type ConfirmStatsData struct {
	HP, MP, STR, INT, AGI, DEF int
	Weapon, Armor              string
}
