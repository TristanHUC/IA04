package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	face, _ := loadFont(14)
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					MaxWidth:  300,
					MaxHeight: 100,
				}),
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(300, 100),
			),
		),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.Black),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(face),

		widget.TextAreaOpts.Text("Beer level : 0\n piss level : 0"),
		//Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(10)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{127, 127, 127, 255}),
				Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Hover: image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{90, 90, 90, 255}),
					Hover:   image.NewNineSliceColor(color.NRGBA{90, 90, 90, 255}),
					Pressed: image.NewNineSliceColor(color.NRGBA{80, 80, 80, 255}),
				},
			),
		),
	)
	// add the textarea as a child of the container
	rootContainer.AddChild(textarea)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(700, 700)
	ebiten.SetWindowTitle("Ebiten UI - TextArea")

	game := game{
		ui: &ui,
	}

	// run Ebiten main loop
	err := ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
	}
}

// Layout implements Game.
func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update implements Game.
func (g *game) Update() error {
	// update the UI
	g.ui.Update()
	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}