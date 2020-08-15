package game

import (
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

func BenchmarkTileMap_CalculateEdges(b *testing.B) {

	tilemap := NewTileMap(nx, ny)
	mutations := tilemapMutations(1000)
	startMut := tilemapMutations(100)
	for j := range startMut {
		tilemap.Set(startMut[j].x, startMut[j].y, startMut[j].value)
		tilemap.CalculateEdges()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range mutations {
			tilemap.Set(mutations[j].x, mutations[j].y, mutations[j].value)
			tilemap.CalculateEdges()
		}
	}
}
