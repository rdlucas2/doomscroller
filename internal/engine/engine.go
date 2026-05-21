package engine

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/states"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type worldTickMsg struct{}

func worldTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return worldTickMsg{}
	})
}

// Engine is the top-level bubbletea model.
type Engine struct {
	fsm      *FSM
	width    int
	height   int
	settings game.Settings

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
	s := game.DefaultSettings()
	e := &Engine{
		fsm:      NewFSM(),
		settings: s,
		width:    80,
		height:   24,
		intro:    &states.IntroState{},
	}
	return e
}

func (e *Engine) Init() tea.Cmd {
	return tea.Batch(e.intro.Init(), worldTick())
}

func (e *Engine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		e.width, e.height = m.Width, m.Height
		return e, nil

	case tea.KeyMsg:
		return e.handleKey(m.String())

	case states.IntroTickMsg:
		done := e.intro.Update(states.IntroTickMsg{})
		if done {
			return e, nil
		}
		return e, e.intro.NextCmd()

	case worldTickMsg:
		if e.fsm.Current() == StateWorld && e.world != nil {
			e.world.TickMessage()
		}
		return e, worldTick()
	}
	return e, nil
}

// ── State entry helpers ───────────────────────────────────────────────────────

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

// ── Key handlers ─────────────────────────────────────────────────────────────

func (e *Engine) handleKey(key string) (tea.Model, tea.Cmd) {
	switch e.fsm.Current() {
	case StateIntro:
		if e.intro.Done() {
			e.enterMainMenu()
		}
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
	return e, nil
}

func (e *Engine) handleMenuKey(key string) (tea.Model, tea.Cmd) {
	m := e.menu
	switch key {
	case "ctrl+c":
		return e, tea.Quit
	case "up", "k":
		m.MoveUp()
	case "down", "j":
		m.MoveDown()
	case "enter", " ":
		switch m.Selected() {
		case states.MenuNewGame:
			e.enterCharCreate()
		case states.MenuContinue:
			e.loadGame()
		case states.MenuSettings:
			e.enterSettings()
		case states.MenuQuit:
			return e, tea.Quit
		}
	}
	return e, nil
}

func (e *Engine) handleSettingsKey(key string) (tea.Model, tea.Cmd) {
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
	return e, nil
}

func (e *Engine) handleCharCreateKey(key string) (tea.Model, tea.Cmd) {
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
		default:
			cc.HandleNameKey(key)
		}
	case states.PhasePickClass:
		switch key {
		case "esc":
			cc.Phase = states.PhasePickName
		case "left", "h":
			cc.MoveClassLeft()
		case "right", "l":
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
				c := cc.BuildCharacter()
				e.startNewGame(c)
			} else {
				cc.Phase = states.PhasePickClass
			}
		}
	}
	return e, nil
}

func (e *Engine) handleWorldKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c":
		return e, tea.Quit
	case "p", "P", "esc":
		e.enterPause()
		return e, nil
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
		return e, nil
	}

	event := e.world.Move(dx, dy)
	return e.handleWorldEvent(event)
}

func (e *Engine) handleWorldEvent(event states.WorldEvent) (tea.Model, tea.Cmd) {
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
		e.world.SetMessage(ui.SuccessStyle.Render("You descend deeper into the tower..."))

	case states.WorldEventWin:
		e.enterWin()
	}
	return e, nil
}

func (e *Engine) handleBattleKey(key string) (tea.Model, tea.Cmd) {
	b := e.battleUI.Battle

	switch b.Phase {
	case game.PhasePerkChoice:
		e.battleUI.HandlePerkChoice(key)
	case game.PhaseVictory:
		if key == "enter" || key == " " {
			e.fsm.Pop()
			// After boss victory, show post-battle story
			if e.battleUI.StoryAfterBattle != nil && !e.battleUI.StoryAfterBattle.Completed {
				beat := e.battleUI.StoryAfterBattle
				e.enterStory(beat)
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
	return e, nil
}

func (e *Engine) handleStoryKey(key string) (tea.Model, tea.Cmd) {
	ss := e.storyUI
	switch key {
	case "enter", " ":
		ss.Advance()
		if ss.Done {
			e.fsm.Pop()
		}
	case "esc":
		// Allow skipping story
		game.MarkCompleted(ss.Beat.ID)
		e.fsm.Pop()
	}
	return e, nil
}

func (e *Engine) handlePauseKey(key string) (tea.Model, tea.Cmd) {
	ps := e.pauseUI
	switch key {
	case "ctrl+c":
		return e, tea.Quit
	case "p", "P", "esc":
		e.fsm.Pop()
	case "up", "k":
		ps.MoveUp()
	case "down", "j":
		ps.MoveDown()
	case "enter", " ":
		switch ps.Selected() {
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
	return e, nil
}

func (e *Engine) handleGameOverKey(key string) (tea.Model, tea.Cmd) {
	g := e.gameOverUI
	switch key {
	case "ctrl+c":
		return e, tea.Quit
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
			return e, tea.Quit
		}
	}
	return e, nil
}

func (e *Engine) handleWinKey(key string) (tea.Model, tea.Cmd) {
	w := e.winUI
	switch key {
	case "ctrl+c":
		return e, tea.Quit
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
			return e, tea.Quit
		}
	}
	return e, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func (e *Engine) View() string {
	w, h := e.width, e.height
	if w < 40 {
		w = 80
	}
	if h < 20 {
		h = 24
	}

	switch e.fsm.Current() {
	case StateIntro:
		return e.intro.View(w, h)
	case StateMainMenu:
		if e.menu != nil {
			return e.menu.View(w, h)
		}
	case StateSettings:
		if e.settingsUI != nil {
			return e.settingsUI.View(w, h)
		}
	case StateCharCreate:
		if e.charCreate != nil {
			return e.charCreate.View(w, h)
		}
	case StateWorld:
		if e.world != nil {
			return e.world.View(w, h)
		}
	case StateBattle:
		if e.battleUI != nil {
			return e.battleUI.View(w, h)
		}
	case StateStory:
		if e.storyUI != nil {
			return e.storyUI.View(w, h)
		}
	case StatePause:
		if e.pauseUI != nil {
			return e.pauseUI.View(w, h)
		}
	case StateGameOver:
		if e.gameOverUI != nil {
			return e.gameOverUI.View(w, h)
		}
	case StateWin:
		if e.winUI != nil {
			return e.winUI.View(w, h)
		}
	}
	return ""
}
