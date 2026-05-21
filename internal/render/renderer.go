// Package render implements all graphical drawing using Ebiten v2.
package render

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rdlucas2/doomscroller/internal/assets"
	"github.com/rdlucas2/doomscroller/internal/game"
)

// Internal rendering resolution (displayed at ×2 scale via Layout).
const (
	ScreenW = 640
	ScreenH = 360
)

// colour palette
var (
	colBlack    = color.RGBA{10, 8, 18, 255}
	colDarkBg   = color.RGBA{14, 10, 28, 255}
	colPanelBg  = color.RGBA{18, 14, 38, 255}
	colWindowBg = color.RGBA{12, 20, 50, 255}
	colGold     = color.RGBA{190, 155, 30, 255}
	colBlue     = color.RGBA{40, 80, 180, 255}
	colWhite    = color.RGBA{240, 238, 232, 255}
	colTextGold = color.RGBA{230, 195, 60, 255}
	colCyan     = color.RGBA{80, 210, 240, 255}
	colGray     = color.RGBA{140, 138, 130, 255}
	colGreen    = color.RGBA{60, 220, 100, 255}
	colRed      = color.RGBA{240, 60, 60, 255}
	colHPFill   = color.RGBA{220, 55, 55, 255}
	colMPFill   = color.RGBA{55, 100, 220, 255}
	colXPFill   = color.RGBA{60, 200, 200, 255}
	colBarBg    = color.RGBA{30, 28, 42, 255}
	colSel      = color.RGBA{240, 220, 60, 255}
	colDimBlue  = color.RGBA{28, 40, 88, 255}
)

// Renderer holds cached Ebiten images and the bitmap font face.
type Renderer struct {
	tileset *ebiten.Image
	sprites *ebiten.Image
	golem   *ebiten.Image
	void    *ebiten.Image
	face    etext.Face
	ready   bool
}

func New() *Renderer { return &Renderer{} }

func (r *Renderer) Init() {
	if r.ready {
		return
	}
	r.tileset = ebiten.NewImageFromImage(assets.GenerateTileset())
	r.sprites = ebiten.NewImageFromImage(assets.GenerateSpriteSheet())
	r.golem = ebiten.NewImageFromImage(assets.GenGolem())
	r.void = ebiten.NewImageFromImage(assets.GenVoidSovereign())
	r.face = etext.NewGoXFace(defaultFace())
	r.ready = true
}

// Draw dispatches to the appropriate screen renderer.
func (r *Renderer) Draw(screen *ebiten.Image, st *EngineState) {
	r.Init()
	screen.Fill(colDarkBg)
	switch st.State {
	case StateIntro:
		r.drawIntro(screen, st)
	case StateMainMenu:
		r.drawMainMenu(screen, st)
	case StateSettings:
		r.drawSettings(screen, st)
	case StateCharCreate:
		r.drawCharCreate(screen, st)
	case StateWorld:
		r.drawWorld(screen, st)
	case StateBattle:
		r.drawBattle(screen, st)
	case StateStory:
		r.drawStory(screen, st)
	case StatePause:
		r.drawPause(screen, st)
	case StateGameOver:
		r.drawGameOver(screen, st)
	case StateWin:
		r.drawWin(screen, st)
	}
}

// ── INTRO ─────────────────────────────────────────────────────────────────────

func (r *Renderer) drawIntro(screen *ebiten.Image, st *EngineState) {
	p := st.IntroProgress
	cx := float64(ScreenW) / 2
	if p >= 0 {
		r.tc(screen, "DOOM SCROLLER", cx, 110, 3, colTextGold)
		r.tc(screen, "Tower of the Void", cx, 148, 2, colCyan)
	}
	if p >= 0.3 {
		r.tc(screen, "~ A JRPG Roguelike ~", cx, 190, 1, colGray)
	}
	if p >= 0.5 && p < 1.0 {
		pct := (p - 0.5) / 0.5
		bx := float64(ScreenW)/2 - 120
		r.drawBar(screen, int(bx), 240, 240, 10, pct, colXPFill, colBarBg)
		r.tc(screen, fmt.Sprintf("Loading... %d%%", int(pct*100)), cx, 258, 1, colGray)
	}
	if p >= 1.0 && int(st.TimeSec*4)%2 == 0 {
		r.tc(screen, "Press any key to continue", cx, 260, 1, colGreen)
	}
}

// ── MAIN MENU ────────────────────────────────────────────────────────────────

func (r *Renderer) drawMainMenu(screen *ebiten.Image, st *EngineState) {
	cx := float64(ScreenW) / 2
	r.tc(screen, "DOOM SCROLLER", cx, 50, 3, colTextGold)
	r.tc(screen, "Tower of the Void", cx, 90, 1, colCyan)
	r.hRule(screen, 115, colGold)
	for i, item := range st.MenuItems {
		c := colWhite
		prefix := "  "
		if i == st.MenuCursor {
			c = colSel
			prefix = "▶ "
		}
		r.tc(screen, prefix+item, cx, float64(145+i*32), 2, c)
	}
	r.tc(screen, "↑↓ Navigate    Enter Select", cx, 330, 1, colGray)
}

// ── SETTINGS ────────────────────────────────────────────────────────────────

func (r *Renderer) drawSettings(screen *ebiten.Image, st *EngineState) {
	cx := float64(ScreenW) / 2
	r.winBox(screen, 80, 30, 480, 300, "SETTINGS")
	for i, label := range st.SettingsLabels {
		y := 80.0 + float64(i)*38
		c := colWhite
		if i == st.SettingsCursor {
			c = colSel
			r.fillRect(screen, 84, int(y)-2, 472, 30, colDimBlue)
			r.t(screen, "▶", 88, y, 1, colSel)
		}
		r.t(screen, label, 110, y, 1, c)
		if i < len(st.SettingsValues) {
			r.t(screen, st.SettingsValues[i], 350, y, 1, colCyan)
		}
	}
	r.tc(screen, "↑↓ Navigate  ◄► Adjust  Enter/Esc Back", cx, 340, 1, colGray)
}

// ── CHARACTER CREATION ────────────────────────────────────────────────────────

func (r *Renderer) drawCharCreate(screen *ebiten.Image, st *EngineState) {
	cx := float64(ScreenW) / 2
	r.tc(screen, "CHARACTER CREATION", cx, 22, 2, colTextGold)

	switch st.CharPhase {
	case 0:
		r.tc(screen, "Enter your hero's name:", cx, 100, 1, colCyan)
		cursor := ""
		if int(st.TimeSec*3)%2 == 0 {
			cursor = "█"
		}
		r.winBox(screen, 200, 128, 240, 44, "")
		r.tc(screen, st.CharName+cursor, cx, 144, 2, colWhite)
		hint := "Type your name and press Enter"
		if st.CharName != "" {
			hint = "Press Enter to continue"
		}
		r.tc(screen, hint, cx, 210, 1, colGray)

	case 1:
		r.tc(screen, "Choose your class:", cx, 90, 1, colCyan)
		n := len(st.ClassCards)
		cardW := 170
		gap := 8
		total := n*(cardW+gap) - gap
		startX := ScreenW/2 - total/2
		for i, card := range st.ClassCards {
			bx := startX + i*(cardW+gap)
			r.classCard(screen, bx, 118, cardW, 195, card, i == st.ClassCursor)
		}
		r.tc(screen, "◄► Change    Enter Confirm", cx, 330, 1, colGray)

	case 2:
		cs := st.ConfirmStats
		r.winBox(screen, 100, 70, 440, 220, st.CharName+" the "+st.ClassName)
		lines := []struct {
			text string
			c    color.RGBA
		}{
			{fmt.Sprintf("HP: %d      MP: %d", cs.HP, cs.MP), colWhite},
			{fmt.Sprintf("STR: %d   INT: %d   AGI: %d   DEF: %d", cs.STR, cs.INT, cs.AGI, cs.DEF), colWhite},
			{"", colGray},
			{"Weapon: " + cs.Weapon, colCyan},
			{"Armor:  " + cs.Armor, colCyan},
		}
		for j, l := range lines {
			r.tc(screen, l.text, cx, float64(118+j*26), 1, l.c)
		}
		confirmItems := []string{"Begin Adventure", "Back"}
		for i, item := range confirmItems {
			c := colWhite
			if i == st.ConfirmCursor {
				c = colSel
			}
			r.tc(screen, item, cx, float64(310+i*28), 2, c)
		}
	}
}

func (r *Renderer) classCard(screen *ebiten.Image, x, y, w, h int, card ClassCardData, selected bool) {
	bc := colBlue
	if selected {
		bc = colGold
	}
	r.fillRect(screen, x, y, w, h, colPanelBg)
	r.border(screen, x, y, w, h, bc)
	cx := float64(x + w/2)
	nc := colWhite
	if selected {
		nc = colSel
	}
	r.tc(screen, card.Name, cx, float64(y+8), 2, nc)
	if selected {
		r.tc(screen, "★", cx, float64(y+32), 2, colTextGold)
	}
	r.tc(screen, fmt.Sprintf("HP:%-3d  MP:%-3d", card.HP, card.MP), cx, float64(y+70), 1, colWhite)
	r.tc(screen, fmt.Sprintf("STR:%-3d INT:%-3d", card.STR, card.INT), cx, float64(y+86), 1, colWhite)
	r.tc(screen, fmt.Sprintf("AGI:%-3d DEF:%-3d", card.AGI, card.DEF), cx, float64(y+102), 1, colWhite)
	desc := wrapText(card.Desc, w/7-1)
	for i, l := range desc {
		r.tc(screen, l, cx, float64(y+124+i*14), 1, colGray)
	}
}

// ── WORLD MAP ────────────────────────────────────────────────────────────────

const (
	mapViewW = 25
	mapViewH = 17
)

func (r *Renderer) drawWorld(screen *ebiten.Image, st *EngineState) {
	if st.Session == nil {
		return
	}
	m := st.Session.Map
	c := st.Session.Character
	hudX := mapViewW*assets.TileSize + 4

	// Tile map
	camX := m.PlayerX - mapViewW/2
	camY := m.PlayerY - mapViewH/2
	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX+mapViewW > m.Width {
		camX = m.Width - mapViewW
	}
	if camY+mapViewH > m.Height {
		camY = m.Height - mapViewH
	}

	for vy := 0; vy < mapViewH; vy++ {
		for vx := 0; vx < mapViewW; vx++ {
			mx, my := camX+vx, camY+vy
			tile := m.At(mx, my)
			if tile == nil {
				continue
			}
			visible := st.Session.Settings.ShowFullMap ||
				tile.Visited ||
				(iabs(mx-m.PlayerX) <= 5 && iabs(my-m.PlayerY) <= 5)
			if !visible {
				continue
			}
			r.drawTile(screen, gameTile(tile.Type), vx*assets.TileSize, vy*assets.TileSize, 1.0)
		}
	}
	// Player @ position
	sprID := st.PlayerSprite
	if st.WalkFrame == 1 {
		sprID++ // alternate frame
	}
	px := (m.PlayerX - camX) * assets.TileSize
	py := (m.PlayerY - camY) * assets.TileSize
	r.drawSprite(screen, sprID, float64(px), float64(py), 1.0)

	// HUD background
	r.fillRect(screen, hudX, 0, ScreenW-hudX, ScreenH, colPanelBg)

	hy := 8.0
	r.t(screen, fmt.Sprintf("Floor %d/%d", st.Session.Floor, game.MaxFloors), float64(hudX+6), hy, 1, colTextGold)
	hy += 18
	r.t(screen, c.Name, float64(hudX+6), hy, 1, colWhite)
	r.t(screen, c.Class.String(), float64(hudX+6+len(c.Name)*8+4), hy, 1, colCyan)
	r.t(screen, fmt.Sprintf("Lv%d", c.Level), float64(ScreenW-28), hy, 1, colTextGold)

	hy += 22
	r.t(screen, "HP", float64(hudX+6), hy, 1, colRed)
	r.drawBar(screen, hudX+22, int(hy), 86, 7, float64(c.HP)/float64(c.MaxHP), colHPFill, colBarBg)
	r.t(screen, fmt.Sprintf("%d/%d", c.HP, c.MaxHP), float64(hudX+112), hy, 1, colRed)

	hy += 14
	r.t(screen, "MP", float64(hudX+6), hy, 1, colMPFill)
	r.drawBar(screen, hudX+22, int(hy), 86, 7, float64(c.MP)/float64(c.MaxMP), colMPFill, colBarBg)
	r.t(screen, fmt.Sprintf("%d/%d", c.MP, c.MaxMP), float64(hudX+112), hy, 1, colMPFill)

	hy += 14
	r.t(screen, "XP", float64(hudX+6), hy, 1, colXPFill)
	r.drawBar(screen, hudX+22, int(hy), 86, 7, float64(c.XP)/float64(c.XPToNext), colXPFill, colBarBg)

	hy += 16
	r.t(screen, fmt.Sprintf("Gold: %d", c.Gold), float64(hudX+6), hy, 1, colTextGold)

	hy += 18
	r.t(screen, "PERKS", float64(hudX+6), hy, 1, colCyan)
	hy += 13
	if len(c.Perks) == 0 {
		r.t(screen, "None", float64(hudX+10), hy, 1, colGray)
	}
	for _, p := range c.Perks {
		r.t(screen, "• "+p.Name, float64(hudX+10), hy, 1, perkCol(p.Type))
		hy += 12
		if hy > float64(ScreenH-50) {
			break
		}
	}

	// Legend at bottom of HUD
	r.t(screen, "> Stairs  ! Story  B Boss", float64(hudX+6), float64(ScreenH-38), 1, colGray)
	r.t(screen, "$ Gold  + Heal", float64(hudX+6), float64(ScreenH-26), 1, colGray)
	r.t(screen, "WASD/Arrows Move  P Pause", float64(hudX+6), float64(ScreenH-13), 1, colGray)

	// Message bar
	if st.WorldMsg != "" {
		r.fillRect(screen, 0, ScreenH-20, mapViewW*assets.TileSize, 20, color.RGBA{5, 30, 5, 200})
		r.t(screen, "▶ "+st.WorldMsg, 6, float64(ScreenH-12), 1, colGreen)
	}
}

// ── BATTLE ───────────────────────────────────────────────────────────────────

func (r *Renderer) drawBattle(screen *ebiten.Image, st *EngineState) {
	b := st.Battle
	if b == nil {
		return
	}

	// Tiled floor bg
	for y := 0; y < ScreenH/assets.TileSize+1; y++ {
		for x := 0; x < ScreenW/assets.TileSize+1; x++ {
			r.drawTile(screen, assets.TileFloor, x*assets.TileSize, y*assets.TileSize, 1.0)
		}
	}
	overlay := ebiten.NewImage(ScreenW, ScreenH)
	overlay.Fill(color.RGBA{0, 0, 0, 155})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	titleStr := "BATTLE"
	titleC := colTextGold
	if b.IsBossEncounter {
		titleStr = "BOSS BATTLE"
		titleC = colRed
	}
	r.tc(screen, titleStr, ScreenW/2, 10, 2, titleC)

	// Enemy display zone (left 2/3)
	eZoneW := ScreenW * 2 / 3
	eCount := len(b.Enemies)
	if eCount == 0 {
		eCount = 1
	}
	eSpacing := eZoneW / (eCount + 1)
	lives := b.LiveEnemies()

	for i, e := range b.Enemies {
		ex := eSpacing*(i+1) - assets.TileSize/2
		ey := 45
		scale := 1.0
		sid := enemySpr(e.Name)
		if e.IsBoss {
			scale = 2.0
			ex -= 16
			ey -= 10
		}
		if e.IsAlive() {
			// Target highlight ring
			isTarget := false
			for ti, le := range lives {
				if le == e && ti == b.SelectedTarget && b.Phase == game.PhasePlayerMenu {
					isTarget = true
				}
			}
			if isTarget {
				sz := int(float64(assets.TileSize)*scale) + 4
				r.border(screen, ex-2, ey-2, sz, sz, colSel)
			}
			r.drawSprite(screen, sid, float64(ex), float64(ey), scale)
			barW := 40
			if e.IsBoss {
				barW = 70
			}
			barY := ey + int(float64(assets.TileSize)*scale) + 3
			r.drawBar(screen, ex, barY, barW, 4, float64(e.HP)/float64(e.MaxHP), colHPFill, colBarBg)
			r.t(screen, e.Name, float64(ex), float64(barY+8), 1, colWhite)
		} else {
			r.t(screen, e.Name+" [X]", float64(ex), float64(ey), 1, colGray)
		}
	}

	// Player area (bottom left)
	pSprScale := 2.0
	pSprID := st.PlayerSprite
	r.drawSprite(screen, pSprID, 20, 155, pSprScale)
	c := b.Player
	px := 60.0
	py := 155.0
	r.t(screen, fmt.Sprintf("%s  Lv%d", c.Name, c.Level), px, py, 1, colWhite)
	r.t(screen, "HP", px, py+14, 1, colRed)
	r.drawBar(screen, int(px)+16, int(py+14), 80, 6, float64(c.HP)/float64(c.MaxHP), colHPFill, colBarBg)
	r.t(screen, fmt.Sprintf("%d/%d", c.HP, c.MaxHP), px+100, py+14, 1, colRed)
	r.t(screen, "MP", px, py+26, 1, colMPFill)
	r.drawBar(screen, int(px)+16, int(py+26), 80, 6, float64(c.MP)/float64(c.MaxMP), colMPFill, colBarBg)
	r.t(screen, fmt.Sprintf("%d/%d", c.MP, c.MaxMP), px+100, py+26, 1, colMPFill)

	// Battle log
	logY := 210
	r.fillRect(screen, 0, logY-2, eZoneW, 70, color.RGBA{5, 5, 20, 200})
	r.border(screen, 0, logY-2, eZoneW, 70, colBlue)
	for i, msg := range b.Log {
		if i >= 5 {
			break
		}
		r.t(screen, msg, 5, float64(logY+2+i*13), 1, colWhite)
	}

	// Action panel (right 1/3)
	mX := eZoneW + 4
	mW := ScreenW - mX - 4

	switch b.Phase {
	case game.PhasePlayerMenu:
		r.winBox(screen, mX, 45, mW, 105, "ACTION")
		for i, name := range game.ActionNames {
			c := colWhite
			if i == b.SelectedAction {
				c = colSel
				r.t(screen, "▶", float64(mX+6), float64(68+i*22), 1, c)
			}
			r.t(screen, name, float64(mX+18), float64(68+i*22), 1, c)
		}
		r.t(screen, "◄►Target ↑↓Action", float64(mX+4), float64(155), 1, colGray)

	case game.PhasePerkChoice:
		r.drawPerkChoice(screen, b)

	case game.PhaseVictory:
		r.winBox(screen, mX, 45, mW, 90, "VICTORY!")
		r.t(screen, fmt.Sprintf("+%d XP", b.TotalXP), float64(mX+8), 76, 1, colXPFill)
		r.t(screen, fmt.Sprintf("+%d Gold", b.TotalGold), float64(mX+8), 92, 1, colTextGold)
		r.t(screen, "Enter continue", float64(mX+8), 114, 1, colGray)

	case game.PhaseDefeat:
		r.winBox(screen, mX, 45, mW, 60, "")
		msg := "Escaped!"
		if b.Player.HP <= 0 {
			msg = "Defeated..."
		}
		r.t(screen, msg, float64(mX+8), 72, 1, colRed)
		r.t(screen, "Enter continue", float64(mX+8), 88, 1, colGray)
	}
}

func (r *Renderer) drawPerkChoice(screen *ebiten.Image, b *game.Battle) {
	r.fillRect(screen, 20, 50, ScreenW-40, 260, colWindowBg)
	r.border(screen, 20, 50, ScreenW-40, 260, colGold)
	r.tc(screen, "Choose a Perk!", ScreenW/2, 62, 2, colTextGold)

	n := len(b.PerkChoices)
	if n == 0 {
		return
	}
	cw := (ScreenW - 60) / n
	for i, p := range b.PerkChoices {
		cx := 30 + i*cw
		sel := i == b.SelectedPerk
		bc := colBlue
		if sel {
			bc = colGold
		}
		r.fillRect(screen, cx, 90, cw-6, 170, colPanelBg)
		r.border(screen, cx, 90, cw-6, 170, bc)
		mid := float64(cx + cw/2 - 3)
		r.tc(screen, p.Name, mid, 102, 1, perkCol(p.Type))
		r.tc(screen, "["+p.Type.String()+"]", mid, 116, 1, colGray)
		lines := wrapText(p.Desc, cw/7-1)
		for j, l := range lines {
			r.tc(screen, l, mid, float64(134+j*14), 1, colWhite)
		}
		if sel {
			r.tc(screen, "▶ Select", mid, 244, 1, colSel)
		}
	}
	r.tc(screen, "◄► Navigate   Enter Select", ScreenW/2, 268, 1, colGray)
}

// ── STORY ─────────────────────────────────────────────────────────────────────

func (r *Renderer) drawStory(screen *ebiten.Image, st *EngineState) {
	beat := st.Story
	if beat == nil {
		return
	}
	screen.Fill(color.RGBA{5, 8, 20, 255})
	r.tc(screen, beat.Name, ScreenW/2, 25, 2, colTextGold)
	r.hRule(screen, 52, colGold)

	if len(beat.Dialogues) == 0 {
		return
	}
	d := beat.Dialogues[st.StoryLine]
	r.winBox(screen, 30, 65, ScreenW-60, 165, d.Speaker)
	for i, l := range wrapText(d.Text, 68) {
		r.t(screen, l, 45, float64(92+i*17), 1, colWhite)
	}
	r.tc(screen, fmt.Sprintf("(%d/%d)", st.StoryLine+1, len(beat.Dialogues)), ScreenW/2, 245, 1, colGray)

	if st.StoryLine == len(beat.Dialogues)-1 {
		rew := beat.Reward
		parts := []string{}
		if rew.Gold > 0 {
			parts = append(parts, fmt.Sprintf("+%d Gold", rew.Gold))
		}
		if rew.XP > 0 {
			parts = append(parts, fmt.Sprintf("+%d XP", rew.XP))
		}
		if rew.ItemName != "" {
			parts = append(parts, rew.ItemName)
		}
		if len(parts) > 0 {
			r.tc(screen, "Reward: "+strings.Join(parts, "  "), ScreenW/2, 268, 1, colTextGold)
		}
		r.tc(screen, "Enter to continue", ScreenW/2, 292, 1, colGreen)
	} else {
		r.tc(screen, "Enter / Space advance   Esc skip", ScreenW/2, 292, 1, colGray)
	}
}

// ── PAUSE ─────────────────────────────────────────────────────────────────────

func (r *Renderer) drawPause(screen *ebiten.Image, st *EngineState) {
	overlay := ebiten.NewImage(ScreenW, ScreenH)
	overlay.Fill(color.RGBA{0, 0, 0, 170})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	r.winBox(screen, ScreenW/2-130, 80, 260, 180, "PAUSED")
	for i, item := range st.PauseItems {
		c := colWhite
		if i == st.PauseCursor {
			c = colSel
		}
		r.tc(screen, item, ScreenW/2, float64(128+i*32), 1, c)
	}
	r.tc(screen, "↑↓ Navigate  Enter Select  P Resume", ScreenW/2, 290, 1, colGray)
}

// ── GAME OVER ────────────────────────────────────────────────────────────────

func (r *Renderer) drawGameOver(screen *ebiten.Image, st *EngineState) {
	cx := float64(ScreenW) / 2
	for y := 0; y < ScreenH; y++ {
		for x := 0; x < ScreenW; x++ {
			dx := float64(x-ScreenW/2) / float64(ScreenW/2)
			dy := float64(y-ScreenH/2) / float64(ScreenH/2)
			v := uint8(iclamp(int((math.Sqrt(dx*dx+dy*dy)*0.5)*80), 0, 80))
			screen.Set(x, y, color.RGBA{v, 0, 0, 255})
		}
	}
	r.tc(screen, "GAME OVER", cx, 75, 4, colRed)
	if st.GameOverChar != nil {
		c := st.GameOverChar
		r.tc(screen, fmt.Sprintf("%s  Floor %d  Level %d  Gold %d", c.Name, st.GameOverFloor, c.Level, c.Gold), cx, 148, 1, colWhite)
	}
	for i, item := range st.GameOverItems {
		c := colWhite
		if i == st.GameOverCursor {
			c = colSel
		}
		r.tc(screen, item, cx, float64(195+i*32), 2, c)
	}
}

// ── WIN ───────────────────────────────────────────────────────────────────────

func (r *Renderer) drawWin(screen *ebiten.Image, st *EngineState) {
	cx := float64(ScreenW) / 2
	t := st.TimeSec
	for y := 0; y < ScreenH; y++ {
		for x := 0; x < ScreenW; x++ {
			dx := float64(x-ScreenW/2) / float64(ScreenW/2)
			dy := float64(y-ScreenH/2) / float64(ScreenH/2)
			d := math.Sqrt(dx*dx + dy*dy)
			pulse := 0.08 * math.Sin(t*2.0+d*3.0)
			v := uint8(fclamp((0.18-d*0.12+pulse)*255, 0, 65))
			screen.Set(x, y, color.RGBA{v, uint8(float64(v) * 0.7), 0, 255})
		}
	}
	r.tc(screen, "VICTORY!", cx, 55, 4, colTextGold)
	r.tc(screen, "The Void Sovereign has fallen.", cx, 110, 1, colCyan)
	if st.WinChar != nil {
		c := st.WinChar
		r.tc(screen, fmt.Sprintf("%s  Lv%d  Gold %d  Perks %d", c.Name, c.Level, c.Gold, len(c.Perks)), cx, 136, 1, colWhite)
	}
	for i, item := range st.WinItems {
		c := colWhite
		if i == st.WinCursor {
			c = colSel
		}
		r.tc(screen, item, cx, float64(185+i*32), 2, c)
	}
}

// ── DRAWING PRIMITIVES ───────────────────────────────────────────────────────

func (r *Renderer) drawTile(dst *ebiten.Image, id, x, y int, scale float64) {
	if id < 0 || id >= assets.TileCount {
		id = 0
	}
	sx := id * assets.TileSize
	sub := r.tileset.SubImage(image.Rect(sx, 0, sx+assets.TileSize, assets.TileSize)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(sub, op)
}

func (r *Renderer) drawSprite(dst *ebiten.Image, id int, x, y, scale float64) {
	switch id {
	case assets.SprGolem:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		dst.DrawImage(r.golem, op)
		return
	case assets.SprVoidSovereign:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		dst.DrawImage(r.void, op)
		return
	}
	if id < 0 || id >= assets.SprCount {
		id = 0
	}
	sx := id * assets.TileSize
	sub := r.sprites.SubImage(image.Rect(sx, 0, sx+assets.TileSize, assets.TileSize)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	dst.DrawImage(sub, op)
}

// t draws text at (x,y) at given scale.
func (r *Renderer) t(dst *ebiten.Image, s string, x, y float64, scale int, c color.RGBA) {
	op := &etext.DrawOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	etext.Draw(dst, s, r.face, op)
}

// tc draws text centered on cx at y.
func (r *Renderer) tc(dst *ebiten.Image, s string, cx, y float64, scale int, c color.RGBA) {
	w, _ := etext.Measure(s, r.face, 0)
	r.t(dst, s, cx-w*float64(scale)/2, y, scale, c)
}

func (r *Renderer) fillRect(dst *ebiten.Image, x, y, w, h int, c color.Color) {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(img, op)
}

func (r *Renderer) border(dst *ebiten.Image, x, y, w, h int, c color.RGBA) {
	for i := 0; i < w; i++ {
		dst.Set(x+i, y, c)
		dst.Set(x+i, y+h-1, c)
	}
	for i := 1; i < h-1; i++ {
		dst.Set(x, y+i, c)
		dst.Set(x+w-1, y+i, c)
	}
}

func (r *Renderer) winBox(dst *ebiten.Image, x, y, w, h int, title string) {
	r.fillRect(dst, x, y, w, h, colWindowBg)
	r.border(dst, x, y, w, h, colGold)
	if title != "" {
		r.tc(dst, " "+title+" ", float64(x+w/2), float64(y+3), 1, colTextGold)
	}
}

func (r *Renderer) drawBar(dst *ebiten.Image, x, y, w, h int, pct float64, fill, bg color.RGBA) {
	r.fillRect(dst, x, y, w, h, bg)
	filled := int(fclamp(pct, 0, 1) * float64(w))
	if filled > 0 {
		r.fillRect(dst, x, y, filled, h, fill)
	}
}

func (r *Renderer) hRule(dst *ebiten.Image, y int, c color.RGBA) {
	for x := 20; x < ScreenW-20; x++ {
		dst.Set(x, y, c)
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func gameTile(tt game.TileType) int {
	switch tt {
	case game.TileFloor:
		return assets.TileFloor
	case game.TileWall:
		return assets.TileWall
	case game.TileCorridorH:
		return assets.TileCorrH
	case game.TileCorridorV:
		return assets.TileCorrV
	case game.TileDoor:
		return assets.TileDoor
	case game.TileStairs:
		return assets.TileStairs
	case game.TileStoryBeat:
		return assets.TileStory
	case game.TileBoss:
		return assets.TileBossRoom
	case game.TileShop:
		return assets.TileShop
	case game.TileHeal:
		return assets.TileHeal
	}
	return assets.TileFloor
}

func enemySpr(name string) int {
	switch strings.ToLower(name) {
	case "goblin":
		return assets.SprGoblin
	case "skeleton":
		return assets.SprSkeleton
	case "slime":
		return assets.SprSlime
	case "shadow bat":
		return assets.SprBat
	case "stone colossus", "floor guardian", "stone golem":
		return assets.SprGolem
	case "void sovereign":
		return assets.SprVoidSovereign
	}
	return assets.SprGoblin
}

func perkCol(pt game.PerkType) color.RGBA {
	switch pt {
	case game.CombatPerk:
		return colRed
	case game.MagicPerk:
		return colMPFill
	case game.UtilityPerk:
		return colGreen
	default:
		return colTextGold
	}
}

func wrapText(s string, maxW int) []string {
	if maxW <= 0 {
		maxW = 40
	}
	words := strings.Fields(s)
	var lines []string
	line := ""
	for _, w := range words {
		if len(line)+len(w)+1 > maxW {
			if line != "" {
				lines = append(lines, line)
			}
			line = w
		} else {
			if line != "" {
				line += " "
			}
			line += w
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func fclamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func iclamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
