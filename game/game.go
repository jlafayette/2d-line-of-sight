package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/jlafayette/2d-line-of-sight/tilemap"
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

type visPolyPoint struct {
	x     float64
	y     float64
	angle float64
}

// Game implements ebiten.Game interface and stores the game state.
type Game struct {
	tilemap       *tilemap.TileMap
	visPolyPoints []visPolyPoint
	debugAngles   []float64
}

// NewGame creates a new Game
func NewGame() *Game {
	tilemap := tilemap.NewTileMap(nx, ny, tilesize, true)
	visPolyPoints := make([]visPolyPoint, 0, 1024)
	debugAngles := make([]float64, 0, 1024)
	return &Game{tilemap, visPolyPoints, debugAngles}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := tilemapIndexOfCursor()
		g.tilemap.Set(x, y, false)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := tilemapIndexOfCursor()
		g.tilemap.Set(x, y, true)
	}

	// cast rays
	mx, my := ebiten.CursorPosition()
	g.calculateVisbilityPolygon(float64(mx), float64(my), 1000)

	return nil
}

func (g *Game) calculateVisbilityPolygon(ox, oy, radius float64) {
	// Clear points (but keep capacity)
	g.visPolyPoints = g.visPolyPoints[:0]
	g.debugAngles = g.debugAngles[:0]

	// iterate over edges
	for _, edge := range g.tilemap.Edges {
		// g.tilemap.Edges[i].Start

		// Run once for start, once for end point
		for i := 0; i < 2; i++ {
			var rayX, rayY, baseAng float64
			switch i {
			case 0:
				rayX = float64(edge.Start.X) - ox
				rayY = float64(edge.Start.Y) - oy
			case 1:
				rayX = float64(edge.End.X) - ox
				rayY = float64(edge.End.Y) - oy
			}
			baseAng = math.Atan(rayY / rayX)
			if rayX < 0 {
				baseAng = baseAng + math.Pi
			}

			g.debugAngles = append(g.debugAngles, baseAng)

			// For each point, cast 3 rays, 1 directly at the point
			// and 1 a little bit to either side.
			var ang float64
			for j := 0; j < 3; j++ {
				switch j {
				case 0:
					ang = baseAng - 0.0001
				case 1:
					ang = baseAng
				case 2:
					ang = baseAng + 0.0001
				}

				// Create ray along angle for required distance.
				rayX = radius * math.Cos(ang)
				rayY = radius * math.Sin(ang)

				var minT1, minPx, minPy, minAng float64
				minT1 = math.Inf(1)

				// Check for intersection between ray and all edges
				// in the tilemap.
				for _, edge2 := range g.tilemap.Edges {
					// line segment vector
					sdx := float64(edge2.End.X - edge2.Start.X)
					sdy := float64(edge2.End.Y - edge2.Start.Y)

					// check for co-linear
					if math.Abs(sdx-rayX) > 0.0 && math.Abs(sdy-rayY) > 0.0 {
						var t2, t1 float64
						t2 = (rayX*(float64(edge2.Start.Y)-oy) + (rayY * (ox - float64(edge2.Start.X)))) / (sdx*rayY - sdy*rayX)
						t1 = (float64(edge2.Start.X) + sdx*t2 - ox) / rayX

						// If intersect point exists along the ray and along the
						// line segment, then intersect point is valid.
						if t1 > 0.0 && t2 >= 0.0 && t2 <= 1.0 {
							// Check if this intersect point is closest.
							if t1 < minT1 {
								minT1 = t1
								minPx = ox + rayX*t1
								minPy = oy + rayY*t1
								minAng = math.Atan((minPy - oy) / (minPx - ox))
								if minPx-ox < 0 {
									minAng = minAng + math.Pi
								}
							}
						}
					}
				}
				if minT1 < math.Inf(1) {
					g.visPolyPoints = append(g.visPolyPoints, visPolyPoint{minPx, minPy, minAng})
				}
			}
		}
	}
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
	for x := range g.tilemap.Tiles {
		for y := range g.tilemap.Tiles[x] {
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
	for _, e := range g.tilemap.Edges {
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

	// Draw visibility points
	mx, my := ebiten.CursorPosition()
	mxf := float64(mx)
	myf := float64(my)
	for _, pt := range g.visPolyPoints {
		// Draw yellow line from origin to the point
		ebitenutil.DrawLine(
			screen,                        // what to draw on
			mxf,                           // x1
			myf,                           // y1
			pt.x,                          // x2
			pt.y,                          // y2
			color.RGBA{200, 200, 10, 255}, // color
		)
	}

	// // Draw debug rays
	// for _, ang := range g.debugAngles {
	// 	r := 200.0
	// 	rayX := math.Cos(ang)
	// 	rayY := math.Sin(ang)
	// 	px := mxf + rayX*r
	// 	py := myf + rayY*r
	// 	ebitenutil.DrawLine(
	// 		screen,                       // what to draw on
	// 		mxf,                          // x1
	// 		myf,                          // y1
	// 		px,                           // x2
	// 		py,                           // y2
	// 		color.RGBA{10, 10, 200, 255}, // color
	// 	)
	// }

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("edges: %d", len(g.tilemap.Edges)), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("rays: %d", len(g.visPolyPoints)), 10, 30)

	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("debug angles: %v", g.debugAngles), 10, 50)
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return W, H
}
