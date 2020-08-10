package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/jlafayette/edges-from-tilemap/game"
)

var (
	// Screen width and height indicates how many pixels to draw, not the
	// window dimensions.
	screenWidth  = 620
	screenHeight = 480
)

func main() {

	// In this test, window size is equal to screen size, so no pixelation
	// or stretching will occur.
	ebiten.SetWindowSize(screenWidth, screenHeight)

	ebiten.SetWindowTitle("Ebiten Test")
	game := game.NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
