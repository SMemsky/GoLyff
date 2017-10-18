package main

import "runtime"

func init() {
	// Make sure main func is called on the main thread
	runtime.LockOSThread()
}

func main() {
	game := NewGame()
	defer DeleteGame(game)

	game.Run()
}
