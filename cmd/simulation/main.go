package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
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
	sim *simulation.Simulation
}

func (v *View) Update() error {
	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	// fill white
	screen.Fill(colornames.White)

	// draw walls (7px thick)
	for _, wall := range v.sim.Environment.Walls {
		ebitenvector.DrawFilledRect(screen, float32(wall[0]*7), float32(wall[1]*7), 7, 7, colornames.Black, false)
	}

	//draw agents, their position and their goals
	for i := 0; i < v.sim.NAgents; i++ {

		// draw red circle for goal (99,99)
		ebitenvector.DrawFilledCircle(screen, float32(v.sim.Environment.Agents[i].Goal.Pos.X*7), float32(v.sim.Environment.Agents[i].Goal.Pos.Y*7), 4, colornames.Red, false)

		// draw agent
		ebitenvector.DrawFilledCircle(screen, float32(v.sim.Environment.Agents[i].X), float32(v.sim.Environment.Agents[i].Y), 4, colornames.Blue, false)

		// draw lines between future waypoints
		for j := v.sim.Environment.Agents[i].CurrentWayPoint; j < len(v.sim.Environment.Agents[i].Path)-1; j++ {
			ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].Path[j].Pos.X*7)+3.5, float32(v.sim.Environment.Agents[i].Path[j].Pos.Y*7)+3.5, float32(v.sim.Environment.Agents[i].Path[j+1].Pos.X*7)+3.5, float32(v.sim.Environment.Agents[i].Path[j+1].Pos.Y*7)+3.5, 1, colornames.Green, false)
		}
		// draw line between agent's projection upon the line between the last waypoint and the next waypoint and the next waypoint

		var currentWayPoint *astar.Node
		if v.sim.Environment.Agents[i].CurrentWayPoint >= len(v.sim.Environment.Agents[i].Path)-1 {
			currentWayPoint = v.sim.Environment.Agents[i].Goal
		} else {
			currentWayPoint = v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint]
		}
		waypointsVectorX := float64(currentWayPoint.Pos.X-v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X)*7 + 3.5
		waypointsVectorY := float64(currentWayPoint.Pos.Y-v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y)*7 + 3.5
		agentVectorX := v.sim.Environment.Agents[i].X - float64(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X*7) - 3.5
		agentVectorY := v.sim.Environment.Agents[i].Y - float64(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y*7) - 3.5
		ProjectedPoint := (agentVectorX*waypointsVectorX + agentVectorY*waypointsVectorY) / (waypointsVectorX*waypointsVectorX + waypointsVectorY*waypointsVectorY)
		ProjectedX := float64(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.X*7) + 3.5 + (ProjectedPoint * waypointsVectorX)
		ProjectedY := float64(v.sim.Environment.Agents[i].Path[v.sim.Environment.Agents[i].CurrentWayPoint-1].Pos.Y*7) + 3.5 + (ProjectedPoint * waypointsVectorY)
		ebitenvector.StrokeLine(screen, float32(ProjectedX), float32(ProjectedY), float32(currentWayPoint.Pos.X*7)+3.5, float32(currentWayPoint.Pos.Y*7)+3.5, 1, colornames.Green, false)
		// draw line between agent and walls that affect it
		for _, mur := range v.sim.Environment.Walls {
			normeEucli := math.Sqrt((float64(mur[0]*7)+3.5-v.sim.Environment.Agents[i].X)*(float64(mur[0]*7)+3.5-v.sim.Environment.Agents[i].X) + (float64(mur[1]*7)+3.5-v.sim.Environment.Agents[i].Y)*(float64(mur[1]*7)+3.5-v.sim.Environment.Agents[i].Y))
			if normeEucli < 50 {
				color := colornames.Blue
				color.A = 50
				ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X), float32(v.sim.Environment.Agents[i].Y), float32(mur[0])*7+3.5, float32(mur[1])*7+3.5, 1, color, false)
			}
		}

		// draw line between agent and agent that affect it
		for _, otherAgent := range v.sim.Environment.Agents {
			normeEucli := math.Sqrt((otherAgent.X-v.sim.Environment.Agents[i].X)*(otherAgent.X-v.sim.Environment.Agents[i].X) + (otherAgent.Y-v.sim.Environment.Agents[i].Y)*(otherAgent.Y-v.sim.Environment.Agents[i].Y))
			if normeEucli < 30 {
				color := colornames.Red
				color.A = 50
				ebitenvector.StrokeLine(screen, float32(v.sim.Environment.Agents[i].X), float32(v.sim.Environment.Agents[i].Y), float32(otherAgent.X), float32(otherAgent.Y), 1, color, false)
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
	err := testmap.LoadFromFile("testmap")
	if err != nil {
		fmt.Println(err)
		return
	}
	m := astar.NewMap(100, 100)
	for _, wall := range testmap.Walls {
		m.SetCell(astar.Position{X: wall[0], Y: wall[1]}, astar.WallCell)
	}

	nAgents := 10

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

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&view); err != nil {
		log.Fatal(err)
	}
}
