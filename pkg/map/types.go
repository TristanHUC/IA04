package _map

import "github.com/hajimehoshi/ebiten/v2"

type Mode int

const (
	ModeMove Mode = iota
	ModeWall
	ModeErase
	ModeBeer
	ModeManWC
	ModeWomanWC
)

var modeToCursor = map[Mode]ebiten.CursorShapeType{
	ModeMove:  ebiten.CursorShapeMove,
	ModeWall:  ebiten.CursorShapeCrosshair,
	ModeErase: ebiten.CursorShapeCrosshair,
	ModeBeer: ebiten.CursorShapeCrosshair,
	ModeManWC: ebiten.CursorShapeCrosshair,
	ModeWomanWC: ebiten.CursorShapeCrosshair,
}

type Map struct {
	Width, Height     int
	Walls             [][2]int
	BarPoints         [][2]int
	ManToiletPoints   [][2]int
	WomanToiletPoints [][2]int
}

type Button struct {
	x, y, width, height int
	text                string
	image               *ebiten.Image
	imageOptions        *ebiten.DrawImageOptions
	selected            bool
	mode                Mode
}
