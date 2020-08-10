package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	// W is screen width
	W = 640
	// H is screen height
	H = 480
	// s is size of a tile
	s = 20
	// number of tiles in x direction
	nx = W / s
	// number of tiles in y direction
	ny = H / s
)

// Game implements ebiten.Game interface and stores the game state.
//
// The methods run in the following order (each one is run once in this order
// if fps is 60 and display is 60 Hz):
//	Update
//	Draw
//	Layout
type Game struct {
	tilemap [][]bool
}

// NewGame creates a new Game
func NewGame() *Game {
	tilemap := make([][]bool, nx)
	for i := range tilemap {
		tilemap[i] = make([]bool, ny)
	}
	return &Game{tilemap}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := tilemapIndexOfCursor()
		g.tilemap[x][y] = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := tilemapIndexOfCursor()
		g.tilemap[x][y] = false
	}
	return nil
}

// Get tilemap x and y indexes of the mouse cursor
func tilemapIndexOfCursor() (int, int) {
	mx, my := ebiten.CursorPosition()
	return mx / s, my / s
}

// Draw is called every frame. The frame frequency depends on the display's
// refresh rate, so if the display is 60 Hz, Draw will be called 60 times per
// second.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	for x := range g.tilemap {
		for y := range g.tilemap[x] {
			if g.tilemap[x][y] {
				ebitenutil.DrawRect(screen, float64(x*s), float64(y*s), s, s, color.RGBA{100, 100, 100, 255})
			}
		}
	}
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return W, H
}
