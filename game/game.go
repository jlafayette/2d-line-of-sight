package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// Game implements ebiten.Game interface and stores the game state.
//
// The methods run in the following order (each one is run once in this order
// if fps is 60 and display is 60 Hz):
//	Update
//	Draw
//	Layout
type Game struct {
	width  int
	height int
}

// NewGame creates a new Game
func NewGame(width, height int) *Game {

	return &Game{width, height}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

// Draw is called every frame. The frame frequency depends on the display's
// refresh rate, so if the display is 60 Hz, Draw will be called 60 times per
// second.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
