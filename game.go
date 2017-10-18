package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	appName string = "Go Lyff"
	screenWidth, screenHeight int = 640, 640
	framesPerSecond = 10

	fieldWidth, fieldHeight int = 64, 64
	updatesPerSecond = 200000	

	colWidth float64 = float64(screenWidth) / float64(fieldWidth)
	rowHeight float64 = float64(screenHeight) / float64(fieldHeight)
)

type MapPreset int

const (
	EmptyMapPreset MapPreset = iota
	StableSquaresPreset
)

func setViewport(width, height int) {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, float64(width), float64(height), 0, -1, 1)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Viewport(0, 0, int32(width), int32(height))
}

type game struct {
	window *glfw.Window

	paused bool
	hideGrid bool

	field [fieldWidth*fieldHeight]bool
	tempField [fieldWidth*fieldHeight]bool
}

func NewGame() *game {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(screenWidth, screenHeight, appName, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(1.0, 1.0, 1.0, 1.0);
	setViewport(screenWidth, screenHeight)

	g := &game{window: window, paused: true}
	g.clearMap(StableSquaresPreset)

	g.field[cellId(3, 2)] = true

	return g
}

func DeleteGame(g *game) {
	glfw.Terminate()
}

func (g *game) Run() {
	lastDraw := time.Now()

	for !g.window.ShouldClose() {
		g.pollEvents()
		g.update()

		time.Sleep(time.Until(lastDraw.Add(time.Duration(time.Second / framesPerSecond))))
		lastDraw = time.Now()
		g.draw()
	}
}

func (g *game) pollEvents() {
	glfw.PollEvents()
}

func (g *game) update() {
	for col := 0; col < fieldWidth; col++ {
		for row := 0; row < fieldHeight; row++ {
			nc := g.countNeighbors(col, row)

			if g.field[col + row * fieldWidth] {
				g.tempField[col + row * fieldWidth] = (nc == 2 || nc == 3)
			} else {
				g.tempField[col + row * fieldWidth] = (nc == 3)
			}
		}
	}

	for i := 0; i < len(g.field); i++ {
		g.field[i] = g.tempField[i]
	}
}

func (g *game) countNeighbors(x, y int) (count uint) {
	if g.field[cellId(x, y + 1)] { count += 1 } //Top
	if g.field[cellId(x, y - 1)] { count += 1 } //Down
	if g.field[cellId(x + 1, y)] { count += 1 } //Right
	if g.field[cellId(x - 1, y)] { count += 1 } //Left

	if g.field[cellId(x + 1, y + 1)] { count += 1 } //Top right
	if g.field[cellId(x - 1, y + 1)] { count += 1 } //Top left
	if g.field[cellId(x + 1, y - 1)] { count += 1 } //Down right
	if g.field[cellId(x - 1, y - 1)] { count += 1 } //Down left

	return
}

func cellId(x, y int) int {
	x %= fieldWidth;
	y %= fieldHeight;

	if x < 0 {
		x += fieldWidth
	}
	if y < 0 {
		y += fieldHeight
	}

	return x + y * fieldWidth
}

func (g *game) draw() {
	fmt.Println("Draw")

	gl.Clear(gl.COLOR_BUFFER_BIT)

	g.drawField()
	g.drawGrid()

	g.window.SwapBuffers()
}

func (g *game) drawField() {
	gl.Begin(gl.QUADS)
	gl.Color3d(0.0, 0.0, 0.0)
	for col := 0; col < fieldWidth; col++ {
		for row := 0; row < fieldHeight; row++ {
			if g.field[col + row * fieldWidth] {
				gl.Vertex2d(float64(col) * colWidth, float64(row) * rowHeight)
				gl.Vertex2d(float64(col+1) * colWidth, float64(row) * rowHeight)
				gl.Vertex2d(float64(col+1) * colWidth, float64(row+1) * rowHeight)
				gl.Vertex2d(float64(col) * colWidth, float64(row+1) * rowHeight)
			}
		}
	}
	gl.End()
}

func (g *game) drawGrid() {
	if g.hideGrid {
		return
	}

	gl.Begin(gl.LINES)
	gl.Color3d(0.5, 0.5, 0.5)
	for col := 0; col < fieldWidth; col++ {
		gl.Vertex2d(float64(col) * colWidth, 0)
		gl.Vertex2d(float64(col) * colWidth, float64(fieldHeight) * rowHeight)
	}
	for row := 0; row < fieldHeight; row++ {
		gl.Vertex2d(0, float64(row) * rowHeight)
		gl.Vertex2d(float64(fieldWidth) * colWidth, float64(row) * rowHeight)
	}
	gl.End()
}

func (g *game) clearMap(preset MapPreset) {
	switch preset {
	case EmptyMapPreset:
		g.setEmptyPreset()
	case StableSquaresPreset:
		g.setStableSquaresPreset()
	default:
		fmt.Println("Reached default")
		panic(1)
	}
}

func (g *game) setEmptyPreset() {
	for i := 0; i < len(g.field); i++ {
		g.field[i] = false
	}
}

func (g *game) setStableSquaresPreset() {
	for col := 0; col < fieldWidth; col++ {
		for row := 0; row < fieldHeight; row++ {
			g.field[col + row * fieldWidth] = !(col % 3 == 0 || row % 3 == 0)
		}
	}
}
