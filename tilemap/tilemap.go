package tilemap

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
	Tiles    [][]tile
	Edges    []Edge
	tilesize int
}

// NewTileMap creates a new TileMap with the given width and height.
func NewTileMap(nx, ny, tilesize int) *TileMap {
	tilemap := make([][]tile, nx)
	for i := range tilemap {
		tilemap[i] = make([]tile, ny)
	}
	return &TileMap{tilemap, make([]Edge, 0, 512), tilesize}
}

// Get a tile's exist value at the given x and y coordinate.
func (t *TileMap) Get(x, y int) bool {
	return t.Tiles[x][y].exist
}

// Set a tile's exist value at the given x and y coordinate.
func (t *TileMap) Set(x, y int, value bool) {
	if x < len(t.Tiles) && x >= 0 && y < len(t.Tiles[0]) && y >= 0 {
		v := t.Get(x, y)
		if v != value {
			t.Tiles[x][y].exist = value
			t.CalculateEdges()
		}
	} // else {
	// 	fmt.Printf("%d, %d is out of bounds\n", x, y)
	// }
}

// CalculateEdges populates the slice of edges based on the existing tiles.
func (t *TileMap) CalculateEdges() {
	// clear everything
	for x := 0; x < len(t.Tiles); x++ {
		for y := 0; y < len(t.Tiles[x]); y++ {
			t.Tiles[x][y].edgeID[0] = 0
			t.Tiles[x][y].edgeID[1] = 0
			t.Tiles[x][y].edgeID[2] = 0
			t.Tiles[x][y].edgeID[3] = 0
			t.Tiles[x][y].edgeExist[0] = false
			t.Tiles[x][y].edgeExist[1] = false
			t.Tiles[x][y].edgeExist[2] = false
			t.Tiles[x][y].edgeExist[3] = false
		}
	}
	// Reset edges list but keep capacity
	t.Edges = t.Edges[:0]

	var nn, sn, en, wn tile

	for x := 0; x < len(t.Tiles); x++ {
		for y := 0; y < len(t.Tiles[x]); y++ {

			// define neighboring tiles with special cases for the edge tiles
			// nn: north neighbor
			// sn: south neighbor
			// en: east neighbor
			// wn: west neighbor
			if t.Tiles[x][y].exist {
				if y > 0 {
					nn = t.Tiles[x][y-1]
				} else {
					nn = tile{}
					nn.exist = true
				}
				if y < len(t.Tiles[x])-1 {
					sn = t.Tiles[x][y+1]
				} else {
					sn = tile{}
					sn.exist = true
				}
				if x < len(t.Tiles)-1 {
					en = t.Tiles[x+1][y]
				} else {
					en = tile{}
					en.exist = true
				}
				if x > 0 {
					wn = t.Tiles[x-1][y]
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
						t.Edges[wn.edgeID[north]].End.X += t.tilesize
						t.Tiles[x][y].edgeID[north] = wn.edgeID[north]
						t.Tiles[x][y].edgeExist[north] = true
					} else {
						// add north edge
						e1 := Edge{
							Point{x * t.tilesize, y * t.tilesize},
							Point{x*t.tilesize + t.tilesize, y * t.tilesize},
						}
						t.Edges = append(t.Edges, e1)
						t.Tiles[x][y].edgeID[north] = len(t.Edges) - 1
						t.Tiles[x][y].edgeExist[north] = true
					}
				}

				if !sn.exist {
					// check if west neighbor has a south edge to extend
					if wn.edgeExist[south] {
						t.Edges[wn.edgeID[south]].End.X += t.tilesize
						t.Tiles[x][y].edgeID[south] = wn.edgeID[south]
						t.Tiles[x][y].edgeExist[south] = true
					} else {
						// add south edge
						e2 := Edge{
							Point{x * t.tilesize, y*t.tilesize + t.tilesize},
							Point{x*t.tilesize + t.tilesize, y*t.tilesize + t.tilesize},
						}
						t.Edges = append(t.Edges, e2)
						t.Tiles[x][y].edgeID[south] = len(t.Edges) - 1
						t.Tiles[x][y].edgeExist[south] = true
					}

				}

				if !en.exist {
					// check if north neighbor has an east edge to extend down
					if nn.edgeExist[east] {
						t.Edges[nn.edgeID[east]].End.Y += t.tilesize
						t.Tiles[x][y].edgeID[east] = nn.edgeID[east]
						t.Tiles[x][y].edgeExist[east] = true
					} else {
						// add east edge
						e3 := Edge{
							Point{x*t.tilesize + t.tilesize, y * t.tilesize},
							Point{x*t.tilesize + t.tilesize, y*t.tilesize + t.tilesize},
						}
						t.Edges = append(t.Edges, e3)
						t.Tiles[x][y].edgeID[east] = len(t.Edges) - 1
						t.Tiles[x][y].edgeExist[east] = true
					}
				}

				if !wn.exist {
					// check if north neighbor has a west edge to extend down
					if nn.edgeExist[west] {
						t.Edges[nn.edgeID[west]].End.Y += t.tilesize
						t.Tiles[x][y].edgeID[west] = nn.edgeID[west]
						t.Tiles[x][y].edgeExist[west] = true
					} else {
						// add west edge
						e4 := Edge{
							Point{x * t.tilesize, y * t.tilesize},
							Point{x * t.tilesize, y*t.tilesize + t.tilesize},
						}
						t.Edges = append(t.Edges, e4)
						t.Tiles[x][y].edgeID[west] = len(t.Edges) - 1
						t.Tiles[x][y].edgeExist[west] = true
					}
				}
			}
		}
	}
}
