package _map

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
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
	}
	wallButton.imageOptions.GeoM.Scale(0.2, 0.2)
	wallButton.imageOptions.GeoM.Translate(float64(wallButton.x), float64(wallButton.y))

	buttons = append(buttons, moveCursorButton, wallButton)
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
