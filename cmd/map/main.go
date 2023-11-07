package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"log"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

func init() {
	_map.InitHud(SCREEN_WIDTH, SCREEN_HEIGHT)
}

func main() {
	m := _map.Map{}
	err := m.LoadFromFile("testmap")
	if err != nil {
		return
	}

	game := &_map.Game{
		ScreenWidth:  SCREEN_WIDTH,
		ScreenHeight: SCREEN_HEIGHT,
		CameraX:      0,
		CameraY:      0,
		Map:          m,
		CurrentMode:  _map.ModeMove,
	}

	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Pic")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	game.Map.SaveToFile("testmap")
}
