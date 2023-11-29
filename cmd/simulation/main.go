package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"golang.org/x/image/colornames"
	"log"
	"math"
)

const (
	ScreenWidth  = 700
	ScreenHeight = 700
)

type View struct {
	sim                   *simulation.Simulation
	showPaths             bool
	showWallInteractions  bool
	showAgentInteractions bool
}

var shownAgent int

func (v *View) Update() error {
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
		maxW := v.sim.Environment.Map.Width
		maxH := v.sim.Environment.Map.Height
		sizeX := float64(ScreenWidth / maxW)
		sizeY := float64(ScreenHeight / maxH)
		mapPosX := float64(x) / sizeX
		mapPosY := float64(y) / sizeY
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
	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	// fill white
	screen.Fill(colornames.White)

	// draw walls (7px thick)
	maxW := v.sim.Environment.Map.Width
	maxH := v.sim.Environment.Map.Height
	sizeX := float32(ScreenWidth / maxW)
	sizeY := float32(ScreenHeight / maxH)
	for _, wall := range v.sim.Environment.Walls {
		ebitenvector.DrawFilledRect(screen, float32(wall[0])*sizeX, float32(wall[1])*sizeY, sizeX, sizeY, colornames.Black, false)
	}

	// draw agents, their position and their goals
	for i := 0; i < v.sim.NAgents; i++ {
		// draw agent
		color := colornames.Blue
		if i == shownAgent {
			color = colornames.Red
		}
		ebitenvector.DrawFilledCircle(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2, sizeX/2, color, false)

		if v.sim.Environment.Agents[i].Path != nil && (v.showPaths || i == shownAgent) {
			// draw red circle for goal (99,99)
			ebitenvector.DrawFilledCircle(screen, float32(v.sim.Environment.Agents[i].Goal.Pos.X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Goal.Pos.Y)*sizeY+sizeY/2, 4, colornames.Red, false)

			// draw lines between future waypoints
			for j := v.sim.Environment.Agents[i].CurrentWayPoint; j < len(v.sim.Environment.Agents[i].Path)-1; j++ {
				ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].Path[j].Pos.X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Path[j].Pos.Y)*sizeY+sizeY/2, float32(v.sim.Environment.Agents[i].Path[j+1].Pos.X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Path[j+1].Pos.Y)*sizeY+sizeY/2, 1, colornames.Green, false)
			}

			// draw line between agent's projection upon the line between the last waypoint and the next waypoint and the next waypoint
			var currentWayPoint *astar.Node
			if v.sim.Environment.Agents[i].CurrentWayPoint >= len(v.sim.Environment.Agents[i].Path)-1 {
				currentWayPoint = v.sim.Environment.Agents[i].Goal
			} else {
				currentWayPoint = v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint]
			}
			//waypointsVectorX := float32(currentWayPoint.Pos.X-v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*sizeX + sizeX/2
			//waypointsVectorY := float32(currentWayPoint.Pos.Y-v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*sizeY + sizeY/2
			//agentVectorX := float32(v.sim.Environment.Agents[i].X) - float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*sizeX - sizeX/2
			//agentVectorY := float32(v.sim.Environment.Agents[i].Y) - float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*sizeY - sizeY/2
			//ProjectedPoint := (agentVectorX*waypointsVectorX + agentVectorY*waypointsVectorY) / (waypointsVectorX*waypointsVectorX + waypointsVectorY*waypointsVectorY)
			//ProjectedX := float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*sizeX + sizeX/2 + (ProjectedPoint * waypointsVectorX)
			//ProjectedY := float32(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*sizeY + sizeY/2 + (ProjectedPoint * waypointsVectorY)
			//ebitenvector.StrokeLine(screen, ProjectedX, ProjectedY, float32(currentWayPoint.Pos.X)*sizeX+sizeX/2, float32(currentWayPoint.Pos.Y)*sizeY+sizeY/2, 1, colornames.Green, false)

			// draw line between agent and next waypoint
			ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2, float32(currentWayPoint.Pos.X)*sizeX+sizeX/2, float32(currentWayPoint.Pos.Y)*sizeY+sizeY/2, 1, colornames.Green, false)
		}

		// draw line between agent and walls that affect it
		if v.showWallInteractions || i == shownAgent {
			for _, mur := range v.sim.Environment.Walls {
				normeEucli := math.Sqrt((float64(mur[0])-v.sim.Environment.Agents[i].X)*(float64(mur[0])-v.sim.Environment.Agents[i].X) + (float64(mur[1])-v.sim.Environment.Agents[i].Y)*(float64(mur[1])-v.sim.Environment.Agents[i].Y))
				if normeEucli < 5 {
					color := colornames.Blue
					color.A = 50
					ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2, float32(mur[0])*sizeX+sizeX/2, float32(mur[1])*sizeY+sizeY/2, 1, color, false)
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
					ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X)*sizeX+sizeX/2, float32(v.sim.Environment.Agents[i].Y)*sizeY+sizeY/2, float32(otherAgent.X)*sizeX+sizeX/2, float32(otherAgent.Y)*sizeY+sizeY/2, 1, color, false)
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
	// load map from file
	testmap := _map.Map{}
	err := testmap.LoadFromFile("pic")
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
	m := astar.NewMap(maxW, maxH)
	for _, wall := range testmap.Walls {
		m.SetCell(astar.Position{X: wall[0], Y: wall[1]}, astar.WallCell)
	}

	nAgents := 500

	env := simulation.NewEnvironment(testmap.Walls, m, nAgents)
	sim := simulation.Simulation{
		Environment: env,
		NAgents:     nAgents,
	}

	view := View{
		sim: &sim,
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
