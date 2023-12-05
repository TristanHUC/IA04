package simulation

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
)

func LoadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.RGBA{R: 170, G: 170, B: 180, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 130, G: 130, B: 150, A: 255})
	pressed := image.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}
