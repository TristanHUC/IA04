package _map

import "github.com/hajimehoshi/ebiten/v2"

type Mode int

const (
	ModeMove Mode = iota
	ModeWall
)

type Map struct {
	Width, Height int
	Walls         [][2]int
}

type Button struct {
	x, y, width, height int
	text                string
	image               *ebiten.Image
	imageOptions        *ebiten.DrawImageOptions
	selected            bool
}
