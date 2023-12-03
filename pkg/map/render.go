package _map

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	Map                       Map
	CurrentMode               Mode
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

var draggingPos [2]int
var draggingCameraPos [2]int

func (g *Game) Update() error {
	// update hud
	stopPropagation := UpdateHud(g)
	// on drag, move camera if move mode
	if g.CurrentMode == ModeMove {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if draggingPos == [2]int{-1, -1} {
				draggingPos = [2]int{x, y}
				draggingCameraPos = [2]int{int(g.CameraX), int(g.CameraY)}
			}
			g.CameraX = draggingCameraPos[0] + draggingPos[0] - x
			g.CameraY = draggingCameraPos[1] + draggingPos[1] - y
		} else {
			draggingPos = [2]int{-1, -1}
		}
	}
	// on click, add wall if wall mode
	if g.CurrentMode == ModeWall && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall does not already exist at this position
			exists := false
			for _, wall := range g.Map.Walls {
				if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, beer := range g.Map.BarPoints {
				if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, ManWCtoilet := range g.Map.ManToiletPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			if !exists {
				g.Map.Walls = append(g.Map.Walls, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}
	// on click, add beer if beer mode
	if g.CurrentMode == ModeBeer && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := false
			for _, wall := range g.Map.Walls {
				if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, beer := range g.Map.BarPoints {
				if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, ManWCtoilet := range g.Map.ManToiletPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			if !exists {
				g.Map.BarPoints = append(g.Map.BarPoints, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}
// on click, add beer if beer mode
	if g.CurrentMode == ModeBeer && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := false
			for _, wall := range g.Map.Walls {
				if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, beer := range g.Map.BarPoints {
				if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, ManWCtoilet := range g.Map.ManToiletPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			if !exists {
				g.Map.BarPoints = append(g.Map.BarPoints, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}
	// on click, add ManToilet if manWC mode
	if g.CurrentMode == ModeManWC && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := false
			for _, wall := range g.Map.Walls {
				if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, beer := range g.Map.BarPoints {
				if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, ManWCtoilet := range g.Map.ManToiletPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			if !exists {
				g.Map.ManToiletPoints = append(g.Map.ManToiletPoints, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}	
	// on click, add WomanToilet if womanWC mode
	if g.CurrentMode == ModeWomanWC && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := false
			for _, wall := range g.Map.Walls {
				if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, beer := range g.Map.BarPoints {
				if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, ManWCtoilet := range g.Map.ManToiletPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			for _, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					exists = true
				}
			}
			if !exists {
				g.Map.WomanToiletPoints = append(g.Map.WomanToiletPoints, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}
	// on click, remove wall if erase mode
	if g.CurrentMode == ModeErase && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall/Toilette/Beer does not already exist at this position
			for i, wall := range g.Map.Walls {
				if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
					g.Map.Walls = append(g.Map.Walls[:i], g.Map.Walls[i+1:]...)
					break
				}
			}
			for i, beer := range g.Map.BarPoints {
				if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
					g.Map.BarPoints = append(g.Map.BarPoints[:i], g.Map.BarPoints[i+1:]...)
					break
				}
			}
			for i, ManWCtoilet := range g.Map.BarPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					g.Map.BarPoints = append(g.Map.ManToiletPoints[:i], g.Map.ManToiletPoints[i+1:]...)
					break
				}
			}
			for i, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					g.Map.WomanToiletPoints = append(g.Map.WomanToiletPoints[:i], g.Map.WomanToiletPoints[i+1:]...)
					break
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.DrawMap(screen)
	g.DrawHud(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

func (g *Game) DrawMap(screen *ebiten.Image) {
	// draw grid lines
	startPosX := g.CameraX % 10
	startPosY := g.CameraY % 10
	for i := 0; i < g.ScreenWidth/10; i++ {
		vector.StrokeLine(
			screen,
			float32(g.CameraX-startPosX+i*10-g.CameraX),
			float32(startPosY),
			float32(g.CameraX-startPosX+i*10-g.CameraX),
			float32(g.ScreenHeight),
			1,
			color.Gray{Y: 240},
			false,
		)
	}
	for i := 0; i < g.ScreenHeight/10; i++ {
		vector.StrokeLine(
			screen,
			0,
			float32(g.CameraY-startPosY+i*10-g.CameraY),
			float32(g.ScreenWidth),
			float32(g.CameraY-startPosY+i*10-g.CameraY),
			1,
			color.Gray{Y: 240},
			false,
		)
	}

	// draw map
	for _, wall := range g.Map.Walls {
		vector.DrawFilledRect(
			screen,
			float32(wall[0])*10-float32(g.CameraX),
			float32(wall[1])*10-float32(g.CameraY),
			10,
			10,
			color.Black,
			false,
		)
	}
	// draw beer
	for _, beer := range g.Map.BarPoints {
		vector.DrawFilledRect(
			screen,
			float32(beer[0])*10-float32(g.CameraX),
			float32(beer[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R:201, G:201, B:0, A:255},
			false,
		)
	}
	// draw WC Woman
	for _, WomanWC := range g.Map.WomanToiletPoints {
		vector.DrawFilledRect(
			screen,
			float32(WomanWC[0])*10-float32(g.CameraX),
			float32(WomanWC[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R:255, G:0, B:200, A:255},
			false,
		)
	}
	// draw WC Man
	for _, ManWC := range g.Map.ManToiletPoints {
		vector.DrawFilledRect(
			screen,
			float32(ManWC[0])*10-float32(g.CameraX),
			float32(ManWC[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R:0, G:200, B:255, A:255},
			false,
		)
	}
}
