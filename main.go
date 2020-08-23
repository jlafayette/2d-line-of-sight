package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/jlafayette/2d-line-of-sight/game"
)

func main() {

	// In this test, window size is equal to screen size, so no pixelation
	// or stretching will occur.
	ebiten.SetWindowSize(game.W, game.H)

	ebiten.SetWindowTitle("Edges from Tilemap")
	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
