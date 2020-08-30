package game

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/jlafayette/2d-line-of-sight/tilemap"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// W is screen width
	W = 1280
	// H is screen height
	H = 960
	// tilesize is size of a tile
	tilesize = 40
	// number of tiles in x direction
	nx = W / tilesize
	// number of tiles in y direction
	ny = H / tilesize
	// Extra Pi constants
	halfPi   = math.Pi / 2
	doublePi = math.Pi * 2
)

var (
	wallColor           = colorful.Hsv(359.0, 0.05, 0.5)
	openColor           = colorful.Hsv(200.0, 0.9, 0.4)
	lightColor          = colorful.Hsv(45.0, 0.9, 0.9)
	debugRayColor       = colorful.Hsv(55.0, 1.0, 1.0)
	debugEdgeColor      = colorful.Hsv(90.0, 1.0, 0.0)
	debugEdgePointColor = colorful.Hsv(90.0, 1.0, 0.0)
	noAlpha             = color.RGBA{0, 0, 0, 0}
)

// mode is an enum for display options
type mode int

const (
	place mode = iota
	light
	debug
)

func (k mode) String() string {
	return [...]string{"place", "light", "debug"}[k]
}

type visPolyPoint struct {
	x     float64
	y     float64
	angle float64
}

func newVisPolyPoint(x, y, angle float64) visPolyPoint {
	return visPolyPoint{x, y, angle}
}

// map angle (radians -Pi..Pi) to the hue of the color (0..360)
func calculateColor(angle float64) color.Color {
	h := (angle + math.Pi) * (180 / math.Pi)
	return colorful.Hsv(h, 1.0, 1.0)
}

type visPolyPoints []visPolyPoint

func (v visPolyPoints) Len() int {
	return len(v)
}
func (v visPolyPoints) Less(i, j int) bool {
	return v[i].angle < v[j].angle
}
func (v visPolyPoints) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// Game implements ebiten.Game interface and stores the game state.
type Game struct {
	mode       mode
	mx         int
	my         int
	tilemap    *tilemap.TileMap
	visPoints  visPolyPoints
	wallImage  *ebiten.Image
	openImage  *ebiten.Image
	lightMask  *ebiten.Image
	lightImage *ebiten.Image
	tileImage  *ebiten.Image
}

// NewGame creates a new Game
func NewGame() *Game {
	tilemap := tilemap.NewTileMap(nx, ny, tilesize, false)
	visPts := make([]visPolyPoint, 0, 1024)

	wallImage, _ := ebiten.NewImage(W, H, ebiten.FilterDefault)
	wallImage.Fill(noAlpha)
	openImage, _ := ebiten.NewImage(W, H, ebiten.FilterDefault)
	openImage.Fill(openColor)
	lightMask, _ := ebiten.NewImage(W, H, ebiten.FilterDefault)
	lightMask.Fill(noAlpha)
	lightImage, _ := ebiten.NewImage(W, H, ebiten.FilterDefault)
	lightImage.Fill(lightColor)

	// For subtracting a tile from alpha
	tileImage, _ := ebiten.NewImage(tilesize, tilesize, ebiten.FilterDefault)
	tileImage.Fill(noAlpha)

	return &Game{place, 0, 0, tilemap, visPts, wallImage, openImage, lightMask, lightImage, tileImage}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update(screen *ebiten.Image) error {
	// Get a bounded cursor postion.
	g.mx, g.my = ebiten.CursorPosition()
	g.mx = max(1, min(W-1, g.mx))
	g.my = max(1, min(H-1, g.my))

	// Handle adding and removing walls from the tilemap.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := g.tilemapIndexOfCursor()
		g.tilemap.Set(x, y, true)

		ebitenutil.DrawRect(
			g.wallImage,         // what to draw on
			float64(x*tilesize), // x pos
			float64(y*tilesize), // y pos
			tilesize,            // width
			tilesize,            // height
			wallColor,           // color
		)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := g.tilemapIndexOfCursor()
		g.tilemap.Set(x, y, false)

		op := &ebiten.DrawImageOptions{}
		op.CompositeMode = ebiten.CompositeModeCopy
		op.GeoM.Translate(float64(x*tilesize), float64(y*tilesize))
		g.wallImage.DrawImage(g.tileImage, op)

	}

	// Handle mode changes.
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.mode = place
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.mode = light
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.mode = debug
	}

	// cast rays
	g.calculateVisbilityPolygon()

	return nil
}

func (g *Game) calculateVisbilityPolygon() {
	ox := float64(g.mx)
	oy := float64(g.my)
	radius := 1000.0
	// Clear points (but keep capacity)
	g.visPoints = g.visPoints[:0]

	// Iterate over edges and cast rays to start and end of each.
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
			baseAng = math.Atan2(rayY, rayX)

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
								minAng = math.Atan2(minPy-oy, minPx-ox)
							}
						}
					}
				}
				if minT1 < math.Inf(1) {
					g.visPoints = append(g.visPoints, newVisPolyPoint(minPx, minPy, minAng))
				}
			}
		}
	}

	// Sort points by angle
	sort.Sort(g.visPoints)
}

// Get tilemap x and y indexes of the mouse cursor
func (g *Game) tilemapIndexOfCursor() (int, int) {
	return g.mx / tilesize, g.my / tilesize
}

// Return if cursor is in the open, accounting the edge pixel.
func (g *Game) cursorInOpen() bool {
	posOpen := !g.tilemap.Get(g.mx/tilesize, g.my/tilesize)
	posOffsetOpen := !g.tilemap.Get((g.mx-1)/tilesize, (g.my-1)/tilesize)
	return posOpen && posOffsetOpen
}

// Draw is called every frame. The frame frequency depends on the display's
// refresh rate, so if the display is 60 Hz, Draw will be called 60 times per
// second.
func (g *Game) Draw(screen *ebiten.Image) {
	mxf := float64(g.mx)
	myf := float64(g.my)

	screen.Clear()
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.openImage, op)
	op = &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeSourceAtop
	screen.DrawImage(g.wallImage, op)

	// Draw light if in light mode and cursor is in the open (not a wall).
	if g.mode >= light && g.cursorInOpen() {
		g.lightMask.Fill(color.RGBA{0, 0, 0, 0})
		// draw light triangles
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Address = ebiten.AddressRepeat
		opt.CompositeMode = ebiten.CompositeModeSourceOut
		for i, pt := range g.visPoints {
			nextPt := g.visPoints[(i+1)%len(g.visPoints)]

			// Draw triangle of area between mouse position, current point, and next point
			g.lightMask.DrawTriangles(
				[]ebiten.Vertex{
					{DstX: float32(mxf), DstY: float32(myf), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
					{DstX: float32(nextPt.x), DstY: float32(nextPt.y), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
					{DstX: float32(pt.x), DstY: float32(pt.y), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
				},
				[]uint16{0, 1, 2},
				g.lightImage,
				opt,
			)
		}
		op = &ebiten.DrawImageOptions{}
		op.CompositeMode = ebiten.CompositeModeSourceAtop
		screen.DrawImage(g.lightMask, op)
	}

	if g.mode == debug {
		// Draw visibility rays
		for _, pt := range g.visPoints {
			ebitenutil.DrawLine(
				screen,        // what to draw on
				mxf,           // x1
				myf,           // y1
				pt.x,          // x2
				pt.y,          // y2
				debugRayColor, // color
			)
		}
	}

	if g.mode == debug || g.mode == place {
		// Draw edges
		for _, e := range g.tilemap.Edges {
			// Draw line for the edge
			ebitenutil.DrawLine(
				screen,             // what to draw on
				float64(e.Start.X), // x1
				float64(e.Start.Y), // y1
				float64(e.End.X),   // x2
				float64(e.End.Y),   // y2
				debugEdgeColor,     // color
			)
			// Draw square at the start
			ebitenutil.DrawRect(
				screen,               // what to draw on
				float64(e.Start.X)-2, // x pos
				float64(e.Start.Y)-2, // y pos
				4,                    // width
				4,                    // height
				debugEdgePointColor,  // color
			)
			// Draw square at the end.
			ebitenutil.DrawRect(
				screen,              // what to draw on
				float64(e.End.X)-2,  // x pos
				float64(e.End.Y)-2,  // y pos
				4,                   // width
				4,                   // height
				debugEdgePointColor, // color
			)
		}
	}

	if g.mode == place || g.mode == debug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("edges: %d", len(g.tilemap.Edges)), 10, 10)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("rays: %d", len(g.visPoints)), 10, 30)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mode: %s", g.mode), 10, 50)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cursor: %d, %d [%d][%d]", g.mx, g.my, g.mx/tilesize, g.my/tilesize), 10, 70)
	}
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return W, H
}
