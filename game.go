package main

import (
	// "fmt"
	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	appName string = "Go Lyff"
	screenWidth, screenHeight int = 640, 640
)

type game struct {
	window *glfw.Window

	paused bool
}

func NewGame() *game {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(screenWidth, screenHeight, appName, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return &game{window: window}
}

func DeleteGame(g *game) {
	glfw.Terminate()
}

func (g *game) Run() {
	for !g.window.ShouldClose() {
		g.pollEvents()
		g.update()
		g.draw()
	}
}

func (g *game) pollEvents() {
	glfw.PollEvents()
}

func (g *game) update() {

}

func (g *game) draw() {
	g.window.SwapBuffers()
}
