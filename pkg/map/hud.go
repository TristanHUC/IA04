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

	beerIconImg, _, err := ebitenutil.NewImageFromFile("assets/beer.png")
	if err != nil {
		log.Fatal(err)
	}

	manWCIconImg, _, err := ebitenutil.NewImageFromFile("assets/ManWC.png")
	if err != nil {
		log.Fatal(err)
	}

	womanWCIconImg, _, err := ebitenutil.NewImageFromFile("assets/WomanWC.png")
	if err != nil {
		log.Fatal(err)
	}

	BarmenAreaIconImg, _, err := ebitenutil.NewImageFromFile("assets/BarmenArea.png")
	if err != nil {
		log.Fatal(err)
	}

	beerTapIconImg, _, err := ebitenutil.NewImageFromFile("assets/BeerTap.png")
	if err != nil {
		log.Fatal(err)
	}

	counterAreaIconImg, _, err := ebitenutil.NewImageFromFile("assets/CounterArea.png")
	if err != nil {
		log.Fatal(err)
	}

	ExitIconImg, _, err := ebitenutil.NewImageFromFile("assets/ExitLogo.png")
	if err != nil {
		log.Fatal(err)
	}

	EnterIconImg, _, err := ebitenutil.NewImageFromFile("assets/EnterLogo.png")
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

	beerButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 200,
		width:        40,
		height:       40,
		text:         "",
		image:        beerIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         ModeBeer,
	}

	beerButton.imageOptions.GeoM.Scale(0.2, 0.2)
	beerButton.imageOptions.GeoM.Translate(float64(beerButton.x), float64(beerButton.y))

	manWCButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 250,
		width:        40,
		height:       40,
		text:         "",
		image:        manWCIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         ModeManWC,
	}

	manWCButton.imageOptions.GeoM.Scale(0.2, 0.2)
	manWCButton.imageOptions.GeoM.Translate(float64(manWCButton.x), float64(manWCButton.y))

	womanWCButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 300,
		width:        40,
		height:       40,
		text:         "",
		image:        womanWCIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         ModeWomanWC,
	}

	womanWCButton.imageOptions.GeoM.Scale(0.2, 0.2)
	womanWCButton.imageOptions.GeoM.Translate(float64(womanWCButton.x), float64(womanWCButton.y))

	BarmenAreaButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 350,
		width:        40,
		height:       40,
		text:         "",
		image:        BarmenAreaIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         BarmenArea,
	}
	BarmenAreaButton.imageOptions.GeoM.Scale(0.2, 0.2)
	BarmenAreaButton.imageOptions.GeoM.Translate(float64(BarmenAreaButton.x), float64(BarmenAreaButton.y))

	bierTapButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 400,
		width:        40,
		height:       40,
		text:         "",
		image:        beerTapIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         ModeBeerTap,
	}
	bierTapButton.imageOptions.GeoM.Scale(0.2, 0.2)
	bierTapButton.imageOptions.GeoM.Translate(float64(bierTapButton.x), float64(bierTapButton.y))

	counterAreaButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 450,
		width:        40,
		height:       40,
		text:         "",
		image:        counterAreaIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         CounterArea,
	}
	counterAreaButton.imageOptions.GeoM.Scale(0.2, 0.2)
	counterAreaButton.imageOptions.GeoM.Translate(float64(counterAreaButton.x), float64(counterAreaButton.y))

	ExitButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 500,
		width:        40,
		height:       40,
		text:         "",
		image:        ExitIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         Exit,
	}
	ExitButton.imageOptions.GeoM.Scale(0.2, 0.2)
	ExitButton.imageOptions.GeoM.Translate(float64(ExitButton.x), float64(ExitButton.y))

	EnterButton := Button{
		x:            screenWidth - 50,
		y:            screenHeight - 550,
		width:        40,
		height:       40,
		text:         "",
		image:        EnterIconImg,
		imageOptions: &ebiten.DrawImageOptions{},
		mode:         Enter,
	}
	EnterButton.imageOptions.GeoM.Scale(0.2, 0.2)
	EnterButton.imageOptions.GeoM.Translate(float64(EnterButton.x), float64(EnterButton.y))

	buttons = append(buttons, moveCursorButton, wallButton, eraseButton, beerButton, manWCButton, womanWCButton, BarmenAreaButton, bierTapButton, counterAreaButton, ExitButton, EnterButton)
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
