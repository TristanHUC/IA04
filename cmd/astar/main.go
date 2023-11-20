package main

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"golang.org/x/image/colornames"
	"log"
	"math"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

type Agent struct {
	x, y, vx, vy, gx, gy, speed, reactivity float64
	tx, ty                                  float64
	controllable                            bool
	path            						[]*astar.Node
	currentWayPoint 						int
	goal 									*astar.Node
	start									*astar.Node			
}

type Simulation struct {
	agents  		[]*Agent
	nAgents 		int
	walls           [][2]int
	m 				*astar.Map
}

func (s *Simulation) Update() error {
	if s.m == nil {
		// find route to goal
		s.m = astar.NewMap(100, 100)
		for _, wall := range s.walls {
			s.m.SetCell(astar.Position{X: wall[0], Y: wall[1]}, astar.WallCell)
		}

		var searcher *astar.JumpPointSearch
		var path []*astar.Node
		var found bool

		for i:=0;i<s.nAgents; i++{
		

			//s.m.SetCell(astar.Position{X: int(s.agentX / 7), Y: int(s.agentY / 7)}, astar.EmptyCell)
			s.agents[i].start = &astar.Node{Pos: astar.Position{X: int(s.agents[i].x / 7), Y: int(s.agents[i].y / 7)}}
			s.agents[i].goal = &astar.Node{Pos: astar.Position{X: 99, Y: 99}}
	
			searcher = astar.NewJumpPointSearch(s.m, s.agents[i].start, s.agents[i].goal)
			path, found = searcher.Search()
			if !found {
				return errors.New("no path found")
			}
			s.agents[i].path = path
			fmt.Printf("path calculated, len=%d\n", len(s.agents[i].path))
			s.agents[i].currentWayPoint = 1
	
		}

	}

	var wayPoint *astar.Node
	
	for i:=0;i<s.nAgents; i++{

		// move agent towards current waypoint at a speed of 2px per frame
		if s.agents[i].currentWayPoint < len(s.agents[i].path) {
			wayPoint = s.agents[i].path[s.agents[i].currentWayPoint]
			s.agents[i].vx = float64(wayPoint.Pos.X*7) + 3.5 - s.agents[i].x
			s.agents[i].vy = float64(wayPoint.Pos.Y*7) + 3.5 - s.agents[i].y
			vNorm := math.Sqrt(s.agents[i].vx*s.agents[i].vx + s.agents[i].vy*s.agents[i].vy)

			s.agents[i].vx = s.agents[i].vx / vNorm
			s.agents[i].vy = s.agents[i].vy / vNorm

			var vectx, vecty, normeEucli, reactionMurX, reactionMurY  float64

			for _, mur := range s.walls {
				vectx = float64(mur[0])*7 + 3.5 - s.agents[i].x
				vecty = float64(mur[1])*7 + 3.5 - s.agents[i].y
				normeEucli = math.Sqrt((float64(mur[0]*7)+3.5-s.agents[i].x)*(float64(mur[0]*7)+3.5-s.agents[i].x) + (float64(mur[1]*7)+3.5-s.agents[i].y)*(float64(mur[1]*7)+3.5-s.agents[i].y))
				if normeEucli > 50 {
					continue
				}
				vectx = vectx / normeEucli
				vecty = vecty / normeEucli

				reactionMurX = vectx * (3 * math.Exp(-normeEucli/2)) * 10
				reactionMurY = vecty * (3 * math.Exp(-normeEucli/2)) * 10

				s.agents[i].vx -= reactionMurX
				s.agents[i].vy -= reactionMurY
			}

			s.agents[i].x += s.agents[i].vx
			s.agents[i].y += s.agents[i].vy

			if math.Sqrt((float64(wayPoint.Pos.X*7)+3.5-s.agents[i].x)*(float64(wayPoint.Pos.X*7)+3.5-s.agents[i].x)+(float64(wayPoint.Pos.Y*7)+3.5-s.agents[i].y)*(float64(wayPoint.Pos.Y*7)+3.5-s.agents[i].y)) < 2 {
				s.agents[i].currentWayPoint++
			}
		}
	}
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	// fill white
	screen.Fill(colornames.White)

	// draw walls (7px thick)
	for _, wall := range s.walls {
		ebitenvector.DrawFilledRect(screen, float32(wall[0]*7), float32(wall[1]*7), 7, 7, colornames.Black, false)
	}

	//draw agents, their position and their goals
	for i:=0;i<s.nAgents; i++{

		// draw red circle for goal (99,99)
		ebitenvector.DrawFilledCircle(screen, float32(s.agents[i].goal.Pos.X*7), float32(s.agents[i].goal.Pos.Y*7), 4, colornames.Red, false)

		// draw agent
		ebitenvector.DrawFilledCircle(screen, float32(s.agents[i].x), float32(s.agents[i].y), 4, colornames.Blue, false)
		
		// draw lines between waypoints
		for j := 0; j < len(s.agents[i].path)-1; j++ {
			ebitenvector.StrokeLine(screen, float32(s.agents[i].path[j].Pos.X*7)+3.5, float32(s.agents[i].path[j].Pos.Y*7)+3.5, float32(s.agents[i].path[j+1].Pos.X*7)+3.5, float32(s.agents[i].path[j+1].Pos.Y*7)+3.5, 1, colornames.Green, false)
		}
		for _, mur := range s.walls {
			normeEucli := math.Sqrt((float64(mur[0]*7)+3.5-s.agents[i].x)*(float64(mur[0]*7)+3.5-s.agents[i].x) + (float64(mur[1]*7)+3.5-s.agents[i].y)*(float64(mur[1]*7)+3.5-s.agents[i].y))
			if normeEucli < 50 {
				color := colornames.Blue
				color.A = 50
				//color.R -= uint8(normeEucli / 5)
				//color.G -= uint8(normeEucli / 5)
				//color.B -= uint8(normeEucli / 5)
				ebitenvector.StrokeLine(screen, float32(s.agents[i].x), float32(s.agents[i].y), float32(mur[0])*7+3.5, float32(mur[1])*7+3.5, 1, color, false)
			}
		}
	}

	
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (s *Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Pic")

	// load map from file
	testmap := _map.Map{}
	err := testmap.LoadFromFile("testmap")
	if err != nil {
		return
	}

	nAgents := 5
	agents := make([]*Agent,0,5)
	for i:=0;i<nAgents; i++{
		agents = append(agents,&Agent{x: float64(70*i), y: float64(70)})
	}

	sim := Simulation{
		agents :		agents,
		nAgents :		nAgents,
		walls:  testmap.Walls,
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
