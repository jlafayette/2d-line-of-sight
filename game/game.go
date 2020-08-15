package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	// W is screen width
	W = 640
	// H is screen height
	H = 480
	// tilesize is size of a tile
	tilesize = 20
	// number of tiles in x direction
	nx = W / tilesize
	// number of tiles in y direction
	ny = H / tilesize
)

// Game implements ebiten.Game interface and stores the game state.
//
// The methods run in the following order (each one is run once in this order
// if fps is 60 and display is 60 Hz):
//	Update
//	Draw
//	Layout
type Game struct {
	tilemap *TileMap
}

// NewGame creates a new Game
func NewGame() *Game {
	tilemap := NewTileMap(nx, ny)
	return &Game{tilemap}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := tilemapIndexOfCursor()
		g.tilemap.Set(x, y, true)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := tilemapIndexOfCursor()
		g.tilemap.Set(x, y, false)
	}
	return nil
}

// Get tilemap x and y indexes of the mouse cursor
func tilemapIndexOfCursor() (int, int) {
	mx, my := ebiten.CursorPosition()
	return mx / tilesize, my / tilesize
}

// Draw is called every frame. The frame frequency depends on the display's
// refresh rate, so if the display is 60 Hz, Draw will be called 60 times per
// second.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	// Draw tiles
	for x := range g.tilemap.m {
		for y := range g.tilemap.m[x] {
			if g.tilemap.Get(x, y) {
				ebitenutil.DrawRect(
					screen,                         // what to draw on
					float64(x*tilesize),            // x pos
					float64(y*tilesize),            // y pos
					tilesize,                       // width
					tilesize,                       // height
					color.RGBA{100, 100, 100, 255}, // color
				)
			}
		}
	}
	// Draw edges
	for _, e := range g.tilemap.edges {
		// Draw orange lines for the edges
		ebitenutil.DrawLine(
			screen,                        // what to draw on
			float64(e.Start.X),            // x1
			float64(e.Start.Y),            // y1
			float64(e.End.X),              // x2
			float64(e.End.Y),              // y2
			color.RGBA{200, 120, 10, 255}, // color
		)
		// Draw red square at the start
		ebitenutil.DrawRect(
			screen,                       // what to draw on
			float64(e.Start.X)-2,         // x pos
			float64(e.Start.Y)-2,         // y pos
			4,                            // width
			4,                            // height
			color.RGBA{200, 10, 10, 255}, // color
		)
		// Draw smaller yellow square at the end.
		ebitenutil.DrawRect(
			screen,                        // what to draw on
			float64(e.End.X)-1,            // x pos
			float64(e.End.Y)-1,            // y pos
			2,                             // width
			2,                             // height
			color.RGBA{200, 200, 10, 255}, // color
		)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", len(g.tilemap.edges)), 10, 10)
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return W, H
}
