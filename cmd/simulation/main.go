package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
}

var shownAgent int
var mplusNormalFont font.Face

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
		if closestAgent != -1 {
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

	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	// fill white
	screen.Fill(colornames.White)

	// write camera pos in top left corner
	textToWrite := fmt.Sprintf("(%d, %d)", int(v.cameraX), int(v.cameraY))
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

	maxW := v.sim.Environment.MapSparse.Width
	maxH := v.sim.Environment.MapSparse.Height
	sizeX := float32(ScreenWidth/maxW) * float32(v.cameraZoom)
	sizeY := float32(ScreenHeight/maxH) * float32(v.cameraZoom)
	for _, wall := range v.sim.Environment.MapSparse.Walls {
		ebitenvector.DrawFilledRect(screen, float32(wall[0])*sizeX-float32(v.cameraX), float32(wall[1])*sizeY-float32(v.cameraY), sizeX, sizeY, colornames.Black, false)
	}

	// draw agents, their position and their goals
	for i := 0; i < v.sim.NAgents; i++ {
		// draw agent
		color := colornames.Blue
		if i == shownAgent {
			color = colornames.Red
		}
		ebitenvector.DrawFilledCircle(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), sizeX/2, color, false)

		if v.sim.Environment.Agents[i].Path != nil && (v.showPaths || i == shownAgent) {
			// draw red circle for goal (99,99)
			ebitenvector.DrawFilledCircle(screen, float32(v.sim.Environment.Agents[i].Goal.GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Goal.GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(4*v.cameraZoom), colornames.Red, false)

			// draw lines between all waypoints
			for j := 0; j < len(v.sim.Environment.Agents[i].Path)-1; j++ {
				ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].Path[j].GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Path[j].GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(v.sim.Environment.Agents[i].Path[j+1].GetCol())*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Path[j+1].GetRow())*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), colornames.Green, false)
			}

			// draw line between agent's projection upon the line between the last waypoint and the next waypoint and the next waypoint
			//var currentWayPoint *jps.Node
			//if v.sim.Environment.Agents[i].CurrentWayPoint >= len(v.sim.Environment.Agents[i].Path)-1 {
			//	currentWayPoint = v.sim.Environment.Agents[i].Goal
			//} else {
			//	currentWayPoint = &v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint]
			//}
			//waypointsVectorX := float32(currentWayPoint.Pos.X-v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*sizeX + sizeX/2
			//waypointsVectorY := float32(currentWayPoint.Pos.Y-v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*sizeY + sizeY/2
			//agentVectorX := float32(v.sim.Environment.Agents[i].X) - float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*sizeX - sizeX/2
			//agentVectorY := float32(v.sim.Environment.Agents[i].Y) - float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*sizeY - sizeY/2
			//ProjectedPoint := (agentVectorX*waypointsVectorX + agentVectorY*waypointsVectorY) / (waypointsVectorX*waypointsVectorX + waypointsVectorY*waypointsVectorY)
			//ProjectedX := float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*sizeX + sizeX/2 + (ProjectedPoint * waypointsVectorX)
			//ProjectedY := float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*sizeY + sizeY/2 + (ProjectedPoint * waypointsVectorY)
			//ebitenvector.StrokeLine(screen, ProjectedX, ProjectedY, float32(currentWayPoint.Pos.X)*sizeX+sizeX/2, float32(currentWayPoint.Pos.Y)*sizeY+sizeY/2, 1, colornames.Green, false)

			// draw line between agent and next waypoint
			//ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2, float32(currentWayPoint.GetCol())*sizeX+sizeX/2, float32(currentWayPoint.GetRow())*sizeY+sizeY/2, 1, colornames.Green, false)
		}

		// draw line between agent and walls that affect it
		if v.showWallInteractions || i == shownAgent {
			for _, mur := range v.sim.Environment.MapSparse.Walls {
				normeEucli := math.Sqrt((float64(mur[0])-v.sim.Environment.Agents[i].X)*(float64(mur[0])-v.sim.Environment.Agents[i].X) + (float64(mur[1])-v.sim.Environment.Agents[i].Y)*(float64(mur[1])-v.sim.Environment.Agents[i].Y))
				if normeEucli < 5 {
					color := colornames.Blue
					color.A = 50
					ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), float32(mur[0])*sizeX+sizeX/2-float32(v.cameraX), float32(mur[1])*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), color, false)
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
					ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2-float32(v.cameraX), float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2-float32(v.cameraY), float32(otherAgent.X)*sizeX+sizeX/2-float32(v.cameraX), float32(otherAgent.Y)*sizeY+sizeY/2-float32(v.cameraY), float32(1*v.cameraZoom), color, false)
				}
			}
		}
	}
}

func (v *View) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Pic")

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
	var testMapDense [][]uint8
	for i := 0; i < maxH; i++ {
		testMapDense = append(testMapDense, make([]uint8, maxW))
	}
	for _, wall := range testmap.Walls {
		testMapDense[wall[1]][wall[0]] = 1
	}

	nAgents := 300

	env := simulation.NewEnvironment(testmap, testMapDense, nAgents)
	sim := simulation.Simulation{
		Environment: env,
		NAgents:     nAgents,
	}

	view := View{
		sim:        &sim,
		cameraZoom: 1,
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
