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
	"reflect"
	"time"
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
	showNames             bool
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
	shownAgent           *simulation.Agent
	agentAnimations      [7][8][3]*ebiten.Image // character models * n°directions * animation steps
	agentAnimationSteps  []float64
	agentLastDirections  []int
	mplusNormalFont      font.Face
	agentNameFont        font.Face
	rootContainer        *widget.Container
	simulationInfoWidget *widget.Container
	agentInfoImages      *widget.Container
	agentNameWidget      *widget.Label
	agentInfoWidget      *widget.Container
	ImageBeer            *widget.Container
	ImageBladder         *widget.Container
	ImageCharacter       *widget.Container
	textarea             *widget.TextArea
	openButton           *widget.Button
	slider               *widget.Slider
	isOpen               bool = true
	EmptyBeerImg         *ebiten.Image
	LogoWC               *ebiten.Image

	alternate int

	WomanToiletTexture *ebiten.Image

	spriteSheet      *ebiten.Image
	groundImg        *ebiten.Image
	leftGroundImg    *ebiten.Image
	topGroundImg     *ebiten.Image
	topLeftGroundImg *ebiten.Image
	cornerGroundImg  *ebiten.Image

	wallImg              *ebiten.Image
	wallLeftImg          *ebiten.Image
	wallRightImg         *ebiten.Image
	wallLeftRightImg     *ebiten.Image
	wallTopImg           *ebiten.Image
	wallLateralImg       *ebiten.Image
	wallLateralUpRight   *ebiten.Image
	wallLateralUpLeft    *ebiten.Image
	wallTopAboveWall     *ebiten.Image
	BarDispenserTexture  *ebiten.Image
	counterCornerUpRight *ebiten.Image
	counterLeftRight     *ebiten.Image
	counterTopDown       *ebiten.Image

	testMapDense    [][]uint8
	SimulationImage *ebiten.Image
	BeerImage       *ebiten.Image
	BladderImage    *ebiten.Image
	CharacterImage  *ebiten.Image

	ActionToName map[simulation.Action]string

	nAgentsWished         int
	lastAgentCreationTime time.Time
	nAgentsMax            int
)

func (v *View) Update() error {
	if v.CurrentMode == ModeMove {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if ((y <= -10) || (y >= 150)) && ((x <= 130) || (x >= 189)) {
				if v.draggingPos == [2]int{-1, -1} {
					v.draggingPos = [2]int{x, y}
					v.draggingCameraPos = [2]int{int(v.cameraX), int(v.cameraY)}
				}
				v.cameraX = v.draggingCameraPos[0] + v.draggingPos[0] - x
				v.cameraY = v.draggingCameraPos[1] + v.draggingPos[1] - y
			}
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
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		v.showNames = !v.showNames
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		v.sim.TogglePause()
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		shownAgent.Vy = -shownAgent.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		shownAgent.Vy = shownAgent.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		shownAgent.Vx = -shownAgent.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		shownAgent.Vx = shownAgent.Speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyG) {
		// speed up simulation
		*v.sim.SimulationSpeed -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		// slow down simulation
		*v.sim.SimulationSpeed += 0.1
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		sizeX := 8.0
		sizeY := 8.0
		mapPosX := (float64(x) + float64(v.cameraX)) / (sizeX * v.cameraZoom)
		mapPosY := (float64(y) + float64(v.cameraY)) / (sizeY * v.cameraZoom)
		// find closest agent
		minDist := math.Inf(1)
		var closestAgent *simulation.Agent
		for _, agent := range v.sim.Environment.Agents {
			dist := math.Sqrt((agent.X-float64(mapPosX))*(agent.X-float64(mapPosX)) + (agent.Y-float64(mapPosY))*(agent.Y-float64(mapPosY)))
			if dist < minDist {
				minDist = dist
				closestAgent = agent
			}
		}
		if closestAgent != nil && minDist < 2 {
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

	v.sim.Environment.Update()
	v.sim.NAgents = len(v.sim.Environment.Agents)

	// if agents wished is bigger than the number of agents, create new agents
	if nAgentsWished > len(v.sim.Environment.Agents) && time.Since(lastAgentCreationTime).Seconds() > 0.1 {
		lastAgentCreationTime = time.Now()
		v.sim.NAgents++
		newAgent := simulation.NewAgent(v.sim.Environment.Agents[len(v.sim.Environment.Agents)-1].ID+1, simulation.ClientBehavior{}, v.sim.Environment.MapDense, &v.sim.Environment.MapSparse, v.sim.Environment.PerceptChannel, true, v.sim.Environment.Counter.BeerCounterChan, v.sim.SimulationSpeed)
		v.sim.Environment.Agents = append(v.sim.Environment.Agents, newAgent)
		go v.sim.Environment.Agents[v.sim.NAgents-1].Run()
	}

	// if agents wished is smaller than the number of agents, remove agents
	nAgentsExiting := 0
	for _, agent := range v.sim.Environment.Agents {
		if agent.Action == simulation.GoToExit {
			nAgentsExiting++
		}
	}
	if nAgentsExiting < len(v.sim.Environment.Agents)-v.sim.NBarmans && nAgentsWished < len(v.sim.Environment.Agents)-nAgentsExiting {
		for _, agent := range v.sim.Environment.Agents {
			if reflect.TypeOf(agent.Behavior) != reflect.TypeOf(simulation.ClientBehavior{}) {
				continue
			}
			if agent.Action != simulation.GoToExit {
				agent.PerceptExitChannel <- simulation.GoToExit
				nAgentsExiting++
			}

			if nAgentsWished >= len(v.sim.Environment.Agents)-nAgentsExiting {
				break
			}
		}
	}

	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	// fill background #404059
	SimulationImage.Fill(color.RGBA{R: 40, G: 40, B: 59, A: 255})

	CharacterImage.Fill(color.NRGBA{157, 157, 157, 230})
	BladderImage.Fill(color.NRGBA{157, 157, 157, 230})
	BeerImage.Fill(color.NRGBA{157, 157, 157, 230})

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
	sizeX := 8 * float32(v.cameraZoom)
	sizeY := 8 * float32(v.cameraZoom)
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
	isCounter := false
	//draw walls
	for _, wall := range v.sim.Environment.MapSparse.Walls {
		isCounter = false
		for _, counter := range v.sim.Environment.MapSparse.CounterArea {
			if wall[0] == counter[0] && wall[1] == counter[1] {
				isCounter = true
			}
		}
		if isCounter {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Scale(float64(sizeX)/float64(counterCornerUpRight.Bounds().Dx()), float64(sizeY)/float64(counterCornerUpRight.Bounds().Dy()))
			options.GeoM.Translate(float64(wall[0])*float64(sizeX)-float64(v.cameraX), float64(wall[1])*float64(sizeY)-float64(v.cameraY))

			//if there is a wall above
			if wall[1] > 0 && v.sim.Environment.MapDense[wall[1]-1][wall[0]] == 1 {
				//if there is a wall on the right
				if wall[0] < maxW-1 && v.sim.Environment.MapDense[wall[1]][wall[0]+1] == 1 {
					SimulationImage.DrawImage(counterCornerUpRight, options)
				} else {
					SimulationImage.DrawImage(counterTopDown, options)
				}
			} else {
				//if there is a wall on the left
				if wall[0] > 0 && v.sim.Environment.MapDense[wall[1]][wall[0]-1] == 1 {
					SimulationImage.DrawImage(counterLeftRight, options)
				} else {
					SimulationImage.DrawImage(counterCornerUpRight, options)
				}
			}

		} else {
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

	for _, barDispenser := range v.sim.Environment.MapSparse.BeerTaps {
		optsBar := &ebiten.DrawImageOptions{}
		optsBar.GeoM.Scale(float64(v.cameraZoom)*0.3, float64(v.cameraZoom)*0.3)
		optsBar.GeoM.Translate(float64(barDispenser[0])*float64(sizeX)-float64(v.cameraX)-0.2, float64(barDispenser[1])*float64(sizeY)-float64(v.cameraY))
		SimulationImage.DrawImage(BarDispenserTexture, optsBar)
	}

	for _, Exit := range v.sim.Environment.MapSparse.Exit {
		ebitenvector.DrawFilledCircle(SimulationImage, float32(Exit[0])*sizeX+sizeX/2-float32(v.cameraX), float32(Exit[1])*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), color.RGBA{R: 255, G: 100, B: 100, A: 255}, false)
	}

	for _, Enter := range v.sim.Environment.MapSparse.Enter {
		ebitenvector.DrawFilledCircle(SimulationImage, float32(Enter[0])*sizeX+sizeX/2-float32(v.cameraX), float32(Enter[1])*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), color.RGBA{R: 100, G: 220, B: 220, A: 255}, false)
	}

	optsImage := &ebiten.DrawImageOptions{}
	optsImage.GeoM.Scale(10, 10)

	// draw agents, their position and their goals
	for i := 0; i < v.sim.NAgents; i++ {
		// draw agent
		//color := colornames.Blue

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-float64(agentAnimations[0][0][0].Bounds().Dx())/2, -float64(agentAnimations[0][0][0].Bounds().Dy())/2)
		opts.GeoM.Scale(float64(sizeX)*1.3/float64(agentAnimations[0][0][0].Bounds().Dx()), float64(sizeY)*1.3/float64(agentAnimations[0][0][0].Bounds().Dy()))
		opts.GeoM.Translate(v.sim.Environment.Agents[i].X*float64(sizeX)+float64(sizeX)/2-float64(v.cameraX), v.sim.Environment.Agents[i].Y*float64(sizeY)+float64(sizeY)/2-float64(v.cameraY))

		speedNorm := simulation.Distance(v.sim.Environment.Agents[i].Vx, v.sim.Environment.Agents[i].Vy, 0, 0)
		model := v.sim.Environment.Agents[i].ID % 7
		animationImage := agentAnimations[model][agentLastDirections[i]][0]
		if speedNorm > 0.01 && !v.sim.Paused {
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

		if v.sim.Environment.Agents[i].ID == shownAgent.ID {

			CharacterImage.DrawImage(
				animationImage,
				optsImage,
			)

			textToWrite = v.sim.Environment.Agents[i].Name
			text.Draw(
				SimulationImage,
				textToWrite,
				mplusNormalFont,
				int((v.sim.Environment.Agents[i].X+1)*float64(sizeX)-float64(v.cameraX)-30),
				int((v.sim.Environment.Agents[i].Y+1)*float64(sizeY)-float64(v.cameraY)-5),
				color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 255,
				},
			)

			//color = colornames.Red
			//textarea.SetText(fmt.Sprintf("Number of agents wanted :%d \n Number of agent currently :%d \n action : %s", v.sim.Environment.Agents[i].DrinkContents, v.sim.Environment.Agents[i].BladderContents, nAgentsWished, v.sim.NAgents, ActionToName[v.sim.Environment.Agents[i].Action]))
			agentNameWidget.Label = v.sim.Environment.Agents[i].Name
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(float64(v.cameraZoom), float64(v.cameraZoom))
			opts.GeoM.Translate((v.sim.Environment.Agents[i].X+1)*float64(sizeX)-float64(v.cameraX), (v.sim.Environment.Agents[i].Y+1)*float64(sizeY)-float64(v.cameraY))

			// draw a yellow rectangle in the beer, depending on the amount of beer
			ebitenvector.DrawFilledRect(SimulationImage, float32(v.sim.Environment.Agents[i].X+1.3)*sizeX-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y+2)*sizeY-float32(v.cameraY), sizeX*0.8, -float32(v.sim.Environment.Agents[i].DrinkContents)/330*sizeY, colornames.Yellow, false)
			// draw a white rectangle on top of the beer
			ebitenvector.DrawFilledRect(SimulationImage, float32(v.sim.Environment.Agents[i].X+1.3)*sizeX-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y+2)*sizeY-float32(v.cameraY)-float32(v.sim.Environment.Agents[i].DrinkContents)/330*sizeY*0.99, sizeX*0.8, -0.15*sizeY, colornames.White, true)
			SimulationImage.DrawImage(EmptyBeerImg, opts)

			// draw a white rectangle on top of the beer of widget
			ebitenvector.DrawFilledRect(BeerImage, float32(30), float32(82)-float32(v.sim.Environment.Agents[i].DrinkContents)/330*82, 50, -9, colornames.White, true)

			// draw a yellow rectangle in the beer of widget, depending on the amount of beer
			ebitenvector.DrawFilledRect(BeerImage, float32(25), float32(80), 55, -float32(v.sim.Environment.Agents[i].DrinkContents)/330*80, colornames.Yellow, false)

			if v.sim.Environment.Agents[i].Action == simulation.Action(2) {
				alternate++
				if alternate == 20 {
					alternate = 0
				}
				if (alternate >= 0) && (alternate < 10) {
					// draw a yellow rectangle in the bladder widget, when need to go to WC
					ebitenvector.DrawFilledRect(BladderImage, 0, 0, 100, 100, colornames.Yellow, false)

				}
			}

			BeerImage.DrawImage(EmptyBeerImg, optsImage)
			BladderImage.DrawImage(LogoWC, nil)
		} else {
			if v.showNames {
				textToWrite = v.sim.Environment.Agents[i].Name
				text.Draw(
					SimulationImage,
					textToWrite,
					mplusNormalFont,
					int((v.sim.Environment.Agents[i].X)*float64(sizeX)-float64(v.cameraX)-25),
					int((v.sim.Environment.Agents[i].Y)*float64(sizeY)-float64(v.cameraY)-5),
					color.RGBA{
						R: 0,
						G: 0,
						B: 0,
						A: 100,
					},
				)
			}

		}
		// ebitenvector.DrawFilledCircle(SimulationImage, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), sizeX/2, color, false)

		if v.sim.Environment.Agents[i].Path != nil && (v.showPaths || v.sim.Environment.Agents[i].ID == shownAgent.ID) {
			// draw red circle for goal (99,99)
			ebitenvector.DrawFilledCircle(SimulationImage, float32(v.sim.Environment.Agents[i].Goal.GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Goal.GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), colornames.Red, false)

			// draw lines between all waypoints
			for j := 0; j < len(v.sim.Environment.Agents[i].Path)-1; j++ {
				ebitenvector.StrokeLine(SimulationImage, float32(v.sim.Environment.Agents[i].Path[j].GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Path[j].GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(v.sim.Environment.Agents[i].Path[j+1].GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Path[j+1].GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), colornames.Green, false)
			}

		}

		// draw line between agent and walls that affect it
		if v.showWallInteractions {
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
		if v.showAgentInteractions {
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

	//image bladder
	nineSliceCharacterImage := image.NewNineSlice(CharacterImage, [3]int{0, CharacterImage.Bounds().Dx(), 0}, [3]int{0, CharacterImage.Bounds().Dy(), 0})
	ImageCharacter.BackgroundImage = nineSliceCharacterImage

	//image bladder
	nineSliceBladderImage := image.NewNineSlice(BladderImage, [3]int{0, BladderImage.Bounds().Dx(), 0}, [3]int{0, BladderImage.Bounds().Dy(), 0})
	ImageBladder.BackgroundImage = nineSliceBladderImage

	//image beer
	nineSliceBeerImage := image.NewNineSlice(BeerImage, [3]int{0, BeerImage.Bounds().Dx(), 0}, [3]int{0, BeerImage.Bounds().Dy(), 0})
	ImageBeer.BackgroundImage = nineSliceBeerImage

	//add the simulation in the background
	nineSliceImage := image.NewNineSlice(SimulationImage, [3]int{0, SimulationImage.Bounds().Dx(), 0}, [3]int{0, SimulationImage.Bounds().Dy(), 0})
	rootContainer.BackgroundImage = nineSliceImage

	v.ui.Draw(screen)

}

func (v *View) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func init() {
	LogoWC, _, _ = ebitenutil.NewImageFromFile("assets/toiletteLogo.png")
	EmptyBeerImg, _, _ = ebitenutil.NewImageFromFile("assets/beerGlass.png")

	WomanToiletTexture, _, _ = ebitenutil.NewImageFromFile("assets/WomanToilet.png")
	BarDispenserTexture, _, _ = ebitenutil.NewImageFromFile("assets/v2dispenser.png")

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

	counterCornerUpRight, _, _ = ebitenutil.NewImageFromFile("assets/counterCornerUpRight.png")
	counterLeftRight, _, _ = ebitenutil.NewImageFromFile("assets/counterLeftRight.png")
	counterTopDown, _, _ = ebitenutil.NewImageFromFile("assets/counterTopDown.png")

}

func main() {

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Pic")

	SimulationImage = ebiten.NewImage(ScreenWidth, ScreenHeight)

	CharacterImage = ebiten.NewImage(200, 200)
	BladderImage = ebiten.NewImage(100, 100)
	BeerImage = ebiten.NewImage(100, 100)

	ActionToName = make(map[simulation.Action]string)
	ActionToName[simulation.Action(0)] = "none"
	ActionToName[simulation.Action(1)] = "GoToRandomSpot"
	ActionToName[simulation.Action(2)] = "GoToToilet"
	ActionToName[simulation.Action(3)] = "GoToBar"
	ActionToName[simulation.Action(4)] = "GoToBeerTap"
	ActionToName[simulation.Action(5)] = "WaitForBeer"
	ActionToName[simulation.Action(6)] = "WaitForClient"
	ActionToName[simulation.Action(7)] = "GoToClient"
	ActionToName[simulation.Action(8)] = "GoToExit"

	nAgentsMax = 1000
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

	agentNameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})

	if err != nil {
		log.Fatal(err)
	}

	// load map from file
	testmap, err := _map.LoadMap("pic")
	if err != nil {
		log.Fatal(err)
	}
	maxW := testmap.Width
	maxH := testmap.Height
	for i := 0; i < maxH; i++ {
		testMapDense = append(testMapDense, make([]uint8, maxW))
	}
	for _, wall := range testmap.Walls {
		testMapDense[wall[1]][wall[0]] = 1
	}

	nAgents := 100
	nAgentsWished = nAgents
	nBarmans := 10

	// initialize animation steps
	agentAnimationSteps = make([]float64, nAgentsMax)

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

	// initialize last directions
	agentLastDirections = make([]int, nAgentsMax)

	SimulationSpeed := float32(1)
	env := simulation.NewEnvironment(testmap, testMapDense, nAgents, nBarmans, &SimulationSpeed)
	shownAgent = env.Agents[0]
	sim := simulation.Simulation{
		Environment:     env,
		NAgents:         nAgents,
		NBarmans:        nBarmans,
		SimulationSpeed: &SimulationSpeed,
	}

	ui := buildUi(nBarmans, nAgents)

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
	fmt.Println(" - E: toggle agents' name")
	fmt.Println(" - G: speed simulation")
	fmt.Println(" - H: slow simulation")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&view); err != nil {
		log.Fatal(err)
	}
}
