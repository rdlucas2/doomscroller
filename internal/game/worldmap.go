package game

import (
	"math/rand"
)

type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileCorridorH
	TileCorridorV
	TileDoor
	TileStairs
	TileStoryBeat
	TileBoss
	TileShop
	TileHeal
	TilePlayer
)

func (t TileType) Rune() rune {
	switch t {
	case TileWall:
		return '#'
	case TileFloor:
		return '.'
	case TileCorridorH:
		return '-'
	case TileCorridorV:
		return '|'
	case TileDoor:
		return '+'
	case TileStairs:
		return '>'
	case TileStoryBeat:
		return '!'
	case TileBoss:
		return 'B'
	case TileShop:
		return '$'
	case TileHeal:
		return 'H'
	case TilePlayer:
		return '@'
	}
	return ' '
}

type Tile struct {
	Type    TileType
	Visited bool
	StoryID string
	BossName string
}

type Room struct {
	X, Y, W, H int
}

func (r Room) CenterX() int { return r.X + r.W/2 }
func (r Room) CenterY() int { return r.Y + r.H/2 }
func (r Room) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

type WorldMap struct {
	Width   int
	Height  int
	Tiles   [][]Tile
	Rooms   []Room
	Floor   int
	PlayerX int
	PlayerY int
	// track which story beat tiles are placed
	StoryTiles map[string][2]int
}

func (m *WorldMap) At(x, y int) *Tile {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return nil
	}
	return &m.Tiles[y][x]
}

func (m *WorldMap) Set(x, y int, tt TileType) {
	if x >= 0 && y >= 0 && x < m.Width && y < m.Height {
		m.Tiles[y][x].Type = tt
	}
}

func GenerateMap(floor int, rng *rand.Rand) *WorldMap {
	w, h := 60, 30
	wm := &WorldMap{
		Width:      w,
		Height:     h,
		Floor:      floor,
		StoryTiles: map[string][2]int{},
	}
	wm.Tiles = make([][]Tile, h)
	for y := range wm.Tiles {
		wm.Tiles[y] = make([]Tile, w)
		for x := range wm.Tiles[y] {
			wm.Tiles[y][x] = Tile{Type: TileWall}
		}
	}

	// BSP-style room placement
	minRooms, maxRooms := 5+floor, 8+floor
	if maxRooms > 14 {
		maxRooms = 14
	}
	numRooms := minRooms + rng.Intn(maxRooms-minRooms+1)

	for i := 0; i < numRooms*10 && len(wm.Rooms) < numRooms; i++ {
		rw := 4 + rng.Intn(7)
		rh := 3 + rng.Intn(5)
		rx := 1 + rng.Intn(w-rw-2)
		ry := 1 + rng.Intn(h-rh-2)
		r := Room{X: rx, Y: ry, W: rw, H: rh}
		if !overlaps(wm.Rooms, r) {
			wm.Rooms = append(wm.Rooms, r)
			carveRoom(wm, r)
		}
	}

	// Connect rooms with corridors
	for i := 1; i < len(wm.Rooms); i++ {
		prev := wm.Rooms[i-1]
		cur := wm.Rooms[i]
		carveCorridor(wm, prev.CenterX(), prev.CenterY(), cur.CenterX(), cur.CenterY())
	}

	// Place player in first room
	if len(wm.Rooms) > 0 {
		wm.PlayerX = wm.Rooms[0].CenterX()
		wm.PlayerY = wm.Rooms[0].CenterY()
	}

	// Place stairs in last room
	if len(wm.Rooms) > 1 {
		last := wm.Rooms[len(wm.Rooms)-1]
		wm.Set(last.CenterX(), last.CenterY(), TileStairs)
	}

	// Place boss in second-to-last room (if available)
	stories := GetStoriesForFloor(floor)
	if len(wm.Rooms) > 2 {
		bossTile := wm.Rooms[len(wm.Rooms)-2]
		bossName := "Floor Guardian"
		storyID := ""
		for _, s := range stories {
			if s.BossName != "" && s.Trigger == TriggerBossDefeated {
				bossName = s.BossName
				storyID = s.ID
				break
			}
		}
		bx, by := bossTile.CenterX(), bossTile.CenterY()
		wm.Set(bx, by, TileBoss)
		wm.At(bx, by).BossName = bossName
		wm.At(bx, by).StoryID = storyID
		if storyID != "" {
			wm.StoryTiles[storyID] = [2]int{bx, by}
		}
	}

	// Scatter story beat, shop, heal tiles in mid rooms
	midRooms := wm.Rooms[1:]
	if len(midRooms) > 2 {
		midRooms = midRooms[:len(midRooms)-2]
	}

	storyIdx := 0
	for _, s := range stories {
		if s.Trigger == TriggerEnterRoom && storyIdx < len(midRooms) {
			r := midRooms[storyIdx]
			sx, sy := r.CenterX(), r.CenterY()
			wm.Set(sx, sy, TileStoryBeat)
			wm.At(sx, sy).StoryID = s.ID
			wm.StoryTiles[s.ID] = [2]int{sx, sy}
			storyIdx++
		}
	}

	// Scatter shop and heal in random mid rooms
	for _, r := range midRooms {
		if rng.Intn(4) == 0 {
			wm.Set(r.CenterX(), r.CenterY(), TileShop)
		} else if rng.Intn(4) == 0 {
			wm.Set(r.CenterX(), r.CenterY(), TileHeal)
		}
	}

	return wm
}

func overlaps(rooms []Room, r Room) bool {
	for _, existing := range rooms {
		if r.X < existing.X+existing.W+2 && r.X+r.W+2 > existing.X &&
			r.Y < existing.Y+existing.H+2 && r.Y+r.H+2 > existing.Y {
			return true
		}
	}
	return false
}

func carveRoom(m *WorldMap, r Room) {
	for y := r.Y; y < r.Y+r.H; y++ {
		for x := r.X; x < r.X+r.W; x++ {
			m.Set(x, y, TileFloor)
		}
	}
}

func carveCorridor(m *WorldMap, x1, y1, x2, y2 int) {
	x, y := x1, y1
	// Horizontal first, then vertical
	for x != x2 {
		if m.At(x, y) != nil && m.At(x, y).Type == TileWall {
			m.Set(x, y, TileCorridorH)
		}
		if x < x2 {
			x++
		} else {
			x--
		}
	}
	for y != y2 {
		if m.At(x, y) != nil && m.At(x, y).Type == TileWall {
			m.Set(x, y, TileCorridorV)
		}
		if y < y2 {
			y++
		} else {
			y--
		}
	}
}

func (m *WorldMap) CanMoveTo(x, y int) bool {
	t := m.At(x, y)
	if t == nil {
		return false
	}
	return t.Type != TileWall
}

func (m *WorldMap) MovePlayer(dx, dy int) (newX, newY int, moved bool) {
	nx, ny := m.PlayerX+dx, m.PlayerY+dy
	if m.CanMoveTo(nx, ny) {
		m.PlayerX, m.PlayerY = nx, ny
		if t := m.At(nx, ny); t != nil {
			t.Visited = true
		}
		return nx, ny, true
	}
	return m.PlayerX, m.PlayerY, false
}

// RandomEncounterChance returns true based on floor difficulty and random roll.
func RandomEncounterChance(floor int, rng *rand.Rand) bool {
	threshold := 15 + floor*2
	return rng.Intn(100) < threshold
}
