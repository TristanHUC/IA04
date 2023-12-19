package main

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"image/color"
)

func buildUi(nBarmen int, nAgents int) ebitenui.UI {

	// Create the rootContainer with the NineSlice image as the background
	rootContainer = widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

	// layout for widget
	layoutWidget = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{157, 157, 157, 230})), // Set NineSlice image as the background
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)
	// layout for image
	layoutImage = widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// container to display beer image
	ImageBeer = widget.NewContainer()

	// container to display beer bladder
	ImageBladder = widget.NewContainer()

	// container to display beer character
	ImageCharacter = widget.NewContainer()

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// construct a textarea
	textarea = widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				// Set the layout data for the textarea
				// including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionStart,
					MaxWidth:  180,
					MaxHeight: 65,
				}),
				// Set the minimum size for the widget
				widget.WidgetOpts.MinSize(180, 65),
			),
		),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(mplusNormalFont),
		widget.TextAreaOpts.Text("Beer level : 0\n piss level : 0"),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(5)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{157, 157, 157, 230}),
				Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 230}),
			}),
		),
	)

	textarea = widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				// Set the layout data for the textarea
				// including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionStart,
					MaxWidth:  180,
					MaxHeight: 80,
				}),
				// Set the minimum size for the widget
				widget.WidgetOpts.MinSize(180, 80),
			),
		),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(mplusNormalFont),
		widget.TextAreaOpts.Text("Beer level : 0\n piss level : 0"),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(5)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{157, 157, 157, 230}),
				Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 230}),
			}),
		),
	)

	buttonImage := simulation.LoadButtonImage()
	openButton = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(" X ", mplusNormalFont, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			isOpen = !isOpen
			if isOpen {
				rootContainer.AddChild(layoutWidget)
				rootContainer.AddChild(slider)
				layoutWidget.AddChild(textarea)
				layoutWidget.AddChild(layoutImage)
				layoutImage.AddChild(ImageBeer)
				layoutImage.AddChild(ImageCharacter)
				layoutImage.AddChild(ImageBladder)
			} else {
				rootContainer.RemoveChild(layoutWidget)
				rootContainer.RemoveChild(slider)
				layoutWidget.RemoveChild(textarea)
				layoutWidget.RemoveChild(layoutImage)
				layoutImage.RemoveChild(ImageBeer)
				layoutImage.RemoveChild(ImageCharacter)
				layoutImage.RemoveChild(ImageBladder)

			}
		}),
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	// construct a slider
	slider = widget.NewSlider(
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Direction(widget.DirectionVertical),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(nBarmen, 600),

		widget.SliderOpts.WidgetOpts(
			// Set the widget's dimensions
			widget.WidgetOpts.MinSize(3, 150),
		),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
		// Set the size of the handle
		widget.SliderOpts.FixedHandleSize(3),
		// Set the offset to display the track
		widget.SliderOpts.TrackOffset(0),
		// Set the size to move the handle
		widget.SliderOpts.PageSizeFunc(func() int {
			return 1
		}),
		// Set the callback to call when the slider value is changed
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			nAgentsWished = 610 - args.Current
		}),
	)
	// Set the current value of the slider
	slider.Current = 610 - nAgents

	rootContainer.AddChild(openButton)
	rootContainer.AddChild(layoutWidget)
	rootContainer.AddChild(slider)

	layoutWidget.AddChild(textarea)
	layoutWidget.AddChild(layoutImage)
	layoutImage.AddChild(ImageBeer)
	layoutImage.AddChild(ImageCharacter)
	layoutImage.AddChild(ImageBladder)

	return ui
}
