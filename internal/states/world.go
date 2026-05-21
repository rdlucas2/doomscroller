package states

import (
	"fmt"
	"github.com/rdlucas2/doomscroller/internal/game"
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
	Session     *game.Session
	StoryQueue  []*game.StoryBeat
	Message     string
	messageTick int
}

func NewWorldState(s *game.Session) *WorldState { return &WorldState{Session: s} }

func (ws *WorldState) SetMessage(msg string) {
	ws.Message = msg
	ws.messageTick = 240 // ~4 seconds at 60fps
}

func (ws *WorldState) TickMessage() {
	if ws.messageTick > 0 {
		ws.messageTick--
		if ws.messageTick == 0 {
			ws.Message = ""
		}
	}
}

func (ws *WorldState) Move(dx, dy int) WorldEvent {
	_, _, moved := ws.Session.Map.MovePlayer(dx, dy)
	if !moved {
		return WorldEventNone
	}
	ws.Session.StepCount++
	ws.TickMessage()

	m := ws.Session.Map
	tile := m.At(m.PlayerX, m.PlayerY)
	if tile == nil {
		return WorldEventNone
	}

	switch tile.Type {
	case game.TileStoryBeat:
		storyID := tile.StoryID
		tile.Type = game.TileFloor
		for _, s := range game.GetStoriesForFloor(ws.Session.Floor) {
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
		ws.SetMessage("Hidden cache! You found 40 gold.")
		return WorldEventNone
	case game.TileHeal:
		tile.Type = game.TileFloor
		heal := ws.Session.Character.MaxHP / 4
		ws.Session.Character.HP += heal
		if ws.Session.Character.HP > ws.Session.Character.MaxHP {
			ws.Session.Character.HP = ws.Session.Character.MaxHP
		}
		ws.SetMessage(fmt.Sprintf("Healing spring! Restored %d HP.", heal))
		return WorldEventNone
	case game.TileStairs:
		if ws.Session.Floor >= game.MaxFloors {
			return WorldEventWin
		}
		return WorldEventNextFloor
	}

	if ws.Session.StepCount%4 == 0 {
		if game.RandomEncounterChance(ws.Session.Floor, ws.Session.Rng) {
			return WorldEventRandomBattle
		}
	}
	return WorldEventNone
}
