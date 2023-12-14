package main

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	goimage "image"
	"image/color"
	"log"
	"math"
)

const (
	ScreenWidth  = 700
	ScreenHeight = 700
)

type Mode int

const (
	ModeMove Mode = iota
	ModeWall
	ModeErase
)

type View struct {
	sim                   *simulation.Simulation
	showPaths             bool
	showWallInteractions  bool
	showAgentInteractions bool
	cameraX, cameraY      int
	draggingPos           [2]int
	draggingCameraPos     [2]int
	CurrentMode           Mode
	cameraZoom            float64
	ui                    *ebitenui.UI
}

var (
	shownAgent          int
	agentAnimations     [7][8][3]*ebiten.Image // character models * n°directions * animation steps
	agentAnimationSteps []float64
	agentLastDirections []int
	mplusNormalFont     font.Face
	rootContainer       *widget.Container
	textarea            *widget.TextArea
	openButton          *widget.Button
	isOpen              bool = true
	FullBeerImg         *ebiten.Image
	EmptyBeerImg        *ebiten.Image
	OneOfFiveBeerImg    *ebiten.Image
	TwoOfFiveBeerImg    *ebiten.Image
	ThreeOfFiveBeerImg  *ebiten.Image
	FourOfFiveBeerImg   *ebiten.Image

	WomanToiletTexture *ebiten.Image

	spriteSheet      *ebiten.Image
	groundImg        *ebiten.Image
	leftGroundImg    *ebiten.Image
	topGroundImg     *ebiten.Image
	topLeftGroundImg *ebiten.Image
	cornerGroundImg  *ebiten.Image

	wallImg            *ebiten.Image
	wallLeftImg        *ebiten.Image
	wallRightImg       *ebiten.Image
	wallLeftRightImg   *ebiten.Image
	wallTopImg         *ebiten.Image
	wallLateralImg     *ebiten.Image
	wallLateralUpRight *ebiten.Image
	wallLateralUpLeft  *ebiten.Image
	wallTopAboveWall   *ebiten.Image

	testMapDense    [][]uint8
	SimulationImage *ebiten.Image
)

func (v *View) Update() error {
	if v.CurrentMode == ModeMove {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if v.draggingPos == [2]int{-1, -1} {
				v.draggingPos = [2]int{x, y}
				v.draggingCameraPos = [2]int{int(v.cameraX), int(v.cameraY)}
			}
			v.cameraX = v.draggingCameraPos[0] + v.draggingPos[0] - x
			v.cameraY = v.draggingCameraPos[1] + v.draggingPos[1] - y
		} else {
			v.draggingPos = [2]int{-1, -1}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		v.showWallInteractions = !v.showWallInteractions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		v.showAgentInteractions = !v.showAgentInteractions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		v.showPaths = !v.showPaths
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		v.sim.Environment.Agents[shownAgent].Vy = -v.sim.Environment.Agents[shownAgent].Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		v.sim.Environment.Agents[shownAgent].Vy = v.sim.Environment.Agents[shownAgent].Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		v.sim.Environment.Agents[shownAgent].Vx = -v.sim.Environment.Agents[shownAgent].Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		v.sim.Environment.Agents[shownAgent].Vx = v.sim.Environment.Agents[shownAgent].Speed
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		maxW := v.sim.Environment.MapSparse.Width
		maxH := v.sim.Environment.MapSparse.Height
		sizeX := float64(ScreenWidth / maxW)
		sizeY := float64(ScreenHeight / maxH)
		mapPosX := (float64(x) + float64(v.cameraX)) / (sizeX * v.cameraZoom)
		mapPosY := (float64(y) + float64(v.cameraY)) / (sizeY * v.cameraZoom)
		// find closest agent
		minDist := math.Inf(1)
		closestAgent := -1
		for i, agent := range v.sim.Environment.Agents {
			dist := math.Sqrt((agent.X-float64(mapPosX))*(agent.X-float64(mapPosX)) + (agent.Y-float64(mapPosY))*(agent.Y-float64(mapPosY)))
			if dist < minDist {
				minDist = dist
				closestAgent = i
			}
		}
		if closestAgent != -1 && minDist < 2 {
			shownAgent = closestAgent
		}
	}

	// scroll to zoom
	_, yWheel := ebiten.Wheel()
	zoomChange := v.cameraZoom * yWheel / 100
	v.cameraZoom += zoomChange
	// also pan when zooming, to keep the same point under the cursor
	x, y := ebiten.CursorPosition()
	v.cameraX += int(float64(x) * zoomChange)
	v.cameraY += int(float64(y) * zoomChange)

	v.ui.Update()

	//agt := v.sim.Environment.Agents[shownAgent]
	// print drink and bladder contents
	//fmt.Printf("Drink: %f, Bladder: %f, Blood alcohol: g/L\n", agt.DrinkContents, agt.BladderContents, agt.BloodAlcoholLevel)

	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	// fill background #404059
	SimulationImage.Fill(color.RGBA{R: 40, G: 40, B: 59, A: 255})

	// write camera pos in top left corner
	textToWrite := fmt.Sprintf("(%d, %d)", int(v.cameraX), int(v.cameraY))
	//textWidth := font.MeasureString(mplusNormalFont, textToWrite).Round()
	textHeight := mplusNormalFont.Metrics().Height.Round()
	text.Draw(
		SimulationImage,
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

	maxW := v.sim.Environment.MapSparse.Width
	maxH := v.sim.Environment.MapSparse.Height
	sizeX := float32(ScreenWidth/maxW) * float32(v.cameraZoom)
	sizeY := float32(ScreenHeight/maxH) * float32(v.cameraZoom)
	// draw ground
	for i := 0; i < maxH; i++ {
		for j := 0; j < maxW; j++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Scale(float64(sizeX)/float64(groundImg.Bounds().Dx()), float64(sizeY)/float64(groundImg.Bounds().Dy()))
			options.GeoM.Translate(float64(j)*float64(sizeX)-float64(v.cameraX), float64(i)*float64(sizeY)-float64(v.cameraY))
			if i > 0 && v.sim.Environment.MapDense[i-1][j] == 1 {
				if j > 0 && v.sim.Environment.MapDense[i][j-1] == 1 {
					SimulationImage.DrawImage(topLeftGroundImg, options)
				} else {
					SimulationImage.DrawImage(topGroundImg, options)
				}
			} else if j > 0 && v.sim.Environment.MapDense[i][j-1] == 1 {
				SimulationImage.DrawImage(leftGroundImg, options)
			} else if i > 0 && j > 0 && v.sim.Environment.MapDense[i-1][j-1] == 1 {
				SimulationImage.DrawImage(cornerGroundImg, options)
			} else {
				SimulationImage.DrawImage(groundImg, options)
			}
		}
	}
	// draw walls
	for _, wall := range v.sim.Environment.MapSparse.Walls {
		//ebitenvector.DrawFilledRect(SimulationImage, float32(wall[0])*sizeX-float32(v.cameraX), float32(wall[1])*sizeY-float32(v.cameraY), sizeX, sizeY, colornames.Black, false)
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(float64(sizeX)/float64(wallImg.Bounds().Dx()), float64(sizeY)/float64(wallImg.Bounds().Dy()))
		options.GeoM.Translate(float64(wall[0])*float64(sizeX)-float64(v.cameraX), float64(wall[1])*float64(sizeY)-float64(v.cameraY))

		rightWallOptions := &ebiten.DrawImageOptions{}
		rightWallOptions.GeoM.Scale(float64(sizeX)/float64(wallImg.Bounds().Dx()), float64(sizeY)/float64(wallImg.Bounds().Dy()))
		rightWallOptions.GeoM.Translate(float64(wall[0])*float64(sizeX)-float64(v.cameraX), float64(wall[1])*float64(sizeY)-float64(v.cameraY))
		rightWallOptions.GeoM.Translate(float64(sizeX)-float64(wallLateralImg.Bounds().Dx())/float64(wallImg.Bounds().Dx())*float64(sizeX), 0)
		// draw full wall if there is no wall down
		if wall[1] < maxH-1 && v.sim.Environment.MapDense[wall[1]+1][wall[0]] == 0 {
			// if nothing left nor right
			if (wall[0] > 0 && v.sim.Environment.MapDense[wall[1]][wall[0]-1] == 0) && (wall[0] < maxW-1 && v.sim.Environment.MapDense[wall[1]][wall[0]+1] == 0) {
				SimulationImage.DrawImage(wallLeftRightImg, options)
			} else if wall[0] > 0 && v.sim.Environment.MapDense[wall[1]][wall[0]-1] == 0 {
				SimulationImage.DrawImage(wallLeftImg, options)
			} else if wall[0] < maxW-1 && v.sim.Environment.MapDense[wall[1]][wall[0]+1] == 0 {
				SimulationImage.DrawImage(wallRightImg, options)
			} else {
				SimulationImage.DrawImage(wallImg, options)
			}
		} else {
			// draw white square #F8F8F8
			ebitenvector.DrawFilledRect(
				SimulationImage,
				float32(wall[0])*sizeX-float32(v.cameraX),
				float32(wall[1])*sizeY-float32(v.cameraY),
				sizeX,
				sizeY,
				color.RGBA{R: 248, G: 248, B: 248, A: 255},
				false,
			)
			// if nothing left or bottom left
			if (wall[0] > 0 && v.sim.Environment.MapDense[wall[1]][wall[0]-1] == 0) || (wall[0] > 0 && wall[1] < maxH-1 && v.sim.Environment.MapDense[wall[1]+1][wall[0]-1] == 0) {
				SimulationImage.DrawImage(wallLateralImg, options)
			}
			// if nothing right or bottom right
			if wall[0] < maxW-1 && v.sim.Environment.MapDense[wall[1]][wall[0]+1] == 0 || (wall[0] < maxW-1 && wall[1] < maxH-1 && v.sim.Environment.MapDense[wall[1]+1][wall[0]+1] == 0) {
				SimulationImage.DrawImage(wallLateralImg, rightWallOptions)
			}

			// if nothing top
			if wall[1] > 0 && v.sim.Environment.MapDense[wall[1]-1][wall[0]] == 0 {
				SimulationImage.DrawImage(wallTopImg, options)

				// if nothing left
				if wall[0] > 0 && v.sim.Environment.MapDense[wall[1]][wall[0]-1] == 0 {
					SimulationImage.DrawImage(wallLateralUpRight, options)
				}
				// if nothing bottom left
				if wall[0] > 0 && wall[1] < maxH-1 && v.sim.Environment.MapDense[wall[1]+1][wall[0]-1] == 0 {
					//SimulationImage.DrawImage(wallLateralUpRight, options)
				}
				// if nothing right or bottom right
				if wall[0] < maxW-1 && v.sim.Environment.MapDense[wall[1]][wall[0]+1] == 0 {
					SimulationImage.DrawImage(wallLateralUpLeft, rightWallOptions)
				}
				if wall[0] < maxW-1 && wall[1] < maxH-1 && v.sim.Environment.MapDense[wall[1]+1][wall[0]+1] == 0 {
					//SimulationImage.DrawImage(wallLateralUpLeft, rightWallOptions)
				}

				// if wall below and left and right
				if wall[1] < maxH-1 && v.sim.Environment.MapDense[wall[1]+1][wall[0]] == 1 && wall[0] > 0 && v.sim.Environment.MapDense[wall[1]][wall[0]-1] == 1 && wall[0] < maxW-1 && v.sim.Environment.MapDense[wall[1]][wall[0]+1] == 1 {
					SimulationImage.DrawImage(wallTopAboveWall, options)
				}
			}
		}
	}
	// draw bar spots and toilet spots
	for _, Beer := range v.sim.Environment.MapSparse.BarPoints {
		ebitenvector.DrawFilledCircle(SimulationImage, float32(Beer[0])*sizeX+sizeX/2-float32(v.cameraX), float32(Beer[1])*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), color.RGBA{R: 201, G: 201, B: 0, A: 255}, false)
	}
	for _, WomanWC := range v.sim.Environment.MapSparse.WomanToiletPoints {
		optsWoman := &ebiten.DrawImageOptions{}
		optsWoman.GeoM.Scale(float64(v.cameraZoom), float64(v.cameraZoom))
		optsWoman.GeoM.Translate(float64(WomanWC[0])*float64(sizeX)-float64(v.cameraX), float64(WomanWC[1])*float64(sizeY)-float64(v.cameraY))
		SimulationImage.DrawImage(WomanToiletTexture, optsWoman)
		//ebitenvector.DrawFilledCircle(SimulationImage, float32(WomanWC[0])*sizeX+sizeX/2-float32(v.cameraX), float32(WomanWC[1])*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), color.RGBA{R: 255, G: 0, B: 200, A: 255}, false)
	}
	for _, ManWC := range v.sim.Environment.MapSparse.ManToiletPoints {
		ebitenvector.DrawFilledCircle(SimulationImage, float32(ManWC[0])*sizeX+sizeX/2-float32(v.cameraX), float32(ManWC[1])*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), color.RGBA{R: 0, G: 200, B: 255, A: 255}, false)
	}

	// draw agents, their position and their goals
	for i := 0; i < v.sim.NAgents; i++ {
		// draw agent
		//color := colornames.Blue

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-float64(agentAnimations[0][0][0].Bounds().Dx())/2, -float64(agentAnimations[0][0][0].Bounds().Dy())/2)
		opts.GeoM.Scale(float64(sizeX)*1.3/float64(agentAnimations[0][0][0].Bounds().Dx()), float64(sizeY)*1.3/float64(agentAnimations[0][0][0].Bounds().Dy()))
		opts.GeoM.Translate(v.sim.Environment.Agents[i].X*float64(sizeX)+float64(sizeX)/2-float64(v.cameraX), v.sim.Environment.Agents[i].Y*float64(sizeY)+float64(sizeY)/2-float64(v.cameraY))

		speedNorm := simulation.Distance(v.sim.Environment.Agents[i].Vx, v.sim.Environment.Agents[i].Vy, 0, 0)
		model := i % 7
		animationImage := agentAnimations[model][agentLastDirections[i]][0]
		if speedNorm > 0.01 {
			angle := simulation.VectToAngle(v.sim.Environment.Agents[i].Vx, -v.sim.Environment.Agents[i].Vy)
			sector := simulation.AngleTo8DirectionsSector(angle)
			// in the spritesheet, the first direction is down and it rotates clockwise
			sector = (8 - sector) % 8
			sector = (sector + 6) % 8
			agentAnimationSteps[i] += speedNorm
			if agentAnimationSteps[i] > 2 {
				agentAnimationSteps[i] = 0
			}
			animationImage = agentAnimations[model][sector][1+int(agentAnimationSteps[i])%2]
			agentLastDirections[i] = sector
		}

		SimulationImage.DrawImage(
			animationImage,
			opts,
		)

		if i == shownAgent {
			//color = colornames.Red
			textarea.SetText(fmt.Sprintf("verre actuel : %.2f \n\n vessie :%.2f ", v.sim.Environment.Agents[i].DrinkContents, v.sim.Environment.Agents[i].BladderContents))
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(float64(v.cameraZoom), float64(v.cameraZoom))
			opts.GeoM.Translate(v.sim.Environment.Agents[i].X*float64(sizeX)+7-float64(v.cameraX), v.sim.Environment.Agents[i].Y*float64(sizeY)-float64(v.cameraY))
			switch {
			case v.sim.Environment.Agents[i].DrinkContents <= 1:
				SimulationImage.DrawImage(EmptyBeerImg, opts)
			case v.sim.Environment.Agents[i].DrinkContents > 1 && v.sim.Environment.Agents[i].DrinkContents < 66:
				SimulationImage.DrawImage(OneOfFiveBeerImg, opts)
			case v.sim.Environment.Agents[i].DrinkContents >= 66 && v.sim.Environment.Agents[i].DrinkContents < 132:
				SimulationImage.DrawImage(TwoOfFiveBeerImg, opts)
			case v.sim.Environment.Agents[i].DrinkContents >= 132 && v.sim.Environment.Agents[i].DrinkContents < 198:
				SimulationImage.DrawImage(ThreeOfFiveBeerImg, opts)
			case v.sim.Environment.Agents[i].DrinkContents >= 198 && v.sim.Environment.Agents[i].DrinkContents < 264:
				SimulationImage.DrawImage(FourOfFiveBeerImg, opts)
			case v.sim.Environment.Agents[i].DrinkContents >= 264:
				SimulationImage.DrawImage(FullBeerImg, opts)
			}
		}
		// ebitenvector.DrawFilledCircle(SimulationImage, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), sizeX/2, color, false)

		if v.sim.Environment.Agents[i].Path != nil && (v.showPaths || i == shownAgent) {
			// draw red circle for goal (99,99)
			ebitenvector.DrawFilledCircle(SimulationImage, float32(v.sim.Environment.Agents[i].Goal.GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Goal.GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), colornames.Red, false)

			// draw lines between all waypoints
			for j := 0; j < len(v.sim.Environment.Agents[i].Path)-1; j++ {
				ebitenvector.StrokeLine(SimulationImage, float32(v.sim.Environment.Agents[i].Path[j].GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Path[j].GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(v.sim.Environment.Agents[i].Path[j+1].GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Path[j+1].GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), colornames.Green, false)
			}

		}

		// draw line between agent and walls that affect it
		if v.showWallInteractions || i == shownAgent {
			for _, mur := range v.sim.Environment.MapSparse.Walls {
				normeEucli := math.Sqrt((float64(mur[0])-v.sim.Environment.Agents[i].X)*(float64(mur[0])-v.sim.Environment.Agents[i].X) + (float64(mur[1])-v.sim.Environment.Agents[i].Y)*(float64(mur[1])-v.sim.Environment.Agents[i].Y))
				if normeEucli < 5 {
					color := colornames.Blue
					color.A = 50
					ebitenvector.StrokeLine(SimulationImage, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), float32(mur[0])*sizeX+sizeX/2-float32(v.cameraX), float32(mur[1])*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), color, false)
				}
			}
		}

		// draw line between agent and agent that affect it
		if v.showAgentInteractions || i == shownAgent {
			for _, otherAgent := range v.sim.Environment.Agents {
				normeEucli := math.Sqrt((otherAgent.X-v.sim.Environment.Agents[i].X)*(otherAgent.X-v.sim.Environment.Agents[i].X) + (otherAgent.Y-v.sim.Environment.Agents[i].Y)*(otherAgent.Y-v.sim.Environment.Agents[i].Y))
				if normeEucli < 5 {
					color := colornames.Red
					color.A = 50
					ebitenvector.StrokeLine(SimulationImage, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), float32(otherAgent.X)*sizeX+sizeX/2-float32(v.cameraX), float32(otherAgent.Y)*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), color, false)
				}
			}
		}
	}

	//add the simulation in the background
	nineSliceImage := image.NewNineSlice(SimulationImage, [3]int{0, SimulationImage.Bounds().Dx(), 0}, [3]int{0, SimulationImage.Bounds().Dy(), 0})
	rootContainer.BackgroundImage = nineSliceImage

	v.ui.Draw(screen)

}

func (v *View) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func init() {
	FullBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/BeerFull.png")
	EmptyBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/BeerEmpty.png")
	OneOfFiveBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/Beer1Of5.png")
	TwoOfFiveBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/Beer2Of5.png")
	ThreeOfFiveBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/Beer3Of5.png")
	FourOfFiveBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/Beer4Of5.png")

	WomanToiletTexture, _, _ = ebitenutil.NewImageFromFile("assets/WomanToilet.png")

	spriteSheet, _, _ = ebitenutil.NewImageFromFile("assets/spritesheet.png")
	groundImg, _, _ = ebitenutil.NewImageFromFile("assets/ground.png")
	leftGroundImg, _, _ = ebitenutil.NewImageFromFile("assets/ground_left.png")
	topGroundImg, _, _ = ebitenutil.NewImageFromFile("assets/ground_top.png")
	topLeftGroundImg, _, _ = ebitenutil.NewImageFromFile("assets/ground_top_left.png")
	cornerGroundImg, _, _ = ebitenutil.NewImageFromFile("assets/ground_corner.png")
	wallImg, _, _ = ebitenutil.NewImageFromFile("assets/wall_new.png")
	wallLeftImg, _, _ = ebitenutil.NewImageFromFile("assets/wall_left.png")
	wallRightImg, _, _ = ebitenutil.NewImageFromFile("assets/wall_right.png")
	wallLeftRightImg, _, _ = ebitenutil.NewImageFromFile("assets/wall_left_right.png")
	wallTopImg, _, _ = ebitenutil.NewImageFromFile("assets/wall_top.png")
	wallLateralImg, _, _ = ebitenutil.NewImageFromFile("assets/wall_lateral.png")
	wallLateralUpLeft, _, _ = ebitenutil.NewImageFromFile("assets/wall_lateral_up_left.png")
	wallLateralUpRight, _, _ = ebitenutil.NewImageFromFile("assets/wall_lateral_up_right.png")
	wallTopAboveWall, _, _ = ebitenutil.NewImageFromFile("assets/wall_top_above_wall.png")
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Pic")

	SimulationImage = ebiten.NewImage(ScreenWidth, ScreenHeight)

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

	// load map from file
	testmap := _map.Map{}
	err = testmap.LoadFromFile("pic")
	if err != nil {
		fmt.Println(err)
		return
	}
	maxW := 0
	maxH := 0
	for _, wall := range testmap.Walls {
		if wall[0] > maxW {
			maxW = wall[0]
		}
		if wall[1] > maxH {
			maxH = wall[1]
		}
	}
	maxW++
	maxH++
	for i := 0; i < maxH; i++ {
		testMapDense = append(testMapDense, make([]uint8, maxW))
	}
	for _, wall := range testmap.Walls {
		testMapDense[wall[1]][wall[0]] = 1
	}

	nAgents := 60

	// initialize animation steps
	agentAnimationSteps = make([]float64, nAgents)

	// initialize animations
	for k := 0; k < 7; k++ {
		agentAnimations[k] = [8][3]*ebiten.Image{}
		for i := 0; i < 8; i++ {
			agentAnimations[k][i] = [3]*ebiten.Image{}
			for j := 0; j < 3; j++ {
				agentAnimations[k][i][j] = ebiten.NewImage(18, 18)
				agentAnimations[k][i][j].DrawImage(spriteSheet.SubImage(goimage.Rect(18*(3*i+j), k*18, 18*(3*i+j+1), (k+1)*18)).(*ebiten.Image), nil)
			}
		}
	}

	// initialize last directions
	agentLastDirections = make([]int, nAgents)

	env := simulation.NewEnvironment(testmap, testMapDense, nAgents)
	sim := simulation.Simulation{
		Environment: env,
		NAgents:     nAgents,
	}

	// Create the rootContainer with the NineSlice image as the background
	rootContainer = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})), // Set NineSlice image as the background

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)
	// construct a textarea
	textarea = widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionStart,
					MaxWidth:  150,
					MaxHeight: 150,
				}),
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(150, 150),
			),
		),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(mplusNormalFont),
		widget.TextAreaOpts.Text("Beer level : 0\n piss level : 0"),
		//Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(5)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{157, 157, 157, 230}),
				Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 230}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 230}),
					Hover: image.NewNineSliceColor(color.NRGBA{200, 200, 200, 230}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{190, 190, 190, 200}),
					Hover:   image.NewNineSliceColor(color.NRGBA{140, 140, 140, 200}),
					Pressed: image.NewNineSliceColor(color.NRGBA{140, 140, 140, 200}),
				},
			),
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
				rootContainer.AddChild(textarea)
			} else {
				rootContainer.RemoveChild(textarea)
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
	// add the textarea as a child of the container
	rootContainer.AddChild(openButton)
	rootContainer.AddChild(textarea)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	view := View{
		sim:        &sim,
		cameraZoom: 1,
		ui:         &ui,
	}

	sim.Start()
	go env.PerceptRequestsHandler()

	fmt.Println("Starting simulation")
	fmt.Println(" - W: toggle showing wall interactions")
	fmt.Println(" - A: toggle showing agent interactions")
	fmt.Println(" - P: toggle showing paths")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&view); err != nil {
		log.Fatal(err)
	}
}
