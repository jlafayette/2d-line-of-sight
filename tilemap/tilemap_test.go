package tilemap

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

const (
	testNumX     = 32
	testNumY     = 24
	testTilesize = 20
)

func tilemapMutations(num int) []mutation {
	rand.Seed(time.Now().UnixNano())
	m := make([]mutation, num)
	for i := 0; i < num; i++ {
		m[i].x = rand.Intn(testNumX)
		m[i].y = rand.Intn(testNumY)
		m[i].value = rand.Intn(2) == 0
	}
	return m
}

func BenchmarkTileMap_CalculateEdges(b *testing.B) {

	tilemap := NewTileMap(testNumX, testNumY, testTilesize, false)
	mutations := tilemapMutations(1000)
	startMut := tilemapMutations(100)
	for j := range startMut {
		tilemap.Set(startMut[j].x, startMut[j].y, startMut[j].value)
		tilemap.CalculateEdges()
	}
	b.ResetTimer()
	j := 0
	for i := 0; i < b.N; i++ {
		j = i % len(mutations)
		tilemap.Set(mutations[j].x, mutations[j].y, mutations[j].value)
		tilemap.CalculateEdges()
	}

}
