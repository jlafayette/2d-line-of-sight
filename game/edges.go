package game

// Point is a 2d coordinate.
type Point struct {
	X, Y int
}

// Edge marks the border of a tile or tiles that exist.
type Edge struct {
	Start Point
	End   Point
}

type tile struct {
	edgeID    [4]int
	edgeExist [4]bool
	exist     bool
}

// Direction indexes to use in the tile edgeID and edgeExist arrays.
const (
	north = iota
	south
	east
	west
)

// TileMap is a grid of tiles with edges along the borders of the existing tiles.
type TileMap struct {
	m     [][]tile
	edges []Edge
}

// NewTileMap creates a new TileMap with the given width and height.
func NewTileMap(nx, ny int) *TileMap {
	tilemap := make([][]tile, nx)
	for i := range tilemap {
		tilemap[i] = make([]tile, ny)
	}
	return &TileMap{tilemap, make([]Edge, 0, 512)}
}

// Get a tile's exist value at the given x and y coordinate.
func (t *TileMap) Get(x, y int) bool {
	return t.m[x][y].exist
}

// Set a tile's exist value at the given x and y coordinate.
func (t *TileMap) Set(x, y int, value bool) {
	if x < len(t.m) && x >= 0 && y < len(t.m[0]) && y >= 0 {
		v := t.Get(x, y)
		if v != value {
			t.m[x][y].exist = value
			t.CalculateEdges()
		}
	} // else {
	// 	fmt.Printf("%d, %d is out of bounds\n", x, y)
	// }
}

// CalculateEdges populates the slice of edges based on the existing tiles.
func (t *TileMap) CalculateEdges() {
	// clear everything
	for x := 0; x < len(t.m); x++ {
		for y := 0; y < len(t.m[x]); y++ {
			t.m[x][y].edgeID[0] = 0
			t.m[x][y].edgeID[1] = 0
			t.m[x][y].edgeID[2] = 0
			t.m[x][y].edgeID[3] = 0
			t.m[x][y].edgeExist[0] = false
			t.m[x][y].edgeExist[1] = false
			t.m[x][y].edgeExist[2] = false
			t.m[x][y].edgeExist[3] = false
		}
	}
	// Reset edges list but keep capacity
	t.edges = t.edges[:0]

	var nn, sn, en, wn tile

	for x := 0; x < len(t.m); x++ {
		for y := 0; y < len(t.m[x]); y++ {

			// define neighboring tiles with special cases for the edge tiles
			// nn: north neighbor
			// sn: south neighbor
			// en: east neighbor
			// wn: west neighbor
			if t.m[x][y].exist {
				if y > 0 {
					nn = t.m[x][y-1]
				} else {
					nn = tile{}
					nn.exist = true
				}
				if y < len(t.m[x])-1 {
					sn = t.m[x][y+1]
				} else {
					sn = tile{}
					sn.exist = true
				}
				if x < len(t.m)-1 {
					en = t.m[x+1][y]
				} else {
					en = tile{}
					en.exist = true
				}
				if x > 0 {
					wn = t.m[x-1][y]
				} else {
					wn = tile{}
					wn.exist = true
				}

				// add or extend neighboring edges for all directions.

				// handle north edge
				if !nn.exist {
					// check if west neighbor has a north edge to extend
					if wn.edgeExist[north] {
						// extend edge
						t.edges[wn.edgeID[north]].End.X += tilesize
						t.m[x][y].edgeID[north] = wn.edgeID[north]
						t.m[x][y].edgeExist[north] = true
					} else {
						// add north edge
						e1 := Edge{
							Point{x * tilesize, y * tilesize},
							Point{x*tilesize + tilesize, y * tilesize},
						}
						t.edges = append(t.edges, e1)
						t.m[x][y].edgeID[north] = len(t.edges) - 1
						t.m[x][y].edgeExist[north] = true
					}
				}

				if !sn.exist {
					// check if west neighbor has a south edge to extend
					if wn.edgeExist[south] {
						t.edges[wn.edgeID[south]].End.X += tilesize
						t.m[x][y].edgeID[south] = wn.edgeID[south]
						t.m[x][y].edgeExist[south] = true
					} else {
						// add south edge
						e2 := Edge{
							Point{x * tilesize, y*tilesize + tilesize},
							Point{x*tilesize + tilesize, y*tilesize + tilesize},
						}
						t.edges = append(t.edges, e2)
						t.m[x][y].edgeID[south] = len(t.edges) - 1
						t.m[x][y].edgeExist[south] = true
					}

				}

				if !en.exist {
					// check if north neighbor has an east edge to extend down
					if nn.edgeExist[east] {
						t.edges[nn.edgeID[east]].End.Y += tilesize
						t.m[x][y].edgeID[east] = nn.edgeID[east]
						t.m[x][y].edgeExist[east] = true
					} else {
						// add east edge
						e3 := Edge{
							Point{x*tilesize + tilesize, y * tilesize},
							Point{x*tilesize + tilesize, y*tilesize + tilesize},
						}
						t.edges = append(t.edges, e3)
						t.m[x][y].edgeID[east] = len(t.edges) - 1
						t.m[x][y].edgeExist[east] = true
					}
				}

				if !wn.exist {
					// check if north neighbor has a west edge to extend down
					if nn.edgeExist[west] {
						t.edges[nn.edgeID[west]].End.Y += tilesize
						t.m[x][y].edgeID[west] = nn.edgeID[west]
						t.m[x][y].edgeExist[west] = true
					} else {
						// add west edge
						e4 := Edge{
							Point{x * tilesize, y * tilesize},
							Point{x * tilesize, y*tilesize + tilesize},
						}
						t.edges = append(t.edges, e4)
						t.m[x][y].edgeID[west] = len(t.edges) - 1
						t.m[x][y].edgeExist[west] = true
					}
				}
			}
		}
	}
}
