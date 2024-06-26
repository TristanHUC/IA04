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
			exists := g.IsCaseTaken(x, y)
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
			exists := g.IsCaseTaken(x, y)
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
			exists := g.IsCaseTaken(x, y)
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
			exists := g.IsCaseTaken(x, y)
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
			exists := g.IsCaseTaken(x, y)
			if !exists {
				g.Map.WomanToiletPoints = append(g.Map.WomanToiletPoints, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}
	// on click, add BarmenArea if BarmanArea mode
	if g.CurrentMode == BarmenArea && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := g.IsCaseTaken(x, y)
			if !exists {
				g.Map.BarmenArea = append(g.Map.BarmenArea, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}

	// on click, add BeerTap if beerTap mode
	if g.CurrentMode == ModeBeerTap && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := g.IsCaseTaken(x, y)
			if !exists {
				g.Map.BeerTaps = append(g.Map.BeerTaps, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}

	// on click, add CounterArea if CounterArea mode
	if g.CurrentMode == CounterArea && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := g.IsCaseTaken(x, y)
			if !exists {
				g.Map.CounterArea = append(g.Map.CounterArea, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
				g.Map.Walls = append(g.Map.Walls, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}

	// on click, add Exit if Exit mode
	if g.CurrentMode == Exit && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := g.IsCaseTaken(x, y)
			if !exists {
				g.Map.Exit = append(g.Map.Exit, [2]int{
					int(float32(x+g.CameraX) / 10),
					int(float32(y+g.CameraY) / 10),
				})
			}

		}
	}

	// on click, add Enter if Enter mode
	if g.CurrentMode == Enter && !stopPropagation {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// if wall,beer,toilet does not already exist at this position
			exists := g.IsCaseTaken(x, y)
			if !exists {
				g.Map.Enter = append(g.Map.Enter, [2]int{
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
			for i, ManWCtoilet := range g.Map.ManToiletPoints {
				if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					g.Map.ManToiletPoints = append(g.Map.ManToiletPoints[:i], g.Map.ManToiletPoints[i+1:]...)
					break
				}
			}
			for i, WomanWCtoilet := range g.Map.WomanToiletPoints {
				if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
					g.Map.WomanToiletPoints = append(g.Map.WomanToiletPoints[:i], g.Map.WomanToiletPoints[i+1:]...)
					break
				}
			}
			for i, BarmanArea := range g.Map.BarmenArea {
				if BarmanArea[0] == int(float32(x+g.CameraX)/10) && BarmanArea[1] == int(float32(y+g.CameraY)/10) {
					g.Map.BarmenArea = append(g.Map.BarmenArea[:i], g.Map.BarmenArea[i+1:]...)
					break
				}
			}
			for i, BeerTap := range g.Map.BeerTaps {
				if BeerTap[0] == int(float32(x+g.CameraX)/10) && BeerTap[1] == int(float32(y+g.CameraY)/10) {
					g.Map.BeerTaps = append(g.Map.BeerTaps[:i], g.Map.BeerTaps[i+1:]...)
					break
				}
			}
			for i, counterArea := range g.Map.CounterArea {
				if counterArea[0] == int(float32(x+g.CameraX)/10) && counterArea[1] == int(float32(y+g.CameraY)/10) {
					g.Map.CounterArea = append(g.Map.CounterArea[:i], g.Map.CounterArea[i+1:]...)
					break
				}
			}
			for i, Exit := range g.Map.Exit {
				if Exit[0] == int(float32(x+g.CameraX)/10) && Exit[1] == int(float32(y+g.CameraY)/10) {
					g.Map.Exit = append(g.Map.Exit[:i], g.Map.Exit[i+1:]...)
					break
				}
			}
			for i, Enter := range g.Map.Enter {
				if Enter[0] == int(float32(x+g.CameraX)/10) && Enter[1] == int(float32(y+g.CameraY)/10) {
					g.Map.Enter = append(g.Map.Enter[:i], g.Map.Enter[i+1:]...)
					break
				}
			}
		}
	}
	return nil
}

func (g *Game) IsCaseTaken(x int, y int) bool {
	for _, wall := range g.Map.Walls {
		if wall[0] == int(float32(x+g.CameraX)/10) && wall[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, beer := range g.Map.BarPoints {
		if beer[0] == int(float32(x+g.CameraX)/10) && beer[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, ManWCtoilet := range g.Map.ManToiletPoints {
		if ManWCtoilet[0] == int(float32(x+g.CameraX)/10) && ManWCtoilet[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, WomanWCtoilet := range g.Map.WomanToiletPoints {
		if WomanWCtoilet[0] == int(float32(x+g.CameraX)/10) && WomanWCtoilet[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, BarmanArea := range g.Map.BarmenArea {
		if BarmanArea[0] == int(float32(x+g.CameraX)/10) && BarmanArea[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, BeerTap := range g.Map.BeerTaps {
		if BeerTap[0] == int(float32(x+g.CameraX)/10) && BeerTap[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, counterArea := range g.Map.CounterArea {
		if counterArea[0] == int(float32(x+g.CameraX)/10) && counterArea[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, Exit := range g.Map.Exit {
		if Exit[0] == int(float32(x+g.CameraX)/10) && Exit[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	for _, Enter := range g.Map.Enter {
		if Enter[0] == int(float32(x+g.CameraX)/10) && Enter[1] == int(float32(y+g.CameraY)/10) {
			return true
		}
	}
	return false
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
	isCounterArea := false
	for _, wall := range g.Map.Walls {
		isCounterArea = false
		for _, counterArea := range g.Map.CounterArea {
			if wall[0] == counterArea[0] && wall[1] == counterArea[1] {
				vector.DrawFilledRect(
					screen,
					float32(counterArea[0])*10-float32(g.CameraX),
					float32(counterArea[1])*10-float32(g.CameraY),
					10,
					10,
					color.RGBA{R: 255, G: 255, B: 0, A: 255},
					false,
				)
				isCounterArea = true
				break
			}
		}
		if !isCounterArea {
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
	}
	// draw beer
	for _, beer := range g.Map.BarPoints {
		vector.DrawFilledRect(
			screen,
			float32(beer[0])*10-float32(g.CameraX),
			float32(beer[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R: 201, G: 201, B: 0, A: 255},
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
			color.RGBA{R: 255, G: 0, B: 200, A: 255},
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
			color.RGBA{R: 0, G: 200, B: 255, A: 255},
			false,
		)
	}
	// draw BarmenArea
	for _, BarmenArea := range g.Map.BarmenArea {
		vector.DrawFilledRect(
			screen,
			float32(BarmenArea[0])*10-float32(g.CameraX),
			float32(BarmenArea[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R: 0, G: 255, B: 0, A: 255},
			false,
		)
	}

	// draw BeerTap
	for _, BeerTap := range g.Map.BeerTaps {
		vector.DrawFilledRect(
			screen,
			float32(BeerTap[0])*10-float32(g.CameraX),
			float32(BeerTap[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R: 255, G: 0, B: 0, A: 255},
			false,
		)
	}
	// draw CounterArea
	for _, CounterArea := range g.Map.CounterArea {
		vector.DrawFilledRect(
			screen,
			float32(CounterArea[0])*10-float32(g.CameraX),
			float32(CounterArea[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R: 255, G: 255, B: 0, A: 255},
			false,
		)
	}

	// draw Exit
	for _, Exit := range g.Map.Exit {
		vector.DrawFilledRect(
			screen,
			float32(Exit[0])*10-float32(g.CameraX),
			float32(Exit[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R: 255, G: 100, B: 100, A: 255},
			false,
		)
	}

	// draw Enter
	for _, Enter := range g.Map.Enter {
		vector.DrawFilledRect(
			screen,
			float32(Enter[0])*10-float32(g.CameraX),
			float32(Enter[1])*10-float32(g.CameraY),
			10,
			10,
			color.RGBA{R: 100, G: 220, B: 220, A: 255},
			false,
		)
	}
}
