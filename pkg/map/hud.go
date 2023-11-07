package _map

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
)

var buttons []Button

var (
	moveCursorImg *ebiten.Image
	wallIconImg   *ebiten.Image
)

var mplusNormalFont font.Face

func InitHud(screenWidth, screenHeight int) {
	// load cursors
	moveCursorImg, _, err := ebitenutil.NewImageFromFile("assets/move.png")
	if err != nil {
		log.Fatal(err)
	}

	wallIconImg, _, err = ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		log.Fatal(err)
	}

	eraseIconImg, _, err := ebitenutil.NewImageFromFile("assets/erase.png")
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetCursorShape(ebiten.CursorShapeMove)

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	// create buttons
	moveCursorButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 50,
		width:        40,
		height:       40,
		text:         "",
		image:        moveCursorImg,
		imageOptions: &ebiten.DrawImageOptions{},
		selected:     true,
		mode:         ModeMove,
	}
	moveCursorButton.imageOptions.GeoM.Scale(0.2, 0.2)
	moveCursorButton.imageOptions.GeoM.Translate(float64(moveCursorButton.x), float64(moveCursorButton.y))

	wallButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 100,
		width:        40,
		height:       40,
		text:         "",
		image:        wallIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         ModeWall,
	}
	wallButton.imageOptions.GeoM.Scale(0.2, 0.2)
	wallButton.imageOptions.GeoM.Translate(float64(wallButton.x), float64(wallButton.y))

	eraseButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 150,
		width:        40,
		height:       40,
		text:         "",
		image:        eraseIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         ModeErase,
	}

	eraseButton.imageOptions.GeoM.Scale(0.2, 0.2)
	eraseButton.imageOptions.GeoM.Translate(float64(eraseButton.x), float64(eraseButton.y))

	buttons = append(buttons, moveCursorButton, wallButton, eraseButton)
}

func DrawButtons(screen *ebiten.Image) {
	for _, button := range buttons {
		// draw external box
		margin := 0
		if button.selected {
			margin = 4
		}
		vector.DrawFilledRect(
			screen,
			float32(button.x),
			float32(button.y),
			float32(button.width),
			float32(button.height),
			color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 255,
			},
			false,
		)
		// draw internal box
		vector.DrawFilledRect(
			screen,
			float32(button.x+1)+float32(margin)/2,
			float32(button.y+1)+float32(margin)/2,
			float32(button.width-2)-float32(margin),
			float32(button.height-2)-float32(margin),
			color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			},
			false,
		)
		// draw text if there is any
		if button.text != "" {
			text.Draw(
				screen,
				button.text,
				mplusNormalFont,
				button.x+5,
				button.y+15,
				color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 255,
				},
			)
		} else {
			screen.DrawImage(button.image, button.imageOptions)
		}
	}
}

func (g *Game) DrawHud(screen *ebiten.Image) {
	// write camera pos in top left corner
	textToWrite := fmt.Sprintf("(%d, %d)", int(g.CameraX), int(g.CameraY))
	//textWidth := font.MeasureString(mplusNormalFont, textToWrite).Round()
	textHeight := mplusNormalFont.Metrics().Height.Round()
	text.Draw(
		screen,
		textToWrite,
		mplusNormalFont,
		0,
		textHeight,
		color.RGBA{
			R: 50,
			G: 100,
			B: 50,
			A: 255,
		},
	)
	DrawButtons(screen)
}

func UpdateHud(g *Game) bool {
	stopPropagation := false
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for i, button := range buttons {
			if x >= button.x && x <= button.x+button.width && y >= button.y && y <= button.y+button.height {
				buttons[i].selected = true
				g.CurrentMode = button.mode
				for j, _ := range buttons {
					if j != i {
						buttons[j].selected = false
					}
				}
			}
		}
	}

	// set cursor to point icon if hovering button
	x, y := ebiten.CursorPosition()
	onButton := false
	for _, button := range buttons {
		if x >= button.x && x <= button.x+button.width && y >= button.y && y <= button.y+button.height {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
			onButton = true
			stopPropagation = true
			break
		}
	}
	if !onButton {
		ebiten.SetCursorShape(modeToCursor[g.CurrentMode])
	}

	return stopPropagation
}
