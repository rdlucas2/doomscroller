package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rdlucas2/doomscroller/internal/engine"
	"github.com/rdlucas2/doomscroller/internal/render"

	_ "github.com/rdlucas2/doomscroller/content/stories"
)

type EbitenGame struct {
	eng      *engine.Engine
	renderer *render.Renderer
	lastTime time.Time
}

func (g *EbitenGame) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	if dt > 0.1 {
		dt = 0.1 // cap to avoid spiral of death on focus loss
	}
	g.lastTime = now
	err := g.eng.Update(dt)
	if engine.IsQuit(err) {
		return ebiten.Termination
	}
	return err
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen, g.eng.BuildState())
}

func (g *EbitenGame) Layout(_, _ int) (int, int) {
	return render.ScreenW, render.ScreenH
}

func main() {
	ebiten.SetWindowTitle("Doom Scroller — Tower of the Void")
	ebiten.SetWindowSize(render.ScreenW*2, render.ScreenH*2)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &EbitenGame{
		eng:      engine.New(),
		renderer: render.New(),
		lastTime: time.Now(),
	}
	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}
