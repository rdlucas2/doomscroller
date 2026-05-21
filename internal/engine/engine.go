// Package engine drives the game loop, FSM, and all state transitions using Ebiten.
package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/render"
	"github.com/rdlucas2/doomscroller/internal/states"
)

// Engine manages all game state and handles input.
// It does NOT implement ebiten.Game; that wrapper lives in cmd/doomscroller/main.go.
type Engine struct {
	fsm      *FSM
	input    *InputState
	settings game.Settings
	timeSec  float64

	// per-screen state objects
	intro      *states.IntroState
	menu       *states.MenuState
	settingsUI *states.SettingsState
	charCreate *states.CharCreateState
	world      *states.WorldState
	battleUI   *states.BattleUIState
	pauseUI    *states.PauseState
	storyUI    *states.StoryState
	gameOverUI *states.GameOverState
	winUI      *states.WinState

	session *game.Session
}

func New() *Engine {
	rand.Seed(time.Now().UnixNano())
	e := &Engine{
		fsm:      NewFSM(),
		input:    NewInputState(),
		settings: game.DefaultSettings(),
		intro:    &states.IntroState{},
	}
	return e
}

// Update is called by the Ebiten game every frame (60 fps). dt is seconds since last frame.
func (e *Engine) Update(dt float64) error {
	e.timeSec += dt
	e.input.Update()

	// Intro advances purely by time
	if e.fsm.Current() == StateIntro {
		e.intro.Update(dt)
		if e.intro.Done() {
			// Any key after loading finishes → go to menu
			for _, k := range e.input.JustPressedKeys() {
				_ = k
				e.enterMainMenu()
				return nil
			}
		}
		return nil
	}

	// Route pressed keys to the current state handler
	for _, key := range e.input.JustPressedKeys() {
		if err := e.handleKey(key); err != nil {
			return err
		}
	}

	// World: collect typed characters for name input (char creation)
	if e.fsm.Current() == StateCharCreate && e.charCreate != nil && e.charCreate.Phase == states.PhasePickName {
		for _, r := range PressedRunes() {
			e.charCreate.HandleRune(r)
		}
	}

	return nil
}

// BuildState produces the snapshot the renderer consumes this frame.
func (e *Engine) BuildState() *render.EngineState {
	st := &render.EngineState{
		State:   render.StateID(e.fsm.Current()),
		TimeSec: e.timeSec,
	}

	switch e.fsm.Current() {
	case StateIntro:
		st.IntroProgress = e.intro.Progress()
		st.IntroDone = e.intro.Done()

	case StateMainMenu:
		if e.menu != nil {
			st.MenuItems = e.menu.Items()
			st.MenuCursor = e.menu.Cursor
		}

	case StateSettings:
		if e.settingsUI != nil {
			labels, values, cur := e.settingsUI.ItemsAndValues()
			st.SettingsLabels = labels
			st.SettingsValues = values
			st.SettingsCursor = cur
		}

	case StateCharCreate:
		if e.charCreate != nil {
			cc := e.charCreate
			st.CharPhase = int(cc.Phase)
			st.CharName = cc.NameInput
			st.ClassCursor = cc.ClassCursor
			classes := []game.Class{game.Warrior, game.Mage, game.Rogue}
			st.ClassCards = make([]render.ClassCardData, len(classes))
			for i, cls := range classes {
				d := game.NewCharacter("x", cls)
				st.ClassCards[i] = render.ClassCardData{
					Name: cls.String(),
					Desc: cls.Description(),
					HP: d.MaxHP, MP: d.MaxMP,
					STR: d.STR, INT: d.INT, AGI: d.AGI, DEF: d.DEF,
				}
			}
			st.ClassName = classes[cc.ClassCursor].String()
			st.ConfirmCursor = cc.ConfirmCursor
			if cc.Phase == states.PhaseConfirm {
				c := cc.BuildCharacter()
				st.ConfirmStats = render.ConfirmStatsData{
					HP: c.MaxHP, MP: c.MaxMP,
					STR: c.STR, INT: c.INT, AGI: c.AGI, DEF: c.DEF,
					Weapon: c.Equipment.Weapon, Armor: c.Equipment.Armor,
				}
			}
		}

	case StateWorld:
		if e.world != nil && e.session != nil {
			st.Session = e.session
			st.WorldMsg = e.world.Message
			st.PlayerSprite = playerSpriteID(e.session.Character.Class, 0)
			// Walk animation: alternate every 8 steps
			st.WalkFrame = (e.session.StepCount / 4) % 2
		}

	case StateBattle:
		if e.battleUI != nil {
			st.Battle = e.battleUI.Battle
			if e.session != nil {
				st.PlayerSprite = playerSpriteID(e.session.Character.Class, 0)
			}
		}

	case StateStory:
		if e.storyUI != nil {
			st.Story = e.storyUI.Beat
			st.StoryLine = e.storyUI.Line
		}

	case StatePause:
		if e.pauseUI != nil {
			st.PauseCursor = e.pauseUI.Cursor
			st.PauseItems = []string{"Resume", "Save & Quit Run", "Settings", "Quit to Main Menu"}
		}

	case StateGameOver:
		if e.gameOverUI != nil {
			st.GameOverChar = e.gameOverUI.Character
			st.GameOverFloor = e.gameOverUI.Floor
			st.GameOverCursor = e.gameOverUI.Cursor
			st.GameOverItems = []string{"New Game", "Main Menu", "Quit"}
		}

	case StateWin:
		if e.winUI != nil {
			st.WinChar = e.winUI.Character
			st.WinCursor = e.winUI.Cursor
			st.WinItems = []string{"New Game", "Main Menu", "Quit"}
		}
	}
	return st
}

func playerSpriteID(cls game.Class, frame int) int {
	base := int(cls) * 2
	return base + frame
}

// ── state entry helpers ───────────────────────────────────────────────────────

func (e *Engine) enterMainMenu() {
	e.menu = states.NewMenuState()
	e.fsm.Transition(StateMainMenu)
}

func (e *Engine) enterSettings() {
	e.settingsUI = states.NewSettingsState(&e.settings)
	e.fsm.Push(StateSettings)
}

func (e *Engine) enterCharCreate() {
	e.charCreate = &states.CharCreateState{}
	e.fsm.Transition(StateCharCreate)
}

func (e *Engine) startNewGame(c *game.Character) {
	game.ResetStories()
	seed := time.Now().UnixNano()
	e.session = game.NewSession(c, seed, e.settings)
	e.world = states.NewWorldState(e.session)
	e.fsm.Transition(StateWorld)
}

func (e *Engine) loadGame() {
	data, err := game.Load()
	if err != nil || data.Character == nil {
		return
	}
	seed := time.Now().UnixNano()
	e.session = game.NewSession(data.Character, seed, e.settings)
	e.session.Floor = data.Floor
	e.session.Map = game.GenerateMap(e.session.Floor, e.session.Rng)
	e.world = states.NewWorldState(e.session)
	e.fsm.Transition(StateWorld)
}

func (e *Engine) enterBattle(enemies []*game.Enemy, storyAfter *game.StoryBeat, isBoss bool) {
	b := game.NewBattle(e.session.Character, enemies, e.session.Rng)
	b.IsBossEncounter = isBoss
	e.battleUI = states.NewBattleUIState(b, storyAfter)
	e.fsm.Push(StateBattle)
}

func (e *Engine) enterStory(beat *game.StoryBeat) {
	e.session.Character.Gold += beat.Reward.Gold
	e.session.Character.AddXP(beat.Reward.XP)
	e.storyUI = states.NewStoryState(beat)
	e.fsm.Push(StateStory)
}

func (e *Engine) enterPause() {
	e.pauseUI = &states.PauseState{}
	e.fsm.Push(StatePause)
}

func (e *Engine) enterGameOver() {
	var c *game.Character
	floor := 1
	if e.session != nil {
		c = e.session.Character
		floor = e.session.Floor
	}
	e.gameOverUI = &states.GameOverState{Character: c, Floor: floor}
	e.session = nil
	e.fsm.Transition(StateGameOver)
}

func (e *Engine) enterWin() {
	var c *game.Character
	if e.session != nil {
		c = e.session.Character
	}
	e.winUI = &states.WinState{Character: c}
	e.session = nil
	e.fsm.Transition(StateWin)
}

// ── key routing ───────────────────────────────────────────────────────────────

var errQuit = fmt.Errorf("quit")

func (e *Engine) handleKey(key string) error {
	switch e.fsm.Current() {
	case StateMainMenu:
		return e.handleMenuKey(key)
	case StateSettings:
		return e.handleSettingsKey(key)
	case StateCharCreate:
		return e.handleCharCreateKey(key)
	case StateWorld:
		return e.handleWorldKey(key)
	case StateBattle:
		return e.handleBattleKey(key)
	case StateStory:
		return e.handleStoryKey(key)
	case StatePause:
		return e.handlePauseKey(key)
	case StateGameOver:
		return e.handleGameOverKey(key)
	case StateWin:
		return e.handleWinKey(key)
	}
	return nil
}

func (e *Engine) handleMenuKey(key string) error {
	switch key {
	case "up", "k":
		e.menu.MoveUp()
	case "down", "j":
		e.menu.MoveDown()
	case "enter", " ":
		switch e.menu.Selected() {
		case states.MenuNewGame:
			e.enterCharCreate()
		case states.MenuContinue:
			e.loadGame()
		case states.MenuSettings:
			e.enterSettings()
		case states.MenuQuit:
			return errQuit
		}
	}
	return nil
}

func (e *Engine) handleSettingsKey(key string) error {
	s := e.settingsUI
	switch key {
	case "esc", "q":
		e.fsm.Pop()
	case "up", "k":
		s.MoveUp()
	case "down", "j":
		s.MoveDown()
	case "left", "h":
		s.AdjustLeft()
	case "right", "l":
		s.AdjustRight()
	case "enter", " ":
		if s.IsBack() {
			e.fsm.Pop()
		}
	}
	return nil
}

func (e *Engine) handleCharCreateKey(key string) error {
	cc := e.charCreate
	switch cc.Phase {
	case states.PhasePickName:
		switch key {
		case "esc":
			e.enterMainMenu()
		case "enter":
			if cc.IsNameReady() {
				cc.Phase = states.PhasePickClass
			}
		case "backspace":
			cc.HandleBackspace()
		}
	case states.PhasePickClass:
		switch key {
		case "esc":
			cc.Phase = states.PhasePickName
		case "left", "h", "a":
			cc.MoveClassLeft()
		case "right", "l", "d":
			cc.MoveClassRight()
		case "enter", " ":
			cc.Phase = states.PhaseConfirm
		}
	case states.PhaseConfirm:
		switch key {
		case "esc":
			cc.Phase = states.PhasePickClass
		case "up", "k":
			if cc.ConfirmCursor > 0 {
				cc.ConfirmCursor--
			}
		case "down", "j":
			if cc.ConfirmCursor < 1 {
				cc.ConfirmCursor++
			}
		case "enter", " ":
			if cc.ConfirmCursor == 0 {
				e.startNewGame(cc.BuildCharacter())
			} else {
				cc.Phase = states.PhasePickClass
			}
		}
	}
	return nil
}

func (e *Engine) handleWorldKey(key string) error {
	switch key {
	case "p", "esc":
		e.enterPause()
		return nil
	}
	var dx, dy int
	switch key {
	case "w", "up":
		dy = -1
	case "s", "down":
		dy = 1
	case "a", "left":
		dx = -1
	case "d", "right":
		dx = 1
	default:
		return nil
	}
	event := e.world.Move(dx, dy)
	return e.handleWorldEvent(event)
}

func (e *Engine) handleWorldEvent(event states.WorldEvent) error {
	w := e.world
	switch event {
	case states.WorldEventRandomBattle:
		enemies := game.RandomEncounterEnemies(e.session.Floor, e.session.Rng)
		e.enterBattle(enemies, nil, false)
	case states.WorldEventStoryBeat:
		if len(w.StoryQueue) > 0 {
			beat := w.StoryQueue[0]
			w.StoryQueue = w.StoryQueue[1:]
			e.enterStory(beat)
		}
	case states.WorldEventBoss:
		px, py := e.session.Map.PlayerX, e.session.Map.PlayerY
		tile := e.session.Map.At(px, py)
		bossName := "Floor Guardian"
		var storyForBoss *game.StoryBeat
		if tile != nil && tile.BossName != "" {
			bossName = tile.BossName
		}
		for _, s := range game.GetStoriesForFloor(e.session.Floor) {
			if s.BossName == bossName {
				storyForBoss = s
				break
			}
		}
		boss := game.BossEnemy(bossName, e.session.Floor)
		e.enterBattle([]*game.Enemy{boss}, storyForBoss, true)
	case states.WorldEventNextFloor:
		e.session.NextFloor()
		e.world = states.NewWorldState(e.session)
		e.world.SetMessage("You descend deeper into the tower...")
	case states.WorldEventWin:
		e.enterWin()
	}
	return nil
}

func (e *Engine) handleBattleKey(key string) error {
	b := e.battleUI.Battle
	switch b.Phase {
	case game.PhasePerkChoice:
		e.battleUI.HandlePerkChoice(key)
	case game.PhaseVictory:
		if key == "enter" || key == " " {
			e.fsm.Pop()
			if e.battleUI.StoryAfterBattle != nil && !e.battleUI.StoryAfterBattle.Completed {
				e.enterStory(e.battleUI.StoryAfterBattle)
			}
		}
	case game.PhaseDefeat:
		if key == "enter" || key == " " {
			e.fsm.Pop()
			if e.session != nil && e.session.Character.HP <= 0 {
				e.enterGameOver()
			}
		}
	default:
		e.battleUI.HandleInput(key)
	}
	return nil
}

func (e *Engine) handleStoryKey(key string) error {
	ss := e.storyUI
	switch key {
	case "enter", " ":
		ss.Advance()
		if ss.Done {
			e.fsm.Pop()
		}
	case "esc":
		game.MarkCompleted(ss.Beat.ID)
		e.fsm.Pop()
	}
	return nil
}

func (e *Engine) handlePauseKey(key string) error {
	ps := e.pauseUI
	switch key {
	case "p", "esc":
		e.fsm.Pop()
	case "up", "k":
		ps.MoveUp()
	case "down", "j":
		ps.MoveDown()
	case "enter", " ":
		switch states.PauseChoice(ps.Cursor) {
		case states.PauseResume:
			e.fsm.Pop()
		case states.PauseSave:
			if e.session != nil {
				_ = game.Save(e.session.ToSaveData())
			}
			e.fsm.Pop()
			if e.world != nil {
				e.world.SetMessage("Game saved!")
			}
		case states.PauseSettings:
			e.enterSettings()
		case states.PauseQuitToMenu:
			e.session = nil
			e.enterMainMenu()
		}
	}
	return nil
}

func (e *Engine) handleGameOverKey(key string) error {
	g := e.gameOverUI
	switch key {
	case "up", "k":
		g.MoveUp()
	case "down", "j":
		g.MoveDown()
	case "enter", " ":
		switch g.Selected() {
		case "New Game":
			e.enterCharCreate()
		case "Main Menu":
			e.enterMainMenu()
		case "Quit":
			return errQuit
		}
	}
	return nil
}

func (e *Engine) handleWinKey(key string) error {
	w := e.winUI
	switch key {
	case "up", "k":
		w.MoveUp()
	case "down", "j":
		w.MoveDown()
	case "enter", " ":
		switch w.Selected() {
		case "New Game":
			e.enterCharCreate()
		case "Main Menu":
			e.enterMainMenu()
		case "Quit":
			return errQuit
		}
	}
	return nil
}

// IsQuit returns true when the engine signalled an exit.
func IsQuit(err error) bool { return err == errQuit }
