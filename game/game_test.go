package game

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

type mutation struct {
	x     int
	y     int
	value bool
}

func tilemapMutations(num int) []mutation {
	rand.Seed(time.Now().UnixNano())
	m := make([]mutation, num)
	for i := 0; i < num; i++ {
		m[i].x = rand.Intn(nx)
		m[i].y = rand.Intn(ny)
		m[i].value = rand.Intn(2) == 0
	}
	return m
}

type mousePos struct {
	x int
	y int
}

func mousePositions() []mousePos {
	// g.mx, g.my = ebiten.CursorPosition()
	// g.mx = max(1, min(W-1, g.mx))
	// g.my = max(1, min(H-1, g.my))
	m := make([]mousePos, W*H)
	for x := 1; x < W; x++ {
		for y := 1; y < H; y++ {
			m = append(m, mousePos{x, y})
		}
	}
	return m
}

func BenchmarkGame_calculateVisbilityPolygon(b *testing.B) {

	// t := tilemap.NewTileMap(testNumX, testNumY, testTilesize, false)
	g := NewGame()
	mutations := tilemapMutations(1000)
	startMut := tilemapMutations(1000)
	for j := range startMut {
		g.tilemap.Set(startMut[j].x, startMut[j].y, startMut[j].value)
		g.tilemap.CalculateEdges()
	}
	positions := mousePositions()
	b.ResetTimer()
	var minE, maxE int
	minE = math.MaxInt32
	maxE = 0
	j := 0
	k := 0
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k = i % len(positions)
		if i%100 == 0 {
			j = i % len(mutations)
			g.tilemap.Set(mutations[j].x, mutations[j].y, mutations[j].value)
			g.tilemap.CalculateEdges()
		}
		minE = min(minE, len(g.tilemap.Edges))
		maxE = max(maxE, len(g.tilemap.Edges))
		b.StartTimer()
		g.mx = positions[k].x
		g.my = positions[k].y
		g.calculateVisbilityPolygon()
	}
	fmt.Printf("\n----  ")
	fmt.Printf("min max: %d - %d", minE, maxE)
	fmt.Printf("  ----\n")
}
