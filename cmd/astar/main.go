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

type Simulation struct {
	agentX          float64
	agentY          float64
	path            []*astar.Node
	currentWayPoint int
	walls           [][2]int
}

func (s *Simulation) Update() error {
	if s.path == nil {
		// find route to goal
		m := astar.NewMap(100, 100)
		for _, wall := range s.walls {
			m.SetCell(astar.Position{X: wall[0], Y: wall[1]}, astar.WallCell)
		}
		//m.SetCell(astar.Position{X: int(s.agentX / 7), Y: int(s.agentY / 7)}, astar.EmptyCell)
		start := &astar.Node{Pos: astar.Position{X: int(s.agentX / 7), Y: int(s.agentY / 7)}}
		goal := &astar.Node{Pos: astar.Position{X: 99, Y: 99}}

		searcher := astar.NewJumpPointSearch(m, start, goal)
		path, found := searcher.Search()
		if !found {
			return errors.New("no path found")
		}
		s.path = path
		fmt.Printf("path calculated, len=%d\n", len(s.path))
		s.currentWayPoint = 1
	}

	// move agent towards current waypoint at a speed of 2px per frame
	if s.currentWayPoint < len(s.path) {
		wayPoint := s.path[s.currentWayPoint]
		vx := float64(wayPoint.Pos.X*7) + 3.5 - s.agentX
		vy := float64(wayPoint.Pos.Y*7) + 3.5 - s.agentY
		vNorm := math.Sqrt(vx*vx + vy*vy)

		vx = vx / vNorm
		vy = vy / vNorm

		for _, mur := range s.walls {
			vectx := float64(mur[0])*7 + 3.5 - s.agentX
			vecty := float64(mur[1])*7 + 3.5 - s.agentY
			normeEucli := math.Sqrt((float64(mur[0]*7)+3.5-s.agentX)*(float64(mur[0]*7)+3.5-s.agentX) + (float64(mur[1]*7)+3.5-s.agentY)*(float64(mur[1]*7)+3.5-s.agentY))
			if normeEucli > 50 {
				continue
			}
			vectx = vectx / normeEucli
			vecty = vecty / normeEucli

			reactionMurX := vectx * (3 * math.Exp(-normeEucli/2)) * 10
			reactionMurY := vecty * (3 * math.Exp(-normeEucli/2)) * 10

			vx -= reactionMurX
			vy -= reactionMurY
		}
		//fmt.Println(s.agentX,s.agentY)
		//fmt.Println(s.path[s.currentWayPoint].Pos)
		s.agentX += vx
		s.agentY += vy
		//fmt.Println(totalReactionX, totalReactionY)

		if math.Sqrt((float64(wayPoint.Pos.X*7)+3.5-s.agentX)*(float64(wayPoint.Pos.X*7)+3.5-s.agentX)+(float64(wayPoint.Pos.Y*7)+3.5-s.agentY)*(float64(wayPoint.Pos.Y*7)+3.5-s.agentY)) < 2 {
			s.currentWayPoint++
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

	// draw red circle for goal (99,99)
	ebitenvector.DrawFilledCircle(screen, 99*7, 99*7, 4, colornames.Red, false)

	// draw agent
	ebitenvector.DrawFilledCircle(screen, float32(s.agentX), float32(s.agentY), 4, colornames.Blue, false)
	// draw lines between waypoints
	for i := 0; i < len(s.path)-1; i++ {
		ebitenvector.StrokeLine(screen, float32(s.path[i].Pos.X*7)+3.5, float32(s.path[i].Pos.Y*7)+3.5, float32(s.path[i+1].Pos.X*7)+3.5, float32(s.path[i+1].Pos.Y*7)+3.5, 1, colornames.Green, false)
	}
	for _, mur := range s.walls {
		normeEucli := math.Sqrt((float64(mur[0]*7)+3.5-s.agentX)*(float64(mur[0]*7)+3.5-s.agentX) + (float64(mur[1]*7)+3.5-s.agentY)*(float64(mur[1]*7)+3.5-s.agentY))
		if normeEucli < 50 {
			color := colornames.Blue
			color.A = 50
			//color.R -= uint8(normeEucli / 5)
			//color.G -= uint8(normeEucli / 5)
			//color.B -= uint8(normeEucli / 5)
			ebitenvector.StrokeLine(screen, float32(s.agentX), float32(s.agentY), float32(mur[0])*7+3.5, float32(mur[1])*7+3.5, 1, color, false)
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

	sim := Simulation{
		agentX: 70,
		agentY: 70,
		walls:  testmap.Walls,
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
