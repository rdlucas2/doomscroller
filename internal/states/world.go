package states

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/game"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

const (
	viewportW = 30
	viewportH = 16
)

type WorldEvent int

const (
	WorldEventNone WorldEvent = iota
	WorldEventRandomBattle
	WorldEventStoryBeat
	WorldEventBoss
	WorldEventShop
	WorldEventHeal
	WorldEventNextFloor
	WorldEventWin
)

type WorldState struct {
	Session    *game.Session
	EventQueue []WorldEvent
	StoryQueue []*game.StoryBeat
	Message    string
	MessageTick int
}

func NewWorldState(session *game.Session) *WorldState {
	return &WorldState{Session: session}
}

func (ws *WorldState) Move(dx, dy int) WorldEvent {
	_, _, moved := ws.Session.Map.MovePlayer(dx, dy)
	if !moved {
		return WorldEventNone
	}

	ws.Session.StepCount++
	m := ws.Session.Map
	tile := m.At(m.PlayerX, m.PlayerY)
	if tile == nil {
		return WorldEventNone
	}

	switch tile.Type {
	case game.TileStoryBeat:
		storyID := tile.StoryID
		tile.Type = game.TileFloor // consume
		stories := game.GetStoriesForFloor(ws.Session.Floor)
		for _, s := range stories {
			if s.ID == storyID && !s.Completed {
				ws.StoryQueue = append(ws.StoryQueue, s)
				return WorldEventStoryBeat
			}
		}
	case game.TileBoss:
		tile.Type = game.TileFloor
		return WorldEventBoss
	case game.TileShop:
		tile.Type = game.TileFloor
		ws.Session.Character.Gold += 40
		ws.SetMessage("A hidden cache! You found 40 gold.")
		return WorldEventNone
	case game.TileHeal:
		tile.Type = game.TileFloor
		heal := ws.Session.Character.MaxHP / 4
		ws.Session.Character.HP += heal
		if ws.Session.Character.HP > ws.Session.Character.MaxHP {
			ws.Session.Character.HP = ws.Session.Character.MaxHP
		}
		ws.SetMessage(fmt.Sprintf("A healing spring! Restored %d HP.", heal))
		return WorldEventNone
	case game.TileStairs:
		if ws.Session.Floor >= game.MaxFloors {
			return WorldEventWin
		}
		return WorldEventNextFloor
	}

	// Random encounter check every few steps
	if ws.Session.StepCount%4 == 0 {
		if game.RandomEncounterChance(ws.Session.Floor, ws.Session.Rng) {
			return WorldEventRandomBattle
		}
	}

	return WorldEventNone
}

func (ws *WorldState) SetMessage(msg string) {
	ws.Message = msg
	ws.MessageTick = 40
}

func (ws *WorldState) TickMessage() {
	if ws.MessageTick > 0 {
		ws.MessageTick--
		if ws.MessageTick == 0 {
			ws.Message = ""
		}
	}
}

func (ws *WorldState) View(width, height int) string {
	s := ws.Session
	c := s.Character
	m := s.Map

	// --- Map viewport ---
	camX := m.PlayerX - viewportW/2
	camY := m.PlayerY - viewportH/2
	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX+viewportW > m.Width {
		camX = m.Width - viewportW
	}
	if camY+viewportH > m.Height {
		camY = m.Height - viewportH
	}

	mapLines := make([]string, 0, viewportH)
	for vy := 0; vy < viewportH; vy++ {
		mapY := camY + vy
		row := ""
		for vx := 0; vx < viewportW; vx++ {
			mapX := camX + vx
			if mapX == m.PlayerX && mapY == m.PlayerY {
				row += ui.SuccessStyle.Render("@")
				continue
			}
			t := m.At(mapX, mapY)
			if t == nil {
				row += " "
				continue
			}
			// Fog of war unless full map is shown
			if !s.Settings.ShowFullMap && !t.Visited && !isNearPlayer(m, mapX, mapY, 4) {
				row += ui.DimStyle.Render("·")
				continue
			}
			switch t.Type {
			case game.TileWall:
				row += ui.DimStyle.Render("#")
			case game.TileFloor:
				row += ui.DimStyle.Render(".")
			case game.TileCorridorH:
				row += ui.DimStyle.Render("-")
			case game.TileCorridorV:
				row += ui.DimStyle.Render("|")
			case game.TileStairs:
				row += ui.GoldStyle.Render(">")
			case game.TileStoryBeat:
				row += ui.XPStyle.Render("!")
			case game.TileBoss:
				row += ui.HPStyle.Render("B")
			case game.TileShop:
				row += ui.GoldStyle.Render("$")
			case game.TileHeal:
				row += ui.SuccessStyle.Render("+")
			default:
				row += " "
			}
		}
		mapLines = append(mapLines, row)
	}

	mapContent := strings.Join(mapLines, "\n")
	mapBox := ui.BoxStyle.Width(viewportW + 2).Render(mapContent)

	// --- Status panel ---
	floorLabel := ui.TitleStyle.Render(fmt.Sprintf("Floor %d / %d", s.Floor, game.MaxFloors))
	nameLabel := fmt.Sprintf("%s  [%s  Lv.%d]", ui.TitleStyle.Render(c.Name), ui.GoldStyle.Render(c.Class.String()), c.Level)

	hpBar := fmt.Sprintf("HP %s %s",
		ui.HPBar(c.HP, c.MaxHP, 12),
		ui.HPStyle.Render(fmt.Sprintf("%d/%d", c.HP, c.MaxHP)))
	mpBar := fmt.Sprintf("MP %s %s",
		ui.MPBar(c.MP, c.MaxMP, 12),
		ui.MPStyle.Render(fmt.Sprintf("%d/%d", c.MP, c.MaxMP)))
	xpBar := fmt.Sprintf("XP %s %s",
		ui.XPBar(c.XP, c.XPToNext, 12),
		ui.XPStyle.Render(fmt.Sprintf("%d/%d", c.XP, c.XPToNext)))
	goldLine := ui.GoldStyle.Render(fmt.Sprintf("Gold: %d", c.Gold))

	perksHeader := ui.SubtitleStyle.Render("Perks:")
	perkLines := []string{}
	for _, p := range c.Perks {
		perkLines = append(perkLines, "  "+perkTypeStyle(p.Type).Render("•"+p.Name))
	}
	if len(perkLines) == 0 {
		perkLines = []string{ui.DimStyle.Render("  None")}
	}

	legendLine := ui.DimStyle.Render("@ You  > Stairs  ! Story  B Boss  $ Gold  + Heal")

	statusContent := strings.Join(append([]string{
		floorLabel, nameLabel, "", hpBar, mpBar, xpBar, goldLine, "", perksHeader,
	}, append(perkLines, "", legendLine)...), "\n")

	statusBox := ui.BoxStyle.Width(30).Render(statusContent)

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, mapBox, "  ", statusBox)

	// --- Bottom bar ---
	var bottomLines []string
	if ws.Message != "" {
		bottomLines = append(bottomLines, ui.SuccessStyle.Render("▶ "+ws.Message))
	}
	bottomLines = append(bottomLines,
		ui.DimStyle.Render("WASD/Arrows Move   P Pause   Q Quit to menu"),
	)
	bottom := strings.Join(bottomLines, "\n")

	return lipgloss.JoinVertical(lipgloss.Left, mainRow, "", bottom)
}

func isNearPlayer(m *game.WorldMap, x, y, radius int) bool {
	dx := x - m.PlayerX
	dy := y - m.PlayerY
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx <= radius && dy <= radius
}

func perkTypeStyle(pt game.PerkType) lipgloss.Style {
	switch pt {
	case game.CombatPerk:
		return ui.PerkCombatStyle
	case game.MagicPerk:
		return ui.PerkMagicStyle
	case game.UtilityPerk:
		return ui.PerkUtilStyle
	default:
		return ui.PerkPassiveStyle
	}
}
