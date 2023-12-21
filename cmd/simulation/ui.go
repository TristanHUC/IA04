package main

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"image/color"
)

var ActionToName = map[simulation.Action]string{
	simulation.None:           "None",
	simulation.GoToRandomSpot: "GoToRandomSpot",
	simulation.GoToToilet:     "GoToToilet",
	simulation.GoToBar:        "GoToBar",
	simulation.GoToBeerTap:    "GoToBeerTap",
	simulation.WaitForBeer:    "WaitForBeer",
	simulation.WaitForClient:  "WaitForClient",
	simulation.GoToClient:     "GoToClient",
	simulation.GoToExit:       "GoToExit",
	simulation.GoWithFriends:  "GoWithFriends",
}

func buildUi(nBarmen int, nAgents int) ebitenui.UI {

	// Create the rootContainer with the NineSlice image as the background
	rootContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(10)),
		)),
	)

	// widget that is only a container for the simulation info and agent info
	rootLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	rootContainer.AddChild(rootLayout)

	// widget for agent info
	agentInfoWidget = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{157, 157, 157, 230})), // Set NineSlice image as the background
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
		),
	)
	rootLayout.AddChild(agentInfoWidget)

	// widget for simulation info
	simulationInfoWidget = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{157, 157, 157, 230})), // Set NineSlice image as the background
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
		)),
	)
	rootLayout.AddChild(simulationInfoWidget)

	agentNameWidget = widget.NewLabel(
		widget.LabelOpts.Text("John Mayer", agentNameFont, &widget.LabelColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
	)
	agentInfoWidget.AddChild(agentNameWidget)

	agentActionLabel = widget.NewLabel(
		widget.LabelOpts.Text("None", mplusNormalFont, &widget.LabelColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
	)
	agentInfoWidget.AddChild(agentActionLabel)

	// layout for agent info images
	agentInfoImages = widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	agentInfoWidget.AddChild(agentInfoImages)

	// container to display beer image
	ImageBeer = widget.NewContainer()
	agentInfoImages.AddChild(ImageBeer)

	// container to display beer bladder
	// beer bladder ??? really titi ???
	ImageBladder = widget.NewContainer()
	agentInfoImages.AddChild(ImageBladder)

	// container to display beer character
	// beer character ??? really titi ???
	ImageCharacter = widget.NewContainer()
	agentInfoImages.AddChild(ImageCharacter)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	//buttonImage := simulation.LoadButtonImage()
	//openButton = widget.NewButton(
	//	widget.ButtonOpts.Image(buttonImage),
	//	widget.ButtonOpts.Text(" X ", mplusNormalFont, &widget.ButtonTextColor{
	//		Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
	//	}),
	//	widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
	//		isOpen = !isOpen
	//		if isOpen {
	//			rootContainer.AddChild(agentInfoWidget)
	//			rootContainer.AddChild(slider)
	//			agentInfoWidget.AddChild(textarea)
	//			agentInfoWidget.AddChild(layoutImage)
	//			layoutImage.AddChild(ImageBeer)
	//			layoutImage.AddChild(ImageCharacter)
	//			layoutImage.AddChild(ImageBladder)
	//		} else {
	//			rootContainer.RemoveChild(agentInfoWidget)
	//			rootContainer.RemoveChild(slider)
	//			agentInfoWidget.RemoveChild(textarea)
	//			agentInfoWidget.RemoveChild(layoutImage)
	//			layoutImage.RemoveChild(ImageBeer)
	//			layoutImage.RemoveChild(ImageCharacter)
	//			layoutImage.RemoveChild(ImageBladder)
	//
	//		}
	//	}),
	//	widget.ButtonOpts.WidgetOpts(
	//		instruct the container's anchor layout to center the button both horizontally and vertically
	//widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
	//	HorizontalPosition: widget.AnchorLayoutPositionStart,
	//	VerticalPosition:   widget.AnchorLayoutPositionStart,
	//}),
	//),
	//)
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
	simulationInfoWidget.AddChild(slider)
	// Set the current value of the slider
	slider.Current = 610 - nAgents

	//rootContainer.AddChild(openButton)
	//rootContainer.AddChild(agentInfoWidget)
	//rootContainer.AddChild(simulationInfoWidget)

	//agentInfoWidget.AddChild(textarea)
	//agentInfoWidget.AddChild(agentInfoImages)
	//simulationInfoWidget.AddChild(slider)

	//agentInfoImages.AddChild(ImageBeer)
	//agentInfoImages.AddChild(ImageCharacter)
	//agentInfoImages.AddChild(ImageBladder)

	return ui
}
