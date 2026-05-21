// Package assets generates all pixel-art sprite and tile images procedurally.
// Every image is a 16x16 (or larger) NRGBA created at startup with no external files.
// To replace with real assets, load PNGs instead of calling the gen* functions.
package assets

import (
	"image"
	"image/color"
)

// Tile IDs
const (
	TileFloor = iota
	TileWall
	TileDoor
	TileStairs
	TileStory
	TileBossRoom
	TileShop
	TileHeal
	TileCorrH
	TileCorrV
	TileCount
)

// Sprite IDs
const (
	SprWarrior1 = iota
	SprWarrior2
	SprMage1
	SprMage2
	SprRogue1
	SprRogue2
	SprGoblin
	SprSkeleton
	SprSlime
	SprBat
	SprGolem
	SprVoidSovereign
	SprCount
)

const TileSize = 16

// ── helpers ──────────────────────────────────────────────────────────────────

func nrgba(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 255} }
func trans() color.NRGBA               { return color.NRGBA{0, 0, 0, 0} }

func newImg(w, h int) *image.NRGBA {
	return image.NewNRGBA(image.Rect(0, 0, w, h))
}

func fill(img *image.NRGBA, c color.NRGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
}

func fillRect(img *image.NRGBA, x, y, w, h int, c color.NRGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.SetNRGBA(x+dx, y+dy, c)
		}
	}
}

func hLine(img *image.NRGBA, y, x, w int, c color.NRGBA) {
	for i := 0; i < w; i++ {
		img.SetNRGBA(x+i, y, c)
	}
}

func vLine(img *image.NRGBA, x, y, h int, c color.NRGBA) {
	for i := 0; i < h; i++ {
		img.SetNRGBA(x, y+i, c)
	}
}

func set(img *image.NRGBA, x, y int, c color.NRGBA) {
	b := img.Bounds()
	if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
		img.SetNRGBA(x, y, c)
	}
}

// drawSprite paints a palette-indexed 16x16 sprite onto img at offset (ox,oy).
// Index 0 = transparent.
func drawSprite(img *image.NRGBA, data [16 * 16]uint8, palette []color.NRGBA, ox, oy int) {
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			idx := data[y*16+x]
			if idx == 0 {
				continue
			}
			if int(idx) < len(palette) {
				img.SetNRGBA(ox+x, oy+y, palette[idx])
			}
		}
	}
}

// ── TILE GENERATORS ──────────────────────────────────────────────────────────

func GenTileFloor() *image.NRGBA {
	img := newImg(16, 16)
	grout := nrgba(35, 38, 42)
	stone := nrgba(72, 79, 86)
	hi := nrgba(98, 108, 116)
	sh := nrgba(52, 57, 62)

	fill(img, stone)
	// grout grid
	hLine(img, 0, 0, 16, grout)
	hLine(img, 8, 0, 16, grout)
	hLine(img, 15, 0, 16, grout)
	vLine(img, 0, 0, 16, grout)
	vLine(img, 8, 0, 16, grout)
	vLine(img, 15, 0, 16, grout)
	// highlights (top-left of each stone)
	hLine(img, 1, 1, 6, hi)
	vLine(img, 1, 1, 6, hi)
	hLine(img, 1, 9, 6, hi)
	vLine(img, 9, 1, 6, hi)
	hLine(img, 9, 1, 6, hi)
	vLine(img, 1, 9, 6, hi)
	hLine(img, 9, 9, 6, hi)
	vLine(img, 9, 9, 6, hi)
	// shadows (bottom-right of each stone)
	hLine(img, 7, 1, 6, sh)
	vLine(img, 7, 1, 6, sh)
	hLine(img, 7, 9, 6, sh)
	vLine(img, 14, 1, 6, sh)
	hLine(img, 14, 1, 6, sh)
	vLine(img, 7, 9, 6, sh)
	hLine(img, 14, 9, 6, sh)
	vLine(img, 14, 9, 6, sh)
	return img
}

func GenTileWall() *image.NRGBA {
	img := newImg(16, 16)
	mortar := nrgba(30, 24, 20)
	brick := nrgba(80, 48, 36)
	brickHi := nrgba(104, 64, 48)
	brickSh := nrgba(55, 33, 24)
	moss := nrgba(45, 62, 40)

	fill(img, mortar)
	// Four rows of bricks; odd rows offset by 4
	offsets := []int{0, 4, 0, 4}
	for row := 0; row < 4; row++ {
		y0 := row * 4
		y1 := y0 + 3
		off := offsets[row]
		// Draw 3 bricks per row (they may wrap)
		for b := -1; b < 3; b++ {
			bx := b*6 + off
			if bx+5 < 0 || bx > 15 {
				continue
			}
			xS := bx
			xE := bx + 4
			if xS < 0 {
				xS = 0
			}
			if xE > 15 {
				xE = 15
			}
			fillRect(img, xS, y0, xE-xS, 3, brick)
			hLine(img, y0, xS, xE-xS, brickHi)
			vLine(img, xS, y0, 3, brickHi)
			hLine(img, y1, xS, xE-xS, brickSh)
		}
	}
	// moss patches (atmospheric detail)
	set(img, 2, 5, moss)
	set(img, 3, 5, moss)
	set(img, 11, 9, moss)
	set(img, 12, 9, moss)
	return img
}

func GenTileDoor() *image.NRGBA {
	img := newImg(16, 16)
	fill(img, nrgba(35, 38, 42))
	frame := nrgba(58, 40, 26)
	wood := nrgba(100, 68, 42)
	woodHi := nrgba(130, 90, 58)
	grain := nrgba(80, 54, 34)
	handle := nrgba(200, 170, 60)

	// Door frame
	fillRect(img, 1, 0, 14, 16, frame)
	// Door panels
	fillRect(img, 2, 1, 12, 14, wood)
	// Wood grain lines
	vLine(img, 7, 1, 14, grain)
	vLine(img, 8, 1, 14, grain)
	hLine(img, 7, 2, 12, grain)
	// Highlight left edge
	vLine(img, 2, 1, 14, woodHi)
	hLine(img, 1, 2, 12, woodHi)
	// Handle
	fillRect(img, 11, 7, 2, 3, handle)
	return img
}

func GenTileStairs() *image.NRGBA {
	img := newImg(16, 16)
	base := nrgba(60, 65, 70)
	step := nrgba(85, 92, 98)
	sh := nrgba(40, 44, 48)
	hi := nrgba(110, 118, 126)
	arrow := nrgba(200, 180, 80)

	fill(img, base)
	// 4 step rows (going down visually = top is back, bottom is front)
	stepH := []int{2, 5, 9, 12}
	for i, y := range stepH {
		w := 4 + i*3
		x := (16 - w) / 2
		fillRect(img, x, y, w, 2, step)
		hLine(img, y, x, w, hi)
		hLine(img, y+2, x, w, sh)
	}
	// Down arrow
	set(img, 7, 7, arrow)
	set(img, 8, 7, arrow)
	set(img, 6, 8, arrow)
	set(img, 9, 8, arrow)
	set(img, 7, 9, arrow)
	set(img, 8, 9, arrow)
	return img
}

func GenTileStory() *image.NRGBA {
	img := newImg(16, 16)
	// stone floor base
	drawImgOver(img, GenTileFloor())
	// glowing rune overlay
	glow := nrgba(80, 200, 255)
	dim := nrgba(40, 120, 180)
	// rune shape (simplified diamond + lines)
	set(img, 8, 4, glow)
	hLine(img, 5, 6, 5, glow)
	hLine(img, 6, 5, 7, glow)
	hLine(img, 7, 4, 9, dim)
	hLine(img, 8, 5, 7, glow)
	hLine(img, 9, 6, 5, glow)
	set(img, 8, 10, glow)
	vLine(img, 8, 4, 8, dim)
	return img
}

func GenTileBossRoom() *image.NRGBA {
	img := newImg(16, 16)
	fill(img, nrgba(28, 14, 14))
	bone := nrgba(210, 196, 172)
	boneSh := nrgba(148, 136, 118)
	red := nrgba(180, 30, 30)
	// skull outline
	fillRect(img, 4, 2, 8, 6, bone)
	fillRect(img, 5, 8, 2, 3, bone)
	fillRect(img, 9, 8, 2, 3, bone)
	// eyes (red glow)
	fillRect(img, 5, 3, 2, 2, red)
	fillRect(img, 9, 3, 2, 2, red)
	// nose
	set(img, 7, 6, boneSh)
	set(img, 8, 6, boneSh)
	// teeth
	for i := 0; i < 3; i++ {
		set(img, 5+i*2, 8, boneSh)
	}
	// dark shadow
	hLine(img, 7, 4, 8, boneSh)
	return img
}

func GenTileShop() *image.NRGBA {
	img := newImg(16, 16)
	drawImgOver(img, GenTileFloor())
	gold := nrgba(210, 175, 50)
	goldHi := nrgba(240, 210, 90)
	goldSh := nrgba(155, 128, 35)
	// coin stack
	for i, y := range []int{11, 9, 7} {
		r := 3 - i/2
		cx := 8
		fillRect(img, cx-r, y, r*2, 2, gold)
		hLine(img, y, cx-r, r*2, goldHi)
		hLine(img, y+1, cx-r, r*2, goldSh)
	}
	// G symbol
	set(img, 8, 4, goldHi)
	set(img, 7, 4, goldHi)
	set(img, 9, 5, goldHi)
	set(img, 8, 5, goldHi)
	set(img, 6, 5, goldHi)
	set(img, 6, 6, goldHi)
	set(img, 9, 6, goldHi)
	return img
}

func GenTileHeal() *image.NRGBA {
	img := newImg(16, 16)
	stone := nrgba(68, 74, 80)
	rim := nrgba(90, 98, 106)
	water := nrgba(40, 120, 200)
	waterHi := nrgba(80, 170, 240)
	ripple := nrgba(120, 195, 255)
	cross := nrgba(220, 60, 60)

	fill(img, stone)
	// Basin rim
	fillRect(img, 2, 4, 12, 9, rim)
	// Water inside basin
	fillRect(img, 3, 5, 10, 7, water)
	hLine(img, 5, 3, 10, waterHi)
	// Ripples
	set(img, 6, 8, ripple)
	set(img, 9, 7, ripple)
	set(img, 11, 9, ripple)
	// Heal cross
	vLine(img, 8, 1, 3, cross)
	hLine(img, 2, 7, 3, cross)
	return img
}

func GenTileCorrH() *image.NRGBA {
	img := newImg(16, 16)
	fill(img, nrgba(35, 38, 42))
	c := nrgba(65, 71, 78)
	hi := nrgba(85, 93, 100)
	fillRect(img, 0, 5, 16, 6, c)
	hLine(img, 5, 0, 16, hi)
	hLine(img, 10, 0, 16, nrgba(48, 52, 58))
	return img
}

func GenTileCorrV() *image.NRGBA {
	img := newImg(16, 16)
	fill(img, nrgba(35, 38, 42))
	c := nrgba(65, 71, 78)
	hi := nrgba(85, 93, 100)
	fillRect(img, 5, 0, 6, 16, c)
	vLine(img, 5, 0, 16, hi)
	vLine(img, 10, 0, 16, nrgba(48, 52, 58))
	return img
}

func drawImgOver(dst, src *image.NRGBA) {
	sb := src.Bounds()
	for y := sb.Min.Y; y < sb.Max.Y; y++ {
		for x := sb.Min.X; x < sb.Max.X; x++ {
			c := src.NRGBAAt(x, y)
			if c.A > 0 {
				dst.SetNRGBA(x, y, c)
			}
		}
	}
}

// GenerateTileset returns a horizontal strip of all TileCount tiles.
func GenerateTileset() *image.NRGBA {
	img := newImg(TileSize*TileCount, TileSize)
	tiles := [TileCount]*image.NRGBA{
		GenTileFloor(),
		GenTileWall(),
		GenTileDoor(),
		GenTileStairs(),
		GenTileStory(),
		GenTileBossRoom(),
		GenTileShop(),
		GenTileHeal(),
		GenTileCorrH(),
		GenTileCorrV(),
	}
	for i, t := range tiles {
		for y := 0; y < TileSize; y++ {
			for x := 0; x < TileSize; x++ {
				img.SetNRGBA(i*TileSize+x, y, t.NRGBAAt(x, y))
			}
		}
	}
	return img
}

// ── CHARACTER SPRITES ─────────────────────────────────────────────────────────

// pallWarrior defines warrior colors
var pallWarrior = []color.NRGBA{
	{0, 0, 0, 0},       // 0 transparent
	{180, 185, 190, 255}, // 1 silver armor
	{230, 185, 155, 255}, // 2 skin
	{28, 58, 195, 255},   // 3 blue cape
	{155, 160, 165, 255}, // 4 chainmail dark
	{75, 52, 32, 255},    // 5 brown boot
	{45, 45, 45, 255},    // 6 outline
	{210, 215, 220, 255}, // 7 silver hi
	{210, 168, 138, 255}, // 8 skin shadow
	{100, 140, 210, 255}, // 9 cape hi
}

var sprWarrior1 = [16 * 16]uint8{
	0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 1, 7, 7, 7, 1, 1, 6, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 2, 2, 2, 2, 2, 1, 6, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 2, 8, 2, 2, 2, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 4, 4, 4, 4, 4, 4, 4, 1, 6, 0, 0, 0,
	0, 0, 6, 4, 4, 6, 4, 4, 4, 6, 4, 4, 6, 0, 0, 0,
	0, 0, 6, 4, 4, 6, 4, 4, 4, 6, 4, 4, 6, 0, 0, 0,
	0, 3, 6, 3, 4, 4, 4, 4, 4, 4, 4, 3, 6, 3, 0, 0,
	3, 3, 6, 3, 3, 3, 3, 3, 3, 3, 3, 3, 6, 3, 3, 0,
	0, 3, 9, 3, 3, 3, 3, 3, 3, 3, 3, 3, 9, 3, 0, 0,
	0, 0, 5, 5, 5, 5, 0, 0, 0, 5, 5, 5, 5, 0, 0, 0,
	0, 0, 5, 5, 5, 5, 0, 0, 0, 5, 5, 5, 5, 0, 0, 0,
	0, 0, 6, 5, 5, 5, 0, 0, 0, 5, 5, 5, 6, 0, 0, 0,
	0, 0, 5, 5, 5, 6, 0, 0, 0, 6, 5, 5, 5, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var sprWarrior2 = [16 * 16]uint8{ // walk frame
	0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 1, 7, 7, 7, 1, 1, 6, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 2, 2, 2, 2, 2, 1, 6, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 2, 8, 2, 2, 2, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 4, 4, 4, 4, 4, 4, 4, 1, 6, 0, 0, 0,
	0, 0, 6, 4, 4, 6, 4, 4, 4, 6, 4, 4, 6, 0, 0, 0,
	0, 0, 6, 4, 4, 6, 4, 4, 4, 6, 4, 4, 6, 0, 0, 0,
	0, 3, 6, 3, 4, 4, 4, 4, 4, 4, 4, 3, 6, 3, 0, 0,
	3, 3, 6, 3, 3, 3, 3, 3, 3, 3, 3, 3, 6, 3, 3, 0,
	0, 3, 9, 3, 3, 3, 3, 3, 3, 3, 3, 3, 9, 3, 0, 0,
	0, 5, 5, 5, 5, 5, 0, 0, 0, 5, 5, 5, 5, 0, 0, 0, // left foot fwd
	0, 5, 5, 5, 5, 5, 0, 0, 0, 5, 5, 5, 5, 0, 0, 0,
	0, 6, 5, 5, 5, 5, 0, 0, 0, 5, 5, 5, 6, 0, 0, 0,
	0, 5, 5, 5, 5, 6, 0, 0, 0, 6, 5, 5, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var pallMage = []color.NRGBA{
	{0, 0, 0, 0},       // 0 transparent
	{80, 28, 120, 255},  // 1 purple robe
	{230, 185, 155, 255}, // 2 skin
	{112, 38, 160, 255},  // 3 robe hi
	{50, 18, 80, 255},    // 4 robe shadow
	{130, 100, 50, 255},  // 5 staff brown
	{45, 45, 45, 255},    // 6 outline
	{155, 90, 200, 255},  // 7 hat mid
	{60, 20, 100, 255},   // 8 hat dark
	{80, 200, 255, 255},  // 9 crystal glow
}

var sprMage1 = [16 * 16]uint8{
	0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 8, 8, 8, 8, 8, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 8, 7, 7, 7, 7, 8, 8, 6, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 2, 2, 2, 2, 2, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 2, 6, 2, 2, 2, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 3, 1, 1, 1, 1, 1, 3, 1, 6, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0,
	0, 0, 6, 4, 1, 1, 1, 1, 1, 1, 1, 4, 6, 0, 0, 0,
	5, 0, 6, 4, 4, 1, 1, 1, 1, 1, 4, 4, 6, 0, 0, 0,
	5, 0, 5, 0, 4, 4, 0, 0, 0, 4, 4, 0, 5, 0, 0, 0,
	5, 0, 5, 0, 4, 4, 0, 0, 0, 4, 4, 0, 0, 0, 0, 0,
	5, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var sprMage2 = [16 * 16]uint8{ // walk frame (staff on right)
	0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 8, 8, 8, 8, 8, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 8, 7, 7, 7, 7, 8, 8, 6, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 2, 2, 2, 2, 2, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 2, 6, 2, 2, 2, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 3, 1, 1, 1, 1, 1, 3, 1, 6, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0,
	0, 0, 6, 4, 1, 1, 1, 1, 1, 1, 1, 4, 6, 5, 0, 0,
	0, 0, 6, 4, 4, 1, 1, 1, 1, 1, 4, 4, 6, 5, 0, 0,
	0, 0, 5, 0, 4, 4, 0, 0, 0, 4, 4, 0, 5, 5, 0, 0,
	0, 0, 5, 0, 4, 4, 0, 0, 0, 4, 4, 0, 5, 9, 0, 0,
	0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var pallRogue = []color.NRGBA{
	{0, 0, 0, 0},        // 0 transparent
	{38, 35, 28, 255},   // 1 dark leather
	{230, 185, 155, 255}, // 2 skin
	{55, 50, 38, 255},   // 3 leather mid
	{22, 20, 16, 255},   // 4 shadow/dark
	{185, 190, 195, 255}, // 5 dagger blade
	{45, 45, 45, 255},   // 6 outline
	{65, 58, 45, 255},   // 7 leather hi
	{210, 168, 138, 255}, // 8 skin shadow
	{140, 144, 148, 255}, // 9 dagger dim
}

var sprRogue1 = [16 * 16]uint8{
	0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 4, 4, 4, 4, 4, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 4, 2, 2, 2, 2, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 4, 2, 8, 2, 2, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 3, 7, 1, 1, 1, 1, 7, 3, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 3, 1, 1, 1, 1, 3, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0,
	5, 6, 1, 4, 1, 1, 1, 1, 1, 4, 1, 6, 5, 0, 0, 0,
	5, 9, 4, 4, 4, 1, 1, 1, 4, 4, 4, 9, 5, 0, 0, 0,
	5, 0, 0, 3, 3, 0, 0, 0, 3, 3, 0, 0, 5, 0, 0, 0,
	0, 0, 0, 3, 3, 0, 0, 0, 3, 3, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 3, 3, 6, 0, 0, 6, 3, 3, 6, 0, 0, 0, 0,
	0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var sprRogue2 = [16 * 16]uint8{ // walk frame
	0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 4, 4, 4, 4, 4, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 4, 2, 2, 2, 2, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 4, 2, 8, 2, 2, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 3, 7, 1, 1, 1, 1, 7, 3, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 3, 1, 1, 1, 1, 3, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0,
	5, 6, 1, 4, 1, 1, 1, 1, 1, 4, 1, 6, 5, 0, 0, 0,
	5, 9, 4, 4, 4, 1, 1, 1, 4, 4, 4, 9, 5, 0, 0, 0,
	5, 0, 3, 3, 0, 0, 0, 0, 0, 3, 3, 0, 5, 0, 0, 0, // legs apart
	0, 3, 3, 0, 0, 0, 0, 0, 0, 0, 3, 3, 0, 0, 0, 0,
	0, 6, 3, 6, 0, 0, 0, 0, 0, 0, 6, 3, 6, 0, 0, 0,
	0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

// ── ENEMY SPRITES ────────────────────────────────────────────────────────────

var pallGoblin = []color.NRGBA{
	{0, 0, 0, 0},        // 0 transparent
	{75, 115, 55, 255},  // 1 green skin
	{50, 80, 35, 255},   // 2 dark green
	{100, 148, 70, 255}, // 3 skin hi
	{240, 210, 30, 255}, // 4 yellow eyes
	{80, 55, 35, 255},   // 5 clothing
	{35, 35, 35, 255},   // 6 outline
	{110, 80, 55, 255},  // 7 clothing hi
	{180, 50, 50, 255},  // 8 weapon red
}

var sprGoblin = [16 * 16]uint8{
	0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 3, 3, 3, 1, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 4, 2, 1, 1, 2, 4, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 2, 1, 1, 2, 2, 1, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0, 0,
	0, 6, 5, 5, 5, 5, 5, 5, 5, 5, 5, 6, 0, 0, 0, 0,
	0, 6, 7, 5, 5, 5, 5, 5, 5, 5, 7, 6, 0, 0, 0, 0,
	8, 6, 5, 5, 5, 5, 5, 5, 5, 5, 5, 6, 0, 0, 0, 0,
	8, 6, 2, 5, 5, 5, 5, 5, 5, 2, 5, 6, 0, 0, 0, 0,
	8, 0, 0, 2, 2, 5, 5, 5, 2, 2, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 2, 2, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 2, 2, 6, 0, 0, 6, 2, 2, 6, 0, 0, 0, 0,
	0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var pallSkeleton = []color.NRGBA{
	{0, 0, 0, 0},        // 0 transparent
	{215, 205, 185, 255}, // 1 bone white
	{155, 145, 128, 255}, // 2 bone shadow
	{240, 230, 210, 255}, // 3 bone highlight
	{55, 130, 220, 255},  // 4 eye glow blue
	{70, 60, 50, 255},    // 5 rusted armor
	{35, 35, 35, 255},    // 6 outline
	{100, 90, 75, 255},   // 7 armor hi
}

var sprSkeleton = [16 * 16]uint8{
	0, 0, 0, 6, 6, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 3, 1, 1, 1, 1, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 4, 1, 1, 1, 4, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 2, 1, 6, 6, 2, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 1, 6, 1, 1, 6, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 6, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 5, 5, 5, 5, 5, 5, 5, 6, 0, 0, 0, 0, 0,
	0, 6, 1, 7, 5, 5, 5, 5, 5, 7, 1, 6, 0, 0, 0, 0,
	0, 6, 1, 5, 5, 5, 5, 5, 5, 5, 1, 6, 0, 0, 0, 0,
	0, 0, 6, 2, 5, 5, 5, 5, 2, 5, 6, 0, 0, 0, 0, 0,
	0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
	0, 0, 1, 2, 1, 0, 0, 1, 2, 1, 0, 0, 0, 0, 0, 0,
	0, 0, 6, 1, 6, 0, 0, 6, 1, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var pallSlime = []color.NRGBA{
	{0, 0, 0, 0},       // 0 transparent
	{80, 195, 70, 255}, // 1 lime green
	{50, 145, 40, 255}, // 2 dark green
	{130, 235, 110, 255}, // 3 highlight
	{40, 40, 40, 255},  // 4 eye
	{155, 240, 140, 255}, // 5 sheen
}

var sprSlime = [16 * 16]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 2, 1, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0,
	0, 0, 0, 2, 1, 3, 5, 1, 1, 1, 1, 2, 0, 0, 0, 0,
	0, 0, 0, 2, 1, 5, 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,
	0, 0, 2, 1, 4, 1, 1, 4, 1, 1, 1, 1, 2, 0, 0, 0,
	0, 0, 2, 1, 4, 4, 1, 4, 4, 1, 1, 1, 2, 0, 0, 0,
	0, 0, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 0, 0, 0,
	0, 0, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 0, 0, 0,
	0, 0, 0, 2, 1, 1, 1, 1, 1, 1, 1, 2, 0, 0, 0, 0,
	0, 0, 0, 2, 2, 1, 1, 1, 1, 2, 2, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var pallBat = []color.NRGBA{
	{0, 0, 0, 0},       // 0 transparent
	{55, 35, 75, 255},  // 1 body purple
	{80, 52, 108, 255}, // 2 wing mid
	{35, 22, 50, 255},  // 3 wing dark
	{220, 40, 40, 255}, // 4 eyes red
	{100, 68, 136, 255}, // 5 wing hi
	{25, 15, 36, 255},  // 6 outline
}

var sprBat = [16 * 16]uint8{
	0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0,
	0, 3, 2, 3, 0, 0, 0, 0, 0, 0, 3, 2, 3, 0, 0, 0,
	3, 2, 5, 2, 3, 0, 6, 6, 6, 0, 3, 5, 2, 3, 0, 0,
	3, 2, 2, 2, 6, 1, 4, 1, 4, 1, 6, 2, 2, 3, 0, 0,
	0, 3, 2, 2, 6, 1, 1, 1, 1, 1, 6, 2, 3, 0, 0, 0,
	0, 0, 3, 2, 6, 1, 6, 1, 6, 1, 6, 3, 0, 0, 0, 0,
	0, 0, 0, 3, 6, 1, 1, 1, 1, 1, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 6, 1, 1, 1, 1, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 6, 6, 6, 6, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

// GenGolem returns a 32x32 boss sprite
func GenGolem() *image.NRGBA {
	img := newImg(32, 32)
	stone := nrgba(108, 114, 120)
	stoneHi := nrgba(145, 153, 160)
	stoneSh := nrgba(68, 72, 78)
	crack := nrgba(44, 48, 52)
	glow := nrgba(220, 110, 30)
	outline := nrgba(30, 30, 30)

	// Body
	fillRect(img, 6, 8, 20, 20, stone)
	// Head
	fillRect(img, 8, 1, 16, 10, stone)
	// Highlights
	hLine(img, 2, 9, 14, stoneHi)
	vLine(img, 9, 2, 6, stoneHi)
	hLine(img, 9, 7, 18, stoneHi)
	vLine(img, 7, 9, 18, stoneHi)
	// Shadows
	hLine(img, 10, 8, 16, stoneSh)
	vLine(img, 23, 2, 6, stoneSh)
	hLine(img, 27, 7, 18, stoneSh)
	vLine(img, 25, 9, 18, stoneSh)
	// Cracks
	vLine(img, 12, 10, 6, crack)
	vLine(img, 19, 12, 8, crack)
	set(img, 13, 15, crack)
	set(img, 14, 16, crack)
	set(img, 18, 20, crack)
	// Glowing eyes
	fillRect(img, 11, 4, 3, 3, glow)
	fillRect(img, 18, 4, 3, 3, glow)
	fillRect(img, 12, 5, 1, 1, nrgba(255, 200, 100))
	fillRect(img, 19, 5, 1, 1, nrgba(255, 200, 100))
	// Fists
	fillRect(img, 2, 16, 5, 6, stone)
	fillRect(img, 25, 16, 5, 6, stone)
	hLine(img, 16, 3, 4, stoneHi)
	hLine(img, 21, 3, 4, stoneSh)
	hLine(img, 16, 26, 4, stoneHi)
	hLine(img, 21, 26, 4, stoneSh)
	// Outline
	for x := 6; x < 26; x++ {
		set(img, x, 7, outline)
		set(img, x, 27, outline)
	}
	for y := 8; y < 28; y++ {
		set(img, 6, y, outline)
		set(img, 25, y, outline)
	}
	return img
}

// GenVoidSovereign returns a 32x40 final boss sprite
func GenVoidSovereign() *image.NRGBA {
	img := newImg(32, 40)
	armor := nrgba(25, 15, 45)
	armorHi := nrgba(48, 30, 80)
	cloak := nrgba(55, 28, 95)
	cloakHi := nrgba(80, 45, 130)
	eyes := nrgba(155, 50, 255)
	eyeGlow := nrgba(200, 120, 255)
	crown := nrgba(180, 140, 20)
	voidEnergy := nrgba(100, 20, 200)

	// Cloak/body (wide, dramatic)
	fillRect(img, 4, 12, 24, 26, cloak)
	fillRect(img, 2, 18, 28, 20, cloak)
	// Cloak edges/folds
	vLine(img, 4, 14, 24, cloakHi)
	vLine(img, 27, 14, 24, nrgba(35, 18, 60))
	hLine(img, 12, 4, 24, cloakHi)
	// Torso/armor
	fillRect(img, 8, 14, 16, 14, armor)
	hLine(img, 14, 9, 14, armorHi)
	vLine(img, 9, 15, 12, armorHi)
	// Head
	fillRect(img, 9, 2, 14, 13, armor)
	hLine(img, 3, 10, 12, armorHi)
	vLine(img, 10, 3, 10, armorHi)
	// Crown
	set(img, 12, 1, crown)
	set(img, 15, 0, crown)
	set(img, 18, 1, crown)
	hLine(img, 2, 11, 10, crown)
	// Eyes
	fillRect(img, 12, 6, 3, 3, eyes)
	fillRect(img, 17, 6, 3, 3, eyes)
	set(img, 13, 7, eyeGlow)
	set(img, 18, 7, eyeGlow)
	// Void energy (aura effect)
	for _, p := range [][2]int{
		{1, 10}, {0, 15}, {1, 20}, {0, 25},
		{30, 10}, {31, 15}, {30, 20}, {31, 25},
		{5, 37}, {8, 38}, {12, 39}, {16, 38}, {20, 39}, {24, 38}, {27, 37},
	} {
		set(img, p[0], p[1], voidEnergy)
	}
	// Hands reaching forward
	fillRect(img, 2, 22, 5, 4, armor)
	fillRect(img, 25, 22, 5, 4, armor)
	set(img, 1, 24, voidEnergy)
	set(img, 30, 24, voidEnergy)
	// Outline
	for x := 8; x < 24; x++ {
		set(img, x, 1, nrgba(15, 8, 25))
	}
	return img
}

// GenerateSpriteSheet returns a horizontal sheet: first 6 chars (16x16 each),
// then enemies (16x16 each).  Golem and Sovereign are kept separate (larger).
func GenerateSpriteSheet() *image.NRGBA {
	sprites := [][16 * 16]uint8{
		sprWarrior1, sprWarrior2,
		sprMage1, sprMage2,
		sprRogue1, sprRogue2,
		sprGoblin, sprSkeleton, sprSlime, sprBat,
	}
	palettes := [][]color.NRGBA{
		pallWarrior, pallWarrior,
		pallMage, pallMage,
		pallRogue, pallRogue,
		pallGoblin, pallSkeleton, pallSlime, pallBat,
	}
	n := len(sprites)
	img := newImg(n*TileSize, TileSize)
	for i, spr := range sprites {
		drawSprite(img, spr, palettes[i], i*TileSize, 0)
	}
	return img
}
