package main

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"image/color"
)

var ActionToName = map[simulation.Action]string{
	simulation.None:                  "None",
	simulation.GoToRandomSpot:        "GoToRandomSpot",
	simulation.GoToToilet:            "GoToToilet",
	simulation.GoToBar:               "GoToBar",
	simulation.GoToBeerTap:           "GoToBeerTap",
	simulation.WaitForBeer:           "WaitForBeer",
	simulation.WaitForClient:         "WaitForClient",
	simulation.GoToClient:            "GoToClient",
	simulation.GoToExit:              "GoToExit",
	simulation.GoWithFriends:         "GoWithFriends",
	simulation.GoFarFromBarAndToilet: "GoFarFromBarAndToilet",
	simulation.WaitingWithFriends:    "WaitingWithFriends",
}

func buildUi(nBarmen int, nAgents int, simSpeedChangeCallback func(int)) ebitenui.UI {

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
			widget.RowLayoutOpts.Spacing(250),
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
				Position: widget.RowLayoutPositionStart,
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

	nAgentsText = widget.NewLabel(
		widget.LabelOpts.Text(fmt.Sprintf("%d/%d Agents", nAgents, nAgentsWished), mplusNormalFont, &widget.LabelColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
	)
	simulationInfoWidget.AddChild(nAgentsText)

	// construct a slider for number of agents
	slider = widget.NewSlider(
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(nBarmen, 600),

		widget.SliderOpts.WidgetOpts(
			// Set the widget's dimensions
			widget.WidgetOpts.MinSize(150, 3),
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
			nAgentsWished = args.Current
		}),
	)
	// Set the current value of the slider
	slider.Current = nAgents
	simulationInfoWidget.AddChild(slider)

	// label for sim speed
	simSpeedText = widget.NewLabel(
		widget.LabelOpts.Text(fmt.Sprintf("Simulation speed: %dx", 1), mplusNormalFont, &widget.LabelColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
	)
	simulationInfoWidget.AddChild(simSpeedText)

	// construct a slider for simulation speed
	slider = widget.NewSlider(
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(nBarmen, 1500),

		widget.SliderOpts.WidgetOpts(
			// Set the widget's dimensions
			widget.WidgetOpts.MinSize(150, 3),
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
			simSpeedChangeCallback(args.Current)
		}),
	)
	// Set the current value of the slider
	slider.Current = 100
	simulationInfoWidget.AddChild(slider)

	return ui
}
